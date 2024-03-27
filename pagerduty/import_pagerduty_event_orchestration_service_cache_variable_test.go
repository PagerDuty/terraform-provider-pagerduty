package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyEventOrchestrationServiceCacheVariable_import(t *testing.T) {
	svc := fmt.Sprintf("tf_service_%s", acctest.RandString(5))
	name := fmt.Sprintf("tf_service_cache_variable_%s", acctest.RandString(5))
	disabled := "false"
	config := `
		configuration {
  		type = "trigger_event_count"
  		ttl_seconds = 60
    }
  `
	cond := ""

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationServiceCacheVariableConfig(svc, name, "svc_1", disabled, config, cond),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_service_cache_variable.cv_1",
				ImportStateIdFunc: testAccPagerDutyEventOrchestrationServiceCacheVariableImportID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPagerDutyEventOrchestrationServiceCacheVariableImportID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v:%v", s.RootModule().Resources["pagerduty_service.svc_1"].Primary.ID, s.RootModule().Resources["pagerduty_event_orchestration_service_cache_variable.cv_1"].Primary.ID), nil
}
