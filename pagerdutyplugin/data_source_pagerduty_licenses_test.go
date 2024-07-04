package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyLicenses_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
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

func TestAccDataSourcePagerDutyLicenses_WithID(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyLicensesConfigWithID(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyLicensesWithID(fmt.Sprintf("data.pagerduty_licenses.%s", name), name),
				),
			},
		},
	})
}

func testLicenses(a map[string]string) error {
	testAttrs := []string{
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

	for _, att := range testAttrs {
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

func testAccDataSourcePagerDutyLicensesWithID(n string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if val, ok := a["id"]; !ok || val != id {
			return fmt.Errorf("Expected id to match provided value: %s\n%#v", id, a)
		}

		return testLicenses(a)
	}
}

func testAccDataSourcePagerDutyLicenses(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		return testLicenses(a)
	}
}

func testAccDataSourcePagerDutyLicensesConfig(name string) string {
	return fmt.Sprintf(`data "pagerduty_licenses" "%s" {}`, name)
}

func testAccDataSourcePagerDutyLicensesConfigWithID(name string) string {
	return fmt.Sprintf(`data "pagerduty_licenses" "%s" {
		id = "%s"
	}`, name, name)
}
