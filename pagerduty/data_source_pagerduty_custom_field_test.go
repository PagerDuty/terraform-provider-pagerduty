package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePagerDutyCustomField(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePagerDutyCustomFieldConfig(fieldName),
				ExpectError: regexp.MustCompile("The standalone custom field feature has been removed"),
			},
		},
	})
}

func testAccDataSourcePagerDutyCustomFieldConfig(name string) string {
	return fmt.Sprintf(`
data "pagerduty_custom_field" "%[1]s" {
  name = "%[1]s"
}
`, name)
}
