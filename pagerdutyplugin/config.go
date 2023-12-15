package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/terraform-providers/terraform-provider-pagerduty/util"
)

// Config defines the configuration options for the PagerDuty client
type Config struct {
	mu sync.Mutex

	// The PagerDuty API URL
	ApiUrl string

	// Override default PagerDuty API URL
	ApiUrlOverride string

	// The PagerDuty APP URL
	AppUrl string

	// The PagerDuty API V2 token
	Token string

	// The PagerDuty User level token for Slack
	UserToken string

	// Skip validation of the token against the PagerDuty API
	SkipCredsValidation bool

	// Target version for terraform
	TerraformVersion string

	// Region where the server of the service is deployed
	ServiceRegion string

	// API wrapper
	client *pagerduty.Client
}

const invalidCreds = `

No valid credentials found for PagerDuty provider.
Please see https://www.terraform.io/docs/providers/pagerduty/index.html
for more information on providing credentials for this provider.
`

// Client returns a PagerDuty client, initializing when necessary.
func (c *Config) Client() (*pagerduty.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Return the previously-configured client if available.
	if c.client != nil {
		return c.client, nil
	}

	// Validate that the PagerDuty token is set
	if c.Token == "" {
		return nil, fmt.Errorf(invalidCreds)
	}

	httpClient := http.DefaultClient
	httpClient.Transport = logging.NewTransport("PagerDuty", http.DefaultTransport)

	apiUrl := c.ApiUrl
	if c.ApiUrlOverride != "" {
		apiUrl = c.ApiUrlOverride
	}

	client := pagerduty.NewClient(c.Token, []pagerduty.ClientOptions{
		pagerduty.WithAPIEndpoint(apiUrl),
		WithHTTPClient(httpClient),
		pagerduty.WithTerraformProvider(c.TerraformVersion),
		// TODO: c.AppOauthScopedTokenParams
		// TODO: c.APITokenType
	}...)

	// TODO: oauth validation
	// if !c.SkipCredsValidation {
	// 	// Validate the credentials by calling the abilities endpoint,
	// 	// if we get a 401 response back we return an error to the user
	// 	if err := client.ValidateAuth(); err != nil {
	// 		return nil, fmt.Errorf(fmt.Sprintf("%s\n%s", err, invalidCreds))
	// 	}
	// }
	c.client = client

	log.Printf("[INFO] PagerDuty plugin client configured")
	return c.client, nil
}

func WithHTTPClient(httpClient pagerduty.HTTPClient) pagerduty.ClientOptions {
	return func(c *pagerduty.Client) {
		if util.IsNilFunc(httpClient) {
			return
		}
		c.HTTPClient = httpClient
	}
}

// ConfigurePagerdutyClient sets a pagerduty API client in a pointer to the
// property of any data source struct from the general configuration of the
// provider.
func ConfigurePagerdutyClient(clientPtr **pagerduty.Client, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pagerduty.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *github.com/PagerDuty/go-pagerduty.Client, got: %T."+
					"Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}
	if clientPtr != nil {
		*clientPtr = client
	}
}
