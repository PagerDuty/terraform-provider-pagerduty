package pagerduty

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// Config defines the configuration options for the PagerDuty client
type Config struct {
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

	// UserAgent for API Client
	UserAgent string
}

const invalidCreds = `

No valid credentials found for PagerDuty provider.
Please see https://www.terraform.io/docs/providers/pagerduty/index.html
for more information on providing credentials for this provider.
`

// Client returns a new PagerDuty client
func (c *Config) Client() (*pagerduty.Client, error) {
	// Validate that the PagerDuty token is set
	if c.Token == "" {
		return nil, fmt.Errorf(invalidCreds)
	}

	var httpClient *http.Client
	httpClient = http.DefaultClient
	httpClient.Transport = logging.NewTransport("PagerDuty", http.DefaultTransport)

	var apiUrl = c.ApiUrl
	if c.ApiUrlOverride != "" {
		apiUrl = c.ApiUrlOverride
	}

	config := &pagerduty.Config{
		BaseURL:    apiUrl,
		Debug:      logging.IsDebugOrHigher(),
		HTTPClient: httpClient,
		Token:      c.Token,
		UserAgent:  c.UserAgent,
	}

	client, err := pagerduty.NewClient(config)
	if err != nil {
		return nil, err
	}

	if !c.SkipCredsValidation {
		// Validate the credentials by calling the abilities endpoint,
		// if we get a 401 response back we return an error to the user
		if err := client.ValidateAuth(); err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%s\n%s", err, invalidCreds))
		}
	}

	log.Printf("[INFO] PagerDuty client configured")

	return client, nil
}

func (c *Config) SlackClient() (*pagerduty.Client, error) {
	// Validate that the user level PagerDuty token is set
	if c.UserToken == "" {
		return nil, fmt.Errorf(invalidCreds)
	}

	var httpClient *http.Client
	httpClient = http.DefaultClient
	httpClient.Transport = logging.NewTransport("PagerDuty", http.DefaultTransport)

	config := &pagerduty.Config{
		BaseURL:    c.AppUrl,
		Debug:      logging.IsDebugOrHigher(),
		HTTPClient: httpClient,
		Token:      c.UserToken,
		UserAgent:  c.UserAgent,
	}

	client, err := pagerduty.NewClient(config)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] PagerDuty client configured for slack")

	return client, nil
}
