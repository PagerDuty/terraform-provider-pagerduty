package pagerduty

import (
	"fmt"
)

// RulesetService handles the communication with rulesets
// related methods of the PagerDuty API.
type RulesetService service

// Ruleset represents a ruleset.
type Ruleset struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Type        string         `json:"type,omitempty"`
	RoutingKeys []string       `json:"routing_keys,omitempty"`
	Team        *RulesetObject `json:"team,omitempty"`
	Updater     *RulesetObject `json:"updater,omitempty"`
	Creator     *RulesetObject `json:"creator,omitempty"`
}

// RulesetObject represents a generic object that is common within a ruleset object
type RulesetObject struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// RulesetPayload represents payload with a ruleset object
type RulesetPayload struct {
	Ruleset *Ruleset `json:"ruleset,omitempty"`
}

// ListRulesetsResponse represents a list response of rulesets.
type ListRulesetsResponse struct {
	Total    int        `json:"total,omitempty"`
	Rulesets []*Ruleset `json:"rulesets,omitempty"`
	Offset   int        `json:"offset,omitempty"`
	More     bool       `json:"more,omitempty"`
	Limit    int        `json:"limit,omitempty"`
}

// RulesetRule represents a Ruleset rule
type RulesetRule struct {
	ID         string            `json:"id,omitempty"`
	Position   *int              `json:"position,omitempty"`
	Disabled   bool              `json:"disabled"`
	Conditions *RuleConditions   `json:"conditions,omitempty"`
	Actions    *RuleActions      `json:"actions,omitempty"`
	Ruleset    *RulesetReference `json:"ruleset,omitempty"`
	Self       string            `json:"self,omitempty"`
	CatchAll   bool              `json:"catch_all,omitempty"`
	TimeFrame  *RuleTimeFrame    `json:"time_frame,omitempty"`
	Variables  []*RuleVariable   `json:"variables,omitempty"`
}

// RulesetRulePayload represents a payload for ruleset rules
type RulesetRulePayload struct {
	Rule *RulesetRule `json:"rule,omitempty"`
}

// RuleConditions represents the conditions field for a Ruleset
type RuleConditions struct {
	Operator          string              `json:"operator,omitempty"`
	RuleSubconditions []*RuleSubcondition `json:"subconditions,omitempty"`
}

// RuleSubcondition represents a subcondition of a ruleset condition
type RuleSubcondition struct {
	Operator   string              `json:"operator,omitempty"`
	Parameters *ConditionParameter `json:"parameters,omitempty"`
}

// ConditionParameter represents  parameters in a rule condition
type ConditionParameter struct {
	Path  string `json:"path,omitempty"`
	Value string `json:"value,omitempty"`
}

// RuleTimeFrame represents a time_frame object on the rule object
type RuleTimeFrame struct {
	ScheduledWeekly *ScheduledWeekly `json:"scheduled_weekly,omitempty"`
	ActiveBetween   *ActiveBetween   `json:"active_between,omitempty"`
}

// RuleVariable represents a rule variable
type RuleVariable struct {
	Name       string                 `json:"name,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Parameters *RuleVariableParameter `json:"parameters,omitempty"`
}

// RuleVariableParameter represents a rule variable parameter
type RuleVariableParameter struct {
	Value string `json:"value"`
	Path  string `json:"path"`
}

// ScheduledWeekly represents a time_frame object for scheduling rules weekly
type ScheduledWeekly struct {
	Weekdays  []int  `json:"weekdays,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	StartTime int    `json:"start_time,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}

// ActiveBetween represents an active_between object for setting a timeline for rules
type ActiveBetween struct {
	StartTime int `json:"start_time,omitempty"`
	EndTime   int `json:"end_time,omitempty"`
}

// ListRulesetRulesResponse represents a list of rules in a ruleset
type ListRulesetRulesResponse struct {
	Total  int            `json:"total,omitempty"`
	Rules  []*RulesetRule `json:"rules,omitempty"`
	Offset int            `json:"offset,omitempty"`
	More   bool           `json:"more,omitempty"`
	Limit  int            `json:"limit,omitempty"`
}

// RuleActions represents a rule action
type RuleActions struct {
	Suppress    *RuleActionSuppress     `json:"suppress,omitempty"`
	Annotate    *RuleActionParameter    `json:"annotate,omitempty"`
	Severity    *RuleActionParameter    `json:"severity,omitempty"`
	Priority    *RuleActionParameter    `json:"priority,omitempty"`
	Route       *RuleActionParameter    `json:"route,omitempty"`
	EventAction *RuleActionParameter    `json:"event_action,omitempty"`
	Extractions []*RuleActionExtraction `json:"extractions,omitempty"`
	Suspend     *RuleActionIntParameter `json:"suspend,omitempty"`
}

// RuleActionParameter represents a string parameter object on a rule action
type RuleActionParameter struct {
	Value string `json:"value,omitempty"`
}

// RuleActionIntParameter represents an integer parameter object on a rule action
type RuleActionIntParameter struct {
	Value int `json:"value"`
}

// RuleActionSuppress represents a rule suppress action object
type RuleActionSuppress struct {
	Value               bool   `json:"value"`
	ThresholdValue      int    `json:"threshold_value,omitempty"`
	ThresholdTimeUnit   string `json:"threshold_time_unit,omitempty"`
	ThresholdTimeAmount int    `json:"threshold_time_amount,omitempty"`
}

// RuleActionExtraction represents a rule extraction action object
type RuleActionExtraction struct {
	Target   string `json:"target,omitempty"`
	Source   string `json:"source,omitempty"`
	Regex    string `json:"regex,omitempty"`
	Template string `json:"template,omitempty"`
}

// List lists existing rulesets.
func (s *RulesetService) List() (*ListRulesetsResponse, *Response, error) {
	u := "/rulesets"
	v := new(ListRulesetsResponse)

	rulesets := make([]*Ruleset, 0)

	// Create a handler closure capable of parsing data from the rulesets endpoint
	// and appending resultant rulesets to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListRulesetsResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		rulesets = append(rulesets, result.Rulesets...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(u, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.Rulesets = rulesets

	return v, nil, nil
}

// Create creates a new ruleset.
func (s *RulesetService) Create(ruleset *Ruleset) (*Ruleset, *Response, error) {
	u := "/rulesets"
	v := new(RulesetPayload)
	p := &RulesetPayload{Ruleset: ruleset}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Ruleset, resp, nil
}

// Get gets a new ruleset.
func (s *RulesetService) Get(ID string) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("/rulesets/%s", ID)
	v := new(RulesetPayload)
	p := &RulesetPayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Ruleset, resp, nil
}

// Delete deletes an existing ruleset.
func (s *RulesetService) Delete(ID string) (*Response, error) {
	u := fmt.Sprintf("/rulesets/%s", ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Update updates an existing ruleset.
func (s *RulesetService) Update(ID string, ruleset *Ruleset) (*Ruleset, *Response, error) {
	u := fmt.Sprintf("/rulesets/%s", ID)
	v := new(RulesetPayload)
	p := RulesetPayload{Ruleset: ruleset}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Ruleset, resp, nil
}

// ListRules Lists Event Rules for Ruleset
func (s *RulesetService) ListRules(rulesetID string) (*ListRulesetRulesResponse, *Response, error) {
	u := fmt.Sprintf("/rulesets/%s/rules", rulesetID)
	v := new(ListRulesetRulesResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// CreateRule for Ruleset
func (s *RulesetService) CreateRule(rulesetID string, rule *RulesetRule) (*RulesetRule, *Response, error) {
	u := fmt.Sprintf("/rulesets/%s/rules", rulesetID)
	v := new(RulesetRulePayload)
	p := RulesetRulePayload{Rule: rule}

	resp, err := s.client.newRequestDo("POST", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Rule, resp, nil
}

// GetRule for Ruleset
func (s *RulesetService) GetRule(rulesetID, ruleID string) (*RulesetRule, *Response, error) {
	u := fmt.Sprintf("/rulesets/%s/rules/%s", rulesetID, ruleID)
	v := new(RulesetRulePayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Rule, resp, nil
}

// UpdateRule for Ruleset
func (s *RulesetService) UpdateRule(rulesetID, ruleID string, rule *RulesetRule) (*RulesetRule, *Response, error) {
	u := fmt.Sprintf("/rulesets/%s/rules/%s", rulesetID, ruleID)
	v := new(RulesetRulePayload)
	p := RulesetRulePayload{Rule: rule}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Rule, resp, nil
}

// DeleteRule deletes an existing rule from the ruleset.
func (s *RulesetService) DeleteRule(rulesetID, ruleID string) (*Response, error) {
	u := fmt.Sprintf("/rulesets/%s/rules/%s", rulesetID, ruleID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
