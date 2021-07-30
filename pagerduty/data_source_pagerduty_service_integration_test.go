package pagerduty

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func setupResources(service, serviceIntegration, escalationPolicy, username, email string) error {
	client, err := pagerduty.NewClient(&pagerduty.Config{Token: os.Getenv("PAGERDUTY_TOKEN")})
	if err != nil {
		return err
	}

	createUser, _, err := client.Users.Create(&pagerduty.User{Name: username, Email: email})
	if err != nil {
		return err
	}

	createEp, _, err := client.EscalationPolicies.Create(&pagerduty.EscalationPolicy{
		Name:            escalationPolicy,
		EscalationRules: []*pagerduty.EscalationRule{{EscalationDelayInMinutes: 10, Targets: []*pagerduty.EscalationTargetReference{{Type: "user_reference", ID: createUser.ID}}}},
	})
	if err != nil {
		return err
	}

	createResponse, _, err := client.Services.Create(&pagerduty.Service{
		Name:             service,
		EscalationPolicy: &pagerduty.EscalationPolicyReference{ID: createEp.ID, Type: "escalation_policy_reference"},
	})
	if err != nil {
		return err
	}

	_, _, err = client.Services.CreateIntegration(createResponse.ID, &pagerduty.Integration{
		Summary: serviceIntegration,
		Type:    "service_integration_reference",
		Service: &pagerduty.ServiceReference{ID: createResponse.ID, Type: "service_reference"},
		Vendor:  &pagerduty.VendorReference{ID: "PAM4FGS", Type: "vendor_reference"},
	})
	if err != nil {
		return fmt.Errorf("error creating integration: %v", err)
	}

	return nil
}

func TestAccDataSourcePagerDutyIntegration_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceIntegration := "Datadog"

	err := setupResources(service, serviceIntegration, escalationPolicy, username, email)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
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

func testAccDataSourcePagerDutyIntegrationConfigStep2(service, serviceIntegration string) string {
	return fmt.Sprintf(`
data "pagerduty_service_integration" "service_integration" {
  service_name = "%s"
  integration_summary = "%s"
}

output "output_id" {
  value = data.pagerduty_service_integration.service_integration.integration_key
}
`, service, serviceIntegration)
}
