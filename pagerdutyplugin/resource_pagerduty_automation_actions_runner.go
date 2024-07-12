package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceAutomationActionsRunner struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure      = (*resourceAutomationActionsRunner)(nil)
	_ resource.ResourceWithImportState    = (*resourceAutomationActionsRunner)(nil)
	_ resource.ResourceWithValidateConfig = (*resourceAutomationActionsRunner)(nil)
)

func (r *resourceAutomationActionsRunner) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_automation_actions_runner"
}

func (r *resourceAutomationActionsRunner) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{Required: true},
			"runner_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("sidecar", "runbook"),
				},
				PlanModifiers: []planmodifier.String{
					// Requires creation of new resource while support for update is not implemented
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description":      schema.StringAttribute{Optional: true},
			"runbook_base_uri": schema.StringAttribute{Optional: true},
			"runbook_api_key":  schema.StringAttribute{Optional: true, Sensitive: true},
			"last_seen": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"creation_time": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *resourceAutomationActionsRunner) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	runnerTypePath := path.Root("runner_type")
	runnerTypeValue := types.StringValue("runbook")
	return []resource.ConfigValidator{
		validate.RequireAIfBEqual(path.Root("description"), runnerTypePath, runnerTypeValue),
		validate.RequireAIfBEqual(path.Root("description"), runnerTypePath, runnerTypeValue),
		validate.RequireAIfBEqual(path.Root("description"), runnerTypePath, runnerTypeValue),
	}
}

func (r *resourceAutomationActionsRunner) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	runnerTypePath := path.Root("runner_type")

	var runnerType types.String
	d := req.Config.GetAttribute(ctx, runnerTypePath, &runnerType)
	if resp.Diagnostics.Append(d...); d.HasError() {
		return
	}

	if runnerType.ValueString() != "runbook" {
		resp.Diagnostics.AddAttributeError(runnerTypePath, "", "only runners of runner_type runbook can be created")
		return
	}
}

func (r *resourceAutomationActionsRunner) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceAutomationActionsRunnerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyAutomationActionsRunner(&model)
	log.Printf("[INFO] Creating PagerDuty automation actions runner %s", plan.Name)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := r.client.CreateAutomationActionsRunnerWithContext(ctx, plan)
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
			fmt.Sprintf("Error creating PagerDuty automation actions runner %s", plan.Name),
			err.Error(),
		)
		return
	}

	model, err = requestGetAutomationActionsRunner(ctx, r.client, plan.ID, &plan.RunbookAPIKey, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation actions runner %s", plan.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAutomationActionsRunner) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	runbookAPIKey := extractString(ctx, req.State, "runbook_api_key", &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty automation actions runner %s", id)

	state, err := requestGetAutomationActionsRunner(ctx, r.client, id.ValueString(), runbookAPIKey, false)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation actions runner %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAutomationActionsRunner) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceAutomationActionsRunnerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyAutomationActionsRunner(&model)
	log.Printf("[INFO] Updating PagerDuty automation actions runner %s", plan.ID)

	runner, err := r.client.UpdateAutomationActionsRunnerWithContext(ctx, plan)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty automation actions runner %s", plan.ID),
			err.Error(),
		)
		return
	}
	model = flattenAutomationActionsRunner(runner, &plan.RunbookAPIKey)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAutomationActionsRunner) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty automation actions runner %s", id)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		if err := r.client.DeleteAutomationActionsRunnerWithContext(ctx, id.ValueString()); err != nil {
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
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty automation actions runner %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceAutomationActionsRunner) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAutomationActionsRunner) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceAutomationActionsRunnerModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	RunnerType     types.String `tfsdk:"runner_type"`
	Description    types.String `tfsdk:"description"`
	RunbookBaseURI types.String `tfsdk:"runbook_base_uri"`
	RunbookAPIKey  types.String `tfsdk:"runbook_api_key"`
	LastSeen       types.String `tfsdk:"last_seen"`
	CreationTime   types.String `tfsdk:"creation_time"`
}

func requestGetAutomationActionsRunner(ctx context.Context, client *pagerduty.Client, id string, runbookAPIKey *string, retryNotFound bool) (resourceAutomationActionsRunnerModel, error) {
	var model resourceAutomationActionsRunnerModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		runner, err := client.GetAutomationActionsRunnerWithContext(ctx, id)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenAutomationActionsRunner(runner, runbookAPIKey)
		return nil
	})

	return model, err
}

func buildPagerdutyAutomationActionsRunner(model *resourceAutomationActionsRunnerModel) pagerduty.AutomationActionsRunner {
	runner := pagerduty.AutomationActionsRunner{
		Name:       model.Name.ValueString(),
		RunnerType: model.RunnerType.ValueString(),
	}

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		runner.Description = model.Description.ValueString()
	}

	if !model.RunbookBaseURI.IsNull() && !model.RunbookBaseURI.IsUnknown() {
		runner.RunbookBaseURI = model.RunbookBaseURI.ValueString()
	}

	if !model.RunbookAPIKey.IsNull() && !model.RunbookAPIKey.IsUnknown() {
		runner.RunbookAPIKey = model.RunbookAPIKey.ValueString()
	}

	runner.ID = model.ID.ValueString()
	return runner
}

func flattenAutomationActionsRunner(response *pagerduty.AutomationActionsRunner, runbookAPIKey *string) resourceAutomationActionsRunnerModel {
	model := resourceAutomationActionsRunnerModel{
		ID:           types.StringValue(response.ID),
		Name:         types.StringValue(response.Name),
		Type:         types.StringValue(response.Type),
		RunnerType:   types.StringValue(response.RunnerType),
		CreationTime: types.StringValue(response.CreationTime),
		LastSeen:     types.StringValue(response.LastSeen),
	}

	if response.Description != "" {
		model.Description = types.StringValue(response.Description)
	}

	if response.RunbookBaseURI != "" {
		model.RunbookBaseURI = types.StringValue(response.RunbookBaseURI)
	}

	if response.RunbookAPIKey != "" {
		model.RunbookAPIKey = types.StringValue(response.RunbookAPIKey)
	} else if runbookAPIKey != nil {
		model.RunbookAPIKey = types.StringValue(*runbookAPIKey)
	}

	return model
}
