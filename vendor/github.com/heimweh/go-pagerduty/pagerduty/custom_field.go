package pagerduty

import (
	"context"
	"errors"
)

const customFieldDeprecationMessage = "standalone custom field functionality has been removed"

func customFieldDeprecationError() error {
	return errors.New(customFieldDeprecationMessage)
}

// CustomFieldService handles the communication with field related methods of the PagerDuty API.
//
// Deprecated: This service should no longer be used. IncidentCustomFieldService provides similar functionality.
type CustomFieldService service

// CustomField represents a custom field.
//
// Deprecated: This struct should no longer be used. IncidentCustomField is similar but not identical.
type CustomField struct {
	ID           string               `json:"id,omitempty"`
	Name         string               `json:"name,omitempty"`
	DisplayName  string               `json:"display_name,omitempty"`
	Type         string               `json:"type,omitempty"`
	Summary      string               `json:"summary,omitempty"`
	Self         string               `json:"self,omitempty"`
	DataType     CustomFieldDataType  `json:"datatype,omitempty"`
	Description  *string              `json:"description,omitempty"`
	MultiValue   bool                 `json:"multi_value"`
	FixedOptions bool                 `json:"fixed_options"`
	FieldOptions []*CustomFieldOption `json:"field_options,omitempty"`
}

// ListCustomFieldResponse represents a list response of fields
//
// Deprecated: This struct should no longer be used.
type ListCustomFieldResponse struct {
	Total  int            `json:"total,omitempty"`
	Fields []*CustomField `json:"fields,omitempty"`
	Offset int            `json:"offset,omitempty"`
	More   bool           `json:"more,omitempty"`
	Limit  int            `json:"limit,omitempty"`
}

// CustomFieldPayload represents payload with a field object
//
// Deprecated: This struct should no longer be used.
type CustomFieldPayload struct {
	Field *CustomField `json:"field,omitempty"`
}

// ListCustomFieldOptions represents options when retrieving a list of fields.
//
// Deprecated: This struct should no longer be used.
type ListCustomFieldOptions struct {
	Offset   int      `url:"offset,omitempty"`
	Limit    int      `url:"limit,omitempty"`
	Total    bool     `url:"total,omitempty"`
	Includes []string `url:"include,brackets,omitempty"`
}

// GetCustomFieldOptions represents options when retrieving a field.
//
// Deprecated: This struct should no longer be used.
type GetCustomFieldOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// List lists existing custom fields. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
//
// Deprecated: Use IncidentCustomFieldService.List
func (s *CustomFieldService) List(o *ListCustomFieldOptions) (*ListCustomFieldResponse, *Response, error) {
	return s.ListContext(context.Background(), o)
}

// ListContext lists existing custom fields. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
//
// Deprecated: Use IncidentCustomFieldService.ListContext
func (s *CustomFieldService) ListContext(_ context.Context, _ *ListCustomFieldOptions) (*ListCustomFieldResponse, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Get gets a custom field.
//
// Deprecated: Use IncidentCustomFieldService.Get
func (s *CustomFieldService) Get(id string, o *GetCustomFieldOptions) (*CustomField, *Response, error) {
	return s.GetContext(context.Background(), id, o)
}

// GetContext gets a custom field.
//
// Deprecated: Use IncidentCustomFieldService.GetContext
func (s *CustomFieldService) GetContext(_ context.Context, _ string, _ *GetCustomFieldOptions) (*CustomField, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Create creates a new custom field.
//
// Deprecated: Use IncidentCustomFieldService.Create
func (s *CustomFieldService) Create(field *CustomField) (*CustomField, *Response, error) {
	return s.CreateContext(context.Background(), field)
}

// CreateContext creates a new custom field.
//
// Deprecated: Use IncidentCustomFieldService.CreateContext
func (s *CustomFieldService) CreateContext(_ context.Context, _ *CustomField) (*CustomField, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Delete removes an existing custom field.
//
// Deprecated: Use IncidentCustomFieldService.Delete
func (s *CustomFieldService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an existing custom field.
//
// Deprecated: Use IncidentCustomFieldService.DeleteContext
func (s *CustomFieldService) DeleteContext(_ context.Context, _ string) (*Response, error) {
	return nil, customFieldDeprecationError()
}

// Update updates an existing custom field.
//
// Deprecated: Use IncidentCustomFieldService.Update
func (s *CustomFieldService) Update(id string, field *CustomField) (*CustomField, *Response, error) {
	return s.UpdateContext(context.Background(), id, field)
}

// UpdateContext updates an existing custom field.
//
// Deprecated: Use IncidentCustomFieldService.UpdateContext
func (s *CustomFieldService) UpdateContext(_ context.Context, _ string, _ *CustomField) (*CustomField, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}
