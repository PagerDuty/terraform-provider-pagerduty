package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyUserNotificationRule_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	usernameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)
	emailUpdated := fmt.Sprintf("%s@foo.com", usernameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserNotificationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserNotificationRuleConfig(username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserNotificationRuleExists("pagerduty_user_notification_rule.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyUserNotificationRuleConfigUpdated(usernameUpdated, emailUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserNotificationRuleExists("pagerduty_user_notification_rule.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyUserNotificationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_user_notification_rule" {
			continue
		}

		userID, ruleID := resourcePagerDutyUserNotificationRuleParseID(r.Primary.ID)

		if _, _, err := client.Users.GetNotificationRule(userID, ruleID); err == nil {
			return fmt.Errorf("User notification rule still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyUserNotificationRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user notification rule ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)

		userID, ruleID := resourcePagerDutyUserNotificationRuleParseID(rs.Primary.ID)

		found, _, err := client.Users.GetNotificationRule(userID, ruleID)
		if err != nil {
			return err
		}

		if found.ID != ruleID {
			return fmt.Errorf("Notification rule not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyUserNotificationRuleConfig(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id = "${pagerduty_user.foo.id}"
  type    = "email_contact_method"
  address = "%[1]v%[2]v"
  label   = "%[1]v"
}

resource "pagerduty_user_notification_rule" "foo" {
	user_id = "${pagerduty_user.foo.id}"
	contact_method_id = "${pagerduty_user_contact_method.foo.contact_method_id}"
	contact_method_type = "${pagerduty_user_contact_method.foo.type}"
	start_delay_in_minutes = 10
	urgency = "high"
}
`, username, email)
}

func testAccCheckPagerDutyUserNotificationRuleConfigUpdated(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[2]v"
  color       = "red"
  role        = "user"
  job_title   = "bar"
  description = "bar"
}

resource "pagerduty_user_contact_method" "foo" {
  user_id = "${pagerduty_user.foo.id}"
  type    = "email_contact_method"
  address = "%[1]v%[2]v"
  label   = "%[1]v"
}

resource "pagerduty_user_notification_rule" "foo" {
	user_id = "${pagerduty_user.foo.id}"
	contact_method_id = "${pagerduty_user_contact_method.foo.contact_method_id}"
	contact_method_type = "${pagerduty_user_contact_method.foo.type}"
	start_delay_in_minutes = 10
	urgency = "high"
}

`, username, email)
}
