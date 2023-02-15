package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_custom_field_schemas", &resource.Sweeper{
		Name:         "pagerduty_custom_field_schemas",
		Dependencies: []string{"pagerduty_custom_field_schema_assignments"},
		F:            testSweepCustomFieldSchema,
	})
}

func testSweepCustomFieldSchema(region string) error {
	return testSweepForEachCustomFieldSchema(region, func(client *pagerduty.Client, fieldSchema *pagerduty.CustomFieldSchema) error {
		log.Printf("Destroying field schema %s. (%s)", fieldSchema.Title, fieldSchema.ID)
		if _, err := client.CustomFieldSchemas.Delete(fieldSchema.ID); err != nil {
			return err
		}
		return nil
	})
}

func testSweepForEachCustomFieldSchema(region string, handler func(client *pagerduty.Client, fieldSchema *pagerduty.CustomFieldSchema) error) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.CustomFieldSchemas.List(nil)
	if err != nil {
		return err
	}

	for _, fieldSchema := range resp.Schemas {
		if strings.HasPrefix(fieldSchema.Title, "tf-") {
			err = handler(client, fieldSchema)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func TestAccPagerDutyCustomFieldSchemas_Basic(t *testing.T) {
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))
	schemaTitleForUpdate := fmt.Sprintf("tf_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyCustomFieldSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyCustomFieldSchemaConfigBasic(schemaTitle),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldSchemaExists("pagerduty_custom_field_schema.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema.test", "title", schemaTitle),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema.test", "description", "some description"),
				),
			},
			{
				Config: testAccCheckPagerDutyCustomFieldSchemaConfigBasic(schemaTitleForUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldSchemaExists("pagerduty_custom_field_schema.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema.test", "title", schemaTitleForUpdate),
					resource.TestCheckResourceAttr(
						"pagerduty_custom_field_schema.test", "description", "some description"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldSchemaConfigBasic(title string) string {
	return fmt.Sprintf(`
resource "pagerduty_custom_field_schema" "test" {
  title = "%[1]s"
  description = "some description"
}
`, title)
}

func testAccCheckPagerDutyCustomFieldSchemaDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_custom_field_schema" {
			continue
		}

		if _, _, err := client.CustomFieldSchemas.Get(r.Primary.ID, nil); err == nil {
			return fmt.Errorf("custom field schema still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyCustomFieldSchemaExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no custom field schema ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.CustomFieldSchemas.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("custom field schema not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
