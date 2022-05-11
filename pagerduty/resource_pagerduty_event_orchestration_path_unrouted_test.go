package pagerduty

import (
	"fmt"
	// "strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO: Double check if sweeper works as expected
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
	orchestration := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchPathType := "unrouted"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathConfig(team, orchestration, orchPathType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "type", orchPathType),
					// resource.TestCheckResourceAttr(
					// 	"pagerduty_event_orchestration_unrouted.unrouted", "self", "https://api.pagerduty.com/event_orchestrations/orch_id/unrouted"),
				),
			},
			// {
			// 	Config: testAccCheckPagerDutyEventOrchestrationPathConfigWithConditions(team, orchestration, orchPathType),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckPagerDutyEventOrchestrationPathExists("pagerduty_event_orchestration_unrouted.unrouted"),
			// 		resource.TestCheckResourceAttr(
			// 			"pagerduty_event_orchestration_unrouted.unrouted", "sets.0.rules.0.conditions.0.expression", "event.summary matches part 'database'"),
			// 	),
			// },
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathDestroy(s *terraform.State) error {
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

func testAccCheckPagerDutyEventOrchestrationPathExists(rn string) resource.TestCheckFunc {
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

		resource "pagerduty_event_orchestration_unrouted" "unrouted" {
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
						route_to = "child-1"
						severity = "info"
						event_action = "trigger"
						variables {
							name = "server_name"
							path = "event.summary"
							type = "regex"
							value = "High CPU on (.*) server"
						}
						extractions {
							target = "event.summary"
							template = "High CPU on variables.hostname server"
						}
					}
				}
			}
			sets {
				id = "child-1"
				rules {
					disabled = false
					label = "rule2 label"
					actions {
						// route_to = pagerduty_service.bar.id
						severity = "warning"
						event_action = "resolve"
						variables {
							name = "server_name"
							path = "event.summary"
							type = "regex"
							value = "High CPU on (.*) server"
						}
						extractions {
							target = "event.summary"
							template = "High CPU on event.custom_details.hostname server"
						}
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

		resource "pagerduty_event_orchestration_unrouted" "unrouted" {
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
						severity = "info"
						event_action = "trigger"
						// variables {
						// 	name = "foo"
						// 	path = "foo.foo"
						// 	type = "foo"
						// 	value = "foo"
						// }
						// extractions {
						// 	target = "foo"
						// 	template = "foo"
						// }
					}
					conditions {
						expression = "event.summary matches part 'database'"
					}
				}
			}
		}
	`, t, o, ptype)
}
