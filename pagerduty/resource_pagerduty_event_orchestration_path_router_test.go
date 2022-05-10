package pagerduty

import (
	"fmt"
	// "strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPagerDutyEventOrchestrationPathRouter_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchPathType := fmt.Sprintf("router")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfig(team, orchestration, orchPathType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "type", orchPathType),
					// resource.TestCheckResourceAttr(
					// 	"pagerduty_event_orchestration_router.router", "self", "https://api.pagerduty.com/event_orchestrations/orch_id/router"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfigWithConditions(team, orchestration, orchPathType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "sets.0.rules.0.conditions.0.expression", "event.summary matches part 'database'"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_path_router" {
			continue
		}

		orch, _ := s.RootModule().Resources["pagerduty_event_orchestration.orch"]

		if _, _, err := client.EventOrchestrationPaths.Get(orch.Primary.ID, "router"); err == nil {
			return fmt.Errorf("Event Orchestration Path still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationPathExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Event Orchestration Path Type is set")
		}

		orch, _ := s.RootModule().Resources["pagerduty_event_orchestration.orch"]
		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.EventOrchestrationPaths.Get(orch.Primary.ID, "router")

		if err != nil {
			return fmt.Errorf("Orchestration Path type not found: %v for orchestration %v", "router", orch.Primary.ID)
		}
		if found.Type != "router" {
			return fmt.Errorf("Event Orchrestration path not found: %v - %v", "router", found)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationPathConfig(t string, o string, ptype string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
	name = "%s"
}

resource "pagerduty_user" "foo" {
	name        = "user"
	email       = "user@pagerduty.com"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "test"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "bar" {
	name = "barService"
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "constant"
		urgency = "high"
	}
}

resource "pagerduty_event_orchestration" "orch" {
	name = "%s"
	team {
		id = pagerduty_team.foo.id
	}
}

resource "pagerduty_event_orchestration_router" "router" {
	type = "%s"
	parent {
        id = pagerduty_event_orchestration.orch.id
		type = "event_orchestration_reference"
		self = "https://api.pagerduty.com/event_orchestrations/orch_id"
    }
	sets {
		id = "start"
		rules {
			disabled = false
			label = "rule1 label"
			actions {
				route_to = pagerduty_service.bar.id
			}
		}
	}
}
`, t, o, ptype)
}

func testAccCheckPagerDutyEventOrchestrationPathConfigWithConditions(t string, o string, ptype string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
	name = "%s"
}

resource "pagerduty_user" "foo" {
	name        = "user"
	email       = "user@pagerduty.com"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "test"
	description = "bar"
	num_loops   = 2

	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}

resource "pagerduty_service" "bar" {
	name = "barService"
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "constant"
		urgency = "high"
	}
}

resource "pagerduty_event_orchestration" "orch" {
	name = "%s"
	team {
		id = pagerduty_team.foo.id
	}
}

resource "pagerduty_event_orchestration_router" "router" {
	type = "%s"
	parent {
        id = pagerduty_event_orchestration.orch.id
		type = "event_orchestration_reference"
		self = "https://api.pagerduty.com/event_orchestrations/orch_id"
    }
	sets {
		id = "start"
		rules {
			disabled = false
			label = "rule1 label"
			actions {
				route_to = pagerduty_service.bar.id
			}
			conditions {
				expression = "event.summary matches part 'database'"
			}
		}
	}
}
`, t, o, ptype)
}
