package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestAccPagerDutyTeamMembership_WithRole(t *testing.T) {
	user := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	role := "manager"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipWithRoleConfig(user, team, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
				),
			},
		},
	})
}

func TestAccPagerDutyTeamMembership_WithRoleConsistentlyAssigned(t *testing.T) {
	user := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	firstRole := "observer"
	secondRole := "responder"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipWithRoleConfig(user, team, firstRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_team_membership.foo", "role", firstRole),
				),
			},
			{
				Config: testAccCheckPagerDutyTeamMembershipWithRoleConfig(user, team, secondRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_team_membership.foo", "role", secondRole),
				),
			},
		},
	})
}

func TestAccPagerDutyTeamMembership_DestroyWithEscalationPolicyDependant(t *testing.T) {
	user := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	role := "manager"
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant(user, team, role, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantUpdated(user, team, role, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipNoExists("pagerduty_team_membership.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTeamMembershipDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
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
		client, _ := testAccProvider.Meta().(*Config).Client()
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		userID := rs.Primary.Attributes["user_id"]
		teamID := rs.Primary.Attributes["team_id"]
		role := rs.Primary.Attributes["role"]

		user, _, err := client.Users.Get(userID, &pagerduty.GetUserOptions{})
		if err != nil {
			return err
		}

		if !isTeamMember(user, teamID) {
			return fmt.Errorf("%s is not a member of: %s", userID, teamID)
		}

		resp, _, err := client.Teams.GetMembers(teamID, &pagerduty.GetMembersOptions{})
		if err != nil {
			return err
		}

		for _, member := range resp.Members {
			if member.User.ID == userID {
				if member.Role != role {
					return fmt.Errorf("%s does not have the role: %s in: %s", userID, role, teamID)
				}
			}
		}

		return nil
	}
}

func testAccCheckPagerDutyTeamMembershipNoExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, _ := testAccProvider.Meta().(*Config).Client()
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return nil
		}

		if rs.Primary.ID == "" {
			return nil
		}

		userID := rs.Primary.Attributes["user_id"]
		teamID := rs.Primary.Attributes["team_id"]

		user, _, err := client.Users.Get(userID, &pagerduty.GetUserOptions{})
		if err != nil {
			return err
		}

		if isTeamMember(user, teamID) {
			return fmt.Errorf("%s is still a member of: %s", userID, teamID)
		}

		return nil
	}
}

func testAccCheckPagerDutyTeamMembershipConfig(user, team string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
}
`, user, team)
}

func testAccCheckPagerDutyTeamMembershipWithRoleConfig(user, team, role string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
  role    = "%[3]v"
}
`, user, team, role)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant(user, team, role, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
  role    = "%[3]v"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, user, team, role, escalationPolicy)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantUpdated(user, team, role, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[4]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, user, team, role, escalationPolicy)
}
