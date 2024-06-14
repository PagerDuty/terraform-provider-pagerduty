package pagerduty

import (
	"context"
	"fmt"
	"log"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util/apiutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceService struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceService)(nil)

func (d *dataSourceService) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_service"
}

func (d *dataSourceService) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                      schema.StringAttribute{Computed: true},
			"name":                    schema.StringAttribute{Required: true},
			"auto_resolve_timeout":    schema.Int64Attribute{Computed: true},
			"acknowledgement_timeout": schema.Int64Attribute{Computed: true},
			"alert_creation":          schema.StringAttribute{Computed: true},
			"description":             schema.StringAttribute{Computed: true},
			"escalation_policy":       schema.StringAttribute{Computed: true},
			"type":                    schema.StringAttribute{Computed: true},
			"teams": schema.ListAttribute{
				Computed:    true,
				Description: "The set of teams associated with the service",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":   types.StringType,
						"name": types.StringType,
					},
				},
			},
		},
	}
}

func (d *dataSourceService) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceService) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Printf("[INFO] Reading PagerDuty service")

	var searchName types.String
	if d := req.Config.GetAttribute(ctx, path.Root("name"), &searchName); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	var found *pagerduty.Service
	err := apiutil.All(ctx, func(offset int) (bool, error) {
		resp, err := d.client.ListServicesWithContext(ctx, pagerduty.ListServiceOptions{
			Query:    searchName.ValueString(),
			Limit:    apiutil.Limit,
			Offset:   uint(offset),
			Includes: []string{"teams"},
		})
		if err != nil {
			return false, err
		}

		for _, service := range resp.Services {
			if service.Name == searchName.ValueString() {
				found = &service
				return false, nil
			}
		}

		return resp.More, nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error searching Service %s", searchName),
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any service with the name: %s", searchName),
			"",
		)
		return
	}

	model := flattenServiceData(found, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceServiceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	AutoResolveTimeout     types.Int64  `tfsdk:"auto_resolve_timeout"`
	AcknowledgementTimeout types.Int64  `tfsdk:"acknowledgement_timeout"`
	AlertCreation          types.String `tfsdk:"alert_creation"`
	Description            types.String `tfsdk:"description"`
	EscalationPolicy       types.String `tfsdk:"escalation_policy"`
	Type                   types.String `tfsdk:"type"`
	Teams                  types.List   `tfsdk:"teams"`
}

func flattenServiceData(service *pagerduty.Service, diags *diag.Diagnostics) dataSourceServiceModel {
	teamObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
	}

	teamsElems := make([]attr.Value, 0, len(service.Teams))
	for _, t := range service.Teams {
		teamObj := types.ObjectValueMust(teamObjectType.AttrTypes, map[string]attr.Value{
			"id":   types.StringValue(t.ID),
			"name": types.StringValue(t.Name),
		})
		teamsElems = append(teamsElems, teamObj)
	}

	teams, d := types.ListValue(teamObjectType, teamsElems)
	if diags.Append(d...); d.HasError() {
		return dataSourceServiceModel{}
	}

	model := dataSourceServiceModel{
		ID:                     types.StringValue(service.ID),
		Name:                   types.StringValue(service.Name),
		Type:                   types.StringValue(service.Type),
		AutoResolveTimeout:     types.Int64Null(),
		AcknowledgementTimeout: types.Int64Null(),
		AlertCreation:          types.StringValue(service.AlertCreation),
		Description:            types.StringValue(service.Description),
		EscalationPolicy:       types.StringValue(service.EscalationPolicy.ID),
		Teams:                  teams,
	}

	if service.AutoResolveTimeout != nil {
		model.AutoResolveTimeout = types.Int64Value(int64(*service.AutoResolveTimeout))
	}
	if service.AcknowledgementTimeout != nil {
		model.AcknowledgementTimeout = types.Int64Value(int64(*service.AcknowledgementTimeout))
	}
	return model
}
