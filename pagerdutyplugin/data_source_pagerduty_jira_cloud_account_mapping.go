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

type dataSourceJiraCloudAccountMapping struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceJiraCloudAccountMapping)(nil)

func (*dataSourceJiraCloudAccountMapping) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_jira_cloud_account_mapping"
}

func (*dataSourceJiraCloudAccountMapping) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"base_url":  schema.StringAttribute{Computed: true},
			"subdomain": schema.StringAttribute{Required: true},
		},
	}
}

func (d *dataSourceJiraCloudAccountMapping) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceJiraCloudAccountMapping) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty jira cloud account mapping")

	var searchSubdomain types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("subdomain"), &searchSubdomain)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.JiraCloudAccountsMapping
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := d.client.ListJiraCloudAccountsMappings(ctx, pagerduty.ListJiraCloudAccountsMappingsOptions{})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		for _, m := range response.AccountsMappings {
			if m.PagerDutyAccount.Subdomain == searchSubdomain.ValueString() {
				found = &m
				break
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty jira cloud account mapping with subdomain %s", searchSubdomain),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any jira cloud account mapping with the subdomain %s", searchSubdomain),
			"",
		)
		return
	}

	model := dataSourceJiraCloudAccountMappingModel{
		ID:        types.StringValue(found.ID),
		BaseURL:   types.StringValue(found.JiraCloudAccount.BaseURL),
		Subdomain: types.StringValue(found.PagerDutyAccount.Subdomain),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceJiraCloudAccountMappingModel struct {
	ID        types.String `tfsdk:"id"`
	BaseURL   types.String `tfsdk:"base_url"`
	Subdomain types.String `tfsdk:"subdomain"`
}
