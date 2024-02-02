package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyEventOrchestrationPathGlobal_import(t *testing.T) {
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orch := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationGlobalPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathGlobalAllActionsConfig(team, escalationPolicy, service, orch),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_global.my_global_orch",
				ImportStateIdFunc: testAccCheckPagerDutyEventOrchestrationPathGlobalID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathGlobalID(s *terraform.State) (string, error) {
	return s.RootModule().Resources["pagerduty_event_orchestration.orch"].Primary.ID, nil
}
