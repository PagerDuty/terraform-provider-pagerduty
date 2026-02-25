package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceScheduleV2 struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceScheduleV2)(nil)
	_ resource.ResourceWithImportState = (*resourceScheduleV2)(nil)
)

func (r *resourceScheduleV2) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_schedulev2"
}

func (r *resourceScheduleV2) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to create and manage on-call schedules using the PagerDuty v3 Schedules API. This is the new version of pagerduty_schedule and supports flexible rotations with per-event assignment strategies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the schedule.",
			},
			"time_zone": schema.StringAttribute{
				Required:    true,
				Description: "The time zone of the schedule (IANA format, e.g. 'America/New_York').",
				Validators:  []validator.String{validate.ValidTimeZone()},
			},
			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "A description of the schedule.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
		Blocks: map[string]schema.Block{
			"rotation": schema.ListNestedBlock{
				Description: "A rotation within the schedule. Each rotation can have multiple events defining on-call periods.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:      true,
							Description:   "The ID of the rotation.",
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
					},
					Blocks: map[string]schema.Block{
						"event": schema.ListNestedBlock{
							Description: "An event within the rotation defining an on-call period.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:      true,
										Description:   "The ID of the event.",
										PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
									},
									"name": schema.StringAttribute{
										Required:    true,
										Description: "The name of the event.",
									},
									"start_time": schema.StringAttribute{
										Required:    true,
										Description: "The shift start time with timezone offset (ISO-8601 format, e.g. '2024-01-01T09:00:00-05:00').",
									},
									"end_time": schema.StringAttribute{
										Required:    true,
										Description: "The shift end time with timezone offset (ISO-8601 format, e.g. '2024-01-01T17:00:00-05:00').",
									},
									"effective_since": schema.StringAttribute{
										Required:    true,
										Description: "When this event configuration starts producing shifts (ISO-8601 UTC). Must be a future time; the API will adjust past times to the current time.",
									},
									"effective_until": schema.StringAttribute{
										Optional:    true,
										Description: "When this event configuration stops producing shifts (ISO-8601 UTC). Null or omitted means indefinite.",
									},
									"recurrence": schema.ListAttribute{
										Required:    true,
										Description: "List of RRULE strings defining the recurrence pattern (RFC 5545, e.g. 'RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR').",
										ElementType: types.StringType,
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"assignment_strategy": schema.ListNestedBlock{
										Description: "Defines how on-call responsibility is assigned for this event.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Required:    true,
													Description: "The assignment strategy type.",
													Validators: []validator.String{
														stringvalidator.OneOf("user_assignment_strategy"),
													},
												},
											},
											Blocks: map[string]schema.Block{
												"member": schema.ListNestedBlock{
													Description: "A member to assign for on-call duty.",
													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The member type.",
																Validators: []validator.String{
																	stringvalidator.OneOf("user_member", "empty_member"),
																},
															},
															"user_id": schema.StringAttribute{
																Optional:    true,
																Description: "The obfuscated user ID. Required when type is 'user_member'.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *resourceScheduleV2) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceScheduleV2) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceScheduleV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleInput := buildScheduleV3Input(&model)
	log.Printf("[INFO] Creating PagerDuty v3 schedule: %s", scheduleInput.Name)

	var scheduleID string
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		schedule, err := r.client.CreateScheduleV3(ctx, scheduleInput)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		scheduleID = schedule.ID
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating PagerDuty v3 schedule", err.Error())
		return
	}

	model.ID = types.StringValue(scheduleID)

	// Create rotations and their events
	resp.Diagnostics.Append(r.createRotationsAndEvents(ctx, scheduleID, model.Rotations, &model.Rotations)...)
	if resp.Diagnostics.HasError() {
		// Attempt cleanup on partial failure
		_ = r.client.DeleteScheduleV3(ctx, scheduleID)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceScheduleV2) createRotationsAndEvents(ctx context.Context, scheduleID string, desired []rotationV2Model, result *[]rotationV2Model) diag.Diagnostics {
	var diags diag.Diagnostics
	updatedRotations := make([]rotationV2Model, 0, len(desired))

	for i, rotModel := range desired {
		var rotationID string
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			rotation, err := r.client.CreateRotationV3(ctx, scheduleID)
			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			rotationID = rotation.ID
			return nil
		})
		if err != nil {
			diags.AddError(fmt.Sprintf("Error creating rotation %d for v3 schedule %s", i, scheduleID), err.Error())
			return diags
		}

		rotModel.ID = types.StringValue(rotationID)
		updatedEvents := make([]eventV2Model, 0, len(rotModel.Events))

		for j, evtModel := range rotModel.Events {
			eventInput, d := buildEventV3Input(ctx, &evtModel)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			var createdEvent *pagerduty.EventV3
			err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
				e, err := r.client.CreateEventV3(ctx, scheduleID, rotationID, *eventInput)
				if err != nil {
					if util.IsBadRequestError(err) {
						return retry.NonRetryableError(err)
					}
					return retry.RetryableError(err)
				}
				createdEvent = e
				return nil
			})
			if err != nil {
				diags.AddError(fmt.Sprintf("Error creating event %d in rotation %s for v3 schedule %s", j, rotationID, scheduleID), err.Error())
				return diags
			}

			flattened := flattenEventV3(ctx, createdEvent, &diags)
			// The API normalizes past effective_since dates to current time. Preserve
			// the config value so the post-apply state matches the plan.
			flattened.EffectiveSince = evtModel.EffectiveSince
			updatedEvents = append(updatedEvents, flattened)
		}

		rotModel.Events = updatedEvents
		updatedRotations = append(updatedRotations, rotModel)
	}

	*result = updatedRotations
	return diags
}

func (r *resourceScheduleV2) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceScheduleV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID := state.ID.ValueString()
	log.Printf("[INFO] Reading PagerDuty v3 schedule: %s", scheduleID)

	var schedule *pagerduty.ScheduleV3
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		s, err := r.client.GetScheduleV3(ctx, scheduleID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				return nil
			}
			return retry.RetryableError(err)
		}
		schedule = s
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error reading PagerDuty v3 schedule %s", scheduleID), err.Error())
		return
	}

	if schedule == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// If the API doesn't return events inline with rotations, fetch them separately
	if len(schedule.Rotations) > 0 && allRotationsHaveNoEvents(schedule.Rotations) {
		for i, rot := range schedule.Rotations {
			fullRot, err := r.client.GetRotationV3(ctx, scheduleID, rot.ID)
			if err == nil && fullRot != nil {
				schedule.Rotations[i].Events = fullRot.Events
			}
		}
	}

	updatedState := flattenScheduleV3(ctx, schedule, state.Rotations, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *resourceScheduleV2) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceScheduleV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID := state.ID.ValueString()
	log.Printf("[INFO] Updating PagerDuty v3 schedule: %s", scheduleID)

	// Update schedule metadata if changed
	if !plan.Name.Equal(state.Name) || !plan.TimeZone.Equal(state.TimeZone) || !plan.Description.Equal(state.Description) {
		scheduleInput := buildScheduleV3Input(&plan)
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			_, err := r.client.UpdateScheduleV3(ctx, scheduleID, scheduleInput)
			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error updating PagerDuty v3 schedule %s", scheduleID), err.Error())
			return
		}
	}

	// Reconcile rotations: compare plan vs state by position
	resp.Diagnostics.Append(r.reconcileRotations(ctx, scheduleID, plan.Rotations, state.Rotations, &plan.Rotations)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceScheduleV2) reconcileRotations(ctx context.Context, scheduleID string, desired, current []rotationV2Model, result *[]rotationV2Model) diag.Diagnostics {
	var diags diag.Diagnostics
	updatedRotations := make([]rotationV2Model, 0, len(desired))

	// Handle rotations that exist in both desired and current (matched by position)
	minLen := len(desired)
	if len(current) < minLen {
		minLen = len(current)
	}

	for i := 0; i < minLen; i++ {
		desiredRot := desired[i]
		currentRot := current[i]
		rotationID := currentRot.ID.ValueString()

		desiredRot.ID = currentRot.ID

		// Reconcile events within this rotation
		diags.Append(r.reconcileEvents(ctx, scheduleID, rotationID, desiredRot.Events, currentRot.Events, &desiredRot.Events)...)
		if diags.HasError() {
			return diags
		}

		updatedRotations = append(updatedRotations, desiredRot)
	}

	// Delete extra rotations that are no longer needed (from the end)
	for i := minLen; i < len(current); i++ {
		rotationID := current[i].ID.ValueString()
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			err := r.client.DeleteRotationV3(ctx, scheduleID, rotationID)
			if err != nil {
				if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil && !util.IsNotFoundError(err) {
			diags.AddError(fmt.Sprintf("Error deleting rotation %s from v3 schedule %s", rotationID, scheduleID), err.Error())
			return diags
		}
	}

	// Create new rotations that are needed beyond what currently exists
	if len(desired) > len(current) {
		newRotations := desired[len(current):]
		newRotationResults := make([]rotationV2Model, 0, len(newRotations))
		diags.Append(r.createRotationsAndEvents(ctx, scheduleID, newRotations, &newRotationResults)...)
		if diags.HasError() {
			return diags
		}
		updatedRotations = append(updatedRotations, newRotationResults...)
	}

	*result = updatedRotations
	return diags
}

func (r *resourceScheduleV2) reconcileEvents(ctx context.Context, scheduleID, rotationID string, desired, current []eventV2Model, result *[]eventV2Model) diag.Diagnostics {
	var diags diag.Diagnostics
	updatedEvents := make([]eventV2Model, 0, len(desired))

	minLen := len(desired)
	if len(current) < minLen {
		minLen = len(current)
	}

	// Update events that exist in both (matched by position)
	for i := 0; i < minLen; i++ {
		desiredEvt := desired[i]
		currentEvt := current[i]
		eventID := currentEvt.ID.ValueString()

		eventInput, d := buildEventV3Input(ctx, &desiredEvt)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		var updatedEvent *pagerduty.EventV3
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			e, err := r.client.UpdateEventV3(ctx, scheduleID, rotationID, eventID, *eventInput)
			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			updatedEvent = e
			return nil
		})
		if err != nil {
			diags.AddError(fmt.Sprintf("Error updating event %s in rotation %s", eventID, rotationID), err.Error())
			return diags
		}

		flattened := flattenEventV3(ctx, updatedEvent, &diags)
		if diags.HasError() {
			return diags
		}
		// The API normalizes past effective_since dates to current time. Preserve
		// the plan value so the post-apply state matches the plan.
		flattened.EffectiveSince = desiredEvt.EffectiveSince
		updatedEvents = append(updatedEvents, flattened)
	}

	// Delete extra events no longer needed
	for i := minLen; i < len(current); i++ {
		eventID := current[i].ID.ValueString()
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			err := r.client.DeleteEventV3(ctx, scheduleID, rotationID, eventID)
			if err != nil {
				if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil && !util.IsNotFoundError(err) {
			diags.AddError(fmt.Sprintf("Error deleting event %s from rotation %s", eventID, rotationID), err.Error())
			return diags
		}
	}

	// Create new events that are needed beyond what currently exists
	for i := minLen; i < len(desired); i++ {
		desiredEvt := desired[i]
		eventInput, d := buildEventV3Input(ctx, &desiredEvt)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		var createdEvent *pagerduty.EventV3
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			e, err := r.client.CreateEventV3(ctx, scheduleID, rotationID, *eventInput)
			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			createdEvent = e
			return nil
		})
		if err != nil {
			diags.AddError(fmt.Sprintf("Error creating event %d in rotation %s", i, rotationID), err.Error())
			return diags
		}

		flattened := flattenEventV3(ctx, createdEvent, &diags)
		if diags.HasError() {
			return diags
		}
		// The API normalizes past effective_since dates to current time. Preserve
		// the plan value so the post-apply state matches the plan.
		flattened.EffectiveSince = desiredEvt.EffectiveSince
		updatedEvents = append(updatedEvents, flattened)
	}

	*result = updatedEvents
	return diags
}

func (r *resourceScheduleV2) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceScheduleV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID := state.ID.ValueString()
	log.Printf("[INFO] Deleting PagerDuty v3 schedule: %s", scheduleID)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		err := r.client.DeleteScheduleV3(ctx, scheduleID)
		if err != nil {
			if util.IsNotFoundError(err) {
				return nil
			}
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error deleting PagerDuty v3 schedule %s", scheduleID), err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *resourceScheduleV2) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	scheduleID := req.ID

	var schedule *pagerduty.ScheduleV3
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		s, err := r.client.GetScheduleV3(ctx, scheduleID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		schedule = s
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error importing PagerDuty v3 schedule %s", scheduleID), err.Error())
		return
	}

	// Fetch events for each rotation if not included inline
	if len(schedule.Rotations) > 0 && allRotationsHaveNoEvents(schedule.Rotations) {
		for i, rot := range schedule.Rotations {
			fullRot, err := r.client.GetRotationV3(ctx, scheduleID, rot.ID)
			if err == nil && fullRot != nil {
				schedule.Rotations[i].Events = fullRot.Events
			}
		}
	}

	state := flattenScheduleV3(ctx, schedule, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --- Model Types ---

type resourceScheduleV2Model struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	TimeZone    types.String      `tfsdk:"time_zone"`
	Description types.String      `tfsdk:"description"`
	Rotations   []rotationV2Model `tfsdk:"rotation"`
}

type rotationV2Model struct {
	ID     types.String   `tfsdk:"id"`
	Events []eventV2Model `tfsdk:"event"`
}

type eventV2Model struct {
	ID                 types.String                `tfsdk:"id"`
	Name               types.String                `tfsdk:"name"`
	StartTime          types.String                `tfsdk:"start_time"`
	EndTime            types.String                `tfsdk:"end_time"`
	EffectiveSince     types.String                `tfsdk:"effective_since"`
	EffectiveUntil     types.String                `tfsdk:"effective_until"`
	Recurrence         types.List                  `tfsdk:"recurrence"`
	AssignmentStrategy []assignmentStrategyV2Model `tfsdk:"assignment_strategy"`
}

type assignmentStrategyV2Model struct {
	Type    types.String    `tfsdk:"type"`
	Members []memberV2Model `tfsdk:"member"`
}

type memberV2Model struct {
	Type   types.String `tfsdk:"type"`
	UserID types.String `tfsdk:"user_id"`
}

// --- Build / Flatten Helpers ---

func buildScheduleV3Input(model *resourceScheduleV2Model) pagerduty.ScheduleV3Input {
	return pagerduty.ScheduleV3Input{
		Name:        model.Name.ValueString(),
		TimeZone:    model.TimeZone.ValueString(),
		Description: model.Description.ValueString(),
	}
}

func buildEventV3Input(ctx context.Context, model *eventV2Model) (*pagerduty.EventV3, diag.Diagnostics) {
	var diags diag.Diagnostics

	var recurrence []string
	diags.Append(model.Recurrence.ElementsAs(ctx, &recurrence, false)...)
	if diags.HasError() {
		return nil, diags
	}

	evt := &pagerduty.EventV3{
		Name:           model.Name.ValueString(),
		StartTime:      pagerduty.EventTimeV3{DateTime: model.StartTime.ValueString()},
		EndTime:        pagerduty.EventTimeV3{DateTime: model.EndTime.ValueString()},
		EffectiveSince: model.EffectiveSince.ValueString(),
		Recurrence:     recurrence,
	}

	if !model.EffectiveUntil.IsNull() && !model.EffectiveUntil.IsUnknown() && model.EffectiveUntil.ValueString() != "" {
		v := model.EffectiveUntil.ValueString()
		evt.EffectiveUntil = &v
	}

	if len(model.AssignmentStrategy) > 0 {
		as := model.AssignmentStrategy[0]
		strategy := pagerduty.AssignmentStrategyV3{
			Type: as.Type.ValueString(),
		}
		for _, m := range as.Members {
			member := pagerduty.MemberV3{
				Type: m.Type.ValueString(),
			}
			if !m.UserID.IsNull() && !m.UserID.IsUnknown() && m.UserID.ValueString() != "" {
				uid := m.UserID.ValueString()
				member.UserID = &uid
			}
			strategy.Members = append(strategy.Members, member)
		}
		evt.AssignmentStrategy = strategy
	}

	return evt, diags
}

func flattenScheduleV3(ctx context.Context, schedule *pagerduty.ScheduleV3, stateRotations []rotationV2Model, diags *diag.Diagnostics) resourceScheduleV2Model {
	// The v3 API normalizes "UTC" to "Etc/UTC". Normalize back to prevent perpetual plan diffs.
	tz := schedule.TimeZone
	if tz == "Etc/UTC" {
		tz = "UTC"
	}

	model := resourceScheduleV2Model{
		ID:          types.StringValue(schedule.ID),
		Name:        types.StringValue(schedule.Name),
		TimeZone:    types.StringValue(tz),
		Description: types.StringValue(schedule.Description),
	}

	rotations := make([]rotationV2Model, 0, len(schedule.Rotations))
	for i, rot := range schedule.Rotations {
		rotModel := rotationV2Model{
			ID: types.StringValue(rot.ID),
		}

		// Try to find the matching state rotation to preserve event ordering context
		var stateEvents []eventV2Model
		if i < len(stateRotations) {
			stateEvents = stateRotations[i].Events
		}

		events := make([]eventV2Model, 0, len(rot.Events))
		for j, evt := range rot.Events {
			evtModel := flattenEventV3(ctx, &evt, diags)
			if diags.HasError() {
				return model
			}

			// The v3 API normalizes start_time/end_time to UTC. Preserve the state value
			// when it represents the same instant to prevent perpetual diffs.
			if j < len(stateEvents) {
				st := stateEvents[j]
				if !st.StartTime.IsNull() && semanticallyEqualTime(st.StartTime.ValueString(), evtModel.StartTime.ValueString()) {
					evtModel.StartTime = st.StartTime
				}
				if !st.EndTime.IsNull() && semanticallyEqualTime(st.EndTime.ValueString(), evtModel.EndTime.ValueString()) {
					evtModel.EndTime = st.EndTime
				}
				// The API normalizes past effective_since dates to current time. Preserve
				// the state value to prevent perpetual plan diffs on subsequent refresh.
				if !st.EffectiveSince.IsNull() && !st.EffectiveSince.IsUnknown() {
					evtModel.EffectiveSince = st.EffectiveSince
				}
			}

			events = append(events, evtModel)
		}

		rotModel.Events = events
		rotations = append(rotations, rotModel)
	}

	model.Rotations = rotations
	return model
}

func flattenEventV3(ctx context.Context, evt *pagerduty.EventV3, diags *diag.Diagnostics) eventV2Model {
	recurrenceVals := make([]attr.Value, 0, len(evt.Recurrence))
	for _, r := range evt.Recurrence {
		recurrenceVals = append(recurrenceVals, types.StringValue(r))
	}
	recurrenceList, d := types.ListValue(types.StringType, recurrenceVals)
	diags.Append(d...)

	evtModel := eventV2Model{
		ID:             types.StringValue(evt.ID),
		Name:           types.StringValue(evt.Name),
		StartTime:      types.StringValue(evt.StartTime.DateTime),
		EndTime:        types.StringValue(evt.EndTime.DateTime),
		EffectiveSince: types.StringValue(evt.EffectiveSince),
		EffectiveUntil: types.StringNull(),
		Recurrence:     recurrenceList,
	}

	if evt.EffectiveUntil != nil && *evt.EffectiveUntil != "" {
		evtModel.EffectiveUntil = types.StringValue(*evt.EffectiveUntil)
	}

	// Flatten assignment strategy.
	// The v3 API normalizes "user_assignment_strategy" to "every_member_assignment_strategy"
	// in responses. Since the schema only supports "user_assignment_strategy", normalize back
	// to prevent perpetual plan diffs.
	if evt.AssignmentStrategy.Type != "" {
		stratType := evt.AssignmentStrategy.Type
		if stratType == "every_member_assignment_strategy" {
			stratType = "user_assignment_strategy"
		}
		as := assignmentStrategyV2Model{
			Type: types.StringValue(stratType),
		}
		members := make([]memberV2Model, 0, len(evt.AssignmentStrategy.Members))
		for _, m := range evt.AssignmentStrategy.Members {
			mem := memberV2Model{
				Type:   types.StringValue(m.Type),
				UserID: types.StringNull(),
			}
			if m.UserID != nil && *m.UserID != "" {
				mem.UserID = types.StringValue(*m.UserID)
			}
			members = append(members, mem)
		}
		as.Members = members
		evtModel.AssignmentStrategy = []assignmentStrategyV2Model{as}
	}

	return evtModel
}

// allRotationsHaveNoEvents returns true when all rotations lack event data,
// indicating the API may not have included events in the schedule response.
func allRotationsHaveNoEvents(rotations []pagerduty.RotationV3) bool {
	for _, r := range rotations {
		if len(r.Events) > 0 {
			return false
		}
	}
	return true
}

// semanticallyEqualTime returns true if two RFC3339 time strings represent the same instant.
// Used to prevent perpetual plan diffs when the v3 API normalizes times to UTC.
func semanticallyEqualTime(a, b string) bool {
	ta, errA := time.Parse(time.RFC3339, a)
	tb, errB := time.Parse(time.RFC3339, b)
	if errA != nil || errB != nil {
		return false
	}
	return ta.Equal(tb)
}
