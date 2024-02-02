package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_automation_actions_runner_team_association", &resource.Sweeper{
		Name: "pagerduty_automation_actions_runner_team_association",
		F:    testSweepAutomationActionsRunnerTeamAssociation,
	})
}

func testSweepAutomationActionsRunnerTeamAssociation(region string) error {
	return nil
}

func TestAccPagerDutyAutomationActionsRunnerTeamAssociation_Basic(t *testing.T) {
	runnerName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationConfig(runnerName, teamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationExists("pagerduty_automation_actions_runner_team_association.foo"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_runner_team_association.foo", "runner_id"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_runner_team_association.foo", "team_id"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_automation_actions_runner_team_association" {
			continue
		}
		runnerID, teamID, err := resourcePagerDutyParseColonCompoundID(r.Primary.ID)
		if err != nil {
			return err
		}

		if _, _, err := client.AutomationActionsRunner.GetAssociationToTeam(runnerID, teamID); err == nil {
			return fmt.Errorf("Automation Actions Runner Team association still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Automation Actions Runner ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		runnerID, teamID, err := resourcePagerDutyParseColonCompoundID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, _, err := client.AutomationActionsRunner.GetAssociationToTeam(runnerID, teamID)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%s:%s", runnerID, found.Team.ID) != rs.Primary.ID {
			return fmt.Errorf("Automation Actions Runner association to team not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationConfig(teamName, runnerName string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "foo"
}

resource "pagerduty_automation_actions_runner" "foo" {
	name = "%s runner"
	description = "Runner created by TF"
	runner_type = "runbook"
	runbook_base_uri = "cat-cat"
	runbook_api_key = "cat-secret"
}

resource "pagerduty_automation_actions_runner_team_association" "foo" {
  runner_id = pagerduty_automation_actions_runner.foo.id
  team_id = pagerduty_team.foo.id
}

`, teamName, runnerName)
}
