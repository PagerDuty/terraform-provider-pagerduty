package pagerduty

import (
	"fmt"
	// "log"
	// "strings"
	"testing"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_service", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_service",
		Dependencies: []string{
			"pagerduty_service",
			// TODO: EP, Schedule, user
		},
	})
}

func TestAccPagerDutyEventOrchestrationPathService_Basic(t *testing.T) {
	//name := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	// description := fmt.Sprintf("tf-description-%s", acctest.RandString(5))
	// nameUpdated := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	// descriptionUpdated := fmt.Sprintf("tf-description-%s", acctest.RandString(5))
	// team1 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	// team2 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	service := "PZ73WUB"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceRequiredFieldsConfig(service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceExists("pagerduty_event_orchestration_service.serviceA"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration_service.serviceA", "sets.#", "1",
					),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.foo", "sets.0.rules.#", "0",
					),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceAllFieldsConfig(service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationServiceExists("pagerduty_event_orchestration_service.serviceA"),
				),
			},
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

func testAccCheckPagerDutyEventOrchestrationServiceRequiredFieldsConfig(sId string) string {
	return fmt.Sprintf(`

resource "pagerduty_event_orchestration_service" "serviceA" {
	parent {
		id = "%s"
	}

	sets {
		id = "start"
	}
}
`, sId)
}

func testAccCheckPagerDutyEventOrchestrationServiceAllFieldsConfig(sId string) string {
	return fmt.Sprintf(`

resource "pagerduty_event_orchestration_service" "serviceA" {
	parent {
		id = "%s"
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
`, sId)
}
