package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

// IncidentCustomFieldService handles the communication with custom field on incidents related methods of the PagerDuty API.
type IncidentCustomFieldService service

// IncidentCustomField represents an incident custom field.
type IncidentCustomField struct {
	ID           string                       `json:"id,omitempty"`
	Name         string                       `json:"name,omitempty"`
	DisplayName  string                       `json:"display_name,omitempty"`
	Type         string                       `json:"type,omitempty"`
	Summary      string                       `json:"summary,omitempty"`
	Self         string                       `json:"self,omitempty"`
	DataType     IncidentCustomFieldDataType  `json:"data_type,omitempty"`
	FieldType    IncidentCustomFieldFieldType `json:"field_type,omitempty"`
	Description  *string                      `json:"description,omitempty"`
	DefaultValue interface{}                  `json:"default_value,omitempty"`
	FieldOptions []*IncidentCustomFieldOption `json:"field_options,omitempty"`
}

type rawIncidentCustomField struct {
	ID           string                       `json:"id,omitempty"`
	Name         string                       `json:"name,omitempty"`
	DisplayName  string                       `json:"display_name,omitempty"`
	Type         string                       `json:"type,omitempty"`
	Summary      string                       `json:"summary,omitempty"`
	Self         string                       `json:"self,omitempty"`
	DataType     IncidentCustomFieldDataType  `json:"data_type,omitempty"`
	FieldType    IncidentCustomFieldFieldType `json:"field_type,omitempty"`
	Description  *string                      `json:"description,omitempty"`
	DefaultValue interface{}                  `json:"default_value,omitempty"`
	FieldOptions []*IncidentCustomFieldOption `json:"field_options,omitempty"`
}

func (d *IncidentCustomField) UnmarshalJSON(data []byte) error {
	var p rawIncidentCustomField
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	*d = IncidentCustomField{
		ID:           p.ID,
		Name:         p.Name,
		DisplayName:  p.DisplayName,
		Type:         p.Type,
		Summary:      p.Summary,
		Self:         p.Self,
		DataType:     p.DataType,
		FieldType:    p.FieldType,
		Description:  p.Description,
		FieldOptions: p.FieldOptions,
	}
	if p.DefaultValue != nil {
		switch p.DataType {
		case IncidentCustomFieldDataTypeInt:
			err := d.convertForInt(p.DefaultValue)
			if err != nil {
				return err
			}
		default:
			d.DefaultValue = p.DefaultValue
		}
	}
	return nil
}

func (d *IncidentCustomField) convertForInt(value interface{}) error {
	switch v := value.(type) {
	case []interface{}:
		if d.FieldType.IsMultiValue() {
			var s []interface{}
			for _, f := range v {
				switch ev := f.(type) {
				case float64:
					s = append(s, int64(math.Round(ev)))
				default:
					return fmt.Errorf("received unexpected %T as an element in a multi-value int", ev)
				}
			}
			d.DefaultValue = s
			return nil
		} else {
			return fmt.Errorf("received unexpected %T for non-multi-value int", v)
		}
	case float64:
		if d.FieldType.IsMultiValue() {
			return fmt.Errorf("received unexpected %T for multi-value int", v)
		} else {
			d.DefaultValue = int64(math.Round(v))
			return nil
		}
	default:
		return fmt.Errorf("received unexpected %T as for an integer default value", v)
	}
}

// ListIncidentCustomFieldResponse represents a list response of fields
type ListIncidentCustomFieldResponse struct {
	Fields []*IncidentCustomField `json:"fields,omitempty"`
}

// IncidentCustomFieldPayload represents payload with a field object
type IncidentCustomFieldPayload struct {
	Field *IncidentCustomField `json:"field,omitempty"`
}

// ListIncidentCustomFieldOptions represents options when retrieving a list of fields.
type ListIncidentCustomFieldOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// GetIncidentCustomFieldOptions represents options when retrieving a field.
type GetIncidentCustomFieldOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// ListContext lists existing custom fields. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
func (s *IncidentCustomFieldService) ListContext(ctx context.Context, o *ListIncidentCustomFieldOptions) (*ListIncidentCustomFieldResponse, *Response, error) {
	u := "/incidents/custom_fields"
	v := new(ListIncidentCustomFieldResponse)

	if o == nil {
		o = &ListIncidentCustomFieldOptions{}
	}

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// GetContext gets a custom field.
func (s *IncidentCustomFieldService) GetContext(ctx context.Context, id string, o *GetIncidentCustomFieldOptions) (*IncidentCustomField, *Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s", id)
	v := new(IncidentCustomFieldPayload)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Field, resp, nil
}

// CreateContext creates a new custom field.
func (s *IncidentCustomFieldService) CreateContext(ctx context.Context, field *IncidentCustomField) (*IncidentCustomField, *Response, error) {
	u := "/incidents/custom_fields"
	v := new(IncidentCustomFieldPayload)

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, &IncidentCustomFieldPayload{Field: field}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Field, resp, nil
}

// DeleteContext removes an existing custom field.
func (s *IncidentCustomFieldService) DeleteContext(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s", id)
	return s.client.newRequestDoContext(ctx, "DELETE", u, nil, nil, nil)
}

// UpdateContext updates an existing custom field.
func (s *IncidentCustomFieldService) UpdateContext(ctx context.Context, id string, field *IncidentCustomField) (*IncidentCustomField, *Response, error) {
	u := fmt.Sprintf("/incidents/custom_fields/%s", id)
	v := new(IncidentCustomFieldPayload)

	resp, err := s.client.newRequestDoContext(ctx, "PUT", u, nil, &IncidentCustomFieldPayload{Field: field}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Field, resp, nil
}
