package pagerduty

import (
	"fmt"
	// "strconv"
	// "log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_service", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_service",
		F:    testSweepEventOrchestration,
		Dependencies: []string{
			"pagerduty_user",
			"pagerduty_escalation_policy",
			"pagerduty_service",
		},
	})
}

func TestAccPagerDutyEventOrchestrationPathService_Basic(t *testing.T) {
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resourceName := "pagerduty_event_orchestration_service.serviceA"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// TODO:
		// CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceRequiredFieldsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
					resource.TestCheckResourceAttr(
						resourceName, "type", "service",
					),
				),
			},
			// Test adding/updating/deleting automation_actions properties
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAutomationActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
					resource.TestCheckResourceAttrSet(
						resourceName, "sets.0.rules.0.id",
					),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAutomationActionsParamsUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
					resource.TestCheckResourceAttr(
						resourceName, "sets.0.rules.0.actions.0.automation_actions.0.auto_send", "false",
					),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAutomationActionsParamsDeleteConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
				),
			},
			// Test adding/updating extractions/variables
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceExtractionsVariablesConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceExtractionsVariablesUpdatedConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
				),
			},
			// Test adding/updating/deleting all actions
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAllActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAllActionsUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAllActionsDeleteConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathExists(resourceName),
					testAccCheckPagerDutyEventOrchestrationServicePathParent(resourceName, "pagerduty_service.bar"),
				),
			},

			// test route to
			// reset rule action items -> should be default
			// reset rule actions -> should be []
			// reset rule conditions -> should be []
			// reset rules
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationServicePathExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		orch, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if orch.Primary.ID == "" {
			return fmt.Errorf("No Event Orchestration Service ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.EventOrchestrationPaths.Get(orch.Primary.ID, "service")
		if err != nil {
			return err
		}
		if found.Parent.ID != orch.Primary.ID {
			return fmt.Errorf("Event Orchrestration Service not found: %v - %v", orch.Primary.ID, found)
		}

		// return fmt.Errorf(">>> attr: %v", orch.Primary.Attributes)

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationServicePathParent(rn, sn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		p, _ := s.RootModule().Resources[rn]
		srv, ok := s.RootModule().Resources[sn]

		if !ok {
			return fmt.Errorf("Service not found: %s", sn)
		}

		var pId = p.Primary.Attributes["parent.0.id"]
		var sId = srv.Primary.Attributes["id"]
		if pId != sId {
			return fmt.Errorf("Event Orchestration Service path parent ID (%v) not matching provided service ID: %v", pId, sId)
		}

		var t = p.Primary.Attributes["parent.0.type"]
		if t != "service_reference" {
			return fmt.Errorf("Event Orchestration Service path parent type (%v) not matching expected type: 'service_reference'", t)
		}

		var self = p.Primary.Attributes["parent.0.self"]
		if !strings.HasSuffix(self, sId) {
			return fmt.Errorf("Event Orchestration Service path parent self URL (%v) not containing expected service ID", self)
		}

		return nil
	}
}

func createBaseServicePathConfig(ep, s string) string {
	return fmt.Sprintf(`
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
	`, ep, s)
}

func testAccCheckPagerDutyEventOrchestrationServiceRequiredFieldsConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAutomationActionsConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					actions {							
							automation_actions {
								name = "test"
								url = "https://test.com"
								auto_send = true
		
								headers {
									key = "foo"
									value = "bar"
								}
								headers {
									key = "baz"
									value = "buz"
								}
		
								parameters {
									key = "source"
									value = "orch"
								}
								parameters {
									key = "region"
									value = "us"
								}
							}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAutomationActionsParamsUpdateConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					actions {							
							automation_actions {
								name = "test1"
								url = "https://test1.com"
		
								headers {
									key = "foo1"
									value = "bar1"
								}
								parameters {
									key = "source_region"
									value = "eu"
								}
							}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAutomationActionsParamsDeleteConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					actions {							
							automation_actions {
								name = "test"
								url = "https://test.com"
							}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceExtractionsVariablesConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					actions {							
						variables {
							name = "server_name_cpu"
							path = "event.summary"
							type = "regex"
							value = "High CPU on (.*) server"
						}
						variables {
							name = "server_name_memory"
							path = "event.custom_details"
							type = "regex"
							value = "High memory usage on (.*) server"
						}
						extractions {
							target = "event.summary"
							template = "High memory usage on variables.hostname server"
						}
						extractions {
							target = "event.custom_details"
							template = "High memory usage on variables.hostname server"
						}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceExtractionsVariablesUpdatedConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					actions {
						variables {
							name = "host_name"
							path = "event.custom_details.memory"
							type = "regex"
							value = "High memory usage on (.*) server"
						}
						extractions {
							target = "event.custom_details.info"
							template = "High memory usage on {{variables.hostname}} server: {{event.custom_details.max_memory}}"
						}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAllActionsConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					conditions {
						expression = "event.summary matches part 'timeout'"
					}
					conditions {
						expression = "event.custom_details.timeout_err exists"
					}
					actions {
						route_to = "set-1"
						priority = "P0IN2KQ"
						annotate = "Routed through an event orchestration"
						pagerduty_automation_actions {
							action_id = "01CSB5SMOKCKVRI5GN0LJG7SMB"
						}
						severity = "critical"
						event_action = "trigger"
					}
				}
			}
			sets {
				id = "set-1"
				rules {
					label = "set-1 rule 1"
					actions {
						suspend = 300
					}
				}
				rules {
					label = "set-1 rule 2"
					conditions {
						expression = "event.source matches part 'stg-'"
					}
					actions {
						suppress = true
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAllActionsUpdateConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1 updated"
					conditions {
						expression = "event.custom_details.timeout_err matches part 'timeout'"
					}
					actions {
						route_to = "set-2"
						priority = "P0IN2KR"
						annotate = "Routed through a service orchestration!"
						pagerduty_automation_actions {
							action_id = "01CSB5SMOKCKVRI5GN0LJG7SMBUPDATED"
						}
						severity = "warning"
						event_action = "resolve"
					}
				}
			}
			sets {
				id = "set-2"
				rules {
					label = "set-2 rule 1"
					actions {
						suspend = 15
					}
				}
				rules {
					label = "set-2 rule 2"
					conditions {
						expression = "event.source matches part 'test-'"
					}
					actions {
						annotate = "Matched set-2 rule 2"
						variables {
							name = "host_name"
							path = "event.custom_details.memory"
							type = "regex"
							value = "High memory usage on (.*) server"
						}
						extractions {
							target = "event.summary"
							template = "High memory usage on {{variables.hostname}} server: {{event.custom_details.max_memory}}"
						}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAllActionsDeleteConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1 updated"
					actions {
						route_to = "set-2"
					}
				}
			}
			sets {
				id = "set-2"
				rules {
					label = "set-2 rule 1"
					actions { }
				}
				rules {
					label = "set-2 rule 2"
					actions { }
				}
			}
		}
	`)
}
