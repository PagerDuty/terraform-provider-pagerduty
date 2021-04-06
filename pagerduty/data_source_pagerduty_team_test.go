package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourcePagerDutyTeam_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := "team description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyTeamConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyTeam("pagerduty_team.test", "data.pagerduty_team.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyTeam(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a user ID from PagerDuty")
		}

		testAtts := []string{"id", "name", "description"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the team %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyTeamConfig(name, description string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "test" {
  name        = "%s"
  description = "%s"
}

data "pagerduty_team" "by_name" {
	name = pagerduty_team.test.name
}
`, name, description)
}
