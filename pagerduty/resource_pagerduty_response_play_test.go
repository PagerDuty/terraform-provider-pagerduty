package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func testSweepResponsePlay(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.ResponsePlays.List()
	if err != nil {
		return err
	}

	for _, rplay := range resp.ResponsePlays {
		if strings.HasPrefix(rplay.Name, "test") || strings.HasPrefix(rplay.Name, "tf-") {
			log.Printf("Destroying response play %s (%s)", rplay.Name, rplay.ID)
			if _, err := client.ResponsePlays.Delete(rplay.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyResponsePlay_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyResponsePlayDestroy, // really?
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyResponsePlayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyResponsePlayExists("pagerduty_response_play.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyResponsePlayConfigUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyResponsePlayExists("pagerduty_response_play.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyResponsePlayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_response_play" {
			continue
		}

		if _, _, err := client.ResponsePlays.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("response play still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyResponsePlayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No response play ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)

		found, _, err := client.ResponsePlays.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("response play not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyResponsePlayConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[1]v@foo.com"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
}

resource "pagerduty_escalation_policy" "foo" {
  name        = "%[1]v"
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

resource "pagerduty_response_play" "foo" {
  name = "%[1]v"
  responder {
	  type = "user_reference"
	  id = pagerduty_user.foo.id
  }
  responder {
	  type = "escalation_policy_reference"
	  id = pagerduty_escalation_policy.foo.id
  }
  runnability = "services"

}
`, name)
}

func testAccCheckPagerDutyResponsePlayConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
	name        = "%[1]v"
	email       = "%[1]v@foo.com"
	color       = "green"
	role        = "user"
	job_title   = "foo"
	description = "foo"
}
	
resource "pagerduty_escalation_policy" "foo" {
	name        = "%[1]v"
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
	
resource "pagerduty_response_play" "foo" {
	name = "%[1]v"
	responder {
		type = "escalation_policy_reference"
		id = pagerduty_escalation_policy.foo.id
	}
	subscriber {
		type = "user_reference"
		id = pagerduty_user.foo.id
	}
	runnability = "services"
}
`, name)
}
