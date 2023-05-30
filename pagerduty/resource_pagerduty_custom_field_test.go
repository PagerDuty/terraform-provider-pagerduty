package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyCustomFields_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	description1 := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldConfig(fieldName, description1, "string"),
				ExpectError: regexp.MustCompile("The standalone custom field feature"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldConfig(name, description, datatype string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  description = "%[2]s" 
  datatype = "%[3]s"
}
`, name, description, datatype)
}
