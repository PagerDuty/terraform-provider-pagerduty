package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyTag_Basic(t *testing.T) {
	tag := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyTagConfig(tag),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyTag("pagerduty_tag.test", "data.pagerduty_tag.by_label"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyTag(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a tag ID from PagerDuty")
		}

		testAtts := []string{"id", "label"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the tag %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyTagConfig(tag string) string {
	return fmt.Sprintf(`
resource "pagerduty_tag" "test" {
    label = "%s"
}

data "pagerduty_tag" "by_label" {
    label = pagerduty_tag.test.label
}
`, tag)
}
