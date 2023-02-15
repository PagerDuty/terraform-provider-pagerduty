package pagerduty

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyCustomFieldOptions_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue2 := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataType := pagerduty.CustomFieldDataTypeString

	testAccExecuteCustomFieldOptionTest(t, fieldName, dataType, fieldOptionValue, fieldOptionValue2)
}

func TestAccPagerDutyCustomFieldOptions_Integer(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("%v", acctest.RandIntRange(1, 500))
	fieldOptionValue2 := fmt.Sprintf("%v", acctest.RandIntRange(1, 500))
	dataType := pagerduty.CustomFieldDataTypeInt

	testAccExecuteCustomFieldOptionTest(t, fieldName, dataType, fieldOptionValue, fieldOptionValue2)
}

func TestAccPagerDutyCustomFieldOptions_Integer_Bad(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataType := pagerduty.CustomFieldDataTypeInt

	testAccExecuteCustomFieldOptionTestError(t, fieldName, dataType, fieldOptionValue,
		regexp.MustCompile(fmt.Sprintf("invalid value for datatype integer: %s", fieldOptionValue)))
}

func TestAccPagerDutyCustomFieldOptions_Integer_Bad_Float(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	floatValue := float64(acctest.RandIntRange(1, 500)) + 0.5
	fieldOptionValue := fmt.Sprintf("%v", floatValue)
	dataType := pagerduty.CustomFieldDataTypeInt

	testAccExecuteCustomFieldOptionTestError(t, fieldName, dataType, fieldOptionValue,
		regexp.MustCompile(fmt.Sprintf("invalid value for datatype integer: %s", fieldOptionValue)))
}

func TestAccPagerDutyCustomFieldOptions_Float(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	floatValue := float64(acctest.RandIntRange(1, 500)) + 0.5
	fieldOptionValue := fmt.Sprintf("%v", floatValue)
	floatValue2 := float64(acctest.RandIntRange(1, 500)) + 0.5
	fieldOptionValue2 := fmt.Sprintf("%v", floatValue2)
	dataType := pagerduty.CustomFieldDataTypeFloat

	testAccExecuteCustomFieldOptionTest(t, fieldName, dataType, fieldOptionValue, fieldOptionValue2)
}

func TestAccPagerDutyCustomFieldOptions_Float_Bad(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	fieldOptionValue := fmt.Sprintf("tf_%s", acctest.RandString(5))
	dataType := pagerduty.CustomFieldDataTypeFloat

	testAccExecuteCustomFieldOptionTestError(t, fieldName, dataType, fieldOptionValue,
		regexp.MustCompile(fmt.Sprintf("invalid value for datatype float: %s", fieldOptionValue)))
}

func testAccExecuteCustomFieldOptionTest(t *testing.T, fieldName string, dataType pagerduty.CustomFieldDataType, fieldOptionValue, fieldOptionValueForUpdate string) {
	var fieldID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)

			field := testAccCreateTestPagerDutyCustomFieldForFieldOption(fieldName, dataType)
			fieldID = field.ID
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckPagerDutyCustomFieldOptionDestroy(state)
			if err != nil {
				return err
			}
			return testAccDeleteTestPagerDutyCustomFieldForFieldOption(fieldID)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyCustomFieldOptionConfig(fieldName, dataType, fieldOptionValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldOptionExists("pagerduty_custom_field_option.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_option.test", "datatype", dataType.String()),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_option.test", "field", fieldID))(state)
					},
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_option.test", "value", fieldOptionValue),
				),
			},
			{
				Config: testAccCheckPagerDutyCustomFieldOptionConfig(fieldName, dataType, fieldOptionValueForUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldOptionExists("pagerduty_custom_field_option.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_option.test", "datatype", dataType.String()),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_option.test", "field", fieldID))(state)
					},
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_option.test", "value", fieldOptionValueForUpdate),
				),
			},
		},
	})
}

func testAccExecuteCustomFieldOptionTestError(t *testing.T, fieldName string, dataType pagerduty.CustomFieldDataType, fieldOptionValue string, errorRegex *regexp.Regexp) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldOptionConfigForErrorCases(fieldName, dataType, fieldOptionValue),
				ExpectError: errorRegex,
			},
		},
	})
}

func testAccCreateTestPagerDutyCustomFieldForFieldOption(fieldName string, dataType pagerduty.CustomFieldDataType) *pagerduty.CustomField {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	field, _, err := client.CustomFields.Create(&pagerduty.CustomField{
		DataType:     dataType,
		Name:         fieldName,
		DisplayName:  fieldName,
		FixedOptions: true,
	})
	if err != nil {
		panic("Unable to create test field to contain option")
	}
	return field
}

func testAccDeleteTestPagerDutyCustomFieldForFieldOption(fieldID string) error {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	_, err := client.CustomFields.Delete(fieldID)
	return err
}

func testAccCheckPagerDutyCustomFieldOptionConfig(name string, dataType pagerduty.CustomFieldDataType, fieldOptionValue string) string {
	return fmt.Sprintf(`
data "pagerduty_custom_field" "input" {
  name = "%s"
}

resource "pagerduty_custom_field_option" "test" {
  field = data.pagerduty_custom_field.input.id
  datatype = "%s"
  value = "%s"
}

`, name, dataType.String(), fieldOptionValue)
}

func testAccCheckPagerDutyCustomFieldOptionConfigForErrorCases(name string, dataType pagerduty.CustomFieldDataType, fieldOptionValue string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "%[2]s"
  fixed_options = true
}

resource "pagerduty_custom_field_option" "test" {
  field = pagerduty_custom_field.input.id
  datatype = "%[2]s"
  value = "%[3]s"
}

`, name, dataType.String(), fieldOptionValue)
}

func testAccCheckPagerDutyCustomFieldOptionDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_custom_field_option" {
			continue
		}

		fieldID := r.Primary.Attributes["field"]
		if _, _, err := client.CustomFields.GetFieldOption(fieldID, r.Primary.ID); err == nil {
			return fmt.Errorf("field still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyCustomFieldOptionExists(n string) resource.TestCheckFunc {
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
		found, _, err := client.CustomFields.GetFieldOption(fieldID, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("field option not found: %v/%v - %v", fieldID, rs.Primary.ID, found)
		}

		return nil
	}
}
