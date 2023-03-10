package pagerduty

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestPagerDutyCustomField_ValidateDefaultFieldValue(t *testing.T) {
	var testData = []struct {
		datatype         pagerduty.CustomFieldDataType
		multiValue       bool
		value            string
		expectError      bool
		expectedErrorStr string
	}{{
		datatype:         pagerduty.CustomFieldDataTypeInt,
		multiValue:       false,
		value:            "default!!!",
		expectedErrorStr: `invalid default value for datatype integer: default!!!`,
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeInt,
		multiValue:  false,
		value:       "42",
		expectError: false,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeInt,
		multiValue:       true,
		value:            "42",
		expectedErrorStr: `invalid default value for datatype integer (multi-value): 42`,
		expectError:      true,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeInt,
		multiValue:       true,
		value:            `[42,"foo"]`,
		expectedErrorStr: `invalid default value for datatype integer (multi-value): [42,"foo"]`,
		expectError:      true,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeFloat,
		multiValue:       true,
		value:            `[50.5,"foo"]`,
		expectedErrorStr: `invalid default value for datatype float (multi-value): [50.5,"foo"]`,
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeFloat,
		multiValue:  true,
		value:       `[50.5,42.3]`,
		expectError: false,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeFloat,
		multiValue:       false,
		value:            "default!!!",
		expectedErrorStr: `invalid default value for datatype float: default!!!`,
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeFloat,
		multiValue:  false,
		value:       "50.5",
		expectError: false,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeFloat,
		multiValue:       true,
		value:            "50.5",
		expectedErrorStr: `invalid default value for datatype float (multi-value): 50.5`,
		expectError:      true,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeFloat,
		multiValue:       true,
		value:            `[50.5,"foo"]`,
		expectedErrorStr: `invalid default value for datatype float (multi-value): [50.5,"foo"]`,
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeFloat,
		multiValue:  true,
		value:       `[50.5,42.3]`,
		expectError: false,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeBool,
		multiValue:       false,
		value:            "default!!!",
		expectedErrorStr: `invalid default value for datatype boolean: default!!!`,
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeBool,
		multiValue:  false,
		value:       "True",
		expectError: false,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeBool,
		multiValue:  false,
		value:       "true",
		expectError: false,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeBool,
		multiValue:  false,
		value:       "false",
		expectError: false,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeUrl,
		multiValue:       false,
		value:            "default!!!",
		expectedErrorStr: `parsed url default value "default!!!" is not an absolute url`,
		expectError:      true,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeUrl,
		multiValue:       false,
		value:            "\x01",
		expectedErrorStr: "invalid control character in URL",
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeUrl,
		multiValue:  false,
		value:       "https://www.pagerduty.com/",
		expectError: false,
	}, {
		datatype:         pagerduty.CustomFieldDataTypeDateTime,
		multiValue:       false,
		value:            "default!!!",
		expectedErrorStr: `parsing time "default!!!"`,
		expectError:      true,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeDateTime,
		multiValue:  false,
		value:       "2022-01-03T15:04:05Z",
		expectError: false,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeFieldOption,
		multiValue:  false,
		value:       "anything goes -- validated server-side",
		expectError: false,
	}, {
		datatype:    pagerduty.CustomFieldDataTypeFieldOption,
		multiValue:  true,
		value:       "something",
		expectError: true,
	}}

	for _, td := range testData {
		name := fmt.Sprintf("Test value: %s for type %v, multi-value: %v", td.value, td.datatype, td.multiValue)
		t.Run(name, func(t *testing.T) {
			err := validateDefaultFieldValue(td.value, td.datatype, td.multiValue)
			if td.expectError {
				if err == nil {
					t.Errorf("expected error, but did not receive one")
				} else if !strings.Contains(err.Error(), td.expectedErrorStr) {
					t.Errorf("expected error to contain %s but did not. was %s", td.expectedErrorStr, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

func TestAccPagerDutyCustomFieldConfiguration_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))
	var fieldID string
	var schemaID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)

			field := testAccCreateTestPagerDutyCustomFieldForFieldConfiguration(fieldName)
			fieldID = field.ID

			schema := testAccCreateTestPagerDutyCustomFieldSchemaForFieldConfiguration(schemaTitle)
			schemaID = schema.ID
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckPagerDutyCustomFieldConfigurationDestroy(state)
			if err != nil {
				return err
			}
			err = testAccDeleteTestPagerDutyCustomFieldForFieldConfiguration(fieldID)
			if err != nil {
				return err
			}
			return testAccDeleteTestPagerDutyCustomFieldSchemaForFieldConfiguration(schemaID)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyCustomFieldConfigurationConfigBasic(fieldName, schemaTitle),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldSchemaFieldConfigurationExists("pagerduty_custom_field_schema_field_configuration.test"),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_schema_field_configuration.test", "field", fieldID))(state)
					},
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_schema_field_configuration.test", "schema", schemaID))(state)
					},
				),
			},
		},
	})
}

func TestAccPagerDutyCustomFieldConfiguration_Required(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))
	var fieldID string
	var schemaID string

	config := fmt.Sprintf(`
data "pagerduty_custom_field" "input" {
  name = "%s"
}

data "pagerduty_custom_field_schema" "input" {
  title = "%s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                  = data.pagerduty_custom_field.input.id
  schema                 = data.pagerduty_custom_field_schema.input.id
  required               = true
  default_value          = "test"
  default_value_datatype = "string"
}

`, fieldName, schemaTitle)

	configForUpdate := fmt.Sprintf(`
data "pagerduty_custom_field" "input" {
  name = "%s"
}

data "pagerduty_custom_field_schema" "input" {
  title = "%s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                  = data.pagerduty_custom_field.input.id
  schema                 = data.pagerduty_custom_field_schema.input.id
  required               = true
  default_value          = "updated"
  default_value_datatype = "string"
}

`, fieldName, schemaTitle)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)

			field := testAccCreateTestPagerDutyCustomFieldForFieldConfiguration(fieldName)
			fieldID = field.ID

			schema := testAccCreateTestPagerDutyCustomFieldSchemaForFieldConfiguration(schemaTitle)
			schemaID = schema.ID
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			err := testAccCheckPagerDutyCustomFieldConfigurationDestroy(state)
			if err != nil {
				return err
			}
			err = testAccDeleteTestPagerDutyCustomFieldForFieldConfiguration(fieldID)
			if err != nil {
				return err
			}
			return testAccDeleteTestPagerDutyCustomFieldSchemaForFieldConfiguration(schemaID)
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldSchemaFieldConfigurationExists("pagerduty_custom_field_schema_field_configuration.test"),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_schema_field_configuration.test", "field", fieldID))(state)
					},
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_schema_field_configuration.test", "schema", schemaID))(state)
					},
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema_field_configuration.test", "required", "true"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema_field_configuration.test", "default_value", "test"),
				),
			},
			{
				Config: configForUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldSchemaFieldConfigurationExists("pagerduty_custom_field_schema_field_configuration.test"),
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_schema_field_configuration.test", "field", fieldID))(state)
					},
					func(state *terraform.State) error {
						return (resource.TestCheckResourceAttr(
							"pagerduty_custom_field_schema_field_configuration.test", "schema", schemaID))(state)
					},
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema_field_configuration.test", "required", "true"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema_field_configuration.test", "default_value", "updated"),
				),
			},
		},
	})
}

func TestAccPagerDutyCustomFieldConfiguration_Required_BadDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	config := fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "string"
}

resource "pagerduty_custom_field_schema" "input" {
  title = "%[2]s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                  = pagerduty_custom_field.input.id
  schema                 = pagerduty_custom_field_schema.input.id
  required               = true
  default_value          = "test"
  default_value_datatype = "garbage"
}

`, fieldName, schemaTitle)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Unknown datatype garbage"),
			},
		},
	})
}

func TestAccPagerDutyCustomFieldConfiguration_Required_Without_Default(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	config := fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "string"
}

resource "pagerduty_custom_field_schema" "input" {
  title = "%[2]s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                  = pagerduty_custom_field.input.id
  schema                 = pagerduty_custom_field_schema.input.id
  required               = true
  default_value          = "test"
}

`, fieldName, schemaTitle)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("required field without default value"),
			},
		},
	})
}

func TestAccPagerDutyCustomFieldConfiguration_Required_Without_DefaultDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	config := fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "string"
}

resource "pagerduty_custom_field_schema" "input" {
  title = "%[2]s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                  = pagerduty_custom_field.input.id
  schema                 = pagerduty_custom_field_schema.input.id
  required               = true
}

`, fieldName, schemaTitle)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("required field without default value"),
			},
		},
	})
}

func TestAccPagerDutyFieldConfiguration_Required_InvalidValue(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	config := fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "string"
}

resource "pagerduty_custom_field_schema" "input" {
  title = "%[2]s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                  = pagerduty_custom_field.input.id
  schema                 = pagerduty_custom_field_schema.input.id
  required               = true
  default_value          = "garbage"
  default_value_datatype = "integer"
}

`, fieldName, schemaTitle)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("invalid default value for datatype integer: garbage"),
			},
		},
	})
}

func TestAccPagerDutyCustomFieldConfiguration_Required_InvalidValue_MultiValue(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	config := fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "string"
}

resource "pagerduty_custom_field_schema" "input" {
  title = "%[2]s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field                     = pagerduty_custom_field.input.id
  schema                    = pagerduty_custom_field_schema.input.id
  required                  = true
  default_value             = "garbage"
  default_value_datatype    = "integer"
  default_value_multi_value = true
}

`, fieldName, schemaTitle)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("invalid default value for datatype integer \\(multi-value\\): garbage"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldSchemaFieldConfigurationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no field configuration ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		schemaID := rs.Primary.Attributes["schema"]
		found, _, err := client.CustomFieldSchemas.GetFieldConfiguration(schemaID, rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("field configuration not found: %v/%v - %v", schemaID, rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCreateTestPagerDutyCustomFieldForFieldConfiguration(fieldName string) *pagerduty.CustomField {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	field, _, err := client.CustomFields.Create(&pagerduty.CustomField{
		DataType:    pagerduty.CustomFieldDataTypeString,
		Name:        fieldName,
		DisplayName: fieldName,
	})
	if err != nil {
		panic("Unable to create test field to use as field configuration")
	}
	return field
}

func testAccDeleteTestPagerDutyCustomFieldForFieldConfiguration(fieldID string) error {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	_, err := client.CustomFields.Delete(fieldID)
	return err
}

func testAccCreateTestPagerDutyCustomFieldSchemaForFieldConfiguration(schemaTitle string) *pagerduty.CustomFieldSchema {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	field, _, err := client.CustomFieldSchemas.Create(&pagerduty.CustomFieldSchema{
		Title: schemaTitle,
	})
	if err != nil {
		panic("Unable to create test field schema to use as field configuration")
	}
	return field
}

func testAccDeleteTestPagerDutyCustomFieldSchemaForFieldConfiguration(schemaID string) error {
	config := Config{
		Token:     os.Getenv("PAGERDUTY_TOKEN"),
		UserToken: os.Getenv("PAGERDUTY_USER_TOKEN"),
	}
	client, _ := (&config).Client()

	_, err := client.CustomFieldSchemas.Delete(schemaID)
	return err
}

func testAccCheckPagerDutyCustomFieldConfigurationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_custom_field_schema_field_configuration" {
			continue
		}

		schemaID := r.Primary.Attributes["schema"]
		if _, _, err := client.CustomFieldSchemas.GetFieldConfiguration(schemaID, r.Primary.ID, nil); err == nil {
			return fmt.Errorf("field configuration still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyCustomFieldConfigurationConfigBasic(name string, schemaTitle string) string {
	return fmt.Sprintf(`
data "pagerduty_custom_field" "input" {
  name = "%s"
}

data "pagerduty_custom_field_schema" "input" {
  title = "%s"
}

resource "pagerduty_custom_field_schema_field_configuration" "test" {
  field  = data.pagerduty_custom_field.input.id
  schema = data.pagerduty_custom_field_schema.input.id
}

`, name, schemaTitle)
}
