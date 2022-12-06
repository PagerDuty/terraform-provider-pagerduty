package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccDataSourcePagerDutyAutomationActionsRunner_Basic(t *testing.T) {
	runnerName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	runnerId := "default"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service),
				// Config: `
				// provider "pagerduty" {}
				// `,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerdutyAutomationActionsRunnerRunnerCreation(runnerName, &runnerId),
				),
			},
			{
				Config: testAccDataSourcePagerDutyAutomationActionsRunnerConfig(runnerName, runnerId),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerdutyAutomationActionsRunner("data.pagerduty_automation_actions_runner.foo"),
				),
			},
		},
	})
}

func testAccDataSourcePagerdutyAutomationActionsRunner(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

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

func testAccDataSourcePagerDutyAutomationActionsRunnerConfig(runnerName, runnerId string) string {
	return fmt.Sprintf(`
data "pagerduty_automation_actions_runner" "tf_test_runner" {
  id = %q
}
`, runnerId)
}

func testAccDataSourcePagerdutyAutomationActionsRunnerRunnerCreation(runnerName string, runnerId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		fmt.Println("runnerId in runner creation:", *runnerId)
		client, _ := testAccProvider.Meta().(*Config).Client()
		input := &pagerduty.AutomationActionsRunner{
			Name:        runnerName,
			Description: "sidecar runner provisioned for tf acceptance test",
			RunnerType:  "sidecar",
		}
		runner, _, err := client.AutomationActionsRunner.Create(input)
		if err != nil {
			return err
		}
		*runnerId = runner.ID
		fmt.Println("runnerId", runnerId)
		fmt.Println("runnerId in runner creation2:", *runnerId)
		return nil
	}
}
