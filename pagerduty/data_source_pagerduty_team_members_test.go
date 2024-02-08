package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyTeamMembers_Basic(t *testing.T) {
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	userName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	userEmail := fmt.Sprintf("%s@foo.test", userName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyTeamMembersConfig(teamName, userName, userEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyTeamMembers("pagerduty_team.test", "pagerduty_user.test", "data.pagerduty_team_members.test"),
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

		teamR := s.RootModule().Resources[teamResource]
		teamRAS := teamR.Primary.Attributes

		teamMembershipDS := s.RootModule().Resources[teamMembershipDataSource]
		as := teamMembershipDS.Primary.Attributes

		if as["id"] == "" {
			return fmt.Errorf("Expected team members ID not to be empty")
		}

		if as["team_id"] != teamRAS["id"] {
			return fmt.Errorf("Expected team ID to be %s, but got %s", teamRAS["id"], as["team_id"])
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
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team_membership" "test" {
  user_id = pagerduty_user.test.id
  team_id = pagerduty_team.test.id
}

data "pagerduty_team_members" "test" {
  depends_on = [
    pagerduty_team_membership.test,
  ]

  team_id = pagerduty_team.test.id
}
`, teamName, teamName, userName, userEmail)
}
