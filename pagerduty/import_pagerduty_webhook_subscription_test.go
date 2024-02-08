package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyWebhookSubscription_import(t *testing.T) {
	description := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyWebhookSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyWebhookSubscriptionConfig(username, email, escalationPolicy, service, description),
			},

			{
				ResourceName:      "pagerduty_webhook_subscription.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
