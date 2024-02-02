package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyIncidentCustomField(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataSourceName := fmt.Sprintf("data.pagerduty_incident_custom_field.%s", fieldName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyIncidentCustomFieldConfig(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", fieldName),
					resource.TestCheckResourceAttr(dataSourceName, "data_type", "string"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyIncidentCustomFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  data_type = "string"
  field_type = "single_value"
}

data "pagerduty_incident_custom_field" "%[1]s" {
  name = pagerduty_incident_custom_field.input.name
}
`, name)
}
