package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyStandardsResourceScores_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyStandardsResourceScoresConfig(
					name, "technical_services", "PR6MHNF",
				),
				Check: testAccCheckAttributes(
					fmt.Sprintf("data.pagerduty_standards_resource_scores.%s", name),
					testStandardsResourceScores,
				),
			},
		},
	})
}

func testStandardsResourceScores(a map[string]string) error {
	testAttrs := []string{
		"id",
		"resource_type",
		"score.passing",
		"score.total",
		"standards.#",
		"standards.0.active",
		"standards.0.description",
		"standards.0.id",
		"standards.0.name",
		"standards.0.pass",
		"standards.0.type",
	}
	for _, attr := range testAttrs {
		if _, ok := a[attr]; !ok {
			return fmt.Errorf("Expected the required attribute %s to exist", attr)
		}
	}
	return nil
}

func testAccDataSourcePagerDutyStandardsResourceScoresConfig(name, rt, id string) string {
	format := `data "pagerduty_standards_resource_scores" "%s" {
  resource_type = "%s"
  id = "%s"
}`
	return fmt.Sprintf(format, name, rt, id)
}
