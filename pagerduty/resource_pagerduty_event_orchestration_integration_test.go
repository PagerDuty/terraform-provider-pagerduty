package pagerduty

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration_integration", &resource.Sweeper{
		Name: "pagerduty_event_orchestration_integration",
		// deleting all test orchestrations will result in integrations being deleted as well
		F: testSweepEventOrchestration,
	})
}

func TestAccPagerDutyEventOrchestrationIntegration_Basic(t *testing.T) {
	onp := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	rn := "pagerduty_event_orchestration_integration.int_1"
	lbl1 := fmt.Sprintf("tf-integration-%s", acctest.RandString(5))
	orn1 := "orch_1"
	lbl2 := fmt.Sprintf("tf-integration-updated-%s", acctest.RandString(5))
	orn2 := "orch_2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationIntegrationConfig(onp, lbl1, orn1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationIntegrationAttr(rn, orn1),
					resource.TestCheckResourceAttr(rn, "label", lbl1),
				),
			},
			// update label and event_orchestration:
			{
				Config: testAccCheckPagerDutyEventOrchestrationIntegrationConfig(onp, lbl2, orn2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationIntegrationAttr(rn, orn2),
					resource.TestCheckResourceAttr(rn, "label", lbl2),
				),
			},
			// update event_orchestration:
			{
				Config: testAccCheckPagerDutyEventOrchestrationIntegrationConfig(onp, lbl2, orn1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationIntegrationAttr(rn, orn1),
					resource.TestCheckResourceAttr(rn, "label", lbl2),
				),
			},
			// update label:
			{
				Config: testAccCheckPagerDutyEventOrchestrationIntegrationConfig(onp, lbl1, orn1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationIntegrationAttr(rn, orn1),
					resource.TestCheckResourceAttr(rn, "label", lbl1),
				),
			},
			// delete integration:
			{
				Config: testAccCheckPagerDutyEventOrchestrationIntegrationDeletedConfig(onp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationIntegrationExistsNot(rn),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationIntegrationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration_integration" {
			continue
		}
		if _, _, err := client.EventOrchestrationIntegrations.GetContext(context.Background(), r.Primary.Attributes["event_orchestration"], r.Primary.ID); err == nil {
			return fmt.Errorf("Event Orchestration Integration still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationIntegrationExistsNot(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[rn]
		if ok {
			return fmt.Errorf("Event Orchestration Integration is not deleted from the state: %s", rn)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationIntegrationAttr(rn, orn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ir, ok := s.RootModule().Resources[rn]
		eor, _ := s.RootModule().Resources[fmt.Sprintf("pagerduty_event_orchestration.%s", orn)]

		if !ok {
			return fmt.Errorf("Event Orchestration Integration resource not found in the state: %s", rn)
		}

		oid := ir.Primary.Attributes["event_orchestration"]
		id := ir.Primary.ID

		client, _ := testAccProvider.Meta().(*Config).Client()
		i, _, err := client.EventOrchestrationIntegrations.GetContext(context.Background(), oid, id)
		eo, _, _ := client.EventOrchestrations.Get(eor.Primary.ID)

		if err != nil {
			return err
		}

		if i.ID != id {
			return fmt.Errorf("Event Orchestration Integration ID does not match the resource ID: %v - %v", i.ID, id)
		}

		if eo.ID != oid {
			return fmt.Errorf("Event Orchestration Integration's parent ID does not match the resource event_orchestration attr: %v - %v", eo.ID, oid)
		}

		lbl := ir.Primary.Attributes["label"]
		if i.Label != lbl {
			return fmt.Errorf("Event Orchestration ID does not match the resource label attr: %v - %v", i.Label, lbl)
		}

		rkey := ir.Primary.Attributes["parameters.0.routing_key"]
		if i.Parameters.RoutingKey != rkey {
			return fmt.Errorf("Event Orchestration Integration routing_key does not match the resource routing_key attr: %v - %v", i.Parameters.RoutingKey, rkey)
		}

		t := ir.Primary.Attributes["parameters.0.type"]
		if i.Parameters.Type != t {
			return fmt.Errorf("Event Orchestration Integration routing_key type does not match the resource routing_key type attr: %v - %v", i.Parameters.Type, t)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationIntegrationConfig(onp, lbl, orn string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch_1" {
			name = "%s-1"
		}

		resource "pagerduty_event_orchestration" "orch_2" {
			name = "%s-2"
		}

		resource "pagerduty_event_orchestration_integration" "int_1" {
			label = "%s"
			event_orchestration = pagerduty_event_orchestration.%s.id
		}
	`, onp, onp, lbl, orn)
}

func testAccCheckPagerDutyEventOrchestrationIntegrationDeletedConfig(onp string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch_1" {
			name = "%s-1"
		}

		resource "pagerduty_event_orchestration" "orch_2" {
			name = "%s-2"
		}
	`, onp, onp)
}
