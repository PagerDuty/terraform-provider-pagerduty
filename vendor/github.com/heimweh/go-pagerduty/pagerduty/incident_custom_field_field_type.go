package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// IncidentCustomFieldFieldType is an enumeration of available fieldtypes for fields.
type IncidentCustomFieldFieldType int64

const (
	IncidentCustomFieldFieldTypeUnknown IncidentCustomFieldFieldType = iota
	IncidentCustomFieldFieldTypeSingleValue
	IncidentCustomFieldFieldTypeSingleValueFixed
	IncidentCustomFieldFieldTypeMultiValue
	IncidentCustomFieldFieldTypeMultiValueFixed
)

func (d IncidentCustomFieldFieldType) String() string {
	return incidentCustomFieldFieldTypeToString[d]
}

func IncidentCustomFieldFieldTypeFromString(s string) IncidentCustomFieldFieldType {
	return incidentCustomFieldFieldTypeFromString[s]
}

var incidentCustomFieldFieldTypeToString = map[IncidentCustomFieldFieldType]string{
	IncidentCustomFieldFieldTypeUnknown:          "unknown",
	IncidentCustomFieldFieldTypeSingleValue:      "single_value",
	IncidentCustomFieldFieldTypeSingleValueFixed: "single_value_fixed",
	IncidentCustomFieldFieldTypeMultiValue:       "multi_value",
	IncidentCustomFieldFieldTypeMultiValueFixed:  "multi_value_fixed",
}

var incidentCustomFieldFieldTypeFromString = map[string]IncidentCustomFieldFieldType{
	"unknown":            IncidentCustomFieldFieldTypeUnknown,
	"single_value":       IncidentCustomFieldFieldTypeSingleValue,
	"single_value_fixed": IncidentCustomFieldFieldTypeSingleValueFixed,
	"multi_value":        IncidentCustomFieldFieldTypeMultiValue,
	"multi_value_fixed":  IncidentCustomFieldFieldTypeMultiValueFixed,
}

func (d IncidentCustomFieldFieldType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf(`"%v"`, d.String()))
	return buffer.Bytes(), nil
}

func (d *IncidentCustomFieldFieldType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*d = IncidentCustomFieldFieldTypeFromString(str)
	return nil
}

func (d *IncidentCustomFieldFieldType) IsKnown() bool {
	return *d != IncidentCustomFieldFieldTypeUnknown
}

func (d *IncidentCustomFieldFieldType) IsMultiValue() bool {
	return *d == IncidentCustomFieldFieldTypeMultiValue || *d == IncidentCustomFieldFieldTypeMultiValueFixed
}
