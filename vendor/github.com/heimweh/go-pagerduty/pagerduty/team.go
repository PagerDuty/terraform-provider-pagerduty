package pagerduty

import "fmt"

// TeamService handles the communication with team
// related methods of the PagerDuty API.
type TeamService service

// Team represents a team.
type Team struct {
	Description string `json:"description,omitempty"`
	HTMLURL     string `json:"html_url,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Self        string `json:"self,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Team        *Team  `json:"team,omitempty"`
	Type        string `json:"type,omitempty"`
}

// ListTeamsOptions represents options when listing teams.
type ListTeamsOptions struct {
	Limit  int    `url:"limit,omitempty"`
	More   bool   `url:"more,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Total  int    `url:"total,omitempty"`
	Query  string `url:"query,omitempty"`
}

// ListTeamsResponse represents a list response of teams.
type ListTeamsResponse struct {
	Limit  int     `url:"limit,omitempty"`
	More   bool    `url:"more,omitempty"`
	Offset int     `url:"offset,omitempty"`
	Total  int     `url:"total,omitempty"`
	Teams  []*Team `json:"teams,omitempty"`
}

// List lists existing teams.
func (s *TeamService) List(o *ListTeamsOptions) (*ListTeamsResponse, *Response, error) {
	u := "/teams"
	v := new(ListTeamsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new team.
func (s *TeamService) Create(team *Team) (*Team, *Response, error) {
	u := "/teams"
	v := new(Team)

	resp, err := s.client.newRequestDo("POST", u, nil, &Team{Team: team}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}

// Delete removes an existing team.
func (s *TeamService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about a team.
func (s *TeamService) Get(id string) (*Team, *Response, error) {
	u := fmt.Sprintf("/teams/%s", id)
	v := new(Team)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}

// Update updates an existing team.
func (s *TeamService) Update(id string, team *Team) (*Team, *Response, error) {
	u := fmt.Sprintf("/teams/%s", id)
	v := new(Team)

	resp, err := s.client.newRequestDo("PUT", u, nil, &Team{Team: team}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}

// RemoveUser removes a user from a team.
func (s *TeamService) RemoveUser(teamID, userID string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s/users/%s", teamID, userID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// AddUser adds a user to a team.
func (s *TeamService) AddUser(teamID, userID string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s/users/%s", teamID, userID)
	return s.client.newRequestDo("PUT", u, nil, nil, nil)
}
