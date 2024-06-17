package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyEventOrchestrationPathUnrouted_import(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationPathUnroutedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationPathUnroutedWithAllConfig(team, escalationPolicy, service, orchestration),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_unrouted.unrouted",
				ImportStateIdFunc: testAccCheckPagerDutyEventOrchestrationPathUnroutedID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"set.0.rule.0.id",
					"set.1.rule.0.id",
					"set.1.rule.1.id",
				},
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathUnroutedID(s *terraform.State) (string, error) {
	return s.RootModule().Resources["pagerduty_event_orchestration.orch"].Primary.ID, nil
}
