package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePagerDutyCustomField(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataSourceName := fmt.Sprintf("data.pagerduty_custom_field.%s", fieldName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyCustomFieldConfig(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", fieldName),
					resource.TestCheckResourceAttr(dataSourceName, "datatype", "string"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyCustomFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "string"
}

data "pagerduty_custom_field" "%[1]s" {
  name = pagerduty_custom_field.input.name
}
`, name)
}
