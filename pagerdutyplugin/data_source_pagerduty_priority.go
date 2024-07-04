package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/apiutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourcePriority struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourcePriority)(nil)

func (*dataSourcePriority) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_priority"
}

func (*dataSourcePriority) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the priority to find in the PagerDuty API",
			},
			"description": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *dataSourcePriority) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourcePriority) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty priority")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.Priority
	err := apiutil.All(ctx, apiutil.AllFunc(func(offset int) (bool, error) {
		list, err := d.client.ListPrioritiesWithContext(ctx, pagerduty.ListPrioritiesOptions{
			Limit:  apiutil.Limit,
			Offset: uint(offset),
		})
		if err != nil {
			return false, err
		}

		for _, priority := range list.Priorities {
			if strings.EqualFold(priority.Name, searchName.ValueString()) {
				found = &priority
				return false, nil
			}
		}

		return list.More, nil
	}))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty priority %s", searchName),
			err.Error(),
		)
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any priority with the name: %s", searchName),
			"",
		)
		return
	}

	model := dataSourcePriorityModel{
		ID:          types.StringValue(found.ID),
		Name:        types.StringValue(found.Name),
		Description: types.StringValue(found.Description),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourcePriorityModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}
