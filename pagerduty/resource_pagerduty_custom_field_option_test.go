package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyCustomFieldOptions_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("tf_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldOptionConfig(fieldName, fieldOptionValue),
				ExpectError: regexp.MustCompile("The standalone custom field feature has been removed"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldOptionConfig(name string, fieldOptionValue string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field_option" "test" {
  field = "%s"
  datatype = "string"
  value = "%s"
}

`, name, fieldOptionValue)
}
