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

type dataSourceIncidentType struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceIncidentType)(nil)

func (*dataSourceIncidentType) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_incident_type"
}

func (*dataSourceIncidentType) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"name":         schema.StringAttribute{Computed: true},
			"type":         schema.StringAttribute{Computed: true},
			"display_name": schema.StringAttribute{Required: true},
			"description":  schema.StringAttribute{Computed: true},
			"parent_type":  schema.StringAttribute{Computed: true},
			"enabled":      schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *dataSourceIncidentType) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceIncidentType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty incident type")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("display_name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.IncidentType
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		response, err := d.client.ListIncidentTypes(ctx, pagerduty.ListIncidentTypesOptions{Filter: "all"})
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		for _, it := range response.IncidentTypes {
			if it.DisplayName == searchName.ValueString() {
				found = &it
				break
			}
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty incident type %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any incident type with the name: %s", searchName),
			"",
		)
		return
	}

	model := dataSourceIncidentTypeModel{
		ID:          types.StringValue(found.ID),
		Name:        types.StringValue(found.Name),
		Type:        types.StringValue(found.Type),
		DisplayName: types.StringValue(found.DisplayName),
		Description: types.StringValue(found.Description),
		Enabled:     types.BoolValue(found.Enabled),
	}
	if found.Parent != nil {
		model.ParentType = types.StringValue(found.Parent.ID)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceIncidentTypeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	ParentType  types.String `tfsdk:"parent_type"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}
