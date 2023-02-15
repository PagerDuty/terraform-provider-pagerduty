package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

type CustomFieldSchemaFieldConfiguration struct {
	ID           string                   `json:"id,omitempty"`
	Type         string                   `json:"type,omitempty"`
	Required     bool                     `json:"required"`
	Field        *CustomField             `json:"field,omitempty"`
	DefaultValue *CustomFieldDefaultValue `json:"default_value,omitempty"`
}

type CustomFieldDefaultValue struct {
	DataType   CustomFieldDataType `json:"datatype,omitempty"`
	MultiValue bool                `json:"multi_value"`
	Value      interface{}         `json:"value,omitempty"`
}

type rawCustomFieldDefaultValue struct {
	DataType   CustomFieldDataType `json:"datatype,omitempty"`
	MultiValue bool                `json:"multi_value"`
	Value      interface{}         `json:"value,omitempty"`
}

func (d *CustomFieldDefaultValue) UnmarshalJSON(data []byte) error {
	var p rawCustomFieldDefaultValue
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	*d = CustomFieldDefaultValue{
		DataType:   p.DataType,
		MultiValue: p.MultiValue,
	}
	if p.Value != nil {
		switch p.DataType {
		case CustomFieldDataTypeInt:
			err := d.convertForInt(p.Value)
			if err != nil {
				return err
			}
		case CustomFieldDataTypeFieldOption:
			m := p.Value.(map[string]interface{})
			d.Value = m["id"]
		default:
			d.Value = p.Value
		}
	}
	return nil
}

func (d *CustomFieldDefaultValue) MarshalJSON() ([]byte, error) {
	var nd rawCustomFieldDefaultValue
	switch d.DataType {
	case CustomFieldDataTypeFieldOption:
		nd = rawCustomFieldDefaultValue{
			DataType:   d.DataType,
			MultiValue: d.MultiValue,
			Value:      map[string]string{"type": "field_option_reference", "id": d.Value.(string)},
		}
	default:
		nd = rawCustomFieldDefaultValue{
			DataType:   d.DataType,
			MultiValue: d.MultiValue,
			Value:      d.Value,
		}
	}
	return json.Marshal(nd)

}

func (d *CustomFieldDefaultValue) convertForInt(value interface{}) error {
	switch v := value.(type) {
	case []interface{}:
		if d.MultiValue {
			var s []interface{}
			for _, f := range v {
				switch ev := f.(type) {
				case float64:
					s = append(s, int64(math.Round(ev)))
				default:
					return fmt.Errorf("received unexpected %T as an element in a multi-value int", ev)
				}
			}
			d.Value = s
			return nil
		} else {
			return fmt.Errorf("Received unexpected %T for non-multi-value int", v)
		}
	case float64:
		if d.MultiValue {
			return fmt.Errorf("Received unexpected %T for multi-value int", v)
		} else {
			d.Value = int64(math.Round(v))
			return nil
		}
	default:
		return fmt.Errorf("received unexpected %T as for an integer default value", v)
	}
}

// ListCustomFieldSchemaConfigurationsOptions represents options when retrieving a list of field schemas.
type ListCustomFieldSchemaConfigurationsOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// GetCustomFieldSchemaConfigurationsOptions represents options when retrieving a field configuration for a schema
type GetCustomFieldSchemaConfigurationsOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

type CustomFieldSchemaFieldConfigurationPayload struct {
	FieldConfiguration *CustomFieldSchemaFieldConfiguration `json:"field_configuration,omitempty"`
}

// ListCustomFieldSchemaFieldConfigurationsResponse represents a list response of field configurations for a schema
type ListCustomFieldSchemaFieldConfigurationsResponse struct {
	FieldConfigurations []*CustomFieldSchemaFieldConfiguration `json:"field_configurations,omitempty"`
}

// ListFieldConfigurations lists field configurations for a schema.
func (s *CustomFieldSchemaService) ListFieldConfigurations(schemaID string, o *ListCustomFieldSchemaConfigurationsOptions) (*ListCustomFieldSchemaFieldConfigurationsResponse, *Response, error) {
	return s.ListFieldConfigurationsContext(context.Background(), schemaID, o)
}

// ListFieldConfigurationsContext lists field configurations for a schema.
func (s *CustomFieldSchemaService) ListFieldConfigurationsContext(ctx context.Context, schemaID string, o *ListCustomFieldSchemaConfigurationsOptions) (*ListCustomFieldSchemaFieldConfigurationsResponse, *Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s/field_configurations", schemaID)
	v := new(ListCustomFieldSchemaFieldConfigurationsResponse)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, o, nil, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// GetFieldConfiguration gets a field configuration in a schema.
func (s *CustomFieldSchemaService) GetFieldConfiguration(schemaID string, configurationID string, o *GetCustomFieldSchemaConfigurationsOptions) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return s.GetFieldConfigurationContext(context.Background(), schemaID, configurationID, o)
}

// GetFieldConfigurationContext gets a field configuration in a schema.
func (s *CustomFieldSchemaService) GetFieldConfigurationContext(ctx context.Context, schemaID string, configurationID string, o *GetCustomFieldSchemaConfigurationsOptions) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s/field_configurations/%s", schemaID, configurationID)
	v := new(CustomFieldSchemaFieldConfigurationPayload)
	p := &CustomFieldSchemaFieldConfigurationPayload{}

	resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, o, p, v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldConfiguration, resp, nil
}

// DeleteFieldConfiguration deletes a field configuration for a schema.
func (s *CustomFieldSchemaService) DeleteFieldConfiguration(schemaID string, configurationID string) (*Response, error) {
	return s.DeleteFieldConfigurationContext(context.Background(), schemaID, configurationID)
}

// DeleteFieldConfigurationContext deletes a field configuration for a schema.
func (s *CustomFieldSchemaService) DeleteFieldConfigurationContext(ctx context.Context, schemaID string, configurationID string) (*Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s/field_configurations/%s", schemaID, configurationID)
	return s.client.newRequestDoOptionsContext(ctx, "DELETE", u, nil, nil, nil, customFieldsEarlyAccessHeader)
}

// CreateFieldConfiguration creates a field configuration in a schema.
func (s *CustomFieldSchemaService) CreateFieldConfiguration(schemaID string, configuration *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return s.CreateFieldConfigurationContext(context.Background(), schemaID, configuration)
}

// CreateFieldConfigurationContext creates a field configuration in a schema.
func (s *CustomFieldSchemaService) CreateFieldConfigurationContext(ctx context.Context, schemaID string, configuration *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s/field_configurations", schemaID)
	v := new(CustomFieldSchemaFieldConfigurationPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "POST", u, nil, &CustomFieldSchemaFieldConfigurationPayload{FieldConfiguration: configuration}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldConfiguration, resp, nil
}

// UpdateFieldConfiguration updates a field configuration in a schema.
func (s *CustomFieldSchemaService) UpdateFieldConfiguration(schemaID string, configurationID string, configuration *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return s.UpdateFieldConfigurationContext(context.Background(), schemaID, configurationID, configuration)
}

// UpdateFieldConfigurationContext updates a field configuration in a schema.
func (s *CustomFieldSchemaService) UpdateFieldConfigurationContext(ctx context.Context, schemaID string, configurationID string, configuration *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s/field_configurations/%s", schemaID, configurationID)
	v := new(CustomFieldSchemaFieldConfigurationPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "PUT", u, nil, CustomFieldSchemaFieldConfigurationPayload{FieldConfiguration: configuration}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.FieldConfiguration, resp, nil
}
