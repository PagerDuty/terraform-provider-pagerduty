package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyAutomationActionsActionTeamAssociation_import(t *testing.T) {
	actionName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionTeamAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionTeamAssociationConfig(actionName, teamName),
			},
			{
				ResourceName:      "pagerduty_automation_actions_action_team_association.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
