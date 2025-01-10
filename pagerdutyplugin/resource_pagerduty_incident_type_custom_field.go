package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type resourceIncidentTypeCustomField struct{ client *pagerduty.Client }

var (
	_ resource.ResourceWithConfigure   = (*resourceIncidentTypeCustomField)(nil)
	_ resource.ResourceWithImportState = (*resourceIncidentTypeCustomField)(nil)
)

func (r *resourceIncidentTypeCustomField) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_incident_type_custom_field"
}

func (r *resourceIncidentTypeCustomField) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model resourceIncidentTypeCustomFieldModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if model.FieldType.ValueString() == "single_value_fixed" || model.FieldType.ValueString() == "multi_value_fixed" {
		if len(model.FieldOptions.Elements()) == 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("field_options"),
				"Invalid Value",
				"`field_options` can't be empty when `field_type` is `single_value_fixed` or `multi_value_fixed`",
			)
		}
	} else if model.FieldType.ValueString() == "single_value" || model.FieldType.ValueString() == "multi_value" {
		if !model.FieldOptions.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("field_options"),
				"Invalid Value",
				"field_options: not allowed for field type "+model.FieldType.ValueString(),
			)
		}
	}

	if model.FieldType.ValueString() != "single_value" {
		// The API error response indicates that data_type != string can be single_value or
		// single_value_fixed but whenever trying, it fails
		if model.DataType.ValueString() != "string" {
			resp.Diagnostics.AddAttributeError(
				path.Root("field_type"),
				"Invalid Value",
				"field_type must be single_value when data_type is "+model.DataType.ValueString(),
			)
		}
	}
}

func (r *resourceIncidentTypeCustomField) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"incident_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:    []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"display_name": schema.StringAttribute{
				Required: true,
			},
			"data_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"field_type": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("single_value", "multi_value", "single_value_fixed", "multi_value_fixed"),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"default_value": schema.StringAttribute{
				Optional:   true,
				CustomType: jsontypes.NormalizedType{},
			},
			"enabled": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"field_options": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"summary": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"self": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *resourceIncidentTypeCustomField) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model resourceIncidentTypeCustomFieldModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Creating PagerDuty incident type custom field %s", model.Name)

	var fieldOptions []pagerduty.IncidentTypeFieldOption
	{
		var target []types.String
		d := model.FieldOptions.ElementsAs(ctx, &target, true)
		if resp.Diagnostics.Append(d...); d.HasError() {
			return
		}
		for _, t := range target {
			opt := pagerduty.IncidentTypeFieldOption{
				Data: &pagerduty.IncidentTypeFieldOptionData{
					Value:    t.ValueString(),
					DataType: model.DataType.ValueString(),
				},
				Type: "field_option",
			}
			fieldOptions = append(fieldOptions, opt)
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var defaultValue any
	if !model.DefaultValue.IsNull() && !model.DefaultValue.IsUnknown() {
		v := model.DefaultValue.ValueString()
		if err := json.Unmarshal([]byte(v), &defaultValue); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_value"),
				"Error parsing",
				err.Error(),
			)
			return
		}
	}

	enabled := model.Enabled.ValueBoolPointer()
	if model.Enabled.IsUnknown() {
		enabled = nil
	}

	plan := pagerduty.CreateIncidentTypeFieldOptions{
		Name:         model.Name.ValueString(),
		DisplayName:  model.DisplayName.ValueString(),
		DataType:     model.DataType.ValueString(),
		FieldType:    model.FieldType.ValueString(),
		DefaultValue: defaultValue,
		Description:  model.Description.ValueStringPointer(),
		Enabled:      enabled,
		FieldOptions: fieldOptions,
	}

	incidentTypeID := model.IncidentType.ValueString()
	var fieldID string

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := r.client.CreateIncidentTypeField(ctx, model.IncidentType.ValueString(), plan)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		fieldID = response.ID
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating PagerDuty incident type custom field %s", plan.Name),
			err.Error(),
		)
		return
	}

	model, err = requestGetIncidentTypeCustomField(ctx, r.client, incidentTypeID, fieldID, false, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type custom field %s for incident type %s", fieldID, incidentTypeID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceIncidentTypeCustomField) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Reading PagerDuty incident type custom field %s", id)

	incidentTypeID, fieldID, err := util.ResourcePagerDutyParseColonCompoundID(id.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}

	state, err := requestGetIncidentTypeCustomField(ctx, r.client, incidentTypeID, fieldID, false, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type custom field %s", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceIncidentTypeCustomField) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model resourceIncidentTypeCustomFieldModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Updating PagerDuty incident type custom field %s", model.ID)

	var defaultValue interface{}
	if !model.DefaultValue.IsNull() && !model.DefaultValue.IsUnknown() {
		v := model.DefaultValue.ValueString()
		json.Unmarshal([]byte(v), &defaultValue)
	}

	plan := pagerduty.UpdateIncidentTypeFieldOptions{
		DisplayName:  model.DisplayName.ValueStringPointer(),
		DefaultValue: &defaultValue,
		Description:  model.Description.ValueStringPointer(),
		Enabled:      model.Enabled.ValueBoolPointer(),
	}

	incidentTypeID, fieldID, err := util.ResourcePagerDutyParseColonCompoundID(model.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}

	var planFieldOptions, stateFieldOptions []string
	{
		var set types.Set
		d := req.State.GetAttribute(ctx, path.Root("field_options"), &set)
		if resp.Diagnostics.Append(d...); d.HasError() {
			return
		}
		d = set.ElementsAs(ctx, &stateFieldOptions, true)
		if resp.Diagnostics.Append(d...); d.HasError() {
			return
		}
	}
	{
		d := model.FieldOptions.ElementsAs(ctx, &planFieldOptions, true)
		if resp.Diagnostics.Append(d...); d.HasError() {
			return
		}
	}
	additions, deletions := util.CalculateDiff(stateFieldOptions, planFieldOptions)

	field, err := r.client.GetIncidentTypeField(ctx, incidentTypeID, fieldID, pagerduty.GetIncidentTypeFieldOptions{
		Includes: []string{"field_options"},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty incident type custom field %s", model.ID),
			err.Error(),
		)
	}

	db := make(map[string]string)
	for _, opt := range field.FieldOptions {
		if opt.Data == nil {
			continue
		}
		db[opt.Data.Value] = opt.ID
	}

	if _, err := r.client.UpdateIncidentTypeField(ctx, incidentTypeID, fieldID, plan); err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating PagerDuty incident type custom field %s", model.ID),
			err.Error(),
		)
		return
	}

	for _, opt := range additions {
		_, err := r.client.CreateIncidentTypeFieldOption(ctx, incidentTypeID, fieldID, pagerduty.CreateIncidentTypeFieldOptionPayload{
			Data: &pagerduty.IncidentTypeFieldOptionData{
				Value:    opt,
				DataType: field.DataType,
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error updating PagerDuty incident type custom field %s", model.ID),
				err.Error(),
			)
		}
	}

	for _, opt := range deletions {
		err := r.client.DeleteIncidentTypeFieldOption(ctx, incidentTypeID, field.ID, db[opt])
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error updating PagerDuty incident type custom field %s", model.ID),
				err.Error(),
			)
		}
	}

	model, err = requestGetIncidentTypeCustomField(ctx, r.client, incidentTypeID, fieldID, false, &resp.Diagnostics)
	if err != nil {
		if util.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type custom field %s", model.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *resourceIncidentTypeCustomField) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[INFO] Deleting PagerDuty incident type custom field %s", id)

	incidentTypeID, fieldID, err := util.ResourcePagerDutyParseColonCompoundID(id.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Invalid Value", err.Error())
		return
	}

	err = r.client.DeleteIncidentTypeField(ctx, incidentTypeID, fieldID)
	if err != nil && !util.IsNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting PagerDuty incident type custom field %s", id),
			err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *resourceIncidentTypeCustomField) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&r.client, req.ProviderData)...)
}

func (r *resourceIncidentTypeCustomField) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type resourceIncidentTypeCustomFieldModel struct {
	ID           types.String         `tfsdk:"id"`
	Enabled      types.Bool           `tfsdk:"enabled"`
	Name         types.String         `tfsdk:"name"`
	Type         types.String         `tfsdk:"type"`
	Self         types.String         `tfsdk:"self"`
	Description  types.String         `tfsdk:"description"`
	FieldType    types.String         `tfsdk:"field_type"`
	DataType     types.String         `tfsdk:"data_type"`
	DisplayName  types.String         `tfsdk:"display_name"`
	DefaultValue jsontypes.Normalized `tfsdk:"default_value"`
	IncidentType types.String         `tfsdk:"incident_type"`
	Summary      types.String         `tfsdk:"summary"`
	FieldOptions types.Set            `tfsdk:"field_options"`
}

func requestGetIncidentTypeCustomField(ctx context.Context, client *pagerduty.Client, incidentTypeID, fieldID string, retryNotFound bool, diags *diag.Diagnostics) (resourceIncidentTypeCustomFieldModel, error) {
	var model resourceIncidentTypeCustomFieldModel

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		field, err := client.GetIncidentTypeField(ctx, incidentTypeID, fieldID, pagerduty.GetIncidentTypeFieldOptions{
			Includes: []string{"field_options"},
		})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			if !retryNotFound && util.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		model = flattenIncidentTypeCustomField(field)
		return nil
	})

	return model, err
}

func flattenIncidentTypeCustomField(response *pagerduty.IncidentTypeField) resourceIncidentTypeCustomFieldModel {
	fieldOptions := types.SetNull(types.StringType)
	if len(response.FieldOptions) > 0 {
		list := make([]attr.Value, 0, len(response.FieldOptions))
		for _, opt := range response.FieldOptions {
			if opt.Data != nil {
				v := types.StringValue(opt.Data.Value)
				list = append(list, v)
			}
		}
		fieldOptions = types.SetValueMust(types.StringType, list)
	}

	model := resourceIncidentTypeCustomFieldModel{
		ID:           types.StringValue(response.IncidentType + ":" + response.ID),
		Enabled:      types.BoolValue(response.Enabled),
		Self:         types.StringNull(),
		Name:         types.StringValue(response.Name),
		Type:         types.StringValue(response.Type),
		FieldType:    types.StringValue(response.FieldType),
		DataType:     types.StringValue(response.DataType),
		DisplayName:  types.StringValue(response.DisplayName),
		DefaultValue: jsontypes.NewNormalizedNull(),
		IncidentType: types.StringValue(response.IncidentType),
		Summary:      types.StringValue(response.Summary),
		FieldOptions: fieldOptions,
	}

	if response.Description != "" {
		model.Description = types.StringValue(response.Description)
	}

	if response.DefaultValue != nil {
		buf, err := json.Marshal(response.DefaultValue)
		if err == nil {
			model.DefaultValue = jsontypes.NewNormalizedValue(string(buf))
		}
	}

	if response.Self != "" {
		model.Self = types.StringValue(response.Self)
	}

	return model
}
