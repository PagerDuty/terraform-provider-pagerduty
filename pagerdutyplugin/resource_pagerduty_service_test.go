package pagerduty

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_service", &resource.Sweeper{
		Name: "pagerduty_service",
		F:    testSweepService,
	})
}

func testSweepService(region string) error {
	ctx := context.Background()
	resp, err := testAccProvider.client.ListServicesWithContext(ctx, pagerduty.ListServiceOptions{})
	if err != nil {
		return err
	}

	for _, service := range resp.Services {
		if strings.HasPrefix(service.Name, "test") || strings.HasPrefix(service.Name, "tf-") {
			log.Printf("Destroying service %s (%s)", service.Name, service.ID)
			if err := testAccProvider.client.DeleteServiceWithContext(ctx, service.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyService_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_timeout"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_service.foo", "html_url"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "type", "service"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigUpdated(username, email, escalationPolicy, serviceUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", serviceUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigUpdatedWithDisabledTimeouts(username, email, escalationPolicy, serviceUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", serviceUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "null"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "null"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_FormatValidation(t *testing.T) {
	service := fmt.Sprintf("ts-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	errMessageMatcher := "Name can not be blank, nor contain non-printable characters. Trailing white spaces are not allowed either."

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			// Just a valid name
			{
				Config:             testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, "DB Technical Service"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Blank Name
			{
				Config:      testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with one white space at the end
			{
				Config:      testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, "this name has a white space at the end "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with multiple white space at the end
			{
				Config:      testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, "this name has white spaces at the end    "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with non printable characters
			{
				Config:      testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, "this name has a non printable\\n character"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Alert grouping parameters "Content Based" type input validation
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "content_based"
            config {}
          }
          `,
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("When using Alert grouping parameters configuration of type \"content_based\" is in use, attributes \"aggregate\" and \"fields\" are required"),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "content_based"
            config {
              aggregate = "all"
              fields    = ["custom_details.source_id"]
            }
          }
          `,
				),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "time"
            config {
              aggregate = "all"
              fields    = ["custom_details.source_id"]
            }
          }
          `,
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Alert grouping parameters configuration attributes \"aggregate\" and \"fields\" are only supported by \"content_based\" type Alert Grouping"),
			},
			// Alert grouping parameters "time" type input validation
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "time"
            config {
              timeout = 5
            }
          }
          `,
				),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "intelligent"
            config {
              timeout = 5
            }
          }
          `,
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Alert grouping parameters configuration attribute \"timeout\" is only supported by \"time\" type Alert Grouping"),
			},
			// Alert grouping parameters "intelligent" type input validation
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "time"
            config {
              time_window = 600
            }
          }
          `,
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Alert grouping parameters configuration attribute \"time_window\" is only supported by \"intelligent\" and \"content-based\" type Alert Grouping"),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "intelligent"
            config {}
          }
          `,
				),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "intelligent"
            config {
              time_window = 5
            }
          }
          `,
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Alert grouping time window value must be between 300 and 3600"),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "intelligent"
            config {
              time_window = 300
            }
          }
          `,
				),
				PlanOnly: true,
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "content_based"
            config {
              time_window = 5
            }
          }
          `,
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Alert grouping time window value must be between 300 and 3600"),
			},
			{
				Config: testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service,
					`
          alert_grouping_parameters {
            type = "content_based"
            config {
              aggregate = "all"
              fields    = ["custom_details.source_id"]
              time_window = 300
            }
          }
          `,
				),
			},
		},
	})
}

func TestAccPagerDutyService_AlertGrouping(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccPreCheckPagerDutyAbility(t, "preview_intelligent_alert_grouping") },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertGrouping(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping", "time"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertGroupingUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping", "intelligent"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_timeout", "1900"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_AlertContentGrouping(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertContentGrouping(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping", "rules"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0.aggregate", "all"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "content_based"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0.fields.0", "custom_details.field1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config:   testAccCheckPagerDutyServiceConfigWithAlertContentGrouping(username, email, escalationPolicy, service),
				PlanOnly: true,
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "intelligent"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertContentGroupingUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertTimeGroupingUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "time"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0.timeout", "5"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertTimeGroupingTimeoutZeroUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "time"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0.timeout", "0"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "intelligent"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingDescriptionUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "intelligent"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingOmittingConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "intelligent"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingTypeNullEmptyConfigConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_AlertContentGroupingIntelligentTimeWindow(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertContentGroupingIntelligentTimeWindow(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "intelligent"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAlertContentGroupingIntelligentTimeWindowUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.type", "intelligent"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_grouping_parameters.0.config.0.time_window", "900"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_AutoPauseNotificationsParameters(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfigWithAutoPauseNotificationsParameters(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.0.timeout", "300"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAutoPauseNotificationsParametersUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.0.enabled", "false"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.0.timeout", "0"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceConfigWithAutoPauseNotificationsParametersRemoved(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.0.enabled", "false"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_pause_notifications_parameters.0.timeout", "0"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_BasicWithIncidentUrgencyRules(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.urgency", "low"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "use_support_hours"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.at.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.at.0.name", "support_hours_start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.to_urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.type", "urgency_change"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.#", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.0", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.1", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.2", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.3", "4"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.4", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.end_time", "17:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.time_zone", "America/Lima"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.type", "fixed_time_per_day"),
				),
			},
			{
				Config:      testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfigError(username, email, escalationPolicy, serviceUpdated),
				ExpectError: regexp.MustCompile("general urgency cannot be set for a use_support_hours incident urgency rule type"),
			},
			{
				Config: testAccCheckPagerDutyServiceWithIncidentUrgencyRulesWithoutScheduledActionsConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.urgency", "severity_based"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "use_support_hours"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.#", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.0", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.1", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.2", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.3", "4"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.4", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.end_time", "17:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.time_zone", "America/Lima"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.type", "fixed_time_per_day"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfigUpdated(username, email, escalationPolicy, serviceUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", serviceUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "bar bar bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.urgency", "low"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "use_support_hours"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.at.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.at.0.name", "support_hours_start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.to_urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.type", "urgency_change"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.#", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.0", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.1", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.2", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.3", "4"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.4", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.end_time", "17:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.time_zone", "America/Lima"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.type", "fixed_time_per_day"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_FromBasicToCustomIncidentUrgencyRules(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfigUpdated(username, email, escalationPolicy, serviceUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", serviceUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "bar bar bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "3600"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.during_support_hours.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.type", "constant"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.outside_support_hours.0.urgency", "low"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "use_support_hours"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.at.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.at.0.name", "support_hours_start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.to_urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "scheduled_actions.0.type", "urgency_change"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.#", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.0", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.1", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.2", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.3", "4"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.days_of_week.4", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.end_time", "17:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.time_zone", "America/Lima"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "support_hours.0.type", "fixed_time_per_day"),
				),
			},
		},
	})
}

func TestAccPagerDutyService_SupportHoursChange(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service_id := ""
	p_service_id := &service_id
	updated_service_id := ""
	p_updated_service_id := &updated_service_id

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfig(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					testAccCheckPagerDutyServiceSaveServiceId(p_service_id, "pagerduty_service.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceWithSupportHoursConfigUpdated(username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					testAccCheckPagerDutyServiceSaveServiceId(p_updated_service_id, "pagerduty_service.foo"),
				),
			},
		},
	})

	if service_id != updated_service_id {
		t.Error(fmt.Errorf("Expected service id to be %s, but found %s", service_id, updated_service_id))
	}
}

func TestAccPagerDutyService_ResponsePlay(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	responsePlay := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceWithResponsePlayConfig(username, email, escalationPolicy, responsePlay, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_timeout"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_service.foo", "html_url"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "type", "service"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_service.foo", "response_play"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceWithNullResponsePlayConfig(username, email, escalationPolicy, responsePlay, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceExists("pagerduty_service.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "auto_resolve_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "acknowledgement_timeout", "1800"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "alert_creation", "create_alerts_and_incidents"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service.foo", "alert_grouping_timeout"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.urgency", "high"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "incident_urgency_rule.0.type", "constant"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_service.foo", "html_url"),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "type", "service"),
					testAccCheckPagerDutyServiceResponsePlayNotExist("pagerduty_service.foo"),
				),
			},
		},
	})

}

func testAccCheckPagerDutyServiceSaveServiceId(p *string, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		ctx := context.Background()
		found, err := testAccProvider.client.GetServiceWithContext(ctx, rs.Primary.ID, &pagerduty.GetServiceOptions{})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Service not found: %v - %v", rs.Primary.ID, found)
		}

		*p = found.ID

		return nil
	}
}

func testAccCheckPagerDutyServiceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pagerduty_service" {
			continue
		}

		ctx := context.Background()
		_, err := testAccProvider.client.GetServiceWithContext(ctx, rs.Primary.ID, &pagerduty.GetServiceOptions{})
		if err == nil {
			return fmt.Errorf("Service still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		ctx := context.Background()
		found, err := testAccProvider.client.GetServiceWithContext(ctx, rs.Primary.ID, &pagerduty.GetServiceOptions{})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Service not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyServiceResponsePlayNotExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		ctx := context.Background()
		found, err := testAccProvider.client.GetServiceWithContext(ctx, rs.Primary.ID, &pagerduty.GetServiceOptions{})
		if err != nil {
			return err
		}

		if found.ID == rs.Primary.ID && found.ResponsePlay != nil {
			return fmt.Errorf("Service %s still has a response play configured", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceAlertGroupingInputValidationConfig(username, email, escalationPolicy, service, alertGroupingParams string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "bar"
  num_loops   = 2
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_service" "foo" {
  name                    = "%s"
  description             = "foo"
  auto_resolve_timeout    = 1800
  acknowledgement_timeout = 1800
  escalation_policy       = pagerduty_escalation_policy.foo.id
  alert_creation          = "create_alerts_and_incidents"
  %s
}
`, username, email, escalationPolicy, service, alertGroupingParams)
}

func testAccCheckPagerDutyServiceConfigWithAlertGrouping(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping          = "time"
	alert_grouping_timeout  = 1800
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertContentGrouping(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
        type = "content_based"
        config {
            aggregate = "all"
            fields = ["custom_details.field1"]
        }
    }
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertContentGroupingIntelligentTimeWindow(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
        type = "intelligent"
    }
}
`, username, email, escalationPolicy, service)
}
func testAccCheckPagerDutyServiceConfigWithAlertContentGroupingIntelligentTimeWindowUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
        type = "intelligent"
        config {
            time_window = 900
        }
    }
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertContentGroupingUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
        type = null
    }
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertTimeGroupingUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
        type = "time"
        config {
          timeout = 5
        }
    }
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertTimeGroupingTimeoutZeroUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
		type = "time"
		config {
			timeout = 0
		}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertGroupingUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping          = "intelligent"
	alert_grouping_timeout  = 1900
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
		type = "intelligent"
		config {
			fields = null
			timeout = 0
		}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingDescriptionUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
		type = "intelligent"
		config {}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingOmittingConfig(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
		type = "intelligent"
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAlertIntelligentGroupingTypeNullEmptyConfigConfig(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	alert_grouping_parameters {
		type = null
		config {}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAutoPauseNotificationsParameters(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	auto_pause_notifications_parameters {
		enabled = true
		timeout = 300
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAutoPauseNotificationsParametersUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
	auto_pause_notifications_parameters {
		enabled = false
		timeout = null
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigWithAutoPauseNotificationsParametersRemoved(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_alerts_and_incidents"
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "bar"
	auto_resolve_timeout    = 3600
	acknowledgement_timeout = 3600

	escalation_policy       = pagerduty_escalation_policy.foo.id
	incident_urgency_rule {
		type    = "constant"
		urgency = "high"
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceConfigUpdatedWithDisabledTimeouts(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "bar"
	auto_resolve_timeout    = "null"
	acknowledgement_timeout = "null"

	escalation_policy       = pagerduty_escalation_policy.foo.id
	incident_urgency_rule {
		type    = "constant"
		urgency = "high"
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfig(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "use_support_hours"

		during_support_hours {
			type    = "constant"
			urgency = "high"
		}
		outside_support_hours {
			type    = "constant"
			urgency = "low"
		}
	}

	support_hours {
		type         = "fixed_time_per_day"
		time_zone    = "America/Lima"
		start_time   = "09:00:00"
		end_time     = "17:00:00"
		days_of_week = [ 1, 2, 3, 4, 5 ]
	}

	scheduled_actions {
		type = "urgency_change"
		to_urgency = "high"
		at {
			type = "named_time"
			name = "support_hours_start"
		}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfigError(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type    = "use_support_hours"
		urgency = "high"
		during_support_hours {
			type    = "constant"
			urgency = "high"
		}
		outside_support_hours {
			type    = "constant"
			urgency = "low"
		}
	}

	support_hours {
		type         = "fixed_time_per_day"
		time_zone    = "America/Lima"
		start_time   = "09:00:00"
		end_time     = "17:00:00"
		days_of_week = [ 1, 2, 3, 4, 5 ]
	}

	scheduled_actions {
		type = "urgency_change"
		to_urgency = "high"
		at {
			type = "named_time"
			name = "support_hours_start"
		}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceWithIncidentUrgencyRulesWithoutScheduledActionsConfig(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "use_support_hours"

		during_support_hours {
			type    = "constant"
			urgency = "high"
		}
		outside_support_hours {
			type    = "constant"
			urgency = "severity_based"
		}
	}

	support_hours {
		type         = "fixed_time_per_day"
		time_zone    = "America/Lima"
		start_time   = "09:00:00"
		end_time     = "17:00:00"
		days_of_week = [ 1, 2, 3, 4, 5 ]
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceWithIncidentUrgencyRulesConfigUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
	resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "bar bar bar"
	auto_resolve_timeout    = 3600
	acknowledgement_timeout = 3600
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "use_support_hours"
		during_support_hours {
			type    = "constant"
			urgency = "high"
		}
		outside_support_hours {
			type    = "constant"
			urgency = "low"
		}
	}

	support_hours {
		type         = "fixed_time_per_day"
		time_zone    = "America/Lima"
		start_time   = "09:00:00"
		end_time     = "17:00:00"
		days_of_week = [ 1, 2, 3, 4, 5 ]
	}

	scheduled_actions {
		type = "urgency_change"
		to_urgency = "high"
		at {
			type = "named_time"
			name = "support_hours_start"
		}
	}
}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceWithSupportHoursConfigUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "constant"
		urgency = "high"
	}

}
`, username, email, escalationPolicy, service)
}

func testAccCheckPagerDutyServiceWithResponsePlayConfig(username, email, escalationPolicy, responsePlay, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "bar"
  num_loops   = 2
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_response_play" "foo" {
  name = "%s"
  from = pagerduty_user.foo.email

  responder {
    type = "escalation_policy_reference"
    id   = pagerduty_escalation_policy.foo.id
  }

  subscriber {
    type = "user_reference"
    id   = pagerduty_user.foo.id
  }

  runnability = "services"
}

resource "pagerduty_service" "foo" {
  name                    = "%s"
  description             = "foo"
  auto_resolve_timeout    = 1800
  acknowledgement_timeout = 1800
  escalation_policy       = pagerduty_escalation_policy.foo.id
  response_play           = pagerduty_response_play.foo.id
}
`, username, email, escalationPolicy, responsePlay, service)
}

func testAccCheckPagerDutyServiceWithNullResponsePlayConfig(username, email, escalationPolicy, responsePlay, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "bar"
  num_loops   = 2
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_response_play" "foo" {
  name = "%s"
  from = pagerduty_user.foo.email

  responder {
    type = "escalation_policy_reference"
    id   = pagerduty_escalation_policy.foo.id
  }

  subscriber {
    type = "user_reference"
    id   = pagerduty_user.foo.id
  }

  runnability = "services"
}

resource "pagerduty_service" "foo" {
  name                    = "%s"
  description             = "foo"
  auto_resolve_timeout    = 1800
  acknowledgement_timeout = 1800
  escalation_policy       = pagerduty_escalation_policy.foo.id
  response_play           = null
}
`, username, email, escalationPolicy, responsePlay, service)
}

func TestFlattenAlertGroupingParameters_Basic(t *testing.T) {
	var diags diag.Diagnostics
	var timeout uint = 1000
	params := &pagerduty.AlertGroupingParameters{
		Type: "foo",
		Config: &pagerduty.AlertGroupParamsConfig{
			Timeout:   &timeout,
			Aggregate: "aggregate",
			Fields:    []string{"a", "b", "c"},
		},
	}
	li := flattenAlertGroupingParameters(context.Background(), params, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected error %s", diags)
	}
	t.Logf("%#v", li)
}
