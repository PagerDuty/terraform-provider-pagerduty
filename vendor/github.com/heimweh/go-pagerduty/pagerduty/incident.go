package pagerduty

import (
	"fmt"
	"log"
)

// IncidentService handles the communication with incident
// related methods of the PagerDuty API.
type IncidentService service

// Incident represents a incident.
type Incident struct {
	ID                   string                      `json:"id,omitempty"`
	Type                 string                      `json:"type,omitempty"`
	Summary              string                      `json:"summary,omitempty"`
	Self                 string                      `json:"self,omitempty"`
	HTMLURL              string                      `json:"html_url,omitempty"`
	IncidentNumber       int                         `json:"incident_number,omitempty"`
	CreatedAt            string                      `json:"created_at,omitempty"`
	Status               string                      `json:"status,omitempty"`
	Title                string                      `json:"title,omitempty"`
	Resolution           string                      `json:"resolution,omitempty"`
	AlertCounts          *AlertCounts                `json:"alert_counts,omitempty"`
	PendingActions       []*PendingAction            `json:"pending_actions,omitempty"`
	IncidentKey          string                      `json:"incident_key,omitempty"`
	Service              *ServiceReference           `json:"service,omitempty"`
	AssignedVia          string                      `json:"assigned_via,omitempty"`
	Assignments          []*IncidentAssignment       `json:"assignments,omitempty"`
	Acknowledgements     []*IncidentAcknowledgement  `json:"acknowledgements,omitempty"`
	LastStatusChangeAt   string                      `json:"last_status_change_at,omitempty"`
	LastStatusChangeBy   *IncidentAttributeReference `json:"last_status_change_by,omitempty"`
	FirstTriggerLogEntry *IncidentAttributeReference `json:"first_trigger_log_entry,omitempty"`
	EscalationPolicy     *EscalationPolicyReference  `json:"escalation_policy,omitempty"`
	Teams                []*TeamReference            `json:"teams,omitempty"`
	Urgency              string                      `json:"urgency,omitempty"`
}

type AlertCounts struct {
	All       int `json:"all"`
	Resolved  int `json:"resolved"`
	Triggered int `json:"triggered"`
}

type PendingAction struct {
	At   string `json:"at"`
	Type string `json:"type"`
}

type IncidentAssignment struct {
	At       string        `json:"at"`
	Assignee UserReference `json:"assignee"`
}

type IncidentAcknowledgement struct {
	At           string                     `json:"at"`
	Acknowledger IncidentAttributeReference `json:"acknowledger"`
}

// IncidentPayload represents an incident.
type IncidentPayload struct {
	Incident *Incident `json:"incident,omitempty"`
}

// ManageIncidentsPayload represents a payload with a list of incidents data.
type ManageIncidentsPayload struct {
	Incidents []*Incident `json:"incidents,omitempty"`
}

// ListIncidentsOptions represents options when listing incidents.
type ListIncidentsOptions struct {
	Limit       int      `url:"limit,omitempty"`
	Offset      int      `url:"offset,omitempty"`
	Total       int      `url:"total,omitempty"`
	DateRange   string   `url:"date_range,omitempty"`
	IncidentKey string   `url:"incident_key,omitempty"`
	Include     []string `url:"include,omitempty,brackets"`
	ServiceIDs  []string `url:"service_ids,omitempty,brackets"`
	Since       string   `url:"since,omitempty"`
	SortBy      []string `url:"sort_by,omitempty,brackets"`
	Statuses    []string `url:"statuses,omitempty,brackets"`
	TeamIDs     []string `url:"team_ids,omitempty,brackets"`
	TimeZone    string   `url:"time_zone,omitempty"`
	Until       string   `url:"until,omitempty"`
	Urgencies   []string `url:"urgencies,omitempty,brackets"`
	UserIDs     []string `url:"user_ids,omitempty,brackets"`
}

// ManageIncidentsOptions represents options when listing incidents.
type ManageIncidentsOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
	Total  int `url:"total,omitempty"`
}

// ListIncidentsResponse represents a list response of incidents.
type ListIncidentsResponse struct {
	Limit     int         `json:"limit,omitempty"`
	More      bool        `json:"more,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	Total     int         `json:"total,omitempty"`
	Incidents []*Incident `json:"incidents,omitempty"`
}

type ManageIncidentsResponse ListIncidentsResponse

// List lists existing incidents.
func (s *IncidentService) List(o *ListIncidentsOptions) (*ListIncidentsResponse, *Response, error) {
	u := "/incidents"
	v := new(ListIncidentsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// ListAll lists all result pages for incidents list.
func (s *IncidentService) ListAll(o *ListIncidentsOptions) ([]*Incident, error) {
	var incidents = make([]*Incident, 0, 25)
	more := true
	offset := 0

	for more {
		log.Printf("==== Getting incidents at offset %d", offset)
		v := new(ListIncidentsResponse)
		_, err := s.client.newRequestDo("GET", "/incidents", o, nil, &v)
		if err != nil {
			return incidents, err
		}
		incidents = append(incidents, v.Incidents...)
		more = v.More
		offset += v.Limit
		o.Offset = offset
	}
	return incidents, nil
}

// ManageIncidents updates existing incidents.
func (s *IncidentService) ManageIncidents(incidents []*Incident, o *ManageIncidentsOptions) (*ManageIncidentsResponse, *Response, error) {
	u := "/incidents"
	v := new(ManageIncidentsResponse)

	resp, err := s.client.newRequestDo("PUT", u, o, &ManageIncidentsPayload{Incidents: incidents}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create an incident
func (s *IncidentService) Create(incident *Incident) (*Incident, *Response, error) {
	u := "/incidents"
	v := new(IncidentPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &IncidentPayload{Incident: incident}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Incident, resp, nil
}

// Get retrieves information about an incident.
func (s *IncidentService) Get(id string) (*Incident, *Response, error) {
	u := fmt.Sprintf("/incidents/%s", id)
	v := new(IncidentPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Incident, resp, nil
}
