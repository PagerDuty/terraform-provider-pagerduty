package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyExtension_import(t *testing.T) {
	extensionName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	url := "https://example.com/receive_a_pagerduty_webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyExtensionConfig(name, extensionName, url, "false", "any"),
			},
			{
				ResourceName:            "pagerduty_extension.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

func TestAccPagerDutyExtension_import_WithoutEndpointURL(t *testing.T) {
	extension_name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	url := "https://example.com/receive_a_pagerduty_webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyExtensionConfig(name, extension_name, url, "false", "any"),
			},
			{
				Config: testAccCheckPagerDutyExtensionConfig_NoEndpointURL(name, extension_name, "false", "any"),
			},
			{
				ResourceName:            "pagerduty_extension.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}
