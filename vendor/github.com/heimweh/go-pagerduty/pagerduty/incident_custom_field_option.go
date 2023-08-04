package pagerduty

import (
	"context"
	"fmt"
)

// IncidentCustomFieldOption represents an option for a fixed-value field.
type IncidentCustomFieldOption struct {
	ID   string                         `json:"id,omitempty"`
	Type string                         `json:"type,omitempty"`
	Data *IncidentCustomFieldOptionData `json:"data,omitempty"`
}

// IncidentCustomFieldOptionData represents the value of a IncidentCustomFieldOption
type IncidentCustomFieldOptionData struct {
	DataType IncidentCustomFieldDataType `json:"data_type,omitempty"`
	Value    interface{}                 `json:"value,omitempty"`
}

// IncidentCustomFieldOptionPayload represents payload with a field option object
type IncidentCustomFieldOptionPayload struct {
	FieldOption *IncidentCustomFieldOption `json:"field_option,omitempty"`
}

// ListIncidentCustomFieldOptionsResponse represents a list response of field options
type ListIncidentCustomFieldOptionsResponse struct {
	FieldOptions []*IncidentCustomFieldOption `json:"field_options,omitempty"`
}

// CreateFieldOptionContext creates a new field option.
func (s *IncidentCustomFieldService) CreateFieldOptionContext(ctx context.Context, fieldID string, fieldOption *IncidentCustomFieldOption) (*IncidentCustomFieldOption, *Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s/field_options", fieldID)
	v := new(IncidentCustomFieldOptionPayload)

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, &IncidentCustomFieldOptionPayload{FieldOption: fieldOption}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldOption, resp, nil
}

// UpdateFieldOptionContext updates an existing field option.
func (s *IncidentCustomFieldService) UpdateFieldOptionContext(ctx context.Context, fieldID string, fieldOptionID string, fieldOption *IncidentCustomFieldOption) (*IncidentCustomFieldOption, *Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s/field_options/%s", fieldID, fieldOptionID)
	v := new(IncidentCustomFieldOptionPayload)

	resp, err := s.client.newRequestDoContext(ctx, "PUT", u, nil, &IncidentCustomFieldOptionPayload{FieldOption: fieldOption}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldOption, resp, nil
}

// GetFieldOptionContext gets a field option.
func (s *IncidentCustomFieldService) GetFieldOptionContext(ctx context.Context, fieldID string, fieldOptionID string) (*IncidentCustomFieldOption, *Response, error) {
	l, resp, err := s.ListFieldOptionsContext(ctx, fieldID)
	if err != nil {
		return nil, nil, err
	}

	for _, o := range l.FieldOptions {
		if o.ID == fieldOptionID {
			return o, resp, nil
		}
	}

	return nil, nil, fmt.Errorf("no field option with ID %s under field %s can be found", fieldOptionID, fieldID)
}

// ListFieldOptionsContext lists the field options for a field.
func (s *IncidentCustomFieldService) ListFieldOptionsContext(ctx context.Context, fieldID string) (*ListIncidentCustomFieldOptionsResponse, *Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s/field_options", fieldID)
	v := new(ListIncidentCustomFieldOptionsResponse)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// DeleteFieldOptionContext disables an existing field option.
func (s *IncidentCustomFieldService) DeleteFieldOptionContext(ctx context.Context, fieldID string, fieldOptionID string) (*Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s/field_options/%s", fieldID, fieldOptionID)
	return s.client.newRequestDoContext(ctx, "DELETE", u, nil, nil, nil)
}
