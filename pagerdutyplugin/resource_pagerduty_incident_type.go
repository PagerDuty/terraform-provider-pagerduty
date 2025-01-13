package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/validate"
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

type resourceIncidentType struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceIncidentType)(nil)
	_ resource.ResourceWithImportState = (*resourceIncidentType)(nil)
)

func (r *resourceIncidentType) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_incident_type"
}

func (r *resourceIncidentType) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validate.StringHasNoSuffix("_default"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validate.StringHasNoPrefix("PD", "PagerDuty", "Default"),
				},
			},
			"parent_type": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{Optional: true},
			"enabled":     schema.BoolAttribute{Optional: true, Computed: true},
			"type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourceIncidentType) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceIncidentTypeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := pagerduty.CreateIncidentTypeOptions{
		Name:        model.Name.ValueString(),
		DisplayName: model.DisplayName.ValueString(),
		ParentType:  model.ParentType.ValueString(),
		Description: model.Description.ValueStringPointer(),
	}
	if !model.Enabled.IsNull() && !model.Enabled.IsUnknown() {
		plan.Enabled = model.Enabled.ValueBoolPointer()
	}
	log.Printf("[INFO] Creating PagerDuty incident type %s", plan.Name)

	if list, err := r.client.ListIncidentTypes(ctx, pagerduty.ListIncidentTypesOptions{
		Filter: "disabled",
	}); err == nil {
		for _, it := range list.IncidentTypes {
			if it.Name == plan.Name {
				resp.Diagnostics.AddWarning(
					"Incident Type disabled",
					fmt.Sprintf("Incident Type with name %q already exists but it is disabled", plan.Name),
				)
			}
		}
	}

	var id string

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := r.client.CreateIncidentType(ctx, plan)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		id = response.ID
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty incident type %s", plan.Name),
			err.Error(),
		)
		return
	}

	model, err = requestGetIncidentType(ctx, r.client, id, plan.ParentType, true, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type %s", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceIncidentType) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id, parent types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty incident type %s", id)

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("parent_type"), &parent)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := requestGetIncidentType(ctx, r.client, id.ValueString(), parent.ValueString(), false, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceIncidentType) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceIncidentTypeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var parent types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("parent_type"), &parent)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := pagerduty.UpdateIncidentTypeOptions{
		DisplayName: model.DisplayName.ValueStringPointer(),
		Enabled:     model.Enabled.ValueBoolPointer(),
		Description: model.Description.ValueStringPointer(),
	}

	id := model.ID.ValueString()
	log.Printf("[INFO] Updating PagerDuty incident type %s", id)

	if parent.ValueString() != model.ParentType.ValueString() {
		resp.Diagnostics.AddWarning(
			"Can not update value of field \"parent_type\"",
			"",
		)

	}

	incidentType, err := r.client.UpdateIncidentType(ctx, id, plan)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty incident type %s", id),
			err.Error(),
		)
		return
	}

	model, err = flattenIncidentType(ctx, r.client, incidentType, model.ParentType.ValueString(), &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceIncidentType) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Cannot delete incident type",
		"This action has no effect, and you might want to disable your incident type by changing it with `enabled = false`. If you want terraform to stop tracking this resource please use `terraform state rm`.",
	)
}

func (r *resourceIncidentType) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceIncidentType) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceIncidentTypeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	ParentType  types.String `tfsdk:"parent_type"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Type        types.String `tfsdk:"type"`
}

func requestGetIncidentType(ctx context.Context, client *pagerduty.Client, id, parent string, retryNotFound bool, diags *diag.Diagnostics) (resourceIncidentTypeModel, error) {
	var model resourceIncidentTypeModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		incidentType, err := client.GetIncidentType(ctx, id, pagerduty.GetIncidentTypeOptions{})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		model, err = flattenIncidentType(ctx, client, incidentType, parent, diags)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	return model, err
}

func flattenIncidentType(ctx context.Context, client *pagerduty.Client, response *pagerduty.IncidentType, parent string, diags *diag.Diagnostics) (resourceIncidentTypeModel, error) {
	if parent != response.Parent.ID {
		incidentType, err := client.GetIncidentType(ctx, parent, pagerduty.GetIncidentTypeOptions{})
		if err != nil {
			return resourceIncidentTypeModel{}, err
		}
		if parent != incidentType.ID && parent != incidentType.Name {
			return resourceIncidentTypeModel{}, fmt.Errorf("parent_type %q was not received, got %q (ID=%s)", parent, incidentType.Name, incidentType.ID)
		}
	}

	model := resourceIncidentTypeModel{
		ID:          types.StringValue(response.ID),
		Name:        types.StringValue(response.Name),
		DisplayName: types.StringValue(response.DisplayName),
		ParentType:  types.StringValue(parent),
		Enabled:     types.BoolValue(response.Enabled),
		Type:        types.StringValue(response.Type),
	}

	if response.Description != "" {
		model.Description = types.StringValue(response.Description)
	}

	return model, nil
}
