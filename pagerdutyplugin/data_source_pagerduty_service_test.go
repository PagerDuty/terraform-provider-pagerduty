package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyService_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceConfig(username, email, service, escalationPolicy, teamname),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyService("pagerduty_service.no_team_service", "data.pagerduty_service.no_team_service"),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyService_HasNoTeam(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceConfig(username, email, service, escalationPolicy, teamname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pagerduty_service.no_team_service", "teams.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyService_HasOneTeam(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyServiceConfig(username, email, service, escalationPolicy, teamname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pagerduty_service.one_team_service", "teams.#", "1"),
					resource.TestCheckResourceAttr("data.pagerduty_service.one_team_service", "teams.0.name", teamname),
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

		testAtts := []string{"id", "name", "type", "auto_resolve_timeout", "acknowledgement_timeout", "alert_creation", "description", "escalation_policy"}

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
resource "pagerduty_team" "team_one" {
  name        = "%s"
  description = "team_one"
}

resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team_membership" "team_membership_one" {
  team_id = pagerduty_team.team_one.id
  user_id = pagerduty_user.test.id
}

resource "pagerduty_escalation_policy" "no_team_ep" {
  name        = "no_team_ep"
  num_loops   = 2
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_escalation_policy" "one_team_ep" {
  depends_on = [pagerduty_team_membership.team_membership_one]
  name        = "%s"
  num_loops   = 2
  teams       = [pagerduty_team.team_one.id]
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_service" "no_team_service" {
  name                    = "no_team_service"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.no_team_ep.id
}

resource "pagerduty_service" "one_team_service" {
  name                    = "%s"
  auto_resolve_timeout    = 14400
  acknowledgement_timeout = 600
  escalation_policy       = pagerduty_escalation_policy.one_team_ep.id
}

data "pagerduty_service" "no_team_service" {
  name = pagerduty_service.no_team_service.name
}

data "pagerduty_service" "one_team_service" {
  name = pagerduty_service.one_team_service.name
}

`, teamname, username, email, service, escalationPolicy)
}
