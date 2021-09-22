package pagerduty

import "fmt"

// AbilityService handles the communication with ability related methods
// of the PagerDuty API.
type AbilityService service

// ListAbilitiesResponse represents a list response of abilities.
type ListAbilitiesResponse struct {
	Abilities []string `json:"abilities,omitempty"`
}

// Test tests whether the account has a given ability.
func (s *AbilityService) Test(id string) (*Response, error) {
	u := fmt.Sprintf("/abilities/%s", id)
	return s.client.newRequestDo("GET", u, nil, nil, nil)
}

// List lists available abilities.
func (s *AbilityService) List() (*ListAbilitiesResponse, *Response, error) {
	u := "/abilities"
	v := new(ListAbilitiesResponse)

	err := cacheGetAbilities(v)
	if err == nil {
		return v, nil, nil
	}

	resp, err := s.client.newRequestDo("GET", u, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
