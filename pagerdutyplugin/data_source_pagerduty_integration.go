package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/apiutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceIntegration struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceIntegration)(nil)

func (*dataSourceIntegration) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_service_integration"
}

func (*dataSourceIntegration) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true},
			"service_name":    schema.StringAttribute{Required: true},
			"integration_key": schema.StringAttribute{Computed: true, Sensitive: true},
			"integration_summary": schema.StringAttribute{
				Required:    true,
				Description: `examples "Amazon CloudWatch", "New Relic"`,
			},
		},
	}
}

func (d *dataSourceIntegration) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceIntegration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty service integration")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service_name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.Service
	err := apiutil.All(ctx, func(offset int) (bool, error) {
		list, err := d.client.ListServicesWithContext(ctx, pagerduty.ListServiceOptions{
			Query:  searchName.ValueString(),
			Limit:  apiutil.Limit,
			Offset: uint(offset),
		})
		if err != nil {
			return false, err
		}

		for _, service := range list.Services {
			if service.Name == searchName.ValueString() {
				found = &service
				return false, nil
			}
		}
		return list.More, nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty service integration %s", searchName),
			err.Error(),
		)
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any service with the name: %s", searchName),
			"",
		)
		return
	}

	var summary types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("integration_summary"), &summary)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundIntegration *pagerduty.Integration
	for _, integration := range found.Integrations {
		if strings.EqualFold(integration.Summary, summary.ValueString()) {
			foundIntegration = &integration
			break
		}
	}

	if foundIntegration == nil {
		resp.Diagnostics.Append(dataSourceIntegrationNotFoundError(nil, searchName, summary))
		return
	}

	var model dataSourceIntegrationModel
	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		details, err := d.client.GetIntegrationWithContext(ctx, found.ID, foundIntegration.ID, pagerduty.GetIntegrationOptions{})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		model = dataSourceIntegrationModel{
			ID:             types.StringValue(foundIntegration.ID),
			ServiceName:    types.StringValue(found.Name),
			IntegrationKey: types.StringValue(details.IntegrationKey),
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.Append(dataSourceIntegrationNotFoundError(err, searchName, summary))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceIntegrationModel struct {
	ID                 types.String `tfsdk:"id"`
	ServiceName        types.String `tfsdk:"service_name"`
	IntegrationKey     types.String `tfsdk:"integration_key"`
	IntegrationSummary types.String `tfsdk:"integration_summary"`
}

func dataSourceIntegrationNotFoundError(err error, service, summary types.String) diag.Diagnostic {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return diag.NewErrorDiagnostic(
		fmt.Sprintf("Unable to locate any integration of type %s on service %s", summary, service),
		errMsg,
	)
}
