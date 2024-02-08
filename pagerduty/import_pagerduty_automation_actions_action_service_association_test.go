package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyAutomationActionsActionServiceAssociation_import(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	actionName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyAutomationActionsActionServiceAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAutomationActionsActionServiceAssociationConfig(username, email, escalationPolicy, serviceName, actionName),
			},
			{
				ResourceName:      "pagerduty_automation_actions_action_service_association.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
