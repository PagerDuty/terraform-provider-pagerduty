package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_incident_type_custom_field", &resource.Sweeper{
		Name: "pagerduty_incident_type_custom_field",
		F: func(_ string) error {
			ctx := context.Background()

			resp1, err := testAccProvider.client.ListIncidentTypes(ctx, pagerduty.ListIncidentTypesOptions{})
			if err != nil {
				return err
			}

			sweepForIncidentType := func(incidentType string) error {
				resp, err := testAccProvider.client.ListIncidentTypeFields(ctx, incidentType, pagerduty.ListIncidentTypeFieldsOptions{})
				if err != nil {
					return err
				}

				for _, f := range resp.Fields {
					if strings.HasPrefix(f.Name, "test") || strings.HasPrefix(f.Name, "tf_") {
						log.Printf("Destroying add-on %s (%s)", f.Name, f.ID)
						if err := testAccProvider.client.DeleteIncidentTypeField(ctx, incidentType, f.ID); err != nil {
							return err
						}
					}
				}

				return nil
			}

			for _, it := range resp1.IncidentTypes {
				if strings.HasPrefix(it.Name, "test") || strings.HasPrefix(it.Name, "tf_") {
					if err := sweepForIncidentType(it.ID); err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

func TestAccPagerDutyIncidentTypeCustomField_Basic(t *testing.T) {
	name := fmt.Sprintf("tf_%s", acctest.RandString(5))
	nameUpdated := fmt.Sprintf("tf_%s", acctest.RandString(5))
	enabledUpdated := "false"
	descriptionUpdated := fmt.Sprintf("tf_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyIncidentTypeCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentTypeCustomFieldConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentTypeCustomFieldExists("pagerduty_incident_type_custom_field.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "display_name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "enabled", "true"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_incident_type_custom_field.foo", "incident_type"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "field_type", "single_value_fixed"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "field_options.0", "hello"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "field_options.1", "hi"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentTypeCustomFieldConfigUpdated(name, nameUpdated, enabledUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentTypeCustomFieldExists("pagerduty_incident_type_custom_field.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "name", name),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "display_name", nameUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "enabled", enabledUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "description", descriptionUpdated),
					resource.TestCheckResourceAttrSet(
						"pagerduty_incident_type_custom_field.foo", "incident_type"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "field_type", "single_value_fixed"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "field_options.0", "hello"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "field_options.1", "hi"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_type_custom_field.foo", "default_value", "\"hi\""),
				),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentTypeCustomFieldDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pagerduty_incident_type_custom_field" {
			continue
		}

		ctx := context.Background()
		parts := strings.Split(rs.Primary.ID, ":")
		incidentType, id := parts[0], parts[1]
		if _, err := testAccProvider.client.GetIncidentTypeField(ctx, incidentType, id, pagerduty.GetIncidentTypeFieldOptions{}); err == nil {
			return fmt.Errorf("Incident type custom field still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyIncidentTypeCustomFieldExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		parts := strings.Split(rs.Primary.ID, ":")
		incidentType, id := parts[0], parts[1]
		_, err := testAccProvider.client.GetIncidentTypeField(ctx, incidentType, id, pagerduty.GetIncidentTypeFieldOptions{})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckPagerDutyIncidentTypeCustomFieldConfig(name string) string {
	incidentType := fmt.Sprintf("tf_%s", acctest.RandString(5))
	return fmt.Sprintf(`
resource "pagerduty_incident_type" "a" {
  name = "%s"
  display_name = "%[1]s"
  parent_type = "incident_default"
}
resource "pagerduty_incident_type_custom_field" "foo" {
  name = "%s"
  display_name = "%[2]s"
  data_type = "string"
  field_options = ["hello", "hi"]
  field_type = "single_value_fixed"
  incident_type = pagerduty_incident_type.a.id
}`, incidentType, name)
}

func testAccCheckPagerDutyIncidentTypeCustomFieldConfigUpdated(name, nameUpdated, enabledUpdated, descriptionUpdated string) string {
	incidentType := fmt.Sprintf("tf_%s", acctest.RandString(5))
	return fmt.Sprintf(`
resource "pagerduty_incident_type" "a" {
  name = "%s"
  display_name = "%[1]s"
  parent_type = "incident_default"
}
resource "pagerduty_incident_type_custom_field" "foo" {
  name = "%s"
  display_name = "%s"
  data_type = "string"
  field_options = ["hello", "hi"]
  field_type = "single_value_fixed"
  incident_type = pagerduty_incident_type.a.id
  enabled = %s
  description = "%s"
  default_value = jsonencode("hi")
}`, incidentType, name, nameUpdated, enabledUpdated, descriptionUpdated)
}
