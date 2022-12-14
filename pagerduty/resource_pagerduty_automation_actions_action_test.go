package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("automation_actions_action", &resource.Sweeper{
		Name: "automation_actions_action",
		F:    testSweepAutomationActionsAction,
	})
}

func testSweepAutomationActionsAction(region string) error {
	return nil
}

func TestAccPagerDutyAutomationActionsAction_Basic(t *testing.T) {
	actionName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionConfig(actionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAutomationActionsActionExists("pagerduty_automation_actions_action.foo"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_action.foo", "name", actionName),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_action.foo", "action_type", "process_automation"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_action.foo", "description", "PA Action created by TF"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_action.foo", "type", "action"),
					resource.TestCheckResourceAttr(
						"pagerduty_automation_actions_action.foo", "action_data_reference.0.process_automation_job_id", "pa_job_id_123"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_action.foo", "id"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_action.foo", "creation_time"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAutomationActionsActionDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_automation_actions_action" {
			continue
		}
		if _, _, err := client.AutomationActionsAction.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Automation Actions Action still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyAutomationActionsActionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Automation Actions Action ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.AutomationActionsAction.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Automation Actions Action not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAutomationActionsActionConfig(actionName string) string {
	return fmt.Sprintf(`
resource "pagerduty_automation_actions_action" "foo" {
	name = "%s"
	description = "PA Action created by TF"
	action_type = "process_automation"
	action_data_reference {
		process_automation_job_id = "pa_job_id_123"
	  }
}
`, actionName)
}
