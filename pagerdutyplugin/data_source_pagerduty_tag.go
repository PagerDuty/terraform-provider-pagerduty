package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceTag struct {
	client *pagerduty.Client
}

var _ datasource.DataSourceWithConfigure = (*dataSourceStandards)(nil)

func (d *dataSourceTag) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_tag"
}

func (d *dataSourceTag) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *dataSourceTag) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceTag) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var searchTag string
	if d := req.Config.GetAttribute(ctx, path.Root("label"), &searchTag); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	log.Printf("[INFO] Reading PagerDuty tag")

	var tags []*pagerduty.Tag
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		list, err := d.client.ListTagsPaginated(ctx, pagerduty.ListTagOptions{Query: searchTag, Limit: 100})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		tags = list
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading list of tags", err.Error())
	}

	var found *pagerduty.Tag
	for _, tag := range tags {
		if tag.Label == searchTag {
			found = tag
			break
		}
	}
	if found == nil {
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
