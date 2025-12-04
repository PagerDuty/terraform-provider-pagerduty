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

type dataSourceServiceCustomField struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceServiceCustomField)(nil)

func (*dataSourceServiceCustomField) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_service_custom_field"
}

func (*dataSourceServiceCustomField) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The human-readable name of the field. This must be unique across an account",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the field. May include ASCII characters, specifically lowercase letters, digits, and underescores. The name for a Field must be unique and cannot be changed once created",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "API object type",
			},
			"summary": schema.StringAttribute{
				Computed:    true,
				Description: "A short-form, server-generated string that provides succinct, important information about an object suitable for primary labeling of an entity in a client. In many cases, this will be identical to display_name",
			},
			"self": schema.StringAttribute{
				Computed:    true,
				Description: "The API show URL at which the object is accessible",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A description of the data this field contains",
			},
			"data_type": schema.StringAttribute{
				Computed:    true,
				Description: "The kind of data the custom field is allowed to contain",
			},
			"field_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of data this field contains. In combination with the data_type field",
			},
			"default_value": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the field is enabled	",
			},
			"field_options": schema.ListAttribute{
				Computed:    true,
				ElementType: serviceCustomFieldObjectType,
				Description: "The options for the custom field. Applies only to single_value_fixed and multi_value_fixed field types",
			},
		},
	}
}

func (d *dataSourceServiceCustomField) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceServiceCustomField) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty service custom field")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("display_name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.ServiceCustomField
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		offset := uint(0)
		more := true

		for more {
			response, err := d.client.ListServiceCustomFields(ctx, pagerduty.ListServiceCustomFieldsOptions{
				Limit:  100,
				Offset: offset,
			})
			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}

			for _, field := range response.Fields {
				if field.DisplayName == searchName.ValueString() {
					found = &field
					return nil
				}
			}

			more = response.More
			offset += response.Limit
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty service custom field %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any service custom field with the display name: %s", searchName),
			"",
		)
		return
	}

	fieldOptions := types.ListNull(serviceCustomFieldObjectType)
	if len(found.FieldOptions) > 0 {
		elements := make([]attr.Value, len(found.FieldOptions))
		fieldOptions = types.ListValueMust(serviceCustomFieldObjectType, elements)
	}

	defaultValue := jsontypes.NewNormalizedNull()
	if found.DefaultValue != nil {
		buf, _ := json.Marshal(found.DefaultValue)
		defaultValue = jsontypes.NewNormalizedValue(string(buf))
	}

	model := dataSourceServiceCustomFieldModel{
		ID:           types.StringValue(found.ID),
		Name:         types.StringValue(found.Name),
		Type:         types.StringValue(found.Type),
		Summary:      types.StringValue(found.Summary),
		Self:         types.StringValue(found.Self),
		DisplayName:  types.StringValue(found.DisplayName),
		Description:  types.StringValue(found.Description),
		DataType:     types.StringValue(string(found.DataType)),
		FieldType:    types.StringValue(string(found.FieldType)),
		Enabled:      types.BoolValue(found.Enabled),
		DefaultValue: defaultValue,
		FieldOptions: fieldOptions,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceServiceCustomFieldModel struct {
	ID           types.String         `tfsdk:"id"`
	Name         types.String         `tfsdk:"name"`
	Type         types.String         `tfsdk:"type"`
	Summary      types.String         `tfsdk:"summary"`
	Self         types.String         `tfsdk:"self"`
	DisplayName  types.String         `tfsdk:"display_name"`
	Description  types.String         `tfsdk:"description"`
	DataType     types.String         `tfsdk:"data_type"`
	FieldType    types.String         `tfsdk:"field_type"`
	DefaultValue jsontypes.Normalized `tfsdk:"default_value"`
	Enabled      types.Bool           `tfsdk:"enabled"`
	FieldOptions types.List           `tfsdk:"field_options"`
}
