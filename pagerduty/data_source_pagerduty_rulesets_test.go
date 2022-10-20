package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyRulesets_Basic(t *testing.T) {
	ruleset := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyRulesetsConfig(ruleset),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyRulesets("pagerduty_ruleset.test", "data.pagerduty_rulesets.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyRulesets(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a ruleset ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			sub_att := fmt.Sprintf("rulesets.0.%s", att)
			if a[sub_att] != srcA[att] {
				return fmt.Errorf("Expected the ruleset %s to be: %s, but got: %s", att, srcA[att], a[sub_att])
			}
		}
		return nil
	}
}

func testAccDataSourcePagerDutyRulesetsConfig(ruleset string) string {
	return fmt.Sprintf(`
resource "pagerduty_ruleset" "test" {
  name                    = "%s"
}

data "pagerduty_rulesets" "by_name" {
  search = pagerduty_ruleset.test.name
}
`, ruleset)
}
