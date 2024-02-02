package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyAutomationActionsRunner_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyAutomationActionsRunnerConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerdutyAutomationActionsRunner("pagerduty_automation_actions_runner.test", "data.pagerduty_automation_actions_runner.foo"),
				),
			},
		},
	})
}

func testAccDataSourcePagerdutyAutomationActionsRunner(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		ds, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		dsA := ds.Primary.Attributes

		if dsA["id"] == "" {
			return fmt.Errorf("No Runner ID is set")
		}

		testAtts := []string{"id", "name", "type", "runner_type", "creation_time", "last_seen", "description", "runbook_base_uri"}

		for _, att := range testAtts {
			if dsA[att] != srcA[att] {
				return fmt.Errorf("Expected the runner %s to be: %s, but got: %s", att, srcA[att], dsA[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyAutomationActionsRunnerConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_automation_actions_runner" "test" {
  name = "%s"
  description = "Runner created by TF"
  runner_type = "runbook"
  runbook_base_uri = "cat-cat"
  runbook_api_key = "secret"
}

data "pagerduty_automation_actions_runner" "foo" {
  id = pagerduty_automation_actions_runner.test.id
}
`, name)
}
