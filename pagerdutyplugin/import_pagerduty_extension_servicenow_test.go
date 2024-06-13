package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyExtensionServiceNow_import(t *testing.T) {
	extensionName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	url := "https://example.com/receive_a_pagerduty_webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyExtensionServiceNowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyExtensionServiceNowConfig(name, extensionName, url, "false", "any"),
			},
			{
				ResourceName:            "pagerduty_extension_servicenow.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}
