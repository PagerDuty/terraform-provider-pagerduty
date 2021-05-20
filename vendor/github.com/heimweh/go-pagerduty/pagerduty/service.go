package pagerduty

import "fmt"

// ServicesService handles the communication with service
// related methods of the PagerDuty API.
type ServicesService service

// At represents when a scheduled action will occur.
type At struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// ScheduledAction contains scheduled actions for the service.
type ScheduledAction struct {
	At        *At    `json:"at,omitempty"`
	ToUrgency string `json:"to_urgency,omitempty"`
	Type      string `json:"type,omitempty"`
}

// IncidentUrgencyType are the incidents urgency during or outside support hours.
type IncidentUrgencyType struct {
	Type    string `json:"type,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

// SupportHours are the support hours for the service.
type SupportHours struct {
	DaysOfWeek []int  `json:"days_of_week,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	TimeZone   string `json:"time_zone,omitempty"`
	Type       string `json:"type,omitempty"`
}

// IncidentUrgencyRule is the default urgency for new incidents.
type IncidentUrgencyRule struct {
	DuringSupportHours  *IncidentUrgencyType `json:"during_support_hours,omitempty"`
	OutsideSupportHours *IncidentUrgencyType `json:"outside_support_hours,omitempty"`
	Type                string               `json:"type,omitempty"`
	Urgency             string               `json:"urgency,omitempty"`
}

// Integration represents a service integration.
type Integration struct {
	CreatedAt        string            `json:"created_at,omitempty"`
	HTMLURL          string            `json:"html_url,omitempty"`
	ID               string            `json:"id,omitempty"`
	Integration      *Integration      `json:"integration,omitempty"`
	IntegrationEmail string            `json:"integration_email,omitempty"`
	IntegrationKey   string            `json:"integration_key,omitempty"`
	Name             string            `json:"name,omitempty"`
	Self             string            `json:"self,omitempty"`
	Service          *ServiceReference `json:"service,omitempty"`
	Summary          string            `json:"summary,omitempty"`
	Type             string            `json:"type,omitempty"`
	Vendor           *VendorReference  `json:"vendor,omitempty"`
}

// Service represents a service.
type Service struct {
	AcknowledgementTimeout *int                       `json:"acknowledgement_timeout"`
	Addons                 []*AddonReference          `json:"addons,omitempty"`
	AlertCreation          string                     `json:"alert_creation,omitempty"`
	AlertGrouping          *string                    `json:"alert_grouping"`
	AlertGroupingTimeout   *int                       `json:"alert_grouping_timeout,omitempty"`
	AutoResolveTimeout     *int                       `json:"auto_resolve_timeout"`
	CreatedAt              string                     `json:"created_at,omitempty"`
	Description            string                     `json:"description,omitempty"`
	EscalationPolicy       *EscalationPolicyReference `json:"escalation_policy,omitempty"`
	HTMLURL                string                     `json:"html_url,omitempty"`
	ID                     string                     `json:"id,omitempty"`
	IncidentUrgencyRule    *IncidentUrgencyRule       `json:"incident_urgency_rule,omitempty"`
	Integrations           []*IntegrationReference    `json:"integrations,omitempty"`
	LastIncidentTimestamp  string                     `json:"last_incident_timestamp,omitempty"`
	Name                   string                     `json:"name,omitempty"`
	ScheduledActions       []*ScheduledAction         `json:"scheduled_actions,omitempty"`
	Self                   string                     `json:"self,omitempty"`
	Service                *Service                   `json:"service,omitempty"`
	Status                 string                     `json:"status,omitempty"`
	Summary                string                     `json:"summary,omitempty"`
	SupportHours           *SupportHours              `json:"support_hours,omitempty"`
	Teams                  []*TeamReference           `json:"teams,omitempty"`
	Type                   string                     `json:"type,omitempty"`
}

// GetIntegrationOptions represents options when retrieving a service integration.
type GetIntegrationOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// ListServicesOptions represents options when listing services.
type ListServicesOptions struct {
	Limit    int      `url:"limit,omitempty"`
	More     bool     `url:"more,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Total    int      `url:"total,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
	Query    string   `url:"query,omitempty"`
	SortBy   string   `url:"sort_by,omitempty"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	TimeZone string   `url:"time_zone,omitempty"`
}

// ListServicesResponse represents a list response of services.
type ListServicesResponse struct {
	Limit    int  `json:"limit,omitempty"`
	More     bool `json:"more,omitempty"`
	Offset   int  `json:"offset,omitempty"`
	Total    int  `json:"total,omitempty"`
	Services []*Service
}

// GetServiceOptions represents options when retrieving a service.
type GetServiceOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// List lists existing services.
func (s *ServicesService) List(o *ListServicesOptions) (*ListServicesResponse, *Response, error) {
	u := "/services"
	v := new(ListServicesResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new service.
func (s *ServicesService) Create(service *Service) (*Service, *Response, error) {
	u := "/services"
	v := new(Service)

	resp, err := s.client.newRequestDo("POST", u, nil, &Service{Service: service}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Service, resp, nil
}

// Delete removes an existing service.
func (s *ServicesService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/services/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about a service.
func (s *ServicesService) Get(id string, o *GetServiceOptions) (*Service, *Response, error) {
	u := fmt.Sprintf("/services/%s", id)
	v := new(Service)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Service, resp, nil
}

// Update updates an existing service.
func (s *ServicesService) Update(id string, service *Service) (*Service, *Response, error) {
	u := fmt.Sprintf("/services/%s", id)
	v := new(Service)

	resp, err := s.client.newRequestDo("PUT", u, nil, &Service{Service: service}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Service, resp, nil
}

// CreateIntegration creates a new service integration.
func (s *ServicesService) CreateIntegration(serviceID string, integration *Integration) (*Integration, *Response, error) {
	u := fmt.Sprintf("/services/%s/integrations", serviceID)
	v := new(Integration)

	resp, err := s.client.newRequestDo("POST", u, nil, &Integration{Integration: integration}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

// GetIntegration retrieves information about a service integration.
func (s *ServicesService) GetIntegration(serviceID, integrationID string, o *GetIntegrationOptions) (*Integration, *Response, error) {
	u := fmt.Sprintf("/services/%s/integrations/%s", serviceID, integrationID)
	v := new(Integration)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

// UpdateIntegration updates an existing service integration.
func (s *ServicesService) UpdateIntegration(serviceID, integrationID string, integration *Integration) (*Integration, *Response, error) {
	u := fmt.Sprintf("/services/%s/integrations/%s", serviceID, integrationID)
	v := new(Integration)

	resp, err := s.client.newRequestDo("PUT", u, nil, &Integration{Integration: integration}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

// DeleteIntegration removes an existing service integration.
func (s *ServicesService) DeleteIntegration(serviceID, integrationID string) (*Response, error) {
	u := fmt.Sprintf("/services/%s/integrations/%s", serviceID, integrationID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
