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

type dataSourceIncidentTypeCustomField struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceIncidentTypeCustomField)(nil)

func (*dataSourceIncidentTypeCustomField) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_incident_type_custom_field"
}

func (*dataSourceIncidentTypeCustomField) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"incident_type": schema.StringAttribute{Required: true},
			"display_name":  schema.StringAttribute{Required: true},
			"data_type":     schema.StringAttribute{Computed: true},
			"default_value": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
			},
			"description": schema.StringAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"field_options": schema.ListAttribute{
				Computed:    true,
				ElementType: incidentTypeFieldOptionObjectType,
			},
			"field_type": schema.StringAttribute{Computed: true},
			"name":       schema.StringAttribute{Computed: true},
			"self":       schema.StringAttribute{Computed: true},
			"summary":    schema.StringAttribute{Computed: true},
			"type":       schema.StringAttribute{Computed: true},
		},
	}
}

func (d *dataSourceIncidentTypeCustomField) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceIncidentTypeCustomField) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var searchName types.String
	diags := req.Config.GetAttribute(ctx, path.Root("display_name"), &searchName)
	if resp.Diagnostics.Append(diags...); diags.HasError() {
		return
	}

	var searchIncidentType types.String
	diags = req.Config.GetAttribute(ctx, path.Root("incident_type"), &searchIncidentType)
	if resp.Diagnostics.Append(diags...); diags.HasError() {
		return
	}
	incidentTypeID := searchIncidentType.ValueString()

	log.Printf("[INFO] Reading PagerDuty incident type custom field %s %s", searchIncidentType, searchName)

	var found *pagerduty.IncidentTypeField
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := d.client.ListIncidentTypeFields(ctx, incidentTypeID, pagerduty.ListIncidentTypeFieldsOptions{
			Includes: []string{"field_options"},
		})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		for _, f := range response.Fields {
			if f.DisplayName == searchName.ValueString() {
				found = &f
				break
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type custom field %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any incident type custom field with the name: %s", searchName),
			"",
		)
		return
	}

	defaultValue, _ := json.Marshal(found.DefaultValue)

	elements := make([]attr.Value, 0, len(found.FieldOptions))
	for _, opt := range found.FieldOptions {
		dataObj := types.ObjectNull(incidentTypeFieldOptionDataObjectType.AttrTypes)
		if opt.Data != nil {
			dataObj = types.ObjectValueMust(incidentTypeFieldOptionDataObjectType.AttrTypes, map[string]attr.Value{
				"value":     types.StringValue(opt.Data.Value),
				"data_type": types.StringValue(opt.Data.DataType),
			})
		}
		obj := types.ObjectValueMust(incidentTypeFieldOptionObjectType.AttrTypes, map[string]attr.Value{
			"id":   types.StringValue(opt.ID),
			"type": types.StringValue(opt.Type),
			"data": dataObj,
		})
		elements = append(elements, obj)
	}
	fieldOptions := types.ListValueMust(incidentTypeFieldOptionObjectType, elements)

	model := dataSourceIncidentTypeCustomFieldModel{
		ID:           types.StringValue(found.ID),
		IncidentType: types.StringValue(found.IncidentType),
		DisplayName:  types.StringValue(found.DisplayName),
		DataType:     types.StringValue(found.DataType),
		DefaultValue: jsontypes.NewNormalizedValue(string(defaultValue)),
		Description:  types.StringValue(found.Description),
		Enabled:      types.BoolValue(found.Enabled),
		FieldOptions: fieldOptions,
		FieldType:    types.StringValue(found.FieldType),
		Name:         types.StringValue(found.Name),
		Self:         types.StringValue(found.Self),
		Summary:      types.StringValue(found.Summary),
		Type:         types.StringValue(found.Type),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceIncidentTypeCustomFieldModel struct {
	ID           types.String         `tfsdk:"id"`
	IncidentType types.String         `tfsdk:"incident_type"`
	DisplayName  types.String         `tfsdk:"display_name"`
	DataType     types.String         `tfsdk:"data_type"`
	DefaultValue jsontypes.Normalized `tfsdk:"default_value"`
	Description  types.String         `tfsdk:"description"`
	Enabled      types.Bool           `tfsdk:"enabled"`
	FieldOptions types.List           `tfsdk:"field_options"`
	FieldType    types.String         `tfsdk:"field_type"`
	Name         types.String         `tfsdk:"name"`
	Self         types.String         `tfsdk:"self"`
	Summary      types.String         `tfsdk:"summary"`
	Type         types.String         `tfsdk:"type"`
}

var incidentTypeFieldOptionDataObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"value":     types.StringType,
		"data_type": types.StringType,
	},
}

var incidentTypeFieldOptionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
		"data": incidentTypeFieldOptionDataObjectType,
	},
}
