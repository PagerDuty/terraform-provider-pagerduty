package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
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

type dataSourceLicense struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceLicense)(nil)

func (*dataSourceLicense) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_license"
}

func (*dataSourceLicense) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: dataSourceLicenseAttributes}
}

var dataSourceLicenseAttributes = map[string]schema.Attribute{
	"id":                    schema.StringAttribute{Optional: true, Computed: true},
	"name":                  schema.StringAttribute{Optional: true, Computed: true},
	"description":           schema.StringAttribute{Optional: true, Computed: true},
	"type":                  schema.StringAttribute{Computed: true},
	"summary":               schema.StringAttribute{Computed: true},
	"role_group":            schema.StringAttribute{Computed: true},
	"current_value":         schema.Int64Attribute{Computed: true},
	"allocations_available": schema.Int64Attribute{Computed: true},
	"valid_roles":           schema.ListAttribute{Computed: true, ElementType: types.StringType},
	"self":                  schema.StringAttribute{Computed: true},
	"html_url":              schema.StringAttribute{Computed: true},
}

func (d *dataSourceLicense) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceLicense) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Fetching PagerDuty licenses")

	var searchName, searchID, searchDescription types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &searchName)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &searchID)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description"), &searchDescription)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.License
	err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		list, err := d.client.ListLicensesWithContext(ctx)
		if err != nil {
			if util.IsBadRequestError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		found = findBestMatchLicense(list.Licenses, searchID.ValueString(), searchName.ValueString(), searchDescription.ValueString())
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading PagerDuty license %s", searchName),
			err.Error(),
		)
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any license with the name: %s", searchName),
			"",
		)
		return
	}

	model := flattenLicense(found)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceLicenseModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Type                 types.String `tfsdk:"type"`
	Summary              types.String `tfsdk:"summary"`
	RoleGroup            types.String `tfsdk:"role_group"`
	AllocationsAvailable types.Int64  `tfsdk:"allocations_available"`
	CurrentValue         types.Int64  `tfsdk:"current_value"`
	ValidRoles           types.List   `tfsdk:"valid_roles"`
	Self                 types.String `tfsdk:"self"`
	HTMLURL              types.String `tfsdk:"html_url"`
}

func flattenLicense(license *pagerduty.License) dataSourceLicenseModel {
	return dataSourceLicenseModel{
		ID:                   types.StringValue(license.ID),
		Name:                 types.StringValue(license.Name),
		Type:                 types.StringValue(license.Type),
		Description:          types.StringValue(license.Description),
		Summary:              types.StringValue(license.Summary),
		RoleGroup:            types.StringValue(license.RoleGroup),
		AllocationsAvailable: types.Int64Value(int64(license.AllocationsAvailable)),
		CurrentValue:         types.Int64Value(int64(license.CurrentValue)),
		Self:                 types.StringValue(license.Self),
		HTMLURL:              types.StringValue(license.HTMLURL),
		ValidRoles:           flattenLicenseValidRoles(license.ValidRoles),
	}
}

func flattenLicenseValidRoles(roles []string) types.List {
	elements := make([]attr.Value, 0, len(roles))
	for _, e := range roles {
		elements = append(elements, types.StringValue(e))
	}
	return types.ListValueMust(types.StringType, elements)
}

func findBestMatchLicense(licenses []pagerduty.License, id, name, description string) *pagerduty.License {
	var found *pagerduty.License
	for _, license := range licenses {
		if licenseIsExactMatch(&license, id, name, description) {
			found = &license
			break
		}
	}

	// If there is no exact match for a license, check for substring matches
	// This allows customers to use a term such as "Full User", which is included
	// in the names of all licenses that support creating full users. However,
	// if id is set then it must match with licenseIsExactMatch
	if id == "" && found == nil {
		for _, license := range licenses {
			if licenseContainsMatch(&license, name, description) {
				found = &license
				break
			}
		}
	}

	return found
}

func licenseIsExactMatch(license *pagerduty.License, id, name, description string) bool {
	if id != "" {
		return license.ID == id &&
			(license.Name == name || name == "") &&
			(license.Description == description || description == "")
	}
	return license.Name == name && license.Description == description
}

func licenseContainsMatch(license *pagerduty.License, name, description string) bool {
	return strings.Contains(license.Name, name) && strings.Contains(license.Description, description)
}
