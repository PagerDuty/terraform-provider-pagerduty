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
	resource.AddTestSweepers("pagerduty_tag", &resource.Sweeper{
		Name: "pagerduty_tag",
		F:    testSweepTag,
	})
}

func testSweepTag(_ string) error {
	client := testAccProvider.client
	ctx := context.Background()

	resp, err := client.ListTags(pagerduty.ListTagOptions{})
	if err != nil {
		return err
	}

	for _, tag := range resp.Tags {
		if strings.HasPrefix(tag.Label, "test") || strings.HasPrefix(tag.Label, "tf-") {
			log.Printf("Destroying tag %s (%s)", tag.Label, tag.ID)
			if err := client.DeleteTagWithContext(ctx, tag.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyTag_Basic(t *testing.T) {
	tagLabel := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagConfig(tagLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTagExists("pagerduty_tag.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_tag.foo", "label", tagLabel),
				),
			},
			// Validating that externally removed tags are detected and planed for
			// re-creation
			{
				Config: testAccCheckPagerDutyTagConfig(tagLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyTag("pagerduty_tag.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckPagerDutyTagDestroy(s *terraform.State) error {
	client := testAccProvider.client
	ctx := context.Background()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_tag" {
			continue
		}
		if _, err := client.GetTagWithContext(ctx, r.Primary.ID); err == nil {
			return fmt.Errorf("Tag still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.client
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag ID is set")
		}

		found, err := client.GetTagWithContext(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Tag not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccExternallyDestroyTag(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.client
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag ID is set")
		}

		if err := client.DeleteTagWithContext(ctx, rs.Primary.ID); err != nil {
			return err
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
