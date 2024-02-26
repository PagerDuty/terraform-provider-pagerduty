package pagerduty

import (
	"fmt"
	"log"
)

// TeamService handles the communication with team
// related methods of the PagerDuty API.
type TeamService service

// Team represents a team.
type Team struct {
	Description string         `json:"description,omitempty"`
	HTMLURL     string         `json:"html_url,omitempty"`
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Self        string         `json:"self,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	Type        string         `json:"type,omitempty"`
	Parent      *TeamReference `json:"parent,omitempty"`
	DefaultRole string         `json:"default_role,omitempty"`
}

// Member represents a team member.
type Member struct {
	User *UserReference `json:"user,omitempty"`
	Role string         `json:"role,omitempty"`
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
	Limit  int     `json:"limit,omitempty"`
	More   bool    `json:"more,omitempty"`
	Offset int     `json:"offset,omitempty"`
	Total  int     `json:"total,omitempty"`
	Teams  []*Team `json:"teams,omitempty"`
}

// GetMembersOptions represents options when getting a list of members.
type GetMembersOptions struct {
	Limit    int      `url:"limit,omitempty"`
	More     bool     `url:"more,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Total    int      `url:"total,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
}

// GetMembersResponse represents a response of a list of members.
type GetMembersResponse struct {
	Limit   int       `json:"limit,omitempty"`
	More    bool      `json:"more,omitempty"`
	Offset  int       `json:"offset,omitempty"`
	Total   int       `json:"total,omitempty"`
	Members []*Member `json:"members,omitempty"`
}

type teamRole struct {
	Role string `json:"role,omitempty"`
}

// TeamPayload represents a team.
type TeamPayload struct {
	Team *Team `json:"team,omitempty"`
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
	v := new(TeamPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &TeamPayload{Team: team}, &v)
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
	v := new(TeamPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}

// Update updates an existing team.
func (s *TeamService) Update(id string, team *Team) (*Team, *Response, error) {
	u := fmt.Sprintf("/teams/%s", id)
	v := new(TeamPayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &TeamPayload{Team: team}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}

// RemoveUser removes a user from a team.
func (s *TeamService) RemoveUser(teamID, userID string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s/users/%s", teamID, userID)
	resp, err := s.client.newRequestDo("DELETE", u, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = cacheDeleteTeamMembership(teamID, userID); err != nil {
		log.Printf("===== Error deleting user %q from team %q cache: %q", userID, teamID, err)
	} else {
		log.Printf("===== Deleted user %q from team %q cache", userID, teamID)
	}

	return resp, nil
}

// AddUser adds a user to a team.
func (s *TeamService) AddUser(teamID, userID string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s/users/%s", teamID, userID)
	return s.client.newRequestDo("PUT", u, nil, nil, nil)
}

// AddUserWithRole adds a user with the specified role (one of observer, manager, or responder[default])
func (s *TeamService) AddUserWithRole(teamID, userID string, role string) (*Response, error) {
	tr := teamRole{Role: role}
	u := fmt.Sprintf("/teams/%s/users/%s", teamID, userID)
	resp, err := s.client.newRequestDo("PUT", u, nil, tr, nil)
	if err != nil {
		return nil, err
	}

	if err = cachePutTeamMembership(teamID, userID, role); err != nil {
		log.Printf("===== Error adding user %q to team %q cache: %q", userID, teamID, err)
	} else {
		log.Printf("===== Added user %q to team %q cache", userID, teamID)
	}

	return resp, nil
}

// GetMembers retrieves information about members on a team.
func (s *TeamService) GetMembers(teamID string, o *GetMembersOptions) (*GetMembersResponse, *Response, error) {
	u := fmt.Sprintf("/teams/%s/members", teamID)
	v := new(GetMembersResponse)

	members := make([]*Member, 0)

	cm := new(GetMembersResponse)
	if err := cacheGetTeamMembers(teamID, cm); err == nil && len(cm.Members) > 0 {
		members = append(members, cm.Members...)
		v.Members = members
		return v, nil, nil
	} else {
		log.Printf("[DEBUG] error retrieving team members %q; %v", teamID, err)
	}

	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result GetMembersResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		members = append(members, result.Members...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(u, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.Members = members

	if err = cachePutTeamMembers(teamID, &GetMembersResponse{Members: members}); err != nil {
		log.Printf("===== Error adding members of team %q to cache: %q", teamID, err)
	} else {
		log.Printf("===== Added members of team %q to cache", teamID)
	}

	return v, nil, nil
}

// RemoveEscalationPolicy removes an escalation policy from a team.
func (s *TeamService) RemoveEscalationPolicy(teamID, escID string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s/escalation_policies/%s", teamID, escID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// AddEscalationPolicy adds an escalation policy to a team.
func (s *TeamService) AddEscalationPolicy(teamID, escID string) (*Response, error) {
	u := fmt.Sprintf("/teams/%s/escalation_policies/%s", teamID, escID)
	return s.client.newRequestDo("PUT", u, nil, nil, nil)
}
