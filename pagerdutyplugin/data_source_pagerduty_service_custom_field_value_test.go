package pagerduty

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyServiceCustomFieldValue_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	fieldName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceCustomFieldValueConfig(username, email, escalationPolicy, service, fieldName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyServiceCustomFieldValue("data.pagerduty_service_custom_field_value.test"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyServiceCustomFieldValue(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a service custom field value ID from PagerDuty")
		}

		if a["service_id"] == "" {
			return fmt.Errorf("Expected to get a service ID from PagerDuty")
		}

		if a["custom_fields.#"] == "0" {
			return fmt.Errorf("Expected to get at least one custom field value")
		}

		// Verify the environment field exists and has the expected value
		foundEnvironmentField := false
		count, _ := strconv.Atoi(a["custom_fields.#"])
		for i := 0; i < count; i++ {
			name := a[fmt.Sprintf("custom_fields.%d.name", i)]
			if name == "environment" {
				foundEnvironmentField = true
				value := a[fmt.Sprintf("custom_fields.%d.value", i)]
				if value == "" {
					return fmt.Errorf("Expected environment field to have a value")
				}
				break
			}
		}

		if !foundEnvironmentField {
			return fmt.Errorf("Expected to find environment field in custom fields")
		}

		return nil
	}
}

func testAccDataSourcePagerDutyServiceCustomFieldValueConfig(username, email, escalationPolicy, service, fieldName string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name      = "%s"
  email     = "%s"
}

resource "pagerduty_escalation_policy" "test" {
  name        = "%s"
  description = "Managed by Terraform"
  num_loops   = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_service" "test" {
  name                    = "%s"
  description             = "Managed by Terraform"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.test.id
}

resource "pagerduty_service_custom_field" "test" {
  name         = "environment"
  display_name = "%s"
  data_type    = "string"
  field_type   = "single_value"
  description  = "Environment where this service is deployed"
}

resource "pagerduty_service_custom_field_value" "test" {
  service_id = pagerduty_service.test.id
  
  custom_fields {
    name  = pagerduty_service_custom_field.test.name
    value = jsonencode("production")
  }
}

data "pagerduty_service_custom_field_value" "test" {
  service_id = pagerduty_service.test.id
  depends_on = [pagerduty_service_custom_field_value.test]
}
`, username, email, escalationPolicy, service, fieldName)
}
