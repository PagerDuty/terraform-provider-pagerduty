package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyEventOrchestrationPathRouter_import(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterConfigWithMultipleRules(team, escalationPolicy, service, orchestration),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_router.router",
				ImportStateIdFunc: testAccCheckPagerDutyEventOrchestrationPathRouterID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"set.0.rule.0.id",
					"set.0.rule.1.id",
				},
			},
		},
	})
}

func TestAccPagerDutyEventOrchestrationPathRouterDynamicRouteTo_import(t *testing.T) {
	team := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	orchestration := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	dynamicRouteToByNameInput := &pagerduty.EventOrchestrationPathDynamicRouteTo{
		LookupBy: "service_name",
		Regex:    ".*",
		Source:   "event.custom_details.pd_service_name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationRouterDynamicRouteToConfig(team, escalationPolicy, service, orchestration, dynamicRouteToByNameInput),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_router.router",
				ImportStateIdFunc: testAccCheckPagerDutyEventOrchestrationPathRouterID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"set.0.rule.0.id",
					"set.0.rule.1.id",
				},
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationPathRouterID(s *terraform.State) (string, error) {
	return s.RootModule().Resources["pagerduty_event_orchestration.orch"].Primary.ID, nil
}
