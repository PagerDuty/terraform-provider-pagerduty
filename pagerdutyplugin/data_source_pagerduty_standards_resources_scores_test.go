package pagerduty

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyStandardsResourcesScores_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyStandardsResourcesScoresConfig(
					name, "technical_services", []string{"P703E9Q", "PR6MHNF", "PQVUB8D"},
				),
				Check: testAccCheckAttributes(
					fmt.Sprintf("data.pagerduty_standards_resources_scores.%s", name),
					testStandardsResourcesScores,
				),
			},
		},
	})
}

func testStandardsResourcesScores(a map[string]string) error {
	testAttrs := []string{
		"ids.#",
		"ids.0",
		"resource_type",
		"resources.#",
		"resources.0.resource_id",
		"resources.0.resource_type",
		"resources.0.score.passing",
		"resources.0.score.total",
		"resources.0.standards.#",
		"resources.0.standards.0.active",
		"resources.0.standards.0.description",
		"resources.0.standards.0.id",
		"resources.0.standards.0.name",
		"resources.0.standards.0.pass",
		"resources.0.standards.0.type",
	}

	for _, attr := range testAttrs {
		if _, ok := a[attr]; !ok {
			return fmt.Errorf("Expected the required attribute %s to exist", attr)
		}
	}

	return nil
}

func testAccDataSourcePagerDutyStandardsResourcesScoresConfig(name, rt string, ids []string) string {
	format := `data "pagerduty_standards_resources_scores" "%s" {
  resource_type = "%s"
  ids = ["%s"]
}`
	idsList := strings.Join(ids, `","`)
	return fmt.Sprintf(format, name, rt, idsList)
}
