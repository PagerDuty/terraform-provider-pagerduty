package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_automation_actions_action_team_association", &resource.Sweeper{
		Name: "pagerduty_automation_actions_action_team_association",
		F:    testSweepAutomationActionsActionTeamAssociation,
	})
}

func testSweepAutomationActionsActionTeamAssociation(region string) error {
	return nil
}

func TestAccPagerDutyAutomationActionsActionTeamAssociation_Basic(t *testing.T) {
	actionName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionTeamAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionTeamAssociationConfig(actionName, teamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAutomationActionsActionTeamAssociationExists("pagerduty_automation_actions_action_team_association.foo"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_action_team_association.foo", "action_id"),
					resource.TestCheckResourceAttrSet("pagerduty_automation_actions_action_team_association.foo", "team_id"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAutomationActionsActionTeamAssociationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_automation_actions_action_team_association" {
			continue
		}
		actionID, teamID, err := resourcePagerDutyParseColonCompoundID(r.Primary.ID)
		if err != nil {
			return err
		}

		if _, _, err := client.AutomationActionsAction.GetAssociationToTeam(actionID, teamID); err == nil {
			return fmt.Errorf("Automation Actions Runner still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyAutomationActionsActionTeamAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Automation Actions Runner ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		actionID, teamID, err := resourcePagerDutyParseColonCompoundID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, _, err := client.AutomationActionsAction.GetAssociationToTeam(actionID, teamID)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%s:%s", actionID, found.Team.ID) != rs.Primary.ID {
			return fmt.Errorf("Automation Actions Action association to team not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAutomationActionsActionTeamAssociationConfig(teamName, actionName string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "foo"
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

resource "pagerduty_automation_actions_action_team_association" "foo" {
  action_id = pagerduty_automation_actions_action.foo.id
  team_id = pagerduty_team.foo.id
}

`, teamName, actionName)
}
