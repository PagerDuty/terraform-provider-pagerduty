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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceServiceCustomFieldValue struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceServiceCustomFieldValue)(nil)

func (*dataSourceServiceCustomFieldValue) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_service_custom_field_value"
}

func (*dataSourceServiceCustomFieldValue) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the custom field values for a PagerDuty service. Note: This is an Early Access feature that requires the X-EARLY-ACCESS header with service-custom-fields-preview value for all API operations.",
		Attributes: map[string]schema.Attribute{
			"service_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the service to get custom field values for.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource.",
			},
			"custom_fields": schema.ListAttribute{
				Computed:    true,
				Description: "The custom field values for the service.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":           types.StringType,
						"name":         types.StringType,
						"display_name": types.StringType,
						"description":  types.StringType,
						"data_type":    types.StringType,
						"field_type":   types.StringType,
						"type":         types.StringType,
						"value":        jsontypes.NormalizedType{},
					},
				},
			},
		},
	}
}

func (d *dataSourceServiceCustomFieldValue) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceServiceCustomFieldValue) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty service custom field values")

	var serviceID types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service_id"), &serviceID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result *pagerduty.ListServiceCustomFieldValuesResponse
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		var err error
		result, err = d.client.GetServiceCustomFieldValues(ctx, serviceID.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty service custom field values for service %s", serviceID.ValueString()),
			err.Error(),
		)
		return
	}

	// Create custom field objects
	customFieldsList := []dataSourceServiceCustomFieldValueItemModel{}
	for _, field := range result.CustomFields {
		// Initialize with null value
		value := jsontypes.NewNormalizedNull()

		// Only set the value if it's not nil
		if field.Value != nil {
			buf, err := json.Marshal(field.Value)
			if err == nil {
				value = jsontypes.NewNormalizedValue(string(buf))
			}
		}

		customFieldsList = append(customFieldsList, dataSourceServiceCustomFieldValueItemModel{
			ID:          types.StringValue(field.ID),
			Name:        types.StringValue(field.Name),
			DisplayName: types.StringValue(field.DisplayName),
			Description: types.StringValue(field.Description),
			DataType:    types.StringValue(string(field.DataType)),
			FieldType:   types.StringValue(string(field.FieldType)),
			Type:        types.StringValue(field.Type),
			Value:       value,
		})
	}

	// Create the list value from the custom fields
	customFields, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           types.StringType,
			"name":         types.StringType,
			"display_name": types.StringType,
			"description":  types.StringType,
			"data_type":    types.StringType,
			"field_type":   types.StringType,
			"type":         types.StringType,
			"value":        jsontypes.NormalizedType{},
		},
	}, customFieldsList)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	model := dataSourceServiceCustomFieldValueModel{
		ID:           types.StringValue(serviceID.ValueString()),
		ServiceID:    serviceID,
		CustomFields: customFields,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceServiceCustomFieldValueModel struct {
	ID           types.String `tfsdk:"id"`
	ServiceID    types.String `tfsdk:"service_id"`
	CustomFields types.List   `tfsdk:"custom_fields"`
}

type dataSourceServiceCustomFieldValueItemModel struct {
	ID          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	DisplayName types.String         `tfsdk:"display_name"`
	Description types.String         `tfsdk:"description"`
	DataType    types.String         `tfsdk:"data_type"`
	FieldType   types.String         `tfsdk:"field_type"`
	Type        types.String         `tfsdk:"type"`
	Value       jsontypes.Normalized `tfsdk:"value"`
}
