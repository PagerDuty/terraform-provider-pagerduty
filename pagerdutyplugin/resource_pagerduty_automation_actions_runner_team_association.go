package pagerduty

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceAutomationActionsRunnerTeamAssociation struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceAutomationActionsRunnerTeamAssociation)(nil)
	_ resource.ResourceWithImportState = (*resourceAutomationActionsRunnerTeamAssociation)(nil)
)

func (r *resourceAutomationActionsRunnerTeamAssociation) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_automation_actions_runner_team_association"
}

func (r *resourceAutomationActionsRunnerTeamAssociation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"runner_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"team_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *resourceAutomationActionsRunnerTeamAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var runnerID, teamID types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("runner_id"), &runnerID)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("team_id"), &teamID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("[INFO] Creating PagerDuty automation runner team association %s:%s", runnerID, teamID)
	id := fmt.Sprintf("%s:%s", runnerID.ValueString(), teamID.ValueString())

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		o := pagerduty.AssociateAutomationActionsRunnerTeamOptions{
			Team: pagerduty.APIReference{
				ID:   teamID.ValueString(),
				Type: "team_reference",
			},
		}
		if _, err := r.client.AssociateAutomationActionsRunnerTeamWithContext(ctx, runnerID.ValueString(), o); err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty automation runner team association %s", id),
			err.Error(),
		)
		return
	}

	model, err := requestGetAutomationActionsRunnerTeamAssociation(ctx, r.client, id, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation runner team association %s", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAutomationActionsRunnerTeamAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty automation runner team association %s", id)

	state, err := requestGetAutomationActionsRunnerTeamAssociation(ctx, r.client, id.ValueString(), false)
	if err != nil {
		if errors.Is(err, errAutomationActionsRunnerTeamNotAssociated) || util.IsNotFoundError(err) {
			log.Printf("[WARN] Removing automation runner team association %s: %s", id, err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation runner team association %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAutomationActionsRunnerTeamAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceAutomationActionsRunnerTeamAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty automation runner team association %s", id)

	runnerID, teamID, err := util.ResourcePagerDutyParseColonCompoundID(id.ValueString())
	if err != nil {
		return
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		err := r.client.DisassociateAutomationActionsRunnerTeamWithContext(ctx, runnerID, teamID)
		if err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty automation runner team association %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceAutomationActionsRunnerTeamAssociation) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAutomationActionsRunnerTeamAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceAutomationActionsRunnerTeamAssociationModel struct {
	ID       types.String `tfsdk:"id"`
	RunnerID types.String `tfsdk:"runner_id"`
	TeamID   types.String `tfsdk:"team_id"`
}

var errAutomationActionsRunnerTeamNotAssociated = errors.New("team is not associated to this runner")

func requestGetAutomationActionsRunnerTeamAssociation(ctx context.Context, client *pagerduty.Client, id string, retryNotFound bool) (resourceAutomationActionsRunnerTeamAssociationModel, error) {
	var model resourceAutomationActionsRunnerTeamAssociationModel

	runnerID, teamID, err := util.ResourcePagerDutyParseColonCompoundID(id)
	if err != nil {
		return model, err
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := client.GetAutomationActionsRunnerTeamWithContext(ctx, runnerID, teamID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		if response.Team.ID != teamID {
			return retry.NonRetryableError(errAutomationActionsRunnerTeamNotAssociated)
		}
		model.ID = types.StringValue(id)
		model.RunnerID = types.StringValue(runnerID)
		model.TeamID = types.StringValue(teamID)
		return nil
	})

	return model, err
}
