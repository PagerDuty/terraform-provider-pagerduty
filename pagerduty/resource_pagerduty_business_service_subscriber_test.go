package pagerduty

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyBusinessServiceSubscriber_User(t *testing.T) {
	businessServiceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyBusinessServiceSubscriberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceSubscriberConfig(businessServiceName, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceSubscriberExists("pagerduty_business_service_subscriber.foo"),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", businessServiceName),
					resource.TestCheckResourceAttr("pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr("pagerduty_user.foo", "email", email),
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
					testAccCheckPagerDutyBusinessServiceSubscriberExists("pagerduty_business_service_subscriber.foo"),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", businessServiceName),
					resource.TestCheckResourceAttr("pagerduty_team.foo", "name", team),
				),
			},
		},
	})
}

func TestAccPagerDutyBusinessServiceSubscriber_TeamUser(t *testing.T) {
	businessServiceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyBusinessServiceSubscriberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceSubscriberTeamUserConfig(businessServiceName, team, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceSubscriberExists("pagerduty_business_service_subscriber.foo"),
					testAccCheckPagerDutyBusinessServiceSubscriberExists("pagerduty_business_service_subscriber.bar"),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", businessServiceName),
					resource.TestCheckResourceAttr("pagerduty_team.foo", "name", team),
					resource.TestCheckResourceAttr("pagerduty_business_service.foo", "name", businessServiceName),
					resource.TestCheckResourceAttr("pagerduty_user.bar", "name", username),
					resource.TestCheckResourceAttr("pagerduty_user.bar", "email", email),
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

func testAccCheckPagerDutyBusinessServiceSubscriberExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Business Service Subscriber ID is set")
		}
		ids := strings.Split(rs.Primary.ID, ".")

		businessServiceId, subscriberType, subscriberId := ids[0], ids[1], ids[2]

		client, _ := testAccProvider.Meta().(*Config).Client()
		response, _, err := client.BusinessServiceSubscribers.List(businessServiceId)
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
			return fmt.Errorf("Business Service %s subscriber not found: %s - %s", businessServiceId, subscriberId, subscriberType)
		}
		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceSubscriberConfig(businessServiceName string, username string, email string) string {
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

func testAccCheckPagerDutyBusinessServiceSubscriberTeamConfig(businessServiceName string, team string) string {
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

func testAccCheckPagerDutyBusinessServiceSubscriberTeamUserConfig(businessServiceName string, team string, username string, email string) string {
	return fmt.Sprintf(`
	resource "pagerduty_business_service" "foo" {
		name = "%s"
	}
	resource "pagerduty_team" "foo" {
		name = "%s"
	}
	resource "pagerduty_user" "bar" {
		name = "%s"
		email = "%s"
	}
	resource "pagerduty_business_service_subscriber" "foo" {
		subscriber_type = "team"
		subscriber_id = pagerduty_team.foo.id
		business_service_id = pagerduty_business_service.foo.id
	}
	resource "pagerduty_business_service_subscriber" "bar" {
		subscriber_type = "user"
		subscriber_id = pagerduty_user.bar.id
		business_service_id = pagerduty_business_service.foo.id
	}
`, businessServiceName, team, username, email)
}
