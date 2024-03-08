package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceTeam struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceTeam)(nil)
	_ resource.ResourceWithImportState = (*resourceTeam)(nil)
)

func (r *resourceTeam) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_team"
}

func (r *resourceTeam) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Managed by Terraform"),
			},
			"html_url": schema.StringAttribute{Computed: true},
			"parent":   schema.StringAttribute{Optional: true},
			"default_role": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

func (r *resourceTeam) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceTeamModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyTeam(&model)
	log.Printf("[INFO] Creating PagerDuty team %s", plan.Name)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := r.client.CreateTeamWithContext(ctx, plan)
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
			fmt.Sprintf("Error creating PagerDuty team %s", plan.Name),
			err.Error(),
		)
		return
	}

	retryNotFound := true
	model, err = requestGetTeam(ctx, r.client, plan, retryNotFound)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty team %s", plan.Name),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceTeam) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTeamModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty team %s", state.ID)

	plan := buildPagerdutyTeam(&state)

	retryNotFound := false
	state, err := requestGetTeam(ctx, r.client, plan, retryNotFound)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty team %s", plan.ID),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceTeam) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceTeamModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyTeam(&model)
	if plan.ID == "" {
		var id string
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		plan.ID = id
	}
	log.Printf("[INFO] Updating PagerDuty team %s", plan.ID)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		team, err := r.client.UpdateTeamWithContext(ctx, plan.ID, plan)
		if err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenTeam(team, plan)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty team %s", plan.ID),
			err.Error(),
		)
		return
	}

	retryNotFound := false
	model, err = requestGetTeam(ctx, r.client, plan, retryNotFound)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty team %s", plan.ID),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceTeam) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty team %s", id)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		err := r.client.DeleteTeamWithContext(ctx, id.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty team %s", id),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *resourceTeam) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceTeam) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceTeamModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DefaultRole types.String `tfsdk:"default_role"`
	Description types.String `tfsdk:"description"`
	HTMLURL     types.String `tfsdk:"html_url"`
	Parent      types.String `tfsdk:"parent"`
}

func requestGetTeam(ctx context.Context, client *pagerduty.Client, plan *pagerduty.Team, retryNotFound bool) (resourceTeamModel, error) {
	var model resourceTeamModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		team, err := client.GetTeamWithContext(ctx, plan.ID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenTeam(team, plan)
		return nil
	})

	return model, err
}

func buildPagerdutyTeam(model *resourceTeamModel) *pagerduty.Team {
	var parent *pagerduty.APIObject
	if !model.Parent.IsNull() && !model.Parent.IsUnknown() {
		parent = &pagerduty.APIObject{
			ID:   model.Parent.ValueString(),
			Type: "team_reference",
		}
	}
	team := pagerduty.Team{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Parent:      parent,
		DefaultRole: model.DefaultRole.ValueString(),
	}
	team.ID = model.ID.ValueString()
	return &team
}

func flattenTeam(response *pagerduty.Team, plan *pagerduty.Team) resourceTeamModel {
	model := resourceTeamModel{
		ID:          types.StringValue(response.ID),
		Name:        types.StringValue(response.Name),
		Description: types.StringValue(response.Description),
		HTMLURL:     types.StringValue(response.HTMLURL),
		DefaultRole: types.StringValue(response.DefaultRole),
	}
	if plan.Parent != nil {
		model.Parent = types.StringValue(plan.Parent.ID)
	}
	if response.Parent != nil {
		model.Parent = types.StringValue(response.Parent.ID)
	}
	return model
}
