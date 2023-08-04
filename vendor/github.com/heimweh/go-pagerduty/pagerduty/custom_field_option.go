package pagerduty

import (
	"context"
)

// CustomFieldOption represents an option for a fixed-value field.
//
// Deprecated: This struct should no longer be used. IncidentCustomFieldOption is similar but not identical.
type CustomFieldOption struct {
	ID   string                 `json:"id,omitempty"`
	Type string                 `json:"type,omitempty"`
	Data *CustomFieldOptionData `json:"data,omitempty"`
}

// CustomFieldOptionData represents the value of a CustomFieldOption
//
// Deprecated: This struct should no longer be used. IncidentCustomFieldOptionData is similar but not identical.
type CustomFieldOptionData struct {
	DataType CustomFieldDataType `json:"datatype,omitempty"`
	Value    interface{}         `json:"value,omitempty"`
}

// CustomFieldOptionPayload represents payload with a field option object
//
// Deprecated: This struct should no longer be used.
type CustomFieldOptionPayload struct {
	FieldOption *CustomFieldOption `json:"field_option,omitempty"`
}

// ListCustomFieldOptionsResponse represents a list response of field options
//
// Deprecated: This struct should no longer be used.
type ListCustomFieldOptionsResponse struct {
	FieldOptions []*CustomFieldOption `json:"field_options,omitempty"`
}

// CreateFieldOption creates a new field option.
//
// Deprecated: Use IncidentCustomFieldService.CreateFieldOption
func (s *CustomFieldService) CreateFieldOption(fieldID string, fieldOption *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	return s.CreateFieldOptionContext(context.Background(), fieldID, fieldOption)
}

// CreateFieldOptionContext creates a new field option.
//
// Deprecated: Use IncidentCustomFieldService.CreateFieldOptionContext
func (s *CustomFieldService) CreateFieldOptionContext(_ context.Context, _ string, _ *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// UpdateFieldOption updates an existing field option.
//
// Deprecated: Use IncidentCustomFieldService.UpdateFieldOption
func (s *CustomFieldService) UpdateFieldOption(fieldID string, fieldOptionID string, fieldOption *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	return s.UpdateFieldOptionContext(context.Background(), fieldID, fieldOptionID, fieldOption)
}

// UpdateFieldOptionContext updates an existing field option.
//
// Deprecated: Use IncidentCustomFieldService.UpdateFieldOptionContext
func (s *CustomFieldService) UpdateFieldOptionContext(_ context.Context, _ string, _ string, _ *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// GetFieldOption gets a field option.
//
// Deprecated: Use IncidentCustomFieldService.GetFieldOption
func (s *CustomFieldService) GetFieldOption(fieldID string, fieldOptionID string) (*CustomFieldOption, *Response, error) {
	return s.GetFieldOptionContext(context.Background(), fieldID, fieldOptionID)
}

// GetFieldOptionContext gets a field option.
//
// Deprecated: Use IncidentCustomFieldService.GetFieldOptionContext
func (s *CustomFieldService) GetFieldOptionContext(_ context.Context, _ string, _ string) (*CustomFieldOption, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// ListFieldOptions lists the field options for a field.
//
// Deprecated: Use IncidentCustomFieldService.ListFieldOptions
func (s *CustomFieldService) ListFieldOptions(fieldID string) (*ListCustomFieldOptionsResponse, *Response, error) {
	return s.ListFieldOptionsContext(context.Background(), fieldID)
}

// ListFieldOptionsContext lists the field options for a field.
//
// Deprecated: Use IncidentCustomFieldService.ListFieldOptionsContext
func (s *CustomFieldService) ListFieldOptionsContext(_ context.Context, _ string) (*ListCustomFieldOptionsResponse, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// DeleteFieldOption deletes an existing field option.
//
// Deprecated: Use IncidentCustomFieldService.DeleteFieldOption
func (s *CustomFieldService) DeleteFieldOption(fieldID string, fieldOptionID string) (*Response, error) {
	return s.DeleteFieldOptionContext(context.Background(), fieldID, fieldOptionID)
}

// DeleteFieldOptionContext disables an existing field option.
//
// Deprecated: Use IncidentCustomFieldService.DeleteFieldOptionContext
func (s *CustomFieldService) DeleteFieldOptionContext(_ context.Context, _ string, _ string) (*Response, error) {
	return nil, customFieldDeprecationError()
}
