package pagerduty

import "fmt"

// ExtensionService handles the communication with extension related methods
// of the PagerDuty API.
type ExtensionService service

// Extension represents an extension.
type Extension struct {
	ID               string                    `json:"id,omitempty"`
	Summary          string                    `json:"summary,omitempty"`
	Type             string                    `json:"type,omitempty"`
	Self             string                    `json:"self,omitempty"`
	HTMLURL          string                    `json:"html_url,omitempty"`
	Name             string                    `json:"name"`
	EndpointURL      string                    `json:"endpoint_url,omitempty"`
	ExtensionObjects []*ServiceReference       `json:"extension_objects,omitempty"`
	ExtensionSchema  *ExtensionSchemaReference `json:"extension_schema"`
	Config           interface{}               `json:"config,omitempty"`
}

// ListExtensionsOptions represents options when listing extensions.
type ListExtensionsOptions struct {
	ExtensionObjectID string   `url:"extension_object_id,omitempty"`
	Query             string   `url:"query,omitempty"`
	ExtensionSchemaID string   `url:"extension_schema_id,omitempty"`
	Include           []string `url:"include,omitempty,brackets"`
}

// ListExtensionsResponse represents a list response of extensions.
type ListExtensionsResponse struct {
	Limit      int          `json:"limit,omitempty"`
	Extensions []*Extension `json:"extensions,omitempty"`
	More       bool         `json:"more,omitempty"`
	Offset     int          `json:"offset,omitempty"`
	Total      int          `json:"total,omitempty"`
}

// ExtensionPayload represents an extension.
type ExtensionPayload struct {
	Extension *Extension `json:"extension"`
}

// List lists existing extensions.
func (s *ExtensionService) List(o *ListExtensionsOptions) (*ListExtensionsResponse, *Response, error) {
	u := "/extensions"
	v := new(ListExtensionsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new extension.
func (s *ExtensionService) Create(extension *Extension) (*Extension, *Response, error) {
	u := "/extensions"
	v := new(ExtensionPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &ExtensionPayload{Extension: extension}, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Extension, resp, nil
}

// Delete removes an existing extension.
func (s *ExtensionService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/extensions/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about an extension.
func (s *ExtensionService) Get(id string) (*Extension, *Response, error) {
	u := fmt.Sprintf("/extensions/%s", id)
	v := new(ExtensionPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Extension, resp, nil
}

// Update updates an existing extension.
func (s *ExtensionService) Update(id string, extension *Extension) (*Extension, *Response, error) {
	u := fmt.Sprintf("/extensions/%s", id)
	v := new(ExtensionPayload)
	resp, err := s.client.newRequestDo("PUT", u, nil, &ExtensionPayload{Extension: extension}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Extension, resp, nil
}
