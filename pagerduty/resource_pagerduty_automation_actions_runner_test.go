package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("automation_actions_runner", &resource.Sweeper{
		Name: "automation_actions_runner",
		F:    testSweepAutomationActionsRunner,
	})
}

func testSweepAutomationActionsRunner(region string) error {
	return nil
}

func TestAccPagerDutyAutomationActionsRunner_Basic(t *testing.T) {
	runnerName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	nameUpdated := fmt.Sprintf("%s-updated", runnerName)
	descriptionUpdated := fmt.Sprintf("Description of %s-updated", runnerName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsRunnerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsRunnerConfig(runnerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAutomationActionsRunnerExists("pagerduty_automation_actions_runner.foo"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "name", runnerName),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "runner_type", "runbook"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "description", "Runner created by TF"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "runbook_base_uri", "cat-cat"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "type", "runner"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_runner.foo", "id"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_runner.foo", "creation_time"),
				),
			},
			{
				Config: testAccCheckPagerDutyAutomationActionsRunnerUpdatedConfig(nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAutomationActionsRunnerExists("pagerduty_automation_actions_runner.foo"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "name", nameUpdated),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "runner_type", "runbook"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "description", descriptionUpdated),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "runbook_base_uri", "cat-cat-updated"),
					resource.TestCheckResourceAttr("pagerduty_automation_actions_runner.foo", "type", "runner"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_runner.foo", "id"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_runner.foo", "creation_time"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAutomationActionsRunnerDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_automation_actions_runner" {
			continue
		}
		if _, _, err := client.AutomationActionsRunner.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Automation Actions Runner still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyAutomationActionsRunnerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Automation Actions Runner ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.AutomationActionsRunner.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Automation Actions Runner not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAutomationActionsRunnerConfig(runnerName string) string {
	return fmt.Sprintf(`
resource "pagerduty_automation_actions_runner" "foo" {
	name = "%s"
	description = "Runner created by TF"
	runner_type = "runbook"
	runbook_base_uri = "cat-cat"
	runbook_api_key = "cat-secret"
}
`, runnerName)
}

func testAccCheckPagerDutyAutomationActionsRunnerUpdatedConfig(runnerName, runnerDescription string) string {
	return fmt.Sprintf(`
resource "pagerduty_automation_actions_runner" "foo" {
	name = "%s"
	description = "%s"
	runner_type = "runbook"
	runbook_base_uri = "cat-cat-updated"
	runbook_api_key = "cat-secret-updated"
}
`, runnerName, runnerDescription)
}
