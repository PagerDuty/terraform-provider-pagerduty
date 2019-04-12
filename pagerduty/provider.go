package pagerduty

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// Provider represents a resource provider in Terraform
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"pagerduty_escalation_policy": dataSourcePagerDutyEscalationPolicy(),
			"pagerduty_schedule":          dataSourcePagerDutySchedule(),
			"pagerduty_user":              dataSourcePagerDutyUser(),
			"pagerduty_team":              dataSourcePagerDutyTeam(),
			"pagerduty_vendor":            dataSourcePagerDutyVendor(),
			"pagerduty_extension_schema":  dataSourcePagerDutyExtensionSchema(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"pagerduty_addon":                  resourcePagerDutyAddon(),
			"pagerduty_escalation_policy":      resourcePagerDutyEscalationPolicy(),
			"pagerduty_maintenance_window":     resourcePagerDutyMaintenanceWindow(),
			"pagerduty_schedule":               resourcePagerDutySchedule(),
			"pagerduty_service":                resourcePagerDutyService(),
			"pagerduty_service_integration":    resourcePagerDutyServiceIntegration(),
			"pagerduty_team":                   resourcePagerDutyTeam(),
			"pagerduty_team_membership":        resourcePagerDutyTeamMembership(),
			"pagerduty_user":                   resourcePagerDutyUser(),
			"pagerduty_user_contact_method":    resourcePagerDutyUserContactMethod(),
			"pagerduty_user_notification_rule": resourcePagerDutyUserNotificationRule(),
			"pagerduty_extension":              resourcePagerDutyExtension(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func handleNotFoundError(err error, d *schema.ResourceData) error {
	if perr, ok := err.(*pagerduty.Error); ok && perr.ErrorResponse.StatusCode == 404 {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return nil
	}

	return fmt.Errorf("Error reading: %s: %s", d.Id(), err)
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	config := Config{
		SkipCredsValidation: data.Get("skip_credentials_validation").(bool),
		Token:               data.Get("token").(string),
	}

	log.Println("[INFO] Initializing PagerDuty client")
	return config.Client()
}
