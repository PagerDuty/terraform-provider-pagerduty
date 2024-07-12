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
	resource.AddTestSweepers("pagerduty_event_orchestration_global", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_global",
		F:    testSweepEventOrchestration,
	})
}

func TestAccPagerDutyEventOrchestrationPathGlobal_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orch := fmt.Sprintf("tf-%s", acctest.RandString(5))

	res := "pagerduty_event_orchestration_global.my_global_orch"
	orchRes := "pagerduty_event_orchestration.orch"

	baseChecks := []resource.TestCheckFunc{
		testAccCheckPagerDutyEventOrchestrationGlobalExists(res),
		testAccCheckPagerDutyEventOrchestrationPathGlobalOrchID(res, orchRes),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationGlobalPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalDefaultConfig(team, escalationPolicy, service, orch),
				Check:  resource.ComposeTestCheckFunc(baseChecks...),
			},
			// Adding/updating/deleting automation_action properties
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAutomationActionsConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttrSet(res, "set.0.rule.0.id"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAutomationActionsParamsUpdateConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttr(
							res, "set.0.rule.0.actions.0.automation_action.0.auto_send", "false",
						),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAutomationActionsParamsDeleteConfig(team, escalationPolicy, service, orch),
				Check:  resource.ComposeTestCheckFunc(baseChecks...),
			},
			// Providing invalid extractions attributes for set rules
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(
					team, escalationPolicy, service, orch, invalidExtractionRegexTemplateNilConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(
					team, escalationPolicy, service, orch, invalidExtractionRegexTemplateValConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(
					team, escalationPolicy, service, orch, invalidExtractionRegexNilSourceConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: source can't be blank"),
			},
			// Providing invalid extractions attributes for the catch_all rule
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(
					team, escalationPolicy, service, orch, "", invalidExtractionRegexTemplateNilConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(
					team, escalationPolicy, service, orch, "", invalidExtractionRegexTemplateValConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(
					team, escalationPolicy, service, orch, "", invalidExtractionRegexNilSourceConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: source can't be blank"),
			},
			// Adding/updating/deleting all actions
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(res, "set.0.rule.0.id"),
							resource.TestCheckResourceAttrSet(res, "set.0.rule.1.id"),
							resource.TestCheckResourceAttrSet(res, "set.1.rule.0.id"),
							resource.TestCheckResourceAttrSet(res, "set.1.rule.1.id"),
						}...,
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsUpdateConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(res, "set.0.rule.0.id"),
							resource.TestCheckResourceAttrSet(res, "set.0.rule.1.id"),
							resource.TestCheckResourceAttrSet(res, "set.1.rule.0.id"),
							resource.TestCheckResourceAttrSet(res, "set.1.rule.1.id"),
							resource.TestCheckResourceAttr(
								res, "set.0.rule.0.actions.0.escalation_policy", "POLICY3",
							),
						}...,
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsDeleteConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						[]resource.TestCheckFunc{
							resource.TestCheckResourceAttrSet(res, "set.0.rule.0.id"),
							resource.TestCheckResourceAttrSet(res, "set.0.rule.1.id"),
							resource.TestCheckResourceAttrSet(res, "set.1.rule.0.id"),
							resource.TestCheckResourceAttrSet(res, "set.1.rule.1.id"),
						}...,
					)...,
				),
			},
			// Deleting sets and the service path resource
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalOneSetNoActionsConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					append(
						baseChecks,
						resource.TestCheckResourceAttrSet(res, "set.0.rule.0.id"),
					)...,
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalResourceDeleteConfig(team, escalationPolicy, service, orch),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServicePathNotExists(res),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationGlobalPathDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_global" {
			continue
		}

		orch, _ := s.RootModule().Resources["pagerduty_event_orchestration.orch"]

		if _, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "global"); err == nil {
			return fmt.Errorf("Event Orchestration Path still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationGlobalExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Global Orchestration ID is not configured")
		}

		orch, _ := s.RootModule().Resources["pagerduty_event_orchestration.orch"]
		client, _ := testAccProvider.Meta().(*Config).Client()
		_, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "global")

		if err != nil {
			return fmt.Errorf("Global Orchestration not found for orchestration %v", orch.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationPathGlobalOrchID(rn, on string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		p, _ := s.RootModule().Resources[rn]
		orch, ok := s.RootModule().Resources[on]

		if !ok {
			return fmt.Errorf("Event Orchestration not found: %s", on)
		}

		var pId = p.Primary.Attributes["event_orchestration"]
		var orchId = orch.Primary.Attributes["id"]
		if pId != orchId {
			return fmt.Errorf("Event Orchestration Global path event_orchestration ID (%v) not matching provided orchestration ID: %v", pId, orchId)
		}

		return nil
	}
}

func createBaseGlobalOrchConfig(t, ep, s, o string) string {
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

func testAccCheckPagerDutyEventOrchestrationGlobalDefaultConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			catch_all {
				actions {}
			}
			set {
				id = "start"
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathGlobalAutomationActionsConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

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

func testAccCheckPagerDutyEventOrchestrationPathGlobalAutomationActionsParamsUpdateConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

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

func testAccCheckPagerDutyEventOrchestrationPathGlobalAutomationActionsParamsDeleteConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

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

func testAccCheckPagerDutyEventOrchestrationPathGlobalInvalidExtractionsConfig(t, ep, s, o, re, cae string) string {
	return fmt.Sprintf(
		"%s%s",
		createBaseGlobalOrchConfig(t, ep, s, o),
		fmt.Sprintf(`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

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

func testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
				id = "start"
				rule {
					label = "start rule 1"
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
				rule {
					label = "start rule 2"
					actions {
						drop_event = true
					}
					condition {
						expression = "event.summary matches part '[test]'"
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
					drop_event = true
					priority = "P0IN2KW"
					escalation_policy = pagerduty_escalation_policy.foo.id
					annotate = "Routed through an event orchestration - catch-all rule"
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

func testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsUpdateConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
				id = "start"
				rule {
					label = "start rule 1 updated"
					condition {
						expression = "event.custom_details.timeout_err matches part 'timeout'"
					}
					actions {
						route_to = "set-2"
						priority = "P0IN2KR"
						escalation_policy = "POLICY3"
						annotate = "Routed through a service orchestration!"
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
				rule {
					label = "start rule 2 updated"
					actions {
						drop_event = false
					}
					condition {
						expression = "event.summary matches '[test - create incident]'"
					}
				}
			}
			set {
				id = "set-2"
				rule {
					label = "set-2 rule 1"
					actions {
						suspend = 15
						escalation_policy = pagerduty_escalation_policy.foo.id
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
					drop_event = false
					priority = "P0IN2KX"
					escalation_policy = "POLICY4"
					annotate = "[UPD] Routed through an event orchestration - catch-all rule"
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

func testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsDeleteConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
				id = "start"
				rule {
					label = "start rule 1 updated"
					actions {
						route_to = "set-2"
					}
				}
				rule {
					label = "start rule 2 updated"
					actions { }
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

func testAccCheckPagerDutyEventOrchestrationPathGlobalOneSetNoActionsConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createBaseGlobalOrchConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_global" "my_global_orch" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
				id = "start"
				rule {
					label = "start rule 1 updated"
					actions {}
				}
			}

			catch_all {
				actions { }
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathGlobalResourceDeleteConfig(t, ep, s, o string) string {
	return createBaseGlobalOrchConfig(t, ep, s, o)
}
