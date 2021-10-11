package pagerduty

import "fmt"

// AddonService handles the communication with add-on related methods
// of the PagerDuty API.
type AddonService service

// Addon represents a PagerDuty add-on.
type Addon struct {
	HTMLURL string `json:"html_url,omitempty"`
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Self    string `json:"self,omitempty"`
	Src     string `json:"src,omitempty"`
	Summary string `json:"summary,omitempty"`
	Type    string `json:"type,omitempty"`
}

// ListAddonsOptions represents options when listing add-ons.
type ListAddonsOptions struct {
	Limit      int      `url:"limit,omitempty"`
	More       bool     `url:"more,omitempty"`
	Offset     int      `url:"offset,omitempty"`
	Total      int      `url:"total,omitempty"`
	Filter     string   `url:"filter,omitempty"`
	Include    []string `url:"include,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
}

// ListAddonsResponse represents a list response of add-ons.
type ListAddonsResponse struct {
	Limit  int      `json:"limit,omitempty"`
	More   bool     `json:"more,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Total  int      `json:"total,omitempty"`
	Addons []*Addon `json:"addons,omitempty"`
}

// AddonPayload represents an addon.
type AddonPayload struct {
	Addon *Addon `json:"addon,omitempty"`
}

// List lists installed add-ons.
func (s *AddonService) List(o *ListAddonsOptions) (*ListAddonsResponse, *Response, error) {
	u := "/addons"
	v := new(ListAddonsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Install installs an add-on.
func (s *AddonService) Install(addon *Addon) (*Addon, *Response, error) {
	u := "/addons"
	v := new(AddonPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &AddonPayload{Addon: addon}, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Addon, resp, nil
}

// Delete removes an existing add-on.
func (s *AddonService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/addons/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about an add-on.
func (s *AddonService) Get(id string) (*Addon, *Response, error) {
	u := fmt.Sprintf("/addons/%s", id)
	v := new(AddonPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Addon, resp, nil
}

// Update updates an existing add-on.
func (s *AddonService) Update(id string, addon *Addon) (*Addon, *Response, error) {
	u := fmt.Sprintf("/addons/%s", id)
	v := new(AddonPayload)
	resp, err := s.client.newRequestDo("PUT", u, nil, &AddonPayload{Addon: addon}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Addon, resp, nil
}
