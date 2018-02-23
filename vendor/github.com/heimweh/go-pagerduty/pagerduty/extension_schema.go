package pagerduty

import "fmt"

// ExtensionSchemaService handles the communication with extension schemas related methods
// of the PagerDuty API.
type ExtensionSchemaService service

// ExtensionSchema represents an extension schema.
type ExtensionSchema struct {
	ExtensionSchema *ExtensionSchema `json:"extension_schema,omitempty"`
	Description     string           `json:"description,omitempty"`
	GuideURL        string           `json:"guide_url,omitempty"`
	HTMLURL         string           `json:"html_url,omitempty"`
	IconURL         string           `json:"icon_url,omitempty"`
	ID              string           `json:"id,omitempty"`
	Key             string           `json:"key,omitempty"`
	Label           string           `json:"label,omitempty"`
	LogoURL         string           `json:"logo_url,omitempty"`
	Self            string           `json:"self,omitempty"`
	SendTypes       []string         `json:"send_types,omitempty"`
	Summary         string           `json:"summary,omitempty"`
	Type            string           `json:"type,omitempty"`
	URL             string           `json:"url,omitempty"`
}

// ListExtensionSchemasResponse represents a list response of extension schemas.
type ListExtensionSchemasResponse struct {
	ExtensionSchemas []*ExtensionSchema `json:"extension_schemas,omitempty"`
	Limit            int                `json:"limit,omitempty"`
	More             bool               `json:"more,omitempty"`
	Offset           int                `json:"offset,omitempty"`
	Total            int                `json:"total,omitempty"`
}

// List lists extension schemas.
func (s *ExtensionSchemaService) List() (*ListExtensionSchemasResponse, *Response, error) {
	u := "/extension_schemas"
	v := new(ListExtensionSchemasResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Get retrieves information about an extension schema.
func (s *ExtensionSchemaService) Get(id string) (*ExtensionSchema, *Response, error) {
	u := fmt.Sprintf("/extension_schemas/%s", id)
	v := new(ExtensionSchema)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.ExtensionSchema, resp, nil
}
