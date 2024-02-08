package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyAutomationActionsAction_import(t *testing.T) {
	actionName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionTypeProcessAutomationConfig(actionName),
			},
			{
				ResourceName:      "pagerduty_automation_actions_action.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
