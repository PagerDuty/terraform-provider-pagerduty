package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var pdClient *Client
var cacheType string
var cacheMongoURL string
var cacheMaxAge, _ = time.ParseDuration("10s")

var mongoClient *mongo.Client

var mongoCache map[string]*mongo.Collection

var memoryCache = map[string]*sync.Map{
	"users":              {},
	"contact_methods":    {},
	"notification_rules": {},
	"misc":               {},
}

type cacheAbilitiesRecord struct {
	ID        string
	Abilities *ListAbilitiesResponse
}

type cacheLastRefreshRecord struct {
	ID        string
	Users     time.Time
	Abilities time.Time
}

// InitCache initializes the cache according to the setting in TF_PAGERDUTY_CACHE
func InitCache(c *Client) {
	pdClient = c
	if cacheMongoURL = os.Getenv("TF_PAGERDUTY_CACHE"); strings.HasPrefix(cacheMongoURL, "mongodb://") {
		log.Printf("===== Enabling PagerDuty Mongo cache at %v", cacheMongoURL)
		cacheType = "mongo"
	} else if cacheMongoURL == "memory" {
		log.Println("===== Enabling PagerDuty memory cache =====")
		cacheType = "memory"
		return
	} else {
		log.Println("===== PagerDuty Cache Skipping Init =====")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, _ = mongo.Connect(ctx, options.Client().ApplyURI(cacheMongoURL))

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("===== PagerDuty Cache couldn't connect to MongoDB at %q, disabling cache =====", cacheMongoURL)
		cacheType = ""
		return
	}

	if os.Getenv("TF_PAGERDUTY_CACHE_MAX_AGE") != "" {
		d, err := time.ParseDuration(os.Getenv("TF_PAGERDUTY_CACHE_MAX_AGE"))
		if err != nil {
			log.Printf("===== PagerDuty Cache couldn't parse max age %q, using the default %v =====", os.Getenv("TF_PAGERDUTY_CACHE_MAX_AGE"), cacheMaxAge)
		} else {
			cacheMaxAge = d
		}
	}

	mongoCache = map[string]*mongo.Collection{
		"users":              mongoClient.Database("pagerduty").Collection("users"),
		"contact_methods":    mongoClient.Database("pagerduty").Collection("contact_methods"),
		"notification_rules": mongoClient.Database("pagerduty").Collection("notification_rules"),
		"misc":               mongoClient.Database("pagerduty").Collection("misc"),
	}
}

// PopulateMemoryCache does initial population of the cache if memory caching is selected
func PopulateMemoryCache() {
	if _, present := os.LookupEnv("TF_PAGERDUTY_CACHE_PREFILL"); present {
		log.Println("===== Prefilling memory cache =====")
		abilities, _, _ := pdClient.Abilities.List()

		abilitiesRecord := &cacheAbilitiesRecord{
			ID:        "abilities",
			Abilities: abilities,
		}
		cachePut("misc", "abilities", abilitiesRecord)
		
		var pdo = ListUsersOptions{
			Include: []string{"contact_methods", "notification_rules"},
			Limit:   100,
		}

		fullUsers, err := pdClient.Users.ListAll(&pdo)
		if err != nil {
			log.Println("===== PopulateMemoryCache: Couldn't load users from PD =====")
			return
		}

		for _, fu := range fullUsers {
			u := new(User)
			b, _ := json.Marshal(fu)
			json.Unmarshal(b, u)

			err = cachePutUser(u)
			if err != nil {
				log.Printf("===== PopulateMemoryCache: Error putting user %v to cache: %v", fu.ID, err)
			} else {
				log.Printf("===== PopulateMemoryCache: Put user %v to cache", fu.ID)
			}

			for _, c := range fu.ContactMethods {
				err = cachePutContactMethod(c)
				if err != nil {
					log.Printf("===== PopulateMemoryCache: Error putting contact method %v to cache: %v", c.ID, err)
				} else {
					log.Printf("===== PopulateMemoryCache: Put contact method %v to cache", c.ID)
				}
			}

			for _, r := range fu.NotificationRules {
				err = cachePutNotificationRule(r)
				if err != nil {
					log.Printf("===== getFullUserToCache: Error putting notification rule %v to cache: %v", r.ID, err)
				} else {
					log.Printf("===== getFullUserToCache: Put notification rule %v to cache", r.ID)
				}
			}
		}
	}
}

// PopulateMongoCache does initial population of the cache if Mongo caching is selected
func PopulateMongoCache() {
	filter := bson.D{primitive.E{Key: "ID", Value: "lastrefresh"}}
	lastRefreshRecord := new(cacheLastRefreshRecord)
	err := mongoCache["misc"].FindOne(context.TODO(), filter).Decode(lastRefreshRecord)
	if err == nil {
		if time.Since(lastRefreshRecord.Users) < cacheMaxAge {
			log.Printf("===== PagerDuty cache was refreshed at %s, not refreshing =====", lastRefreshRecord.Users.Format(time.RFC3339))
			return
		}
		log.Printf("===== PagerDuty cache was refreshed at %s, refreshing =====", lastRefreshRecord.Users.Format(time.RFC3339))
	}

	var pdo = ListUsersOptions{
		Include: []string{"contact_methods", "notification_rules"},
		Limit:   100,
	}

	fullUsers, err := pdClient.Users.ListAll(&pdo)
	if err != nil {
		log.Println("===== Couldn't load users =====")
		return
	}

	users := make([]interface{}, len(fullUsers))
	var contactMethods []interface{}
	var notificationRules []interface{}
	for i := 0; i < len(fullUsers); i++ {
		user := new(User)
		b, _ := json.Marshal(fullUsers[i])
		json.Unmarshal(b, user)
		users[i] = &user

		for j := 0; j < len(fullUsers[i].ContactMethods); j++ {
			contactMethods = append(contactMethods, &(fullUsers[i].ContactMethods[j]))
		}

		for j := 0; j < len(fullUsers[i].NotificationRules); j++ {
			notificationRules = append(notificationRules, &(fullUsers[i].NotificationRules[j]))
		}
	}

	abilities, _, _ := pdClient.Abilities.List()

	abilitiesRecord := &cacheAbilitiesRecord{
		ID:        "abilities",
		Abilities: abilities,
	}

	mongoCache["users"].Drop(context.TODO())
	if len(users) > 0 {
		res, err := mongoCache["users"].InsertMany(context.TODO(), users)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted %d users", len(res.InsertedIDs))
	}

	mongoCache["contact_methods"].Drop(context.TODO())
	if len(contactMethods) > 0 {
		res, err := mongoCache["contact_methods"].InsertMany(context.TODO(), contactMethods)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted %d contact methods", len(res.InsertedIDs))
	}

	mongoCache["notification_rules"].Drop(context.TODO())
	if len(notificationRules) > 0 {
		res, err := mongoCache["notification_rules"].InsertMany(context.TODO(), notificationRules)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted %d notification rules", len(res.InsertedIDs))
	}

	mongoCache["misc"].Drop(context.TODO())
	ares, err := mongoCache["misc"].InsertOne(context.TODO(), &abilitiesRecord)
	log.Println(ares)
	if err != nil {
		log.Fatal(err)
	}

	cacheLastRefreshRecord := &cacheLastRefreshRecord{
		ID:        "lastrefresh",
		Users:     time.Now(),
		Abilities: time.Now(),
	}
	cres, err := mongoCache["misc"].InsertOne(context.TODO(), &cacheLastRefreshRecord)
	log.Println(cres)
	if err != nil {
		log.Fatal(err)
	}
}

// PopulateCache does initial population of the cache
func PopulateCache() {
	if cacheType == "mongo" {
		PopulateMongoCache()
	} else if cacheType == "memory" {
		PopulateMemoryCache()
	}
}

func getFullUserToCache(id string, v interface{}) error {
	fu, _, err := pdClient.Users.GetFull(id)
	if err != nil {
		log.Printf("===== getFullUserToCache: Error getting user %v from PD: %v", id, err)
		return err
	}

	u := new(User)
	b, _ := json.Marshal(fu)
	json.Unmarshal(b, u)
	json.Unmarshal(b, v)
	err = cachePutUser(u)
	if err != nil {
		log.Printf("===== getFullUserToCache: Error putting user %v to cache: %v", id, err)
		return err
	}
	log.Printf("===== getFullUserToCache: Put user %v to cache", id)

	for _, c := range fu.ContactMethods {
		err = cachePutContactMethod(c)
		if err != nil {
			log.Printf("===== getFullUserToCache: Error putting contact method %v to cache: %v", c.ID, err)
			return err
		}
		log.Printf("===== getFullUserToCache: Put contact method %v to cache", c.ID)
	}

	for _, r := range fu.NotificationRules {
		err = cachePutNotificationRule(r)
		if err != nil {
			log.Printf("===== getFullUserToCache: Error putting notification rule %v to cache: %v", r.ID, err)
			return err
		}
		log.Printf("===== getFullUserToCache: Put notification rule %v to cache", r.ID)
	}
	return nil
}

func memoryCacheGet(collectionName string, id string, v interface{}) error {
	log.Printf("===== memoryCacheGet %v from %v", id, collectionName)
	if collection, ok := memoryCache[collectionName]; ok {
		if item, ok := collection.Load(id); ok {
			err := json.Unmarshal(item.([]byte), v)
			if err != nil {
				log.Printf("===== memoryCacheGet Error unmarshaling JSON getting %v from %q: %v", id, collectionName, err)
				return err
			}
			log.Printf("===== memoryCacheGet Got %v from %q cache", id, collectionName)
			return nil
		} else if collectionName == "users" {
			// special case for filling users into memory cache on demand
			return getFullUserToCache(id, v)
		} else {
			return fmt.Errorf("memoryCacheGet Item %q is not in %q hash", id, collectionName)
		}
	} else {
		return fmt.Errorf("memoryCacheGet No such collection: %q", collectionName)
	}
}

func mongoCacheGet(collectionName string, id string, v interface{}) error {
	if collection, ok := mongoCache[collectionName]; ok {
		filter := bson.D{primitive.E{Key: "id", Value: id}}
		r := collection.FindOne(context.TODO(), filter)
		err := r.Decode(v)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("mongoCacheGet No such collection: %q", collectionName)
}

func cacheGet(collectionName string, id string, v interface{}) error {
	if cacheType == "mongo" {
		return mongoCacheGet(collectionName, id, v)
	} else if cacheType == "memory" {
		return memoryCacheGet(collectionName, id, v)
	}
	return fmt.Errorf("cacheGet Cache is not enabled")
}

func mongoCachePut(collectionName string, id string, v interface{}) error {
	if collection, ok := mongoCache[collectionName]; ok {
		filter := bson.D{primitive.E{Key: "id", Value: id}}
		opts := options.Replace().SetUpsert(true)
		res, err := collection.ReplaceOne(context.TODO(), filter, &v, opts)
		if err != nil {
			log.Printf("===== Error updating %v: %q", collectionName, err)
			return err
		}
		if res.MatchedCount != 0 {
			log.Printf("===== replaced an existing item %q in %v cache", id, collectionName)
			return nil
		}
		if res.UpsertedCount != 0 {
			log.Printf("===== inserted a new item %q in %v cache", id, collectionName)
		}
		return nil
	}
	return fmt.Errorf("no such collection %q", collectionName)
}

func memoryCachePut(collectionName string, id string, v interface{}) error {
	if collection, ok := memoryCache[collectionName]; ok {
		b, _ := json.Marshal(v)
		collection.Store(id, b)
		return nil
	}
	return fmt.Errorf("no such collection: %q", collectionName)
}

func cachePut(collectionName string, id string, v interface{}) error {
	if cacheType == "mongo" {
		return mongoCachePut(collectionName, id, v)
	} else if cacheType == "memory" {
		return memoryCachePut(collectionName, id, v)
	}
	return fmt.Errorf("cachePut Cache is not enabled")
}

func mongoCacheDelete(collectionName string, id string) error {
	if collection, ok := mongoCache[collectionName]; ok {
		filter := bson.D{primitive.E{Key: "id", Value: id}}
		_, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Printf("===== mongoCacheDelete mongo error: %q", err)
			return err
		}
		log.Printf("===== mongoCacheDetele deleted item %v from %q", id, collectionName)
		return nil
	}
	return fmt.Errorf("mongoCacheDelete No such collection %q", collectionName)
}

func memoryCacheDelete(collectionName string, id string) error {
	if collection, ok := memoryCache[collectionName]; ok {
		collection.Delete(id)
		log.Printf("===== memoryCacheDelete deleted item %v from %q", id, collectionName)
		return nil
	}
	return fmt.Errorf("memoryCacheDelete No such collection: %q", collectionName)
}

func cacheDelete(collectionName string, id string) error {
	if cacheType == "mongo" {
		return mongoCacheDelete(collectionName, id)
	} else if cacheType == "memory" {
		return memoryCacheDelete(collectionName, id)
	}
	return fmt.Errorf("cacheDelete Cache is not enabled")
}

func cacheGetAbilities(v interface{}) error {
	r := new(cacheAbilitiesRecord)
	err := cacheGet("misc", "abilities", r)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(r)
	json.Unmarshal(b, v)
	return nil
}

func cacheGetUser(id string, v interface{}) error {
	return cacheGet("users", id, v)
}

func cachePutUser(u *User) error {
	return cachePut("users", u.ID, u)
}

func cacheDeleteUser(id string) error {
	return cacheDelete("users", id)
}

func cacheGetContactMethod(id string, v interface{}) error {
	return cacheGet("contact_methods", id, v)
}

func cachePutContactMethod(c *ContactMethod) error {
	return cachePut("contact_methods", c.ID, c)
}

func cacheDeleteContactMethod(id string) error {
	return cacheDelete("contact_methods", id)
}

func cacheGetNotificationRule(id string, v interface{}) error {
	return cacheGet("notification_rules", id, v)
}

func cachePutNotificationRule(r *NotificationRule) error {
	return cachePut("notification_rules", r.ID, r)
}

func cacheDeleteNotificationRule(id string) error {
	return cacheDelete("notification_rules", id)
}
