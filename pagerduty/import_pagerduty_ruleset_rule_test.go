package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyRulesetRule_import(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyRulesetRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyRulesetRuleConfig(ruleset, teamName, rule),
			},

			{
				ResourceName:      "pagerduty_ruleset_rule.foo",
				ImportStateIdFunc: testAccCheckPagerDutyRulesetRuleID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyRulesetRule_importCatchAll(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	rule1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	catchAllRule := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyRulesetRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyRulesetRuleConfigCatchAllRule(team, ruleset, rule1, catchAllRule),
			},
			{
				ResourceName:      "pagerduty_ruleset_rule.catch_all",
				ImportStateIdFunc: testAccCheckPagerDutyRulesetRuleCatchAllID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyRulesetRuleID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v.%v", s.RootModule().Resources["pagerduty_ruleset.foo"].Primary.ID, s.RootModule().Resources["pagerduty_ruleset_rule.foo"].Primary.ID), nil
}

func testAccCheckPagerDutyRulesetRuleCatchAllID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v.%v", s.RootModule().Resources["pagerduty_ruleset.foo"].Primary.ID, s.RootModule().Resources["pagerduty_ruleset_rule.catch_all"].Primary.ID), nil
}
