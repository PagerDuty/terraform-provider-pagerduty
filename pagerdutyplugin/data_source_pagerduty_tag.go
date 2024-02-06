package pagerduty

import (
	"context"
	"fmt"
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceTag struct {
	client *pagerduty.Client
}

var _ datasource.DataSourceWithConfigure = (*dataSourceStandards)(nil)

func (d *dataSourceTag) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceTag) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_tag"
}

func (d *dataSourceTag) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				Required:    true,
				Description: "The label of the tag to find in the PagerDuty API",
			},
			"id": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *dataSourceTag) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var searchTag string
	if d := req.Config.GetAttribute(ctx, path.Root("label"), &searchTag); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	log.Printf("[INFO] Reading PagerDuty tag")

	// TODO: retry
	list, err := d.client.ListTags(pagerduty.ListTagOptions{Query: searchTag})
	if err != nil {
		// TODO: if 400 non retryable
		resp.Diagnostics.AddError("Error calling ListTags", err.Error())
		// TODO: wait 30 + retry
		return
	}

	var found *pagerduty.Tag

	for _, tag := range list.Tags {
		if tag.Label == searchTag {
			found = tag
			break
		}
	}

	if found == nil {
		// return retry.NonRetryableError(
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any tag with label: %s", searchTag),
			"",
		)
		return
	}

	model := dataSourceTagModel{
		ID:    types.StringValue(found.ID),
		Label: types.StringValue(found.Label),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceTagModel struct {
	ID    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}
