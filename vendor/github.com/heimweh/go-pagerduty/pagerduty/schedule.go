package pagerduty

import (
	"fmt"
)

// ScheduleService handles the communication with schedule
// related methods of the PagerDuty API.
type ScheduleService service

// Override represents an override
type Override struct {
	ID    string         `json:"id,omitempty"`
	Start string         `json:"start,omitempty"`
	End   string         `json:"end,omitempty"`
	User  *UserReference `json:"user,omitempty"`
}

// Schedule represents a schedule.
type Schedule struct {
	Description          string                       `json:"description,omitempty"`
	EscalationPolicies   []*EscalationPolicyReference `json:"escalation_policies,omitempty"`
	FinalSchedule        *SubSchedule                 `json:"final_schedule,omitempty"`
	HTMLURL              string                       `json:"html_url,omitempty"`
	ID                   string                       `json:"id,omitempty"`
	Name                 string                       `json:"name,omitempty"`
	OverridesSubSchedule *SubSchedule                 `json:"overrides_subschedule,omitempty"`
	ScheduleLayers       []*ScheduleLayer             `json:"schedule_layers,omitempty"`
	Self                 string                       `json:"self,omitempty"`
	Summary              string                       `json:"summary,omitempty"`
	TimeZone             string                       `json:"time_zone,omitempty"`
	Type                 string                       `json:"type,omitempty"`
	Users                []*UserReference             `json:"users,omitempty"`
	Teams                []*TeamReference             `json:"teams,omitempty"`
}

// SubSchedule represents a sub-schedule of a schedule.
type SubSchedule struct {
	Name                       string                `json:"name,omitempty"`
	RenderedCoveragePercentage float64               `json:"rendered_coverage_percentage,omitempty"`
	RenderedScheduleEntries    []*ScheduleLayerEntry `json:"rendered_schedule_entries,omitempty"`
}

// Restriction represents a schedule layer restriction.
type Restriction struct {
	DurationSeconds int    `json:"duration_seconds,omitempty"`
	StartDayOfWeek  int    `json:"start_day_of_week,omitempty"`
	StartTimeOfDay  string `json:"start_time_of_day,omitempty"`
	Type            string `json:"type,omitempty"`
}

// ScheduleLayerEntry represents a rendered schedule layer entry.
type ScheduleLayerEntry struct {
	End   string         `json:"end,omitempty"`
	Start string         `json:"start,omitempty"`
	User  *UserReference `json:"user,omitempty"`
}

// ScheduleLayer represents a schedule layer in a schedule
type ScheduleLayer struct {
	End                        string                  `json:"end,omitempty"`
	ID                         string                  `json:"id,omitempty"`
	Name                       string                  `json:"name,omitempty"`
	RenderedCoveragePercentage float64                 `json:"rendered_coverage_percentage,omitempty"`
	RenderedScheduleEntries    []*ScheduleLayerEntry   `json:"rendered_schedule_entries,omitempty"`
	Restrictions               []*Restriction          `json:"restrictions,omitempty"`
	RotationTurnLengthSeconds  int                     `json:"rotation_turn_length_seconds,omitempty"`
	RotationVirtualStart       string                  `json:"rotation_virtual_start,omitempty"`
	Start                      string                  `json:"start,omitempty"`
	Users                      []*UserReferenceWrapper `json:"users,omitempty"`
}

// ListSchedulesOptions represents options when listing schedules.
type ListSchedulesOptions struct {
	Limit  int    `url:"limit,omitempty"`
	More   bool   `url:"more,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Query  string `url:"query,omitempty"`
	Total  int    `url:"total,omitempty"`
}

// ListSchedulesResponse represents a list response of schedules.
type ListSchedulesResponse struct {
	Limit     int         `json:"limit,omitempty"`
	More      bool        `json:"more,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	Schedules []*Schedule `json:"schedules,omitempty"`
	Total     int         `json:"total,omitempty"`
}

// ListOnCallsOptions represents options when listing on calls.
type ListOnCallsOptions struct {
	ID    string `url:"id,omitempty"`
	Since string `url:"since,omitempty"`
	Until string `url:"until,omitempty"`
}

// ListOnCallsResponse represents a list response of on calls.
type ListOnCallsResponse struct {
	Users []*User `json:"users,omitempty"`
}

// ListOverridesOptions represents options when listing overrides.
type ListOverridesOptions struct {
	Editable bool   `url:"editable,omitempty"`
	ID       string `url:"id,omitempty"`
	Overflow bool   `url:"overflow,omitempty"`
	Since    string `url:"since,omitempty"`
	Until    string `url:"until,omitempty"`
}

// ListOverridesResponse represents a list response of schedules.
type ListOverridesResponse struct {
	Limit     int         `json:"limit,omitempty"`
	More      bool        `json:"more,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	Overrides []*Override `json:"overrides,omitempty"`
	Total     int         `json:"total,omitempty"`
}

// GetScheduleOptions represents options when retrieving a schedule.
type GetScheduleOptions struct {
	Since    string `url:"since,omitempty"`
	TimeZone string `url:"time_zone,omitempty"`
	Until    string `url:"until,omitempty"`
}

// CreateScheduleOptions represents options when creating a schedule.
type CreateScheduleOptions struct {
	Overflow bool `url:"overflow,omitempty"`
}

// UpdateScheduleOptions represents options when updating a schedule.
type UpdateScheduleOptions struct {
	Overflow bool `url:"overflow,omitempty"`
}

// SchedulePayload represents a schedule.
type SchedulePayload struct {
	Schedule *Schedule `json:"schedule,omitempty"`
}

// OverridePayload represents an override.
type OverridePayload struct {
	Override *Override `json:"override,omitempty"`
}

// List lists existing schedules.
func (s *ScheduleService) List(o *ListSchedulesOptions) (*ListSchedulesResponse, *Response, error) {
	u := "/schedules"
	v := new(ListSchedulesResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new schedule.
func (s *ScheduleService) Create(schedule *Schedule, o *CreateScheduleOptions) (*Schedule, *Response, error) {
	u := "/schedules"
	v := new(SchedulePayload)

	resp, err := s.client.newRequestDo("POST", u, o, &SchedulePayload{Schedule: schedule}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Schedule, resp, nil
}

// Delete removes an existing schedule.
func (s *ScheduleService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/schedules/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about a schedule.
func (s *ScheduleService) Get(id string, o *GetScheduleOptions) (*Schedule, *Response, error) {
	u := fmt.Sprintf("/schedules/%s", id)
	v := new(SchedulePayload)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Schedule, resp, nil
}

// Update updates an existing schedule.
func (s *ScheduleService) Update(id string, schedule *Schedule, o *UpdateScheduleOptions) (*Schedule, *Response, error) {
	u := fmt.Sprintf("/schedules/%s", id)
	v := new(SchedulePayload)

	resp, err := s.client.newRequestDo("PUT", u, o, &SchedulePayload{Schedule: schedule}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Schedule, resp, nil
}

// ListOnCalls lists all of the users on call in a given schedule for a given time range.
func (s *ScheduleService) ListOnCalls(scheduleID string, o *ListOnCallsOptions) (*ListOnCallsResponse, *Response, error) {
	u := fmt.Sprintf("/schedules/%s/users", scheduleID)
	v := new(ListOnCallsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// ListOverrides lists existing overrides.
func (s *ScheduleService) ListOverrides(scheduleID string, o *ListOverridesOptions) (*ListOverridesResponse, *Response, error) {
	u := fmt.Sprintf("/schedules/%s/overrides", scheduleID)
	v := new(ListOverridesResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// CreateOverride creates an override for a specific user covering the specified time range.
func (s *ScheduleService) CreateOverride(id string, override *Override) (*Override, *Response, error) {
	u := fmt.Sprintf("/schedules/%s/overrides", id)
	v := new(OverridePayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &OverridePayload{Override: override}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Override, resp, nil
}

// DeleteOverride deletes an override.
func (s *ScheduleService) DeleteOverride(id string, overrideID string) (*Response, error) {
	u := fmt.Sprintf("/schedules/%s/overrides/%s", id, overrideID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
