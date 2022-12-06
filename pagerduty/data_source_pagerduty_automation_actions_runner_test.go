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
	runnerId := "default"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyAutomationActionsRunnerConfig_1(username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerdutyAutomationActionsRunnerRunnerCreation_1(runnerName, &runnerId),
				),
			},
			{
				Config: testAccDataSourcePagerDutyAutomationActionsRunnerConfig(runnerName, &runnerId),
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

func testAccDataSourcePagerDutyAutomationActionsRunnerConfig(runnerName string, runnerId *string) string {

	fmt.Println("testAccDataSourcePagerDutyAutomationActionsRunnerConfig runnerId:", *runnerId)

	return fmt.Sprintf(`
data "pagerduty_automation_actions_runner" "tf_test_runner" {
  id = %q
}
`, *runnerId)
}

func testAccDataSourcePagerDutyAutomationActionsRunnerConfig_1(username, email string) string {

	fmt.Println("Called: testAccDataSourcePagerDutyAutomationActionsRunnerConfig_1")

	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}
`, username, email)
}

func testAccDataSourcePagerdutyAutomationActionsRunnerRunnerCreation_1(runnerName string, runnerId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		fmt.Println("testAccDataSourcePagerdutyAutomationActionsRunnerRunnerCreation_1 runnerId:", *runnerId)

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

		fmt.Println("POST testAccDataSourcePagerdutyAutomationActionsRunnerRunnerCreation_1 runnerId:", *runnerId)

		return nil
	}
}
