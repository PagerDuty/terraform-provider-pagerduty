package pagerduty

// OnCallService handles the communication with team
// related methods of the PagerDuty API.
type OnCallService service

// OnCall represents an oncall.
type OnCall struct {
	User             *UserReference             `json:"user,omitempty"`
	Schedule         *ScheduleReference         `json:"schedule,omitemtpy"`
	EscalationPolicy *EscalationPolicyReference `json:"escalation_policy,omitempty"`
	EscalationLevel  int                        `json:"escalation_level"`
	Start            *string                    `json:"start"`
	End              *string                    `json:"end"`
}

// ListOnCallOptions represents options when listing oncalls.
type ListOnCallOptions struct {
	Limit               int      `url:"limit,omitempty"`
	Offset              int      `url:"offset,omitempty"`
	Total               bool     `url:"total,omitempty"`
	Earliest            bool     `url:"earliest,omitempty"`
	EscalationPolicyIds []string `url:"escalation_policy_ids,omitempty"`
	Includes            []string `url:"include,omitempty"`
	ScheduleIds         []string `url:"schedule_ids,omitempty"`
	UserIds             []string `url:"user_ids,brackets,omitempty"`
	Since               string   `url:"since,omitempty"`
	TimeZone            string   `url:"time_zone,omitempty"`
	Until               string   `url:"until,omitempty"`
}

// ListOnCallResponse represents a list response of oncalls.
type ListOnCallResponse struct {
	Oncalls []*OnCall `json:"oncalls,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	More    bool      `json:"more,omitempty"`
	Offset  int       `json:"offset,omitempty"`
	Total   int       `json:"total,omitempty"`
}

// List lists existing oncalls.
func (s *OnCallService) List(o *ListOnCallOptions) (*ListOnCallResponse, *Response, error) {
	u := "/oncalls"
	v := new(ListOnCallResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
