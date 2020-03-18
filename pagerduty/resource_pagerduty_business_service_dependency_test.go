package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyBusinessServiceDependency_Basic(t *testing.T) {
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	businessService := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	businessServiceUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyBusinessServiceDependencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceExists("pagerduty_business_service_dependency.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service_dependency.foo", "relationship.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service_dependency.foo", "relationship.0.supporting_service.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service_dependency.foo", "relationship.0.dependent_service.#", "1"),
				),
			},
			{
				Config: testAccCheckPagerDutyBusinessServiceDependencyConfigUpdated(service, businessServiceUpdated, username, email, escalationPolicy, service2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyBusinessServiceExists("pagerduty_business_service_dependency.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service_dependency.foo", "relationship.#", "2"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service_dependency.foo", "relationship.0.supporting_service.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_business_service_dependency.foo", "relationship.0.dependent_service.#", "1"),
				),
			},
		},
	})
}
func testAccCheckPagerDutyBusinessServiceDependencyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Business Service ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)

		found, _, err := client.BusinessServices.GetDependencies()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Business Service not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyBusinessServiceDependencyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_business_service" {
			continue
		}

		if _, _, err := client.BusinessServices.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Business service still exists")
		}

	}
	return nil
}
func testAccCheckPagerDutyBusinessServiceDependencyConfig(service, businessService, username, email, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_business_service" "foo" {
	name = "%s"
}

resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}
resource "pagerduty_service" "foo" {
	name = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_business_service_dependency" "foo" {
	relationship {
		supporting_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		dependent_service {
			id = pagerduty_service.foo.id
			type = "service"
		}
	}
}
`, businessService, username, email, escalationPolicy, service)
}
func testAccCheckPagerDutyBusinessServiceDependencyConfigUpdated(service, businessService, username, email, escalationPolicy, service2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_business_service" "foo" {
	name = "%s"
}

resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "bar"
	num_loops   = 2
	rule {
		escalation_delay_in_minutes = 10
		target {
			type = "user_reference"
			id   = pagerduty_user.foo.id
		}
	}
}
resource "pagerduty_service" "foo" {
	name = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}

resource "pagerduty_service" "two" {
	name = "%s"
	description             = "two"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_business_service_dependency" "foo" {
	relationship {
		supporting_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		dependent_service {
			id = pagerduty_service.foo.id
			type = "service"
		}
	}
	relationship {
		supporting_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		dependent_service {
			id = pagerduty_service.two.id
			type = "service"
		}
	}
}
`, businessService, username, email, escalationPolicy, service, service2)
}
