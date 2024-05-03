package pagerduty

import (
	"context"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceStandardsResourcesScores struct {
	client *pagerduty.Client
}

var _ datasource.DataSource = (*dataSourceStandardsResourcesScores)(nil)

func (d *dataSourceStandardsResourcesScores) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_standards_resources_scores"
}

func (d *dataSourceStandardsResourcesScores) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"resource_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("technical_services"),
				},
			},
			"resources": schema.ListAttribute{
				ElementType: resourceStandardScoreObjectType,
				Computed:    true,
			},
		},
	}
}

func (d *dataSourceStandardsResourcesScores) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceStandardsResourcesScores) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceStandardsResourcesScoresModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	rt := data.ResourceType.ValueString()
	ids := make([]string, 0)
	resp.Diagnostics.Append(data.IDs.ElementsAs(ctx, &ids, true)...)

	opt := pagerduty.ListMultiResourcesStandardScoresOptions{IDs: ids}
	scores, err := d.client.ListMultiResourcesStandardScores(ctx, rt, opt)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(
			"Error calling ListResourceStandardScores",
			err.Error(),
		))
		return
	}

	resources, di := resourceStandardScoresToModel(scores.Resources)
	resp.Diagnostics.Append(di...)
	data.Resources = resources

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type dataSourceStandardsResourcesScoresModel struct {
	IDs          types.List   `tfsdk:"ids"`
	ResourceType types.String `tfsdk:"resource_type"`
	Resources    types.List   `tfsdk:"resources"`
}

func resourceStandardScoresToModel(data []pagerduty.ResourceStandardScore) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var list []attr.Value
	for i := range data {
		resourceStandardScore, di := resourceStandardScoreToModel(&data[i])
		diags.Append(di...)
		list = append(list, resourceStandardScore)
	}
	listValue, di := types.ListValue(resourceStandardScoreObjectType, list)
	diags.Append(di...)
	return listValue, diags
}

var resourceStandardScoreObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"resource_id":   types.StringType,
		"resource_type": types.StringType,
		"standards":     types.ListType{ElemType: resourceStandardObjectType},
		"score":         resourceScoresObjectType,
	},
}

func resourceStandardScoreToModel(data *pagerduty.ResourceStandardScore) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	standards, di := resourceStandardsToModel(data.Standards)
	diags.Append(di...)

	score, di := resourceScoreToModel(data.Score)
	diags.Append(di...)

	standardScore, di := types.ObjectValue(resourceStandardScoreObjectType.AttrTypes, map[string]attr.Value{
		"resource_id":   types.StringValue(data.ResourceID),
		"resource_type": types.StringValue(data.ResourceType),
		"score":         score,
		"standards":     standards,
	})
	diags.Append(di...)

	return standardScore, diags
}
