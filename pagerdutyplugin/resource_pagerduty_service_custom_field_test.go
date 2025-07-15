package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_service_custom_field", &resource.Sweeper{
		Name: "pagerduty_service_custom_field",
		F:    testSweepServiceCustomField,
	})
}

func testSweepServiceCustomField(_ string) error {
	ctx := context.Background()

	options := pagerduty.ListServiceCustomFieldsOptions{
		Include: []string{"field_option"},
	}

	resp, err := testAccProvider.client.ListServiceCustomFields(ctx, options)
	if err != nil {
		return err
	}

	for _, field := range resp.Fields {
		if strings.HasPrefix(field.Name, "tf_") || strings.HasPrefix(field.Name, "test_") {
			log.Printf("Destroying service custom field %s (%s)", field.Name, field.ID)
			if err := testAccProvider.client.DeleteServiceCustomField(ctx, field.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyServiceCustomField_Basic(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	displayName := fmt.Sprintf("TF Test %s", acctest.RandString(5))
	updatedDisplayName := fmt.Sprintf("TF Test Updated %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyServiceCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfig(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "field_type", "single_value"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "description", "Test service custom field"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigUpdated(name, updatedDisplayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "summary", updatedDisplayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "field_type", "single_value"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "description", "Updated test service custom field"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test", "enabled", "false"),
				),
			},
			{
				ResourceName:      "pagerduty_service_custom_field.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyServiceCustomField_WithOptions(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	displayName := fmt.Sprintf("TF Test %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigWithOptions(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "field_type", "single_value_fixed"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "field_option.#", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "default_value", `"production"`),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigWithUpdatedOptions(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyServiceCustomFieldExists("pagerduty_service_custom_field.test_options"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_options", "field_option.#", "4"),
					resource.TestCheckNoResourceAttr(
						"pagerduty_service_custom_field.test_options", "default_value"),
				),
			},
			{
				ResourceName:      "pagerduty_service_custom_field.test_options",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyServiceCustomField_MultiValueFixed(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	displayName := fmt.Sprintf("TF Test %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigMultiValueFixed(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_type", "multi_value_fixed"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.0.value", "us-east-1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.1.value", "us-west-1"),
				),
			},
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigMultiValueFixedUpdated(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_type", "multi_value_fixed"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.#", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.0.value", "us-east-1"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.1.value", "us-east-2"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_multi", "field_option.2.value", "us-west-2"),
				),
			},
		},
	})
}

func TestAccPagerDutyServiceCustomField_BooleanType(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	displayName := fmt.Sprintf("TF Test Boolean %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigBooleanType(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_boolean", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_boolean", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_boolean", "data_type", "boolean"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_boolean", "field_type", "single_value"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_boolean", "default_value", "true"),
				),
			},
		},
	})
}

func TestAccPagerDutyServiceCustomField_IntegerType(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	displayName := fmt.Sprintf("TF Test Integer %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyServiceCustomFieldConfigIntegerType(name, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_integer", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_integer", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_integer", "data_type", "integer"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_integer", "field_type", "single_value"),
					resource.TestCheckResourceAttr(
						"pagerduty_service_custom_field.test_integer", "default_value", "42"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyServiceCustomFieldDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_service_custom_field" {
			continue
		}

		ctx := context.Background()
		options := pagerduty.ListServiceCustomFieldsOptions{
			Include: []string{"field_option"},
		}

		_, err := testAccProvider.client.GetServiceCustomField(ctx, r.Primary.ID, options)
		if err == nil {
			return fmt.Errorf("Service custom field still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyServiceCustomFieldExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No service custom field ID is set")
		}

		ctx := context.Background()
		options := pagerduty.ListServiceCustomFieldsOptions{
			Include: []string{"field_option"},
		}

		found, err := testAccProvider.client.GetServiceCustomField(ctx, rs.Primary.ID, options)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Service custom field not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyServiceCustomFieldConfig(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test" {
  name         = "%s"
  display_name = "%s"
  data_type    = "string"
  field_type   = "single_value"
  description  = "Test service custom field"
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigUpdated(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test" {
  name         = "%s"
  display_name = "%s"
  data_type    = "string"
  field_type   = "single_value"
  description  = "Updated test service custom field"
  enabled      = false
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigWithOptions(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test_options" {
  name         = "%s"
  display_name = "%s"
  data_type    = "string"
  field_type   = "single_value_fixed"
  description  = "Test service custom field with options"
  enabled      = true
  default_value = jsonencode("production")

  field_option {
    value = "production"
    data_type = "string"
  }

  field_option {
    value = "staging"
    data_type = "string"
  }

  field_option {
    value = "development"
    data_type = "string"
  }
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigWithUpdatedOptions(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test_options" {
  name         = "%s"
  display_name = "%s"
  data_type    = "string"
  field_type   = "single_value_fixed"
  description  = "Test service custom field with options"
  enabled      = true

  field_option {
    data_type = "string"
    value = "production"
  }

  field_option {
    data_type = "string"
    value = "staging"
  }

  field_option {
    data_type = "string"
    value = "development"
  }

  field_option {
    data_type = "string"
    value = "testing"
  }
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigMultiValueFixed(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test_multi" {
  name         = "%s"
  display_name = "%s"
  data_type    = "string"
  field_type   = "multi_value_fixed"
  description  = "Test multi-value fixed service custom field"
  enabled      = true

  field_option {
    value = "us-east-1"
    data_type = "string"
  }

  field_option {
    value = "us-west-1"
    data_type = "string"
  }
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigMultiValueFixedUpdated(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test_multi" {
  name         = "%s"
  display_name = "%s"
  data_type    = "string"
  field_type   = "multi_value_fixed"
  description  = "Test multi-value fixed service custom field"
  enabled      = true

  field_option {
    value = "us-east-1"
    data_type = "string"
  }

  field_option {
    value = "us-east-2"
    data_type = "string"
  }

  field_option {
    value = "us-west-2"
    data_type = "string"
  }
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigBooleanType(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test_boolean" {
  name          = "%s"
  display_name  = "%s"
  data_type     = "boolean"
  field_type    = "single_value"
  description   = "Test boolean service custom field"
  default_value = jsonencode(true)
  enabled       = true
}
`, name, displayName)
}

func testAccCheckPagerDutyServiceCustomFieldConfigIntegerType(name, displayName string) string {
	return fmt.Sprintf(`
resource "pagerduty_service_custom_field" "test_integer" {
  name          = "%s"
  display_name  = "%s"
  data_type     = "integer"
  field_type    = "single_value"
  description   = "Test integer service custom field"
  default_value = jsonencode(42)
  enabled       = true
}
`, name, displayName)
}
