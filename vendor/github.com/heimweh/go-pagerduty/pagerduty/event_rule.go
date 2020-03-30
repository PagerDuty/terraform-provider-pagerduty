package pagerduty

import "fmt"

// EventRuleService handles the communication with event rules
// related methods of the PagerDuty API.
type EventRuleService service

// EventRule represents an event rule.
type EventRule struct {
	Actions           []interface{} `json:"actions,omitempty"`
	AdvancedCondition []interface{} `json:"advanced_condition,omitempty"`
	CatchAll          bool          `json:"catch_all,omitempty"`
	Condition         []interface{} `json:"condition,omitempty"`
	ID                string        `json:"id,omitempty"`
}

// ListEventRulesResponse represents a list response of event rules.
type ListEventRulesResponse struct {
	ExternalID    string       `json:"external_id,omitempty"`
	ObjectVersion string       `json:"object_version,omitempty"`
	FormatVersion int          `json:"format_version,string,omitempty"`
	EventRules    []*EventRule `json:"rules,omitempty"`
}

// List lists existing event rules.
func (s *EventRuleService) List() (*ListEventRulesResponse, *Response, error) {
	u := "/event_rules"
	v := new(ListEventRulesResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new event rule.
func (s *EventRuleService) Create(eventRule *EventRule) (*EventRule, *Response, error) {
	u := "/event_rules"
	v := new(EventRule)

	resp, err := s.client.newRequestDo("POST", u, nil, eventRule, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Delete deletes an existing event rule.
func (s *EventRuleService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/event_rules/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Update updates an existing event rule.
func (s *EventRuleService) Update(id string, eventRule *EventRule) (*EventRule, *Response, error) {
	u := fmt.Sprintf("/event_rules/%s", id)
	v := new(EventRule)

	resp, err := s.client.newRequestDo("PUT", u, nil, eventRule, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
