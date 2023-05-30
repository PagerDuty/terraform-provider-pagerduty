package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePagerDutyCustomFieldSchema(t *testing.T) {
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePagerDutyFieldSchemaConfig(schemaTitle),
				ExpectError: regexp.MustCompile("The custom field schema feature has been removed"),
			},
		},
	})
}

func testAccDataSourcePagerDutyFieldSchemaConfig(title string) string {
	return fmt.Sprintf(`
data "pagerduty_custom_field_schema" "%[1]s" {
  title = "%[1]s"
}
`, title)
}
