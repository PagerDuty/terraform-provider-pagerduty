package pagerduty

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyCustomFieldSchemaAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldSchemaAssignment(),
				ExpectError: regexp.MustCompile("The custom field schema feature has been removed"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldSchemaAssignment() string {
	return `
resource "pagerduty_custom_field_schema_assignment" "test" {
  schema    = "test1"
  service   = "test2"
}
`
}
