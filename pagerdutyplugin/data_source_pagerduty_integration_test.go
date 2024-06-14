package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyServiceIntegration_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := "Datadog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceIntegrationConfigStep1(service, serviceIntegration, email, escalationPolicy),
				Check: func(state *terraform.State) error {
					resource.Test(t, resource.TestCase{
						ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
						Steps: []resource.TestStep{
							{
								Config: testAccDataSourcePagerDutyServiceIntegrationConfigStep2(service, serviceIntegration),
								Check:  verifyOutput("output_id"),
							},
						},
					})
					return nil
				},
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

func testAccDataSourcePagerDutyServiceIntegrationConfigStep1(service, serviceIntegration, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "pagerduty_user" {
  email = "%s"
  name = "test user"
}

resource "pagerduty_escalation_policy" "escalation_policy" {
  name = "%s"
  rule {
    escalation_delay_in_minutes = 5
    target {
      type = "user_reference"
      id = pagerduty_user.pagerduty_user.id
    }
  }

}

resource "pagerduty_service" "pagerduty_service" {
  name = "%s"
  escalation_policy = pagerduty_escalation_policy.escalation_policy.id
}

resource "pagerduty_service_integration" "service_integration" {
  name    = "%s"
  service = pagerduty_service.pagerduty_service.id
  vendor  = data.pagerduty_vendor.datadog.id
}

data "pagerduty_vendor" "datadog" {
  name = "datadog"
}


`, email, escalationPolicy, service, serviceIntegration)
}

func testAccDataSourcePagerDutyServiceIntegrationConfigStep2(service, serviceIntegration string) string {
	return fmt.Sprintf(`

data "pagerduty_service_integration" "service_integration" {
 service_name = "%s"
 integration_summary = "%s"
}

output "output_id" {
 value     = data.pagerduty_service_integration.service_integration.integration_key
 sensitive = true
}
`, service, serviceIntegration)
}
