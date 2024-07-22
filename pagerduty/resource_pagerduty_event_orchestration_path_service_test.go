package pagerduty

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	// according to the official docs (https://pkg.go.dev/github.com/hashicorp/terraform-plugin-testing/helper/resource#TestCheckResourceAttr)
	// "State value checking is only recommended for testing Computed attributes and attribute defaults."
	baseChecks := []resource.TestCheckFunc{
		testAccCheckPagerDutyEventOrchestrationPathServiceExists(resourceName),
		testAccCheckPagerDutyEventOrchestrationPathServiceServiceID(resourceName, serviceResourceName),
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
			// Adding/updating/deleting automation_action properties
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttrSet(resourceName, "set.0.rule.0.id"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAutomationActionsParamsUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(
							resourceName, "set.0.rule.0.actions.0.automation_action.0.auto_send", "false",
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
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, invalidExtractionRegexTemplateValConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, invalidExtractionRegexNilSourceConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: source can't be blank"),
			},
			// Providing invalid extractions attributes for the catch_all rule
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, "", invalidExtractionRegexTemplateNilConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, "", invalidExtractionRegexTemplateValConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceInvalidExtractionsConfig(
					escalationPolicy, service, "", invalidExtractionRegexNilSourceConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: source can't be blank"),
			},
			// Adding/updating/deleting all actions
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(resourceName, "set.0.rule.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "set.1.rule.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "set.1.rule.1.id"),
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
							resource.TestCheckResourceAttrSet(resourceName, "set.0.rule.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "set.1.rule.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "set.1.rule.1.id"),
							resource.TestCheckResourceAttr(
								resourceName, "set.0.rule.0.actions.0.escalation_policy", "POLICY3",
							),
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
							resource.TestCheckResourceAttrSet(resourceName, "set.0.rule.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "set.1.rule.0.id"),
							resource.TestCheckResourceAttrSet(resourceName, "set.1.rule.1.id"),
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
						resource.TestCheckResourceAttrSet(resourceName, "set.0.rule.0.id"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceResourceDeleteConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathNotExists(resourceName),
				),
			},
			// Adding/Updating/Removing `enable_event_orchestration_for_service` attribute
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceDefaultConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckNoResourceAttr(resourceName, "enable_event_orchestration_for_service"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceEnableEOForServiceEnableUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(resourceName, "enable_event_orchestration_for_service", "true"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceEnableEOForServiceDisableUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(resourceName, "enable_event_orchestration_for_service", "false"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceEnableEOForServiceEnableUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(resourceName, "enable_event_orchestration_for_service", "true"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceResourceDeleteConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathNotExists(resourceName),
				),
			},
			// Disabling Service Orchestration at creation by setting
			// `enable_event_orchestration_for_service`  attribute to false
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceEnableEOForServiceDisableUpdateConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(resourceName, "enable_event_orchestration_for_service", "false"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceDefaultConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(resourceName, "enable_event_orchestration_for_service", "false"),
					)...,
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

		if _, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), srv.Primary.ID, "service"); err == nil {
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
		found, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "service")
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

func testAccCheckPagerDutyEventOrchestrationPathServiceServiceID(rn, sn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		p, _ := s.RootModule().Resources[rn]
		srv, ok := s.RootModule().Resources[sn]

		if !ok {
			return fmt.Errorf("Service not found: %s", sn)
		}

		var pId = p.Primary.Attributes["service"]
		var sId = srv.Primary.Attributes["id"]
		if pId != sId {
			return fmt.Errorf("Event Orchestration Service path service ID (%v) not matching provided service ID: %v", pId, sId)
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
			service = pagerduty_service.bar.id

			set {
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
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
					label = "rule 1"
					actions {
							automation_action {
								name = "test"
								url = "https://test.com"
								auto_send = true

								header {
									key = "foo"
									value = "bar"
								}
								header {
									key = "baz"
									value = "buz"
								}

								parameter {
									key = "source"
									value = "orch"
								}
								parameter {
									key = "region"
									value = "us"
								}
							}
					}
				}
			}

			catch_all {
				actions {
					automation_action {
						name = "catch-all test"
						url = "https://catch-all-test.com"
						auto_send = true

						header {
							key = "foo1"
							value = "bar1"
						}
						header {
							key = "baz1"
							value = "buz1"
						}

						parameter {
							key = "source1"
							value = "orch1"
						}
						parameter {
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
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
					label = "rule 1"
					actions {
							automation_action {
								name = "test1"
								url = "https://test1.com"

								header {
									key = "foo1"
									value = "bar1"
								}
								parameter {
									key = "source_region"
									value = "eu"
								}
							}
					}
				}
			}

			catch_all {
				actions {
					automation_action {
						name = "catch-all test upd"
						url = "https://catch-all-test-upd.com"

						header {
							key = "baz2"
							value = "buz2"
						}

						parameter {
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
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
					label = "rule 1"
					actions {
							automation_action {
								name = "test"
								url = "https://test.com"
							}
					}
				}
			}

			catch_all {
				actions {
					automation_action {
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
				service = pagerduty_service.bar.id

				set {
					id = "start"
					rule {
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

func testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
					label = "rule 1"
					condition {
						expression = "event.summary matches part 'timeout'"
					}
					condition {
						expression = "event.custom_details.timeout_err exists"
					}
					actions {
						route_to = "set-1"
						priority = "P0IN2KQ"
						escalation_policy = pagerduty_escalation_policy.foo.id
						annotate = "Routed through an event orchestration"
						pagerduty_automation_action {
							action_id = "01CSB5SMOKCKVRI5GN0LJG7SMB"
						}
						severity = "critical"
						event_action = "trigger"
						variable {
							name = "hostname"
							path = "event.source"
							type = "regex"
							value = "Source host: (.*)"
						}
						variable {
							name = "cpu_val"
							path = "event.custom_details.cpu"
							type = "regex"
							value = "(.*)"
						}
						extraction {
							target = "event.summary"
							template = "High CPU usage on {{variables.hostname}}"
						}
						extraction {
							regex = ".*"
							source = "event.group"
							target = "event.custom_details.message"
						}
						incident_custom_field_update {
							id = "PIJ90N7"
							value = "foo"
						}
					}
				}
			}
			set {
				id = "set-1"
				rule {
					label = "set-1 rule 1"
					actions {
						suspend = 300
					}
				}
				rule {
					label = "set-1 rule 2"
					condition {
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
					escalation_policy = pagerduty_escalation_policy.foo.id
					annotate = "Routed through an event orchestration - catch-all rule"
					pagerduty_automation_action {
						action_id = "01CSB5SMOKCKVRI5GN0LJG7SMC"
					}
					severity = "warning"
					event_action = "trigger"
					variable {
						name = "user_id"
						path = "event.custom_details.user_id"
						type = "regex"
						value = "Source host: (.*)"
					}
					variable {
						name = "updated_at"
						path = "event.custom_details.updated_at"
						type = "regex"
						value = "(.*)"
					}
					extraction {
						target = "event.custom_details.message"
						template = "Last modified by {{variables.user_id}} on {{variables.updated_at}}"
					}
					extraction {
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
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
					label = "rule 1 updated"
					condition {
						expression = "event.custom_details.timeout_err matches part 'timeout'"
					}
					actions {
						route_to = "set-2"
						priority = "P0IN2KR"
						escalation_policy = "POLICY3"
						annotate = "Routed through a service orchestration!"
						pagerduty_automation_action {
							action_id = "01CSB5SMOKCKVRI5GN0LJG7SMBUPDATED"
						}
						severity = "warning"
						event_action = "resolve"
						variable {
							name = "cpu_val_upd"
							path = "event.custom_details.cpu_upd"
							type = "regex"
							value = "CPU:(.*)"
						}
						extraction {
							regex = ".*"
							source = "event.custom_details.region_upd"
							target = "event.source"
						}
						extraction {
							target = "event.custom_details.message_upd"
							template = "[UPD] High CPU usage on {{variables.hostname}}: {{variables.cpu_val}}"
						}
						incident_custom_field_update {
							id = "PIJ90N7"
							value = "bar"
						}
					}
				}
			}
			set {
				id = "set-2"
				rule {
					label = "set-2 rule 1"
					actions {
						suspend = 15
					}
				}
				rule {
					label = "set-2 rule 2"
					condition {
						expression = "event.source matches part 'test-'"
					}
					actions {
						annotate = "Matched set-2 rule 2"
						variable {
							name = "host_name"
							path = "event.custom_details.memory"
							type = "regex"
							value = "High memory usage on (.*) server"
						}
						extraction {
							target = "event.summary"
							template = "High memory usage on {{variables.hostname}} server: {{event.custom_details.max_memory}}"
						}
						extraction {
							regex = ".*"
							source = "event.custom_details.region"
							target = "event.group"
						}
						extraction {
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
					priority = "P0IN2KX"
					escalation_policy = "POLICY4"
					annotate = "[UPD] Routed through an event orchestration - catch-all rule"
					pagerduty_automation_action {
						action_id = "01CSB5SMOKCKVRI5GN0LJG7SMD"
					}
					severity = "info"
					event_action = "resolve"
					variable {
						name = "updated_at_upd"
						path = "event.custom_details.updated_at"
						type = "regex"
						value = "UPD (.*)"
					}
					extraction {
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
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
					label = "rule 1 updated"
					actions {
						route_to = "set-2"
					}
				}
			}
			set {
				id = "set-2"
				rule {
					label = "set-2 rule 1"
					actions { }
				}
				rule {
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
			service = pagerduty_service.bar.id

			set {
				id = "start"
				rule {
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

func testAccCheckPagerDutyEventOrchestrationPathServiceEnableEOForServiceEnableUpdateConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			service = pagerduty_service.bar.id
      enable_event_orchestration_for_service = true

			set {
				id = "start"
			}

			catch_all {
				actions { }
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathServiceEnableEOForServiceDisableUpdateConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			service = pagerduty_service.bar.id
      enable_event_orchestration_for_service = false

			set {
				id = "start"
			}

			catch_all {
				actions { }
			}
		}
	`)
}
