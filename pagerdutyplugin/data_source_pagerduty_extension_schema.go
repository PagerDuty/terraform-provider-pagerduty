package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type dataSourceExtensionSchema struct{ client *pagerduty.Client }

var _ datasource.DataSourceWithConfigure = (*dataSourceExtensionSchema)(nil)

func (*dataSourceExtensionSchema) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "pagerduty_extension_schema"
}

func (*dataSourceExtensionSchema) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
			"type": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *dataSourceExtensionSchema) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.Append(ConfigurePagerdutyClient(&d.client, req.ProviderData)...)
}

func (d *dataSourceExtensionSchema) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	log.Println("[INFO] Reading PagerDuty extension schema")

	var searchName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &searchName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *pagerduty.ExtensionSchema
	offset := 0
	more := true
	for more {
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			o := pagerduty.ListExtensionSchemaOptions{Limit: 20, Offset: uint(offset), Total: true}
			list, err := d.client.ListExtensionSchemasWithContext(ctx, o)
			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}

			for _, extensionSchema := range list.ExtensionSchemas {
				if strings.EqualFold(extensionSchema.Label, searchName.ValueString()) {
					found = &extensionSchema
					more = false
					return nil
				}
			}

			more = list.More
			offset += len(list.ExtensionSchemas)
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading PagerDuty extension schema %s", searchName),
				err.Error(),
			)
		}
	}

	if found == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to locate any extension schema with the name: %s", searchName),
			"",
		)
		return
	}

	model := dataSourceExtensionSchemaModel{
		ID:   types.StringValue(found.ID),
		Name: types.StringValue(found.Label),
		Type: types.StringValue(found.Type),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type dataSourceExtensionSchemaModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}
