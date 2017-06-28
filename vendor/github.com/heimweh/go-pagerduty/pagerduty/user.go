package pagerduty

import "fmt"

// UserService handles the communication with user
// related methods of the PagerDuty API.
type UserService service

// NotificationRule represents a user notification rule.
type NotificationRule struct {
	ContactMethod       *ContactMethodReference `json:"contact_method,omitempty"`
	HTMLURL             string                  `json:"html_url,omitempty"`
	ID                  string                  `json:"id,omitempty"`
	Self                string                  `json:"self,omitempty"`
	StartDelayInMinutes int                     `json:"start_delay_in_minutes,omitempty"`
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

// ListUsersOptions represents options when listing users.
type ListUsersOptions struct {
	*Pagination
	Include []string `url:"include,omitempty,brackets"`
	Query   string   `url:"query,omitempty"`
	TeamIDs []string `url:"team_ids,omitempty,brackets"`
}

// ListUsersResponse represents a list response of users.
type ListUsersResponse struct {
	*Pagination
	Users []*User `json:"users,omitempty"`
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
