package pagerduty

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPagerDutyBusinessServiceSubscriber_User(t *testing.T) {
	businessServiceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyBusinessServiceSubscriberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceSubscriberConfig(businessServiceName, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceSubscriberExists("pagerduty_business_service_subscriber.foo", "pagerduty_business_service.foo", "user"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service.foo", "name", businessServiceName),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", email),
				),
			},
		},
	})
}
func TestAccPagerDutyBusinessServiceSubscriber_Team(t *testing.T) {
	businessServiceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyBusinessServiceSubscriberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceSubscriberTeamConfig(businessServiceName, team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceSubscriberExists("pagerduty_business_service_subscriber.foo", "pagerduty_business_service.foo", "team"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service.foo", "name", businessServiceName),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", team),
				),
			},
		},
	})
}

func testAccCheckPagerDutyBusinessServiceSubscriberDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_business_service" {
			continue
		}
		ids := strings.Split(r.Primary.ID, ".")

		businessServiceID := ids[0]

		response, _, err := client.BusinessServiceSubscribers.List(businessServiceID)
		if err != nil {
			// if there are no subscriber for the entity that's okay
			return nil
		}
		// find subscriber the test created
		for _, subscriber := range response.BusinessServiceSubscribers {
			if subscriber.ID != "" {
				return fmt.Errorf("Subscriber %s still exists and is connected to ID %s", subscriber.ID, businessServiceID)
			}
		}
	}
	return nil
}

func testAccCheckPagerDutyBusinessServiceSubscriberExists(n, b, subscriberType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		bs, ok := s.RootModule().Resources[b]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Business Service Subscriber ID is set")
		}
		if bs.Primary.ID == "" {
			return fmt.Errorf("No Business Service ID is set")
		}

		subscriberId := rs.Primary.ID
		businessServiceID := bs.Primary.ID

		client, _ := testAccProvider.Meta().(*Config).Client()
		response, _, err := client.BusinessServiceSubscribers.List(businessServiceID)
		if err != nil {
			return err
		}
		// find tag the test created
		var isFound bool = false
		for _, subscriber := range response.BusinessServiceSubscribers {
			if subscriber.ID == subscriberId && subscriber.Type == subscriberType {
				isFound = true
				break
			}
		}
		if !isFound {
			return fmt.Errorf("Subscriber %s type %s still exists and is connected to ID %s", subscriberId, subscriberType, businessServiceID)
		}
		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceSubscriberConfig(businessServiceName, username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_business_service" "foo" {
	name = "%s"
}
resource "pagerduty_user" "foo" {
	name = "%s"
	email = "%s"
}
resource "pagerduty_business_service_subscriber" "foo" {
	subscriber_type = "user"
	subscriber_id = pagerduty_user.foo.id
	business_service_id = pagerduty_business_service.foo.id
}
`, businessServiceName, username, email)
}

func testAccCheckPagerDutyBusinessServiceSubscriberTeamConfig(businessServiceName, team string) string {
	return fmt.Sprintf(`
	resource "pagerduty_business_service" "foo" {
		name = "%s"
	}
	resource "pagerduty_team" "foo" {
		name = "%s"
	}
	resource "pagerduty_business_service_subscriber" "foo" {
		subscriber_type = "team"
		subscriber_id = pagerduty_team.foo.id
		business_service_id = pagerduty_business_service.foo.id
	}
`, businessServiceName, team)
}
