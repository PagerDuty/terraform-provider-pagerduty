package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_custom_field_schema_assignments", &resource.Sweeper{
		Name: "pagerduty_custom_field_schema_assignments",
		F:    testSweepCustomFieldSchemaAssignment,
	})
}

func testSweepCustomFieldSchemaAssignment(region string) error {
	return testSweepForEachCustomFieldSchema(region, func(client *pagerduty.Client, fieldSchema *pagerduty.CustomFieldSchema) error {
		resp, _, err := client.CustomFieldSchemaAssignments.ListForSchema(fieldSchema.ID, nil)
		if err != nil {
			return err
		}
		for _, a := range resp.SchemaAssignments {
			_, err = client.CustomFieldSchemaAssignments.Delete(a.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func TestAccPagerDutyCustomFieldSchemaAssignment(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	schemaTitle := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckCustomFieldTests(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyCustomFieldSchemaAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyCustomFieldSchemaAssignment(schemaTitle, username, email, escalationPolicy, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyCustomFieldSchemaAssignmentExists("pagerduty_custom_field_schema.test", "pagerduty_service.foo"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyCustomFieldSchemaAssignment(schemaTitle, username, email, escalationPolicy, service string) string {
	serviceConfig := testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service)
	schemaConfig := testAccCheckPagerDutyCustomFieldSchemaConfigBasic(schemaTitle)
	return fmt.Sprintf(`
%s

%s

resource "pagerduty_custom_field_schema_assignment" "test" {
  schema        = pagerduty_custom_field_schema.test.id
  service   = pagerduty_service.foo.id
}
`, serviceConfig, schemaConfig)
}

func testAccCheckPagerDutyCustomFieldSchemaAssignmentDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_custom_field_schema_assignment" {
			continue
		}
		schemaId := r.Primary.Attributes["schema"]
		serviceID := r.Primary.Attributes["service"]

		if as, _, err := client.CustomFieldSchemaAssignments.ListForService(serviceID, nil); err == nil && len(as.SchemaAssignments) != 0 &&
			testAccCustomFieldSchemaAssignmentsContainSchema(as.SchemaAssignments, schemaId) {
			return fmt.Errorf("field schema assignment still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyCustomFieldSchemaAssignmentExists(schemaName, serviceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		schemaResource, ok := s.RootModule().Resources[schemaName]
		if !ok {
			return fmt.Errorf("schema not found: %s", schemaName)
		}
		if schemaResource.Primary.ID == "" {
			return fmt.Errorf("no field schema ID is set")
		}

		serviceResource, ok := s.RootModule().Resources[serviceName]
		if !ok {
			return fmt.Errorf("service not found: %s", serviceName)
		}
		if serviceResource.Primary.ID == "" {
			return fmt.Errorf("no service ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		as, _, err := client.CustomFieldSchemaAssignments.ListForService(serviceResource.Primary.ID, nil)
		if err != nil {
			return err
		}

		if !testAccCustomFieldSchemaAssignmentsContainSchema(as.SchemaAssignments, schemaResource.Primary.ID) {
			return fmt.Errorf("field schema %v not found in associations of %v", schemaResource.Primary.ID, serviceResource.Primary.ID)
		}

		return nil
	}
}

func testAccCustomFieldSchemaAssignmentsContainSchema(as []*pagerduty.CustomFieldSchemaAssignment, id string) bool {
	for _, a := range as {
		if a.Schema.ID == id {
			return true
		}
	}
	return false
}
