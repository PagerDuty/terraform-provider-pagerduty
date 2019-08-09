package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_rule", &resource.Sweeper{
		Name: "pagerduty_event_rule",
		F:    testSweepEventRule,
		Dependencies: []string{
			"pagerduty_service",
		},
	})
}

func testSweepEventRule(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.EventRules.List()
	if err != nil {
		return err
	}

	for _, rule := range resp.EventRules {
		if strings.HasPrefix(rule.ID, "test") || strings.HasPrefix(rule.ID, "tf-") {
			log.Printf("Destroying event rule %s", rule.ID)
			if _, err := client.EventRules.Delete(rule.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyEventRule_Basic(t *testing.T) {
	eventRule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	eventRuleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventRuleConfig(eventRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventRuleExists("pagerduty_event_rule.first"),
				),
			},

			{
				Config: testAccCheckPagerDutyEventRuleConfigUpdated(eventRuleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventRuleExists("pagerduty_event_rule.foo_resource_updated"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_rule" {
			continue
		}
		// get list of event_rules and then check that list.
		resp, _, err := client.EventRules.List()
		if err != nil {
			return err
		}
		for _, er := range resp.EventRules {
			if er.ID == r.Primary.ID {
				return fmt.Errorf("Event Rule still exists")
			}
		}
	}
	return nil
}

func testAccCheckPagerDutyEventRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Escalation Policy ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)
		resp, _, err := client.EventRules.List()
		if err != nil {
			return err
		}
		var found *pagerduty.EventRule

		for _, rule := range resp.EventRules {
			if rule.ID == rs.Primary.ID {
				found = rule
			}
		}
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Escalation policy not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventRuleConfig(eventRule string) string {
	return fmt.Sprintf(`
variable "action_list" {
	default = [["route","P5DTL0K"],["severity","warning"],["annotate","%s"],["priority","PL451DT"]]
}
variable "condition_list" {
	default = ["and",["contains",["path","payload","source"],"website"]]
}
variable "advanced_condition_list" {
    default = [
                [
                    "scheduled-weekly",
                    1565392127032,
                    3600000,
                    "America/Los_Angeles",
                    [
                        1,
                        3,
                        5,
                        7
                    ]
                ]
	]
}
resource "pagerduty_event_rule" "first" {
	action_json = jsonencode(var.action_list)
	condition_json = jsonencode(var.condition_list)
	advanced_condition_json = jsonencode(var.advanced_condition_list)
}
`, eventRule)
}

func testAccCheckPagerDutyEventRuleConfigUpdated(eventRule string) string {
	return fmt.Sprintf(`
variable "action_list" {
	default = [["route","P5DTL0K"],["severity","warning"],["annotate","%s"],["priority","PL451DT"]]
}
variable "condition_list" {
	default = ["and",["contains",["path","payload","source"],"website"],["contains",["path","headers","from","0","address"],"homer"]]
}
variable "advanced_condition_list" {
    default = [
                [
                    "scheduled-weekly",
                    1565392127032,
                    3600000,
                    "America/Los_Angeles",
                    [
                        1,
                        3,
                        5,
                        7
                    ]
                ]
	]
}
resource "pagerduty_event_rule" "foo_resource_updated" {
	action_json = jsonencode(var.action_list)
	condition_json = jsonencode(var.condition_list)
	advanced_condition_json = jsonencode(var.advanced_condition_list)
}
`, eventRule)
}
