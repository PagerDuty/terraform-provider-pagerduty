package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyEventOrchestrationPathService_import(t *testing.T) {
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationServicePathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathServiceAllActionsConfig(escalationPolicy, service),
			},
			{
				ResourceName:            "pagerduty_event_orchestration_service.serviceA",
				ImportStateIdFunc:       testAccCheckPagerDutyEventOrchestrationPathServiceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_event_orchestration_for_service"},
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathServiceID(s *terraform.State) (string, error) {
	return s.RootModule().Resources["pagerduty_service.bar"].Primary.ID, nil
}
