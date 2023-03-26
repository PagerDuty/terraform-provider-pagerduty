package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyService_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceConfig(username, email, service, escalationPolicy, teamname),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyService("pagerduty_service.test", "data.pagerduty_service.by_name"),
					resource.TestCheckResourceAttr("data.pagerduty_service.by_name", "teams.#", "1"),
					resource.TestCheckResourceAttr("data.pagerduty_service.by_name", "teams.0.name", teamname),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyService(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a service ID from PagerDuty")
		}

		testAtts := []string{"teams", "id", "name", "type", "auto_resolve_timeout", "acknowledgement_timeout", "alert_creation", "description", "escalation_policy"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the service %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyServiceConfig(username, email, service, escalationPolicy, teamname string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "test" {
	name        = "%s"
	description = "test"
}

resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team_membership" "test" {
	team_id = pagerduty_team.test.id
	user_id = pagerduty_user.test.id
}

resource "pagerduty_escalation_policy" "test" {
  name        = "%s"
  num_loops   = 2
  teams       = [pagerduty_team.test.id]
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
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.test.id
  alert_creation          = "create_incidents"
}

data "pagerduty_service" "by_name" {
  depends_on = [pagerduty_team_membership.test, pagerduty_service.test]
  name = pagerduty_service.test.name
}
`, teamname, username, email, service, escalationPolicy)
}
