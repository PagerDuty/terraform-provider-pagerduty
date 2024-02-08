package pagerduty

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_incident_custom_fields", &resource.Sweeper{
		Name: "pagerduty_incident_custom_fields",
		F:    testSweepIncidentCustomField,
	})
}

func testSweepIncidentCustomField(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.IncidentCustomFields.ListContext(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, customField := range resp.Fields {
		if strings.HasPrefix(customField.Name, "tf_") {
			log.Printf("Destroying field %s (%s)", customField.Name, customField.ID)
			if _, err := client.IncidentCustomFields.DeleteContext(context.Background(), customField.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyIncidentCustomFields_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	description1 := acctest.RandString(10)
	description2 := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentCustomFieldConfig(fieldName, description1, "string"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentCustomFieldExists("pagerduty_incident_custom_field.input"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "name", fieldName),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "description", description1),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "data_type", "string"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentCustomFieldConfig(fieldName, description2, "string"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentCustomFieldExists("pagerduty_incident_custom_field.input"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "description", description2),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "data_type", "string"),
				),
			},
		},
	})
}

func TestAccPagerDutyIncidentCustomField_BasicWithDescription(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	description := acctest.RandString(30)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentCustomFieldConfigWithDescription(fieldName, description, "string"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentCustomFieldExists("pagerduty_incident_custom_field.input"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "name", fieldName),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field.input", "description", description),
				),
			},
		},
	})
}

func TestAccPagerDutyIncidentCustomFields_UnknownDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyIncidentCustomFieldConfig(fieldName, "", "garbage"),
				ExpectError: regexp.MustCompile("Unknown data_type garbage"),
			},
		},
	})
}

func TestAccPagerDutyIncidentCustomFields_IllegalDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyIncidentCustomFieldConfig(fieldName, "", pagerduty.IncidentCustomFieldDataTypeUnknown.String()),
				ExpectError: regexp.MustCompile("Unknown data_type unknown"),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentCustomFieldConfig(name, description, datatype string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  description = "%[2]s" 
  data_type = "%[3]s"
  field_type = "single_value_fixed"
}
`, name, description, datatype)
}

func testAccCheckPagerDutyIncidentCustomFieldConfigNoDescription(name, datatype string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  data_type = "%[2]s"
  field_type = "single_value_fixed"
}
`, name, datatype)
}

func testAccCheckPagerDutyIncidentCustomFieldConfigWithDescription(name, description, datatype string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  data_type = "%[2]s"
  description = "%[3]s"
  field_type = "single_value_fixed"
}
`, name, datatype, description)
}

func testAccCheckPagerDutyIncidentCustomFieldDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_incident_custom_field" {
			continue
		}

		if _, _, err := client.IncidentCustomFields.GetContext(context.Background(), r.Primary.ID, nil); err == nil {
			return fmt.Errorf("field still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyIncidentCustomFieldExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no field ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.IncidentCustomFields.GetContext(context.Background(), rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("field not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccPreCheckIncidentCustomFieldTests(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_INCIDENT_CUSTOM_FIELDS"); v == "" {
		t.Skip("PAGERDUTY_ACC_INCIDENT_CUSTOM_FIELDS not set. Skipping Incident Custom Field-related test")
	}
}
