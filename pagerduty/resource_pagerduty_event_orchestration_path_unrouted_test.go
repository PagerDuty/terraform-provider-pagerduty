package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_unrouted", &resource.Sweeper{
		Name: "unrouted",
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

func TestAccPagerDutyEventOrchestrationPathUnrouted_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	escalation_policy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationPathUnroutedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigNoRules(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "type", "unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.#", "0"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigWithConditions(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.0.conditions.0.expression", "event.summary matches part 'rds'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigWithMultipleRules(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.0.conditions.0.expression", "event.summary matches part 'rds'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.0.conditions.1.expression", "event.severity matches part 'warning'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.1.conditions.0.expression", "event.severity matches part 'info'"),
				),
			},

			// {
			// 	Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigDeleteAllRulesInSet(team, escalation_policy, service, orchestration),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
			// 		resource.TestCheckResourceAttr(
			// 			"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.#", "0"),
			// 		testAccCheckPagerDutyEventOrchestrationRouterPathRouteToMatch(
			// 			"pagerduty_event_orchestration_unrouted.unrouted", "pagerduty_service.bar", true),
			// 	),
			// },
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigDelete(team, escalation_policy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedNotExists("pagerduty_event_orchestration_unrouted.unrouted"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_path_unrouted" {
			continue
		}

		orch := s.RootModule().Resources["pagerduty_event_orchestration.orch"]

		if _, _, err := client.EventOrchestrationPaths.Get(orch.Primary.ID, "unrouted"); err == nil {
			return fmt.Errorf("Event Orchestration Path still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Event Orchestration Path Type is set")
		}

		orch := s.RootModule().Resources["pagerduty_event_orchestration.orch"]
		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.EventOrchestrationPaths.Get(orch.Primary.ID, "unrouted")

		if err != nil {
			return fmt.Errorf("Orchestration Path type not found: %v for orchestration %v", "unrouted", orch.Primary.ID)
		}
		if found.Type != "unrouted" {
			return fmt.Errorf("Event Orchrestration path not found: %v - %v", "unrouted", found)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedNotExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[rn]
		if ok {
			return fmt.Errorf("Event Orchestration Unrouted Path is not deleted: %s", rn)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigDelete(t, ep, s, o string) string {
	return fmt.Sprintf(createUnroutedBaseConfig(t, ep, s, o))
}

func createUnroutedBaseConfig(t, ep, s, o string) string {
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
			team {
				id = pagerduty_team.foo.id
			}
		}
	`, t, ep, s, o)
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigNoRules(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createUnroutedBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
			type = "unrouted"
			parent {
				id = pagerduty_event_orchestration.orch.id
			}
			sets {
				id = "start"
			}
			catch_all {
				actions { }
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigWithConditions(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createUnroutedBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
			type = "unrouted"
			parent {
				id = pagerduty_event_orchestration.orch.id
			}
			sets {
				id = "start"
				rules {
					disabled = false
					label = "rule1 label"
					actions { }
					conditions {
						expression = "event.summary matches part 'rds'"
					}
				}
			}
			catch_all {
				actions { }
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigWithMultipleRules(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createUnroutedBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
			type = "unrouted"
			parent {
				id = pagerduty_event_orchestration.orch.id
			}
			sets {
				id = "start"
				rules {
					disabled = false
					label = "rule1 label"
					actions { }
					conditions {
						expression = "event.summary matches part 'rds'"
					}
					conditions {
						expression = "event.severity matches part 'warning'"
					}
				}

				rules {
					disabled = false
					label = "rule2 label"
					actions { }
					conditions {
						expression = "event.severity matches part 'info'"
					}
				}
			}
			catch_all {
				actions { }
			}
		}
`)
}

// func testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigDeleteAllRulesInSet(t, ep, s, o string) string {
// 	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
// 		`
// resource "pagerduty_event_orchestration_router" "router" {
// 	type = "router"
// 	parent {
//         id = pagerduty_event_orchestration.orch.id
//     }
// 	catch_all {
// 		actions {
// 			route_to = pagerduty_service.bar.id
// 		}
// 	}
// 	sets {
// 		id = "start"
// 	}
// }
// `)
// }

// func testAccCheckPagerDutyEventOrchestrationPathConfig(t, ep, s, o string) string {
// 	return fmt.Sprintf("%s%s", createBaseConfig(t, ep, s, o),
// 		`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
// 			type = "unrouted"
// 			parent {
// 				id = pagerduty_event_orchestration.orch.id
// 			}
// 			sets {
// 				id = "start"
// 				rules {
// 					disabled = false
// 					label = "rule1 label"
// 					actions {
// 						route_to = "child-1"
// 						severity = "info"
// 						event_action = "trigger"
// 						variables {
// 							name = "server_name"
// 							path = "event.summary"
// 							type = "regex"
// 							value = "High CPU on (.*) server"
// 						}
// 						extractions {
// 							target = "event.summary"
// 							template = "High CPU on variables.hostname server"
// 						}
// 					}
// 				}
// 			}
// 			sets {
// 				id = "child-1"
// 				rules {
// 					disabled = false
// 					label = "rule2 label"
// 					actions {
// 						severity = "warning"
// 						event_action = "resolve"
// 						variables {
// 							name = "server_name"
// 							path = "event.summary"
// 							type = "regex"
// 							value = "High CPU on (.*) server"
// 						}
// 						extractions {
// 							target = "event.summary"
// 							template = "High CPU on event.custom_details.hostname server"
// 						}
// 					}
// 				}
// 			}
// 			catch_all {
// 				actions {

// 				}
// 			}
// 		}
// 	`)
// }

// func testAccCheckPagerDutyEventOrchestrationPathConfigWithConditions(t string, o string, ptype string) string {
// 	return fmt.Sprintf(`
// 		resource "pagerduty_team" "foo" {
// 			name = "%s"
// 		}

// 		resource "pagerduty_user" "foo" {
// 			name        = "user"
// 			email       = "user@pagerduty.com"
// 			color       = "green"
// 			role        = "user"
// 			job_title   = "foo"
// 			description = "foo"
// 		}

// 		resource "pagerduty_escalation_policy" "foo" {
// 			name        = "test"
// 			description = "bar"
// 			num_loops   = 2
// 			rule {
// 				escalation_delay_in_minutes = 10
// 				target {
// 					type = "user_reference"
// 					id   = pagerduty_user.foo.id
// 				}
// 			}
// 		}

// 		resource "pagerduty_service" "bar" {
// 			name = "barService"
// 			escalation_policy       = pagerduty_escalation_policy.foo.id
// 			incident_urgency_rule {
// 				type = "constant"
// 				urgency = "high"
// 			}
// 		}

// 		resource "pagerduty_event_orchestration" "orch" {
// 			name = "%s"
// 			team {
// 				id = pagerduty_team.foo.id
// 			}
// 		}

// 		resource "pagerduty_event_orchestration_unrouted" "unrouted" {
// 			type = "%s"
// 			parent {
// 				id = pagerduty_event_orchestration.orch.id
// 				type = "event_orchestration_reference"
// 				self = "https://api.pagerduty.com/event_orchestrations/orch_id"
// 			}
// 			sets {
// 				id = "start"
// 				rules {
// 					disabled = false
// 					label = "rule1 label"
// 					actions {
// 						route_to = pagerduty_service.bar.id
// 						severity = "info"
// 						event_action = "trigger"
// 						// variables {
// 						// 	name = "foo"
// 						// 	path = "foo.foo"
// 						// 	type = "foo"
// 						// 	value = "foo"
// 						// }
// 						// extractions {
// 						// 	target = "foo"
// 						// 	template = "foo"
// 						// }
// 					}
// 					conditions {
// 						expression = "event.summary matches part 'database'"
// 					}
// 				}
// 			}
// 		}
// 	`, t, o, ptype)
// }
