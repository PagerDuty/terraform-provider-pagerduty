package pagerduty

import (
	"context"
	"encoding/json"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceExtensionServiceNow struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceExtensionServiceNow)(nil)
	_ resource.ResourceWithImportState = (*resourceExtensionServiceNow)(nil)
)

func (r *resourceExtensionServiceNow) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_extension_servicenow"
}

func (r *resourceExtensionServiceNow) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"name":         schema.StringAttribute{Optional: true, Computed: true},
			"html_url":     schema.StringAttribute{Computed: true},
			"type":         schema.StringAttribute{Optional: true, Computed: true},
			"endpoint_url": schema.StringAttribute{Optional: true, Sensitive: true},
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
			"snow_user":     schema.StringAttribute{Required: true},
			"snow_password": schema.StringAttribute{Required: true, Sensitive: true},
			"summary":       schema.StringAttribute{Optional: true, Computed: true},
			"sync_options": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("manual_sync", "sync_all"),
				},
			},
			"target":    schema.StringAttribute{Required: true},
			"task_type": schema.StringAttribute{Required: true},
			"referer":   schema.StringAttribute{Required: true},
		},
	}
}

func (r *resourceExtensionServiceNow) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceExtensionServiceNowModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyExtensionServiceNow(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Creating extension service now %s", plan.Name)

	extension, err := r.client.CreateExtensionWithContext(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating extension service now %s", plan.Name),
			err.Error(),
		)
		return
	}
	plan.ID = extension.ID

	model, err = r.requestGetExtensionServiceNow(ctx, requestGetExtensionServiceNowOptions{
		ID:            plan.ID,
		RetryNotFound: false,
		SnowPassword:  extractString(ctx, req.Plan, "snow_password", &resp.Diagnostics),
		EndpointURL:   extractString(ctx, req.Plan, "endpoint_url", &resp.Diagnostics),
		Diagnostics:   &resp.Diagnostics,
	})
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceExtensionServiceNow) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceExtensionServiceNowModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading extension service now %s", state.ID)
	id := state.ID.ValueString()

	state, err := r.requestGetExtensionServiceNow(ctx, requestGetExtensionServiceNowOptions{
		ID:            id,
		RetryNotFound: false,
		SnowPassword:  extractString(ctx, req.State, "snow_password", &resp.Diagnostics),
		EndpointURL:   extractString(ctx, req.State, "endpoint_url", &resp.Diagnostics),
		Diagnostics:   &resp.Diagnostics,
	})
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceExtensionServiceNow) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceExtensionServiceNowModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyExtensionServiceNow(ctx, &model, &resp.Diagnostics)
	if plan.ID == "" {
		var id string
		req.State.GetAttribute(ctx, path.Root("id"), &id)
		plan.ID = id
	}
	log.Printf("[INFO] Updating extension service now %s", plan.ID)

	_, err := r.client.UpdateExtensionWithContext(ctx, plan.ID, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating extension service now %s", plan.ID),
			err.Error(),
		)
		return
	}

	model, err = r.requestGetExtensionServiceNow(ctx, requestGetExtensionServiceNowOptions{
		ID:            plan.ID,
		RetryNotFound: true,
		SnowPassword:  extractString(ctx, req.Plan, "snow_password", &resp.Diagnostics),
		EndpointURL:   extractString(ctx, req.Plan, "endpoint_url", &resp.Diagnostics),
		Diagnostics:   &resp.Diagnostics,
	})
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceExtensionServiceNow) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting extension service now %s", id)

	err := r.client.DeleteExtensionWithContext(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting extension service now %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceExtensionServiceNow) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceExtensionServiceNow) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	model, err := r.requestGetExtensionServiceNow(ctx, requestGetExtensionServiceNowOptions{
		ID:            req.ID,
		RetryNotFound: false,
		Diagnostics:   &resp.Diagnostics,
	})
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
		}
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type resourceExtensionServiceNowModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	HTMLURL          types.String `tfsdk:"html_url"`
	Type             types.String `tfsdk:"type"`
	EndpointURL      types.String `tfsdk:"endpoint_url"`
	ExtensionObjects types.Set    `tfsdk:"extension_objects"`
	ExtensionSchema  types.String `tfsdk:"extension_schema"`
	SnowUser         types.String `tfsdk:"snow_user"`
	SnowPassword     types.String `tfsdk:"snow_password"`
	Summary          types.String `tfsdk:"summary"`
	SyncOptions      types.String `tfsdk:"sync_options"`
	Target           types.String `tfsdk:"target"`
	TaskType         types.String `tfsdk:"task_type"`
	Referer          types.String `tfsdk:"referer"`
}

type requestGetExtensionServiceNowOptions struct {
	ID            string
	RetryNotFound bool
	SnowPassword  *string
	EndpointURL   *string
	Diagnostics   *diag.Diagnostics
}

func (r *resourceExtensionServiceNow) requestGetExtensionServiceNow(ctx context.Context, opts requestGetExtensionServiceNowOptions) (resourceExtensionServiceNowModel, error) {
	var model resourceExtensionServiceNowModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		extensionServiceNow, err := r.client.GetExtensionWithContext(ctx, opts.ID)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !opts.RetryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenExtensionServiceNow(extensionServiceNow, opts.SnowPassword, opts.EndpointURL)
		return nil
	})
	if err != nil && (opts.RetryNotFound || !util.IsNotFoundError(err)) {
		opts.Diagnostics.AddError(
			fmt.Sprintf("Error reading extension service now %s", opts.ID),
			err.Error(),
		)
	}

	return model, err
}

func buildPagerdutyExtensionServiceNow(ctx context.Context, model *resourceExtensionServiceNowModel, diags *diag.Diagnostics) *pagerduty.Extension {
	config := &pagerDutyExtensionServiceNowConfig{
		User:        model.SnowUser.ValueString(),
		Password:    model.SnowPassword.ValueString(),
		SyncOptions: model.SyncOptions.ValueString(),
		Target:      model.Target.ValueString(),
		TaskType:    model.TaskType.ValueString(),
		Referer:     model.Referer.ValueString(),
	}
	extensionServiceNow := pagerduty.Extension{
		APIObject: pagerduty.APIObject{
			ID:   model.ID.ValueString(),
			Type: "extension",
		},
		Name:        model.Name.ValueString(),
		EndpointURL: model.EndpointURL.ValueString(),
		ExtensionSchema: pagerduty.APIObject{
			ID:   model.ExtensionSchema.ValueString(),
			Type: "extension_schema_reference",
		},
		ExtensionObjects: buildExtensionServiceNowObjects(ctx, model.ExtensionObjects, diags),
		Config:           config,
	}
	return &extensionServiceNow
}

func buildExtensionServiceNowObjects(ctx context.Context, set types.Set, diags *diag.Diagnostics) []pagerduty.APIObject {
	var target []string
	diags.Append(set.ElementsAs(ctx, &target, false)...)

	list := []pagerduty.APIObject{}
	for _, s := range target {
		list = append(list, pagerduty.APIObject{Type: "service_reference", ID: s})
	}

	return list
}

func flattenExtensionServiceNow(src *pagerduty.Extension, snowPassword *string, endpointURL *string) resourceExtensionServiceNowModel {
	model := resourceExtensionServiceNowModel{
		ID:               types.StringValue(src.ID),
		Name:             types.StringValue(src.Name),
		HTMLURL:          types.StringValue(src.HTMLURL),
		ExtensionSchema:  types.StringValue(src.ExtensionSchema.ID),
		ExtensionObjects: flattenExtensionServiceNowObjects(src.ExtensionObjects),
	}

	b, _ := json.Marshal(src.Config)
	var config pagerDutyExtensionServiceNowConfig
	_ = json.Unmarshal(b, &config)

	model.SnowUser = types.StringValue(config.User)
	if snowPassword != nil {
		model.SnowPassword = types.StringValue(*snowPassword)
	} else if config.Password != "" {
		model.SnowPassword = types.StringValue(config.Password)
	}
	if endpointURL != nil {
		model.EndpointURL = types.StringValue(*endpointURL)
	} else if src.EndpointURL != "" {
		model.EndpointURL = types.StringValue(src.EndpointURL)
	}
	model.SyncOptions = types.StringValue(config.SyncOptions)
	model.Target = types.StringValue(config.Target)
	model.TaskType = types.StringValue(config.TaskType)
	model.Referer = types.StringValue(config.Referer)
	return model
}

func flattenExtensionServiceNowObjects(list []pagerduty.APIObject) types.Set {
	elements := make([]attr.Value, 0, len(list))
	for _, s := range list {
		if s.Type == "service_reference" {
			elements = append(elements, types.StringValue(s.ID))
		}
	}
	return types.SetValueMust(types.StringType, elements)
}

type pagerDutyExtensionServiceNowConfig struct {
	User        string `json:"snow_user"`
	Password    string `json:"snow_password,omitempty"`
	SyncOptions string `json:"sync_options"`
	Target      string `json:"target"`
	TaskType    string `json:"task_type"`
	Referer     string `json:"referer"`
}
