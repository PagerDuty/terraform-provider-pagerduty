package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_extension_servicenow", &resource.Sweeper{
		Name: "pagerduty_extension_servicenow",
		F:    testSweepExtensionServiceNow,
	})
}

func testSweepExtensionServiceNow(_ string) error {
	ctx := context.Background()

	resp, err := testAccProvider.client.ListExtensionsWithContext(ctx, pagerduty.ListExtensionOptions{})
	if err != nil {
		return err
	}

	for _, extension := range resp.Extensions {
		if strings.HasPrefix(extension.Name, "test") || strings.HasPrefix(extension.Name, "tf-") {
			log.Printf("Destroying extension %s (%s)", extension.Name, extension.ID)
			if err := testAccProvider.client.DeleteExtensionWithContext(ctx, extension.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyExtensionServiceNow_Basic(t *testing.T) {
	extensionName := id.PrefixedUniqueId("tf-")
	extensionNameUpdated := id.PrefixedUniqueId("tf-")
	name := id.PrefixedUniqueId("tf-")
	url := "https://example.com/receive_a_pagerduty_webhook"
	urlUpdated := "https://example.com/webhook_foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyExtensionServiceNowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyExtensionServiceNowConfig(name, extensionName, url, "false", "any"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyExtensionServiceNowExists("pagerduty_extension_servicenow.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "name", extensionName),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "extension_schema", "PJFWPEP"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "endpoint_url", url),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "html_url", ""),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "snow_user", "meeps"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "snow_password", "zorz"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "sync_options", "manual_sync"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "target", "foo.servicenow.com/webhook_foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "task_type", "incident"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "referer", "None"),
				),
			},
			{
				Config: testAccCheckPagerDutyExtensionServiceNowConfig(name, extensionNameUpdated, urlUpdated, "true", "pd-users"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyExtensionServiceNowExists("pagerduty_extension_servicenow.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "name", extensionNameUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "extension_schema", "PJFWPEP"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "endpoint_url", urlUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "html_url", ""),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "snow_user", "meeps"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "snow_password", "zorz"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "sync_options", "manual_sync"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "target", "foo.servicenow.com/webhook_foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "task_type", "incident"),
					resource.TestCheckResourceAttr(
						"pagerduty_extension_servicenow.foo", "referer", "None"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyExtensionServiceNowDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_extension_servicenow" {
			continue
		}

		if _, err := testAccProvider.client.GetExtensionWithContext(ctx, r.Primary.ID); err == nil {
			return fmt.Errorf("Extension still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyExtensionServiceNowExists(n string) resource.TestCheckFunc {
	ctx := context.Background()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No extension ID is set")
		}

		found, err := testAccProvider.client.GetExtensionWithContext(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Extension not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyExtensionServiceNowConfig(name string, extensionName string, url string, _ string, _ string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%[1]v"
  email       = "%[1]v@foo.test"
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

resource "pagerduty_service" "foo" {
  name                    = "%[1]v"
  description             = "foo"
  auto_resolve_timeout    = 1800
  acknowledgement_timeout = 1800
  escalation_policy       = pagerduty_escalation_policy.foo.id

  incident_urgency_rule {
    type    = "constant"
    urgency = "high"
  }
}

data "pagerduty_extension_schema" "foo" {
	name = "Generic V2 Webhook"
}

resource "pagerduty_extension_servicenow" "foo"{
  name = "%s"
  endpoint_url = "%s"
  extension_schema = data.pagerduty_extension_schema.foo.id
  extension_objects = [pagerduty_service.foo.id]
  snow_user = "meeps"
  snow_password = "zorz"
  sync_options = "manual_sync"
  target = "foo.servicenow.com/webhook_foo"
  task_type = "incident"
  referer = "None"
}

`, name, extensionName, url)
}
