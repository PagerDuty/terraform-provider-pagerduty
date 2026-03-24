package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// scheduleV3Headers returns the HTTP headers required by PagerDuty's v3 Schedules API.
func scheduleV3Headers() map[string]string {
	return map[string]string{
		"Accept": "application/json",
		// The v3 API requires Accept: application/json and the flexible-schedules early access header.
		"X-Early-Access": "flexible-schedules-early-access",
	}
}

// ScheduleV3 represents a schedule in PagerDuty's v3 API.
type ScheduleV3 struct {
	ID                 string            `json:"id,omitempty"`
	Type               string            `json:"type,omitempty"`
	Name               string            `json:"name"`
	TimeZone           string            `json:"time_zone"`
	Description        string            `json:"description,omitempty"`
	EscalationPolicies []APIObject       `json:"escalation_policies,omitempty"`
	Teams              []TeamReferenceV3 `json:"teams,omitempty"`
	Users              []APIObject       `json:"users,omitempty"`
	Rotations          []RotationV3      `json:"rotations,omitempty"`
	FinalSchedule      *FinalScheduleV3  `json:"final_schedule,omitempty"`
	HTTPCalURL         string            `json:"http_cal_url,omitempty"`
	WebCalURL          string            `json:"web_cal_url,omitempty"`
	Self               string            `json:"self,omitempty"`
	HTMLURL            string            `json:"html_url,omitempty"`
}

// RotationV3 represents a rotation within a v3 schedule.
type RotationV3 struct {
	ID      string    `json:"id,omitempty"`
	Type    string    `json:"type,omitempty"`
	Events  []EventV3 `json:"events,omitempty"`
	Self    string    `json:"self,omitempty"`
	HTMLURL string    `json:"html_url,omitempty"`
}

// EventTimeV3 represents a time field in a v3 schedule event.
// The v3 API uses an object {"date_time": "...", "time_zone": "..."} rather than a plain string.
type EventTimeV3 struct {
	DateTime string `json:"date_time"`
	TimeZone string `json:"time_zone,omitempty"`
}

// EventV3 represents an on-call event configuration within a rotation.
type EventV3 struct {
	ID                 string               `json:"id,omitempty"`
	Type               string               `json:"type,omitempty"`
	Name               string               `json:"name"`
	StartTime          EventTimeV3          `json:"start_time"`
	EndTime            EventTimeV3          `json:"end_time"`
	EffectiveSince     string               `json:"effective_since"`
	EffectiveUntil     *string              `json:"effective_until"`
	Recurrence         []string             `json:"recurrence"`
	AssignmentStrategy AssignmentStrategyV3 `json:"assignment_strategy"`
	Self               string               `json:"self,omitempty"`
	HTMLURL            string               `json:"html_url,omitempty"`
}

// AssignmentStrategyV3 defines how on-call responsibility is assigned within an event.
type AssignmentStrategyV3 struct {
	Type            string     `json:"type"`
	ShiftsPerMember *int       `json:"shifts_per_member,omitempty"`
	Members         []MemberV3 `json:"members"`
}

// TeamReferenceV3 represents a team associated with a v3 schedule.
type TeamReferenceV3 struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// MemberV3 represents a member in an assignment strategy or shift.
type MemberV3 struct {
	Type   string  `json:"type"`
	UserID *string `json:"user_id,omitempty"`
}

// FinalScheduleV3 holds computed on-call assignments for a requested time range.
// Only present when include[]=final_schedule is passed to GetScheduleV3 with since/until.
type FinalScheduleV3 struct {
	Type                       string                      `json:"type"`
	RenderedCoveragePercentage float64                     `json:"rendered_coverage_percentage"`
	ComputedShiftAssignments   []ComputedShiftAssignmentV3 `json:"computed_shift_assignments"`
}

// ComputedShiftAssignmentV3 represents a single computed on-call interval within a FinalScheduleV3.
type ComputedShiftAssignmentV3 struct {
	Type      string        `json:"type"`
	StartTime string        `json:"start_time"`
	EndTime   string        `json:"end_time"`
	Member    MemberV3      `json:"member"`
	Source    ShiftSourceV3 `json:"source"`
}

// ShiftSourceV3 identifies where a computed shift assignment originated.
type ShiftSourceV3 struct {
	Type       string `json:"type"`
	RotationID string `json:"rotation_id,omitempty"`
	ShiftID    string `json:"shift_id,omitempty"`
	OverrideID string `json:"override_id,omitempty"`
}

// ShiftAssignmentV3 assigns a member to a custom shift slot (response form includes an ID).
type ShiftAssignmentV3 struct {
	ID     string   `json:"id,omitempty"`
	Type   string   `json:"type"`
	Member MemberV3 `json:"member"`
}

// CustomShiftV3 represents an ad-hoc one-off shift outside of rotation events.
type CustomShiftV3 struct {
	ID          string              `json:"id,omitempty"`
	Type        string              `json:"type,omitempty"`
	StartTime   string              `json:"start_time"`
	EndTime     string              `json:"end_time"`
	Assignments []ShiftAssignmentV3 `json:"assignments"`
	Self        string              `json:"self,omitempty"`
	HTMLURL     string              `json:"html_url,omitempty"`
}

// CustomShiftInputV3 is the creation payload for a single custom shift.
type CustomShiftInputV3 struct {
	Type        string              `json:"type"`
	StartTime   string              `json:"start_time"`
	EndTime     string              `json:"end_time"`
	Assignments []ShiftAssignmentV3 `json:"assignments"`
}

// CustomShiftUpdateV3 is the update payload for a custom shift (all fields optional).
type CustomShiftUpdateV3 struct {
	StartTime   string              `json:"start_time,omitempty"`
	EndTime     string              `json:"end_time,omitempty"`
	Assignments []ShiftAssignmentV3 `json:"assignments,omitempty"`
}

// OverrideShiftV3 temporarily replaces a scheduled on-call member for a specific time period.
type OverrideShiftV3 struct {
	ID               string   `json:"id,omitempty"`
	Type             string   `json:"type,omitempty"`
	RotationID       string   `json:"rotation_id,omitempty"`
	CustomShiftID    string   `json:"custom_shift_id,omitempty"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	OverriddenMember MemberV3 `json:"overridden_member"`
	OverridingMember MemberV3 `json:"overriding_member"`
	Self             string   `json:"self,omitempty"`
	HTMLURL          string   `json:"html_url,omitempty"`
}

// OverrideShiftInputV3 is the creation payload for a single override.
type OverrideShiftInputV3 struct {
	Type             string   `json:"type"`
	RotationID       string   `json:"rotation_id,omitempty"`
	CustomShiftID    string   `json:"custom_shift_id,omitempty"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	OverriddenMember MemberV3 `json:"overridden_member"`
	OverridingMember MemberV3 `json:"overriding_member"`
}

// OverrideShiftUpdateV3 is the update payload for an override (all fields optional).
type OverrideShiftUpdateV3 struct {
	StartTime        string    `json:"start_time,omitempty"`
	EndTime          string    `json:"end_time,omitempty"`
	OverridingMember *MemberV3 `json:"overriding_member,omitempty"`
}

// ScheduleV3Input contains the mutable fields for creating or updating a v3 schedule.
type ScheduleV3Input struct {
	Name        string            `json:"name"`
	TimeZone    string            `json:"time_zone"`
	Description string            `json:"description,omitempty"`
	Teams       []TeamReferenceV3 `json:"teams,omitempty"`
}

// scheduleV3Payload is the request body shape for v3 schedule create/update.
// The rotations field must be present (even as []) per the v3 API validation.
type scheduleV3Payload struct {
	Name        string            `json:"name"`
	TimeZone    string            `json:"time_zone"`
	Description string            `json:"description,omitempty"`
	Teams       []TeamReferenceV3 `json:"teams,omitempty"`
	Rotations   []RotationV3      `json:"rotations"`
}

type createScheduleV3Request struct {
	Schedule scheduleV3Payload `json:"schedule"`
}

type updateScheduleV3Request struct {
	Schedule scheduleV3Payload `json:"schedule"`
}

type scheduleV3Response struct {
	Schedule ScheduleV3 `json:"schedule"`
}

// ListSchedulesV3Response is the response for listing v3 schedules.
type ListSchedulesV3Response struct {
	Schedules []APIObject `json:"schedules"`
	Limit     int         `json:"limit,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	More      bool        `json:"more,omitempty"`
}

// ListSchedulesV3Options are query parameters for listing v3 schedules.
type ListSchedulesV3Options struct {
	Query  string `url:"query,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// GetScheduleV3Options are optional query parameters for GetScheduleV3.
// Pass include[]=final_schedule with Since/Until to retrieve computed on-call assignments.
type GetScheduleV3Options struct {
	Since    string   `url:"since,omitempty"`
	Until    string   `url:"until,omitempty"`
	TimeZone string   `url:"time_zone,omitempty"`
	Overflow string   `url:"overflow,omitempty"`
	Include  []string `url:"include,omitempty,brackets"`
}

// GetRotationV3Options are optional query parameters for GetRotationV3.
type GetRotationV3Options struct {
	Since string `url:"since,omitempty"`
	Until string `url:"until,omitempty"`
}

// GetEventV3Options are optional query parameters for GetEventV3.
type GetEventV3Options struct {
	Since string `url:"since,omitempty"`
	Until string `url:"until,omitempty"`
}

// ListRotationsV3Options are query parameters for listing rotations.
type ListRotationsV3Options struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

// ListEventsV3Options are query parameters for listing events within a rotation.
type ListEventsV3Options struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

// ListCustomShiftsV3Options are query parameters for listing custom shifts.
// Since and Until are required by the API.
type ListCustomShiftsV3Options struct {
	Since    string `url:"since"`
	Until    string `url:"until"`
	TimeZone string `url:"time_zone,omitempty"`
	Overflow string `url:"overflow,omitempty"`
	Limit    int    `url:"limit,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

// ListOverridesV3Options are query parameters for listing overrides.
// Since and Until are required by the API.
type ListOverridesV3Options struct {
	Since    string `url:"since"`
	Until    string `url:"until"`
	TimeZone string `url:"time_zone,omitempty"`
	Overflow string `url:"overflow,omitempty"`
	Limit    int    `url:"limit,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

type rotationV3Response struct {
	Rotation RotationV3 `json:"rotation"`
}

type eventV3Response struct {
	Event EventV3 `json:"event"`
}

type createEventV3Request struct {
	Event EventV3 `json:"event"`
}

type updateEventV3Request struct {
	Event EventV3 `json:"event"`
}

// ListRotationsV3Response is the response for listing rotations.
type ListRotationsV3Response struct {
	Rotations []RotationV3 `json:"rotations"`
	Limit     int          `json:"limit,omitempty"`
	Offset    int          `json:"offset,omitempty"`
	More      bool         `json:"more,omitempty"`
}

// ListEventsV3Response is the response for listing events within a rotation.
type ListEventsV3Response struct {
	Events []EventV3 `json:"events"`
	Limit  int       `json:"limit,omitempty"`
	Offset int       `json:"offset,omitempty"`
	More   bool      `json:"more,omitempty"`
}

// ListCustomShiftsV3Response is the response for listing custom shifts.
type ListCustomShiftsV3Response struct {
	CustomShifts []CustomShiftV3 `json:"custom_shifts"`
	Limit        int             `json:"limit,omitempty"`
	Offset       int             `json:"offset,omitempty"`
	More         bool            `json:"more,omitempty"`
}

type createCustomShiftsV3Request struct {
	CustomShifts []CustomShiftInputV3 `json:"custom_shifts"`
}

type createCustomShiftsV3Response struct {
	CustomShifts []CustomShiftV3 `json:"custom_shifts"`
}

type customShiftV3Response struct {
	CustomShift CustomShiftV3 `json:"custom_shift"`
}

type updateCustomShiftV3Request struct {
	CustomShift CustomShiftUpdateV3 `json:"custom_shift"`
}

// ListOverridesV3Response is the response for listing overrides.
type ListOverridesV3Response struct {
	Overrides []OverrideShiftV3 `json:"overrides"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
	More      bool              `json:"more,omitempty"`
}

type createOverridesV3Request struct {
	Overrides []OverrideShiftInputV3 `json:"overrides"`
}

type createOverridesV3Response struct {
	Overrides []OverrideShiftV3 `json:"overrides"`
}

type overrideShiftV3Response struct {
	Override OverrideShiftV3 `json:"override"`
}

type updateOverrideV3Request struct {
	Override OverrideShiftUpdateV3 `json:"override"`
}

// buildPath appends non-empty query string to path.
func buildV3Path(base string, v interface{}) (string, error) {
	vals, err := query.Values(v)
	if err != nil {
		return "", err
	}
	if encoded := vals.Encode(); encoded != "" {
		return base + "?" + encoded, nil
	}
	return base, nil
}

// ListSchedulesV3 retrieves a paginated list of v3 schedules.
func (c *Client) ListSchedulesV3(ctx context.Context, o ListSchedulesV3Options) (*ListSchedulesV3Response, error) {
	path, err := buildV3Path("/v3/schedules", o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result ListSchedulesV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetScheduleV3 retrieves a v3 schedule by ID, including its rotations and events.
// Optionally pass GetScheduleV3Options with Include: []string{"final_schedule"} and Since/Until
// to also receive computed on-call assignments.
func (c *Client) GetScheduleV3(ctx context.Context, id string, opts ...GetScheduleV3Options) (*ScheduleV3, error) {
	var o GetScheduleV3Options
	if len(opts) > 0 {
		o = opts[0]
	}
	path, err := buildV3Path("/v3/schedules/"+id, o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result scheduleV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Schedule, nil
}

// CreateScheduleV3 creates a new v3 schedule with metadata only.
// Rotations and events must be added via separate API calls.
func (c *Client) CreateScheduleV3(ctx context.Context, s ScheduleV3Input) (*ScheduleV3, error) {
	d := createScheduleV3Request{Schedule: scheduleV3Payload{
		Name:        s.Name,
		TimeZone:    s.TimeZone,
		Description: s.Description,
		Teams:       s.Teams,
		Rotations:   []RotationV3{},
	}}

	resp, err := c.post(ctx, "/v3/schedules", d, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 schedule, HTTP status code: %d", resp.StatusCode)
	}

	var result scheduleV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Schedule, nil
}

// UpdateScheduleV3 updates a v3 schedule's metadata (name, time_zone, description).
func (c *Client) UpdateScheduleV3(ctx context.Context, id string, s ScheduleV3Input) (*ScheduleV3, error) {
	d := updateScheduleV3Request{Schedule: scheduleV3Payload{
		Name:        s.Name,
		TimeZone:    s.TimeZone,
		Description: s.Description,
		Teams:       s.Teams,
		Rotations:   []RotationV3{},
	}}

	resp, err := c.put(ctx, "/v3/schedules/"+id, d, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result scheduleV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Schedule, nil
}

// DeleteScheduleV3 soft-deletes a v3 schedule and all its rotations and events.
func (c *Client) DeleteScheduleV3(ctx context.Context, id string) error {
	// Use do() directly since delete() does not accept custom headers
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+id, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListRotationsV3 retrieves all rotations for a v3 schedule.
func (c *Client) ListRotationsV3(ctx context.Context, scheduleID string, o ListRotationsV3Options) (*ListRotationsV3Response, error) {
	path, err := buildV3Path("/v3/schedules/"+scheduleID+"/rotations", o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result ListRotationsV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateRotationV3 creates a new empty rotation for a v3 schedule.
// Events are added to the rotation via CreateEventV3.
func (c *Client) CreateRotationV3(ctx context.Context, scheduleID string) (*RotationV3, error) {
	resp, err := c.post(ctx, "/v3/schedules/"+scheduleID+"/rotations", map[string]interface{}{}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 rotation, HTTP status code: %d", resp.StatusCode)
	}

	var result rotationV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Rotation, nil
}

// GetRotationV3 retrieves a rotation by ID for a given schedule.
// Optionally pass GetRotationV3Options with Since/Until to filter the events returned.
func (c *Client) GetRotationV3(ctx context.Context, scheduleID, rotationID string, opts ...GetRotationV3Options) (*RotationV3, error) {
	var o GetRotationV3Options
	if len(opts) > 0 {
		o = opts[0]
	}
	path, err := buildV3Path("/v3/schedules/"+scheduleID+"/rotations/"+rotationID, o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result rotationV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Rotation, nil
}

// DeleteRotationV3 soft-deletes a rotation from a v3 schedule.
func (c *Client) DeleteRotationV3(ctx context.Context, scheduleID, rotationID string) error {
	// Use do() directly since delete() does not accept custom headers
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListEventsV3 retrieves all events for a v3 rotation, ordered by start time.
func (c *Client) ListEventsV3(ctx context.Context, scheduleID, rotationID string, o ListEventsV3Options) (*ListEventsV3Response, error) {
	path, err := buildV3Path("/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events", o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result ListEventsV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateEventV3 creates a new event within a v3 rotation.
func (c *Client) CreateEventV3(ctx context.Context, scheduleID, rotationID string, e EventV3) (*EventV3, error) {
	resp, err := c.post(ctx, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events", createEventV3Request{Event: e}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 event, HTTP status code: %d", resp.StatusCode)
	}

	var result eventV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Event, nil
}

// GetEventV3 retrieves a single event from a v3 rotation.
// Optionally pass GetEventV3Options with Since/Until to filter the event's computed data.
func (c *Client) GetEventV3(ctx context.Context, scheduleID, rotationID, eventID string, opts ...GetEventV3Options) (*EventV3, error) {
	var o GetEventV3Options
	if len(opts) > 0 {
		o = opts[0]
	}
	path, err := buildV3Path("/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events/"+eventID, o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result eventV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Event, nil
}

// UpdateEventV3 updates an event within a v3 rotation.
func (c *Client) UpdateEventV3(ctx context.Context, scheduleID, rotationID, eventID string, e EventV3) (*EventV3, error) {
	resp, err := c.put(ctx, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events/"+eventID, updateEventV3Request{Event: e}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result eventV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Event, nil
}

// DeleteEventV3 deletes an event from a v3 rotation.
func (c *Client) DeleteEventV3(ctx context.Context, scheduleID, rotationID, eventID string) error {
	// Use do() directly since delete() does not accept custom headers
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+scheduleID+"/rotations/"+rotationID+"/events/"+eventID, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListCustomShiftsV3 retrieves custom shifts for a schedule within the given time range.
// Since and Until in the options are required by the API.
func (c *Client) ListCustomShiftsV3(ctx context.Context, scheduleID string, o ListCustomShiftsV3Options) (*ListCustomShiftsV3Response, error) {
	path, err := buildV3Path("/v3/schedules/"+scheduleID+"/custom_shifts", o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result ListCustomShiftsV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateCustomShiftsV3 creates one or more ad-hoc custom shifts for a schedule.
func (c *Client) CreateCustomShiftsV3(ctx context.Context, scheduleID string, shifts []CustomShiftInputV3) ([]CustomShiftV3, error) {
	resp, err := c.post(ctx, "/v3/schedules/"+scheduleID+"/custom_shifts", createCustomShiftsV3Request{CustomShifts: shifts}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 custom shifts, HTTP status code: %d", resp.StatusCode)
	}

	var result createCustomShiftsV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.CustomShifts, nil
}

// GetCustomShiftV3 retrieves a single custom shift by ID.
func (c *Client) GetCustomShiftV3(ctx context.Context, scheduleID, customShiftID string) (*CustomShiftV3, error) {
	resp, err := c.get(ctx, "/v3/schedules/"+scheduleID+"/custom_shifts/"+customShiftID, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result customShiftV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.CustomShift, nil
}

// UpdateCustomShiftV3 updates an existing custom shift.
// If the shift has already started, only EndTime can be modified.
func (c *Client) UpdateCustomShiftV3(ctx context.Context, scheduleID, customShiftID string, update CustomShiftUpdateV3) (*CustomShiftV3, error) {
	resp, err := c.put(ctx, "/v3/schedules/"+scheduleID+"/custom_shifts/"+customShiftID, updateCustomShiftV3Request{CustomShift: update}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result customShiftV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.CustomShift, nil
}

// DeleteCustomShiftV3 deletes a custom shift by ID.
func (c *Client) DeleteCustomShiftV3(ctx context.Context, scheduleID, customShiftID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+scheduleID+"/custom_shifts/"+customShiftID, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListOverridesV3 retrieves overrides for a schedule within the given time range.
// Since and Until in the options are required by the API.
func (c *Client) ListOverridesV3(ctx context.Context, scheduleID string, o ListOverridesV3Options) (*ListOverridesV3Response, error) {
	path, err := buildV3Path("/v3/schedules/"+scheduleID+"/overrides", o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, path, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result ListOverridesV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateOverridesV3 creates one or more overrides that temporarily replace scheduled on-call members.
// Each override must reference either a RotationID or a CustomShiftID (not both).
func (c *Client) CreateOverridesV3(ctx context.Context, scheduleID string, overrides []OverrideShiftInputV3) ([]OverrideShiftV3, error) {
	resp, err := c.post(ctx, "/v3/schedules/"+scheduleID+"/overrides", createOverridesV3Request{Overrides: overrides}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create v3 overrides, HTTP status code: %d", resp.StatusCode)
	}

	var result createOverridesV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Overrides, nil
}

// GetOverrideV3 retrieves a single override by ID.
func (c *Client) GetOverrideV3(ctx context.Context, scheduleID, overrideID string) (*OverrideShiftV3, error) {
	resp, err := c.get(ctx, "/v3/schedules/"+scheduleID+"/overrides/"+overrideID, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result overrideShiftV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Override, nil
}

// UpdateOverrideV3 updates an existing override.
func (c *Client) UpdateOverrideV3(ctx context.Context, scheduleID, overrideID string, update OverrideShiftUpdateV3) (*OverrideShiftV3, error) {
	resp, err := c.put(ctx, "/v3/schedules/"+scheduleID+"/overrides/"+overrideID, updateOverrideV3Request{Override: update}, scheduleV3Headers())
	if err != nil {
		return nil, err
	}

	var result overrideShiftV3Response
	if err = c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Override, nil
}

// DeleteOverrideV3 deletes an override by ID.
func (c *Client) DeleteOverrideV3(ctx context.Context, scheduleID, overrideID string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/v3/schedules/"+scheduleID+"/overrides/"+overrideID, nil, scheduleV3Headers())
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
