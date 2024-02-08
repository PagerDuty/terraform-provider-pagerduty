package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyStandardsResourceScores_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyStandardsResourceScoresConfig(name),
				Check: testAccCheckAttributes(
					fmt.Sprintf("data.pagerduty_standards_resource_scores.%s", name),
					testStandardsResourceScores,
				),
			},
		},
	})
}

func testStandardsResourceScores(a map[string]string) error {
	testAttrs := []string{
		"id",
		"resource_type",
		"score.passing",
		"score.total",
		"standards.#",
		"standards.0.active",
		"standards.0.description",
		"standards.0.id",
		"standards.0.name",
		"standards.0.pass",
		"standards.0.type",
	}
	for _, attr := range testAttrs {
		if _, ok := a[attr]; !ok {
			return fmt.Errorf("Expected the required attribute %s to exist", attr)
		}
	}
	return nil
}

func testAccDataSourcePagerDutyStandardsResourceScoresConfig(name string) string {
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

data "pagerduty_standards_resource_scores" "%s" {
  resource_type = "technical_services"
  id            = pagerduty_service.example.id
}`, name)
}
