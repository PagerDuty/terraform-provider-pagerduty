package pagerduty

import (
	"fmt"
)

// UserService handles the communication with user
// related methods of the PagerDuty API.
type UserService service

// NotificationRule represents a user notification rule.
type NotificationRule struct {
	NotificationRule    *NotificationRule       `json:"notification_rule,omitempty"`
	ContactMethod       *ContactMethodReference `json:"contact_method,omitempty"`
	HTMLURL             string                  `json:"html_url,omitempty"`
	ID                  string                  `json:"id,omitempty"`
	Self                string                  `json:"self,omitempty"`
	StartDelayInMinutes int                     `json:"start_delay_in_minutes"`
	Summary             string                  `json:"summary,omitempty"`
	Type                string                  `json:"type,omitempty"`
	Urgency             string                  `json:"urgency,omitempty"`
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
	User              *User                     `json:"user,omitempty"`
}

// ContactMethod represents a contact method for a user.
type ContactMethod struct {
	ContactMethod *ContactMethod `json:"contact_method,omitempty"`
	ID            string         `json:"id,omitempty"`
	Summary       string         `json:"summary,omitempty"`
	Type          string         `json:"type,omitempty"`
	Self          string         `json:"self,omitempty"`
	HTMLURL       string         `json:"html_url,omitempty"`
	Label         string         `json:"label,omitempty"`
	Address       string         `json:"address,omitempty"`
	BlackListed   bool           `json:"blacklisted,omitempty"`

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

// Create creates a new user.
func (s *UserService) Create(user *User) (*User, *Response, error) {
	u := "/users"
	v := new(User)

	resp, err := s.client.newRequestDo("POST", u, nil, &User{User: user}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.User, resp, nil
}

// Delete removes an existing user.
func (s *UserService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about a user.
func (s *UserService) Get(id string, o *GetUserOptions) (*User, *Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	v := new(User)

	resp, err := s.client.newRequestDo("GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.User, resp, nil
}

// Update updates an existing user.
func (s *UserService) Update(id string, user *User) (*User, *Response, error) {
	u := fmt.Sprintf("/users/%s", id)
	v := new(User)

	resp, err := s.client.newRequestDo("PUT", u, nil, &User{User: user}, &v)
	if err != nil {
		return nil, nil, err
	}

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
	v := new(ContactMethod)

	resp, err := s.client.newRequestDo("POST", u, nil, &ContactMethod{ContactMethod: contactMethod}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.ContactMethod, resp, nil
}

// GetContactMethod retrieves a contact method for a user.
func (s *UserService) GetContactMethod(userID string, contactMethodID string) (*ContactMethod, *Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods/%s", userID, contactMethodID)
	v := new(ContactMethod)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.ContactMethod, resp, nil
}

// UpdateContactMethod updates a contact method for a user.
func (s *UserService) UpdateContactMethod(userID, contactMethodID string, contactMethod *ContactMethod) (*ContactMethod, *Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods/%s", userID, contactMethodID)
	v := new(ContactMethod)

	resp, err := s.client.newRequestDo("PUT", u, nil, &ContactMethod{ContactMethod: contactMethod}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.ContactMethod, resp, nil
}

// DeleteContactMethod deletes a contact method for a user.
func (s *UserService) DeleteContactMethod(userID, contactMethodID string) (*Response, error) {
	u := fmt.Sprintf("/users/%s/contact_methods/%s", userID, contactMethodID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// CreateNotificationRule creates a new notification rule for a user.
func (s *UserService) CreateNotificationRule(userID string, rule *NotificationRule) (*NotificationRule, *Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules", userID)
	v := new(NotificationRule)

	resp, err := s.client.newRequestDo("POST", u, nil, &NotificationRule{NotificationRule: rule}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.NotificationRule, resp, nil
}

// GetNotificationRule retrieves a notification rule for a user.
func (s *UserService) GetNotificationRule(userID string, ruleID string) (*NotificationRule, *Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules/%s", userID, ruleID)
	v := new(NotificationRule)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.NotificationRule, resp, nil
}

// UpdateNotificationRule updates a notification rulefor a user.
func (s *UserService) UpdateNotificationRule(userID, ruleID string, rule *NotificationRule) (*NotificationRule, *Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules/%s", userID, ruleID)
	v := new(NotificationRule)

	resp, err := s.client.newRequestDo("PUT", u, nil, &NotificationRule{NotificationRule: rule}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.NotificationRule, resp, nil
}

// DeleteNotificationRule deletes a notification rule for a user.
func (s *UserService) DeleteNotificationRule(userID, ruleID string) (*Response, error) {
	u := fmt.Sprintf("/users/%s/notification_rules/%s", userID, ruleID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
