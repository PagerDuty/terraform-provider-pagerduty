package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CustomFieldDataType is an enumeration of available datatypes for fields.
//
// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
type CustomFieldDataType int64

const (
	CustomFieldDataTypeUnknown CustomFieldDataType = iota
	CustomFieldDataTypeString
	CustomFieldDataTypeInt
	CustomFieldDataTypeFloat
	CustomFieldDataTypeBool
	CustomFieldDataTypeUrl
	CustomFieldDataTypeDateTime
	CustomFieldDataTypeFieldOption
)

// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
func (d CustomFieldDataType) String() string {
	return "unknown"
}

// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
func CustomFieldDataTypeFromString(s string) CustomFieldDataType {
	return CustomFieldDataTypeUnknown
}

// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
func (d CustomFieldDataType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf(`"%v"`, d.String()))
	return buffer.Bytes(), nil
}

// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
func (d *CustomFieldDataType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*d = CustomFieldDataTypeFromString(str)
	return nil
}

// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
func (d *CustomFieldDataType) IsKnown() bool {
	return *d != CustomFieldDataTypeUnknown
}

// IsAllowedOnField determines if the CustomFieldDataType is a legal value for fields. This enables field_option to be a defined datatype
// (as is necessary for default values on field configurations) but not on fields.
//
// Deprecated: Use IncidentCustomFieldDataType for Incident Custom Fields
func (d *CustomFieldDataType) IsAllowedOnField() bool {
	return d.IsKnown() && *d != CustomFieldDataTypeFieldOption
}
