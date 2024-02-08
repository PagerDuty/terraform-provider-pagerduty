package pagerduty

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyIncidentCustomFieldOptions_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue2 := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataType := pagerduty.IncidentCustomFieldDataTypeString

	testAccExecuteIncidentCustomFieldOptionTest(t, fieldName, dataType, fieldOptionValue, fieldOptionValue2)
}

func TestAccPagerDutyIncidentCustomFieldOptions_InvalidDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataType := pagerduty.IncidentCustomFieldDataTypeInt

	testAccExecuteIncidentCustomFieldOptionTestError(t, fieldName, dataType, fieldOptionValue,
		regexp.MustCompile(`Error: "integer" is an invalid value. Must be one of \[]string{"string"}`))
}

func testAccExecuteIncidentCustomFieldOptionTest(t *testing.T, fieldName string, dataType pagerduty.IncidentCustomFieldDataType, fieldOptionValue, fieldOptionValueForUpdate string) {
	var fieldID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)

			field := testAccCreateTestPagerDutyIncidentCustomFieldForFieldOption(fieldName, dataType)
			fieldID = field.ID
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckPagerDutyIncidentCustomFieldOptionDestroy(state)
			if err != nil {
				return err
			}
			return testAccDeleteTestPagerDutyIncidentCustomFieldForFieldOption(fieldID)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentCustomFieldOptionConfig(fieldName, dataType, fieldOptionValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentCustomFieldOptionExists("pagerduty_incident_custom_field_option.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field_option.test", "data_type", dataType.String()),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_incident_custom_field_option.test", "field", fieldID))(state)
					},
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field_option.test", "value", fieldOptionValue),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentCustomFieldOptionConfig(fieldName, dataType, fieldOptionValueForUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentCustomFieldOptionExists("pagerduty_incident_custom_field_option.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field_option.test", "data_type", dataType.String()),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_incident_custom_field_option.test", "field", fieldID))(state)
					},
					resource.TestCheckResourceAttr(
						"pagerduty_incident_custom_field_option.test", "value", fieldOptionValueForUpdate),
				),
			},
		},
	})
}

func testAccExecuteIncidentCustomFieldOptionTestError(t *testing.T, fieldName string, dataType pagerduty.IncidentCustomFieldDataType, fieldOptionValue string, errorRegex *regexp.Regexp) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyIncidentCustomFieldOptionConfigForErrorCases(fieldName, dataType, fieldOptionValue),
				ExpectError: errorRegex,
			},
		},
	})
}

func testAccCreateTestPagerDutyIncidentCustomFieldForFieldOption(fieldName string, dataType pagerduty.IncidentCustomFieldDataType) *pagerduty.IncidentCustomField {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	field, _, err := client.IncidentCustomFields.CreateContext(context.Background(), &pagerduty.IncidentCustomField{
		DataType:    dataType,
		Name:        fieldName,
		DisplayName: fieldName,
		FieldType:   pagerduty.IncidentCustomFieldFieldTypeSingleValueFixed,
	})
	if err != nil {
		panic("Unable to create test field to contain option")
	}
	return field
}

func testAccDeleteTestPagerDutyIncidentCustomFieldForFieldOption(fieldID string) error {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	_, err := client.IncidentCustomFields.DeleteContext(context.Background(), fieldID)
	return err
}

func testAccCheckPagerDutyIncidentCustomFieldOptionConfig(name string, dataType pagerduty.IncidentCustomFieldDataType, fieldOptionValue string) string {
	return fmt.Sprintf(`
data "pagerduty_incident_custom_field" "input" {
  name = "%s"
}

resource "pagerduty_incident_custom_field_option" "test" {
  field = data.pagerduty_incident_custom_field.input.id
  data_type = "%s"
  value = "%s"
}

`, name, dataType.String(), fieldOptionValue)
}

func testAccCheckPagerDutyIncidentCustomFieldOptionConfigForErrorCases(name string, dataType pagerduty.IncidentCustomFieldDataType, fieldOptionValue string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  data_type = "%[2]s"
  field_type = "single_value_fixed"
}

resource "pagerduty_incident_custom_field_option" "test" {
  field = pagerduty_incident_custom_field.input.id
  data_type = "%[2]s"
  value = "%[3]s"
}

`, name, dataType.String(), fieldOptionValue)
}

func testAccCheckPagerDutyIncidentCustomFieldOptionDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_incident_custom_field_option" {
			continue
		}

		fieldID := r.Primary.Attributes["field"]
		if _, _, err := client.IncidentCustomFields.GetFieldOptionContext(context.Background(), fieldID, r.Primary.ID); err == nil {
			return fmt.Errorf("field still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyIncidentCustomFieldOptionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no field option ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		fieldID := rs.Primary.Attributes["field"]
		found, _, err := client.IncidentCustomFields.GetFieldOptionContext(context.Background(), fieldID, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("field option not found: %v/%v - %v", fieldID, rs.Primary.ID, found)
		}

		return nil
	}
}
