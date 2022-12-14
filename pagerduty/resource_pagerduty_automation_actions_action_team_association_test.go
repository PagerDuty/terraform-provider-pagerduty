package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionTeamAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionTeamAssociationConfig(teamName),
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
		actionID, teamID := resourcePagerDutyParseColonCompoundID(r.Primary.ID)
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
		actionID, teamID := resourcePagerDutyParseColonCompoundID(rs.Primary.ID)
		found, _, err := client.AutomationActionsAction.GetAssociationToTeam(actionID, teamID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Automation Actions Action association to team not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAutomationActionsActionTeamAssociationConfig(teamName string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "foo"
}

# expecting action resource to be here

resource "pagerduty_automation_actions_action_team_association" "foo" {
  action_id = "action_id_will_be_here"
  team_id = pagerduty_team.foo.id
}

`, teamName)
}
