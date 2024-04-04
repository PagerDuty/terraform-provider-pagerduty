package pagerduty

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/apiutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceEventOrchestrations struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceEventOrchestrations)(nil)

func (*dataSourceEventOrchestrations) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_event_orchestrations"
}

func (*dataSourceEventOrchestrations) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name_filter": schema.StringAttribute{Required: true},
			"event_orchestrations": schema.ListAttribute{
				Computed:    true,
				ElementType: eventOrchestrationObjectType,
			},
		},
	}
}

func (d *dataSourceEventOrchestrations) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceEventOrchestrations) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty event orchestrations")

	var nameFilter types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name_filter"), &nameFilter)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameFilterRe, err := regexp.Compile(nameFilter.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("invalid regexp for name_filter provided %s", nameFilter),
			err.Error(),
		)
		return
	}

	var found []pagerduty.Orchestration
	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := d.client.ListOrchestrationsWithContext(ctx, pagerduty.ListOrchestrationsOptions{
			Limit: apiutil.Limit,
		})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		for _, orchestration := range response.Orchestrations {
			if nameFilterRe.MatchString(orchestration.Name) {
				found = append(found, orchestration)
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty event orchestrations matching %s", nameFilter),
			err.Error(),
		)
		return
	}

	if len(found) == 0 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any Event Orchestration matching the expression: %s", nameFilter),
			"",
		)
		return
	}

	// TODO: get

	model := dataSourceEventOrchestrationsModel{
		NameFilter:          nameFilter,
		EventOrchestrations: flattenEventOrchestrations(found),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceEventOrchestrationsModel struct {
	NameFilter          types.String `tfsdk:"name_filter"`
	EventOrchestrations types.List   `tfsdk:"event_orchestrations"`
}

var eventOrchestrationObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"integration": types.ListType{ElemType: eventOrchestrationIntegrationObjectType},
	},
}

func flattenEventOrchestrations(list []pagerduty.Orchestration) types.List {
	elements := []attr.Value{}
	for _, o := range list {
		obj := types.ObjectValueMust(
			eventOrchestrationObjectType.AttrTypes,
			map[string]attr.Value{
				"id":          types.StringValue(o.ID),
				"name":        types.StringValue(o.Name),
				"integration": flattenEventOrchestrationIntegrations(o.Integrations),
			},
		)
		elements = append(elements, obj)
	}
	return types.ListValueMust(eventOrchestrationObjectType, elements)
}
