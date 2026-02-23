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

type dataSourceScheduleV2 struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceScheduleV2)(nil)

func (*dataSourceScheduleV2) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_schedulev2"
}

func (*dataSourceScheduleV2) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a specific PagerDuty v3 schedule so that you can reference it in other resources.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, Description: "The ID of the schedule."},
			"name": schema.StringAttribute{Required: true, Description: "The name of the schedule to search for."},
		},
	}
}

func (d *dataSourceScheduleV2) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceScheduleV2) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty v3 schedule")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := pagerduty.ListSchedulesV3Options{Query: searchName.ValueString()}

	var found *pagerduty.APIObject
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := d.client.ListSchedulesV3(ctx, opts)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		for i, schedule := range response.Schedules {
			if schedule.Summary == searchName.ValueString() {
				found = &response.Schedules[i]
				break
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty v3 schedule %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any v3 schedule with the name: %s", searchName),
			"",
		)
		return
	}

	model := dataSourceScheduleV2Model{
		ID:   types.StringValue(found.ID),
		Name: types.StringValue(found.Summary),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceScheduleV2Model struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
