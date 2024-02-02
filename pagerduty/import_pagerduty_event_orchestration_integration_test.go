package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyEventOrchestrationIntegration_import(t *testing.T) {
	onp := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	lbl := fmt.Sprintf("tf-integration-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationIntegrationConfig(onp, lbl, "orch_1"),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_integration.int_1",
				ImportStateIdFunc: testAccPagerDutyEventOrchestrationIntegrationImportID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPagerDutyEventOrchestrationIntegrationImportID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v:%v", s.RootModule().Resources["pagerduty_event_orchestration.orch_1"].Primary.ID, s.RootModule().Resources["pagerduty_event_orchestration_integration.int_1"].Primary.ID), nil
}
