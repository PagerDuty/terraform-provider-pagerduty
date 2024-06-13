package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceAddon struct{ client *pagerduty.Client }

var (
	_ resource.Resource                = (*resourceAddon)(nil)
	_ resource.ResourceWithConfigure   = (*resourceAddon)(nil)
	_ resource.ResourceWithImportState = (*resourceAddon)(nil)
)

func (r *resourceAddon) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_addon"
}

func (r *resourceAddon) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{Required: true},
			"src":  schema.StringAttribute{Required: true},
		},
	}
}

func (r *resourceAddon) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceAddonModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	addon := buildAddon(model)
	log.Printf("[INFO] Creating PagerDuty add-on %s", addon.Name)

	addonResp, err := r.client.InstallAddonWithContext(ctx, addon)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating add-on %s", model.Name),
			err.Error(),
		)
		return
	}
	model = requestGetAddon(ctx, r.client, addonResp.ID, nil, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAddon) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty add-on %s", id)

	removeNotFound := func(err error) *retry.RetryError {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
		}
		return retry.RetryableError(err)
	}
	model := requestGetAddon(ctx, r.client, id.ValueString(), removeNotFound, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAddon) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceAddonModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	addon := buildAddon(model)

	if addon.ID == "" {
		var id types.String
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
		if resp.Diagnostics.HasError() {
			return
		}
		addon.ID = id.ValueString()
	}
	log.Printf("[INFO] Updating PagerDuty add-on %s", addon.ID)

	addonResp, err := r.client.UpdateAddonWithContext(ctx, addon.ID, addon)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating addon %s", model.Name),
			err.Error(),
		)
		return
	}

	model = flattenAddon(addonResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAddon) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty add-on %s", id)

	err := r.client.DeleteAddonWithContext(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating addon %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceAddon) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAddon) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceAddonModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Source types.String `tfsdk:"src"`
}

func requestGetAddon(ctx context.Context, client *pagerduty.Client, id string, handleErr func(error) *retry.RetryError, diags *diag.Diagnostics) resourceAddonModel {
	var addon *pagerduty.Addon
	err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		var err error
		addon, err = client.GetAddonWithContext(ctx, id)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if handleErr != nil {
				return handleErr(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading addon %s", id),
			err.Error(),
		)
		return resourceAddonModel{}
	}
	model := flattenAddon(addon)
	return model
}

func buildAddon(model resourceAddonModel) pagerduty.Addon {
	addon := pagerduty.Addon{
		Name: model.Name.ValueString(),
		Src:  model.Source.ValueString(),
	}
	addon.ID = model.ID.ValueString()
	addon.Type = "full_page_addon"
	return addon
}

func flattenAddon(addon *pagerduty.Addon) resourceAddonModel {
	model := resourceAddonModel{
		ID:     types.StringValue(addon.ID),
		Name:   types.StringValue(addon.Name),
		Source: types.StringValue(addon.Src),
	}
	return model
}
