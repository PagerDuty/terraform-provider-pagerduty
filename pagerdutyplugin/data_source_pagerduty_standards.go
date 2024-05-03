package pagerduty

import (
	"context"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceStandards struct {
	client *pagerduty.Client
}

var _ datasource.DataSourceWithConfigure = (*dataSourceStandards)(nil)

func (d *dataSourceStandards) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_standards"
}

func (d *dataSourceStandards) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_type": schema.StringAttribute{Optional: true},
			"standards": schema.ListAttribute{
				ElementType: standardObjectType,
				Computed:    true,
			},
		},
	}
}

func (d *dataSourceStandards) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceStandards) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceStandardsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	opts := pagerduty.ListStandardsOptions{}
	if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
		opts.ResourceType = data.ResourceType.ValueString()
	}

	list, err := d.client.ListStandards(ctx, opts)
	if err != nil {
		resp.Diagnostics.AddError("Error calling ListStandards", err.Error())
		return
	}

	standards, diags := flattenStandards(ctx, list.Standards)
	resp.Diagnostics.Append(diags...)
	data.Standards = standards
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenStandards(ctx context.Context, list []pagerduty.Standard) (types.List, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	mapList := make([]types.Object, 0, len(list))
	for _, standard := range list {
		exclusions := make([]types.Object, 0, len(standard.Exclusions))
		for _, exc := range standard.Exclusions {
			item, diags := types.ObjectValue(
				standardReferenceObjectType.AttrTypes,
				map[string]attr.Value{
					"id":   types.StringValue(exc.ID),
					"type": types.StringValue(exc.Type),
				},
			)
			diagnostics.Append(diags...)
			exclusions = append(exclusions, item)
		}
		exclusionsValue, diags := types.ListValueFrom(ctx, standardReferenceObjectType, exclusions)
		diagnostics.Append(diags...)

		inclusions := make([]types.Object, 0, len(standard.Exclusions))
		for _, inc := range standard.Exclusions {
			item, diags := types.ObjectValue(
				standardReferenceObjectType.AttrTypes,
				map[string]attr.Value{
					"id":   types.StringValue(inc.ID),
					"type": types.StringValue(inc.Type),
				},
			)
			diagnostics.Append(diags...)
			inclusions = append(inclusions, item)
		}
		inclusionsValue, diags := types.ListValueFrom(ctx, standardReferenceObjectType, inclusions)
		diagnostics.Append(diags...)

		item, diags := types.ObjectValue(
			standardObjectType.AttrTypes,
			map[string]attr.Value{
				"active":        types.BoolValue(standard.Active),
				"description":   types.StringValue(standard.Description),
				"id":            types.StringValue(standard.ID),
				"name":          types.StringValue(standard.Name),
				"type":          types.StringValue(standard.Type),
				"resource_type": types.StringValue(standard.ResourceType),
				"exclusions":    exclusionsValue,
				"inclusions":    inclusionsValue,
			},
		)
		diagnostics.Append(diags...)
		mapList = append(mapList, item)
	}
	listValue, diags := types.ListValueFrom(ctx, standardObjectType, mapList)
	diagnostics.Append(diags...)
	return listValue, diagnostics
}

type dataSourceStandardsModel struct {
	ResourceType types.String `tfsdk:"resource_type"`
	Standards    types.List   `tfsdk:"standards"`
}

var standardReferenceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	},
}

var standardObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"active":        types.BoolType,
		"description":   types.StringType,
		"id":            types.StringType,
		"name":          types.StringType,
		"type":          types.StringType,
		"resource_type": types.StringType,
		"exclusions": types.ListType{
			ElemType: standardReferenceObjectType,
		},
		"inclusions": types.ListType{
			ElemType: standardReferenceObjectType,
		},
	},
}
