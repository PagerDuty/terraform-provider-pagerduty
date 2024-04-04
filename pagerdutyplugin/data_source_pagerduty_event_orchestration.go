package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceEventOrchestration struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceEventOrchestration)(nil)

func (*dataSourceEventOrchestration) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_event_orchestration"
}

func (*dataSourceEventOrchestration) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
			"integration": schema.ListAttribute{
				Computed:    true,
				ElementType: eventOrchestrationIntegrationObjectType,
			},
		},
	}
}

func (d *dataSourceEventOrchestration) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceEventOrchestration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty event orchestration")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.Orchestration
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := d.client.ListOrchestrationsWithContext(ctx, pagerduty.ListOrchestrationsOptions{})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		for _, orchestration := range response.Orchestrations {
			if orchestration.Name == searchName.ValueString() {
				found = &orchestration
				break
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty event orchestration %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any event orchestration with the name: %s", searchName),
			"",
		)
		return
	}

	// TODO: Get

	model := dataSourceEventOrchestrationModel{
		ID:          types.StringValue(found.ID),
		Name:        types.StringValue(found.Name),
		Integration: flattenEventOrchestrationIntegrations(found.Integrations),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceEventOrchestrationModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Integration types.List   `tfsdk:"integration"`
}

var eventOrchestrationIntegrationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"label":      types.StringType,
		"parameters": types.ListType{ElemType: eventOrchestrationParameterObjectType},
	},
}

var eventOrchestrationParameterObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"routing_key": types.StringType,
		"type":        types.StringType,
	},
}

func flattenEventOrchestrationIntegrations(list []*pagerduty.OrchestrationIntegration) types.List {
	elements := make([]attr.Value, 0, len(list))
	for _, integration := range list {
		obj := types.ObjectValueMust(eventOrchestrationIntegrationObjectType.AttrTypes, map[string]attr.Value{
			"id":         types.StringValue(integration.ID),
			"label":      types.StringNull(),
			"parameters": flattenEventOrchestrationIntegrationParameters(integration.Parameters),
		})
		elements = append(elements, obj)
	}
	return types.ListValueMust(eventOrchestrationIntegrationObjectType, elements)
}

func flattenEventOrchestrationIntegrationParameters(p *pagerduty.OrchestrationIntegrationParameters) types.List {
	obj := types.ObjectValueMust(eventOrchestrationParameterObjectType.AttrTypes, map[string]attr.Value{
		"routing_key": types.StringValue(p.RoutingKey),
		"type":        types.StringValue(p.Type),
	})
	return types.ListValueMust(eventOrchestrationParameterObjectType, []attr.Value{obj})
}
