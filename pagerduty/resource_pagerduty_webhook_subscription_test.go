package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_webhook_subscription", &resource.Sweeper{
		Name: "pagerduty_webhook_subscription",
		F:    testSweepWebhookSubscription,
	})
}

func testSweepWebhookSubscription(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.WebhookSubscriptions.List()
	if err != nil {
		return err
	}

	for _, webhook := range resp.WebhookSubscriptions {
		if strings.HasPrefix(webhook.Description, "tf-test-") {
			log.Printf("Destroying webhook subscription %s ", webhook.ID)
			if _, err := client.WebhookSubscriptions.Delete(webhook.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
func TestAccPagerDutyWebhookSubscription_Basic(t *testing.T) {
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyWebhookSubscriptionExists("pagerduty_webhook_subscription.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_webhook_subscription.foo", "description", description),
					resource.TestCheckResourceAttr(
						"pagerduty_webhook_subscription.foo", "events.#", "13"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyWebhookSubscriptionDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_webhook_subscription" {
			continue
		}
		if _, _, err := client.WebhookSubscriptions.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Webhook subscription still exists")
		}
	}
	return nil
}
func testAccCheckPagerDutyWebhookSubscriptionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Webhook Subscription ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.WebhookSubscriptions.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Webhook subscription not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyWebhookSubscriptionConfig(username, useremail, escalationPolicy, service, description string) string {
	return fmt.Sprintf(`
	resource "pagerduty_user" "foo" {
		name        = "%s"
		email       = "%s"
	}

	resource "pagerduty_escalation_policy" "foo" {
		name        = "%s"
		description = "foo"
		num_loops   = 1

		rule {
			escalation_delay_in_minutes = 10

			target {
				type = "user_reference"
				id   = pagerduty_user.foo.id
			}
		}
	}

	resource "pagerduty_service" "foo" {
		name                    = "%s"
		description             = "foo"
		auto_resolve_timeout    = 1800
		acknowledgement_timeout = 1800
		escalation_policy       = pagerduty_escalation_policy.foo.id

		incident_urgency_rule {
			type = "constant"
			urgency = "high"
		}
	}

	resource "pagerduty_webhook_subscription" "foo" {
		delivery_method {
			type = "http_delivery_method"
			url = "https://example.com/receive_a_pagerduty_webhook"
			custom_header {
				name = "X-Foo"
				value = "foo"
			}
		}
		description = "%s"
		events = [
            "incident.acknowledged",
            "incident.annotated",
            "incident.delegated",
            "incident.escalated",
            "incident.priority_updated",
            "incident.reassigned",
            "incident.reopened",
            "incident.resolved",
            "incident.responder.added",
            "incident.responder.replied",
            "incident.status_update_published",
            "incident.triggered",
            "incident.unacknowledged"
		]
		active = true
		filter {
			id = pagerduty_service.foo.id
			type = "service_reference"
		}
		type = "webhook_subscription"
	}
	`, username, useremail, escalationPolicy, service, description)
}
