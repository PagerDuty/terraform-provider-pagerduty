package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePagerDutyCustomFieldSchema(t *testing.T) {
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))
	dataSourceName := fmt.Sprintf("data.pagerduty_custom_field_schema.%s", schemaTitle)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyFieldSchemaConfig(schemaTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "title", schemaTitle),
					resource.TestCheckResourceAttr(dataSourceName, "description", "some description"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyFieldSchemaConfig(title string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field_schema" "input" {
  title = "%[1]s"
  description = "some description"
}

data "pagerduty_custom_field_schema" "%[1]s" {
  title = pagerduty_custom_field_schema.input.title
}
`, title)
}
