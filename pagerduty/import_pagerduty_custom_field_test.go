package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyField_import(t *testing.T) {
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
				Config: testAccCheckPagerDutyCustomFieldConfig(fieldName, "", "string"),
			},
			{
				ResourceName:      "pagerduty_custom_field.input",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
