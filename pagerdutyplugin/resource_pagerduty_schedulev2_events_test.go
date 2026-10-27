package pagerduty

import (
	"context"
	"reflect"
	"testing"
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func scheduleV2TestTime(offset time.Duration) string {
	return time.Now().UTC().Add(offset).Truncate(time.Second).Format(time.RFC3339)
}

// testScheduleV2Event builds an event model with a weekly recurrence and one
// rotating member. An empty effectiveUntil leaves it null.
func testScheduleV2Event(id, name, effectiveSince, effectiveUntil, userID string) eventV2Model {
	recurrence, _ := types.ListValue(types.StringType, []attr.Value{types.StringValue("RRULE:FREQ=WEEKLY")})
	until := types.StringNull()
	if effectiveUntil != "" {
		until = types.StringValue(effectiveUntil)
	}
	return eventV2Model{
		ID:             types.StringValue(id),
		Name:           types.StringValue(name),
		StartTime:      types.StringValue("2026-01-05T13:00:00Z"),
		EndTime:        types.StringValue("2026-01-12T13:00:00Z"),
		EffectiveSince: types.StringValue(effectiveSince),
		EffectiveUntil: until,
		Recurrence:     recurrence,
		AssignmentStrategy: []assignmentStrategyV2Model{{
			Type:            types.StringValue("rotating_member_assignment_strategy"),
			ShiftsPerMember: types.Int64Value(1),
			Members:         []memberV2Model{{Type: types.StringValue("user_member"), UserID: types.StringValue(userID)}},
		}},
	}
}

func testScheduleV2APIEvent(evt eventV2Model) pagerduty.EventV3 {
	one := 1
	userID := evt.AssignmentStrategy[0].Members[0].UserID.ValueString()
	apiEvt := pagerduty.EventV3{
		ID:             evt.ID.ValueString(),
		Name:           evt.Name.ValueString(),
		StartTime:      pagerduty.EventTimeV3{DateTime: evt.StartTime.ValueString()},
		EndTime:        pagerduty.EventTimeV3{DateTime: evt.EndTime.ValueString()},
		EffectiveSince: evt.EffectiveSince.ValueString(),
		Recurrence:     []string{"RRULE:FREQ=WEEKLY"},
		AssignmentStrategy: pagerduty.AssignmentStrategyV3{
			Type:            "rotating_member_assignment_strategy",
			ShiftsPerMember: &one,
			Members:         []pagerduty.MemberV3{{Type: "user_member", UserID: &userID}},
		},
	}
	if !evt.EffectiveUntil.IsNull() {
		until := evt.EffectiveUntil.ValueString()
		apiEvt.EffectiveUntil = &until
	}
	return apiEvt
}

func TestScheduleV2EventHasEnded(t *testing.T) {
	cases := map[string]struct {
		until string
		want  bool
	}{
		"no effective_until": {"", false},
		"until in the past":  {scheduleV2TestTime(-time.Hour), true},
		"until in future":    {scheduleV2TestTime(time.Hour), false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			evt := testScheduleV2Event("E1", "Layer 1", scheduleV2TestTime(-48*time.Hour), tc.until, "U1")
			if got := eventHasEnded(evt); got != tc.want {
				t.Errorf("eventHasEnded() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestScheduleV2EquivalentEffectiveSince(t *testing.T) {
	past1 := scheduleV2TestTime(-48 * time.Hour)
	past2 := scheduleV2TestTime(-24 * time.Hour)
	future1 := scheduleV2TestTime(24 * time.Hour)
	future2 := scheduleV2TestTime(48 * time.Hour)
	cases := map[string]struct {
		a, b string
		want bool
	}{
		"same instant, different offsets": {"2099-01-01T08:00:00-05:00", "2099-01-01T13:00:00Z", true},
		"two past values":                 {past1, past2, true},
		"past and future":                 {past1, future1, false},
		"two future values":               {future1, future2, false},
		"unparseable":                     {"not a time", past1, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := equivalentEffectiveSince(types.StringValue(tc.a), types.StringValue(tc.b)); got != tc.want {
				t.Errorf("equivalentEffectiveSince(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestScheduleV2PairEvents(t *testing.T) {
	ctx := context.Background()
	past := scheduleV2TestTime(-10 * 24 * time.Hour)
	cutover := scheduleV2TestTime(-24 * time.Hour)
	future := scheduleV2TestTime(24 * time.Hour)

	ended := testScheduleV2Event("E1", "Layer 1", past, cutover, "U1")
	live := testScheduleV2Event("E2", "Layer 1 (expanded)", cutover, "", "U2")
	// The same event as live, with the effective_since of the ended event that
	// positional matching on read used to copy onto it.
	liveMisattributed := testScheduleV2Event("E2", "Layer 1 (expanded)", past, "", "U2")
	upcoming := testScheduleV2Event("E3", "Layer 1 (next)", future, "", "U3")

	cases := map[string]struct {
		desired, current []eventV2Model
		want             eventPairing
	}{
		"unchanged handover keeps both events": {
			desired: []eventV2Model{ended, live},
			current: []eventV2Model{ended, live},
			want:    eventPairing{current: []int{0, 1}},
		},
		"removing the ended event leaves the live event alone": {
			desired: []eventV2Model{live},
			current: []eventV2Model{ended, live},
			want:    eventPairing{current: []int{1}},
		},
		"removing the ended event with a misattributed effective_since": {
			desired: []eventV2Model{live},
			current: []eventV2Model{liveMisattributed},
			want:    eventPairing{current: []int{0}},
		},
		"changed event is paired in order": {
			desired: []eventV2Model{testScheduleV2Event("", "Layer 1 (expanded)", cutover, "", "U9")},
			current: []eventV2Model{ended, live},
			want:    eventPairing{current: []int{1}},
		},
		"inserting an event before an unchanged one creates it": {
			desired: []eventV2Model{upcoming, live},
			current: []eventV2Model{live},
			want:    eventPairing{current: []int{-1, 0}},
		},
		"removing a live event deletes it": {
			desired: []eventV2Model{live},
			current: []eventV2Model{live, upcoming},
			want:    eventPairing{current: []int{0}, stale: []int{1}},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := pairEvents(ctx, tc.desired, tc.current)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("pairEvents() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestScheduleV2FlattenKeepsEndedEventsAndMatchesByID(t *testing.T) {
	ctx := context.Background()
	past := scheduleV2TestTime(-10 * 24 * time.Hour)
	cutover := scheduleV2TestTime(-24 * time.Hour)

	ended := testScheduleV2Event("E1", "Layer 1", past, cutover, "U1")
	live := testScheduleV2Event("E2", "Layer 1 (expanded)", cutover, "", "U2")
	deleted := testScheduleV2Event("E3", "Removed in the web app", past, "", "U3")

	liveFromAPI := testScheduleV2APIEvent(live)
	// The API moves a past effective_since to the time of the request.
	liveFromAPI.EffectiveSince = scheduleV2TestTime(0)
	outOfBand := testScheduleV2Event("E4", "Added in the web app", past, "", "U4")

	schedule := &pagerduty.ScheduleV3{
		ID:       "S1",
		Name:     "Example",
		TimeZone: "UTC",
		Rotations: []pagerduty.RotationV3{{
			ID:     "R1",
			Events: []pagerduty.EventV3{testScheduleV2APIEvent(outOfBand), liveFromAPI},
		}},
	}
	stateRotations := []rotationV2Model{{
		ID:     types.StringValue("R1"),
		Events: []eventV2Model{ended, deleted, live},
	}}

	var diags diag.Diagnostics
	got := flattenScheduleV3(ctx, schedule, stateRotations, types.ListNull(types.StringType), &diags)
	if diags.HasError() {
		t.Fatalf("flattenScheduleV3() diagnostics: %v", diags)
	}

	events := got.Rotations[0].Events
	var ids []string
	for _, evt := range events {
		ids = append(ids, evt.ID.ValueString())
	}
	if want := []string{"E1", "E2", "E4"}; !reflect.DeepEqual(ids, want) {
		t.Fatalf("event IDs = %v, want %v", ids, want)
	}
	if !eventModelsEqual(ctx, events[0], ended) {
		t.Errorf("ended event = %+v, want it kept from state as %+v", events[0], ended)
	}
	if events[1].EffectiveSince.ValueString() != cutover {
		t.Errorf("live event effective_since = %s, want its own state value %s", events[1].EffectiveSince.ValueString(), cutover)
	}
}
