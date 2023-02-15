package pagerduty

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_custom_fields", &resource.Sweeper{
		Name:         "pagerduty_custom_fields",
		Dependencies: []string{"pagerduty_custom_field_schemas"},
		F:            testSweepCustomField,
	})
}

func testSweepCustomField(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.CustomFields.List(nil)
	if err != nil {
		return err
	}

	for _, customField := range resp.Fields {
		if strings.HasPrefix(customField.Name, "tf_") {
			log.Printf("Destroying field %s (%s)", customField.Name, customField.ID)
			if _, err := client.CustomFields.Delete(customField.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyCustomFields_Basic(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	description1 := acctest.RandString(10)
	description2 := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyCustomFieldConfig(fieldName, description1, "string"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldExists("pagerduty_custom_field.input"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "name", fieldName),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "description", description1),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "datatype", "string"),
				),
			},
			{
				Config: testAccCheckPagerDutyCustomFieldConfig(fieldName, description2, "string"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldExists("pagerduty_custom_field.input"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "description", description2),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "datatype", "string"),
				),
			},
		},
	})
}

func TestAccPagerDutyCustomField_BasicWithDescription(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))
	description := acctest.RandString(30)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyCustomFieldConfigWithDescription(fieldName, description, "string"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldExists("pagerduty_custom_field.input"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "name", fieldName),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "datatype", "string"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field.input", "description", description),
				),
			},
		},
	})
}

func TestAccPagerDutyCustomFields_UnknownDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldConfig(fieldName, "", "garbage"),
				ExpectError: regexp.MustCompile("Unknown datatype garbage"),
			},
		},
	})
}

func TestAccPagerDutyCustomFields_IllegalDataType(t *testing.T) {
	fieldName := fmt.Sprintf("tf_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyCustomFieldDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPagerDutyCustomFieldConfig(fieldName, "", pagerduty.CustomFieldDataTypeFieldOption.String()),
				ExpectError: regexp.MustCompile("Datatype field_option is not allowed on fields"),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldConfig(name, description, datatype string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  description = "%[2]s" 
  datatype = "%[3]s"
}
`, name, description, datatype)
}

func testAccCheckPagerDutyCustomFieldConfigWithDescription(name, description, datatype string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field" "input" {
  name = "%[1]s"
  display_name = "%[1]s"
  datatype = "%[2]s"
  description = "%[3]s"
}
`, name, datatype, description)
}

func testAccCheckPagerDutyCustomFieldDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_custom_field" {
			continue
		}

		if _, _, err := client.CustomFields.Get(r.Primary.ID, nil); err == nil {
			return fmt.Errorf("field still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyCustomFieldExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no field ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.CustomFields.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("field not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccPreCheckCustomFieldTests(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_CUSTOM_FIELDS"); v == "" {
		t.Skip("PAGERDUTY_ACC_CUSTOM_FIELDS not set. Skipping Custom Field-related test")
	}
}
