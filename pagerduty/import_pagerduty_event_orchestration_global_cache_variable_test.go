package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyEventOrchestrationGlobalCacheVariable_import(t *testing.T) {
	orch := fmt.Sprintf("tf_orchestration_%s", acctest.RandString(5))
	name := fmt.Sprintf("tf_global_cache_variable_%s", acctest.RandString(5))
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
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationGlobalCacheVariableConfig(orch, name, "orch_1", disabled, config, cond),
			},
			{
				ResourceName:      "pagerduty_event_orchestration_global_cache_variable.cv_1",
				ImportStateIdFunc: testAccPagerDutyEventOrchestrationGlobalCacheVariableImportID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPagerDutyEventOrchestrationGlobalCacheVariableImportID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v:%v", s.RootModule().Resources["pagerduty_event_orchestration.orch_1"].Primary.ID, s.RootModule().Resources["pagerduty_event_orchestration_global_cache_variable.cv_1"].Primary.ID), nil
}
