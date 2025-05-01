package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyServiceCustomField_Basic(t *testing.T) {
	displayName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceCustomFieldConfig(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyServiceCustomField("pagerduty_service_custom_field.test", "data.pagerduty_service_custom_field.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyServiceCustomField(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a service custom field ID from PagerDuty")
		}

		testAtts := []string{
			"id", "display_name", "name", "type", "summary",
			"self", "description", "data_type", "field_type",
			"default_value", "enabled", "field_options",
		}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the service custom field %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyServiceCustomFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test" {
  name         = "regions"
  display_name = "%s"
  data_type    = "string"
  field_type   = "multi_value_fixed"
  description  = "AWS regions where this service is deployed"

  field_option {
    value     = "us-east-1"
    data_type = "string"
  }

  field_option {
    value     = "us-west-1"
    data_type = "string"
  }
}

data "pagerduty_service_custom_field" "by_name" {
    display_name = "%[1]s"
    depends_on = [pagerduty_service_custom_field.test]
}
`, name)
}
