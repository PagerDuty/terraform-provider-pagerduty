package pagerduty

import (
	"fmt"
)

type EventOrchestrationPathService service

type EventOrchestrationPath struct {
	Type      string                           `json:"type,omitempty"`
	Self      string                           `json:"self,omitempty"`
	Parent    *EventOrchestrationPathReference `json:"parent,omitempty"`
	Sets      []*EventOrchestrationPathSet     `json:"sets,omitempty"`
	CatchAll  *EventOrchestrationPathCatchAll  `json:"catch_all,omitempty"`
	CreatedAt string                           `json:"created_at,omitempty"`
	CreatedBy *EventOrchestrationPathReference `json:"created_by,omitempty"`
	UpdatedAt string                           `json:"updated_at,omitempty"`
	UpdatedBy *EventOrchestrationPathReference `json:"updated_by,omitempty"`
	Version   string                           `json:"version,omitempty"`
}

// A reference to a related object (e.g. an EventOrchestration, User, Team, etc)
type EventOrchestrationPathReference struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Self string `json:"self,omitempty"`
}

type EventOrchestrationPathSet struct {
	ID    string                        `json:"id,omitempty"`
	Rules []*EventOrchestrationPathRule `json:"rules"`
}

type EventOrchestrationPathRule struct {
	ID         string                                 `json:"id,omitempty"`
	Label      string                                 `json:"label,omitempty"`
	Conditions []*EventOrchestrationPathRuleCondition `json:"conditions"`
	Actions    *EventOrchestrationPathRuleActions     `json:"actions,omitempty"`
	Disabled   bool                                   `json:"disabled,omitempty"`
}

type EventOrchestrationPathRuleCondition struct {
	// A PCL string: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview
	Expression string `json:"expression,omitempty"`
}

// See the full list of supported actions for path types:
// Router: https://developer.pagerduty.com/api-reference/f0fae270c70b3-get-the-router-for-a-global-event-orchestration
// Service: https://developer.pagerduty.com/api-reference/179537b835e2d-get-the-service-orchestration-for-a-service
// Unrouted: https://developer.pagerduty.com/api-reference/70aa1139e1013-get-the-unrouted-orchestration-for-a-global-event-orchestration
type EventOrchestrationPathRuleActions struct {
	RouteTo                    string                                             `json:"route_to"`
	Suppress                   bool                                               `json:"suppress"`
	Suspend                    *int                                               `json:"suspend"`
	Priority                   string                                             `json:"priority"`
	Annotate                   string                                             `json:"annotate"`
	PagerdutyAutomationActions []*EventOrchestrationPathPagerdutyAutomationAction `json:"pagerduty_automation_actions"`
	AutomationActions          []*EventOrchestrationPathAutomationAction          `json:"automation_actions"`
	Severity                   string                                             `json:"severity"`
	EventAction                string                                             `json:"event_action"`
	Variables                  []*EventOrchestrationPathActionVariables           `json:"variables"`
	Extractions                []*EventOrchestrationPathActionExtractions         `json:"extractions"`
}

type EventOrchestrationPathPagerdutyAutomationAction struct {
	ActionId string `json:"action_id,omitempty"`
}

type EventOrchestrationPathAutomationAction struct {
	Name       string                                          `json:"name,omitempty"`
	Url        string                                          `json:"url,omitempty"`
	AutoSend   bool                                            `json:"auto_send,omitempty"`
	Headers    []*EventOrchestrationPathAutomationActionObject `json:"headers"`
	Parameters []*EventOrchestrationPathAutomationActionObject `json:"parameters"`
}

type EventOrchestrationPathAutomationActionObject struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type EventOrchestrationPathActionVariables struct {
	Name  string `json:"name,omitempty"`
	Path  string `json:"path,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type EventOrchestrationPathActionExtractions struct {
	Target   string `json:"target,omitempty"`
	Regex    string `json:"regex,omitempty"`
	Template string `json:"template,omitempty"`
	Source   string `json:"source,omitempty"`
}

type EventOrchestrationPathServiceActiveStatus struct {
	Active bool `json:"active"`
}

type EventOrchestrationPathCatchAll struct {
	Actions *EventOrchestrationPathRuleActions `json:"actions,omitempty"`
}

type EventOrchestrationPathPayload struct {
	OrchestrationPath *EventOrchestrationPath `json:"orchestration_path,omitempty"`
}

const PathTypeRouter string = "router"
const PathTypeService string = "service"
const PathTypeUnrouted string = "unrouted"

func orchestrationPathUrlBuilder(id string, pathType string) string {
	switch {
	case pathType == PathTypeService:
		return fmt.Sprintf("%s/services/%s", eventOrchestrationBaseUrl, id)
	case pathType == PathTypeUnrouted:
		return fmt.Sprintf("%s/%s/unrouted", eventOrchestrationBaseUrl, id)
	case pathType == PathTypeRouter:
		return fmt.Sprintf("%s/%s/router", eventOrchestrationBaseUrl, id)
	default:
		return ""
	}
}

// Get for EventOrchestrationPath
func (s *EventOrchestrationPathService) Get(id string, pathType string) (*EventOrchestrationPath, *Response, error) {
	u := orchestrationPathUrlBuilder(id, pathType)
	v := new(EventOrchestrationPathPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)

	if err != nil {
		return nil, nil, err
	}

	return v.OrchestrationPath, resp, nil
}

// GetServiceActiveStatus for EventOrchestrationPath
func (s *EventOrchestrationPathService) GetServiceActiveStatus(id string) (*EventOrchestrationPathServiceActiveStatus, *Response, error) {
	u := fmt.Sprintf("%s/services/%s/active", eventOrchestrationBaseUrl, id)
	v := new(EventOrchestrationPathServiceActiveStatus)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)

	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Update for EventOrchestrationPath
func (s *EventOrchestrationPathService) Update(id string, pathType string, orchestration_path *EventOrchestrationPath) (*EventOrchestrationPath, *Response, error) {
	u := orchestrationPathUrlBuilder(id, pathType)
	v := new(EventOrchestrationPathPayload)
	p := EventOrchestrationPathPayload{OrchestrationPath: orchestration_path}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.OrchestrationPath, resp, nil
}

// UpdateServiceActiveStatus for EventOrchestrationPath
func (s *EventOrchestrationPathService) UpdateServiceActiveStatus(id string, isActive bool) (*EventOrchestrationPathServiceActiveStatus, *Response, error) {
	u := fmt.Sprintf("%s/services/%s/active", eventOrchestrationBaseUrl, id)
	v := new(EventOrchestrationPathServiceActiveStatus)
	p := EventOrchestrationPathServiceActiveStatus{Active: isActive}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, err
}
