package pagerduty

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

type ErrorResponse struct {
	ShouldReturn bool
	ReturnVal    *resource.RetryError
}

// Provider represents a resource provider in Terraform
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"skip_credentials_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PAGERDUTY_TOKEN", nil),
			},

			"user_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PAGERDUTY_USER_TOKEN", nil),
			},

			"service_region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"api_url_override": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"pagerduty_escalation_policy":   dataSourcePagerDutyEscalationPolicy(),
			"pagerduty_escalation_policies": dataSourcePagerDutyEscalationPolicies(),
			"pagerduty_schedule":            dataSourcePagerDutySchedule(),
			"pagerduty_user":                dataSourcePagerDutyUser(),
			"pagerduty_user_contact_method": dataSourcePagerDutyUserContactMethod(),
			"pagerduty_team":                dataSourcePagerDutyTeam(),
			"pagerduty_vendor":              dataSourcePagerDutyVendor(),
			"pagerduty_vendors":             dataSourcePagerDutyVendors(),
			"pagerduty_extension_schema":    dataSourcePagerDutyExtensionSchema(),
			"pagerduty_service":             dataSourcePagerDutyService(),
			"pagerduty_service_integration": dataSourcePagerDutyServiceIntegration(),
			"pagerduty_business_service":    dataSourcePagerDutyBusinessService(),
			"pagerduty_priority":            dataSourcePagerDutyPriority(),
			"pagerduty_ruleset":             dataSourcePagerDutyRuleset(),
			"pagerduty_tag":                 dataSourcePagerDutyTag(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"pagerduty_addon":                       resourcePagerDutyAddon(),
			"pagerduty_escalation_policy":           resourcePagerDutyEscalationPolicy(),
			"pagerduty_maintenance_window":          resourcePagerDutyMaintenanceWindow(),
			"pagerduty_schedule":                    resourcePagerDutySchedule(),
			"pagerduty_service":                     resourcePagerDutyService(),
			"pagerduty_service_integration":         resourcePagerDutyServiceIntegration(),
			"pagerduty_team":                        resourcePagerDutyTeam(),
			"pagerduty_team_membership":             resourcePagerDutyTeamMembership(),
			"pagerduty_user":                        resourcePagerDutyUser(),
			"pagerduty_user_contact_method":         resourcePagerDutyUserContactMethod(),
			"pagerduty_user_notification_rule":      resourcePagerDutyUserNotificationRule(),
			"pagerduty_extension":                   resourcePagerDutyExtension(),
			"pagerduty_extension_servicenow":        resourcePagerDutyExtensionServiceNow(),
			"pagerduty_event_rule":                  resourcePagerDutyEventRule(),
			"pagerduty_ruleset":                     resourcePagerDutyRuleset(),
			"pagerduty_ruleset_rule":                resourcePagerDutyRulesetRule(),
			"pagerduty_business_service":            resourcePagerDutyBusinessService(),
			"pagerduty_service_dependency":          resourcePagerDutyServiceDependency(),
			"pagerduty_response_play":               resourcePagerDutyResponsePlay(),
			"pagerduty_tag":                         resourcePagerDutyTag(),
			"pagerduty_tag_assignment":              resourcePagerDutyTagAssignment(),
			"pagerduty_service_event_rule":          resourcePagerDutyServiceEventRule(),
			"pagerduty_slack_connection":            resourcePagerDutySlackConnection(),
			"pagerduty_business_service_subscriber": resourcePagerDutyBusinessServiceSubscriber(),
			"pagerduty_webhook_subscription":        resourcePagerDutyWebhookSubscription(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			terraformVersion = "0.12+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func isErrCode(err error, code int) bool {
	currentErr := err
	for errors.Unwrap(currentErr) != nil {
		currentErr = errors.Unwrap(currentErr)

		if e, ok := currentErr.(*pagerduty.Error); ok && e.ErrorResponse.Response.StatusCode == code {
			log.Printf("[INFO] Error code matches expected %d", code)
			return true
		}
	}

	log.Printf("[INFO] Error code doesn't match expected %d", code)

	return false
}

func genError(err error, d *schema.ResourceData) error {
	resId := "<ID missing>"
	if d != nil && d.Id() != "" {
		resId = d.Id()
	}
	errStr := "unknown error"
	if err != nil {
		errStr = err.Error()
	}
	return fmt.Errorf("error reading: %s: %s", resId, errStr)
}

func handleGenericErrors(err error, d *schema.ResourceData) ErrorResponse {
	if err == nil {
		return ErrorResponse{ShouldReturn: false}
	}

	if isErrCode(err, 429) {
		// Delaying retry by 30s as recommended by PagerDuty
		// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
		time.Sleep(30 * time.Second)
		return ErrorResponse{ShouldReturn: true, ReturnVal: resource.RetryableError(err)}
	}

	if isErrCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return ErrorResponse{ShouldReturn: true, ReturnVal: nil}
	}

	generatedError := genError(err, d)
	if generatedError != nil {
		return ErrorResponse{ShouldReturn: true, ReturnVal: resource.NonRetryableError(generatedError)}
	}

	return ErrorResponse{ShouldReturn: false}
}

func handleRateLimitErrors(err error, d *schema.ResourceData) ErrorResponse {
	if err == nil {
		return ErrorResponse{ShouldReturn: false}
	}

	if isErrCode(err, 429) {
		// Delaying retry by 30s as recommended by PagerDuty
		// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
		time.Sleep(30 * time.Second)
		return ErrorResponse{ShouldReturn: true, ReturnVal: resource.RetryableError(err)}
	}
	generatedError := genError(err, d)
	if generatedError != nil {
		return ErrorResponse{ShouldReturn: true, ReturnVal: resource.NonRetryableError(generatedError)}
	}

	return ErrorResponse{ShouldReturn: false}
}

func getErrorHandler(shouldHandle404Errors bool) func(err error, d *schema.ResourceData) ErrorResponse {
	if shouldHandle404Errors {
		return handleGenericErrors
	}
	return handleRateLimitErrors
}

func providerConfigure(data *schema.ResourceData, terraformVersion string) (interface{}, error) {
	var ServiceRegion = strings.ToLower(data.Get("service_region").(string))

	if ServiceRegion == "us" || ServiceRegion == "" {
		ServiceRegion = ""
	} else {
		ServiceRegion = ServiceRegion + "."
	}

	config := Config{
		ApiUrl:              "https://api." + ServiceRegion + "pagerduty.com",
		AppUrl:              "https://app." + ServiceRegion + "pagerduty.com",
		SkipCredsValidation: data.Get("skip_credentials_validation").(bool),
		Token:               data.Get("token").(string),
		UserToken:           data.Get("user_token").(string),
		UserAgent:           fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, terraformVersion),
		ApiUrlOverride:      data.Get("api_url_override").(string),
	}

	log.Println("[INFO] Initializing PagerDuty client")
	return &config, nil
}
