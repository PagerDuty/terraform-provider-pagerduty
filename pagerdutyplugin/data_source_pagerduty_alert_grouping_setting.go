package pagerduty

import (
	"context"
	"fmt"
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/apiutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceAlertGroupingSetting struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceAlertGroupingSetting)(nil)

func (*dataSourceAlertGroupingSetting) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_alert_grouping_setting"
}

func (*dataSourceAlertGroupingSetting) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{Computed: true},
			"type":        schema.StringAttribute{Computed: true},
			"services": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"config": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"timeout": schema.Int64Attribute{
						Computed: true,
					},
					"time_window": schema.Int64Attribute{
						Computed: true,
					},
					"aggregate": schema.StringAttribute{
						Optional: true,
					},
					"fields": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *dataSourceAlertGroupingSetting) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceAlertGroupingSetting) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty alert grouping setting")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cursorAfter string
	var found *pagerduty.AlertGroupingSetting
	err := apiutil.All(ctx, func(offset int) (bool, error) {
		resp, err := d.client.ListAlertGroupingSettings(ctx, pagerduty.ListAlertGroupingSettingsOptions{
			After: cursorAfter,
			Limit: 100,
		})
		if err != nil {
			return false, err
		}

		for _, alertGroupingSetting := range resp.AlertGroupingSettings {
			if alertGroupingSetting.Name == searchName.ValueString() {
				found = &alertGroupingSetting
				break
			}
		}

		cursorAfter = resp.After
		return resp.After != "", nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty alert grouping setting %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any alert grouping setting with the name: %s", searchName),
			"",
		)
		return
	}

	model := flattenAlertGroupingSetting(found)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
