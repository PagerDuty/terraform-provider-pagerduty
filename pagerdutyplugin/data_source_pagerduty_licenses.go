package pagerduty

import (
	"context"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceLicenses struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceLicenses)(nil)

func (*dataSourceLicenses) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_licenses"
}

func (*dataSourceLicenses) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Optional: true},
			"licenses": schema.ListAttribute{
				Computed:    true,
				ElementType: licenseObjectType,
			},
		},
	}
}

func (d *dataSourceLicenses) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceLicenses) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model dataSourceLicensesModel
	log.Println("[INFO] Reading PagerDuty licenses")

	diags := req.Config.Get(ctx, &model)
	if resp.Diagnostics.Append(diags...); diags.HasError() {
		return
	}

	uid := ""
	if model.ID.IsNull() {
		uid = id.UniqueId()
	} else {
		uid = model.ID.ValueString()
	}

	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		list, err := d.client.ListLicensesWithContext(ctx)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		model = flattenLicenses(uid, list.Licenses, &resp.Diagnostics)
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading PagerDuty licenses",
			err.Error(),
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceLicensesModel struct {
	ID       types.String `tfsdk:"id"`
	Licenses types.List   `tfsdk:"licenses"`
}

func flattenLicenses(uid string, licenses []pagerduty.License, diags *diag.Diagnostics) dataSourceLicensesModel {
	elements := make([]attr.Value, 0, len(licenses))
	for _, license := range licenses {
		model := flattenLicense(&license)
		e, d := types.ObjectValue(licenseObjectType.AttrTypes, map[string]attr.Value{
			"id":                    model.ID,
			"name":                  model.Name,
			"description":           model.Description,
			"type":                  model.Type,
			"summary":               model.Summary,
			"role_group":            model.RoleGroup,
			"current_value":         model.CurrentValue,
			"allocations_available": model.AllocationsAvailable,
			"valid_roles":           model.ValidRoles,
			"self":                  model.Self,
			"html_url":              model.HTMLURL,
		})
		diags.Append(d...)
		if d.HasError() {
			continue
		}
		elements = append(elements, e)
	}

	return dataSourceLicensesModel{
		Licenses: types.ListValueMust(licenseObjectType, elements),
		ID:       types.StringValue(uid),
	}
}

var licenseObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"description":           types.StringType,
		"type":                  types.StringType,
		"summary":               types.StringType,
		"role_group":            types.StringType,
		"current_value":         types.Int64Type,
		"allocations_available": types.Int64Type,
		"valid_roles":           types.ListType{ElemType: types.StringType},
		"self":                  types.StringType,
		"html_url":              types.StringType,
	},
}
