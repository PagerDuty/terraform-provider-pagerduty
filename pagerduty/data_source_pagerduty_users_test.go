package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyUsers_Basic(t *testing.T) {
	teamname := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	username1 := fmt.Sprintf("tf-user-%s", acctest.RandString(5))
	email1 := fmt.Sprintf("%s@foo.test", username1)
	username2 := fmt.Sprintf("tf-user-%s", acctest.RandString(5))
	email2 := fmt.Sprintf("%s@foo.test", username2)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyUsersConfig(teamname, username1, email1, username2, email2),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyUsersExists("data.pagerduty_users.test_all_users"),
					testAccDataSourcePagerDutyUsersExists("data.pagerduty_users.test_by_team"),
					resource.TestCheckResourceAttrSet(
						"data.pagerduty_users.test_all_users", "users.#"),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_team", "users.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.pagerduty_users.test_by_team", "users.0.id"),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_team", "users.0.name", username2),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_team", "users.0.email", email2),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyUsersExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a user ID from PagerDuty")
		}

		return nil
	}
}

func testAccDataSourcePagerDutyUsersConfig(teamname, username1, email1, username2, email2 string) string {
	return fmt.Sprintf(`
    resource "pagerduty_team" "test" {
        name        = "%s"
    }
    resource "pagerduty_user" "test_wo_team" {
      name = "%s"
      email = "%s"
    }
    resource "pagerduty_user" "test_w_team" {
      name = "%s"
      email = "%s"
    }
    resource "pagerduty_team_membership" "test" {
      team_id = pagerduty_team.test.id
      user_id = pagerduty_user.test_w_team.id
    }

    data "pagerduty_users" "test_all_users" {
      depends_on = [pagerduty_user.test_w_team, pagerduty_user.test_wo_team]
    }

    data "pagerduty_users" "test_by_team" {
      depends_on = [pagerduty_team_membership.test]
      team_ids = [pagerduty_team.test.id]
    }
`, teamname, username1, email1, username2, email2)
}
