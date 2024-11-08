package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceAlertGroupingSetting struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure      = (*resourceAlertGroupingSetting)(nil)
	_ resource.ResourceWithImportState    = (*resourceAlertGroupingSetting)(nil)
	_ resource.ResourceWithValidateConfig = (*resourceAlertGroupingSetting)(nil)
)

func (r *resourceAlertGroupingSetting) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_alert_grouping_setting"
}

func (r *resourceAlertGroupingSetting) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model resourceAlertGroupingSettingModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t := pagerduty.AlertGroupingSettingType(model.Type.ValueString())

	if t == pagerduty.AlertGroupingSettingContentBasedIntelligentType || t == pagerduty.AlertGroupingSettingIntelligentType || t == pagerduty.AlertGroupingSettingTimeType {
		if len(model.Services.Elements()) > 1 {
			resp.Diagnostics.AddAttributeError(
				path.Root("services"),
				"Invalid configuration",
				fmt.Sprintf("Setting of type %q allows for only one service in the array", t),
			)
			return
		}
	}
}

func (r *resourceAlertGroupingSetting) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("Managed by Terraform"),
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"content_based",
						"content_based_intelligent",
						"intelligent",
						"time",
					),
				},
			},
			"services": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplaceIf(
						checkAlertGroupingSettingServicesRequiresReplace,
						"Requires replace when no service from previous configuration was reused.",
						"Requires replace when no service from previous configuration was reused.",
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"config": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"timeout": schema.Int64Attribute{
						Optional:   true,
						Computed:   true,
						Validators: []validator.Int64{int64validator.NoneOf(0)},
					},
					"time_window": schema.Int64Attribute{
						Optional:   true,
						Computed:   true,
						Validators: []validator.Int64{int64validator.NoneOf(0)},
					},
					"aggregate": schema.StringAttribute{
						Optional: true,
					},
					"fields": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r *resourceAlertGroupingSetting) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceAlertGroupingSettingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan := buildPagerdutyAlertGroupingSetting(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Creating PagerDuty alert grouping setting %s", plan.Name)

	r.validateServicesReuse(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := r.client.CreateAlertGroupingSetting(ctx, plan)
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
			fmt.Sprintf("Error creating PagerDuty alert grouping setting %s", plan.Name),
			err.Error(),
		)
		return
	}

	model, err = requestGetAlertGroupingSetting(ctx, r.client, plan.ID, true, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty alert grouping setting %s", plan.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAlertGroupingSetting) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty alert grouping setting %s", id)

	state, err := requestGetAlertGroupingSetting(ctx, r.client, id.ValueString(), false, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty alert grouping setting %s", id),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceAlertGroupingSetting) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceAlertGroupingSettingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := buildPagerdutyAlertGroupingSetting(ctx, &model, &resp.Diagnostics)
	log.Printf("[INFO] Updating PagerDuty alert grouping setting %s", plan.ID)

	r.validateServicesReuse(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	alertGroupingSetting, err := r.client.UpdateAlertGroupingSetting(ctx, plan)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty alert grouping setting %s", plan.Name),
			err.Error(),
		)
		return
	}
	model = flattenAlertGroupingSetting(alertGroupingSetting)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceAlertGroupingSetting) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty alert grouping setting %s", id)

	err := r.client.DeleteAlertGroupingSetting(ctx, id.ValueString())
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty alert grouping setting %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceAlertGroupingSetting) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceAlertGroupingSetting) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceAlertGroupingSetting) validateServicesReuse(ctx context.Context, plan pagerduty.AlertGroupingSetting, diags *diag.Diagnostics) {
	serviceIDs := make([]string, len(plan.Services))
	for i, s := range plan.Services {
		serviceIDs[i] = s.ID
	}

	list, err := r.client.ListAlertGroupingSettings(ctx, pagerduty.ListAlertGroupingSettingsOptions{
		ServiceIDs: serviceIDs,
	})
	if err != nil {
		diags.AddError(
			"Unable to obtain list of alert grouping settings",
			err.Error(),
		)
		return
	}

	var reused []pagerduty.AlertGroupingSetting
	if plan.ID == "" {
		for _, a := range list.AlertGroupingSettings {
			reused = append(reused, a)
		}
	} else {
		for _, a := range list.AlertGroupingSettings {
			if a.ID != plan.ID {
				reused = append(reused, a)
			}
		}
	}

	if len(reused) > 0 {
		for _, a := range reused {
			type usage struct {
				At       int
				By, ByID string
			}
			bad := []usage{}
			for _, s := range a.Services {
				for i, sid := range serviceIDs {
					if s.ID == sid {
						bad = append(bad, usage{
							At:   i,
							By:   a.Name,
							ByID: a.ID,
						})
					}
				}
			}
			for _, b := range bad {
				var agsString string
				if b.By == "" {
					agsString = fmt.Sprintf("id=%s", b.ByID)
				} else {
					agsString = fmt.Sprintf("%q [id=%s]", b.By, b.ByID)
				}
				diags.AddAttributeError(
					path.Root("services").AtListIndex(b.At),
					"This service is associated to another alert grouping setting",
					fmt.Sprintf("Alert grouping setting %s", agsString),
				)
			}
		}
	}
}

type resourceAlertGroupingSettingModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Config      types.Object `tfsdk:"config"`
	Services    types.Set    `tfsdk:"services"`
}

func requestGetAlertGroupingSetting(ctx context.Context, client *pagerduty.Client, id string, retryNotFound bool, diags *diag.Diagnostics) (resourceAlertGroupingSettingModel, error) {
	var model resourceAlertGroupingSettingModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		alertGroupingSetting, err := client.GetAlertGroupingSetting(ctx, id)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenAlertGroupingSetting(alertGroupingSetting)
		return nil
	})

	return model, err
}

func buildPagerdutyAlertGroupingSetting(ctx context.Context, model *resourceAlertGroupingSettingModel, diags *diag.Diagnostics) pagerduty.AlertGroupingSetting {
	alertGroupingSetting := pagerduty.AlertGroupingSetting{
		ID:          model.ID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Type:        pagerduty.AlertGroupingSettingType(model.Type.ValueString()),
		Config:      buildPagerdutyAlertGroupingSettingConfig(ctx, model, diags),
		Services:    buildPagerdutyAlertGroupingSettingServices(model),
	}
	return alertGroupingSetting
}

func buildPagerdutyAlertGroupingSettingConfig(ctx context.Context, model *resourceAlertGroupingSettingModel, diags *diag.Diagnostics) interface{} {
	var target struct {
		Timeout    types.Int64  `tfsdk:"timeout"`
		TimeWindow types.Int64  `tfsdk:"time_window"`
		Aggregate  types.String `tfsdk:"aggregate"`
		Fields     types.Set    `tfsdk:"fields"`
	}

	switch model.Type.ValueString() {
	case string(pagerduty.AlertGroupingSettingContentBasedType), string(pagerduty.AlertGroupingSettingContentBasedIntelligentType):
		d := model.Config.As(ctx, &target, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
		if diags.Append(d...); d.HasError() {
			return pagerduty.AlertGroupingSettingConfigContentBased{}
		}
		fields := []string{}
		diags.Append(target.Fields.ElementsAs(ctx, &fields, false)...)
		return pagerduty.AlertGroupingSettingConfigContentBased{
			TimeWindow: uint(target.TimeWindow.ValueInt64()),
			Aggregate:  target.Aggregate.ValueString(),
			Fields:     fields,
		}

	case string(pagerduty.AlertGroupingSettingIntelligentType):
		diags.Append(model.Config.As(ctx, &target, basetypes.ObjectAsOptions{})...)
		return pagerduty.AlertGroupingSettingConfigIntelligent{
			TimeWindow: uint(target.TimeWindow.ValueInt64()),
		}

	case string(pagerduty.AlertGroupingSettingTimeType):
		diags.Append(model.Config.As(ctx, &target, basetypes.ObjectAsOptions{})...)
		return pagerduty.AlertGroupingSettingConfigTime{
			Timeout: uint(target.Timeout.ValueInt64()),
		}
	}

	return nil
}

func buildPagerdutyAlertGroupingSettingServices(model *resourceAlertGroupingSettingModel) []pagerduty.AlertGroupingSettingService {
	elements := model.Services.Elements()
	list := make([]pagerduty.AlertGroupingSettingService, 0, len(elements))
	for _, e := range elements {
		v, _ := e.(types.String)
		list = append(list, pagerduty.AlertGroupingSettingService{ID: v.ValueString()})
	}
	return list
}

func flattenAlertGroupingSetting(response *pagerduty.AlertGroupingSetting) resourceAlertGroupingSettingModel {
	model := resourceAlertGroupingSettingModel{
		ID:          types.StringValue(response.ID),
		Name:        types.StringValue(response.Name),
		Description: types.StringValue(response.Description),
		Type:        types.StringValue(string(response.Type)),
		Config:      flattenAlertGroupingSettingConfig(response),
		Services:    flattenAlertGroupingSettingServices(response),
	}
	return model
}

func flattenAlertGroupingSettingConfig(response *pagerduty.AlertGroupingSetting) types.Object {
	var alertGroupingSettingConfigAttrTypes = map[string]attr.Type{
		"timeout":     types.Int64Type,
		"time_window": types.Int64Type,
		"aggregate":   types.StringType,
		"fields":      types.SetType{ElemType: types.StringType},
	}

	var obj map[string]attr.Value

	switch c := response.Config.(type) {
	case pagerduty.AlertGroupingSettingConfigContentBased:
		fields := make([]attr.Value, 0, len(c.Fields))
		for _, f := range c.Fields {
			fields = append(fields, types.StringValue(f))
		}
		tw := types.Int64Value(int64(c.TimeWindow))
		obj = map[string]attr.Value{
			"timeout":     types.Int64Null(),
			"time_window": tw,
			"aggregate":   types.StringValue(c.Aggregate),
			"fields":      types.SetValueMust(types.StringType, fields),
		}

	case pagerduty.AlertGroupingSettingConfigIntelligent:
		obj = map[string]attr.Value{
			"timeout":     types.Int64Null(),
			"time_window": types.Int64Value(int64(c.TimeWindow)),
			"aggregate":   types.StringNull(),
			"fields":      types.SetNull(types.StringType),
		}

	case pagerduty.AlertGroupingSettingConfigTime:
		obj = map[string]attr.Value{
			"timeout":     types.Int64Value(int64(c.Timeout)),
			"time_window": types.Int64Null(),
			"aggregate":   types.StringNull(),
			"fields":      types.SetNull(types.StringType),
		}
	}

	return types.ObjectValueMust(alertGroupingSettingConfigAttrTypes, obj)
}

func flattenAlertGroupingSettingServices(response *pagerduty.AlertGroupingSetting) types.Set {
	serviceIDs := make([]attr.Value, 0, len(response.Services))
	for _, s := range response.Services {
		serviceIDs = append(serviceIDs, types.StringValue(s.ID))
	}
	return types.SetValueMust(types.StringType, serviceIDs)
}

// checkAlertGroupingSettingServicesRequiresReplace forces the resource to be
// recreated when no service from previous configuration was reused.
func checkAlertGroupingSettingServicesRequiresReplace(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifier.RequiresReplaceIfFuncResponse) {
	// TODO: check other resources to see if they are also failing because
	// the API is silently triggering a deletion of a resourced planned to
	// be updated when an object referenced inside a list or set attribute
	// is deleted earlier in a `terraform apply` execution.
	noneReused := true

	var stateIDs []string
	d := req.StateValue.ElementsAs(ctx, &stateIDs, false)
	if resp.Diagnostics.Append(d...); d.HasError() {
		return
	}

	var planIDs []types.String
	d = req.PlanValue.ElementsAs(ctx, &planIDs, false)
	if resp.Diagnostics.Append(d...); d.HasError() {
		return
	}

outerLoop:
	for _, pID := range planIDs {
		for _, sID := range stateIDs {
			if pID.ValueString() == sID {
				noneReused = false
				break outerLoop
			}
		}
	}

	resp.RequiresReplace = noneReused
}
