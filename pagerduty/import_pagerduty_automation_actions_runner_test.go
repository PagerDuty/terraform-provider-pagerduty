package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyAutomationActionsRunner_import(t *testing.T) {
	runnerName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsRunnerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsRunnerConfig(runnerName),
			},

			{
				Config: testAccCheckPagerDutyAutomationActionsRunnerConfig2(runnerName),
			},

			{
				ResourceName:      "pagerduty_automation_actions_runner.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyAutomationActionsRunnerConfig2(runnerName string) string {
	return fmt.Sprintf(`
resource "pagerduty_automation_actions_runner" "foo" {
	name = "%s"
	description = "Runner created by TF"
	runner_type = "runbook"
	runbook_base_uri = "cat-cat"
    # runbook_api_key = "cat-secret"
}
`, runnerName)
}
