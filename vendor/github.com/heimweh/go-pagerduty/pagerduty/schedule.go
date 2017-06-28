package pagerduty

import "fmt"

// ScheduleService handles the communication with schedule
// related methods of the PagerDuty API.
type ScheduleService service

// Schedule represents a schedule.
type Schedule struct {
	Description          string                       `json:"description,omitempty"`
	EscalationPolicies   []*EscalationPolicyReference `json:"escalation_policies,omitempty"`
	FinalSchedule        *SubSchedule                 `json:"final_schedule,omitempty"`
	HTMLURL              string                       `json:"html_url,omitempty"`
	ID                   string                       `json:"id,omitempty"`
	Name                 string                       `json:"name,omitempty"`
	OverridesSubSchedule *SubSchedule                 `json:"overrides_subschedule,omitempty"`
	Schedule             *Schedule                    `json:"schedule,omitempty"`
	ScheduleLayers       []*ScheduleLayer             `json:"schedule_layers,omitempty"`
	Self                 string                       `json:"self,omitempty"`
	Summary              string                       `json:"summary,omitempty"`
	TimeZone             string                       `json:"time_zone,omitempty"`
	Type                 string                       `json:"type,omitempty"`
	Users                []*UserReference             `json:"users,omitempty"`
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
	*Pagination
	Query string `url:"query,omitempty"`
}

// ListSchedulesResponse represents a list response of schedules.
type ListSchedulesResponse struct {
	*Pagination
	Schedules []*Schedule `json:"schedules,omitempty"`
}

// GetScheduleOptions represents options when retrieving a schedule.
type GetScheduleOptions struct {
	Since    string `url:"since,omitempty"`
	TimeZone string `url:"time_zone,omitempty"`
	Until    string `url:"until,omitempty"`
}

// UpdateScheduleOptions represents options when updating a schedule.
type UpdateScheduleOptions struct {
	Overflow bool `url:"overflow,omitempty"`
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
func (s *ScheduleService) Create(schedule *Schedule) (*Schedule, *Response, error) {
	u := "/schedules"
	v := new(Schedule)

	resp, err := s.client.newRequestDo("POST", u, nil, &Schedule{Schedule: schedule}, &v)
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
	v := new(Schedule)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Schedule, resp, nil
}

// Update updates an existing schedule.
func (s *ScheduleService) Update(id string, schedule *Schedule) (*Schedule, *Response, error) {
	u := fmt.Sprintf("/schedules/%s", id)
	v := new(Schedule)

	resp, err := s.client.newRequestDo("PUT", u, nil, &Schedule{Schedule: schedule}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Schedule, resp, nil
}
