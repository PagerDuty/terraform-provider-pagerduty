package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyAutomationActionsRunnerTeamAssociation_import(t *testing.T) {
	runnerName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsRunnerTeamAssociationConfig(runnerName, teamName),
			},
			{
				ResourceName:      "pagerduty_automation_actions_runner_team_association.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
