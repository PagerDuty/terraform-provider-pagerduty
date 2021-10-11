package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPagerDutyTagAssignment_import(t *testing.T) {
	tag := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTagAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagAssignmentTeamConfig(tag, team),
			},

			{
				ResourceName:      "pagerduty_tag_assignment.foo",
				ImportStateIdFunc: testAccCheckPagerDutyTagAssignmentID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyTagAssignmentID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v.%v.%v", "teams", s.RootModule().Resources["pagerduty_team.foo"].Primary.ID, s.RootModule().Resources["pagerduty_tag.foo"].Primary.ID), nil
}
