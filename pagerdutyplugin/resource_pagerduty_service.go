package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/enumtypes"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/rangetypes"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/tztypes"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceService struct {
	client *pagerduty.Client
}

var (
	_ resource.ResourceWithConfigure        = (*resourceService)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceService)(nil)
	_ resource.ResourceWithImportState      = (*resourceService)(nil)
)

func (r *resourceService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceService) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		validate.Require(
			path.Root("incident_urgency_rule").AtListIndex(0).AtName("type"),
		),
		validate.RequireAIfBEqual(
			path.Root("support_hours"),
			path.Root("incident_urgency_rule").AtListIndex(0).AtName("type"),
			types.StringValue("use_support_hours"),
		),
		validate.RequireList(path.Root("alert_grouping_parameters").AtListIndex(0).AtName("config")),
		validate.RequireList(path.Root("incident_urgency_rule").AtListIndex(0).AtName("during_support_hours")),
		validate.RequireList(path.Root("incident_urgency_rule").AtListIndex(0).AtName("outside_support_hours")),
		validate.RequireList(path.Root("support_hours").AtListIndex(0).AtName("days_of_week")), // TODO at most 7
	}
}

func (r *resourceService) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_service"
}

func (r *resourceService) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validate.IsAllowedString(util.NoNonPrintableChars),
				},
			},

			"acknowledgement_timeout": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("1800"),
			},

			"alert_creation": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("create_alerts_and_incidents"),
				Validators: []validator.String{
					stringvalidator.OneOf("create_alerts_and_incidents", "create_incidents"),
				},
			},

			"alert_grouping": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("time", "intelligent", "content_based"),
					stringvalidator.ConflictsWith(path.MatchRoot("alert_grouping_parameters")),
				},
				DeprecationMessage: "Use `alert_grouping_parameters.type`",
			},

			"alert_grouping_timeout": schema.StringAttribute{
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "Use `alert_grouping_parameters.config.timeout`",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("alert_grouping_parameters")),
				},
			},

			"auto_resolve_timeout": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("14400"),
			},

			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Managed by Terraform"),
			},

			"id":                      schema.StringAttribute{Computed: true},
			"created_at":              schema.StringAttribute{Computed: true},
			"escalation_policy":       schema.StringAttribute{Required: true},
			"html_url":                schema.StringAttribute{Computed: true},
			"last_incident_timestamp": schema.StringAttribute{Computed: true},
			"response_play":           schema.StringAttribute{Computed: true, Optional: true},
			"status":                  schema.StringAttribute{Computed: true},
			"type":                    schema.StringAttribute{Computed: true},

			"alert_grouping_parameters": schema.ListAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 1),
					listvalidator.ConflictsWith(path.MatchRoot("alert_grouping")),
					listvalidator.ConflictsWith(path.MatchRoot("alert_grouping_timeout")),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type": alertGroupingParametersTypeType,
						"config": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"timeout":     types.Int64Type,
									"fields":      types.ListType{ElemType: types.StringType},
									"aggregate":   alertGroupingParametersConfigAggregateType,
									"time_window": alertGroupingParametersConfigTimeWindowType,
								},
							},
						},
					},
				},
			},

			"auto_pause_notifications_parameters": schema.ListAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"enabled": types.BoolType,
						"timeout": autoPauseNotificationsParametersTimeoutType,
					},
				},
			},

			"incident_urgency_rule": schema.ListAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":    types.StringType,
						"urgency": types.StringType,
						"during_support_hours": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type":    types.StringType, // require
									"urgency": types.StringType,
								},
							},
						},
						"outside_support_hours": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type":    types.StringType, // require
									"urgency": types.StringType,
								},
							},
						},
					},
				},
			},

			"scheduled_actions": schema.ListAttribute{
				Optional: true,
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":       types.StringType,
						"to_urgency": types.StringType,
						"at": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type": types.StringType,
									"name": types.StringType,
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 1),
				},
			},

			"support_hours": schema.ListAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":       types.StringType,
						"start_time": types.StringType,
						"end_time":   types.StringType,
						"time_zone":  tztypes.StringType{},
						"days_of_week": types.ListType{
							ElemType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *resourceService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceServiceModel
	if d := req.Plan.Get(ctx, &model); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	serviceBody := buildService(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Creating PagerDuty service %s", serviceBody.Name)

	service, err := r.client.CreateServiceWithContext(ctx, serviceBody)
	if err != nil {
		resp.Diagnostics.AddError("Error calling CreateServiceWithContext", err.Error())
		return
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		serviceResponse, err := r.client.GetServiceWithContext(ctx, service.ID, &pagerduty.GetServiceOptions{
			Includes: []string{"auto_pause_notifications_parameters"},
		})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenService(ctx, serviceResponse, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("%#v", resp.Diagnostics))
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error calling GetServiceWithContext", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	if d := req.State.GetAttribute(ctx, path.Root("id"), &id); d.HasError() {
		resp.Diagnostics.Append(d...)
	}
	log.Printf("[INFO] Reading PagerDuty service %s", id)

	if id.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	var model resourceServiceModel
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		serviceResponse, err := r.client.GetServiceWithContext(ctx, id.ValueString(), &pagerduty.GetServiceOptions{
			Includes: []string{"auto_pause_notifications_parameters"},
		})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				resp.State.RemoveResource(ctx)
				return nil
			}
			return retry.RetryableError(err)
		}
		model = flattenService(ctx, serviceResponse, &resp.Diagnostics)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error calling GetServiceWithContext", err.Error())
		return
	}
	resp.State.Set(ctx, &model)
}

func (r *resourceService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	if d := req.State.GetAttribute(ctx, path.Root("id"), &id); d.HasError() {
		resp.Diagnostics.Append(d...)
	}
	log.Printf("[INFO] Deleting PagerDuty service %s", id)

	if id.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := r.client.DeleteServiceWithContext(ctx, id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error calling DeleteServiceWithContext", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

type resourceServiceModel struct {
	ID                               types.String `tfsdk:"id"`
	AcknowledgementTimeout           types.String `tfsdk:"acknowledgement_timeout"`
	AlertCreation                    types.String `tfsdk:"alert_creation"`
	AlertGrouping                    types.String `tfsdk:"alert_grouping"`
	AlertGroupingTimeout             types.String `tfsdk:"alert_grouping_timeout"`
	AutoResolveTimeout               types.String `tfsdk:"auto_resolve_timeout"`
	CreatedAt                        types.String `tfsdk:"created_at"`
	Description                      types.String `tfsdk:"description"`
	EscalationPolicy                 types.String `tfsdk:"escalation_policy"`
	HtmlUrl                          types.String `tfsdk:"html_url"`
	LastIncidentTimestamp            types.String `tfsdk:"last_incident_timestamp"`
	Name                             types.String `tfsdk:"name"`
	ResponsePlay                     types.String `tfsdk:"response_play"`
	Status                           types.String `tfsdk:"status"`
	Type                             types.String `tfsdk:"type"`
	AlertGroupingParameters          types.List   `tfsdk:"alert_grouping_parameters"`
	AutoPauseNotificationsParameters types.List   `tfsdk:"auto_pause_notifications_parameters"`
	IncidentUrgencyRule              types.List   `tfsdk:"incident_urgency_rule"`
	ScheduledActions                 types.List   `tfsdk:"scheduled_actions"`
	SupportHours                     types.List   `tfsdk:"support_hours"`
}

func buildService(ctx context.Context, model *resourceServiceModel, diags *diag.Diagnostics) pagerduty.Service {
	service := pagerduty.Service{
		Name:          model.Name.ValueString(),
		Description:   model.Description.ValueString(),
		AlertCreation: model.AlertCreation.ValueString(),
		AlertGrouping: model.AlertGrouping.ValueString(),
	}

	u := util.StringToUintPointer(path.Root("auto_resolve_timeout"), model.AutoResolveTimeout, diags)
	service.AutoResolveTimeout = u

	u = util.StringToUintPointer(path.Root("acknowledgement_timeout"), model.AcknowledgementTimeout, diags)
	service.AcknowledgementTimeout = u

	u = util.StringToUintPointer(path.Root("alert_grouping_timeout"), model.AlertGroupingTimeout, diags)
	service.AlertGroupingTimeout = u

	service.EscalationPolicy.ID = model.EscalationPolicy.ValueString()
	service.EscalationPolicy.Type = "escalation_policy_reference"

	service.AlertGroupingParameters = buildAlertGroupingParameters(ctx, model.AlertGroupingParameters, diags)
	service.AutoPauseNotificationsParameters = buildAutoPauseNotificationsParameters(ctx, model.AutoPauseNotificationsParameters, diags)
	service.IncidentUrgencyRule = buildIncidentUrgencyRule(ctx, model.IncidentUrgencyRule, diags)
	service.ScheduledActions = buildScheduledActions(ctx, model.ScheduledActions, diags)
	service.SupportHours = buildSupportHours(ctx, model.SupportHours, diags)

	if !model.ResponsePlay.IsNull() && !model.ResponsePlay.IsUnknown() {
		service.ResponsePlay = &pagerduty.APIObject{
			ID:   model.ResponsePlay.ValueString(),
			Type: "response_play_reference",
		}
	}

	return service
}

func buildAlertGroupingParameters(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.AlertGroupingParameters {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var target []struct {
		Type   types.String `tfsdk:"type"`
		Config types.List   `tfsdk:"config"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return nil
	}
	obj := target[0]
	return &pagerduty.AlertGroupingParameters{
		Type:   obj.Type.ValueString(),
		Config: buildAlertGroupingConfig(ctx, obj.Config, diags),
	}
}

func buildAlertGroupingConfig(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.AlertGroupParamsConfig {
	var target []struct {
		Timeout    types.Int64  `tfsdk:"timeout"`
		Aggregate  types.String `tfsdk:"aggregate"`
		Fields     types.List   `tfsdk:"fields"`
		TimeWindow types.Int64  `tfsdk:"time_window"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return nil
	}
	obj := target[0]

	ut := uint(obj.Timeout.ValueInt64())

	var fields []string
	if d := obj.Fields.ElementsAs(ctx, &fields, false); d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &pagerduty.AlertGroupParamsConfig{
		Timeout:   &ut,
		Aggregate: obj.Aggregate.ValueString(),
		Fields:    fields,
	}
}

func buildAutoPauseNotificationsParameters(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.AutoPauseNotificationsParameters {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var target []struct {
		Timeout types.Int64 `tfsdk:"timeout"`
		Enabled types.Bool  `tfsdk:"enabled"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return nil
	}
	obj := target[0]

	return &pagerduty.AutoPauseNotificationsParameters{
		Enabled: obj.Enabled.ValueBool(),
		Timeout: uint(obj.Timeout.ValueInt64()),
	}
}

func buildIncidentUrgencyRule(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.IncidentUrgencyRule {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var target []struct {
		Type                types.String `tfsdk:"type"`
		Urgency             types.String `tfsdk:"urgency"`
		DuringSupportHours  types.List   `tfsdk:"during_support_hours"`
		OutsideSupportHours types.List   `tfsdk:"outside_support_hours"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return nil
	}
	obj := target[0]
	incidentUrgencyRule := &pagerduty.IncidentUrgencyRule{
		Type:    obj.Type.ValueString(),
		Urgency: obj.Urgency.ValueString(),
	}
	incidentUrgencyRule.DuringSupportHours = buildIncidentUrgencyType(ctx, obj.DuringSupportHours, diags)
	incidentUrgencyRule.OutsideSupportHours = buildIncidentUrgencyType(ctx, obj.OutsideSupportHours, diags)
	return incidentUrgencyRule
}

func buildIncidentUrgencyType(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.IncidentUrgencyType {
	var target []struct {
		Type    types.String `tfsdk:"type"`
		Urgency types.String `tfsdk:"urgency"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
	}
	if len(target) < 1 {
		return nil
	}
	obj := target[0]
	return &pagerduty.IncidentUrgencyType{
		Type:    obj.Type.ValueString(),
		Urgency: obj.Urgency.ValueString(),
	}
}

func buildScheduledActions(ctx context.Context, list types.List, diags *diag.Diagnostics) []pagerduty.ScheduledAction {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var target []struct {
		Type      types.String `tfsdk:"type"`
		ToUrgency types.String `tfsdk:"to_urgency"`
		At        types.List   `tfsdk:"at"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
	}
	scheduledActions := []pagerduty.ScheduledAction{}
	for _, src := range target {
		dst := pagerduty.ScheduledAction{
			Type:      src.Type.ValueString(),
			ToUrgency: src.ToUrgency.ValueString(),
			At:        buildScheduledActionAt(ctx, src.At, diags),
		}
		scheduledActions = append(scheduledActions, dst)
	}
	return scheduledActions
}

func buildScheduledActionAt(ctx context.Context, list types.List, diags *diag.Diagnostics) pagerduty.InlineModel {
	var target []struct {
		Type types.String `tfsdk:"type"`
		Name types.String `tfsdk:"name"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return pagerduty.InlineModel{}
	}
	obj := target[0]
	return pagerduty.InlineModel{
		Type: obj.Type.ValueString(),
		Name: obj.Name.ValueString(),
	}
}

func buildSupportHours(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.SupportHours {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var target []struct {
		Type       types.String `tfsdk:"type"`
		Timezone   types.String `tfsdk:"time_zone"`
		StartTime  types.String `tfsdk:"start_time"`
		EndTime    types.String `tfsdk:"end_time"`
		DaysOfWeek types.List   `tfsdk:"days_of_week"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return nil
	}
	obj := target[0]
	supportHours := &pagerduty.SupportHours{
		Type:      obj.Type.ValueString(),
		Timezone:  obj.Timezone.ValueString(),
		StartTime: obj.StartTime.ValueString(),
		EndTime:   obj.EndTime.ValueString(),
	}

	if !obj.DaysOfWeek.IsNull() {
		daysOfWeekStr := []string{}
		if d := obj.DaysOfWeek.ElementsAs(ctx, &daysOfWeekStr, false); d.HasError() {
			diags.Append(d...)
			return nil
		}
		daysOfWeek := make([]uint, 0, len(daysOfWeekStr))
		for _, s := range daysOfWeekStr {
			v, err := strconv.Atoi(s)
			if err != nil {
				continue
			}
			daysOfWeek = append(daysOfWeek, uint(v))
		}
		supportHours.DaysOfWeek = daysOfWeek
	}
	return supportHours
}

var (
	alertGroupingParametersTypeType = enumtypes.StringType{
		OneOf: []string{"time", "intelligent", "content_based"}}
	alertGroupingParametersConfigAggregateType = enumtypes.StringType{
		OneOf: []string{"all", "any"}}
	alertGroupingParametersConfigTimeWindowType = rangetypes.Int64Type{
		Start: 300, End: 3600}
	autoPauseNotificationsParametersTimeoutType = enumtypes.Int64Type{
		OneOf: []int64{120, 180, 300, 600, 900}}
)

func flattenService(ctx context.Context, service *pagerduty.Service, diags *diag.Diagnostics) resourceServiceModel {
	model := resourceServiceModel{
		ID:                    types.StringValue(service.ID),
		AlertCreation:         types.StringValue(service.AlertCreation),
		CreatedAt:             types.StringValue(service.CreateAt),
		Description:           types.StringValue(service.Description),
		EscalationPolicy:      types.StringValue(service.EscalationPolicy.ID),
		HtmlUrl:               types.StringValue(service.HTMLURL),
		LastIncidentTimestamp: types.StringValue(service.LastIncidentTimestamp),
		Name:                  types.StringValue(service.Name),
		Status:                types.StringValue(service.Status),
		Type:                  types.StringValue(service.Type),
	}

	if service.AcknowledgementTimeout != nil {
		s := strconv.Itoa(int(*service.AcknowledgementTimeout))
		model.AcknowledgementTimeout = types.StringValue(s)
	}

	if service.AutoResolveTimeout != nil {
		s := strconv.Itoa(int(*service.AutoResolveTimeout))
		model.AutoResolveTimeout = types.StringValue(s)
	}

	if service.AlertGrouping != "" {
		model.AlertGrouping = types.StringValue(service.AlertGrouping)
	}

	if service.AlertGroupingTimeout != nil {
		s := strconv.Itoa(int(*service.AlertGroupingTimeout))
		model.AlertGroupingTimeout = types.StringValue(s)
	}

	model.AlertGroupingParameters = flattenAlertGroupingParameters(ctx, service.AlertGroupingParameters, diags)
	model.AutoPauseNotificationsParameters = flattenAutoPauseNotificationsParameters(service.AutoPauseNotificationsParameters, diags)
	model.IncidentUrgencyRule = flattenIncidentUrgencyRule(service.IncidentUrgencyRule, diags)

	if service.ResponsePlay != nil {
		model.ResponsePlay = types.StringValue(service.ResponsePlay.ID)
	}

	model.ScheduledActions = flattenScheduledActions(service.ScheduledActions, diags)
	model.SupportHours = flattenSupportHours(service.SupportHours, diags)

	return model
}

func flattenAlertGroupingParameters(ctx context.Context, params *pagerduty.AlertGroupingParameters, diags *diag.Diagnostics) types.List {
	alertGroupParamsConfigObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"aggregate":   alertGroupingParametersConfigAggregateType,
			"fields":      types.ListType{ElemType: types.StringType},
			"timeout":     types.Int64Type,
			"time_window": alertGroupingParametersConfigTimeWindowType,
		},
	}
	alertGroupingParametersObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":   alertGroupingParametersTypeType,
			"config": types.ListType{ElemType: alertGroupParamsConfigObjectType},
		},
	}

	nullList := types.ListNull(alertGroupingParametersObjectType)
	if params == nil {
		return nullList
	}

	configList := types.ListNull(alertGroupParamsConfigObjectType)
	if params.Config != nil {
		fieldsList, d := types.ListValueFrom(ctx, types.StringType, params.Config.Fields)
		if d.HasError() {
			diags.Append(d...)
			return nullList
		}

		var timeout types.Int64
		if params.Config.Timeout != nil {
			timeout = types.Int64Value(int64(*params.Config.Timeout))
		}

		aggregate := enumtypes.NewStringNull(alertGroupingParametersConfigAggregateType)
		if params.Config.Aggregate != "" {
			aggregate = enumtypes.NewStringValue(params.Config.Aggregate, alertGroupingParametersConfigAggregateType)
		}

		configObj, d := types.ObjectValue(alertGroupParamsConfigObjectType.AttrTypes, map[string]attr.Value{
			"aggregate":   aggregate,
			"fields":      fieldsList,
			"timeout":     timeout,
			"time_window": types.Int64Null(), // TODO
		})
		if d.HasError() {
			diags.Append(d...)
			return nullList
		}
		configList, d = types.ListValue(alertGroupParamsConfigObjectType, []attr.Value{configObj})
		if d.HasError() {
			diags.Append(d...)
			return nullList
		}
	}

	obj, d := types.ObjectValue(alertGroupingParametersObjectType.AttrTypes, map[string]attr.Value{
		"type":   enumtypes.NewStringValue(params.Type, alertGroupingParametersTypeType),
		"config": configList,
	})
	diags.Append(d...)
	if d.HasError() {
		return nullList
	}

	list, d := types.ListValue(alertGroupingParametersObjectType, []attr.Value{obj})
	diags.Append(d...)
	if d.HasError() {
		return nullList
	}

	return list
}

func flattenAutoPauseNotificationsParameters(params *pagerduty.AutoPauseNotificationsParameters, diags *diag.Diagnostics) types.List {
	autoPauseNotificationsParametersObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
			"timeout": autoPauseNotificationsParametersTimeoutType,
		},
	}

	nullList := types.ListNull(autoPauseNotificationsParametersObjectType)
	if params == nil {
		return nullList
	}

	timeout := enumtypes.NewInt64Null(autoPauseNotificationsParametersTimeoutType)
	if params.Enabled {
		timeout = enumtypes.NewInt64Value(
			int64(params.Timeout),
			autoPauseNotificationsParametersTimeoutType,
		)
	}

	obj, d := types.ObjectValue(autoPauseNotificationsParametersObjectType.AttrTypes, map[string]attr.Value{
		"enabled": types.BoolValue(params.Enabled),
		"timeout": timeout,
	})
	if d.HasError() {
		diags.Append(d...)
		return nullList
	}

	list, d := types.ListValue(autoPauseNotificationsParametersObjectType, []attr.Value{obj})
	if d.HasError() {
		diags.Append(d...)
		return nullList
	}

	return list
}

var incidentUrgencyTypeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":    types.StringType,
		"urgency": types.StringType,
	},
}

func flattenIncidentUrgencyRule(rule *pagerduty.IncidentUrgencyRule, diags *diag.Diagnostics) types.List {
	incidentUrgencyRuleObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":                  types.StringType,
			"urgency":               types.StringType,
			"during_support_hours":  types.ListType{ElemType: incidentUrgencyTypeObjectType},
			"outside_support_hours": types.ListType{ElemType: incidentUrgencyTypeObjectType},
		},
	}
	nullList := types.ListNull(incidentUrgencyTypeObjectType)
	if rule == nil {
		return nullList
	}

	objValues := map[string]attr.Value{
		"type":                  types.StringValue(rule.Type),
		"urgency":               types.StringNull(),
		"during_support_hours":  types.ListNull(incidentUrgencyTypeObjectType),
		"outside_support_hours": types.ListNull(incidentUrgencyTypeObjectType),
	}
	if rule.Urgency != "" {
		objValues["urgency"] = types.StringValue(rule.Urgency)
	}
	if rule.DuringSupportHours != nil {
		objValues["during_support_hours"] = flattenIncidentUrgencyType(rule.DuringSupportHours, diags)
	}
	if rule.OutsideSupportHours != nil {
		objValues["outside_support_hours"] = flattenIncidentUrgencyType(rule.OutsideSupportHours, diags)
	}
	if diags.HasError() {
		return nullList
	}

	obj, d := types.ObjectValue(incidentUrgencyRuleObjectType.AttrTypes, objValues)
	if d.HasError() {
		diags.Append(d...)
		return nullList
	}

	list, d := types.ListValue(incidentUrgencyRuleObjectType, []attr.Value{obj})
	diags.Append(d...)
	return list
}

func flattenIncidentUrgencyType(urgency *pagerduty.IncidentUrgencyType, diags *diag.Diagnostics) types.List {
	obj, d := types.ObjectValue(incidentUrgencyTypeObjectType.AttrTypes, map[string]attr.Value{
		"type":    types.StringValue(urgency.Type),
		"urgency": types.StringValue(urgency.Urgency),
	})
	diags.Append(d...)
	if d.HasError() {
		return types.List{}
	}
	list, d := types.ListValue(incidentUrgencyTypeObjectType, []attr.Value{obj})
	diags.Append(d...)
	return list
}

var scheduledActionAtObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"name": types.StringType,
	},
}

func flattenScheduledActions(actions []pagerduty.ScheduledAction, diags *diag.Diagnostics) types.List {
	scheduledActionObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":       types.StringType,
			"to_urgency": types.StringType,
			"at":         types.ListType{ElemType: scheduledActionAtObjectType},
		},
	}
	nullList := types.ListNull(scheduledActionObjectType)
	if len(actions) == 0 {
		return nullList
	}

	elements := []attr.Value{}
	for _, action := range actions {
		obj, d := types.ObjectValue(scheduledActionObjectType.AttrTypes, map[string]attr.Value{
			"type":       types.StringValue(action.Type),
			"to_urgency": types.StringValue(action.ToUrgency),
			"at":         flattenScheduledActionAt(action.At, diags),
		})
		diags.Append(d...)
		if diags.HasError() {
			return nullList
		}
		elements = append(elements, obj)
	}

	list, d := types.ListValue(scheduledActionObjectType, elements)
	diags.Append(d...)
	return list
}

func flattenScheduledActionAt(at pagerduty.InlineModel, diags *diag.Diagnostics) types.List {
	obj, d := types.ObjectValue(scheduledActionAtObjectType.AttrTypes, map[string]attr.Value{
		"type": types.StringValue(at.Type),
		"name": types.StringValue(at.Name),
	})
	if d.HasError() {
		diags.Append(d...)
		return types.List{}
	}
	list, d := types.ListValue(scheduledActionAtObjectType, []attr.Value{obj})
	diags.Append(d...)
	return list
}

func flattenSupportHours(hours *pagerduty.SupportHours, diags *diag.Diagnostics) types.List {
	supportHoursObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":         types.StringType,
			"start_time":   types.StringType,
			"end_time":     types.StringType,
			"time_zone":    tztypes.StringType{},
			"days_of_week": types.ListType{ElemType: types.StringType},
		},
	}
	nullList := types.ListNull(supportHoursObjectType)
	if hours == nil {
		return nullList
	}

	daysOfWeek := []attr.Value{}
	for _, dow := range hours.DaysOfWeek {
		v := strconv.FormatInt(int64(dow), 10)
		daysOfWeek = append(daysOfWeek, types.StringValue(v))
	}

	dowList, d := types.ListValue(types.StringType, daysOfWeek)
	diags.Append(d...)

	obj, d := types.ObjectValue(supportHoursObjectType.AttrTypes, map[string]attr.Value{
		"type":         types.StringValue(hours.Type),
		"start_time":   types.StringValue(hours.StartTime),
		"end_time":     types.StringValue(hours.EndTime),
		"time_zone":    tztypes.NewStringValue(hours.Timezone),
		"days_of_week": dowList,
	})
	if d.HasError() {
		diags.Append(d...)
		return nullList
	}

	list, d := types.ListValue(supportHoursObjectType, []attr.Value{obj})
	diags.Append(d...)
	return list
}
