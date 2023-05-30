package pagerduty

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyCustomFieldConfiguration_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldConfigurationConfigBasic(),
				ExpectError: regexp.MustCompile("The standalone custom field schema feature"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldConfigurationConfigBasic() string {
	return `
resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field  = "foo"
  schema = "bar"
}
`
}
