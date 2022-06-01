package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_router", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_router",
		F:    testSweepEventOrchestration,
	})
}

func TestAccPagerDutyEventOrchestrationPathRouter_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigNoRules(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "unrouted", true), //test for catch_all route_to prop, by default it should be unrouted
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.#", "0"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfig(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "pagerduty_service.bar", false), // test for rule action route_to
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "unrouted", true), //test for catch_all route_to prop, by default it should be unrouted
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigWithConditions(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.0.condition.0.expression", "event.summary matches part 'database'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigWithMultipleRules(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.0.condition.0.expression", "event.summary matches part 'database'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.1.condition.0.expression", "event.severity matches part 'critical'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigWithCatchAllToService(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.#", "1"),
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "pagerduty_service.bar", true), //test for catch_all routing to service if provided
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigNoConditions(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.0.condition.#", "0"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigDeleteAllRulesInSet(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.#", "0"),
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "pagerduty_service.bar", true),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigDelete(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterNotExists("pagerduty_event_orchestration_router.router"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationRouterDestroy(s *terraform.State) error {
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

func testAccCheckPagerDutyEventOrchestrationRouterExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Event Orchestration Router is set")
		}

		orch, _ := s.RootModule().Resources["pagerduty_event_orchestration.orch"]
		client, _ := testAccProvider.Meta().(*Config).Client()
		_, _, err := client.EventOrchestrationPaths.Get(orch.Primary.ID, "router")

		if err != nil {
			return fmt.Errorf("Orchestration Path type not found: %v for orchestration %v", "router", orch.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationRouterNotExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[rn]
		if ok {
			return fmt.Errorf("Event Orchestration Router Path is not deleted: %s", rn)
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
		name        = "tf-user"
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
		team = pagerduty_team.foo.id
	}
	`, t, ep, s, o)
}

func testAccCheckPagerDutyEventOrchestrationRouterConfigNoRules(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			catch_all {
				actions {
					route_to = "unrouted"
				}
			}
			set {
				id = "start"
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationRouterConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			catch_all {
				actions {
					route_to = "unrouted"
				}
			}
			set {
				id = "start"
				rule {
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

func testAccCheckPagerDutyEventOrchestrationRouterConfigWithConditions(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			catch_all {
				actions {
					route_to = "unrouted"
				}
			}
			set {
				id = "start"
				rule {
					disabled = false
					label = "rule1 label"
					actions {
						route_to = pagerduty_service.bar.id
					}
					condition {
						expression = "event.summary matches part 'database'"
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationRouterConfigWithMultipleRules(t, ep, s, o string) string {
	return fmt.Sprintf(
		"%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_service" "bar2" {
			name = "tf-barService2"
			escalation_policy       = pagerduty_escalation_policy.foo.id

			incident_urgency_rule {
				type = "constant"
				urgency = "high"
			}
		}

		resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			catch_all {
				actions {
					route_to = "unrouted"
				}
			}
			set {
				id = "start"
				rule {
					disabled = false
					label = "rule1 label"
					actions {
						route_to = pagerduty_service.bar.id
					}
					condition {
						expression = "event.summary matches part 'database'"
					}
					condition {
						expression = "event.severity matches part 'critical'"
					}
				}

				rule {
					disabled = false
					label = "rule2 label"
					actions {
						route_to = pagerduty_service.bar2.id
					}
					condition {
						expression = "event.severity matches part 'critical'"
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationRouterConfigNoConditions(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			catch_all {
				actions {
					route_to = pagerduty_service.bar.id
				}
			}
			set {
				id = "start"
				rule {
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

func testAccCheckPagerDutyEventOrchestrationRouterConfigWithCatchAllToService(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			
			catch_all {
				actions {
					route_to = pagerduty_service.bar.id
				}
			}
			set {
				id = "start"
				rule {
					disabled = false
					label = "rule1 label"
					actions {
						route_to = pagerduty_service.bar.id
					}
					condition {
						expression = "event.severity matches part 'critical'"
					}
				}
			}
		}
		`)
}

func testAccCheckPagerDutyEventOrchestrationRouterConfigDeleteAllRulesInSet(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_router" "router" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			
			catch_all {
				actions {
					route_to = pagerduty_service.bar.id
				}
			}
			set {
				id = "start"
			}
		}
		`)
}

func testAccCheckPagerDutyEventOrchestrationRouterConfigDelete(t, ep, s, o string) string {
	return createBaseConfig(t, ep, s, o)
}

func testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(router, service string, catchAll bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, rOk := s.RootModule().Resources[router]
		if !rOk {
			return fmt.Errorf("Not found: %s", router)
		}

		var rRouteToId = ""
		if catchAll == true {
			rRouteToId = r.Primary.Attributes["catch_all.0.actions.0.route_to"]
		} else {
			rRouteToId = r.Primary.Attributes["set.0.rule.0.actions.0.route_to"]
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
