package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
						"pagerduty_ruleset_rule.foo", "disabled", "true"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "variable.#", "2"),
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
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "actions.0.extractions.1.template", "{{VAR1}} | {{VAR2}}"),
				),
			},
			{
				Config: testAccCheckPagerDutyRulesetRuleConfigUpdated(team, ruleset, ruleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "position", "0"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "disabled", "false"),
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

func TestAccPagerDutyRulesetRule_MultipleRules(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule3 := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyRulesetRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyRulesetRuleConfigMultipleRules(team, ruleset, rule1, rule2, rule3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.foo"),
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "position", "0"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.bar", "position", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.baz", "position", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "disabled", "false"),
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
						"pagerduty_ruleset_rule.foo", "actions.0.annotate.0.value", rule1),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.bar", "actions.0.annotate.0.value", rule2),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.baz", "actions.0.annotate.0.value", rule3),
				),
			},
		},
	})
}

func TestAccPagerDutyRulesetRule_CatchAllRule(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	catch_all_rule := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyRulesetRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyRulesetRuleConfigCatchAllRule(team, ruleset, rule1, catch_all_rule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.foo"),
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.catch_all"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "position", "0"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "position", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "disabled", "false"),
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
						"pagerduty_ruleset_rule.foo", "actions.0.annotate.0.value", rule1),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "catch_all", "true"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "actions.0.annotate.0.value", catch_all_rule),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "actions.0.suppress.0.value", "true"),
				),
			},
		},
	})
}

func TestAccPagerDutyRulesetRule_CatchAllRuleRoute(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	catch_all_rule := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyRulesetRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyRulesetRuleConfigCatchAllRuleRoute(team, ruleset, rule1, catch_all_rule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.foo"),
					testAccCheckPagerDutyRulesetRuleExists("pagerduty_ruleset_rule.catch_all"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "position", "0"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "position", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.foo", "disabled", "false"),
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
						"pagerduty_ruleset_rule.foo", "actions.0.annotate.0.value", rule1),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "catch_all", "true"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "actions.0.annotate.0.value", catch_all_rule),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "actions.0.suppress.0.value", "false"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "actions.0.route.0.value", "P5DTL0K"),
					resource.TestCheckResourceAttr(
						"pagerduty_ruleset_rule.catch_all", "actions.0.severity.0.value", "info"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyRulesetRuleDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
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

		client, _ := testAccProvider.Meta().(*Config).Client()
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
	disabled = true
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
		suppress {
			value = true
		}
		extractions {
			target = "dedup_key"
			source = "details.host"
			regex = "(.*)"
		}
		extractions {
			target   = "summary"
			template = "{{VAR1}} | {{VAR2}}"
		}
	}
	variable {
		type = "regex"
		parameters {
		  value = "another.*regex"
		  path = "custom_details.path.to.field"
		}
		name = "VAR2"
	}
	variable {
		type = "regex"
		parameters {
			value = ".*"
			path = "class"
		}
		name = "VAR1"
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
	disabled = false
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
		suppress {
			value = false
		}
		extractions {
			target = "dedup_key"
			source = "details.host"
			regex = "(.*)"
		}
	}
	variable {
		type = "regex"
		parameters {
		  value = "another.*regex"
		  path = "custom_details.path.to.field"
		}
		name = "VAR2"
	}
	variable {
		type = "regex"
		parameters {
			value = ".*"
			path = "class"
		}
		name = "VAR1"
	}
}
`, team, ruleset, rule)
}

func testAccCheckPagerDutyRulesetRuleConfigMultipleRules(team, ruleset, rule1, rule2, rule3 string) string {
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
	disabled = false
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
				path = "summary"
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
			source = "source"
			regex = "(.*)"
		}
	}
}
resource "pagerduty_ruleset_rule" "bar" {
	ruleset = pagerduty_ruleset.foo.id
	position = 1
	disabled = true
	conditions {
		operator = "and"
		subconditions {
			operator = "contains"
			parameter {
				value = "cpu spike"
				path = "summary"
			}
		}
	}
	actions {
		annotate {
			value = "%s"
		}
	}
}
resource "pagerduty_ruleset_rule" "baz" {
	ruleset = pagerduty_ruleset.foo.id
	position = 2
	disabled = true
	conditions {
		operator = "and"
		subconditions {
			operator = "contains"
			parameter {
				value = "slow database connection"
				path = "summary"
			}
		}
	}
	actions {
		annotate {
			value = "%s"
		}
	}
}
`, team, ruleset, rule1, rule2, rule3)
}

func testAccCheckPagerDutyRulesetRuleConfigCatchAllRule(team, ruleset, rule1, catch_all_rule string) string {
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
	disabled = false
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
				path = "summary"
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
			source = "source"
			regex = "(.*)"
		}
	}
}
resource "pagerduty_ruleset_rule" "catch_all" {
	ruleset = pagerduty_ruleset.foo.id
	position = 1
	catch_all = true
	actions {
		annotate {
			value = "%s"
		}
		suppress {
			value = true
		}
	}
}
`, team, ruleset, rule1, catch_all_rule)
}

func testAccCheckPagerDutyRulesetRuleConfigCatchAllRuleRoute(team, ruleset, rule1, catch_all_rule string) string {
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
	disabled = false
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
				path = "summary"
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
			source = "source"
			regex = "(.*)"
		}
	}
}
resource "pagerduty_ruleset_rule" "catch_all" {
	ruleset = pagerduty_ruleset.foo.id
	position = 1
	catch_all = true
	actions {
		annotate {
			value = "%s"
		}
		suppress {
			value = false
		}
		route {
			value = "P5DTL0K"
		}
		severity  {
			value = "info"
		}
	}
}
`, team, ruleset, rule1, catch_all_rule)
}
