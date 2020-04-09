package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyRulesetRule_Basic(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	ruleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyRulesetRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyRulesetRuleConfig(team, ruleset, rule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "position", "0"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.operator", "and"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.operator", "contains"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.parameter.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.parameter.0.value", "disk space"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "actions.0.annotate.0.value", rule),
				),
			},
			{
				Config: testAccCheckPagerDutyRulesetRuleConfigUpdated(team, ruleset, ruleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "position", "0"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.operator", "and"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.operator", "contains"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.parameter.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.parameter.0.path", "payload.summary"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "actions.0.annotate.0.value", ruleUpdated),
				),
			},
		},
	})
}

func testAccCheckPagerDutyRulesetRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_ruleset_rule" {
			continue
		}

		ruleset, _ := s.RootModule().Resources["pagerduty_ruleset.foo"]

		if _, _, err := client.Rulesets.GetRule(ruleset.Primary.ID, r.Primary.ID); err == nil {
			return fmt.Errorf("Ruleset Rule still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyRulesetRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ruleset Rule ID is set")
		}

		ruleset, _ := s.RootModule().Resources["pagerduty_ruleset.foo"]

		client := testAccProvider.Meta().(*pagerduty.Client)
		found, _, err := client.Rulesets.GetRule(ruleset.Primary.ID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Ruleset Rule not found: %v", rs.Primary.ID)
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Ruleset Rule not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyRulesetRuleConfig(team, ruleset, rule string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
	name = "%s"
}

resource "pagerduty_ruleset" "foo" {
	name = "%s"
	team { 
		id = pagerduty_team.foo.id
	}
}
resource "pagerduty_ruleset_rule" "foo" {
	ruleset = pagerduty_ruleset.foo.id
	position = 0
	disabled = "false"
	conditions {
		operator = "and"
		subconditions {
			operator = "contains"
			parameter {
				value = "disk space"
				path = "payload.summary"
			}
		}
	}
	actions {
		route {
			value = "P5DTL0K"
		}
		annotate {
			value = "%s"
		}
		extractions {
			target = "dedup_key"
			source = "details.host"
			regex = "(.*)"
		}
	}
}
`, team, ruleset, rule)
}

func testAccCheckPagerDutyRulesetRuleConfigUpdated(team, ruleset, rule string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
	name = "%s"
}

resource "pagerduty_ruleset" "foo" {
	name = "%s"
	team { 
		id = pagerduty_team.foo.id
	}
}
resource "pagerduty_ruleset_rule" "foo" {
	ruleset = pagerduty_ruleset.foo.id
	position = 0
	disabled = "false"
	time_frame {
		scheduled_weekly {
			weekdays = [3,7]
			timezone = "America/Los_Angeles"
			start_time = "1000000"
			duration = "3600000"

		}
	}
	conditions {
		operator = "and"
		subconditions {
			operator = "contains"
			parameter {
				value = "disk space"
				path = "payload.summary"
			}
		}
	}
	actions {
		route {
			value = "P5DTL0K"
		}
		severity  {
			value = "warning"
		}
		annotate {
			value = "%s"
		}
		extractions {
			target = "dedup_key"
			source = "details.host"
			regex = "(.*)"
		}
	}
}
`, team, ruleset, rule)
}
