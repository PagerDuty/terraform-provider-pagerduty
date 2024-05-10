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
	_ resource.ResourceWithValidateConfig   = (*resourceService)(nil)
)

func (r *resourceService) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_service"
}

func (r *resourceService) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

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
				Computed:      true,
				Optional:      true,
				Default:       stringdefault.StaticString("1800"),
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"alert_creation": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("create_alerts_and_incidents"),
				Validators: []validator.String{
					stringvalidator.OneOf("create_alerts_and_incidents", "create_incidents"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"alert_grouping": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("time", "intelligent", "content_based"),
					stringvalidator.ConflictsWith(path.MatchRoot("alert_grouping_parameters")),
				},
				DeprecationMessage: "Use `alert_grouping_parameters.type`",
				PlanModifiers:      []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"alert_grouping_timeout": schema.StringAttribute{
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "Use `alert_grouping_parameters.config.timeout`",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("alert_grouping_parameters")),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"auto_resolve_timeout": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				Default:       stringdefault.StaticString("14400"),
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Default:       stringdefault.StaticString("Managed by Terraform"),
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"created_at": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"escalation_policy": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"html_url": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"last_incident_timestamp": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"response_play": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"status": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

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
									"time_window": types.Int64Type,
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
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
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
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
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
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
					listplanmodifier.UseStateForUnknown(),
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
					listplanmodifier.UseStateForUnknown(),
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
		validate.ForbidAIfBEqualWithMessage(
			path.Root("incident_urgency_rule").AtListIndex(0).AtName("urgency"),
			path.Root("incident_urgency_rule").AtListIndex(0).AtName("type"),
			types.StringValue("use_support_hours"),
			"general urgency cannot be set for a use_support_hours incident urgency rule type",
		),
		validate.RequireList(path.Root("alert_grouping_parameters").AtListIndex(0).AtName("config")),
		validate.RequireList(path.Root("incident_urgency_rule").AtListIndex(0).AtName("during_support_hours")),
		validate.RequireList(path.Root("incident_urgency_rule").AtListIndex(0).AtName("outside_support_hours")),
		validate.RequireList(path.Root("support_hours").AtListIndex(0).AtName("days_of_week")), // TODO at most 7
	}
}

func (r *resourceService) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	configPath := path.Root("alert_grouping_parameters").AtListIndex(0).AtName("config").AtListIndex(0)

	// Validate time window
	var timeWindow types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, configPath.AtName("time_window"), &timeWindow)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !timeWindow.IsNull() && !timeWindow.IsUnknown() {
		if tw := timeWindow.ValueInt64(); tw < 300 || tw > 3600 {
			resp.Diagnostics.AddAttributeError(
				configPath.AtName("time_window"),
				"Alert grouping time window value must be between 300 and 3600",
				fmt.Sprintf("Current setting is %d", tw),
			)
		}
	}

	// Validate Alert Grouping Parameters
	var aggregate types.String
	var fields types.List
	var timeout types.Int64
	var pType types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("alert_grouping_parameters").AtListIndex(0).AtName("type"), &pType)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, configPath.AtName("aggregate"), &aggregate)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, configPath.AtName("fields"), &fields)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, configPath.AtName("timeout"), &timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if pType.ValueString() == "content_based" && (aggregate.ValueString() == "" || len(fields.Elements()) == 0) {
		resp.Diagnostics.AddError(`When using Alert grouping parameters configuration of type "content_based" is in use, attributes "aggregate" and "fields" are required`, "")
		return
	}

	if !pType.IsNull() && pType.ValueString() != "content_based" && (aggregate.ValueString() != "" || len(fields.Elements()) > 0) {
		resp.Diagnostics.AddError(`Alert grouping parameters configuration attributes "aggregate" and "fields" are only supported by "content_based" type Alert Grouping`, "")
		return
	}

	if !pType.IsNull() && pType.ValueString() != "time" && timeout.ValueInt64() > 0 {
		resp.Diagnostics.AddError(`Alert grouping parameters configuration attribute "timeout" is only supported by "time" type Alert Grouping`, "")
		return
	}

	if !pType.IsNull() && (pType.ValueString() != "intelligent" && pType.ValueString() != "content_based") && timeWindow.ValueInt64() > 300 {
		resp.Diagnostics.AddError(`Alert grouping parameters configuration attribute "time_window" is only supported by "intelligent" and "content-based" type Alert Grouping`, "")
		return
	}
}

func (r *resourceService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config resourceServiceModel
	var model resourceServiceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty service %s", serviceBody.Name),
			err.Error(),
		)
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
		model = flattenService(ctx, serviceResponse, config, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("%#v", resp.Diagnostics))
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty service %s", service.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

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
		model = flattenService(ctx, serviceResponse, state, &resp.Diagnostics)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty service %s", id),
			err.Error(),
		)
		return
	}
	resp.State.Set(ctx, &model)
}

func (r *resourceService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state resourceServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	var model resourceServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildService(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Updating PagerDuty service %s", plan.ID)

	service, err := r.client.UpdateServiceWithContext(ctx, plan)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty service %s", plan.ID),
			err.Error(),
		)
		return
	}
	model = flattenService(ctx, service, state, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
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
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty service %s", id),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *resourceService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

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

	timeout := uint(obj.Timeout.ValueInt64())
	timeWindow := uint(obj.TimeWindow.ValueInt64())

	var fields []string
	if d := obj.Fields.ElementsAs(ctx, &fields, false); d.HasError() {
		diags.Append(d...)
		return nil
	}

	return &pagerduty.AlertGroupParamsConfig{
		Timeout:    &timeout,
		Aggregate:  obj.Aggregate.ValueString(),
		Fields:     fields,
		TimeWindow: &timeWindow,
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
	defaultValue := &pagerduty.IncidentUrgencyRule{Type: "constant", Urgency: "high"}

	if list.IsNull() || list.IsUnknown() {
		return defaultValue
	}

	var target []struct {
		Type                types.String `tfsdk:"type"`
		Urgency             types.String `tfsdk:"urgency"`
		DuringSupportHours  types.List   `tfsdk:"during_support_hours"`
		OutsideSupportHours types.List   `tfsdk:"outside_support_hours"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return defaultValue
	}
	obj := target[0]

	return &pagerduty.IncidentUrgencyRule{
		Type:                obj.Type.ValueString(),
		Urgency:             obj.Urgency.ValueString(),
		DuringSupportHours:  buildIncidentUrgencyType(ctx, obj.DuringSupportHours, diags),
		OutsideSupportHours: buildIncidentUrgencyType(ctx, obj.OutsideSupportHours, diags),
	}
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
	scheduledActions := []pagerduty.ScheduledAction{}
	if list.IsNull() || list.IsUnknown() {
		return scheduledActions
	}
	var target []struct {
		Type      types.String `tfsdk:"type"`
		ToUrgency types.String `tfsdk:"to_urgency"`
		At        types.List   `tfsdk:"at"`
	}
	if d := list.ElementsAs(ctx, &target, false); d.HasError() {
		diags.Append(d...)
		return scheduledActions
	}
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

func flattenService(ctx context.Context, service *pagerduty.Service, state resourceServiceModel, diags *diag.Diagnostics) resourceServiceModel {
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
	} else if state.AcknowledgementTimeout.ValueString() == "null" {
		model.AcknowledgementTimeout = state.AcknowledgementTimeout
	}

	if service.AutoResolveTimeout != nil {
		s := strconv.Itoa(int(*service.AutoResolveTimeout))
		model.AutoResolveTimeout = types.StringValue(s)
	} else if state.AutoResolveTimeout.ValueString() == "null" {
		model.AutoResolveTimeout = state.AutoResolveTimeout
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
	nullList := types.ListNull(alertGroupingParametersObjectType)
	if params == nil {
		return nullList
	}

	configList := types.ListNull(alertGroupingParametersConfigObjectType)
	log.Printf("[CG] config %#v", params.Config)
	if params.Config != nil {
		fieldsList, d := types.ListValueFrom(ctx, types.StringType, params.Config.Fields)
		if d.HasError() {
			diags.Append(d...)
			return nullList
		}

		timeout := types.Int64Null()
		if params.Config.Timeout != nil {
			timeout = types.Int64Value(int64(*params.Config.Timeout))
		}

		timeWindow := types.Int64Null()
		if params.Config.TimeWindow != nil {
			timeWindow = types.Int64Value(int64(*params.Config.TimeWindow))
		}

		aggregate := enumtypes.NewStringNull(alertGroupingParametersConfigAggregateType)
		if params.Config.Aggregate != "" {
			aggregate = enumtypes.NewStringValue(params.Config.Aggregate, alertGroupingParametersConfigAggregateType)
		}

		configObj, d := types.ObjectValue(alertGroupingParametersConfigObjectType.AttrTypes, map[string]attr.Value{
			"aggregate":   aggregate,
			"fields":      fieldsList,
			"timeout":     timeout,
			"time_window": timeWindow,
		})
		if d.HasError() {
			diags.Append(d...)
			return nullList
		}
		configList, d = types.ListValue(alertGroupingParametersConfigObjectType, []attr.Value{configObj})
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

var (
	alertGroupingParametersTypeType = enumtypes.StringType{
		OneOf: []string{"time", "intelligent", "content_based"}}
	alertGroupingParametersConfigAggregateType = enumtypes.StringType{
		OneOf: []string{"all", "any"}}
	autoPauseNotificationsParametersTimeoutType = enumtypes.Int64Type{
		OneOf: []int64{120, 180, 300, 600, 900}}

	alertGroupingParametersConfigObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"aggregate":   alertGroupingParametersConfigAggregateType,
			"fields":      types.ListType{ElemType: types.StringType},
			"timeout":     types.Int64Type,
			"time_window": types.Int64Type,
		},
	}
	alertGroupingParametersObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":   alertGroupingParametersTypeType,
			"config": types.ListType{ElemType: alertGroupingParametersConfigObjectType},
		},
	}

	alertGroupingParametersPath = path.Root("alert_grouping_parameters").AtListIndex(0)

	incidentUrgencyTypeObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":    types.StringType,
			"urgency": types.StringType,
		},
	}

	scheduledActionAtObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"name": types.StringType,
		},
	}
)
