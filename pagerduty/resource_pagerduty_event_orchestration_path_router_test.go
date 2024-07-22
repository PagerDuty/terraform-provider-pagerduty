package pagerduty

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
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

	dynamicRouteToByNameInput := &pagerduty.EventOrchestrationPathDynamicRouteTo{
		LookupBy: "service_name",
		Regex:    ".*",
		Source:   "event.custom_details.pd_service_name",
	}
	dynamicRouteToByIDInput := &pagerduty.EventOrchestrationPathDynamicRouteTo{
		LookupBy: "service_id",
		Regex:    "ID:(.*)",
		Source:   "event.custom_details.pd_service_id",
	}
	invalidDynamicRouteToPlacementMessage := "Invalid Dynamic Routing rule configuration:\n- A Router can have at most one Dynamic Routing rule; Rules with the dynamic_route_to action found at indexes: 1, 2\n- The Dynamic Routing rule must be the first rule in a Router"
	invalidDynamicRouteToConfigMessage := "Invalid Dynamic Routing rule configuration:\n- Dynamic Routing rules cannot have conditions\n- Dynamic Routing rules cannot have the `route_to` action"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationRouterDestroy,
		Steps: []resource.TestStep{
			// Invalid Dynamic Routing rule config for a new resource: multiple Dynamic Routing rules, Dynamic Routing rule not the first rule in the Router:
			{
				Config:      testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRulePlacement(team, escalationPolicy, service, orchestration),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(invalidDynamicRouteToPlacementMessage),
			},
			// Invalid Dynamic Routing rule config for a new resource: Dynamic Routing rule with conditions and the interpolated route_to action:
			{
				Config:      testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRuleConfig(team, escalationPolicy, service, orchestration, "pagerduty_service.bar.id"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(invalidDynamicRouteToConfigMessage),
			},
			// Invalid Dynamic Routing rule config for a new resource: Dynamic Routing rule with conditions and the hard-coded route_to action:
			{
				Config:      testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRuleConfig(team, escalationPolicy, service, orchestration, "\"PARASOL\""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(invalidDynamicRouteToConfigMessage),
			},
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
			// Configure a Dynamic Routing rule:
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterDynamicRouteToConfig(team, escalationPolicy, service, orchestration, dynamicRouteToByNameInput),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					testAccCheckPagerDutyEventOrchestrationRouterPathDynamicRouteToMatch("pagerduty_event_orchestration_router.router", dynamicRouteToByNameInput),
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "unrouted", true),
				),
			},
			// Invalid Dynamic Routing rule config for an existing resource: multiple Dynamic Routing rules, Dynamic Routing rule not the first rule in the Router:
			{
				Config:      testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRulePlacement(team, escalationPolicy, service, orchestration),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(invalidDynamicRouteToPlacementMessage),
			},
			// Invalid Dynamic Routing rule config for an existing resource: Dynamic Routing rule with conditions and the interpolated route_to action:
			{
				Config:      testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRuleConfig(team, escalationPolicy, service, orchestration, "pagerduty_service.bar.id"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(invalidDynamicRouteToConfigMessage),
			},
			// Invalid Dynamic Routing rule config for an existing resource: Dynamic Routing rule with conditions and the hard-coded route_to action:
			{
				Config:      testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRuleConfig(team, escalationPolicy, service, orchestration, "\"PARASOL\""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(invalidDynamicRouteToConfigMessage),
			},
			// Update the Dynamic Routing rule:
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterDynamicRouteToConfig(team, escalationPolicy, service, orchestration, dynamicRouteToByIDInput),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					testAccCheckPagerDutyEventOrchestrationRouterPathDynamicRouteToMatch("pagerduty_event_orchestration_router.router", dynamicRouteToByIDInput),
					testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
						"pagerduty_event_orchestration_router.router", "unrouted", true),
				),
			},
			// Delete the Dynamic Routing rule:
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

func TestAccPagerDutyEventOrchestrationPathRouter_EnableRoutingRule(t *testing.T) {
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
				Config: testAccCheckPagerDutyEventOrchestrationRouterEnableRoutingRuleConfig(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.0.disabled", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterEnableRoutingRuleConfigUpdated(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationRouterExists("pagerduty_event_orchestration_router.router"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_router.router", "set.0.rule.0.disabled", "false"),
				),
				// This is unnecessary, because this is the default behaviour of all
				// tests, it is only here to explicitely state that this is the expected
				// outcome from test.
				ExpectNonEmptyPlan: false,
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

		if _, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "router"); err == nil {
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
		_, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "router")

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

func testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRulePlacement(t, ep, s, o string) string {
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
					label = "static routing rule 1"
					actions {
						route_to = pagerduty_service.bar.id
					}
				}
				rule {
					disabled = false
					label = "dynamic routing rule 1"
					actions {
						dynamic_route_to {
							lookup_by = "service_id"
							regex = ".*"
							source = "event.custom_details.pd_service_id"
						}
					}
				}
				rule {
					label = "dynamic routing rule 2"
					actions {
						dynamic_route_to {
							lookup_by = "service_name"
							regex = ".*"
							source = "event.custom_details.pd_service_name"
						}
					}
				}
				rule {
					label = "static routing rule 2"
					actions {
						route_to = "P1B2C23"
					}
				}
				
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationRouterInvalidDynamicRoutingRuleConfig(t, ep, s, o, routeTo string) string {
	routerConfig := fmt.Sprintf(
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
					label = "dynamic routing rule 1"
					condition {
						expression = "event.summary matches part 'production'"
					}
					actions {
						dynamic_route_to {
							lookup_by = "service_id"
							regex = ".*"
							source = "event.custom_details.pd_service_id"
						}
						route_to = %s
					}
				}
				rule {
					label = "static routing rule 1"
					actions {
						route_to = "P1B2C23"
					}
				}
				
			}
		}
	`, routeTo)
	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o), routerConfig)
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

func testAccCheckPagerDutyEventOrchestrationRouterDynamicRouteToConfig(t, ep, s, o string, dynamicRouteToByNameInput *pagerduty.EventOrchestrationPathDynamicRouteTo) string {
	routerConfig := fmt.Sprintf(
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
					label = "dynamic routing rule"
					actions {
						dynamic_route_to {
							lookup_by = "%s"
							regex = "%s"
							source = "%s"
						}
					}
				}
				rule {
					label = "static routing rule"
					actions {
						route_to = pagerduty_service.bar.id
					}
				}
			}
		}
	`, dynamicRouteToByNameInput.LookupBy, dynamicRouteToByNameInput.Regex, dynamicRouteToByNameInput.Source)

	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o), routerConfig)
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

func testAccCheckPagerDutyEventOrchestrationRouterEnableRoutingRuleConfig(t, ep, s, o string) string {
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
					disabled = true
					label = "rule1 label"
					actions {
						route_to = pagerduty_service.bar.id
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationRouterEnableRoutingRuleConfigUpdated(t, ep, s, o string) string {
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

func testAccCheckPagerDutyEventOrchestrationRouterPathDynamicRouteToMatch(router string, expectedDynamicRouteTo *pagerduty.EventOrchestrationPathDynamicRouteTo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, rOk := s.RootModule().Resources[router]
		if !rOk {
			return fmt.Errorf("Not found: %s", router)
		}

		rLookupBy := r.Primary.Attributes["set.0.rule.0.actions.0.dynamic_route_to.0.lookup_by"]
		rRegex := r.Primary.Attributes["set.0.rule.0.actions.0.dynamic_route_to.0.regex"]
		rSource := r.Primary.Attributes["set.0.rule.0.actions.0.dynamic_route_to.0.source"]

		if rLookupBy != expectedDynamicRouteTo.LookupBy {
			return fmt.Errorf("Event Orchestration Router `dynamic_route_to.lookup_by` (%v) does not match expected value: %v", rLookupBy, expectedDynamicRouteTo.LookupBy)
		}
		if rRegex != expectedDynamicRouteTo.Regex {
			return fmt.Errorf("Event Orchestration Router `dynamic_route_to.regex` (%v) does not match expected value: %v", rRegex, expectedDynamicRouteTo.Regex)
		}
		if rSource != expectedDynamicRouteTo.Source {
			return fmt.Errorf("Event Orchestration Router `dynamic_route_to.source` (%v) does not match expected value: %v", rSource, expectedDynamicRouteTo.Source)
		}

		return nil
	}
}
