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
	resource.AddTestSweepers("pagerduty_tag", &resource.Sweeper{
		Name: "pagerduty_tag",
		F:    testSweepTag,
	})
}

func testSweepTag(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.Tags.List(&pagerduty.ListTagsOptions{})
	if err != nil {
		return err
	}

	for _, tag := range resp.Tags {
		if strings.HasPrefix(tag.Label, "test") || strings.HasPrefix(tag.Label, "tf-") {
			log.Printf("Destroying tag %s (%s)", tag.Label, tag.ID)
			if _, err := client.Tags.Delete(tag.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyTag_Basic(t *testing.T) {
	tagLabel := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagConfig(tagLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTagExists("pagerduty_tag.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_tag.foo", "label", tagLabel),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTagDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_tag" {
			continue
		}
		if _, _, err := client.Tags.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Tag still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.Tags.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Tag not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyTagConfig(tagLabel string) string {
	return fmt.Sprintf(`
resource "pagerduty_tag" "foo" {
	label = "%s"
}
`, tagLabel)
}
