package pagerduty

import "fmt"

// EscalationPolicyService handles the communication with escalation policy
// related methods of the PagerDuty API.
type EscalationPolicyService service

// EscalationRuleAssignmentStrategy represents an Escalation rule assignment
// strategy
type EscalationRuleAssignmentStrategy struct {
	Type string `json:"type,omitempty"`
}

// EscalationRule represents an escalation rule.
type EscalationRule struct {
	EscalationDelayInMinutes         int                               `json:"escalation_delay_in_minutes,omitempty"`
	EscalationRuleAssignmentStrategy *EscalationRuleAssignmentStrategy `json:"escalation_rule_assignment_strategy,omitempty"`
	ID                               string                            `json:"id,omitempty"`
	Targets                          []*EscalationTargetReference      `json:"targets,omitempty"`
}

// EscalationPolicy represents an escalation policy.
type EscalationPolicy struct {
	Description     string              `json:"description,omitempty"`
	EscalationRules []*EscalationRule   `json:"escalation_rules,omitempty"`
	HTMLURL         string              `json:"html_url,omitempty"`
	ID              string              `json:"id,omitempty"`
	Name            string              `json:"name,omitempty"`
	NumLoops        *int                `json:"num_loops,omitempty"`
	RepeatEnabled   bool                `json:"repeat_enabled,omitempty"`
	Self            string              `json:"self,omitempty"`
	Services        []*ServiceReference `json:"services,omitempty"`
	Summary         string              `json:"summary,omitempty"`
	Teams           []*TeamReference    `json:"teams"`
	Type            string              `json:"type,omitempty"`
}

// ListEscalationPoliciesResponse represents a list response of escalation policies.
type ListEscalationPoliciesResponse struct {
	Limit              int                 `json:"limit,omitempty"`
	More               bool                `json:"more,omitempty"`
	Offset             int                 `json:"offset,omitempty"`
	Total              int                 `json:"total,omitempty"`
	EscalationPolicies []*EscalationPolicy `json:"escalation_policies,omitempty"`
}

// ListEscalationRulesResponse represents a list response of escalation rules.
type ListEscalationRulesResponse struct {
	Limit           int               `json:"limit,omitempty"`
	More            bool              `json:"more,omitempty"`
	Offset          int               `json:"offset,omitempty"`
	Total           int               `json:"total,omitempty"`
	EscalationRules []*EscalationRule `json:"escalation_rules,omitempty"`
}

// ListEscalationPoliciesOptions represents options when listing escalation policies.
type ListEscalationPoliciesOptions struct {
	Limit    int      `url:"limit,omitempty"`
	More     bool     `url:"more,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Total    int      `url:"total,omitempty"`
	Includes []string `url:"include,omitempty,brackets"`
	Query    string   `url:"query,omitempty"`
	SortBy   string   `url:"sort_by,omitempty"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	UserIDs  []string `url:"user_ids,omitempty,brackets"`
}

// GetEscalationRuleOptions represents options when retrieving an escalation rule.
type GetEscalationRuleOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetEscalationPolicyOptions represents options when retrieving an escalation policy.
type GetEscalationPolicyOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// List lists existing escalation policies.
func (s *EscalationPolicyService) List(o *ListEscalationPoliciesOptions) (*ListEscalationPoliciesResponse, *Response, error) {
	u := "/escalation_policies"
	v := new(ListEscalationPoliciesResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// EscalationPolicyPayload represents an escalation policy.
type EscalationPolicyPayload struct {
	EscalationPolicy *EscalationPolicy `json:"escalation_policy"`
}

// Create creates a new escalation policy.
func (s *EscalationPolicyService) Create(escalationPolicy *EscalationPolicy) (*EscalationPolicy, *Response, error) {
	u := "/escalation_policies"
	v := new(EscalationPolicyPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &EscalationPolicyPayload{EscalationPolicy: escalationPolicy}, v)
	if err != nil {
		return nil, nil, err
	}

	return v.EscalationPolicy, resp, nil
}

// Delete deletes an existing escalation policy.
func (s *EscalationPolicyService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/escalation_policies/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about an escalation policy.
func (s *EscalationPolicyService) Get(id string, o *GetEscalationPolicyOptions) (*EscalationPolicy, *Response, error) {
	u := fmt.Sprintf("/escalation_policies/%s", id)
	v := new(EscalationPolicyPayload)

	resp, err := s.client.newRequestDo("GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.EscalationPolicy, resp, nil
}

// Update updates an existing escalation policy.
func (s *EscalationPolicyService) Update(id string, escalationPolicy *EscalationPolicy) (*EscalationPolicy, *Response, error) {
	u := fmt.Sprintf("/escalation_policies/%s", id)
	v := new(EscalationPolicyPayload)

	resp, err := s.client.newRequestDo("PUT", u, nil, &EscalationPolicyPayload{EscalationPolicy: escalationPolicy}, v)
	if err != nil {
		return nil, nil, err
	}

	return v.EscalationPolicy, resp, nil
}
