package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_addon", &resource.Sweeper{
		Name: "pagerduty_addon",
		F:    testSweepAddon,
	})
}

func testSweepAddon(_ string) error {
	ctx := context.Background()

	resp, err := testAccProvider.client.ListAddonsWithContext(ctx, pagerduty.ListAddonOptions{})
	if err != nil {
		return err
	}

	for _, addon := range resp.Addons {
		if strings.HasPrefix(addon.Name, "test") || strings.HasPrefix(addon.Name, "tf-") {
			log.Printf("Destroying add-on %s (%s)", addon.Name, addon.ID)
			if err := testAccProvider.client.DeleteAddonWithContext(ctx, addon.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyAddon_Basic(t *testing.T) {
	addon := fmt.Sprintf("tf-%s", acctest.RandString(5))
	addonUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyAddonDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyAddonConfig(addon),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAddonExists("pagerduty_addon.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_addon.foo", "name", addon),
					resource.TestCheckResourceAttr(
						"pagerduty_addon.foo", "src", "https://intranet.foo.test/status"),
				),
			},
			{
				Config: testAccCheckPagerDutyAddonConfigUpdated(addonUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyAddonExists("pagerduty_addon.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_addon.foo", "name", addonUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_addon.foo", "src", "https://intranet.bar.com/status"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyAddonDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_addon" {
			continue
		}

		ctx := context.Background()

		if _, err := testAccProvider.client.GetAddonWithContext(ctx, r.Primary.ID); err == nil {
			return fmt.Errorf("Add-on still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyAddonExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No add-on ID is set")
		}

		found, err := testAccProvider.client.GetAddonWithContext(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Add-on not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyAddonConfig(addon string) string {
	return fmt.Sprintf(`
resource "pagerduty_addon" "foo" {
  name = "%s"
  src  = "https://intranet.foo.test/status"
}
`, addon)
}

func testAccCheckPagerDutyAddonConfigUpdated(addon string) string {
	return fmt.Sprintf(`
resource "pagerduty_addon" "foo" {
  name = "%s"
  src  = "https://intranet.bar.com/status"
}
`, addon)
}
