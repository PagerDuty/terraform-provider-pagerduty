package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceAutomationActionsAction struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure      = (*resourceAutomationActionsAction)(nil)
	_ resource.ResourceWithImportState    = (*resourceAutomationActionsAction)(nil)
	_ resource.ResourceWithValidateConfig = (*resourceAutomationActionsAction)(nil)
)

func (r *resourceAutomationActionsAction) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_automation_actions_action"
}

func (r *resourceAutomationActionsAction) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":        schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{Optional: true},
			"action_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("script", "process_automation"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runner_id": schema.StringAttribute{Optional: true},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"action_classification": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("diagnostic", "remediation"),
				},
			},
			"runner_type":   schema.StringAttribute{Optional: true, Computed: true},
			"creation_time": schema.StringAttribute{Computed: true},
			"modify_time":   schema.StringAttribute{Computed: true},
			"action_data_reference": schema.ListAttribute{
				Required:    true,
				ElementType: actionDataReferenceObjectType,
			},
		},
	}
}

func (r *resourceAutomationActionsAction) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var actionType types.String

	d := req.Config.GetAttribute(ctx, path.Root("action_type"), &actionType)
	if resp.Diagnostics.Append(d...); d.HasError() {
		return
	}

	if actionType.ValueString() == "script" {
		var script types.String
		scriptPath := path.Root("action_data_reference").AtListIndex(0).AtName("script")
		d := req.Config.GetAttribute(ctx, scriptPath, &script)
		if resp.Diagnostics.Append(d...); d.HasError() {
			return
		}
		if script.IsNull() {
			resp.Diagnostics.AddAttributeError(scriptPath, "action_data_reference script is required", "")
			return
		}
	}

	if actionType.ValueString() == "process_automation" {
		var script types.String
		scriptPath := path.Root("action_data_reference").AtListIndex(0).AtName("process_automation_job_id")
		d := req.Config.GetAttribute(ctx, scriptPath, &script)
		if resp.Diagnostics.Append(d...); d.HasError() {
			return
		}
		if script.IsNull() {
			resp.Diagnostics.AddAttributeError(scriptPath, "action_data_reference process_automation_job_id is required", "")
			return
		}
	}
}

var actionDataReferenceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"process_automation_job_id":        types.StringType,
		"process_automation_job_arguments": types.StringType,
		"process_automation_node_filter":   types.StringType,
		"script":                           types.StringType,
		"invocation_command":               types.StringType,
	},
}

func (r *resourceAutomationActionsAction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceAutomationActionsActionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyAutomationActionsAction(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Creating PagerDuty automation action %s", plan.Name)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		o := pagerduty.CreateAutomationActionOptions{Action: *plan}
		response, err := r.client.CreateAutomationActionWithContext(ctx, o)
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
			fmt.Sprintf("Error creating PagerDuty automation action %s", plan.Name),
			err.Error(),
		)
		return
	}

	model, err = requestGetAutomationActionsAction(ctx, r.client, plan.ID, true, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action %s", plan.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAutomationActionsAction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty automation action %s", id)

	state, err := requestGetAutomationActionsAction(ctx, r.client, id.ValueString(), false, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAutomationActionsAction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceAutomationActionsActionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyAutomationActionsAction(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Updating PagerDuty automation action %s", plan.ID)

	o := pagerduty.UpdateAutomationActionOptions{Action: *plan}
	_, err := r.client.UpdateAutomationActionWithContext(ctx, o)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty automation action %s", plan.ID),
			err.Error(),
		)
		return
	}

	state, err := requestGetAutomationActionsAction(ctx, r.client, plan.ID, false, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action %s", plan.ID),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAutomationActionsAction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty automation action %s", id)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		err := r.client.DeleteAutomationActionWithContext(ctx, id.ValueString())
		if err != nil && !util.IsNotFoundError(err) {
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty automation action %s", id),
			err.Error(),
		)
	}

	resp.State.RemoveResource(ctx)
}

func (r *resourceAutomationActionsAction) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAutomationActionsAction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceAutomationActionsActionModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ActionType           types.String `tfsdk:"action_type"`
	RunnerID             types.String `tfsdk:"runner_id"`
	Type                 types.String `tfsdk:"type"`
	ActionClassification types.String `tfsdk:"action_classification"`
	RunnerType           types.String `tfsdk:"runner_type"`
	CreationTime         types.String `tfsdk:"creation_time"`
	ModifyTime           types.String `tfsdk:"modify_time"`
	ActionDataReference  types.List   `tfsdk:"action_data_reference"`
}

func requestGetAutomationActionsAction(ctx context.Context, client *pagerduty.Client, id string, retryNotFound bool, diags *diag.Diagnostics) (resourceAutomationActionsActionModel, error) {
	var model resourceAutomationActionsActionModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		action, err := client.GetAutomationActionWithContext(ctx, id)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenAutomationActionsAction(action, diags)
		return nil
	})

	return model, err
}

func buildPagerdutyAutomationActionsAction(ctx context.Context, model *resourceAutomationActionsActionModel, diags *diag.Diagnostics) *pagerduty.AutomationAction {
	action := pagerduty.AutomationAction{
		APIObject:  pagerduty.APIObject{ID: model.ID.ValueString()},
		Name:       model.Name.ValueString(),
		ActionType: model.ActionType.ValueString(),
	}

	if !model.ActionDataReference.IsNull() || !model.ActionDataReference.IsUnknown() {
		action.ActionDataReference = buildActionDataReference(ctx, model.ActionDataReference, diags)
	}

	if !model.Description.IsNull() || !model.Description.IsUnknown() {
		action.Description = model.Description.ValueString()
	} else {
		diags.AddAttributeError(
			path.Root("description"),
			"action description must be specified when creating an action",
			"",
		)
	}

	if !model.RunnerID.IsNull() || !model.RunnerID.IsUnknown() {
		action.Runner = model.RunnerID.ValueString()
	}

	if !model.Type.IsNull() || !model.Type.IsUnknown() {
		action.Type = model.Type.ValueString()
	}

	if !model.ActionClassification.IsNull() || !model.ActionClassification.IsUnknown() {
		action.ActionClassification = model.ActionClassification.ValueStringPointer()
	}

	if !model.RunnerType.IsNull() || !model.RunnerType.IsUnknown() {
		action.RunnerType = model.RunnerType.ValueString()
	}

	if !model.CreationTime.IsNull() || !model.CreationTime.IsUnknown() {
		action.CreationTime = model.CreationTime.ValueString()
	}

	if !model.ModifyTime.IsNull() || !model.ModifyTime.IsUnknown() {
		action.ModifyTime = model.ModifyTime.ValueString()
	}

	return &action
}

func buildActionDataReference(ctx context.Context, list types.List, diags *diag.Diagnostics) *pagerduty.ActionDataReference {
	var target []struct {
		ProcessAutomationJobID        types.String `tfsdk:"process_automation_job_id"`
		ProcessAutomationJobArguments types.String `tfsdk:"process_automation_job_arguments"`
		ProcessAutomationNodeFilter   types.String `tfsdk:"process_automation_node_filter"`
		Script                        types.String `tfsdk:"script"`
		InvocationCommand             types.String `tfsdk:"invocation_command"`
	}

	d := list.ElementsAs(ctx, &target, false)
	if diags.Append(d...); d.HasError() {
		return nil
	}
	obj := target[0]

	var script, invocationCommand *string
	s1 := obj.Script.ValueString()
	script = &s1
	s2 := obj.InvocationCommand.ValueString()
	invocationCommand = &s2

	var jobID, jobArguments, nodeFilter *string
	s3 := obj.ProcessAutomationJobID.ValueString()
	jobID = &s3
	s4 := obj.ProcessAutomationJobArguments.ValueString()
	jobArguments = &s4
	s5 := obj.ProcessAutomationNodeFilter.ValueString()
	nodeFilter = &s5

	out := &pagerduty.ActionDataReference{
		Script:                        script,
		InvocationCommand:             invocationCommand,
		ProcessAutomationJobID:        jobID,
		ProcessAutomationJobArguments: jobArguments,
		ProcessAutomationNodeFilter:   nodeFilter,
	}
	return out
}

func flattenAutomationActionsAction(response *pagerduty.AutomationAction, diags *diag.Diagnostics) resourceAutomationActionsActionModel {
	model := resourceAutomationActionsActionModel{
		ID:                  types.StringValue(response.ID),
		Name:                types.StringValue(response.Name),
		Type:                types.StringValue(response.Type),
		ActionType:          types.StringValue(response.ActionType),
		CreationTime:        types.StringValue(response.CreationTime),
		ModifyTime:          types.StringValue(response.ModifyTime),
		Description:         types.StringValue(response.Description),
		ActionDataReference: flattenActionDataReference(response.ActionDataReference, diags),
	}

	if response.Runner != "" {
		model.RunnerID = types.StringValue(response.Runner)
	}

	if response.RunnerType != "" {
		model.RunnerType = types.StringValue(response.RunnerType)
	}

	if response.ActionClassification != nil {
		model.ActionClassification = types.StringValue(*response.ActionClassification)
	}

	return model
}

func flattenActionDataReference(ref *pagerduty.ActionDataReference, diags *diag.Diagnostics) types.List {
	values := map[string]attr.Value{
		"script":                           types.StringNull(),
		"invocation_command":               types.StringNull(),
		"process_automation_job_id":        types.StringNull(),
		"process_automation_job_arguments": types.StringNull(),
		"process_automation_node_filter":   types.StringNull(),
	}

	for k, v := range map[string]*string{
		"script":                           ref.Script,
		"invocation_command":               ref.InvocationCommand,
		"process_automation_job_id":        ref.ProcessAutomationJobID,
		"process_automation_job_arguments": ref.ProcessAutomationJobArguments,
		"process_automation_node_filter":   ref.ProcessAutomationNodeFilter,
	} {
		if v != nil {
			values[k] = types.StringValue(*v)
		}
	}

	obj, d := types.ObjectValue(actionDataReferenceObjectType.AttrTypes, values)
	if diags.Append(d...); d.HasError() {
		return types.ListNull(actionDataReferenceObjectType)
	}

	list, d := types.ListValue(actionDataReferenceObjectType, []attr.Value{obj})
	if diags.Append(d...); d.HasError() {
		return types.ListNull(actionDataReferenceObjectType)
	}

	return list
}
