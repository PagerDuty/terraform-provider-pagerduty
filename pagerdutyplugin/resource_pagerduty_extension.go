package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

type resourceExtension struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceExtension)(nil)
	_ resource.ResourceWithImportState = (*resourceExtension)(nil)
)

func (r *resourceExtension) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_extension"
}

func (r *resourceExtension) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Optional: true, Computed: true},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"html_url": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional: true, Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_url": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"summary": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"extension_objects": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"extension_schema": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.StringAttribute{
				Optional:   true,
				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *resourceExtension) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceExtensionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyExtension(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Creating PagerDuty extension %s", plan.Name)

	extension, err := r.client.CreateExtensionWithContext(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating extension %s", plan.Name),
			err.Error(),
		)
		return
	}
	plan.ID = extension.ID

	accessToken := buildExtensionConfigAccessToken(model.Config, &resp.Diagnostics)
	model = requestGetExtension(ctx, r.client, plan.ID, accessToken, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceExtension) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceExtensionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty extension %s", state.ID)

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		extension, err := r.client.GetExtensionWithContext(ctx, state.ID.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) || util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		accessToken := buildExtensionConfigAccessToken(state.Config, &resp.Diagnostics)
		state = flattenExtension(extension, accessToken, &resp.Diagnostics)
		return nil
	})
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading extension %s", state.ID),
			err.Error(),
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceExtension) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceExtensionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyExtension(ctx, &model, &resp.Diagnostics)
	if plan.ID == "" {
		var id string
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		plan.ID = id
	}
	log.Printf("[INFO] Updating PagerDuty extension %s", plan.ID)

	_, err := r.client.UpdateExtensionWithContext(ctx, plan.ID, plan)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating extension %s", plan.ID),
			err.Error(),
		)
		return
	}

	accessToken := buildExtensionConfigAccessToken(model.Config, &resp.Diagnostics)
	model = requestGetExtension(ctx, r.client, plan.ID, accessToken, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceExtension) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty extension %s", id)

	err := r.client.DeleteExtensionWithContext(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting extension %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceExtension) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceExtension) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	model := requestGetExtension(ctx, r.client, req.ID, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type resourceExtensionModel struct {
	Name             types.String         `tfsdk:"name"`
	Config           jsontypes.Normalized `tfsdk:"config"`
	EndpointURL      types.String         `tfsdk:"endpoint_url"`
	ExtensionObjects types.Set            `tfsdk:"extension_objects"`
	ExtensionSchema  types.String         `tfsdk:"extension_schema"`
	HTMLURL          types.String         `tfsdk:"html_url"`
	ID               types.String         `tfsdk:"id"`
	Summary          types.String         `tfsdk:"summary"`
	Type             types.String         `tfsdk:"type"`
}

func requestGetExtension(ctx context.Context, client *pagerduty.Client, id string, accessToken *string, diags *diag.Diagnostics) resourceExtensionModel {
	var model resourceExtensionModel
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		extension, err := client.GetExtensionWithContext(ctx, id)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenExtension(extension, accessToken, diags)
		return nil
	})
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading extension %s", id),
			err.Error(),
		)
	}
	return model
}

func buildPagerdutyExtension(ctx context.Context, model *resourceExtensionModel, diags *diag.Diagnostics) *pagerduty.Extension {
	extension := pagerduty.Extension{
		Name:             model.Name.ValueString(),
		Config:           buildExtensionConfig(model.Config, diags),
		EndpointURL:      model.EndpointURL.ValueString(),
		ExtensionObjects: buildExtensionObjects(ctx, model.ExtensionObjects, diags),
		ExtensionSchema:  buildExtensionSchema(model.ExtensionSchema),
	}
	extension.ID = model.ID.ValueString()
	extension.Type = "extension"
	return &extension
}

func buildExtensionSchema(s types.String) pagerduty.APIObject {
	if s.IsNull() || s.IsUnknown() {
		return pagerduty.APIObject{}
	}
	return pagerduty.APIObject{Type: "extension_schema_reference", ID: s.ValueString()}
}

func buildExtensionObjects(ctx context.Context, list types.Set, diags *diag.Diagnostics) []pagerduty.APIObject {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	values := make([]string, 0, len(list.Elements()))
	d := list.ElementsAs(ctx, &values, false)
	diags.Append(d...)

	objects := make([]pagerduty.APIObject, 0, len(values))
	for _, v := range values {
		objects = append(objects, pagerduty.APIObject{
			ID:   v,
			Type: "service_reference",
		})
	}

	return objects
}

func buildExtensionConfig(s jsontypes.Normalized, diags *diag.Diagnostics) interface{} {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}

	var config interface{}
	if err := json.Unmarshal([]byte(s.ValueString()), &config); err != nil {
		diags.AddError(
			"Could not unmarshal extension config",
			fmt.Sprintf("%v\nIn value %s", err, s),
		)
		return nil
	}

	return config
}

func buildExtensionConfigAccessToken(s jsontypes.Normalized, diags *diag.Diagnostics) *string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}

	var config interface{}
	if err := json.Unmarshal([]byte(s.ValueString()), &config); err != nil {
		diags.AddError(
			"Could not unmarshal extension config",
			fmt.Sprintf("%v\nIn value %s", err, s),
		)
		return nil
	}

	if c, ok := config.(map[string]interface{}); ok {
		if v, ok := c["access_token"].(string); ok {
			return &v
		}
	}
	return nil
}

func flattenExtension(response *pagerduty.Extension, accessToken *string, diags *diag.Diagnostics) resourceExtensionModel {
	model := resourceExtensionModel{
		ID:               types.StringValue(response.ID),
		Name:             types.StringValue(response.Name),
		HTMLURL:          types.StringValue(response.HTMLURL),
		Type:             types.StringValue(response.Type),
		Summary:          types.StringValue(response.Summary),
		EndpointURL:      types.StringValue(response.EndpointURL),
		Config:           flattenExtensionConfig(response.Config, accessToken, diags),
		ExtensionSchema:  types.StringValue(response.ExtensionSchema.ID),
		ExtensionObjects: flattenExtensionObjects(response.ExtensionObjects, diags),
	}
	return model
}

func flattenExtensionConfig(config interface{}, accessToken *string, diags *diag.Diagnostics) jsontypes.Normalized {
	if c, ok := config.(map[string]interface{}); ok {
		if accessToken == nil {
			delete(c, "access_token")
		} else {
			c["access_token"] = *accessToken
		}
		config = c
	}

	buf, err := json.Marshal(config)
	if err != nil {
		diags.AddError(
			"Could not marshal extension config",
			fmt.Sprintf("%v\n%v", err, config),
		)
		return jsontypes.NewNormalizedNull()
	}
	return jsontypes.NewNormalizedValue(string(buf))
}

func flattenExtensionObjects(objects []pagerduty.APIObject, diags *diag.Diagnostics) types.Set {
	values := []attr.Value{}
	for _, o := range objects {
		// only flatten service_reference types, because that's all we
		// send at this time
		if o.Type == "service_reference" {
			values = append(values, types.StringValue(o.ID))
		}
	}

	list, d := types.SetValue(types.StringType, values)
	diags.Append(d...)
	return list
}
