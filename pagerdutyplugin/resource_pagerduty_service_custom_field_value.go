package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ServiceCustomFieldValueResource{}
	_ resource.ResourceWithConfigure   = &ServiceCustomFieldValueResource{}
	_ resource.ResourceWithImportState = &ServiceCustomFieldValueResource{}
)

func serviceCustomFieldValueResource() resource.Resource {
	return &ServiceCustomFieldValueResource{}
}

type ServiceCustomFieldValueResource struct {
	client *pagerduty.Client
}

type ServiceCustomFieldValueResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ServiceID    types.String `tfsdk:"service_id"`
	CustomFields types.List   `tfsdk:"custom_fields"`
}

type ServiceCustomFieldValueModel struct {
	ID    types.String         `tfsdk:"id"`
	Name  types.String         `tfsdk:"name"`
	Value jsontypes.Normalized `tfsdk:"value"`
}

func (r *ServiceCustomFieldValueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "pagerduty_service_custom_field_value"
}

func (r *ServiceCustomFieldValueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A resource to manage PagerDuty Service Custom Field Values. Note: This is an Early Access feature that requires the X-EARLY-ACCESS header with service-custom-fields-preview value for all API operations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Description: "The ID of the service to set custom field values for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"custom_fields": schema.ListAttribute{
				Description: "The custom field values to set for the service.",
				Required:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":    types.StringType,
						"name":  types.StringType,
						"value": jsontypes.NormalizedType{},
					},
				},
			},
		},
	}
}

func (r *ServiceCustomFieldValueResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServiceCustomFieldValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceCustomFieldValueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := plan.ServiceID.ValueString()

	// Extract custom fields from plan
	var customFieldModels []ServiceCustomFieldValueModel
	resp.Diagnostics.Append(plan.CustomFields.ElementsAs(ctx, &customFieldModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare custom fields for API request
	pdCustomFields := make([]pagerduty.ServiceCustomFieldValue, 0, len(customFieldModels))
	for _, field := range customFieldModels {
		customField := pagerduty.ServiceCustomFieldValue{}

		// For API requests, we prioritize using name to identify fields
		// This helps avoid inconsistencies with IDs during state updates
		if !field.Name.IsNull() && field.Name.ValueString() != "" {
			customField.Name = field.Name.ValueString()
		} else if !field.ID.IsNull() && field.ID.ValueString() != "" {
			customField.ID = field.ID.ValueString()
		} else {
			resp.Diagnostics.AddError(
				"Missing Field Identifier",
				"Either name or id must be provided for each custom field",
			)
			return
		}

		d := field.Value.Unmarshal(&customField.Value)
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}

		pdCustomFields = append(pdCustomFields, customField)
	}

	// Update custom field values
	result, err := r.client.UpdateServiceCustomFieldValues(ctx, serviceID, &pagerduty.ListServiceCustomFieldValuesResponse{
		CustomFields: pdCustomFields,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting PagerDuty Service Custom Field Values",
			fmt.Sprintf("Could not set custom field values for service %s: %s", serviceID, err),
		)
		return
	}

	// Set ID to service_id since this resource doesn't have its own ID
	plan.ID = types.StringValue(serviceID)

	// Map response back to model
	if err := r.mapResponseToModel(ctx, result, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Mapping Response",
			fmt.Sprintf("Could not map response: %s", err),
		)
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceCustomFieldValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceCustomFieldValueResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := state.ServiceID.ValueString()

	// Get current custom field values
	result, err := r.client.GetServiceCustomFieldValues(ctx, serviceID)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading PagerDuty Service Custom Field Values",
			fmt.Sprintf("Could not read custom field values for service %s: %s", serviceID, err),
		)
		return
	}

	// Map response to model
	if err := r.mapResponseToModel(ctx, result, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error Mapping Response",
			fmt.Sprintf("Could not map response: %s", err),
		)
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceCustomFieldValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceCustomFieldValueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := plan.ServiceID.ValueString()

	// Extract custom fields from plan
	var customFieldModels []ServiceCustomFieldValueModel
	resp.Diagnostics.Append(plan.CustomFields.ElementsAs(ctx, &customFieldModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare custom fields for API request
	pdCustomFields := make([]pagerduty.ServiceCustomFieldValue, 0, len(customFieldModels))
	for _, field := range customFieldModels {
		customField := pagerduty.ServiceCustomFieldValue{}

		// For API requests, we prioritize using name to identify fields
		// This helps avoid inconsistencies with IDs during state updates
		if !field.Name.IsNull() && field.Name.ValueString() != "" {
			customField.Name = field.Name.ValueString()
		} else if !field.ID.IsNull() && field.ID.ValueString() != "" {
			customField.ID = field.ID.ValueString()
		} else {
			resp.Diagnostics.AddError(
				"Missing Field Identifier",
				"Either name or id must be provided for each custom field",
			)
			return
		}

		d := field.Value.Unmarshal(&customField.Value)
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}

		pdCustomFields = append(pdCustomFields, customField)
	}

	// Update custom field values
	result, err := r.client.UpdateServiceCustomFieldValues(ctx, serviceID, &pagerduty.ListServiceCustomFieldValuesResponse{
		CustomFields: pdCustomFields,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating PagerDuty Service Custom Field Values",
			fmt.Sprintf("Could not update custom field values for service %s: %s", serviceID, err),
		)
		return
	}

	// Map response to model
	if err := r.mapResponseToModel(ctx, result, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Mapping Response",
			fmt.Sprintf("Could not map response: %s", err),
		)
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceCustomFieldValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This is a no-op as there's no explicit delete operation for custom field values
	// They are tied to the service and will be deleted when the service is deleted
}

func (r *ServiceCustomFieldValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is the service ID
	resource.ImportStatePassthroughID(ctx, path.Root("service_id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	// For import, we'll use the Read method to populate the state
	// This ensures consistency with how the resource normally behaves
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State: resp.State,
	}

	r.Read(ctx, readReq, &readResp)

	// Copy any diagnostics from the read response
	resp.Diagnostics.Append(readResp.Diagnostics...)

	// Copy the state from the read response
	resp.State = readResp.State
}

func (r *ServiceCustomFieldValueResource) mapResponseToModel(ctx context.Context, result *pagerduty.ListServiceCustomFieldValuesResponse, model *ServiceCustomFieldValueResourceModel) error {
	// First, extract the current custom fields from the model to maintain the same order and selection
	var currentFields []ServiceCustomFieldValueModel
	if !model.CustomFields.IsNull() {
		diags := model.CustomFields.ElementsAs(ctx, &currentFields, false)
		if diags.HasError() {
			return fmt.Errorf("error extracting current custom fields: %v", diags)
		}
	}

	// Create a map of field name/id to field value for quick lookup
	fieldMap := make(map[string]pagerduty.ServiceCustomFieldValue)
	for _, field := range result.CustomFields {
		// Use both ID and name as keys for lookup
		if field.ID != "" {
			fieldMap[field.ID] = field
		}
		if field.Name != "" {
			fieldMap[field.Name] = field
		}
	}

	// Update the current fields with values from the API response
	updatedFields := make([]ServiceCustomFieldValueModel, 0, len(currentFields))
	for _, field := range currentFields {
		updatedField := field

		// Look up the field in the API response by ID or name
		var apiField pagerduty.ServiceCustomFieldValue
		var found bool

		// Prioritize lookup by name to maintain consistency
		if !field.Name.IsNull() && field.Name.ValueString() != "" {
			apiField, found = fieldMap[field.Name.ValueString()]
		}

		if !found && !field.ID.IsNull() && field.ID.ValueString() != "" {
			apiField, found = fieldMap[field.ID.ValueString()]
		}

		// If we found the field in the API response, update the value
		if found {
			// Only set the ID if it was already set in the plan
			// This prevents Terraform from seeing an inconsistency when the ID wasn't specified
			if !field.ID.IsNull() {
				updatedField.ID = types.StringValue(apiField.ID)
			}

			// Format the value based on its type

			buf, _ := json.Marshal(apiField.Value)
			updatedField.Value = jsontypes.NewNormalizedValue(string(buf))
		}

		updatedFields = append(updatedFields, updatedField)
	}

	// Create the list value from the updated fields
	fieldsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"name":  types.StringType,
			"value": jsontypes.NormalizedType{},
		},
	}, updatedFields)

	if diags.HasError() {
		return fmt.Errorf("error creating custom fields list: %v", diags)
	}

	model.CustomFields = fieldsList
	return nil
}
