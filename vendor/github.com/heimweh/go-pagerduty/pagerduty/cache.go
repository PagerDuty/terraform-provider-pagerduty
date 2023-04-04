package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
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
	"team_members":       {},
	"misc":               {},
}

type cacheAbilitiesRecord struct {
	ID        string
	Abilities *ListAbilitiesResponse
}

type cacheTeamMembersRecord struct {
	TeamID string
	UserID string
	Member *Member
}

type cacheLastRefreshRecord struct {
	ID          string
	Users       time.Time
	Abilities   time.Time
	TeamMembers time.Time
}

// InitCache initializes the cache according to the setting in TF_PAGERDUTY_CACHE
func InitCache(c *Client) {
	pdClient = c
	cacheMongoURL = os.Getenv("TF_PAGERDUTY_CACHE")
	re := regexp.MustCompile("^mongodb+(\\+srv)?://")
	isMongodbURL := re.Match([]byte(cacheMongoURL))
	if isMongodbURL {
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
		"team_members":       mongoClient.Database("pagerduty").Collection("team_members"),
		"misc":               mongoClient.Database("pagerduty").Collection("misc"),
	}
}

// PopulateMemoryCache does initial population of the cache if memory caching is selected
func PopulateMemoryCache() {
	if _, present := os.LookupEnv("TF_PAGERDUTY_CACHE_PREFILL"); !present {
		return
	}

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

// PopulateMongoCache does initial population of the cache if Mongo caching is selected
func PopulateMongoCache() {
	cacheTypeRefreshedAt := map[string]*time.Time{
		"Users":       nil,
		"Abilities":   nil,
		"TeamMembers": nil,
	}
	filter := bson.D{primitive.E{Key: "id", Value: "lastrefresh"}}
	lastRefreshRecord := new(cacheLastRefreshRecord)
	err := mongoCache["misc"].FindOne(context.TODO(), filter).Decode(lastRefreshRecord)
	if err == nil {
		cacheTypeRefreshedAt["Users"] = &lastRefreshRecord.Users
		cacheTypeRefreshedAt["Abilities"] = &lastRefreshRecord.Abilities
		cacheTypeRefreshedAt["TeamMembers"] = &lastRefreshRecord.TeamMembers
	}

	cachedTypeRefreshed := map[string]bool{
		"Users":       true,
		"Abilities":   true,
		"TeamMembers": true,
	}

	err = refreshMongoUsersCache(cacheTypeRefreshedAt["Users"])
	if err != nil {
		log.Printf("===== PagerDuty Mongo cache, error while refreshing Users cache: %v", err)
		cachedTypeRefreshed["Users"] = false
	}
	err = nil

	err = refreshMongoAbilitiesCache(cacheTypeRefreshedAt["Abilities"])
	if err != nil {
		log.Printf("===== PagerDuty Mongo cache, error while refreshing Abilities cache: %v", err)
		cachedTypeRefreshed["Abilities"] = false
	}
	err = nil

	err = refreshMongoTeamMembersCache(cacheTypeRefreshedAt["TeamMembers"])
	if err != nil {
		log.Printf("===== PagerDuty Mongo cache, error while refreshing Team Members cache: %v", err)
		cachedTypeRefreshed["TeamMembers"] = false
	}
	err = nil

	cacheLastRefreshRecord := &cacheLastRefreshRecord{
		ID: "lastrefresh",
	}
	for k, cacheRefreshed := range cachedTypeRefreshed {
		if k == "Users" && cacheRefreshed {
			cacheLastRefreshRecord.Users = time.Now()
			continue
		}
		if k == "Abilities" && cacheRefreshed {
			cacheLastRefreshRecord.Abilities = time.Now()
			continue
		}
		if k == "TeamMembers" && cacheRefreshed {
			cacheLastRefreshRecord.TeamMembers = time.Now()
			continue
		}
	}
	opts := options.Replace().SetUpsert(true)
	cres, err := mongoCache["misc"].ReplaceOne(context.TODO(), filter, &cacheLastRefreshRecord, opts)
	log.Println(cres)
	if err != nil {
		log.Fatal(err)
	}
}

func refreshMongoUsersCache(lastRefreshed *time.Time) error {
	if needToRefreshCache := needToRefreshMongoCacheType("Users", lastRefreshed); !needToRefreshCache {
		return nil
	}

	var pdo = ListUsersOptions{
		Include: []string{"contact_methods", "notification_rules"},
		Limit:   100,
	}

	fullUsers, err := pdClient.Users.ListAll(&pdo)
	if err != nil {
		log.Println("===== Couldn't load users =====")
		return err
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

	mongoCache["users"].Drop(context.TODO())
	if len(users) > 0 {
		res, err := mongoCache["users"].InsertMany(context.TODO(), users)
		if err != nil {
			return err
		}
		log.Printf("Inserted %d users", len(res.InsertedIDs))
	}

	mongoCache["contact_methods"].Drop(context.TODO())
	if len(contactMethods) > 0 {
		res, err := mongoCache["contact_methods"].InsertMany(context.TODO(), contactMethods)
		if err != nil {
			return err
		}
		log.Printf("Inserted %d contact methods", len(res.InsertedIDs))
	}

	mongoCache["notification_rules"].Drop(context.TODO())
	if len(notificationRules) > 0 {
		res, err := mongoCache["notification_rules"].InsertMany(context.TODO(), notificationRules)
		if err != nil {
			return err
		}
		log.Printf("Inserted %d notification rules", len(res.InsertedIDs))
	}

	return nil
}

func refreshMongoAbilitiesCache(lastRefreshed *time.Time) error {
	if needToRefreshCache := needToRefreshMongoCacheType("Abilities", lastRefreshed); !needToRefreshCache {
		return nil
	}

	abilities, _, _ := pdClient.Abilities.List()

	abilitiesRecord := &cacheAbilitiesRecord{
		ID:        "abilities",
		Abilities: abilities,
	}

	mongoCache["misc"].Drop(context.TODO())
	ares, err := mongoCache["misc"].InsertOne(context.TODO(), &abilitiesRecord)
	log.Println(ares)
	if err != nil {
		return err
	}

	return nil
}

func refreshMongoTeamMembersCache(lastRefreshed *time.Time) error {
	if needToRefreshCache := needToRefreshMongoCacheType("Team Members", lastRefreshed); !needToRefreshCache {
		return nil
	}

	// Since `team_members` doesn't need to cache pre-fill, then It's only
	// needed to remove the staled entries.
	mongoCache["team_members"].Drop(context.TODO())

	return nil
}

func needToRefreshMongoCacheType(name string, lastRefreshed *time.Time) bool {
	if lastRefreshed != nil {
		if time.Since(*lastRefreshed) < cacheMaxAge {
			log.Printf("===== PagerDuty Mongo cache for %s was refreshed at %s, not refreshing =====", name, lastRefreshed.Format(time.RFC3339))
			return false
		}
		log.Printf("===== PagerDuty Mongo cache for %s was refreshed at %s, refreshing =====", name, lastRefreshed.Format(time.RFC3339))
		return true
	}

	log.Printf("===== PagerDuty Mongo cache for %s refreshing =====", name)
	return true
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
		} else if collectionName == "team_members" {
			return memoryCacheGetTeamMembers(collection, id, v)
		} else {
			return fmt.Errorf("memoryCacheGet Item %q is not in %q hash", id, collectionName)
		}
	} else {
		return fmt.Errorf("memoryCacheGet No such collection: %q", collectionName)
	}
}

func memoryCacheGetTeamMembers(collection *sync.Map, id string, v interface{}) (err error) {
	var r []*cacheTeamMembersRecord
	collection.Range(func(key, value interface{}) bool {
		parts := strings.Split(key.(string), ":")
		teamId := parts[0]
		if teamId == id {
			member := new(cacheTeamMembersRecord)
			err = json.Unmarshal(value.([]byte), member)
			if err != nil {
				log.Printf("===== memoryCacheGet Error unmarshaling JSON getting %v from %q: %v", id, "team_members", err)
				return false
			}
			r = append(r, member)
			log.Printf("===== memoryCacheGet Got %d items of team %q from %q cache", len(r), id, "team_members")
		}
		return true
	})

	b, _ := json.Marshal(r)
	json.Unmarshal(b, v)
	return err
}

func mongoCacheGet(collectionName string, id string, v interface{}) error {
	if collection, ok := mongoCache[collectionName]; ok {
		if collectionName == "team_members" {
			filter := bson.D{primitive.E{Key: "teamid", Value: id}}
			cur, err := collection.Find(context.TODO(), filter)
			if err != nil {
				defer cur.Close(context.TODO())
				return err
			}
			var results []bson.M
			if err = cur.All(context.TODO(), &results); err != nil {
				return err
			}
			b, _ := json.Marshal(results)
			json.Unmarshal(b, v)

			return nil
		}

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
		if collectionName == "team_members" {
			return mongoCahePutMany(collection, v.([]interface{}))
		}

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

func mongoCahePutMany(collection *mongo.Collection, entries []interface{}) error {
	opts := options.InsertMany()
	res, err := collection.InsertMany(context.TODO(), entries, opts)
	if err != nil {
		log.Printf("===== Error updating %v: %q", collection.Name(), err)
		return err
	}
	if len(res.InsertedIDs) != 0 {
		log.Printf("===== inserted %d items in %v cache", len(res.InsertedIDs), collection.Name())
		return nil
	}
	return nil
}

func memoryCachePut(collectionName string, id string, v interface{}) error {
	if collection, ok := memoryCache[collectionName]; ok {
		if collectionName == "team_members" {
			return memoryCachePutTeamMembers(collection, v.([]interface{}))
		}

		b, _ := json.Marshal(v)
		collection.Store(id, b)
		return nil
	}
	return fmt.Errorf("no such collection: %q", collectionName)
}

func memoryCachePutTeamMembers(collection *sync.Map, v []interface{}) error {
	for _, entry := range v {
		var id string
		member := entry.(*cacheTeamMembersRecord)
		id = fmt.Sprintf("%s:%s", member.TeamID, member.UserID)
		b, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		collection.Store(id, b)
	}
	return nil
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

		if collectionName == "team_members" {
			parts := strings.Split(id, ":")
			teamID := parts[0]
			userID := parts[1]
			filter = bson.D{
				{Key: "teamid", Value: teamID},
				{Key: "userid", Value: userID},
			}
		}
		_, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Printf("===== mongoCacheDelete mongo error: %q", err)
			return err
		}
		log.Printf("===== mongoCacheDelete deleted item %v from %q", id, collectionName)
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

func cacheGetTeamMembers(id string, v interface{}) error {
	r := []*cacheTeamMembersRecord{}
	err := cacheGet("team_members", id, &r)
	if err != nil {
		return err
	}
	members := &GetMembersResponse{}
	for _, m := range r {
		members.Members = append(members.Members, m.Member)
	}
	b, _ := json.Marshal(members)
	json.Unmarshal(b, v)
	return nil
}

func cachePutTeamMembers(id string, m *GetMembersResponse) error {
	var members []interface{}
	for _, member := range m.Members {
		members = append(members, &cacheTeamMembersRecord{TeamID: id, UserID: member.User.ID, Member: member})
	}
	return cachePut("team_members", id, members)
}

func cachePutTeamMembership(teamID, userID, role string) error {
	cm := new(GetMembersResponse)
	members := []*Member{
		{
			Role: role,
			User: &UserReference{
				ID:   userID,
				Type: "user_reference",
			},
		},
	}
	cm.Members = members
	return cachePutTeamMembers(teamID, cm)
}

func cacheDeleteTeamMembership(teamID, userID string) error {
	return cacheDelete("team_members", fmt.Sprintf("%s:%s", teamID, userID))
}
