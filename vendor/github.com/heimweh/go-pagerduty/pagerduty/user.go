package pagerduty

import (
	"fmt"
	"log"
)

// UserService handles the communication with user
// related methods of the PagerDuty API.
type UserService service

// NotificationRule represents a user notification rule.
type NotificationRule struct {
	ContactMethod       *ContactMethodReference `json:"contact_method,omitempty"`
	HTMLURL             string                  `json:"html_url,omitempty"`
	ID                  string                  `json:"id,omitempty"`
	Self                string                  `json:"self,omitempty"`
	StartDelayInMinutes int                     `json:"start_delay_in_minutes"`
	Summary             string                  `json:"summary,omitempty"`
	Type                string                  `json:"type,omitempty"`
	Urgency             string                  `json:"urgency,omitempty"`
}

// NotificationRulePayload represents a notification rule.
type NotificationRulePayload struct {
	NotificationRule *NotificationRule `json:"notification_rule,omitempty"`
}

// User represents a user.
type User struct {
	AvatarURL         string                    `json:"avatar_url,omitempty"`
	Color             string                    `json:"color,omitempty"`
	ContactMethods    []*ContactMethodReference `json:"contact_methods,omitempty"`
	Description       string                    `json:"description,omitempty"`
	Email             string                    `json:"email,omitempty"`
	HTMLURL           string                    `json:"html_url,omitempty"`
	ID                string                    `json:"id,omitempty"`
	InvitationSent    bool                      `json:"invitation_sent,omitempty"`
	JobTitle          string                    `json:"job_title,omitempty"`
	Name              string                    `json:"name,omitempty"`
	NotificationRules []*NotificationRule       `json:"notification_rules,omitempty"`
	Role              string                    `json:"role,omitempty"`
	Self              string                    `json:"self,omitempty"`
	Summary           string                    `json:"summary,omitempty"`
	Teams             []*TeamReference          `json:"teams,omitempty"`
	TimeZone          string                    `json:"time_zone,omitempty"`
	Type              string                    `json:"type,omitempty"`
}

// UserPayload represents a user.
type UserPayload struct {
	User *User `json:"user,omitempty"`
}

// FullUser represents a user fetched with include[]=contact_methods,notification_rules.
// This is only used when caching is enabled
type FullUser struct {
	AvatarURL         string              `json:"avatar_url,omitempty"`
	Color             string              `json:"color,omitempty"`
	ContactMethods    []*ContactMethod    `json:"contact_methods,omitempty"`
	Description       string              `json:"description,omitempty"`
	Email             string              `json:"email,omitempty"`
	HTMLURL           string              `json:"html_url,omitempty"`
	ID                string              `json:"id,omitempty"`
	InvitationSent    bool                `json:"invitation_sent,omitempty"`
	JobTitle          string              `json:"job_title,omitempty"`
	Name              string              `json:"name,omitempty"`
	NotificationRules []*NotificationRule `json:"notification_rules,omitempty"`
	Role              string              `json:"role,omitempty"`
	Self              string              `json:"self,omitempty"`
	Summary           string              `json:"summary,omitempty"`
	Teams             []*Team             `json:"teams,omitempty"`
	TimeZone          string              `json:"time_zone,omitempty"`
	Type              string              `json:"type,omitempty"`
}

// FullUserPayload represents a user.
type FullUserPayload struct {
	User *FullUser `json:"user,omitempty"`
}

// ContactMethod represents a contact method for a user.
type ContactMethod struct {
	ID          string `json:"id,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Type        string `json:"type,omitempty"`
	Self        string `json:"self,omitempty"`
	HTMLURL     string `json:"html_url,omitempty"`
	Label       string `json:"label,omitempty"`
	Address     string `json:"address,omitempty"`
	BlackListed bool   `json:"blacklisted,omitempty"`

	// Email contact method options
	SendShortEmail bool `json:"send_short_email,omitempty"`

	// Phone contact method options
	CountryCode int  `json:"country_code,omitempty"`
	Enabled     bool `json:"enabled,omitempty"`

	// Push contact method options
	DeviceType string                    `json:"device_type,omitempty"`
	Sounds     []*PushContactMethodSound `json:"sounds,omitempty"`
	CreatedAt  string                    `json:"created_at,omitempty"`
}

// ContactMethodPayload represents a contact method.
type ContactMethodPayload struct {
	ContactMethod *ContactMethod `json:"contact_method"`
}

// PushContactMethodSound represents a sound for a push contact method.
type PushContactMethodSound struct {
	Type string `json:"type,omitempty"`
	File string `json:"file,omitempty"`
}

// ListContactMethodsResponse represents
type ListContactMethodsResponse struct {
	Limit          int              `json:"limit,omitempty"`
	More           bool             `json:"more,omitempty"`
	Offset         int              `json:"offset,omitempty"`
	Total          int              `json:"total,omitempty"`
	ContactMethods []*ContactMethod `json:"contact_methods,omitempty"`
}

// ListUsersOptions represents options when listing users.
type ListUsersOptions struct {
	Limit   int      `url:"limit,omitempty"`
	More    bool     `url:"more,omitempty"`
	Offset  int      `url:"offset,omitempty"`
	Total   int      `url:"total,omitempty"`
	Include []string `url:"include,omitempty,brackets"`
	Query   string   `url:"query,omitempty"`
	TeamIDs []string `url:"team_ids,omitempty,brackets"`
}

// ListUsersResponse represents a list response of users.
type ListUsersResponse struct {
	Limit  int     `json:"limit,omitempty"`
	More   bool    `json:"more,omitempty"`
	Offset int     `json:"offset,omitempty"`
	Total  int     `json:"total,omitempty"`
	Users  []*User `json:"users,omitempty"`
}

// ListFullUsersResponse represents a list response containing FullUser objects.
type ListFullUsersResponse struct {
	Limit  int         `json:"limit,omitempty"`
	More   bool        `json:"more,omitempty"`
	Offset int         `json:"offset,omitempty"`
	Total  int         `json:"total,omitempty"`
	Users  []*FullUser `json:"users,omitempty"`
}

// GetUserOptions represents options when retrieving a user.
type GetUserOptions struct {
	Include []string `url:"include,omitempty,brackets"`
}

// List lists existing users.
func (s *UserService) List(o *ListUsersOptions) (*ListUsersResponse, *Response, error) {
	u := "/users"
	v := new(ListUsersResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// ListAll lists users into FullUser objects
func (s *UserService) ListAll(o *ListUsersOptions) ([]*FullUser, error) {
	var users = make([]*FullUser, 0, 25)
	var v *ListFullUsersResponse
	more := true
	offset := 0

	for more {
		log.Printf("==== Getting users at offset %d", offset)
		v = new(ListFullUsersResponse)
		_, err := s.client.newRequestDo("GET", "/users", o, nil, &v)
		if err != nil {
			return users, err
		}
		users = append(users, v.Users...)
		more = v.More
		offset += v.Limit
		o.Offset = offset
	}
	return users, nil
}

// Create creates a new user.
func (s *UserService) Create(user *User) (*User, *Response, error) {
	u := "/users"
	v := new(UserPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &UserPayload{User: user}, &v)
	if err != nil {
		return nil, nil, err
	}

	if err = cachePutUser(v.User); err != nil {
		log.Printf("===== Error adding user %q to cache: %q", v.User.ID, err)
	} else {
		log.Printf("===== Added user %q to cache", v.User.ID)
	}

	return v.User, resp, nil
}

// Delete removes an existing user.
func (s *UserService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	resp, err := s.client.newRequestDo("DELETE", u, nil, nil, nil)

	if cerr := cacheDeleteUser(id); cerr != nil {
		log.Printf("===== Error deleting user %q from cache: %q", id, cerr)
	} else {
		log.Printf("===== Deleted user %q from cache", id)
	}

	return resp, err
}

// Get retrieves information about a user.
func (s *UserService) Get(id string, o *GetUserOptions) (*User, *Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	v := new(UserPayload)

	if err := cacheGetUser(id, v); err == nil {
		return v.User, nil, nil
	}

	resp, err := s.client.newRequestDo("GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.User, resp, nil
}

// GetFull retrieves information about a user including contact methods and notification rules.
func (s *UserService) GetFull(id string) (*FullUser, *Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	v := new(FullUserPayload)
	o := &GetUserOptions{
		Include: []string{"contact_methods", "notification_rules"},
	}

	resp, err := s.client.newRequestDo("GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.User, resp, nil
}

// Update updates an existing user.
func (s *UserService) Update(id string, user *User) (*User, *Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	v := new(UserPayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &UserPayload{User: user}, &v)
	if err != nil {
		return nil, nil, err
	}

	cachePutUser(v.User)

	return v.User, resp, nil
}

// ListContactMethods lists contact methods for a user.
func (s *UserService) ListContactMethods(userID string) (*ListContactMethodsResponse, *Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods", userID)
	v := new(ListContactMethodsResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// CreateContactMethod creates a new contact method for a user.
func (s *UserService) CreateContactMethod(userID string, contactMethod *ContactMethod) (*ContactMethod, *Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods", userID)
	v := new(ContactMethodPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &ContactMethodPayload{ContactMethod: contactMethod}, &v)
	if err != nil {
		return nil, nil, err
	}

	if err = cachePutContactMethod(v.ContactMethod); err != nil {
		log.Printf("===== Error adding contact method %q to cache: %q", v.ContactMethod.ID, err)
	} else {
		log.Printf("===== Added contact method %q to cache", v.ContactMethod.ID)
	}

	return v.ContactMethod, resp, nil
}

// GetContactMethod retrieves a contact method for a user.
func (s *UserService) GetContactMethod(userID string, contactMethodID string) (*ContactMethod, *Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods/%s", userID, contactMethodID)
	v := new(ContactMethodPayload)

	if err := cacheGetContactMethod(contactMethodID, v); err == nil {
		return v.ContactMethod, nil, nil
	}

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.ContactMethod, resp, nil
}

// UpdateContactMethod updates a contact method for a user.
func (s *UserService) UpdateContactMethod(userID, contactMethodID string, contactMethod *ContactMethod) (*ContactMethod, *Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods/%s", userID, contactMethodID)
	v := new(ContactMethodPayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &ContactMethodPayload{ContactMethod: contactMethod}, &v)
	if err != nil {
		return nil, nil, err
	}

	cachePutContactMethod(v.ContactMethod)

	return v.ContactMethod, resp, nil
}

// DeleteContactMethod deletes a contact method for a user.
func (s *UserService) DeleteContactMethod(userID, contactMethodID string) (*Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods/%s", userID, contactMethodID)
	resp, err := s.client.newRequestDo("DELETE", u, nil, nil, nil)

	if cerr := cacheDeleteContactMethod(contactMethodID); cerr != nil {
		log.Printf("===== Error deleting contact method %q from cache: %q", contactMethodID, cerr)
	} else {
		log.Printf("===== Deleted contact method %q from cache", contactMethodID)
	}

	return resp, err
}

// CreateNotificationRule creates a new notification rule for a user.
func (s *UserService) CreateNotificationRule(userID string, rule *NotificationRule) (*NotificationRule, *Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules", userID)
	v := new(NotificationRulePayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &NotificationRulePayload{NotificationRule: rule}, &v)
	if err != nil {
		return nil, nil, err
	}

	if err = cachePutNotificationRule(v.NotificationRule); err != nil {
		log.Printf("===== Error adding notification rule %q to cache: %q", v.NotificationRule.ID, err)
	} else {
		log.Printf("===== Added notification rule %q to cache", v.NotificationRule.ID)
	}

	return v.NotificationRule, resp, nil
}

// GetNotificationRule retrieves a notification rule for a user.
func (s *UserService) GetNotificationRule(userID string, ruleID string) (*NotificationRule, *Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules/%s", userID, ruleID)
	v := new(NotificationRulePayload)

	if err := cacheGetNotificationRule(ruleID, v); err == nil {
		return v.NotificationRule, nil, nil
	}

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.NotificationRule, resp, nil
}

// UpdateNotificationRule updates a notification rulefor a user.
func (s *UserService) UpdateNotificationRule(userID, ruleID string, rule *NotificationRule) (*NotificationRule, *Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules/%s", userID, ruleID)
	v := new(NotificationRulePayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &NotificationRulePayload{NotificationRule: rule}, &v)
	if err != nil {
		return nil, nil, err
	}

	cachePutNotificationRule(v.NotificationRule)

	return v.NotificationRule, resp, nil
}

// DeleteNotificationRule deletes a notification rule for a user.
func (s *UserService) DeleteNotificationRule(userID, ruleID string) (*Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules/%s", userID, ruleID)
	resp, err := s.client.newRequestDo("DELETE", u, nil, nil, nil)

	if cerr := cacheDeleteNotificationRule(ruleID); cerr != nil {
		log.Printf("===== Error deleting notification rule %q from cache: %q", ruleID, cerr)
	} else {
		log.Printf("===== Deleted notification rule %q from cache", ruleID)
	}

	return resp, err
}
