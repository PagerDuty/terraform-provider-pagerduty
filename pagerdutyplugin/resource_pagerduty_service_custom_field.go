package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ServiceCustomFieldResource{}
	_ resource.ResourceWithConfigure   = &ServiceCustomFieldResource{}
	_ resource.ResourceWithImportState = &ServiceCustomFieldResource{}
)

type ServiceCustomFieldResource struct {
	client *pagerduty.Client
}

type ServiceCustomFieldResourceModel struct {
	ID           types.String         `tfsdk:"id"`
	Name         types.String         `tfsdk:"name"`
	DisplayName  types.String         `tfsdk:"display_name"`
	Description  types.String         `tfsdk:"description"`
	DataType     types.String         `tfsdk:"data_type"`
	FieldType    types.String         `tfsdk:"field_type"`
	DefaultValue jsontypes.Normalized `tfsdk:"default_value"`
	Enabled      types.Bool           `tfsdk:"enabled"`
	FieldOptions types.List           `tfsdk:"field_option"`
}

type ServiceCustomFieldOptionModel struct {
	ID       types.String `tfsdk:"id"`
	Value    types.String `tfsdk:"value"`
	DataType types.String `tfsdk:"data_type"`
}

func (r *ServiceCustomFieldResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_service_custom_field"
}

func (r *ServiceCustomFieldResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A resource to manage PagerDuty Service Custom Fields. Note: This is an Early Access feature that requires the X-EARLY-ACCESS header with service-custom-fields-preview value for all API operations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the field. May include ASCII characters, specifically lowercase letters, digits, and underescores. The `name` for a Field must be unique and cannot be changed once created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The human-readable name of the field. This must be unique across an account.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the data this field contains.",
				Optional:    true,
			},
			"data_type": schema.StringAttribute{
				Description: "The kind of data the custom field is allowed to contain.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"field_type": schema.StringAttribute{
				Description: "The type of data this field contains. In combination with the `data_type` field.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_value": schema.StringAttribute{
				Description: "Default value for the field.",
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the field is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
		Blocks: map[string]schema.Block{
			"field_option": schema.ListNestedBlock{
				Description: "The options for the custom field. Applies only to `single_value_fixed` and `multi_value_fixed` field types.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"data_type": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{stringvalidator.OneOf("string")},
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *ServiceCustomFieldResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pagerduty.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pagerduty.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ServiceCustomFieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceCustomFieldResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	field := pagerduty.ServiceCustomField{
		Name:        plan.Name.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		DataType:    pagerduty.ServiceCustomFieldDataType(plan.DataType.ValueString()),
		FieldType:   pagerduty.ServiceCustomFieldType(plan.FieldType.ValueString()),
		Enabled:     plan.Enabled.ValueBool(),
	}

	if !plan.Description.IsNull() {
		field.Description = plan.Description.ValueString()
	}

	if !plan.DefaultValue.IsNull() {
		var v any
		err := json.Unmarshal([]byte(plan.DefaultValue.ValueString()), &v)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_value"),
				"Invalid value",
				"Invalid JSON string: "+err.Error(),
			)
			return
		}
		field.DefaultValue = v
	}

	if !plan.FieldOptions.IsNull() {
		var fieldOptions []ServiceCustomFieldOptionModel
		resp.Diagnostics.Append(plan.FieldOptions.ElementsAs(ctx, &fieldOptions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		pdFieldOptions := make([]pagerduty.ServiceCustomFieldOption, 0, len(fieldOptions))
		for _, option := range fieldOptions {
			pdOption := pagerduty.ServiceCustomFieldOption{
				Data: pagerduty.ServiceCustomFieldOptionData{
					DataType: pagerduty.ServiceCustomFieldDataType(plan.DataType.ValueString()),
					Value:    option.Value.ValueString(),
				},
			}
			pdFieldOptions = append(pdFieldOptions, pdOption)
		}
		field.FieldOptions = pdFieldOptions
	}

	createdField, err := r.client.CreateServiceCustomField(ctx, &field)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating PagerDuty Service Custom Field",
			fmt.Sprintf("Could not create service custom field: %s", err),
		)
		return
	}

	// Map response body to model
	plan.ID = types.StringValue(createdField.ID)
	plan.Name = types.StringValue(createdField.Name)
	plan.DisplayName = types.StringValue(createdField.DisplayName)
	plan.DataType = types.StringValue(string(createdField.DataType))
	plan.FieldType = types.StringValue(string(createdField.FieldType))
	plan.Enabled = types.BoolValue(createdField.Enabled)

	if createdField.Description != "" {
		plan.Description = types.StringValue(createdField.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if createdField.DefaultValue != nil {
		buf, _ := json.Marshal(createdField.DefaultValue)
		plan.DefaultValue = jsontypes.NewNormalizedValue(string(buf))
	} else {
		plan.DefaultValue = jsontypes.NewNormalizedNull()
	}

	if len(createdField.FieldOptions) > 0 {
		fieldOptions := make([]ServiceCustomFieldOptionModel, 0, len(createdField.FieldOptions))
		for _, option := range createdField.FieldOptions {
			fieldOptions = append(fieldOptions, ServiceCustomFieldOptionModel{
				ID:       types.StringValue(option.ID),
				Value:    types.StringValue(option.Data.Value),
				DataType: types.StringValue(string(option.Data.DataType)),
			})
		}

		fieldOptionsList, diags := types.ListValueFrom(ctx, serviceCustomFieldObjectType, fieldOptions)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.FieldOptions = fieldOptionsList
	} else {
		plan.FieldOptions = types.ListNull(serviceCustomFieldObjectType)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceCustomFieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceCustomFieldResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	options := pagerduty.ListServiceCustomFieldsOptions{
		Include: []string{"field_options"},
	}

	field, err := r.client.GetServiceCustomField(ctx, state.ID.ValueString(), options)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading PagerDuty Service Custom Field",
			fmt.Sprintf("Could not read service custom field %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(field.ID)
	state.Name = types.StringValue(field.Name)
	state.DisplayName = types.StringValue(field.DisplayName)
	state.DataType = types.StringValue(string(field.DataType))
	state.FieldType = types.StringValue(string(field.FieldType))
	state.Enabled = types.BoolValue(field.Enabled)

	if field.Description != "" {
		state.Description = types.StringValue(field.Description)
	} else {
		state.Description = types.StringNull()
	}

	if field.DefaultValue != nil {
		buf, _ := json.Marshal(field.DefaultValue)
		state.DefaultValue = jsontypes.NewNormalizedValue(string(buf))
	} else {
		state.DefaultValue = jsontypes.NewNormalizedNull()
	}

	if len(field.FieldOptions) > 0 {
		fieldOptions := make([]ServiceCustomFieldOptionModel, 0, len(field.FieldOptions))
		for _, option := range field.FieldOptions {
			fieldOptions = append(fieldOptions, ServiceCustomFieldOptionModel{
				ID:       types.StringValue(option.ID),
				Value:    types.StringValue(option.Data.Value),
				DataType: types.StringValue(string(option.Data.DataType)),
			})
		}

		fieldOptionsList, diags := types.ListValueFrom(ctx, serviceCustomFieldObjectType, fieldOptions)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.FieldOptions = fieldOptionsList
	} else {
		state.FieldOptions = types.ListNull(serviceCustomFieldObjectType)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceCustomFieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceCustomFieldResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	field := pagerduty.ServiceCustomField{
		APIObject: pagerduty.APIObject{
			ID: plan.ID.ValueString(),
		},
		DisplayName: plan.DisplayName.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	if !plan.Description.IsNull() {
		field.Description = plan.Description.ValueString()
	}

	if !plan.DefaultValue.IsNull() {
		var v any
		err := json.Unmarshal([]byte(plan.DefaultValue.ValueString()), &v)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_value"),
				"Invalid value",
				"Invalid JSON string: "+err.Error(),
			)
			return
		}
		field.DefaultValue = v
	}

	if !plan.FieldOptions.IsNull() {
		var fieldOptions []ServiceCustomFieldOptionModel
		resp.Diagnostics.Append(plan.FieldOptions.ElementsAs(ctx, &fieldOptions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		pdFieldOptions := make([]pagerduty.ServiceCustomFieldOption, 0, len(fieldOptions))
		for _, option := range fieldOptions {
			pdOption := pagerduty.ServiceCustomFieldOption{
				Data: pagerduty.ServiceCustomFieldOptionData{
					DataType: pagerduty.ServiceCustomFieldDataType(plan.DataType.ValueString()),
					Value:    option.Value.ValueString(),
				},
			}
			if !option.ID.IsNull() && option.ID.ValueString() != "" {
				pdOption.ID = option.ID.ValueString()
			}
			pdFieldOptions = append(pdFieldOptions, pdOption)
		}
		field.FieldOptions = pdFieldOptions
	}

	updatedField, err := r.client.UpdateServiceCustomField(ctx, &field)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating PagerDuty Service Custom Field",
			fmt.Sprintf("Could not update service custom field %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Map response body to model
	plan.ID = types.StringValue(updatedField.ID)
	plan.Name = types.StringValue(updatedField.Name)
	plan.DisplayName = types.StringValue(updatedField.DisplayName)
	plan.DataType = types.StringValue(string(updatedField.DataType))
	plan.FieldType = types.StringValue(string(updatedField.FieldType))
	plan.Enabled = types.BoolValue(updatedField.Enabled)

	if updatedField.Description != "" {
		plan.Description = types.StringValue(updatedField.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if updatedField.DefaultValue != nil {
		buf, _ := json.Marshal(updatedField.DefaultValue)
		plan.DefaultValue = jsontypes.NewNormalizedValue(string(buf))
	} else {
		plan.DefaultValue = jsontypes.NewNormalizedNull()
	}

	if len(updatedField.FieldOptions) > 0 {
		fieldOptions := make([]ServiceCustomFieldOptionModel, 0, len(updatedField.FieldOptions))
		for _, option := range updatedField.FieldOptions {
			fieldOptions = append(fieldOptions, ServiceCustomFieldOptionModel{
				ID:       types.StringValue(option.ID),
				Value:    types.StringValue(option.Data.Value),
				DataType: types.StringValue(string(option.Data.DataType)),
			})
		}

		fieldOptionsList, diags := types.ListValueFrom(ctx, serviceCustomFieldObjectType, fieldOptions)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.FieldOptions = fieldOptionsList
	} else {
		plan.FieldOptions = types.ListNull(serviceCustomFieldObjectType)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceCustomFieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceCustomFieldResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServiceCustomField(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PagerDuty Service Custom Field",
			fmt.Sprintf("Could not delete service custom field %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *ServiceCustomFieldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

var serviceCustomFieldObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":        types.StringType,
		"value":     types.StringType,
		"data_type": types.StringType,
	},
}
