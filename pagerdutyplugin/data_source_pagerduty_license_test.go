package pagerduty

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyLicense_Basic(t *testing.T) {
	reference := "full_user"
	description := ""
	name := os.Getenv("PAGERDUTY_ACC_LICENSE_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckLicenseNameTests(t, name)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
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
	reference := "testing_reference_missing_license"
	expectedErrorString := "Unable to locate any license"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePagerDutyLicenseConfigError(reference),
				ExpectError: regexp.MustCompile(expectedErrorString),
			},
		},
	})
}

func TestAccDataSourcePagerDutyLicense_ErrorWithID(t *testing.T) {
	reference := "testing_reference_missing_license"
	expectedErrorString := "Unable to locate any license"
	// Even with an expected name, if the configured ID is not found there
	// should be an error
	name := "Full User"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePagerDutyLicenseConfigErrorWithID(reference, name),
				ExpectError: regexp.MustCompile(expectedErrorString),
			},
		},
	})
}

func testAccPreCheckLicenseNameTests(t *testing.T, name string) {
	if name == "" {
		t.Skip("PAGERDUTY_ACC_LICENSE_NAME not set. Skipping tests requiring license names")
	}
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

func testAccDataSourcePagerDutyLicenseConfigErrorWithID(reference, name string) string {
	return fmt.Sprintf(`
data "pagerduty_license" "%s" {
	id = "%s"
	name = "%s"
}
`, reference, reference, name)
}
