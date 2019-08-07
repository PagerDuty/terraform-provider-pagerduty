package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourcePagerDutyEventRule_Basic(t *testing.T) {
	eventRule := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEventRuleConfig(eventRule),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventRule("pagerduty_event_rule.test_data_source", "data.pagerduty_event_rule.by_id"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyEventRule(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get an event rule ID from PagerDuty")
		}

		testAtts := []string{"id"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the event rule %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEventRuleConfig(eventRule string) string {
	return fmt.Sprintf(`
variable "action_list" {
	default = [
		[
			"route",
			"P5DTL0K"
		],
		[
			"severity",
			"warning"
		],
		[
			"annotate",
			"foo bar"
		],
		[
			"priority",
			"PL451DT"
		]
	]
}
variable "condition_list" {
	default = [
		"and",
		["contains",["path","payload","source"],"website"]		]
}
variable "action_list_suppress" {
	default = [["suppress",true]]
}
resource "pagerduty_event_rule" "test_data_source" {
	action_json = jsonencode(var.action_list)
	condition_json = jsonencode(var.condition_list)
	catch_all = false
	advanced_condition_json = [%s]
}

resource "pagerduty_event_rule" "foo_data_source" {
	action_json = jsonencode(var.action_list)
	condition_json = jsonencode(var.condition_list)
}

resource "pagerduty_event_rule" "test_catchall" {
	catch_all = true,
	action_json = jsonencode(var.action_list_suppress)
	condition_json = []
}

data "pagerduty_event_rule" "by_adv_cond" {
  advanced_condition_json = pagerduty_event_rule.test_data_source.advanced_condition_json
}
`, eventRule)
}
