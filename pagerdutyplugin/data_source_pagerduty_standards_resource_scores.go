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

type dataSourceStandardsResourceScores struct {
	client *pagerduty.Client
}

var _ datasource.DataSource = (*dataSourceStandardsResourceScores)(nil)

func (d *dataSourceStandardsResourceScores) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_standards_resource_scores"
}

func (d *dataSourceStandardsResourceScores) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"resource_type": schema.StringAttribute{Required: true},
			"score": schema.ObjectAttribute{
				AttributeTypes: resourceScoresObjectType.AttrTypes,
				Computed:       true,
			},
			"standards": schema.ListAttribute{
				ElementType: resourceStandardObjectType,
				Computed:    true,
			},
		},
	}
}

func (d *dataSourceStandardsResourceScores) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ConfigurePagerdutyClient(&d.client, req, resp)
}

func (d *dataSourceStandardsResourceScores) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceStandardsResourceScoresModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	id := data.ID.ValueString()
	rt := data.ResourceType.ValueString()

	scores, err := d.client.ListResourceStandardScores(ctx, id, rt)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(
			"Error calling ListResourceStandardScores",
			err.Error(),
		))
		return
	}

	data.ID = types.StringValue(scores.ResourceID)
	data.ResourceType = types.StringValue(scores.ResourceType)

	standards, diags := resourceStandardsToModel(scores.Standards, &data)
	resp.Diagnostics.Append(diags...)
	data.Standards = standards

	score, diags := resourceScoreToModel(scores.Score)
	resp.Diagnostics.Append(diags...)
	data.Score = score

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type dataSourceStandardsResourceScoresModel struct {
	ID           types.String `tfsdk:"id"`
	ResourceType types.String `tfsdk:"resource_type"`
	Standards    types.List   `tfsdk:"standards"`
	Score        types.Object `tfsdk:"score"`
}

var resourceScoresObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"passing": types.Int64Type,
		"total":   types.Int64Type,
	},
}

func resourceScoreToModel(score *pagerduty.ResourceScore) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(resourceScoresObjectType.AttrTypes, map[string]attr.Value{
		"passing": types.Int64Value(int64(score.Passing)),
		"total":   types.Int64Value(int64(score.Total)),
	})
}

var resourceStandardObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"active":      types.BoolType,
		"description": types.StringType,
		"id":          types.StringType,
		"name":        types.StringType,
		"type":        types.StringType,
		"pass":        types.BoolType,
	},
}

func resourceStandardsToModel(standards []pagerduty.ResourceStandard, data *dataSourceStandardsResourceScoresModel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var list []attr.Value
	for _, standard := range standards {
		value, di := types.ObjectValue(resourceStandardObjectType.AttrTypes, map[string]attr.Value{
			"active":      types.BoolValue(standard.Active),
			"description": types.StringValue(standard.Description),
			"id":          types.StringValue(standard.ID),
			"name":        types.StringValue(standard.Name),
			"pass":        types.BoolValue(standard.Pass),
			"type":        types.StringValue(standard.Type),
		})
		diags.Append(di...)
		list = append(list, value)
	}
	modelStandards, di := types.ListValue(resourceStandardObjectType, list)
	diags.Append(di...)
	return modelStandards, diags
}
