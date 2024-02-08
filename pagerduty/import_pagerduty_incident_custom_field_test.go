package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyIncidentCustomField_import(t *testing.T) {
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
				Config: testAccCheckPagerDutyIncidentCustomFieldConfig(fieldName, "test description", "string"),
			},
			{
				ResourceName:      "pagerduty_incident_custom_field.input",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyIncidentCustomField_import_no_description(t *testing.T) {
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
				Config: testAccCheckPagerDutyIncidentCustomFieldConfigNoDescription(fieldName, "string"),
			},
			{
				ResourceName:      "pagerduty_incident_custom_field.input",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
