package pagerduty

// PriorityService handles the communication with priority related methods
// of the PagerDuty API.
type PriorityService service

// ListPrioritiesResponse represents a list response of abilities.
type ListPrioritiesResponse struct {
	Total      int         `json:"total,omitempty"`
	Offset     int         `json:"offset,omitempty"`
	More       bool        `json:"more,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	Priorities []*Priority `json:"priorities,omitempty"`
}

// Priority represents a priority
type Priority struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// List lists available priorities.
func (s *PriorityService) List() (*ListPrioritiesResponse, *Response, error) {
	u := "/priorities"
	v := new(ListPrioritiesResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
