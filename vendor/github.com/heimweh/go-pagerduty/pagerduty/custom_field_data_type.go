package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CustomFieldDataType is an enumeration of available datatypes for fields.
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

func (d CustomFieldDataType) String() string {
	return customFieldDataTypeToString[d]
}

func CustomFieldDataTypeFromString(s string) CustomFieldDataType {
	return customFieldDataTypeFromString[s]
}

var customFieldDataTypeToString = map[CustomFieldDataType]string{
	CustomFieldDataTypeUnknown:     "unknown",
	CustomFieldDataTypeString:      "string",
	CustomFieldDataTypeInt:         "integer",
	CustomFieldDataTypeFloat:       "float",
	CustomFieldDataTypeBool:        "boolean",
	CustomFieldDataTypeUrl:         "url",
	CustomFieldDataTypeDateTime:    "datetime",
	CustomFieldDataTypeFieldOption: "field_option",
}

var customFieldDataTypeFromString = map[string]CustomFieldDataType{
	"unknown":      CustomFieldDataTypeUnknown,
	"string":       CustomFieldDataTypeString,
	"integer":      CustomFieldDataTypeInt,
	"float":        CustomFieldDataTypeFloat,
	"boolean":      CustomFieldDataTypeBool,
	"url":          CustomFieldDataTypeUrl,
	"datetime":     CustomFieldDataTypeDateTime,
	"field_option": CustomFieldDataTypeFieldOption,
}

func (d CustomFieldDataType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf(`"%v"`, d.String()))
	return buffer.Bytes(), nil
}

func (d *CustomFieldDataType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*d = CustomFieldDataTypeFromString(str)
	return nil
}

func (d *CustomFieldDataType) IsKnown() bool {
	return *d != CustomFieldDataTypeUnknown
}

// IsAllowedOnField determines if the CustomFieldDataType is a legal value for fields. This enables field_option to be a defined datatype
// (as is necessary for default values on field configurations) but not on fields.
func (d *CustomFieldDataType) IsAllowedOnField() bool {
	return d.IsKnown() && *d != CustomFieldDataTypeFieldOption
}
