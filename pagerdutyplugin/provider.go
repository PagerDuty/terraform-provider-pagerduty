package pagerduty

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Provider struct {
	client *pagerduty.Client
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pagerduty"
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	useAppOauthScopedTokenBlock := schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"pd_client_id":     schema.StringAttribute{Optional: true},
				"pd_client_secret": schema.StringAttribute{Optional: true},
				"pd_subdomain":     schema.StringAttribute{Optional: true},
			},
		},
	}
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url_override":            schema.StringAttribute{Optional: true},
			"service_region":              schema.StringAttribute{Optional: true},
			"skip_credentials_validation": schema.BoolAttribute{Optional: true},
			"token":                       schema.StringAttribute{Optional: true},
			"user_token":                  schema.StringAttribute{Optional: true},
		},
		Blocks: map[string]schema.Block{
			"use_app_oauth_scoped_token": useAppOauthScopedTokenBlock,
		},
	}
}

func (p *Provider) DataSources(ctx context.Context) [](func() datasource.DataSource) {
	return [](func() datasource.DataSource){
		func() datasource.DataSource { return &dataSourceBusinessService{} },
		func() datasource.DataSource { return &dataSourceStandardsResourceScores{} },
		func() datasource.DataSource { return &dataSourceStandardsResourcesScores{} },
		func() datasource.DataSource { return &dataSourceStandards{} },
	}
}

func (p *Provider) Resources(ctx context.Context) [](func() resource.Resource) {
	return [](func() resource.Resource){
		func() resource.Resource { return &resourceBusinessService{} },
		func() resource.Resource { return &resourceServiceDependency{} },
	}
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var args providerArguments
	resp.Diagnostics.Append(req.Config.Get(ctx, &args)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceRegion := args.ServiceRegion.ValueString()
	if serviceRegion == "" {
		if v, ok := os.LookupEnv("PAGERDUTY_SERVICE_REGION"); ok && v != "" {
			serviceRegion = v
		} else {
			serviceRegion = "us"
		}
	}

	var regionApiUrl string
	if serviceRegion == "us" {
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

	if !args.UseAppOauthScopedToken.IsNull() {
		blockList := []UseAppOauthScopedToken{}
		resp.Diagnostics.Append(args.UseAppOauthScopedToken.ElementsAs(ctx, &blockList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.AppOauthScopedToken = &AppOauthScopedToken{
			ClientId:     blockList[0].PdClientId.ValueString(),
			ClientSecret: blockList[0].PdClientSecret.ValueString(),
			Subdomain:    blockList[0].PdSubdomain.ValueString(),
		}
	}

	if args.UseAppOauthScopedToken.IsNull() {
		if config.Token == "" {
			config.Token = os.Getenv("PAGERDUTY_TOKEN")
		}
		if config.UserToken == "" {
			config.UserToken = os.Getenv("PAGERDUTY_USER_TOKEN")
		}
	} else {
		if config.AppOauthScopedToken.ClientId == "" {
			config.AppOauthScopedToken.ClientId = os.Getenv("PAGERDUTY_CLIENT_ID")
		}
		if config.AppOauthScopedToken.ClientSecret == "" {
			config.AppOauthScopedToken.ClientSecret = os.Getenv("PAGERDUTY_CLIENT_SECRET")
		}
		if config.AppOauthScopedToken.Subdomain == "" {
			config.AppOauthScopedToken.Subdomain = os.Getenv("PAGERDUTY_SUBDOMAIN")
		}
	}

	if config.AppOauthScopedToken != nil {
		// While doing migration to terraform plugin framework, because
		// of a limitation of the provider mux
		// https://github.com/hashicorp/terraform-plugin-framework/issues/539
		// We had to define pd_client_id, pd_client_secret, and pd_subdomain
		// as Optional and manually check its presence here.
		li := []string{}
		if config.AppOauthScopedToken.ClientId == "" {
			li = append(li, "pd_client_id")
		}
		if config.AppOauthScopedToken.ClientSecret == "" {
			li = append(li, "pd_client_secret")
		}
		if config.AppOauthScopedToken.Subdomain == "" {
			li = append(li, "pd_subdomain")
		}
		if len(li) > 0 {
			resp.Diagnostics.AddError(
				fmt.Sprintf(`Missing required arguments: "%v"`, strings.Join(li, `", "`)),
				"Despite being defined as Optional, its value is required and no definition was found.",
			)
			return
		}
	}

	log.Println("[INFO] Initializing PagerDuty plugin client")

	client, err := config.Client(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Cannot obtain plugin client", err.Error())
	}
	p.client = client
	resp.DataSourceData = client
	resp.ResourceData = client
}

type UseAppOauthScopedToken struct {
	PdClientId     types.String `tfsdk:"pd_client_id"`
	PdClientSecret types.String `tfsdk:"pd_client_secret"`
	PdSubdomain    types.String `tfsdk:"pd_subdomain"`
}

type providerArguments struct {
	Token                     types.String `tfsdk:"token"`
	UserToken                 types.String `tfsdk:"user_token"`
	SkipCredentialsValidation types.Bool   `tfsdk:"skip_credentials_validation"`
	ServiceRegion             types.String `tfsdk:"service_region"`
	ApiUrlOverride            types.String `tfsdk:"api_url_override"`
	UseAppOauthScopedToken    types.List   `tfsdk:"use_app_oauth_scoped_token"`
}
