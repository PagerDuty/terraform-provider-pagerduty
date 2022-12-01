package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// IncidentWorkflowTriggerType is an enumeration of available types for incident workflow triggers.
type IncidentWorkflowTriggerType int64

const (
	IncidentWorkflowTriggerTypeUnknown IncidentWorkflowTriggerType = iota
	IncidentWorkflowTriggerTypeManual
	IncidentWorkflowTriggerTypeConditional
)

func (d IncidentWorkflowTriggerType) String() string {
	return incidentWorkflowTriggerTypeToString[d]
}

func IncidentWorkflowTriggerTypeFromString(s string) IncidentWorkflowTriggerType {
	return incidentWorkflowTriggerTypeFromString[s]
}

var incidentWorkflowTriggerTypeToString = map[IncidentWorkflowTriggerType]string{
	IncidentWorkflowTriggerTypeUnknown:     "unknown",
	IncidentWorkflowTriggerTypeManual:      "manual",
	IncidentWorkflowTriggerTypeConditional: "conditional",
}

var incidentWorkflowTriggerTypeFromString = map[string]IncidentWorkflowTriggerType{
	"unknown":     IncidentWorkflowTriggerTypeUnknown,
	"manual":      IncidentWorkflowTriggerTypeManual,
	"conditional": IncidentWorkflowTriggerTypeConditional,
}

func (t IncidentWorkflowTriggerType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf(`"%v"`, t.String()))
	return buffer.Bytes(), nil
}

func (t *IncidentWorkflowTriggerType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*t = IncidentWorkflowTriggerTypeFromString(str)
	return nil
}

func (t *IncidentWorkflowTriggerType) IsKnown() bool {
	return *t != IncidentWorkflowTriggerTypeUnknown
}
