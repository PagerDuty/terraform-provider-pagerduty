package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyStandards_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyStandardsConfig(name),
				Check: testAccCheckAttributes(
					fmt.Sprintf("data.pagerduty_standards.%s", name),
					testStandards,
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyStandards_WithResourceType(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	path := fmt.Sprintf("data.pagerduty_standards.%s", name)
	resourceType := "technical_service"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyStandardsConfigWithResourceType(name, resourceType),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyStandardsWithResourceType(path, resourceType),
				),
			},
		},
	})
}

func testStandards(a map[string]string) error {
	testAttrs := []string{
		"id",
		"name",
		"active",
		"description",
		"type",
		"resource_type",
		"exclusions.#",
		"inclusions.#",
	}

	if val, ok := a["standards.#"]; !ok || val == "0" {
		return fmt.Errorf("Expected standards.standards to have at least 1 standard")
	}

	for _, att := range testAttrs {
		requiredSubAttr := fmt.Sprintf("standards.0.%s", att)
		if _, ok := a[requiredSubAttr]; !ok {
			return fmt.Errorf("Expected the required attribute %s to exist", requiredSubAttr)
		}
	}

	return nil
}

func testAccDataSourcePagerDutyStandardsWithResourceType(path, rt string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[path]
		a := r.Primary.Attributes

		if val, ok := a["resource_type"]; !ok || val != rt {
			return fmt.Errorf("Expected %s to match provided value: %s", val, rt)
		}

		return testStandards(a)
	}
}

func testAccDataSourcePagerDutyStandardsConfig(name string) string {
	return fmt.Sprintf(`data "pagerduty_standards" "%s" {}`, name)
}

func testAccDataSourcePagerDutyStandardsConfigWithResourceType(name, rt string) string {
	return fmt.Sprintf(`data "pagerduty_standards" "%s" { resource_type = "%s" }`, name, rt)
}
