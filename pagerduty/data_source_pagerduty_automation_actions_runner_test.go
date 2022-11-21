package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyAutomationActionsRunner_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyAutomationActionsRunnerConfig("01DCTHG8L7X4BDEQG3OQO2HZCN"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerdutyAutomationActionsRunnerExists("data.pagerduty_automation_actions_runner.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerdutyAutomationActionsRunnerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// client, _ := testAccProvider.Meta().(*Config).Client()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Runner ID is set")
		}

		return nil
	}
}

func testAccDataSourcePagerDutyAutomationActionsRunnerConfig(id string) string {
	return fmt.Sprintf(`
data "pagerduty_automation_actions_runner" "foo" {
  id = "%s"
}
`, id)
}
