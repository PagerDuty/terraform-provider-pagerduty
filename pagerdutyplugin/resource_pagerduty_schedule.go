package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/apiutil"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/rangetypes"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/tztypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceSchedule struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceSchedule)(nil)
	_ resource.ResourceWithImportState = (*resourceSchedule)(nil)
)

func (r *resourceSchedule) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_schedule"
}

func (r *resourceSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{Optional: true},
			"time_zone": schema.StringAttribute{
				Required:   true,
				CustomType: tztypes.StringType{},
			},
			"overflow": schema.BoolAttribute{Optional: true},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Managed by terraform"),
			},
			"layer": schema.ListAttribute{
				Required:    true,
				ElementType: scheduleLayerObjectType,
			},
		},
	}
}

var scheduleLayerObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                           types.StringType,
		"name":                         types.StringType,
		"start":                        tztypes.RFC3339Type{}, // required, rfc3339, suppressScheduleLayerStartDiff,
		"end":                          tztypes.RFC3339Type{},
		"rotation_virtual_start":       tztypes.RFC3339Type{},
		"rotation_turn_length_seconds": scheduleLayerRotationTurnLengthSecondsType, // required
		"users":                        types.ListType{ElemType: types.StringType}, // required, min 1
		"rendered_coverage_percentage": types.StringType,
		"restriction": types.ListType{
			ElemType: scheduleLayerRestrictionObjectType,
		},
		"teams": types.ListType{
			ElemType: types.StringType,
		},
		"final_schedule": types.ListType{
			ElemType: scheduleFinalScheduleObjectType,
		},
	},
}

type scheduleLayerModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Start                      types.String `tfsdk:"start"`
	End                        types.String `tfsdk:"end"`
	RenderedCoveragePercentage types.String `tfsdk:"rendered_coverage_percentage"`
	RotationTurnLengthSeconds  types.Int64  `tfsdk:"rotation_turn_length_seconds"`
	RotationVirtualStart       types.String `tfsdk:"rotation_virtual_start"`
	Restriction                types.List   `tfsdk:"restriction"`
	Users                      types.List   `tfsdk:"users"`
	Teams                      types.List   `tfsdk:"teams"`
	FinalSchedule              types.List   `tfsdk:"final_schedule"`
}

var scheduleLayerRestrictionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":              types.StringType, // required. "daily_restriction", "weekly_restriction",
		"start_time_of_day": types.StringType, // required. validation.StringMatch(regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]`), "must be of 00:00:00 format"),
		"start_day_of_week": types.Int64Type,  // required. [1,7]
		"duration_seconds":  types.Int64Type,  // required. [1, 7*24*3600 - 1]
	},
}

var scheduleFinalScheduleObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":                         types.StringType,
		"rendered_coverage_percentage": types.StringType,
	},
}

func (r *resourceSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceScheduleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutySchedule(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Creating PagerDuty schedule %s", plan.Name)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		// TODO: add overflow query param
		response, err := r.client.CreateScheduleWithContext(ctx, plan)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		plan.ID = response.ID
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty schedule %s", plan.Name),
			err.Error(),
		)
		return
	}

	schedule, err := fetchPagerdutySchedule(ctx, r.client, plan.ID, RetryNotFound)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty schedule %s", plan.ID),
			err.Error(),
		)
		return
	}
	model = flattenSchedule(schedule, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty schedule %s", id)

	schedule, err := fetchPagerdutySchedule(ctx, r.client, id.ValueString(), !RetryNotFound)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty schedule %s", id),
			err.Error(),
		)
		return
	}
	state := flattenSchedule(schedule, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceSchedule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateModel resourceScheduleModel
	var planModel resourceScheduleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateModel)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := buildPagerdutySchedule(ctx, &stateModel, &resp.Diagnostics)
	plan := buildPagerdutySchedule(ctx, &planModel, &resp.Diagnostics)
	log.Printf("[INFO] Updating PagerDuty schedule %s", plan.ID)

	// if !reflect.DeepEqual(state.ScheduleLayers, plan.ScheduleLayers) {
	for _, stateLayer := range state.ScheduleLayers {
		found := false
		for _, planLayer := range plan.ScheduleLayers {
			if stateLayer.ID == planLayer.ID {
				found = true
			}
		}
		if !found {
			stateLayer.End = time.Now().UTC().String()
			plan.ScheduleLayers = append(plan.ScheduleLayers, stateLayer)
		}
	}
	// }

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		// TODO: add overflow query param
		schedule, err := r.client.UpdateScheduleWithContext(ctx, plan.ID, plan)
		if err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		planModel = flattenSchedule(schedule, &resp.Diagnostics)
		return nil
	})
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty schedule %s", plan.Name),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
}

func (r *resourceSchedule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty schedule %s", id)

	isScheduleUsedByEP := false
	isScheduleWithOpenOrphanIncidents := false

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		if err := r.client.DeleteScheduleWithContext(ctx, id.ValueString()); err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}

			isScheduleUsedByEP = strings.Contains(err.Error(), "Schedule can't be deleted if it's being used by escalation policies")
			isScheduleWithOpenOrphanIncidents = strings.Contains(err.Error(), "Schedule can't be deleted if it's being used by an escalation policy snapshot with open incidents")
			if isScheduleUsedByEP || isScheduleWithOpenOrphanIncidents {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil && !util.IsNotFoundError(err) {
		deleteErr := err

		schedule, err := fetchPagerdutySchedule(ctx, r.client, id.ValueString(), !RetryNotFound)
		if err != nil {
			if util.IsNotFoundError(err) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading PagerDuty schedule %s", id.ValueString()),
				err.Error(),
			)
			return
		}

		// When isScheduleWithOpenOrphanIncidents just return the error, but in the case of
		// isScheduleUsedByEP we need to check if there is open incidents before prompting a
		// request to delete the escalation policies.
		if !isScheduleWithOpenOrphanIncidents && isScheduleUsedByEP {
			incidents, d := fetchPagerdutyIncidentsOpenWithSchedule(ctx, r.client, schedule)
			if resp.Diagnostics.Append(d...); d.HasError() {
				return
			}

			msg := deleteErr.Error()
			if len(incidents) > 0 {
				msg = msgForScheduleWithOpenIncidents(schedule, incidents)
			} else if len(schedule.EscalationPolicies) > 0 {
				msg = msgForScheduleUsedByEP(schedule)
			}

			resp.Diagnostics.AddError(
				fmt.Sprintf("Schedule %q couldn't be deleted", schedule.ID),
				msg,
			)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty schedule %s", id),
			deleteErr.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *resourceSchedule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceSchedule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceScheduleModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	TimeZone    types.String `tfsdk:"time_zone"`
	Layer       types.List   `tfsdk:"layer"`
	Description types.String `tfsdk:"description"`
	Overflow    types.Bool   `tfsdk:"overflow"`
}

func fetchPagerdutySchedule(ctx context.Context, client *pagerduty.Client, id string, retryNotFound bool) (*pagerduty.Schedule, error) {
	var schedule *pagerduty.Schedule

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var err error
		o := pagerduty.GetScheduleOptions{}
		schedule, err = client.GetScheduleWithContext(ctx, id, o)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})

	return schedule, err
}

func buildPagerdutySchedule(ctx context.Context, model *resourceScheduleModel, diags *diag.Diagnostics) pagerduty.Schedule {
	return pagerduty.Schedule{
		Description:    model.Description.ValueString(),
		Name:           model.Name.ValueString(),
		ScheduleLayers: buildScheduleLayers(ctx, model.Layer, diags),
		Teams:          buildScheduleTeams(ctx, model.Layer, diags),
		TimeZone:       model.TimeZone.ValueString(),
	}
}

func buildScheduleLayers(ctx context.Context, list types.List, diags *diag.Diagnostics) []pagerduty.ScheduleLayer {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var target []scheduleLayerModel
	d := list.ElementsAs(ctx, &target, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	scheduleLayers := make([]pagerduty.ScheduleLayer, 0, len(target))
	for _, item := range target {
		// This is a temporary fix to prevent getting back the wrong rotation_virtual_start time.
		// The background here is that if a user specifies a rotation_virtual_start time to be:
		// "2017-09-01T10:00:00+02:00" the API returns back "2017-09-01T12:00:00+02:00".
		// With this fix in place, we get the correct rotation_virtual_start time, thus
		// eliminating the diff issues we've been seeing in the past.
		// This has been confirmed working by PagerDuty support.
		rvs, err := util.TimeToUTC(item.RotationVirtualStart.ValueString())
		if err != nil {
			diags.AddAttributeError(
				path.Root("rotation_virtual_start"),
				"Cannot convert to UTC",
				err.Error(),
			)
			return nil
		}

		layer := pagerduty.ScheduleLayer{
			APIObject: pagerduty.APIObject{
				ID: item.ID.ValueString(),
			},
			Name:                      item.Name.ValueString(),
			Start:                     item.Start.ValueString(),
			End:                       item.End.ValueString(),
			RotationVirtualStart:      rvs.Format(time.RFC3339),
			RotationTurnLengthSeconds: uint(item.RotationTurnLengthSeconds.ValueInt64()),
		}

		userList := buildPagerdutyAPIObjectFromIDs(ctx, item.Users, "user", diags)
		for _, user := range userList {
			layer.Users = append(layer.Users, pagerduty.UserReference{User: user})
		}

		var restrictionList []struct {
			Type            types.String `tfsdk:"type"`
			StartTimeOfDay  types.String `tfsdk:"start_time_of_day"`
			StartDayOfWeek  types.Int64  `tfsdk:"start_day_of_week"`
			DurationSeconds types.Int64  `tfsdk:"duration_seconds"`
		}
		diags.Append(item.Restriction.ElementsAs(ctx, &restrictionList, false)...)
		if diags.HasError() {
			return nil
		}

		for _, restriction := range restrictionList {
			layer.Restrictions = append(layer.Restrictions, pagerduty.Restriction{
				Type:            restriction.Type.ValueString(),
				StartTimeOfDay:  restriction.StartTimeOfDay.ValueString(),
				StartDayOfWeek:  uint(restriction.StartDayOfWeek.ValueInt64()),
				DurationSeconds: uint(restriction.DurationSeconds.ValueInt64()),
			})
		}
		scheduleLayers = append(scheduleLayers, layer)
	}

	return scheduleLayers
}

func buildScheduleTeams(ctx context.Context, layerList types.List, diags *diag.Diagnostics) []pagerduty.APIObject {
	var target []scheduleLayerModel
	d := layerList.ElementsAs(ctx, &target, true)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}
	obj := target[0]
	return buildPagerdutyAPIObjectFromIDs(ctx, obj.Teams, "team_reference", diags)
}

func flattenAPIObjectIDs(objects []pagerduty.APIObject) types.List {
	elements := make([]attr.Value, 0, len(objects))
	for _, obj := range objects {
		elements = append(elements, types.StringValue(obj.ID))
	}
	return types.ListValueMust(types.StringType, elements)
}

func flattenSchedule(response *pagerduty.Schedule, diags *diag.Diagnostics) resourceScheduleModel {
	model := resourceScheduleModel{
		ID:          types.StringValue(response.ID),
		Name:        types.StringValue(response.Name),
		TimeZone:    types.StringValue(response.TimeZone),
		Description: types.StringValue(response.Description),
		Layer: flattenScheduleLayers(
			response.ScheduleLayers, response.Teams, response.FinalSchedule, diags,
		),
	}
	return model
}

func flattenFinalSchedule(response pagerduty.ScheduleLayer, diags *diag.Diagnostics) types.List {
	obj, d := types.ObjectValue(scheduleFinalScheduleObjectType.AttrTypes, map[string]attr.Value{
		"name": types.StringValue(response.Name),
		"rendered_coverage_percentage": types.StringValue(util.RenderRoundedPercentage(
			response.RenderedCoveragePercentage,
		)),
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(scheduleFinalScheduleObjectType)
	}

	list, d := types.ListValue(scheduleFinalScheduleObjectType, []attr.Value{obj})
	diags.Append(d...)
	return list
}

func flattenScheduleLayers(scheduleLayers []pagerduty.ScheduleLayer, teams []pagerduty.APIObject, finalSchedule pagerduty.ScheduleLayer, diags *diag.Diagnostics) types.List {
	var elements []attr.Value
nextLayer:
	for _, layer := range scheduleLayers {
		// A schedule layer can never be removed but it can be ended.
		// Here we check each layer and if it has been ended we don't
		// read it back because it's not relevant anymore.
		if layer.End != "" {
			end, err := util.TimeToUTC(layer.End)
			if err != nil {
				diags.AddError(err.Error(), "")
				continue
			}
			if time.Now().UTC().After(end) {
				continue
			}
		}

		usersElems := make([]attr.Value, 0, len(layer.Users))
		for _, u := range layer.Users {
			usersElems = append(usersElems, types.StringValue(u.User.ID))
		}
		users, d := types.ListValue(types.StringType, usersElems)
		if d.HasError() {
			continue
		}

		restrictionsElems := make([]attr.Value, 0, len(layer.Restrictions))
		for _, r := range layer.Restrictions {
			sdow := types.Int64Null()
			if r.StartDayOfWeek > 0 {
				sdow = types.Int64Value(int64(r.StartDayOfWeek))
			}
			rst, d := types.ObjectValue(scheduleLayerRestrictionObjectType.AttrTypes, map[string]attr.Value{
				"duration_seconds":  types.Int64Value(int64(r.DurationSeconds)),
				"start_time_of_day": types.StringValue(r.StartTimeOfDay),
				"type":              types.StringValue(r.Type),
				"start_day_of_week": sdow,
			})
			if d.HasError() {
				continue nextLayer
			}
			restrictionsElems = append(restrictionsElems, rst)
		}
		restrictions, d := types.ListValue(scheduleLayerRestrictionObjectType, restrictionsElems)
		if d.HasError() {
			continue
		}

		obj, d := types.ObjectValue(scheduleLayerObjectType.AttrTypes, map[string]attr.Value{
			"id":                           types.StringValue(layer.ID),
			"name":                         types.StringValue(layer.Name),
			"end":                          tztypes.NewRFC3339Value(layer.End),
			"start":                        tztypes.NewRFC3339Value(layer.Start),
			"rotation_virtual_start":       tztypes.NewRFC3339Value(layer.RotationVirtualStart),
			"rotation_turn_length_seconds": types.Int64Value(int64(layer.RotationTurnLengthSeconds)),
			"rendered_coverage_percentage": types.StringValue(util.RenderRoundedPercentage(layer.RenderedCoveragePercentage)),
			"users":                        users,
			"restriction":                  restrictions,
			"teams":                        flattenAPIObjectIDs(teams),
			"final_schedule":               flattenFinalSchedule(finalSchedule, diags),
		})
		diags.Append(d...)
		if d.HasError() {
			continue
		}

		elements = append(elements, obj)
	}

	reversedElems := make([]attr.Value, 0, len(elements))
	for i, l := 0, len(elements); i < l; i++ {
		reversedElems = append(reversedElems, elements[l-i-1])
	}

	list, d := types.ListValue(scheduleLayerObjectType, reversedElems)
	diags.Append(d...)
	if d.HasError() {
		return types.ListNull(scheduleLayerObjectType)
	}

	return list
}

func msgForScheduleUsedByEP(schedule *pagerduty.Schedule) string {
	var links []string
	for _, ep := range schedule.EscalationPolicies {
		links = append(links, fmt.Sprintf("\t* %s", ep.HTMLURL))
	}
	return fmt.Sprintf(
		"Please remove this Schedule from the following Escalation Policies in order to unblock the Schedule removal:\n"+
			"%s\n"+
			"After completing, come back to continue with the destruction of Schedule.",
		strings.Join(links, "\n"),
	)
}

func msgForScheduleWithOpenIncidents(schedule *pagerduty.Schedule, incidents []pagerduty.Incident) string {
	links := make([]string, 0, len(incidents))
	for _, inc := range incidents {
		links = append(links, fmt.Sprintf("\t* %s", inc.HTMLURL))
	}
	return fmt.Sprintf(
		"Before destroying Schedule %q you must first resolve or reassign "+
			"the following incidents related with the Escalation Policies using "+
			"this Schedule:\n%s",
		schedule.ID, strings.Join(links, "\n"),
	)
}

func fetchPagerdutyIncidentsOpenWithSchedule(ctx context.Context, client *pagerduty.Client, schedule *pagerduty.Schedule) ([]pagerduty.Incident, diag.Diagnostics) {
	var diags diag.Diagnostics

	var incidents []pagerduty.Incident

	err := apiutil.All(ctx, func(offset int) (bool, error) {
		resp, err := client.ListIncidentsWithContext(ctx, pagerduty.ListIncidentsOptions{
			DateRange: "all",
			Statuses:  []string{"triggered", "acknowledged"},
			Limit:     apiutil.Limit,
			Offset:    uint(offset),
		})
		if err != nil {
			return false, err
		}

		incidents = append(incidents, resp.Incidents...)
		return resp.More, nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading PagerDuty incidents for schedule %s", schedule.ID),
			err.Error(),
		)
		return nil, diags
	}

	db := make(map[string]struct{})
	for _, ep := range schedule.EscalationPolicies {
		db[ep.ID] = struct{}{}
	}

	var output []pagerduty.Incident
	for _, inc := range incidents {
		if _, ok := db[inc.EscalationPolicy.ID]; ok {
			output = append(output, inc)
		}
	}

	return output, diags
}

var (
	scheduleLayerRotationTurnLengthSecondsType = rangetypes.Int64Type{Start: 3600, End: 365 * 24 * 3600}
)

/*
func resourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
			ln := diff.Get("layer.#").(int)
			for li := 0; li <= ln; li++ {
				rn := diff.Get(fmt.Sprintf("layer.%d.restriction.#", li)).(int)
				for ri := 0; ri <= rn; ri++ {
					t := diff.Get(fmt.Sprintf("layer.%d.restriction.%d.type", li, ri)).(string)
					isStartDayOfWeekSetWhenDailyRestrictionType := t == "daily_restriction" && diff.Get(fmt.Sprintf("layer.%d.restriction.%d.start_day_of_week", li, ri)).(int) != 0
					if isStartDayOfWeekSetWhenDailyRestrictionType {
						return fmt.Errorf("start_day_of_week must only be set for a weekly_restriction schedule restriction type")
					}
					isStartDayOfWeekNotSetWhenWeeklyRestrictionType := t == "weekly_restriction" && diff.Get(fmt.Sprintf("layer.%d.restriction.%d.start_day_of_week", li, ri)).(int) == 0
					if isStartDayOfWeekNotSetWhenWeeklyRestrictionType {
						return fmt.Errorf("start_day_of_week must be set for a weekly_restriction schedule restriction type")
					}
					ds := diff.Get(fmt.Sprintf("layer.%d.restriction.%d.duration_seconds", li, ri)).(int)
					if t == "daily_restriction" && ds >= 3600*24 {
						return fmt.Errorf("duration_seconds for a daily_restriction schedule restriction type must be shorter than a day")
					}
				}
			}
			return nil
		},
	}
}
*/
