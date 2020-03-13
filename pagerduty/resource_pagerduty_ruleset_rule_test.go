package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_ruleset_rule", &resource.Sweeper{
		Name: "pagerduty_ruleset_rule",
		F:    testSweepRuleset,
	})
}

func testSweepRulesetRule(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}
	// todo: need to fix this before pushing

	// resp, _, err := client.Rulesets.ListRules()
	// if err != nil {
	// 	return err
	// }
	// for _, ruleset := range resp.Rules {
	// 	if strings.HasPrefix(ruleset.Name, "test") || strings.HasPrefix(ruleset.Name, "tf-") {
	// 		log.Printf("Destroying ruleset %s (%s)", ruleset.Name, ruleset.ID)
	// 		if _, err := client.Rulesets.Delete(ruleset.ID); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
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
						"pagerduty_ruleset_rule.foo", "action.#", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "action.1.parameters.0.value", rule),
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
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.parameter.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "conditions.0.subconditions.0.parameter.0.path", "payload.summary"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "action.#", "4"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "action.2.parameters.0.value", ruleUpdated),
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
	action {
		action = "route"
		parameters {
			value = "P5DTL0K"
		}
	}
	action {
		action = "annotate"
		parameters {
			value = "%s"
		}
	}
	action {
		action = "extract"
		parameters {
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
	conditions {
		operator = "and"
		subconditions {
			operator = "contains"
			parameter {
				value = "disk space"
				path = "payload.summary"
			}
			parameter {
				value = "db"
				path = "payload.source"
			}
		}
	}
	action {
		action = "route"
		parameters {
			value = "P5DTL0K"
		}
	}
	action {
		action = "severity"
		parameters {
			value = "warning"
		}
	}
	action {
		action = "annotate"
		parameters {
			value = "%s"
		}
	}
	action {
		action = "extract"
		parameters {
			target = "dedup_key"
			source = "details.host"
			regex = "(.*)"
		}
	}
}
`, team, ruleset, rule)
}
