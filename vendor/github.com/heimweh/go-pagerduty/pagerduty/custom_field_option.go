package pagerduty

import (
	"context"
	"fmt"
)

// CustomFieldOption represents an option for a fixed-value field.
type CustomFieldOption struct {
	ID   string                 `json:"id,omitempty"`
	Type string                 `json:"type,omitempty"`
	Data *CustomFieldOptionData `json:"data,omitempty"`
}

// CustomFieldOptionData represents the value of a CustomFieldOption
type CustomFieldOptionData struct {
	DataType CustomFieldDataType `json:"datatype,omitempty"`
	Value    interface{}         `json:"value,omitempty"`
}

// CustomFieldOptionPayload represents payload with a field option object
type CustomFieldOptionPayload struct {
	FieldOption *CustomFieldOption `json:"field_option,omitempty"`
}

// ListCustomFieldOptionsResponse represents a list response of field options
type ListCustomFieldOptionsResponse struct {
	FieldOptions []*CustomFieldOption `json:"field_options,omitempty"`
}

// CreateFieldOption creates a new field option.
func (s *CustomFieldService) CreateFieldOption(fieldID string, fieldOption *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	return s.CreateFieldOptionContext(context.Background(), fieldID, fieldOption)
}

// CreateFieldOptionContext creates a new field option.
func (s *CustomFieldService) CreateFieldOptionContext(ctx context.Context, fieldID string, fieldOption *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s/field_options", fieldID)
	v := new(CustomFieldOptionPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "POST", u, nil, &CustomFieldOptionPayload{FieldOption: fieldOption}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldOption, resp, nil
}

// UpdateFieldOption updates an existing field option.
func (s *CustomFieldService) UpdateFieldOption(fieldID string, fieldOptionID string, fieldOption *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	return s.UpdateFieldOptionContext(context.Background(), fieldID, fieldOptionID, fieldOption)
}

// UpdateFieldOptionContext updates an existing field option.
func (s *CustomFieldService) UpdateFieldOptionContext(ctx context.Context, fieldID string, fieldOptionID string, fieldOption *CustomFieldOption) (*CustomFieldOption, *Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s/field_options/%s", fieldID, fieldOptionID)
	v := new(CustomFieldOptionPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "PUT", u, nil, &CustomFieldOptionPayload{FieldOption: fieldOption}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldOption, resp, nil
}

// GetFieldOption gets a field option.
func (s *CustomFieldService) GetFieldOption(fieldID string, fieldOptionID string) (*CustomFieldOption, *Response, error) {
	return s.GetFieldOptionContext(context.Background(), fieldID, fieldOptionID)
}

// GetFieldOptionContext gets a field option.
func (s *CustomFieldService) GetFieldOptionContext(ctx context.Context, fieldID string, fieldOptionID string) (*CustomFieldOption, *Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s/field_options/%s", fieldID, fieldOptionID)
	v := new(CustomFieldOptionPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, nil, nil, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldOption, resp, nil
}

// ListFieldOptions lists the field options for a field.
func (s *CustomFieldService) ListFieldOptions(fieldID string) (*ListCustomFieldOptionsResponse, *Response, error) {
	return s.ListFieldOptionsContext(context.Background(), fieldID)
}

// ListFieldOptionsContext lists the field options for a field.
func (s *CustomFieldService) ListFieldOptionsContext(ctx context.Context, fieldID string) (*ListCustomFieldOptionsResponse, *Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s/field_options", fieldID)
	v := new(ListCustomFieldOptionsResponse)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, nil, nil, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// DeleteFieldOption deletes an existing field option.
func (s *CustomFieldService) DeleteFieldOption(fieldID string, fieldOptionID string) (*Response, error) {
	return s.DeleteFieldOptionContext(context.Background(), fieldID, fieldOptionID)
}

// DeleteFieldOptionContext disables an existing field option.
func (s *CustomFieldService) DeleteFieldOptionContext(ctx context.Context, fieldID string, fieldOptionID string) (*Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s/field_options/%s", fieldID, fieldOptionID)
	return s.client.newRequestDoOptionsContext(ctx, "DELETE", u, nil, nil, nil, customFieldsEarlyAccessHeader)
}
