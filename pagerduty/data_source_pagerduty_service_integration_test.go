package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourcePagerDutyIntegration_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyIntegrationConfigStep1(username, email, service, escalationPolicy, serviceIntegration),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyIntegration("pagerduty_service.test", "data.pagerduty_service.by_name"),
				),
			},
			{
				Config: testAccDataSourcePagerDutyIntegrationConfigStep2(service, serviceIntegration),
				Check:  verifyOutput("output_id"),
			},
		},
	})
}

func verifyOutput(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Type != "string" {
			return fmt.Errorf("expected an error: %#v", rs)
		}

		return nil
	}
}

func testAccDataSourcePagerDutyIntegration(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a service ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the service %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyIntegrationConfigStep1(username, email, service, escalationPolicy, serviceIntegration string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_service" "test" {
  name                    = "%s"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.test.id
  alert_creation          = "create_incidents"
}

resource "pagerduty_escalation_policy" "test" {
  name        = "%s"
  num_loops   = 2
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_service_integration" "foo" {
  name    = "%s"
  service = pagerduty_service.test.id
  vendor  = data.pagerduty_vendor.datadog.id
}

data "pagerduty_vendor" "datadog" {
  name = "datadog"
}

data "pagerduty_service" "by_name" {
 name = pagerduty_service.test.name
}
`, username, email, service, escalationPolicy, serviceIntegration)
}

func testAccDataSourcePagerDutyIntegrationConfigStep2(service, serviceIntegration string) string {
	return fmt.Sprintf(`
data "pagerduty_service_integration" "service_integration" {
  service_name = "%s"
  integration_type = "%s"
}

output "output_id" {
  value = data.pagerduty_service_integration.service_integration.integration_key
}
`, service, serviceIntegration)
}
