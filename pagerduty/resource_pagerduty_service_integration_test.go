package pagerduty

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyServiceIntegration_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegrationUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyServiceIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceIntegrationConfig(username, email, escalationPolicy, service, serviceIntegration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceIntegrationExists("pagerduty_service_integration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "name", serviceIntegration),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_events_api_inbound_integration"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "vendor", "PAM4FGS"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceIntegrationConfigUpdated(username, email, escalationPolicy, service, serviceIntegrationUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceIntegrationExists("pagerduty_service_integration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "name", serviceIntegrationUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_events_api_inbound_integration"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "vendor", "PAM4FGS"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_service_integration.foo", "html_url"),
				),
			},
		},
	})
}

func TestAccPagerDutyServiceIntegrationGeneric_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegrationUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyServiceIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceIntegrationGenericConfig(username, email, escalationPolicy, service, serviceIntegration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceIntegrationExists("pagerduty_service_integration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "name", serviceIntegration),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_events_api_inbound_integration"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceIntegrationGenericConfigUpdated(username, email, escalationPolicy, service, serviceIntegrationUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceIntegrationExists("pagerduty_service_integration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "name", serviceIntegrationUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_events_api_inbound_integration"),
				),
			},
			{
				Config:      testAccCheckPagerDutyServiceIntegrationGenericEmail(username, email, escalationPolicy, service, serviceIntegration, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("integration_email attribute must be set for an integration type generic_email_inbound_integration"),
			},
			{
				Config: testAccCheckPagerDutyServiceIntegrationGenericEmail(username, email, escalationPolicy, service, serviceIntegration, "user@pagerduty.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_email_inbound_integration"),
				),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccCheckPagerDutyServiceIntegrationGenericEmailWithKnownAfterApply(username, email, escalationPolicy, service, serviceIntegration, "pagerduty.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_email_inbound_integration"),
				),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccPagerDutyServiceIntegrationEmail_Filters(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegrationUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyServiceIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceIntegrationEmailFiltersConfig(username, email, escalationPolicy, service, serviceIntegration, testAccGetPagerDutyAccountDomain(t)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceIntegrationExists("pagerduty_service_integration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "name", serviceIntegration),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "integration_email", fmt.Sprintf("s1@%s", testAccGetPagerDutyAccountDomain(t))),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_incident_creation", "use_rules"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter_mode", "and-rules-email"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parsing_fallback", "open_new_incident"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.body_mode", "always"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.from_email_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.from_email_regex", "(@foo.test*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.subject_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.subject_regex", "(CRITICAL*)"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.body_mode", "always"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.from_email_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.from_email_regex", "(@bar.com*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.subject_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.subject_regex", "(CRITICAL*)"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.action", "resolve"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.type", "any"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.0.matcher", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.0.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.0.type", "contains"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.type", "not"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.predicate.0.matcher", "(bar*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.predicate.0.part", "body"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.predicate.0.type", "regex"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.value_name", "incident_key"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.value_name", "FieldName1"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.action", "trigger"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.type", "all"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.type", "not"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.predicate.0.matcher", "(foo*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.predicate.0.part", "from_addresses"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.predicate.0.type", "exactly"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.1.matcher", "Bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.1.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.1.type", "contains"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.part", "body"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.value_name", "incident_key"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.value_name", "FieldName1"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceIntegrationEmailFiltersConfigUpdated(username, email, escalationPolicy, service, serviceIntegrationUpdated, testAccGetPagerDutyAccountDomain(t)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceIntegrationExists("pagerduty_service_integration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "name", serviceIntegrationUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "integration_email", fmt.Sprintf("s11@%s", testAccGetPagerDutyAccountDomain(t))),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_incident_creation", "use_rules"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter_mode", "and-rules-email"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parsing_fallback", "open_new_incident"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.body_mode", "always"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.from_email_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.from_email_regex", "(@foo.test*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.subject_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.0.subject_regex", "(CRITICAL*)"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.body_mode", "always"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.from_email_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.from_email_regex", "(@bar.com*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.subject_mode", "match"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_filter.1.subject_regex", "(CRITICAL*)"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.action", "resolve"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.type", "any"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.0.matcher", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.0.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.0.type", "contains"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.type", "not"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.predicate.0.matcher", "(bar*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.predicate.0.part", "body"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.match_predicate.0.predicate.1.predicate.0.type", "regex"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.0.value_name", "incident_key"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.0.value_extractor.1.value_name", "FieldName1"),

					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.action", "trigger"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.type", "all"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.type", "not"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.predicate.0.matcher", "(foo1*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.predicate.0.part", "from_addresses"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.0.predicate.0.type", "exactly"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.1.matcher", "Bar1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.1.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.match_predicate.0.predicate.1.type", "contains"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.part", "body"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.0.value_name", "incident_key"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.ends_before", "end"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.starts_after", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.type", "between"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.1.value_name", "FieldName11"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.2.part", "subject"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.2.regex", "(bar*)"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.2.type", "regex"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "email_parser.1.value_extractor.2.value_name", "FieldName2"),
				),
			},
		},
	})
}

func TestAccPagerDutyServiceIntegration_GenericEmailNoFilters(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckServiceIntegrationGenericEmailNoFilters(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyServiceIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceIntegrationGenericEmailNoFilters(username, email, escalationPolicy, service, serviceIntegration),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_email_inbound_integration"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.bar", "type", "generic_email_inbound_integration"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceIntegrationMultipleGenericEmailNoFilters(username, email, escalationPolicy, service, serviceIntegration),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_integration.foo", "type", "generic_email_inbound_integration"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyServiceIntegrationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_service_integration" {
			continue
		}

		service, _ := s.RootModule().Resources["pagerduty_service.foo"]

		if _, _, err := client.Services.GetIntegration(service.Primary.ID, r.Primary.ID, &pagerduty.GetIntegrationOptions{}); err == nil {
			return fmt.Errorf("Service Integration still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyServiceIntegrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Integration ID is set")
		}

		service, _ := s.RootModule().Resources["pagerduty_service.foo"]

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.Services.GetIntegration(service.Primary.ID, rs.Primary.ID, &pagerduty.GetIntegrationOptions{})
		if err != nil {
			return fmt.Errorf("Service integration not found: %v", rs.Primary.ID)
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Service Integration not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyServiceIntegrationConfig(username, email, escalationPolicy, service, serviceIntegration string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1

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

data "pagerduty_vendor" "datadog" {
  name = "datadog"
}

resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.foo.id
  vendor  = data.pagerduty_vendor.datadog.id
}
`, username, email, escalationPolicy, service, serviceIntegration)
}

func testAccCheckPagerDutyServiceIntegrationConfigUpdated(username, email, escalationPolicy, service, serviceIntegration string) string {
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

data "pagerduty_vendor" "datadog" {
  name = "datadog"
}

resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.foo.id
  vendor  = data.pagerduty_vendor.datadog.id
}
`, username, email, escalationPolicy, service, serviceIntegration)
}

func testAccCheckPagerDutyServiceIntegrationGenericConfig(username, email, escalationPolicy, service, serviceIntegration string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1

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

resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.foo.id
  type    = "generic_events_api_inbound_integration"
}
`, username, email, escalationPolicy, service, serviceIntegration)
}

func testAccCheckPagerDutyServiceIntegrationGenericConfigUpdated(username, email, escalationPolicy, service, serviceIntegration string) string {
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

resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.foo.id
  type    = "generic_events_api_inbound_integration"
}
`, username, email, escalationPolicy, service, serviceIntegration)
}

func testAccCheckPagerDutyServiceIntegrationGenericEmail(username, email, escalationPolicy, service, serviceIntegration, integrationEmail string) string {
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

resource "pagerduty_service_integration" "foo" {
  name              = "%s"
  service           = pagerduty_service.foo.id
  type              = "generic_email_inbound_integration"
  integration_email = "%s"
}
`, username, email, escalationPolicy, service, serviceIntegration, integrationEmail)
}

func testAccCheckPagerDutyServiceIntegrationGenericEmailWithKnownAfterApply(username, email, escalationPolicy, service, serviceIntegration, integrationEmailDomain string) string {
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

resource "pagerduty_service_integration" "foo" {
  name              = "%s"
  service           = pagerduty_service.foo.id
  type              = "generic_email_inbound_integration"
  # Instead of 'pagerduty_escalation_policy.foo.id', people could use
  # 'random_id.foo.hex' which has the same effect.
  integration_email = "service-${pagerduty_escalation_policy.foo.id}@%s"
}
`, username, email, escalationPolicy, service, serviceIntegration, integrationEmailDomain)
}

func testAccCheckPagerDutyServiceIntegrationEmailFiltersConfig(username, email, escalationPolicy, service, serviceIntegration string, accountDomain string) string {
	return fmt.Sprintf(`
data "pagerduty_vendor" "email" {
  name = "Email"
}
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}
resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1
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
resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.foo.id
  vendor  = data.pagerduty_vendor.email.id
  integration_email       = "s1@%s"
  email_incident_creation = "use_rules"
  email_filter_mode       = "and-rules-email"
  email_filter {
    body_mode        = "always"
    body_regex       = null
    from_email_mode  = "match"
    from_email_regex = "(@foo.test*)"
    subject_mode     = "match"
    subject_regex    = "(CRITICAL*)"
  }
  email_filter {
    body_mode        = "always"
    body_regex       = null
    from_email_mode  = "match"
    from_email_regex = "(@bar.com*)"
    subject_mode     = "match"
    subject_regex    = "(CRITICAL*)"
  }
  email_parser {
    action = "resolve"
    match_predicate {
      type = "any"
      predicate {
        matcher = "foo"
        part    = "subject"
        type    = "contains"
      }
      predicate {
        type = "not"
        predicate {
          matcher = "(bar*)"
          part    = "body"
          type    = "regex"
        }
      }
    }
    value_extractor {
      ends_before  = "end"
      part         = "subject"
      starts_after = "start"
      type         = "between"
      value_name   = "incident_key"
    }
    value_extractor {
      ends_before  = "end"
      part         = "subject"
      starts_after = "start"
      type         = "between"
      value_name   = "FieldName1"
    }
  }
  email_parser {
    action = "trigger"
    match_predicate {
      type = "all"
      predicate {
        type = "not"
        predicate {
          matcher = "(foo*)"
          part    = "from_addresses"
          type    = "exactly"
        }
      }
      predicate {
        matcher = "Bar"
        part    = "subject"
        type    = "contains"
      }
    }
    value_extractor {
      ends_before  = "end"
      part         = "body"
      starts_after = "start"
      type         = "between"
      value_name   = "incident_key"
    }
    value_extractor {
      ends_before  = "end"
      part         = "subject"
      starts_after = "start"
      type         = "between"
      value_name   = "FieldName1"
    }
  }
  email_parsing_fallback = "open_new_incident"
}
`, username, email, escalationPolicy, service, serviceIntegration, accountDomain)
}

func testAccCheckPagerDutyServiceIntegrationEmailFiltersConfigUpdated(username, email, escalationPolicy, service, serviceIntegration string, accountDomain string) string {
	return fmt.Sprintf(`
data "pagerduty_vendor" "email" {
  name = "Email"
}
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}
resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1
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
resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.foo.id
  vendor  = data.pagerduty_vendor.email.id
  integration_email       = "s11@%s"
  email_incident_creation = "use_rules"
  email_filter_mode       = "and-rules-email"
  email_filter {
    body_mode        = "always"
    body_regex       = null
    from_email_mode  = "match"
    from_email_regex = "(@foo.test*)"
    subject_mode     = "match"
    subject_regex    = "(CRITICAL*)"
  }
  email_filter {
    body_mode        = "always"
    body_regex       = null
    from_email_mode  = "match"
    from_email_regex = "(@bar.com*)"
    subject_mode     = "match"
    subject_regex    = "(CRITICAL*)"
  }
  email_parser {
    action = "resolve"
    match_predicate {
      type = "any"
      predicate {
        matcher = "foo"
        part    = "subject"
        type    = "contains"
      }
      predicate {
        type = "not"
        predicate {
          matcher = "(bar*)"
          part    = "body"
          type    = "regex"
        }
      }
    }
    value_extractor {
      ends_before  = "end"
      part         = "subject"
      starts_after = "start"
      type         = "between"
      value_name   = "incident_key"
    }
    value_extractor {
      ends_before  = "end"
      part         = "subject"
      starts_after = "start"
      type         = "between"
      value_name   = "FieldName1"
    }
  }
  email_parser {
    action = "trigger"
    match_predicate {
      type = "all"
      predicate {
        type = "not"
        predicate {
          matcher = "(foo1*)"
          part    = "from_addresses"
          type    = "exactly"
        }
      }
      predicate {
        matcher = "Bar1"
        part    = "subject"
        type    = "contains"
      }
    }
    value_extractor {
      ends_before  = "end"
      part         = "body"
      starts_after = "start"
      type         = "between"
      value_name   = "incident_key"
    }
    value_extractor {
      ends_before  = "end"
      part         = "subject"
      starts_after = "start"
      type         = "between"
      value_name   = "FieldName11"
    }
    value_extractor {
      part       = "subject"
      regex      = "(bar*)"
      type       = "regex"
      value_name = "FieldName2"
    }
  }
  email_parsing_fallback = "open_new_incident"
}
`, username, email, escalationPolicy, service, serviceIntegration, accountDomain)
}

func testAccCheckPagerDutyServiceIntegrationGenericEmailNoFilters(username, email, escalationPolicy, service, serviceIntegration string) string {
	integrationEmail := os.Getenv("PAGERDUTY_ACC_SERVICE_INTEGRATION_GENERIC_EMAIL_NO_FILTERS")

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

resource "pagerduty_service_integration" "foo" {
  name              = "%[5]s"
  service           = pagerduty_service.foo.id
  type              = "generic_email_inbound_integration"
  integration_email = "%[6]s"
}

resource "pagerduty_service_integration" "bar" {
  name              = "%[5]s"
  service           = pagerduty_service.foo.id
  type              = "generic_email_inbound_integration"
  integration_email = "%[6]s"
  email_filter {}
}
`, username, email, escalationPolicy, service, serviceIntegration, integrationEmail)
}

func testAccCheckPagerDutyServiceIntegrationMultipleGenericEmailNoFilters(username, email, escalationPolicy, service, serviceIntegration string) string {
	integrationEmail := os.Getenv("PAGERDUTY_ACC_SERVICE_INTEGRATION_GENERIC_EMAIL_NO_FILTERS")

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

resource "pagerduty_service_integration" "foo" {
  name              = "%[5]s"
  service           = pagerduty_service.foo.id
  type              = "generic_email_inbound_integration"
  integration_email = "%[6]s"
  email_filter {
    body_mode        = "always"
    body_regex       = null
    from_email_mode  = "match"
    from_email_regex = "(@bar.com*)"
    subject_mode     = "match"
    subject_regex    = "(CRITICAL*)"
  }
  email_filter {}
  email_filter {
    body_mode        = "always"
    body_regex       = null
    from_email_mode  = "match"
    from_email_regex = "(@bar.com*)"
    subject_mode     = "match"
    subject_regex    = "(CRITICAL*)"
  }
}
`, username, email, escalationPolicy, service, serviceIntegration, integrationEmail)
}

func testAccPreCheckServiceIntegrationGenericEmailNoFilters(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SERVICE_INTEGRATION_GENERIC_EMAIL_NO_FILTERS"); v == "" {
		t.Skip("PAGERDUTY_ACC_SERVICE_INTEGRATION_GENERIC_EMAIL_NO_FILTERS not set. Skipping Service Integration related test")
	}
}
