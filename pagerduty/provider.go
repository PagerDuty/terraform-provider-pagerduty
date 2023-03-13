package pagerduty

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

var (
	version = "dev"
	commit  = "none"
)

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
			"pagerduty_escalation_policy":         dataSourcePagerDutyEscalationPolicy(),
			"pagerduty_schedule":                  dataSourcePagerDutySchedule(),
			"pagerduty_user":                      dataSourcePagerDutyUser(),
			"pagerduty_users":                     dataSourcePagerDutyUsers(),
			"pagerduty_user_contact_method":       dataSourcePagerDutyUserContactMethod(),
			"pagerduty_team":                      dataSourcePagerDutyTeam(),
			"pagerduty_vendor":                    dataSourcePagerDutyVendor(),
			"pagerduty_extension_schema":          dataSourcePagerDutyExtensionSchema(),
			"pagerduty_service":                   dataSourcePagerDutyService(),
			"pagerduty_service_integration":       dataSourcePagerDutyServiceIntegration(),
			"pagerduty_business_service":          dataSourcePagerDutyBusinessService(),
			"pagerduty_priority":                  dataSourcePagerDutyPriority(),
			"pagerduty_ruleset":                   dataSourcePagerDutyRuleset(),
			"pagerduty_tag":                       dataSourcePagerDutyTag(),
			"pagerduty_event_orchestration":       dataSourcePagerDutyEventOrchestration(),
			"pagerduty_event_orchestrations":      dataSourcePagerDutyEventOrchestrations(),
			"pagerduty_automation_actions_runner": dataSourcePagerDutyAutomationActionsRunner(),
			"pagerduty_automation_actions_action": dataSourcePagerDutyAutomationActionsAction(),
			"pagerduty_incident_workflow":         dataSourcePagerDutyIncidentWorkflow(),
			"pagerduty_custom_field":              dataSourcePagerDutyField(),
			"pagerduty_custom_field_schema":       dataSourcePagerDutyFieldSchema(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"pagerduty_addon":                                         resourcePagerDutyAddon(),
			"pagerduty_escalation_policy":                             resourcePagerDutyEscalationPolicy(),
			"pagerduty_maintenance_window":                            resourcePagerDutyMaintenanceWindow(),
			"pagerduty_schedule":                                      resourcePagerDutySchedule(),
			"pagerduty_service":                                       resourcePagerDutyService(),
			"pagerduty_service_integration":                           resourcePagerDutyServiceIntegration(),
			"pagerduty_team":                                          resourcePagerDutyTeam(),
			"pagerduty_team_membership":                               resourcePagerDutyTeamMembership(),
			"pagerduty_user":                                          resourcePagerDutyUser(),
			"pagerduty_user_contact_method":                           resourcePagerDutyUserContactMethod(),
			"pagerduty_user_notification_rule":                        resourcePagerDutyUserNotificationRule(),
			"pagerduty_extension":                                     resourcePagerDutyExtension(),
			"pagerduty_extension_servicenow":                          resourcePagerDutyExtensionServiceNow(),
			"pagerduty_event_rule":                                    resourcePagerDutyEventRule(),
			"pagerduty_ruleset":                                       resourcePagerDutyRuleset(),
			"pagerduty_ruleset_rule":                                  resourcePagerDutyRulesetRule(),
			"pagerduty_business_service":                              resourcePagerDutyBusinessService(),
			"pagerduty_service_dependency":                            resourcePagerDutyServiceDependency(),
			"pagerduty_response_play":                                 resourcePagerDutyResponsePlay(),
			"pagerduty_tag":                                           resourcePagerDutyTag(),
			"pagerduty_tag_assignment":                                resourcePagerDutyTagAssignment(),
			"pagerduty_service_event_rule":                            resourcePagerDutyServiceEventRule(),
			"pagerduty_slack_connection":                              resourcePagerDutySlackConnection(),
			"pagerduty_business_service_subscriber":                   resourcePagerDutyBusinessServiceSubscriber(),
			"pagerduty_webhook_subscription":                          resourcePagerDutyWebhookSubscription(),
			"pagerduty_event_orchestration":                           resourcePagerDutyEventOrchestration(),
			"pagerduty_event_orchestration_router":                    resourcePagerDutyEventOrchestrationPathRouter(),
			"pagerduty_event_orchestration_unrouted":                  resourcePagerDutyEventOrchestrationPathUnrouted(),
			"pagerduty_event_orchestration_service":                   resourcePagerDutyEventOrchestrationPathService(),
			"pagerduty_automation_actions_runner":                     resourcePagerDutyAutomationActionsRunner(),
			"pagerduty_automation_actions_action":                     resourcePagerDutyAutomationActionsAction(),
			"pagerduty_automation_actions_action_team_association":    resourcePagerDutyAutomationActionsActionTeamAssociation(),
			"pagerduty_automation_actions_runner_team_association":    resourcePagerDutyAutomationActionsRunnerTeamAssociation(),
			"pagerduty_incident_workflow":                             resourcePagerDutyIncidentWorkflow(),
			"pagerduty_incident_workflow_trigger":                     resourcePagerDutyIncidentWorkflowTrigger(),
			"pagerduty_automation_actions_action_service_association": resourcePagerDutyAutomationActionsActionServiceAssociation(),
			"pagerduty_custom_field":                                  resourcePagerDutyCustomField(),
			"pagerduty_custom_field_option":                           resourcePagerDutyCustomFieldOption(),
			"pagerduty_custom_field_schema":                           resourcePagerDutyCustomFieldSchema(),
			"pagerduty_custom_field_schema_field_configuration":       resourcePagerDutyCustomFieldSchemaFieldConfiguration(),
			"pagerduty_custom_field_schema_assignment":                resourcePagerDutyCustomFieldSchemaAssignment(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		var ServiceRegion = strings.ToLower(d.Get("service_region").(string))

		if ServiceRegion == "us" || ServiceRegion == "" {
			ServiceRegion = ""
		} else {
			ServiceRegion = ServiceRegion + "."
		}

		config := Config{
			ApiUrl:              "https://api." + ServiceRegion + "pagerduty.com",
			AppUrl:              "https://app." + ServiceRegion + "pagerduty.com",
			SkipCredsValidation: d.Get("skip_credentials_validation").(bool),
			Token:               d.Get("token").(string),
			UserToken:           d.Get("user_token").(string),
			UserAgent:           p.UserAgent("terraform-provider-pagerduty", version),
			ApiUrlOverride:      d.Get("api_url_override").(string),
		}

		log.Println("[INFO] Initializing PagerDuty client")
		return &config, nil
	}

	return p
}

func isErrCode(err error, code int) bool {
	if e, ok := err.(*pagerduty.Error); ok && e.ErrorResponse.Response.StatusCode == code {
		return true
	}

	return false
}

func genError(err error, d *schema.ResourceData) error {
	return fmt.Errorf("Error reading: %s: %s", d.Id(), err)
}

func handleNotFoundError(err error, d *schema.ResourceData) error {
	if isErrCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return nil
	}
	return genError(err, d)
}
