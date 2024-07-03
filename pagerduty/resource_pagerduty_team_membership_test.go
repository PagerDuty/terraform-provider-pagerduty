package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccPagerDutyTeamMembership_basic(t *testing.T) { // from gpt and modified
	user := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	role := "manager"
	escalationPolicy1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy2 := fmt.Sprintf("tf-%s", acctest.RandString(5))

	fmt.Print("starting tests\n")
	fmt.Printf("starting tests\n\n")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant_UnrelatedEPs(user, team1, team2, role, escalationPolicy1, escalationPolicy2),
				Check: resource.ComposeTestCheckFunc(
					printdebug("checkpoint1-1\n\n"),
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
					printdebug("checkpoint1-2\n\n"),

					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.bar"),
					printdebug("checkpoint1-3\n\n"),

					resource.TestCheckResourceAttr("pagerduty_team_membership.bar", "role", "manager"),
					printdebug("checkpoint1-4\n\n"),
				),
			},
			{
				// ResourceName: "pagerduty_team_membership.foo",
				// ImportState:  true,
				// RefreshState: true,
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant_UnrelatedEPsUpdated(user, team1, team2, role, escalationPolicy1, escalationPolicy2),
				Check: resource.ComposeTestCheckFunc(
					printdebug("checkpoint2-1\n\n"),

					testAccCheckPagerDutyTeamMembershipDestroyState("pagerduty_team_membership.foo"),
					printdebug("checkpoint2-2\n\n"),

					testAccCheckPagerDutyTeamExists("pagerduty_team.foobar"),
					printdebug("checkpoint2-3\n\n"),

					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					printdebug("checkpoint2-4\n\n"),

					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.bar"),
					printdebug("checkpoint2-5\n\n"),

					testAccCheckPagerDutyEscalationPolicyTeamsFieldMatches("pagerduty_escalation_policy.foo", "foo"),
					printdebug("checkpoint2-6\n\n"),

					testAccCheckPagerDutyEscalationPolicyTeamsFieldMatches("pagerduty_escalation_policy.bar", "bar"),
					printdebug("checkpoint2-7\n\n"),
				),
			},
		},
	})
}

func printdebug(msg string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Printf(msg)
		return nil
	}
}

func testAccCheckPagerDutyEscalationPolicyTeamsFieldMatches(n, expectedTeam string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		teams := rs.Primary.Attributes["teams.#"]
		if teams != "1" {
			return fmt.Errorf("Expected teams to have 1 element, got: %s", teams)
		}

		team := rs.Primary.Attributes["teams.0"]
		if team != expectedTeam {
			return fmt.Errorf("Expected team to be %s, got: %s", expectedTeam, team)
		}

		return fmt.Errorf("dummy error")
		// return nil
	}
}

func testAccCheckPagerDutyTeamMembershipDestroyState(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Resource still exists: %s", n)
		}
		return fmt.Errorf("dummy error: %s", "wow")
		// return nil
	}
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant_UnrelatedEPs(user, team1, team2, role, escalationPolicy1, escalationPolicy2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_team" "bar" {
  name        = "%[5]v"
  description = "bar"
}

resource "pagerduty_team_membership" "foo" {
  user_id = data.pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
  role    = "%[3]v"
}

resource "pagerduty_team_membership" "bar" {
  user_id = data.pagerduty_user.foo.id
  team_id = pagerduty_team.bar.id
  role    = "%[3]v"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[4]v"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]
  description = "foo"

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = data.pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "%[6]v"
  num_loops = 2
  teams     = [pagerduty_team.bar.id]
  description = "bar"

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = data.pagerduty_user.foo.id
    }
  }
}
`, user, team1, role, escalationPolicy1, team2, escalationPolicy2)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant_UnrelatedEPsUpdated(user, team1, team2, role, escalationPolicy1, escalationPolicy2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[2]v"
  description = "foo"
}

resource "pagerduty_team" "bar" {
  name        = "%[5]v"
  description = "bar"
}

// now the user should be on a single team
resource "pagerduty_team_membership" "bar" {
  user_id = data.pagerduty_user.foo.id
  team_id = pagerduty_team.bar.id
  role    = "%[3]v"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[4]v"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = data.pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "%[6]v"
  num_loops = 2
  description = "bar"
  teams     = [pagerduty_team.bar.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = data.pagerduty_user.foo.id
    }
  }
}
`, user, team1, role, escalationPolicy1, team2, escalationPolicy2)
}
