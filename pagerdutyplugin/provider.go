package pagerduty

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Provider struct{}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	config, diags := ReadConfig(ctx, req)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	client, err := config.Client()
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(
			"Cannot obtain plugin client",
			err.Error(),
		))
	}
	resp.DataSourceData = client
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pagerduty"
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url_override":            schema.StringAttribute{Optional: true},
			"service_region":              schema.StringAttribute{Optional: true},
			"skip_credentials_validation": schema.BoolAttribute{Optional: true},
			"token":                       schema.StringAttribute{Optional: true},
			"user_token":                  schema.StringAttribute{Optional: true},
		},

		Blocks: map[string]schema.Block{
			"use_app_oauth_scoped_token": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"pd_client_id": schema.StringAttribute{
							Required: true,
							// DefaultFunc: schema.EnvDefaultFunc("PAGERDUTY_CLIENT_ID", nil),
						},
						"pd_client_secret": schema.StringAttribute{
							Required: true,
							// DefaultFunc: schema.EnvDefaultFunc("PAGERDUTY_CLIENT_SECRET", nil),
						},
						"pd_subdomain": schema.StringAttribute{
							Required: true,
							// DefaultFunc: schema.EnvDefaultFunc("PAGERDUTY_SUBDOMAIN", nil),
						},
					},
				},
			},
		},
	}
}

func (p *Provider) DataSources(ctx context.Context) [](func() datasource.DataSource) {
	return [](func() datasource.DataSource){
		func() datasource.DataSource { return &dataSourceStandards{} },
	}
}

func (p *Provider) Resources(ctx context.Context) [](func() resource.Resource) {
	return [](func() resource.Resource){}
}

func New() provider.Provider {
	return &Provider{}
}

type providerArguments struct {
	Token                     types.String `tfsdk:"token"`
	UserToken                 types.String `tfsdk:"user_token"`
	SkipCredentialsValidation types.Bool   `tfsdk:"skip_credentials_validation"`
	ServiceRegion             types.String `tfsdk:"service_region"`
	ApiUrlOverride            types.String `tfsdk:"api_url_override"`
	UseAppOauthScopedToken    *struct {
		PdClientId     types.String `tfsdk:"pd_client_id"`
		PdClientSecret types.String `tfsdk:"pd_client_secret"`
		PdDomain       types.String `tfsdk:"pd_domain"`
	} `tfsdk:"use_app_oauth_scoped_token"`
}

func ReadConfig(ctx context.Context, req provider.ConfigureRequest) (*Config, diag.Diagnostics) {
	var diags diag.Diagnostics
	var args providerArguments
	diags.Append(req.Config.Get(ctx, &args)...)

	serviceRegion := args.ServiceRegion.ValueString()
	var regionApiUrl string
	if serviceRegion == "us" || serviceRegion == "" {
		regionApiUrl = ""
	} else {
		regionApiUrl = serviceRegion + "."
	}

	skipCredentialsValidation := args.SkipCredentialsValidation.Equal(types.BoolValue(true))

	config := Config{
		ApiUrl:              "https://api." + regionApiUrl + "pagerduty.com",
		AppUrl:              "https://app." + regionApiUrl + "pagerduty.com",
		SkipCredsValidation: skipCredentialsValidation,
		Token:               args.Token.ValueString(),
		UserToken:           args.UserToken.ValueString(),
		TerraformVersion:    req.TerraformVersion,
		ApiUrlOverride:      args.ApiUrlOverride.ValueString(),
		ServiceRegion:       serviceRegion,
	}

	if config.Token == "" {
		config.Token = os.Getenv("PAGERDUTY_TOKEN")
	}
	if config.UserToken == "" {
		config.UserToken = os.Getenv("PAGERDUTY_USER_TOKEN")
	}

	log.Println("[INFO] Initializing PagerDuty plugin client")
	return &config, diags
}
