package pagerduty

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyUserHandoffNotificationRule_Basic(t *testing.T) {
	userName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	invalidContactMethodType := "phone"
	contactMethodType := "phone_contact_method"
	invalidHandoffType := "on-call"
	handoffType := "both"
	updatedHandoffType := "oncall"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyUserHandoffNotificationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyUserHandoffNotificationRuleConfig(userName, invalidContactMethodType, handoffType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"phone_contact_method" "phone_contact_method_reference"`),
			},
			{
				Config:      testAccCheckPagerDutyUserHandoffNotificationRuleConfig(userName, contactMethodType, invalidHandoffType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"both" "oncall" "offcall"`),
			},
			{
				Config: testAccCheckPagerDutyUserHandoffNotificationRuleConfig(userName, contactMethodType, handoffType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserHandoffNotificationRuleExists("pagerduty_user_handoff_notification_rule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user_handoff_notification_rule.foo", "handoff_type", handoffType),
					resource.TestCheckResourceAttr(
						"pagerduty_user_handoff_notification_rule.foo", "contact_method.0.type", contactMethodType),
				),
			},
			{
				Config: testAccCheckPagerDutyUserHandoffNotificationRuleConfig(userName, contactMethodType, updatedHandoffType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_user_handoff_notification_rule.foo", "handoff_type", updatedHandoffType),
				),
			},
		},
	})
}

func testAccCheckPagerDutyUserHandoffNotificationRuleDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_user_handoff_notification_rule" {
			continue
		}
		userID := r.Primary.Attributes["user_id"]
		ruleID := r.Primary.ID

		ctx := context.Background()
		if _, err := testAccProvider.client.GetUserOncallHandoffNotificationRuleWithContext(ctx, userID, ruleID); err == nil {
			return fmt.Errorf("user handoff notification rule still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyUserHandoffNotificationRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		userID := rs.Primary.Attributes["user_id"]
		ruleID := rs.Primary.ID

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user handoff notification rule ID is set")
		}

		ctx := context.Background()
		found, err := testAccProvider.client.GetUserOncallHandoffNotificationRuleWithContext(ctx, userID, ruleID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("user handoff notification rule not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyUserHandoffNotificationRuleConfig(userName, contactMethodType, handoffType string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[1]v@foo.test"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_user_contact_method" "phone" {
  user_id      = pagerduty_user.foo.id
  type         = "phone_contact_method"
  country_code = "+1"
  address      = "2025550199"
  label        = "Work"
}

resource "pagerduty_user_handoff_notification_rule" "foo" {
  user_id                   = pagerduty_user.foo.id
  handoff_type              = "%[3]v"
  notify_advance_in_minutes = 180
  contact_method {
    id   = pagerduty_user_contact_method.phone.id
    type = "%[2]v"
  }
}
`, userName, contactMethodType, handoffType)
}
