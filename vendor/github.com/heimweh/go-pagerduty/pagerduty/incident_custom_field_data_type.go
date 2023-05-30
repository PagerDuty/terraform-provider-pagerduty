package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// IncidentCustomFieldDataType is an enumeration of available datatypes for fields.
type IncidentCustomFieldDataType int64

const (
	IncidentCustomFieldDataTypeUnknown IncidentCustomFieldDataType = iota
	IncidentCustomFieldDataTypeString
	IncidentCustomFieldDataTypeInt
	IncidentCustomFieldDataTypeFloat
	IncidentCustomFieldDataTypeBool
	IncidentCustomFieldDataTypeUrl
	IncidentCustomFieldDataTypeDateTime
)

func (d IncidentCustomFieldDataType) String() string {
	return incidentCustomFieldDataTypeToString[d]
}

func IncidentCustomFieldDataTypeFromString(s string) IncidentCustomFieldDataType {
	return incidentCustomFieldDataTypeFromString[s]
}

var incidentCustomFieldDataTypeToString = map[IncidentCustomFieldDataType]string{
	IncidentCustomFieldDataTypeUnknown:  "unknown",
	IncidentCustomFieldDataTypeString:   "string",
	IncidentCustomFieldDataTypeInt:      "integer",
	IncidentCustomFieldDataTypeFloat:    "float",
	IncidentCustomFieldDataTypeBool:     "boolean",
	IncidentCustomFieldDataTypeUrl:      "url",
	IncidentCustomFieldDataTypeDateTime: "datetime",
}

var incidentCustomFieldDataTypeFromString = map[string]IncidentCustomFieldDataType{
	"unknown":  IncidentCustomFieldDataTypeUnknown,
	"string":   IncidentCustomFieldDataTypeString,
	"integer":  IncidentCustomFieldDataTypeInt,
	"float":    IncidentCustomFieldDataTypeFloat,
	"boolean":  IncidentCustomFieldDataTypeBool,
	"url":      IncidentCustomFieldDataTypeUrl,
	"datetime": IncidentCustomFieldDataTypeDateTime,
}

func (d IncidentCustomFieldDataType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf(`"%v"`, d.String()))
	return buffer.Bytes(), nil
}

func (d *IncidentCustomFieldDataType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*d = IncidentCustomFieldDataTypeFromString(str)
	return nil
}

func (d *IncidentCustomFieldDataType) IsKnown() bool {
	return *d != IncidentCustomFieldDataTypeUnknown
}
