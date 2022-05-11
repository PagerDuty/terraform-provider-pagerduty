package pagerduty

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_router", &resource.Sweeper{
		Name: "router",
		Dependencies: []string{
			"pagerduty_schedule",
			"pagerduty_team",
			"pagerduty_user",
			"pagerduty_escalation_policy",
			"pagerduty_service",
			"pagerduty_event_orchestration",
		},
	})
}

func TestAccPagerDutyEventOrchestrationPathRouter_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	escalation_policy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfig(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "type", "router"),
					testAccCheckPagerDutyEventOrchestrationPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "pagerduty_service.bar", false), // test for rule action route_to
					testAccCheckPagerDutyEventOrchestrationPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "unrouted", true), //test for catch_all route_to prop, by default it should be unrouted
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfigWithConditions(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "sets.0.rules.0.conditions.0.expression", "event.summary matches part 'database'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfigWithMultipleRules(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "sets.0.rules.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "sets.0.rules.0.conditions.0.expression", "event.summary matches part 'database'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "sets.0.rules.1.conditions.0.expression", "event.severity matches part 'critical'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfigWithCatchAllToService(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "sets.0.rules.#", "1"),
					testAccCheckPagerDutyEventOrchestrationPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "pagerduty_service.bar", true), //test for catch_all routing to service if provided
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

func createBaseConfig(t, ep, s, o string) string {
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
		name        = "%s"
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
		name = "%s"
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
	`, t, ep, s, o)
}

func testAccCheckPagerDutyEventOrchestrationPathConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
		type = "router"
		parent {
			id = pagerduty_event_orchestration.orch.id
			type = "event_orchestration_reference"
			self = "https://api.pagerduty.com/event_orchestrations/orch_id"
		}
		catch_all {
			actions {
				route_to = "unrouted"
			}
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
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathConfigWithConditions(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`
resource "pagerduty_event_orchestration_router" "router" {
	type = "router"
	parent {
        id = pagerduty_event_orchestration.orch.id
		type = "event_orchestration_reference"
		self = "https://api.pagerduty.com/event_orchestrations/orch_id"
    }
	catch_all {
		actions {
			route_to = "unrouted"
		}
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
`)
}

func testAccCheckPagerDutyEventOrchestrationPathConfigWithMultipleRules(t, ep, s, o string) string {
	return fmt.Sprintf(
		"%s%s", createBaseConfig(t, ep, s, o),
		`
resource "pagerduty_service" "bar2" {
	name = "tf-barService2"
	escalation_policy       = pagerduty_escalation_policy.foo.id

	incident_urgency_rule {
		type = "constant"
		urgency = "high"
	}
}

resource "pagerduty_event_orchestration_router" "router" {
	type = "router"
	parent {
        id = pagerduty_event_orchestration.orch.id
		type = "event_orchestration_reference"
		self = "https://api.pagerduty.com/event_orchestrations/orch_id"
    }
	catch_all {
		actions {
			route_to = "unrouted"
		}
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
			conditions {
				expression = "event.severity matches part 'critical'"
			}
		}

		rules {
			disabled = false
			label = "rule2 label"
			actions {
				route_to = pagerduty_service.bar2.id
			}
			conditions {
				expression = "event.severity matches part 'critical'"
			}
		}
	}
}
`)
}

func testAccCheckPagerDutyEventOrchestrationPathConfigWithCatchAllToService(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`
resource "pagerduty_event_orchestration_router" "router" {
	type = "router"
	parent {
        id = pagerduty_event_orchestration.orch.id
		type = "event_orchestration_reference"
		self = "https://api.pagerduty.com/event_orchestrations/orch_id"
    }
	catch_all {
		actions {
			route_to = pagerduty_service.bar.id
		}
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
				expression = "event.severity matches part 'critical'"
			}
		}
	}
}
`)
}

func testAccCheckPagerDutyEventOrchestrationPathRouteToMatch(router, service string, catchAll bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, rOk := s.RootModule().Resources[router]
		if !rOk {
			return fmt.Errorf("Not found: %s", router)
		}

		var rRouteToId = ""
		if catchAll == true {
			rRouteToId = r.Primary.Attributes["catch_all.0.actions.0.route_to"]
		} else {
			rRouteToId = r.Primary.Attributes["sets.0.rules.0.actions.0.route_to"]
		}

		var sId = ""
		if service == "unrouted" {
			sId = "unrouted"
		} else {
			svc, sOk := s.RootModule().Resources[service]
			if !sOk {
				return fmt.Errorf("Not found: %s", service)
			}
			sId = svc.Primary.Attributes["id"]
		}

		if rRouteToId != sId {
			return fmt.Errorf("Event Orchestration Router Route to Service ID (%v) not matching provided service ID: %v", rRouteToId, sId)
		}

		return nil
	}
}
