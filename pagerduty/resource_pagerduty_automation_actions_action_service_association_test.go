package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_automation_actions_action_service_association", &resource.Sweeper{
		Name: "pagerduty_automation_actions_action_service_association",
		F:    testSweepAutomationActionsActionServiceAssociation,
	})
}

func testSweepAutomationActionsActionServiceAssociation(region string) error {
	return nil
}

func TestAccPagerDutyAutomationActionsActionServiceAssociation_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	actionName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionServiceAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionServiceAssociationConfig(username, email, escalationPolicy, serviceName, actionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAutomationActionsActionServiceAssociationExists("pagerduty_automation_actions_action_service_association.foo"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_action_service_association.foo", "action_id"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_action_service_association.foo", "service_id"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAutomationActionsActionServiceAssociationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_automation_actions_action_service_association" {
			continue
		}
		actionID, serviceID, err := resourcePagerDutyParseColonCompoundID(r.Primary.ID)
		if err != nil {
			return err
		}

		if _, _, err := client.AutomationActionsAction.GetAssociationToService(actionID, serviceID); err == nil {
			return fmt.Errorf("Automation Actions Runner still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyAutomationActionsActionServiceAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Automation Actions Runner ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		actionID, serviceID, err := resourcePagerDutyParseColonCompoundID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, _, err := client.AutomationActionsAction.GetAssociationToService(actionID, serviceID)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%s:%s", actionID, found.Service.ID) != rs.Primary.ID {
			return fmt.Errorf("Automation Actions Action association to service not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAutomationActionsActionServiceAssociationConfig(username, email, escalationPolicy, serviceName, actionName string) string {
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
	alert_creation          = "create_incidents"
}

resource "pagerduty_automation_actions_action" "foo" {
	name = "%s"
	description = "PA Action created by TF"
	action_type = "script"
	action_data_reference {
		script = "java --version"
		invocation_command = "/bin/bash"
	  }
}

resource "pagerduty_automation_actions_action_service_association" "foo" {
  action_id = pagerduty_automation_actions_action.foo.id
  service_id = pagerduty_service.foo.id
}

`, username, email, escalationPolicy, serviceName, actionName)
}
