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

type dataSourceAutomationActionsRunner struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceAutomationActionsRunner)(nil)

func (*dataSourceAutomationActionsRunner) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_automation_actions_runner"
}

func (*dataSourceAutomationActionsRunner) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Required: true},
			"name":             schema.StringAttribute{Computed: true},
			"type":             schema.StringAttribute{Computed: true},
			"runner_type":      schema.StringAttribute{Computed: true},
			"creation_time":    schema.StringAttribute{Computed: true},
			"last_seen":        schema.StringAttribute{Computed: true, Optional: true},
			"description":      schema.StringAttribute{Computed: true, Optional: true},
			"runbook_base_uri": schema.StringAttribute{Computed: true, Optional: true},
		},
	}
}

func (d *dataSourceAutomationActionsRunner) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceAutomationActionsRunner) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty automation actions runner")

	var searchID types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &searchID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var model dataSourceAutomationActionsRunnerModel
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		runner, err := d.client.GetAutomationActionsRunnerWithContext(ctx, searchID.ValueString())
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		model = dataSourceAutomationActionsRunnerModel{
			ID:           types.StringValue(runner.ID),
			Name:         types.StringValue(runner.Name),
			Type:         types.StringValue(runner.Type),
			RunnerType:   types.StringValue(runner.RunnerType),
			CreationTime: types.StringValue(runner.CreationTime),
		}

		if runner.Description != "" {
			model.Description = types.StringValue(runner.Description)
		}

		if runner.RunbookBaseURI != "" {
			model.RunbookBaseURI = types.StringValue(runner.RunbookBaseURI)
		}

		if runner.LastSeen != "" {
			model.LastSeen = types.StringValue(runner.LastSeen)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty automation actions runner %s", searchID),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceAutomationActionsRunnerModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	RunnerType     types.String `tfsdk:"runner_type"`
	CreationTime   types.String `tfsdk:"creation_time"`
	LastSeen       types.String `tfsdk:"last_seen"`
	Description    types.String `tfsdk:"description"`
	RunbookBaseURI types.String `tfsdk:"runbook_base_uri"`
}
