package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourcePagerDutyUsers_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)
	team := "12345"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyUsersConfig(username, email, team),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyUsers("pagerduty_user.test", "data.pagerduty_users.by_teamid"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyUsers(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		fmt.Println("Should be here")
		fmt.Println("\nsrc: " + src)
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a user ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the user %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyUsersConfig(username, email, team string) string {
	fmt.Println("Made it into the config")
	setup := fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name = "%s"
  email = "%s"
	teams = ["%s"]
}

data "pagerduty_users" "by_teamid" {
	team = "%s"
}
`, username, email, team, team)
	fmt.Println("Built the config")
	return setup
}
