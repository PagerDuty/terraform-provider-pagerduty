package pagerduty

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_escalation_policy", &resource.Sweeper{
		Name: "pagerduty_escalation_policy",
		F:    testSweepEscalationPolicy,
		Dependencies: []string{
			"pagerduty_service",
		},
	})
}

func testSweepEscalationPolicy(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.EscalationPolicies.List(&pagerduty.ListEscalationPoliciesOptions{})
	if err != nil {
		return err
	}

	for _, escalation := range resp.EscalationPolicies {
		if strings.HasPrefix(escalation.Name, "test") || strings.HasPrefix(escalation.Name, "tf-") {
			log.Printf("Destroying escalation policy %s (%s)", escalation.Name, escalation.ID)
			if _, err := client.EscalationPolicies.Delete(escalation.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyEscalationPolicy_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicyUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEscalationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEscalationPolicyConfig(username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "num_loops", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_delay_in_minutes", "10"),
				),
			},
			{
				Config: testAccCheckPagerDutyEscalationPolicyConfigUpdated(username, email, escalationPolicyUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicyUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "num_loops", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_delay_in_minutes", "10"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.1.escalation_delay_in_minutes", "20"),
				),
			},
			// Validating that externally removed escalation policies are detected and
			// planed for re-creation
			{
				Config: testAccCheckPagerDutyEscalationPolicyConfigUpdated(username, email, escalationPolicyUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyEscalationPolicy("pagerduty_escalation_policy.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPagerDutyEscalationPolicyWithRoundRobinAssignmentStrategy(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEscalationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEscalationPolicyConfig(username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_rule_assignment_strategy.0.type", "assign_to_everyone"),
				),
			},
			{
				Config:      testAccCheckPagerDutyEscalationPolicyWithRoundRoundAssignmentStrategyConfig(username, email, escalationPolicy, "not_valid_strategy"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Must be one of \[\]string{"assign_to_everyone", "round_robin"}`),
			},
			{
				Config: testAccCheckPagerDutyEscalationPolicyWithRoundRoundAssignmentStrategyConfig(username, email, escalationPolicy, "round_robin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "num_loops", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_rule_assignment_strategy.0.type", "round_robin"),
				),
			},
			{
				Config: testAccCheckPagerDutyEscalationPolicyWithRoundRoundAssignmentStrategyConfig(username, email, escalationPolicy, "assign_to_everyone"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_rule_assignment_strategy.0.type", "assign_to_everyone"),
				),
			},
		},
	})
}

func TestAccPagerDutyEscalationPolicy_FormatValidation(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	errMessageMatcher := "Name can not be blank, nor contain the characters.*, or any non-printable characters. Trailing white spaces are not allowed either."

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEscalationPolicyDestroy,
		Steps: []resource.TestStep{
			// Just a valid name
			{
				Config:             testAccCheckPagerDutyEscalationPolicyConfig(username, email, "SRE Escalation Policy"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Blank Name
			{
				Config:      testAccCheckPagerDutyEscalationPolicyConfig(username, email, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with & in it
			{
				Config:      testAccCheckPagerDutyEscalationPolicyConfig(username, email, "this name has an ampersand (&)"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with one white space at the end
			{
				Config:      testAccCheckPagerDutyEscalationPolicyConfig(username, email, "this name has a white space at the end "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with multiple white space at the end
			{
				Config:      testAccCheckPagerDutyEscalationPolicyConfig(username, email, "this name has white spaces at the end    "),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
			// Name with non printable characters
			{
				Config:      testAccCheckPagerDutyEscalationPolicyConfig(username, email, "this name has a non printable\\n character"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(errMessageMatcher),
			},
		},
	})
}

func TestAccPagerDutyEscalationPolicyWithTeams_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicyUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEscalationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEscalationPolicyWithTeamsConfig(username, email, team, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "num_loops", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_delay_in_minutes", "10"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "teams.#", "1"),
				),
			},
			{
				Config: testAccCheckPagerDutyEscalationPolicyWithTeamsConfigUpdated(username, email, team, escalationPolicyUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEscalationPolicyExists("pagerduty_escalation_policy.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicyUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "num_loops", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.0.escalation_delay_in_minutes", "10"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "rule.1.escalation_delay_in_minutes", "20"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "teams.#", "0"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEscalationPolicyWithRoundRoundAssignmentStrategyConfig(name, email, escalationPolicy, strategy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1

  rule {
    escalation_delay_in_minutes = 10
    escalation_rule_assignment_strategy {
      type = "%s"
    }

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, name, email, escalationPolicy, strategy)
}

func testAccCheckPagerDutyEscalationPolicyDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_escalation_policy" {
			continue
		}

		if _, _, err := client.EscalationPolicies.Get(r.Primary.ID, &pagerduty.GetEscalationPolicyOptions{}); err == nil {
			return fmt.Errorf("Escalation Policy still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEscalationPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Escalation Policy ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.EscalationPolicies.Get(rs.Primary.ID, &pagerduty.GetEscalationPolicyOptions{})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Escalation policy not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccExternallyDestroyEscalationPolicy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		_, err := client.EscalationPolicies.Delete(rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckPagerDutyEscalationPolicyConfig(name, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, name, email, escalationPolicy)
}

func testAccCheckPagerDutyEscalationPolicyConfigUpdated(name, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "bar"
  num_loops   = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }

  rule {
    escalation_delay_in_minutes = 20

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, name, email, escalationPolicy)
}

func testAccCheckPagerDutyEscalationPolicyWithTeamsConfig(name, email, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "foo"
  num_loops   = 1
	teams       = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, name, email, team, escalationPolicy)
}

func testAccCheckPagerDutyEscalationPolicyWithTeamsConfigUpdated(name, email, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%s"
  description = "bar"
  num_loops   = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }

  rule {
    escalation_delay_in_minutes = 20

    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, name, email, team, escalationPolicy)
}
