package pagerduty

import (
	"context"
	"fmt"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyIncidentType_Basic(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	displayName := fmt.Sprintf("Terraform Test Incident Type %s", acctest.RandString(5))
	parentType := "incident_default"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentTypeConfig(name, displayName, parentType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentTypeExists("pagerduty_incident_type.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "display_name", displayName),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "parent_type", parentType),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentTypeConfigUpdated(name, displayName+"_updated", parentType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentTypeExists("pagerduty_incident_type.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "display_name", displayName+"_updated"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "parent_type", parentType),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentTypeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Incident Type ID is set")
		}

		client := testAccProvider.client
		ctx := context.Background()
		found, err := client.GetIncidentType(ctx, rs.Primary.ID, pagerduty.GetIncidentTypeOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Incident Type not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyIncidentTypeConfig(name, displayName, parentType string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_type" "test" {
  name         = "%s"
  display_name = "%s"
  parent_type  = "%s"
  description  = "Terraform test incident type"
  enabled      = true
}
`, name, displayName, parentType)
}

func testAccCheckPagerDutyIncidentTypeConfigUpdated(name, displayName, parentType string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_type" "test" {
  name         = "%s"
  display_name = "%s"
  parent_type  = "%s"
  description  = "Terraform test incident type updated"
  enabled      = false
}
`, name, displayName, parentType)
}
