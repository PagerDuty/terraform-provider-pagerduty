package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyUserHandoffNotificationRule_import(t *testing.T) {
	userName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	contactMethodType := "phone_contact_method"
	handoffType := "both"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyUserHandoffNotificationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserHandoffNotificationRuleConfig(userName, contactMethodType, handoffType),
			},
			{
				ResourceName:      "pagerduty_user_handoff_notification_rule.foo",
				ImportStateIdFunc: testAccCheckPagerDutyUserHandoffNotificationRuleID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyUserHandoffNotificationRuleID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v.%v", s.RootModule().Resources["pagerduty_user.foo"].Primary.ID, s.RootModule().Resources["pagerduty_user_handoff_notification_rule.foo"].Primary.ID), nil
}
