package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceTagAssignment struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceTagAssignment)(nil)
	_ resource.ResourceWithImportState = (*resourceTagAssignment)(nil)
)

func (r *resourceTagAssignment) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_tag_assignment"
}

func (r *resourceTagAssignment) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"entity_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("users", "teams", "escalation_policies"),
				},
			},
			"entity_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tag_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *resourceTagAssignment) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceTagAssignmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	assign := buildTagAssignment(&model)
	log.Printf("[INFO] Creating PagerDuty tag assignment with tagID %s for %s entity with ID %s", assign.TagID, assign.EntityType, assign.EntityID)

	assignments := &pagerduty.TagAssignments{
		Add: []*pagerduty.TagAssignment{
			{Type: "tag_reference", TagID: assign.TagID},
		},
	}

	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		err := r.client.AssignTagsWithContext(ctx, assign.EntityType, assign.EntityID, assignments)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model.ID = flattenTagAssignmentID(assign.EntityID, assign.TagID)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty tag assignment with tagID %s for %s entity with ID %s", assign.TagID, assign.EntityType, assign.EntityID),
			err.Error(),
		)
		return
	}

	isFound := r.requestGetTagAssignents(ctx, model, &resp.Diagnostics)
	if !isFound {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceTagAssignment) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTagAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty tag assignment %s", state.ID)

	isFound := r.requestGetTagAssignents(ctx, state, &resp.Diagnostics)
	if !isFound {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceTagAssignment) requestGetTagAssignents(ctx context.Context, model resourceTagAssignmentModel, diags *diag.Diagnostics) bool {
	assign := buildTagAssignment(&model)

	isFound := r.isFoundTagAssignment(ctx, assign.EntityType, assign.EntityID, diags)
	if !isFound {
		return false
	}

	isFound = false
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		opts := pagerduty.ListTagOptions{}
		response, err := r.client.GetTagsForEntity(assign.EntityType, assign.EntityID, opts)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		for _, tag := range response.Tags {
			if tag.ID == assign.TagID {
				isFound = true
				break
			}
		}
		return nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading tags for %s entity with ID %s", assign.EntityType, assign.EntityID),
			err.Error(),
		)
	}
	if !isFound {
		return false
	}
	return true
}

func (r *resourceTagAssignment) isFoundTagAssignment(ctx context.Context, entityType, entityID string, diags *diag.Diagnostics) bool {
	isFound := false

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var err error

		switch entityType {
		case "users":
			opts := pagerduty.GetUserOptions{}
			_, err = r.client.GetUserWithContext(ctx, entityID, opts)
		case "teams":
			_, err = r.client.GetTeamWithContext(ctx, entityID)
		case "escalation_policies":
			opts := pagerduty.GetEscalationPolicyOptions{}
			_, err = r.client.GetEscalationPolicyWithContext(ctx, entityID, &opts)
		}

		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				return nil
			}
			return retry.RetryableError(err)
		}

		isFound = true
		return nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error finding %s entity with ID %s", entityType, entityID),
			err.Error(),
		)
		isFound = false
	}

	return isFound
}

func (r *resourceTagAssignment) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *resourceTagAssignment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model resourceTagAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assign := buildTagAssignment(&model)
	log.Printf("[INFO] Deleting PagerDuty tag assignment with tagID %s for entityID %s", assign.TagID, assign.EntityID)

	assignments := &pagerduty.TagAssignments{
		Remove: []*pagerduty.TagAssignment{
			{Type: "tag_reference", TagID: assign.TagID},
		},
	}

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		err := r.client.AssignTagsWithContext(ctx, assign.EntityType, assign.EntityID, assignments)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if util.IsNotFoundError(err) {
				return nil
			}
			return retry.RetryableError(err)
		}
		model.ID = flattenTagAssignmentID(assign.EntityID, assign.TagID)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty tag assignment with tagID %s for %s entity with ID %s", assign.TagID, assign.EntityType, assign.EntityID),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceTagAssignment) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceTagAssignment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, ".")
	if len(ids) != 3 {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_tag_assignment",
			"Expecting an importation ID formed as '<entity_type>.<entity_id>.<tag_id>'",
		)
		return
	}
	entityType, entityID, tagID := ids[0], ids[1], ids[2]

	tagResponse, err := r.client.GetTagsForEntity(entityType, entityID, pagerduty.ListTagOptions{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing pagerduty_tag_assignment",
			err.Error(),
		)
	}

	isFound := false
	for _, tag := range tagResponse.Tags {
		if tag.ID == tagID {
			isFound = true
			break
		}
	}
	if !isFound {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Error importing pagerduty_tag_assignment", "Tag not found for entity")
		return
	}

	state := flattenTagAssignment(entityType, entityID, tagID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

type resourceTagAssignmentModel struct {
	ID         types.String `tfsdk:"id"`
	EntityID   types.String `tfsdk:"entity_id"`
	EntityType types.String `tfsdk:"entity_type"`
	TagID      types.String `tfsdk:"tag_id"`
}

type tagAssignment struct {
	ID         string
	EntityID   string
	EntityType string
	TagID      string
}

func buildTagAssignment(model *resourceTagAssignmentModel) tagAssignment {
	return tagAssignment{
		ID:         model.ID.ValueString(),
		EntityID:   model.EntityID.ValueString(),
		EntityType: model.EntityType.ValueString(),
		TagID:      model.TagID.ValueString(),
	}
}

func flattenTagAssignment(entityType, entityID, tagID string) resourceTagAssignmentModel {
	model := resourceTagAssignmentModel{
		ID:         flattenTagAssignmentID(entityID, tagID),
		EntityID:   types.StringValue(entityID),
		EntityType: types.StringValue(entityType),
		TagID:      types.StringValue(tagID),
	}
	return model
}

func flattenTagAssignmentID(entityID, tagID string) types.String {
	// TODO: i think this should have entityType for consistency with import
	return types.StringValue(fmt.Sprintf("%v.%v", entityID, tagID))
}
