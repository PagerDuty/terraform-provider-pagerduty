package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyTeamMembership_Basic(t *testing.T) {
	user := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipConfig(user, team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTeamMembershipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_team_membership" {
			continue
		}

		user, _, err := client.Users.Get(r.Primary.Attributes["user_id"], &pagerduty.GetUserOptions{})
		if err == nil {
			if isTeamMember(user, r.Primary.Attributes["team_id"]) {
				return fmt.Errorf("%s is still a member of: %s", user.ID, r.Primary.Attributes["team_id"])
			}
		}
	}

	return nil
}

func testAccCheckPagerDutyTeamMembershipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*pagerduty.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		userID := rs.Primary.Attributes["user_id"]
		teamID := rs.Primary.Attributes["team_id"]

		user, _, err := client.Users.Get(userID, &pagerduty.GetUserOptions{})
		if err != nil {
			return err
		}

		if !isTeamMember(user, teamID) {
			return fmt.Errorf("%s is not a member of: %s", userID, teamID)
		}

		return nil
	}
}

func testAccCheckPagerDutyTeamMembershipConfig(user, team string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.com"
  teams = ["${pagerduty_team.foo.id}"]
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = "${pagerduty_user.foo.id}"
  team_id = "${pagerduty_team.foo.id}"
}
`, user, team)
}
