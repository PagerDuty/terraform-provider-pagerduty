package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyStandardsResourcesScores_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyStandardsResourcesScoresConfig(name),
				Check: testAccCheckAttributes(
					fmt.Sprintf("data.pagerduty_standards_resources_scores.%s", name),
					testStandardsResourcesScores,
				),
			},
		},
	})
}

func testStandardsResourcesScores(a map[string]string) error {
	testAttrs := []string{
		"ids.#",
		"ids.0",
		"resource_type",
		"resources.#",
		"resources.0.resource_id",
		"resources.0.resource_type",
		"resources.0.score.passing",
		"resources.0.score.total",
		"resources.0.standards.#",
		"resources.0.standards.0.active",
		"resources.0.standards.0.description",
		"resources.0.standards.0.id",
		"resources.0.standards.0.name",
		"resources.0.standards.0.pass",
		"resources.0.standards.0.type",
	}

	for _, attr := range testAttrs {
		if _, ok := a[attr]; !ok {
			return fmt.Errorf("Expected the required attribute %s to exist", attr)
		}
	}

	return nil
}

func testAccDataSourcePagerDutyStandardsResourcesScoresConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "Earline Greenholt"
  email = "125.greenholt.earline@graham.name"
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "Testing Escalation Policy"
  num_loops = 2
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_service" "example" {
  name                    = "My Web App test"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.bar.id
  alert_creation          = "create_alerts_and_incidents"
  auto_pause_notifications_parameters {
    enabled = true
    timeout = 300
  }
}

data "pagerduty_standards_resources_scores" "%s" {
  resource_type = "technical_services"
  ids           = [pagerduty_service.example.id]
}`, name)
}
