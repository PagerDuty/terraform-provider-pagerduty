package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyLicense_Basic(t *testing.T) {
	reference := "full_user"
	name := "User"
	description := ""

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyLicenseConfig(reference, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyLicense(fmt.Sprintf("data.pagerduty_license.%s", reference)),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyLicense_Empty(t *testing.T) {
	// Note that this test does not actually set any values for the name or
	// description of the license. An accounts license data changes over time and
	// per account. So, this test only verifies that a license can be found with
	// an empty pagerduty_license datasource
	reference := "full_user"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEmptyPagerDutyLicenseConfig(reference),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyLicense(fmt.Sprintf("data.pagerduty_license.%s", reference)),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyLicense_Error(t *testing.T) {
	reference := acctest.RandString(30)
	expectedErrorString := fmt.Sprintf("Unable to locate any license with the configured id: '', name: '%s' or description: ''", reference)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePagerDutyLicenseConfigError(reference),
				ExpectError: regexp.MustCompile(expectedErrorString),
			},
		},
	})
}

func testAccDataSourcePagerDutyLicense(n string) resource.TestCheckFunc {
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

		for _, att := range testAtts {
			if _, ok := a[att]; !ok {
				return fmt.Errorf("Expected the required attribute %s to exist", att)
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyLicenseConfig(reference string, name string, description string) string {
	return fmt.Sprintf(`
data "pagerduty_license" "%s" {
	name = "%s"
	description = "%s"
}
`, reference, name, description)
}

func testAccDataSourceEmptyPagerDutyLicenseConfig(reference string) string {
	return fmt.Sprintf(`
data "pagerduty_license" "%s" {}
`, reference)
}

func testAccDataSourcePagerDutyLicenseConfigError(reference string) string {
	return fmt.Sprintf(`
data "pagerduty_license" "%s" {
	name = "%s"
}
`, reference, reference)
}
