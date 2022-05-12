package pagerduty

import (
	"fmt"
	"log"
	// "strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	// "github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_service", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_service",
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
					testAccCheckPagerDutyEventOrchestrationServiceExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "sets.#", "1",
					),
					resource.TestCheckResourceAttr(
						resourceName, "sets.0.rules.#", "0",
					),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAllFieldsConfig(escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "sets.#", "1",
					),
					resource.TestCheckResourceAttr(
						resourceName, "sets.0.rules.#", "1",
					),
					testAccCheckPagerDutyEventOrchestrationServiceRuleActions(resourceName, "sets.0.rules.0.actions.0", map[string]interface{}{
						"pagerduty_automation_actions": map[string]interface{}{"action_id": "SOME_ACTION_ID"},
					}),
				),
			},
			// update all fields
			// reset rule action items -> should be default
			// reset rule actions -> should be []
			// reset rule conditions -> should be []
			// reset rules
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationServiceExists(rn string) resource.TestCheckFunc {
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

func testAccCheckPagerDutyEventOrchestrationServiceRuleActions(rn, path string, a map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}

		attr := r.Primary.Attributes[path]

		log.Printf(">>> attr: %v", attr)
		log.Printf(">>> a: %v", a)

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
				id = "pagerduty_service.bar.id"
			}
		
			sets {
				id = "start"
			}
		}
	`)
}

func testAccCheckPagerDutyEventOrchestrationServiceAllFieldsConfig(ep, s string) string {
	return fmt.Sprintf("%s%s", createBaseServicePathConfig(ep, s),
		`resource "pagerduty_event_orchestration_service" "serviceA" {
			parent {
				id = "pagerduty_service.bar.id"
			}
		
			sets {
				id = "start"
				rules {
					label = "rule 1"
					actions {
							pagerduty_automation_actions {
								action_id = "SOME_ACTION_ID"
							}
							automation_actions {
								name = "Reboot me"
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
									value = "orch_rule"
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
