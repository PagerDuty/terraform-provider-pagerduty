package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyTeamMembers_Basic(t *testing.T) {
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	userName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	userEmail := fmt.Sprintf("%s@pagerduty.com", userName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyTeamMembersConfig(teamName, userName, userEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pagerduty_team_members.test", "members.#", "1"),
					resource.TestCheckResourceAttr("data.pagerduty_team_members.test", "members.0.summary", userName),
					resource.TestCheckResourceAttr("data.pagerduty_team_members.test", "members.0.role", "manager"),
					resource.TestCheckResourceAttr("data.pagerduty_team_members.test", "members.0.type", "user_reference"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyTeamMembers(teamResource, userResource, teamMembershipDataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		userR := s.RootModule().Resources[userResource]
		userRAs := userR.Primary.Attributes

		teamMembershipDS := s.RootModule().Resources[teamMembershipDataSource]
		as := teamMembershipDS.Primary.Attributes

		if as["id"] != "" {
			return fmt.Errorf("Expected team members ID not to be empty")
		}

		if as["members.0.id"] != userRAs["id"] {
			return fmt.Errorf("Expected team member ID to match user ID")
		}

		return nil
	}
}

func testAccDataSourcePagerDutyTeamMembersConfig(teamName, userName, userEmail string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "test" {
  name        = "%s"
  description = "%s"
}

resource "pagerduty_user" "test" {
  name = "%s"
  email = "%s"
}

resource "pagerduty_team_membership" "test" {
  user_id = pagerduty_user.test.id
  team_id = pagerduty_team.test.id
}

data "pagerduty_team_members" "test" {
	team_id = pagerduty_team.test.id
}
`, teamName, teamName, userName, userEmail)
}
