package pagerduty

import (
	"fmt"
	"regexp"
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
	})
}

func TestAccPagerDutyEventOrchestrationPathService_Basic(t *testing.T) {
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resourceName := "pagerduty_event_orchestration_service.serviceA"
	serviceResourceName := "pagerduty_service.bar"

	// Checks that run on every step except the last one. These checks that verify the existance of the resource
	// and computed/default attributes. We're not checking individual resource attributes because
	// according to the official docs (https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource#TestCheckResourceAttr)
	// "State value checking is only recommended for testing Computed attributes and attribute defaults."
	baseChecks := []resource.TestCheckFunc{
		testAccCheckPagerDutyEventOrchestrationPathServiceExists(resourceName),
		testAccCheckPagerDutyEventOrchestrationPathServiceParent(resourceName, serviceResourceName),
		resource.TestCheckResourceAttr(resourceName, "type", "service"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationServicePathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceDefaultConfig(escalationPolicy, service),
				Check:  resource.ComposeTestCheckFunc(baseChecks...),
			},
			// Adding/updating/deleting automation_actions properties
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttrSet(resourceName, "sets.0.rules.0.id"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsParamsUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(
							resourceName, "sets.0.rules.0.actions.0.automation_actions.0.auto_send", "false",
						),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsParamsDeleteConfig(escalationPolicy, service),
				Check:  resource.ComposeTestCheckFunc(baseChecks...),
			},
			// Providing invalid extractions attributes for set rules
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, invalidExtractionRegexTemplateNilConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in sets.0.rules.0.actions.0.extractions.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, invalidExtractionRegexTemplateValConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in sets.0.rules.0.actions.0.extractions.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, invalidExtractionRegexNilSourceConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in sets.0.rules.0.actions.0.extractions.0: source can't be blank"),
			},
			// Providing invalid extractions attributes for the catch_all rule
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, "", invalidExtractionRegexTemplateNilConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extractions.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, "", invalidExtractionRegexTemplateValConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extractions.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, "", invalidExtractionRegexNilSourceConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extractions.0: source can't be blank"),
			},
			// Adding/updating/deleting all actions
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "sets.0.rules.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "sets.1.rules.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "sets.1.rules.1.id"),
						}...,
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "sets.0.rules.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "sets.1.rules.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "sets.1.rules.1.id"),
						}...,
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsDeleteConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "sets.0.rules.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "sets.1.rules.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "sets.1.rules.1.id"),
						}...,
					)...,
				),
			},
			// Deleting sets and the service path resource
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceOneSetNoActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttrSet(resourceName, "sets.0.rules.0.id"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceResourceDeleteConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathNotExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationServicePathDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_path_service" {
			continue
		}

		srv := s.RootModule().Resources["pagerduty_service.bar"]

		if _, _, err := client.EventOrchestrationPaths.Get(srv.Primary.ID, "service"); err == nil {
			return fmt.Errorf("Event Orchestration Service Path still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationPathServiceExists(rn string) resource.TestCheckFunc {
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

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationServicePathNotExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[rn]
		if ok {
			return fmt.Errorf("Event Orchestration Service Path is not deleted from the state: %s", rn)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationPathServiceParent(rn, sn string) resource.TestCheckFunc {
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

func testAccCheckPagerDutyEventOrchestrationPathServiceDefaultConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
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

func testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsConfig(ep, s string) string {
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

			catch_all {
				actions {
					automation_actions {
						name = "catch-all test"
						url = "https://catch-all-test.com"
						auto_send = true

						headers {
							key = "foo1"
							value = "bar1"
						}
						headers {
							key = "baz1"
							value = "buz1"
						}

						parameters {
							key = "source1"
							value = "orch1"
						}
						parameters {
							key = "region1"
							value = "us1"
						}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsParamsUpdateConfig(ep, s string) string {
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

			catch_all {
				actions {
					automation_actions {
						name = "catch-all test upd"
						url = "https://catch-all-test-upd.com"

						headers {
							key = "baz2"
							value = "buz2"
						}

						parameters {
							key = "source2"
							value = "orch2"
						}
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsParamsDeleteConfig(ep, s string) string {
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

			catch_all {
				actions {
					automation_actions {
						name = "catch-all test upd"
						url = "https://catch-all-test-upd.com"
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(ep, s, re, cae string) string {
	return fmt.Sprintf(
		"%s%s",
		createBaseServicePathConfig(ep, s),
		fmt.Sprintf(`resource "pagerduty_event_orchestration_service" "serviceA" {
				parent {
					id = pagerduty_service.bar.id
				}			
				sets {
					id = "start"
					rules {
						actions {
							%s
						}
					}
				}
				catch_all {
					actions {
						%s
					}
				}
			}
		`, re, cae),
	)
}

func invalidExtractionRegexTemplateNilConfig() string {
	return `
		extractions {
			target = "event.summary"
		}`
}

func invalidExtractionRegexTemplateValConfig() string {
	return `
		extractions {
			regex = ".*"
			template = "hi"
			target = "event.summary"
		}`
}

func invalidExtractionRegexNilSourceConfig() string {
	return `
		extractions {
			regex = ".*"
			target = "event.summary"
		}`
}

func testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsConfig(ep, s string) string {
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
						variables {
							name = "hostname"
							path = "event.source"
							type = "regex"
							value = "Source host: (.*)"
						}
						variables {
							name = "cpu_val"
							path = "event.custom_details.cpu"
							type = "regex"
							value = "(.*)"
						}
						extractions {
							target = "event.summary"
							template = "High CPU usage on {{variables.hostname}}"
						}
						extractions {
							regex = ".*"
							source = "event.group"
							target = "event.custom_details.message"
						}
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

			catch_all {
				actions {
					suspend = 120
					priority = "P0IN2KW"
					annotate = "Routed through an event orchestration - catch-all rule"
					pagerduty_automation_actions {
						action_id = "01CSB5SMOKCKVRI5GN0LJG7SMC"
					}
					severity = "warning"
					event_action = "trigger"
					variables {
						name = "user_id"
						path = "event.custom_details.user_id"
						type = "regex"
						value = "Source host: (.*)"
					}
					variables {
						name = "updated_at"
						path = "event.custom_details.updated_at"
						type = "regex"
						value = "(.*)"
					}
					extractions {
						target = "event.custom_details.message"
						template = "Last modified by {{variables.user_id}} on {{variables.updated_at}}"
					}
					extractions {
						regex = ".*"
						source = "event.custom_details.region"
						target = "event.group"
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsUpdateConfig(ep, s string) string {
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
						variables {
							name = "cpu_val_upd"
							path = "event.custom_details.cpu_upd"
							type = "regex"
							value = "CPU:(.*)"
						}
						extractions {
							regex = ".*"
							source = "event.custom_details.region_upd"
							target = "event.source"
						}
						extractions {
							target = "event.custom_details.message_upd"
							template = "[UPD] High CPU usage on {{variables.hostname}}: {{variables.cpu_val}}"
						}
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
						extractions {
							regex = ".*"
							source = "event.custom_details.region"
							target = "event.group"
						}
						extractions {
							regex = ".*"
							source = "event.custom_details.hostname"
							target = "event.source"
						}
					}
				}
			}

			catch_all {
				actions {
					suspend = 360
					suppress = true
					priority = "P0IN2KX"
					annotate = "[UPD] Routed through an event orchestration - catch-all rule"
					pagerduty_automation_actions {
						action_id = "01CSB5SMOKCKVRI5GN0LJG7SMD"
					}
					severity = "info"
					event_action = "resolve"
					variables {
						name = "updated_at_upd"
						path = "event.custom_details.updated_at"
						type = "regex"
						value = "UPD (.*)"
					}					
					extractions {
						regex = ".*"
						source = "event.custom_details.region_upd"
						target = "event.class"
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsDeleteConfig(ep, s string) string {
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

			catch_all {
				actions { }
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceOneSetNoActionsConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = pagerduty_service.bar.id
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1 updated"
					actions {}
				}
			}

			catch_all {
				actions { }
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceResourceDeleteConfig(ep, s string) string {
	return createBaseServicePathConfig(ep, s)
}
