package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyLicenses_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyLicensesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyLicenses(fmt.Sprintf("data.pagerduty_licenses.%s", name)),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyLicenses(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		testAtts := []string{
			"id",
			"name",
			"summary",
			"description",
			"role_group",
			"current_value",
			"allocations_available",
			"html_url",
			"self",
		}

		if val, ok := a["licenses.#"]; !ok || val == "0" {
			return fmt.Errorf("Expected licenses.licenses to have at least 1 license")
		}

		for _, att := range testAtts {
			required_sub_attr := fmt.Sprintf("licenses.0.%s", att)
			if _, ok := a[required_sub_attr]; !ok {
				return fmt.Errorf("Expected the required attribute %s to exist", required_sub_attr)
			}
		}

		if val, ok := a["licenses.0.valid_roles.#"]; !ok || val == "0" {
			return fmt.Errorf("Expected licenses.licenses[0].valid_roles to have at least 1 role")
		}

		return nil
	}
}

func testAccDataSourcePagerDutyLicensesConfig(name string) string {
	return fmt.Sprintf(`data "pagerduty_licenses" "%s" {}`, name)
}
