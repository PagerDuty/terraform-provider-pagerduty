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

type resourceAutomationActionsActionTeamAssociation struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceAutomationActionsActionTeamAssociation)(nil)
	_ resource.ResourceWithImportState = (*resourceAutomationActionsActionTeamAssociation)(nil)
)

func (r *resourceAutomationActionsActionTeamAssociation) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_automation_actions_action_team_association"
}

func (r *resourceAutomationActionsActionTeamAssociation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"action_id": schema.StringAttribute{
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

func (r *resourceAutomationActionsActionTeamAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var actionID, teamID types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("action_id"), &actionID)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("team_id"), &teamID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("[INFO] Creating PagerDuty automation action team association %s:%s", actionID, teamID)
	id := fmt.Sprintf("%s:%s", actionID.ValueString(), teamID.ValueString())

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		o := pagerduty.AssociateAutomationActionTeamOptions{
			Team: pagerduty.APIReference{
				ID:   teamID.ValueString(),
				Type: "team_reference",
			},
		}
		if _, err := r.client.AssociateAutomationActionTeamWithContext(ctx, actionID.ValueString(), o); err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty automation action team association %s", id),
			err.Error(),
		)
		return
	}

	model, err := requestGetAutomationActionsActionTeamAssociation(ctx, r.client, id, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action team association %s", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAutomationActionsActionTeamAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty automation action team association %s", id)

	state, err := requestGetAutomationActionsActionTeamAssociation(ctx, r.client, id.ValueString(), false)
	if err != nil {
		if errors.Is(err, errAutomationActionTeamNotAssociated) || util.IsNotFoundError(err) {
			log.Printf("[WARN] Removing automation action team association %s: %s", id, err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation action team association %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAutomationActionsActionTeamAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceAutomationActionsActionTeamAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty automation action team association %s", id)

	actionID, teamID, err := util.ResourcePagerDutyParseColonCompoundID(id.ValueString())
	if err != nil {
		return
	}

	err = r.client.DisassociateAutomationActionTeamWithContext(ctx, actionID, teamID)
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty automation action team association %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceAutomationActionsActionTeamAssociation) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAutomationActionsActionTeamAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceAutomationActionsActionTeamAssociationModel struct {
	ID       types.String `tfsdk:"id"`
	ActionID types.String `tfsdk:"action_id"`
	TeamID   types.String `tfsdk:"team_id"`
}

var errAutomationActionTeamNotAssociated = errors.New("team is not associated to this action")

func requestGetAutomationActionsActionTeamAssociation(ctx context.Context, client *pagerduty.Client, id string, retryNotFound bool) (resourceAutomationActionsActionTeamAssociationModel, error) {
	var model resourceAutomationActionsActionTeamAssociationModel

	actionID, teamID, err := util.ResourcePagerDutyParseColonCompoundID(id)
	if err != nil {
		return model, err
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := client.GetAutomationActionTeamWithContext(ctx, actionID, teamID)
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
			return retry.NonRetryableError(errAutomationActionTeamNotAssociated)
		}
		model.ID = types.StringValue(id)
		model.ActionID = types.StringValue(actionID)
		model.TeamID = types.StringValue(teamID)
		return nil
	})

	return model, err
}
