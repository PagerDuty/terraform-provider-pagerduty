package pagerduty

import (
	"fmt"
)

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

// AlertGroupingConfig - populate timeout if AlertGroupingParameters Type is 'time', populate Aggregate & Fields if Type is 'content_grouping'
type AlertGroupingConfig struct {
	Timeout    *int     `json:"timeout,omitempty"`
	TimeWindow *int     `json:"time_window,omitempty"`
	Aggregate  *string  `json:"aggregate,omitempty"`
	Fields     []string `json:"fields,omitempty"`
}

// AlertGroupingParameters defines how alerts are grouped into incidents
type AlertGroupingParameters struct {
	Type   *string              `json:"type,omitempty"`
	Config *AlertGroupingConfig `json:"config,omitempty"`
}

// AutoPauseNotificationsParameters defines how alerts on this service are automatically suspended for a period of time before triggering, when identified as likely being transient.
type AutoPauseNotificationsParameters struct {
	Enabled bool `json:"enabled"`
	Timeout *int `json:"timeout"`
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
	CreatedAt             string            `json:"created_at,omitempty"`
	EmailIncidentCreation string            `json:"email_incident_creation,omitempty"`
	EmailFilterMode       string            `json:"email_filter_mode,omitempty"`
	EmailParsers          []*EmailParser    `json:"email_parsers,omitempty"`
	EmailParsingFallback  string            `json:"email_parsing_fallback,omitempty"`
	EmailFilters          []*EmailFilter    `json:"email_filters,omitempty"`
	HTMLURL               string            `json:"html_url,omitempty"`
	ID                    string            `json:"id,omitempty"`
	Integration           *Integration      `json:"integration,omitempty"`
	IntegrationEmail      string            `json:"integration_email,omitempty"`
	IntegrationKey        string            `json:"integration_key,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Self                  string            `json:"self,omitempty"`
	Service               *ServiceReference `json:"service,omitempty"`
	Summary               string            `json:"summary,omitempty"`
	Type                  string            `json:"type,omitempty"`
	Vendor                *VendorReference  `json:"vendor,omitempty"`
}

// EmailFilter represents a integration email filters
type EmailFilter struct {
	BodyMode       string `json:"body_mode,omitempty"`
	BodyRegex      string `json:"body_regex,omitempty"`
	FromEmailMode  string `json:"from_email_mode,omitempty"`
	FromEmailRegex string `json:"from_email_regex,omitempty"`
	ID             string `json:"id,omitempty"`
	SubjectMode    string `json:"subject_mode,omitempty"`
	SubjectRegex   string `json:"subject_regex,omitempty"`
}

// EmailParser represents a integration email parsers
type EmailParser struct {
	Action          string            `json:"action,omitempty"`
	ID              *int              `json:"id,omitempty"`
	MatchPredicate  *MatchPredicate   `json:"match_predicate,omitempty"`
	ValueExtractors []*ValueExtractor `json:"value_extractors,omitempty"`
}

// MatchPredicate represents a integration email MatchPredicate
type MatchPredicate struct {
	Predicates []*Predicate `json:"children,omitempty"`
	Type       string       `json:"type,omitempty"`
}

// Predicate represents a integration email Predicate
type Predicate struct {
	Matcher    string       `json:"matcher,omitempty"`
	Part       string       `json:"part,omitempty"`
	Predicates []*Predicate `json:"children,omitempty"`
	Type       string       `json:"type,omitempty"`
}

// ValueExtractor represents a integration email ValueExtractor
type ValueExtractor struct {
	ValueName   string `json:"value_name,omitempty"`
	Part        string `json:"part,omitempty"`
	StartsAfter string `json:"starts_after"`
	EndsBefore  string `json:"ends_before"`
	Type        string `json:"type,omitempty"`
	Regex       string `json:"regex,omitempty"`
}

// Service represents a service.
type Service struct {
	AcknowledgementTimeout           *int                              `json:"acknowledgement_timeout"`
	Addons                           []*AddonReference                 `json:"addons,omitempty"`
	AlertCreation                    string                            `json:"alert_creation,omitempty"`
	AlertGrouping                    *string                           `json:"alert_grouping"`
	AlertGroupingTimeout             *int                              `json:"alert_grouping_timeout,omitempty"`
	AlertGroupingParameters          *AlertGroupingParameters          `json:"alert_grouping_parameters,omitempty"`
	AutoPauseNotificationsParameters *AutoPauseNotificationsParameters `json:"auto_pause_notifications_parameters,omitempty"`
	AutoResolveTimeout               *int                              `json:"auto_resolve_timeout"`
	CreatedAt                        string                            `json:"created_at,omitempty"`
	Description                      string                            `json:"description,omitempty"`
	EscalationPolicy                 *EscalationPolicyReference        `json:"escalation_policy,omitempty"`
	ResponsePlay                     *ResponsePlayReference            `json:"response_play"`
	HTMLURL                          string                            `json:"html_url,omitempty"`
	ID                               string                            `json:"id,omitempty"`
	IncidentUrgencyRule              *IncidentUrgencyRule              `json:"incident_urgency_rule,omitempty"`
	Integrations                     []*IntegrationReference           `json:"integrations,omitempty"`
	LastIncidentTimestamp            string                            `json:"last_incident_timestamp,omitempty"`
	Name                             string                            `json:"name,omitempty"`
	ScheduledActions                 []*ScheduledAction                `json:"scheduled_actions,omitempty"`
	Self                             string                            `json:"self,omitempty"`
	Status                           string                            `json:"status,omitempty"`
	Summary                          string                            `json:"summary,omitempty"`
	SupportHours                     *SupportHours                     `json:"support_hours,omitempty"`
	Teams                            []*TeamReference                  `json:"teams,omitempty"`
	Type                             string                            `json:"type,omitempty"`
}

// ServicePayload represents a service.
type ServicePayload struct {
	Service *Service `json:"service,omitempty"`
}

// ServiceEventRule represents a service event rule
type ServiceEventRule struct {
	ID         string            `json:"id,omitempty"`
	Self       string            `json:"self,omitempty"`
	Disabled   bool              `json:"disabled"`
	Conditions *RuleConditions   `json:"conditions,omitempty"`
	TimeFrame  *RuleTimeFrame    `json:"time_frame,omitempty"`
	Variables  []*RuleVariable   `json:"variables,omitempty"`
	Position   *int              `json:"position,omitempty"`
	Actions    *RuleActions      `json:"actions,omitempty"`
	Service    *ServiceReference `json:"service_id,omitempty"`
}

// IntegrationPayload represents an integration.
type IntegrationPayload struct {
	Integration *Integration `json:"integration,omitempty"`
}

// ServiceEventRulePayload represents a payload for service event rules
type ServiceEventRulePayload struct {
	Rule *ServiceEventRule `json:"rule,omitempty"`
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

// ListServiceEventRuleOptions represents options when retrieving a list of event rules for a service
type ListServiceEventRuleOptions struct {
	Limit  int  `json:"limit,omitempty"`
	More   bool `json:"more,omitempty"`
	Offset int  `json:"offset,omitempty"`
	Total  int  `json:"total,omitempty"`
}

// ListServiceEventRuleResponse represents a list of event rules for a service
type ListServiceEventRuleResponse struct {
	Limit      int                 `json:"limit,omitempty"`
	More       bool                `json:"more,omitempty"`
	Offset     int                 `json:"offset,omitempty"`
	Total      int                 `json:"total,omitempty"`
	EventRules []*ServiceEventRule `json:"rules,omitempty"`
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
	v := new(ServicePayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &ServicePayload{Service: service}, &v)
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
	v := new(ServicePayload)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Service, resp, nil
}

// Update updates an existing service.
func (s *ServicesService) Update(id string, service *Service) (*Service, *Response, error) {
	u := fmt.Sprintf("/services/%s", id)
	v := new(ServicePayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &ServicePayload{Service: service}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Service, resp, nil
}

// CreateIntegration creates a new service integration.
func (s *ServicesService) CreateIntegration(serviceID string, integration *Integration) (*Integration, *Response, error) {
	u := fmt.Sprintf("/services/%s/integrations", serviceID)
	v := new(IntegrationPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &IntegrationPayload{Integration: integration}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

// GetIntegration retrieves information about a service integration.
func (s *ServicesService) GetIntegration(serviceID, integrationID string, o *GetIntegrationOptions) (*Integration, *Response, error) {
	u := fmt.Sprintf("/services/%s/integrations/%s", serviceID, integrationID)
	v := new(IntegrationPayload)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

// UpdateIntegration updates an existing service integration.
func (s *ServicesService) UpdateIntegration(serviceID, integrationID string, integration *Integration) (*Integration, *Response, error) {
	u := fmt.Sprintf("/services/%s/integrations/%s", serviceID, integrationID)
	v := new(IntegrationPayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &IntegrationPayload{Integration: integration}, &v)
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

// ListEventRules lists existing service event rules.
func (s *ServicesService) ListEventRules(serviceID string, o *ListServiceEventRuleOptions) (*ListServiceEventRuleResponse, *Response, error) {
	u := fmt.Sprintf("/services/%s/rules", serviceID)
	v := new(ListServiceEventRuleResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// CreateEventRule creates a new service event rule.
func (s *ServicesService) CreateEventRule(serviceID string, eventRule *ServiceEventRule) (*ServiceEventRule, *Response, error) {
	u := fmt.Sprintf("/services/%s/rules", serviceID)
	v := new(ServiceEventRulePayload)
	p := ServiceEventRulePayload{Rule: eventRule}

	resp, err := s.client.newRequestDo("POST", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Rule, resp, nil
}

// GetEventRule retrieves information about a service event rule.
func (s *ServicesService) GetEventRule(serviceID, ruleID string) (*ServiceEventRule, *Response, error) {
	u := fmt.Sprintf("/services/%s/rules/%s", serviceID, ruleID)
	v := new(ServiceEventRulePayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Rule, resp, nil
}

// UpdateEventRule updates an existing service event rule.
func (s *ServicesService) UpdateEventRule(serviceID, ruleID string, eventRule *ServiceEventRule) (*ServiceEventRule, *Response, error) {
	u := fmt.Sprintf("/services/%s/rules/%s", serviceID, ruleID)
	v := new(ServiceEventRulePayload)
	p := ServiceEventRulePayload{Rule: eventRule}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Rule, resp, nil
}

// DeleteEventRule removes an existing service event rule.
func (s *ServicesService) DeleteEventRule(serviceID, ruleID string) (*Response, error) {
	u := fmt.Sprintf("/services/%s/rules/%s", serviceID, ruleID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
