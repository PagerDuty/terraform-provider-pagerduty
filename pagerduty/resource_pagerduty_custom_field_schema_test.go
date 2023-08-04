package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyCustomFieldSchemas_Basic(t *testing.T) {
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldSchemaConfigBasic(schemaTitle),
				ExpectError: regexp.MustCompile("The standalone custom field schema feature"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldSchemaConfigBasic(title string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field_schema" "test" {
  title = "%[1]s"
  description = "some description"
}
`, title)
}
