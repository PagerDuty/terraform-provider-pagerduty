package pagerduty

import (
	"fmt"
	"regexp"
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
	userFoo := fmt.Sprintf("tf-%s", acctest.RandString(5))
	userBar := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	role := "manager"
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant(userFoo, userBar, team, role, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.foo"),
				),
			},
			{
				// This test case is expected to fail because userFoo is a member of the
				// escalation policy foo
				Config:      testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantUpdated(userFoo, userBar, team, role, escalationPolicy),
				ExpectError: regexp.MustCompile("User \".*\" can't be removed from Team \".*\" as they belong to an Escalation Policy on this team"),
			},
			{
				// This test case is expected to pass because userFoo is being removed
				// from escalation policy as remediation measure to unblock the team
				// membership removal
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAfterRemediation(userFoo, userBar, team, role, escalationPolicy),
			},
			{
				// This test case is expected to pass because userFoo is no longer a
				// member of the escalation policy
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantUpdated(userFoo, userBar, team, role, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipNoExists("pagerduty_team_membership.foo"),
				),
			},
		},
	})
}

func TestAccPagerDutyTeamMembership_DestroyWithEscalationPolicyDependantAndMultipleTeams(t *testing.T) {
	userOne := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamOne := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamTwo := fmt.Sprintf("tf-%s", acctest.RandString(5))
	role := "manager"
	escalationPolicyOne := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicyTwo := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeams(userOne, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipExists("pagerduty_team_membership.one"),
				),
			},
			{
				// This test case is expected to fail because userOne is a member of the
				// teamOne which is associated with escalation policyOne
				Config:      testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeamsUpdated(userOne, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo),
				ExpectError: regexp.MustCompile("User \".*\" can't be removed from Team \".*\" as they belong to an Escalation Policy on this team"),
			},
			{
				// This test case is expected to pass because teamOne is being removed
				// from escalation policy policyOne as remediation measure to unblock the
				// team membership removal
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeamsAfterRemediation(userOne, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo),
			},
			{
				// This test case is expected to pass because teamOne is no longer a
				// associated to the escalation policyOne
				Config: testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeamsUpdated(userOne, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamMembershipNoExists("pagerduty_team_membership.one"),
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

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependant(userFoo, userBar, team, role, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[1]s@foo.test"
}

resource "pagerduty_user" "bar" {
  name  = "%[2]s"
  email = "%[2]s@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[3]s"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
  role    = "%[4]s"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[5]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "user_reference"
      id   = pagerduty_user.bar.id
    }
  }
}
`, userFoo, userBar, team, role, escalationPolicy)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantUpdated(userFoo, userBar, team, role, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[1]s@foo.test"
}

resource "pagerduty_user" "bar" {
  name  = "%[2]s"
  email = "%[2]s@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[3]s"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[5]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "user_reference"
      id   = pagerduty_user.bar.id
    }
  }
}
`, userFoo, userBar, team, role, escalationPolicy)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAfterRemediation(userFoo, userBar, team, role, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[1]s@foo.test"
}

resource "pagerduty_user" "bar" {
  name  = "%[2]s"
  email = "%[2]s@foo.test"
}

resource "pagerduty_team" "foo" {
  name        = "%[3]s"
  description = "foo"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
  role    = "%[4]s"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[5]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.bar.id
    }
  }
}`, userFoo, userBar, team, role, escalationPolicy)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeams(user, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "one" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "one" {
  name        = "%[2]v"
  description = "team_one"
}

resource "pagerduty_team" "two" {
  name        = "%[3]v"
  description = "team_two"
}

resource "pagerduty_team_membership" "one" {
  user_id = pagerduty_user.one.id
  team_id = pagerduty_team.one.id
  role    = "%[4]v"
}

resource "pagerduty_team_membership" "two" {
  user_id = pagerduty_user.one.id
  team_id = pagerduty_team.two.id
  role    = "%[4]v"
}

resource "pagerduty_escalation_policy" "one" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.one.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.one.id
    }
  }
}

resource "pagerduty_escalation_policy" "two" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.two.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.one.id
    }
  }
}
`, user, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeamsAfterRemediation(user, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "one" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "one" {
  name        = "%[2]v"
  description = "team_one"
}

resource "pagerduty_team" "two" {
  name        = "%[3]v"
  description = "team_two"
}

resource "pagerduty_team_membership" "one" {
  user_id = pagerduty_user.one.id
  team_id = pagerduty_team.one.id
  role    = "%[4]v"
}

resource "pagerduty_team_membership" "two" {
  user_id = pagerduty_user.one.id
  team_id = pagerduty_team.two.id
  role    = "%[4]v"
}

resource "pagerduty_escalation_policy" "one" {
  name      = "%s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.one.id
    }
  }
}

resource "pagerduty_escalation_policy" "two" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.two.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.one.id
    }
  }
}
`, user, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo)
}

func testAccCheckPagerDutyTeamMembershipDestroyWithEscalationPolicyDependantAndMultipleTeamsUpdated(user, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "one" {
  name = "%[1]v"
  email = "%[1]v@foo.test"
}

resource "pagerduty_team" "one" {
  name        = "%[2]v"
  description = "team_one"
}

resource "pagerduty_team" "two" {
  name        = "%[3]v"
  description = "team_two"
}

resource "pagerduty_team_membership" "two" {
  user_id = pagerduty_user.one.id
  team_id = pagerduty_team.two.id
  role    = "%[4]v"
}

resource "pagerduty_escalation_policy" "one" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.one.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.one.id
    }
  }
}

resource "pagerduty_escalation_policy" "two" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.two.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.one.id
    }
  }
}
`, user, teamOne, teamTwo, role, escalationPolicyOne, escalationPolicyTwo)
}
