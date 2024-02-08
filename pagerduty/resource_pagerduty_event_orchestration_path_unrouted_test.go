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
	resource.AddTestSweepers("pagerduty_event_orchestration_unrouted", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_unrouted",
		F:    testSweepEventOrchestration,
	})
}

func TestAccPagerDutyEventOrchestrationPathUnrouted_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationPathUnroutedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigNoRules(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.#", "0"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigWithConditions(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.condition.0.expression", "event.summary matches part 'rds'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigWithMultipleRules(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.condition.0.expression", "event.summary matches part 'rds'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.condition.1.expression", "event.severity matches part 'warning'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.1.condition.0.expression", "event.severity matches part 'info'"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedWithAllConfig(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.#", "2"),
					//Set #1
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.id", "start"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.condition.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.actions.0.route_to", "child-1"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.condition.0.expression", "event.severity matches part 'info'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.condition.1.expression", "event.severity matches part 'warning'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.actions.0.severity", "info"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.actions.0.event_action", "trigger"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.actions.0.variable.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"set.0.rule.0.actions.0.variable.*",
						map[string]string{
							"name":  "server_name_cpu",
							"path":  "event.summary",
							"type":  "regex",
							"value": "High CPU on (.*) server",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"set.0.rule.0.actions.0.variable.*",
						map[string]string{
							"name":  "server_name_memory",
							"path":  "event.custom_details",
							"type":  "regex",
							"value": "High memory usage on (.*) server",
						},
					),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.0.actions.0.extraction.#", "2"),

					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"set.0.rule.0.actions.0.extraction.*",
						map[string]string{
							"target":   "event.summary",
							"template": "High memory usage on {{variables.hostname}} server",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"set.0.rule.0.actions.0.extraction.*",
						map[string]string{
							"target":   "event.custom_details",
							"template": "High memory usage on {{variables.hostname}} server",
						},
					),
					//Set #2
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.1.id", "child-1"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.1.rule.0.condition.0.expression", "event.severity matches part 'warning'"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.1.rule.0.actions.0.event_action", "resolve"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.1.rule.1.condition.0.expression", "event.severity matches part 'critical'"),
					// Catch All
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "catch_all.0.actions.0.severity", "critical"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "catch_all.0.actions.0.event_action", "trigger"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "catch_all.0.actions.0.variable.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"catch_all.0.actions.0.variable.*",
						map[string]string{
							"name":  "server_name_cpu",
							"path":  "event.summary",
							"type":  "regex",
							"value": "High CPU on (.*) server",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"catch_all.0.actions.0.variable.*",
						map[string]string{
							"name":  "server_name_memory",
							"path":  "event.custom_details",
							"type":  "regex",
							"value": "High memory usage on (.*) server",
						},
					),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "catch_all.0.actions.0.extraction.#", "2"),

					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"catch_all.0.actions.0.extraction.*",
						map[string]string{
							"target":   "event.summary",
							"template": "High memory usage on {{variables.hostname}} server",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						"pagerduty_event_orchestration_unrouted.unrouted",
						"catch_all.0.actions.0.extraction.*",
						map[string]string{
							"target":   "event.custom_details",
							"template": "High memory usage on {{variables.hostname}} server",
						},
					),
				),
			},
			// Providing invalid extractions attributes for set rules
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(
					team, escalationPolicy, service, orchestration, invalidExtractionRegexTemplateValConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(
					team, escalationPolicy, service, orchestration, invalidExtractionRegexTemplateValConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(
					team, escalationPolicy, service, orchestration, invalidExtractionRegexNilSourceConfig(), "",
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in set.0.rule.0.actions.0.extraction.0: source can't be blank"),
			},
			// Providing invalid extractions attributes for the catch_all rule
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(
					team, escalationPolicy, service, orchestration, "", invalidExtractionRegexTemplateNilConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: regex and template cannot both be null"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(
					team, escalationPolicy, service, orchestration, "", invalidExtractionRegexTemplateValConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: regex and template cannot both have values"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(
					team, escalationPolicy, service, orchestration, "", invalidExtractionRegexNilSourceConfig(),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid configuration in catch_all.0.actions.0.extraction.0: source can't be blank"),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigNoRules(team, escalationPolicy, service, orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationPathUnroutedExists("pagerduty_event_orchestration_unrouted.unrouted"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_unrouted.unrouted", "set.0.rule.#", "0"),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigDelete(team, escalationPolicy, service, orchestration),
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

		if _, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "unrouted"); err == nil {
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
			return fmt.Errorf("No Event Orchestration Unrouted Path is set")
		}

		orch := s.RootModule().Resources["pagerduty_event_orchestration.orch"]
		client, _ := testAccProvider.Meta().(*Config).Client()
		_, _, err := client.EventOrchestrationPaths.GetContext(context.Background(), orch.Primary.ID, "unrouted")

		if err != nil {
			return fmt.Errorf("Event Orchestration Unrouted Path not found: %v for orchestration %v", "unrouted", orch.Primary.ID)
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
	return createUnroutedBaseConfig(t, ep, s, o)
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
			team = pagerduty_team.foo.id
		}
	`, t, ep, s, o)
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedConfigNoRules(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createUnroutedBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
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
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
				id = "start"
				rule {
					disabled = false
					label = "rule1 label"
					actions { }
					condition {
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
			event_orchestration = pagerduty_event_orchestration.orch.id

			set {
				id = "start"
				rule {
					disabled = false
					label = "rule1 label"
					actions { }
					condition {
						expression = "event.summary matches part 'rds'"
					}
					condition {
						expression = "event.severity matches part 'warning'"
					}
				}

				rule {
					disabled = false
					label = "rule2 label"
					actions { }
					condition {
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

func testAccCheckPagerDutyEventOrchestrationPathUnroutedWithAllConfig(t, ep, s, o string) string {
	return fmt.Sprintf("%s%s", createUnroutedBaseConfig(t, ep, s, o),
		`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			
			set {
				id = "start"
				rule {
					disabled = false
					label = "rule1 label"
					condition {
						expression = "event.severity matches part 'info'"
					}
					condition {
						expression = "event.severity matches part 'warning'"
					}
					actions {
						route_to = "child-1"
						severity = "info"
						event_action = "trigger"
						variable {
							name = "server_name_cpu"
							path = "event.summary"
							type = "regex"
							value = "High CPU on (.*) server"
						}
						variable {
							name = "server_name_memory"
							path = "event.custom_details"
							type = "regex"
							value = "High memory usage on (.*) server"
						}
						extraction {
							target = "event.summary"
							template = "High memory usage on {{variables.hostname}} server"
						}
						extraction {
							target = "event.custom_details"
							template = "High memory usage on {{variables.hostname}} server"
						}
					}
				}
			}
			set {
				id = "child-1"
				rule {
					disabled = false
					label = "rule2 label1"
					condition {
						expression = "event.severity matches part 'warning'"
					}
					actions {
						severity = "warning"
						event_action = "resolve"
						variable {
							name = "server_name_cpu"
							path = "event.summary"
							type = "regex"
							value = "High CPU on (.*) server"
						}
						extraction {
							target = "event.summary"
							template = "High CPU on {{event.custom_details.hostname}} server"
						}
					}
				}
				rule {
					disabled = false
					label = "rule2 label2"
					condition {
						expression = "event.severity matches part 'critical'"
					}
					actions {
						severity = "warning"
						event_action = "trigger"
						variable {
							name = "server_name_cpu"
							path = "event.summary"
							type = "regex"
							value = "High CPU on (.*) server"
						}
						extraction {
							target = "event.summary"
							template = "High CPU on {{event.custom_details.hostname}} server"
						}
					}
				}
			}
			catch_all {
				actions {
					severity = "critical"
					event_action = "trigger"
					variable {
						name = "server_name_cpu"
						path = "event.summary"
						type = "regex"
						value = "High CPU on (.*) server"
					}
					variable {
						name = "server_name_memory"
						path = "event.custom_details"
						type = "regex"
						value = "High memory usage on (.*) server"
					}
					extraction {
						target = "event.summary"
						template = "High memory usage on {{variables.hostname}} server"
					}
					extraction {
						target = "event.custom_details"
						template = "High memory usage on {{variables.hostname}} server"
					}
				}
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedInvalidExtractionsConfig(t, ep, s, o, re, cae string) string {
	return fmt.Sprintf(
		"%s%s",
		createUnroutedBaseConfig(t, ep, s, o),
		fmt.Sprintf(`resource "pagerduty_event_orchestration_unrouted" "unrouted" {
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
