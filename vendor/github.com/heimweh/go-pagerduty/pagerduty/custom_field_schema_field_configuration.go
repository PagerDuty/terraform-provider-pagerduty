package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

// CustomFieldSchemaFieldConfiguration represents a field configuration in a field schema.
//
// Deprecated: This struct should no longer be used
type CustomFieldSchemaFieldConfiguration struct {
	ID           string                   `json:"id,omitempty"`
	Type         string                   `json:"type,omitempty"`
	Required     bool                     `json:"required"`
	Field        *CustomField             `json:"field,omitempty"`
	DefaultValue *CustomFieldDefaultValue `json:"default_value,omitempty"`
}

// CustomFieldDefaultValue
//
// Deprecated: This struct should no longer be used
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
//
// Deprecated: This struct should no longer be used
type ListCustomFieldSchemaConfigurationsOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// GetCustomFieldSchemaConfigurationsOptions represents options when retrieving a field configuration for a schema
//
// Deprecated: This struct should no longer be used
type GetCustomFieldSchemaConfigurationsOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// Deprecated: This struct should no longer be used
type CustomFieldSchemaFieldConfigurationPayload struct {
	FieldConfiguration *CustomFieldSchemaFieldConfiguration `json:"field_configuration,omitempty"`
}

// ListCustomFieldSchemaFieldConfigurationsResponse represents a list response of field configurations for a schema
//
// Deprecated: This struct should no longer be used
type ListCustomFieldSchemaFieldConfigurationsResponse struct {
	FieldConfigurations []*CustomFieldSchemaFieldConfiguration `json:"field_configurations,omitempty"`
}

// ListFieldConfigurations lists field configurations for a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) ListFieldConfigurations(schemaID string, o *ListCustomFieldSchemaConfigurationsOptions) (*ListCustomFieldSchemaFieldConfigurationsResponse, *Response, error) {
	return s.ListFieldConfigurationsContext(context.Background(), schemaID, o)
}

// ListFieldConfigurationsContext lists field configurations for a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) ListFieldConfigurationsContext(_ context.Context, _ string, _ *ListCustomFieldSchemaConfigurationsOptions) (*ListCustomFieldSchemaFieldConfigurationsResponse, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// GetFieldConfiguration gets a field configuration in a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) GetFieldConfiguration(schemaID string, configurationID string, o *GetCustomFieldSchemaConfigurationsOptions) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return s.GetFieldConfigurationContext(context.Background(), schemaID, configurationID, o)
}

// GetFieldConfigurationContext gets a field configuration in a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) GetFieldConfigurationContext(_ context.Context, _ string, _ string, _ *GetCustomFieldSchemaConfigurationsOptions) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// DeleteFieldConfiguration deletes a field configuration for a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) DeleteFieldConfiguration(schemaID string, configurationID string) (*Response, error) {
	return s.DeleteFieldConfigurationContext(context.Background(), schemaID, configurationID)
}

// DeleteFieldConfigurationContext deletes a field configuration for a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) DeleteFieldConfigurationContext(_ context.Context, _ string, _ string) (*Response, error) {
	return nil, customFieldDeprecationError()
}

// CreateFieldConfiguration creates a field configuration in a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) CreateFieldConfiguration(schemaID string, configuration *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return s.CreateFieldConfigurationContext(context.Background(), schemaID, configuration)
}

// CreateFieldConfigurationContext creates a field configuration in a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) CreateFieldConfigurationContext(_ context.Context, _ string, _ *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// UpdateFieldConfiguration updates a field configuration in a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) UpdateFieldConfiguration(schemaID string, configurationID string, configuration *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return s.UpdateFieldConfigurationContext(context.Background(), schemaID, configurationID, configuration)
}

// UpdateFieldConfigurationContext updates a field configuration in a schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) UpdateFieldConfigurationContext(_ context.Context, _ string, _ string, _ *CustomFieldSchemaFieldConfiguration) (*CustomFieldSchemaFieldConfiguration, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}
