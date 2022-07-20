package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyUsers_Basic(t *testing.T) {
	teamname1 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	teamname2 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	username1 := fmt.Sprintf("tf-user1-%s", acctest.RandString(5))
	email1 := fmt.Sprintf("%s@foo.test", username1)
	username2 := fmt.Sprintf("tf-user2-%s", acctest.RandString(5))
	email2 := fmt.Sprintf("%s@foo.test", username2)
	username3 := fmt.Sprintf("tf-user3-%s", acctest.RandString(5))
	email3 := fmt.Sprintf("%s@foo.test", username3)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyUsersConfig(teamname1, teamname2, username1, email1, username2, email2, username3, email3),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyUsersExists("data.pagerduty_users.test_all_users"),
					testAccDataSourcePagerDutyUsersExists("data.pagerduty_users.test_by_1_team"),
					testAccDataSourcePagerDutyUsersExists("data.pagerduty_users.test_by_2_team"),
					resource.TestCheckResourceAttrSet(
						"data.pagerduty_users.test_all_users", "users.#"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.pagerduty_users.test_all_users",
						"users.*",
						map[string]string{
							"name": username1,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.pagerduty_users.test_all_users",
						"users.*",
						map[string]string{
							"name": username2,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.pagerduty_users.test_all_users",
						"users.*",
						map[string]string{
							"name": username3,
						}),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_1_team", "users.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.pagerduty_users.test_by_1_team", "users.0.id"),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_1_team", "users.0.name", username2),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_1_team", "users.0.email", email2),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_2_team", "users.#", "2"),
					resource.TestCheckResourceAttrSet(
						"data.pagerduty_users.test_by_2_team", "users.0.id"),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_2_team", "users.0.name", username2),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_2_team", "users.0.email", email2),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_2_team", "users.1.name", username3),
					resource.TestCheckResourceAttr(
						"data.pagerduty_users.test_by_2_team", "users.1.email", email3),
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

func testAccDataSourcePagerDutyUsersConfig(teamname1, teamname2, username1, email1, username2, email2, username3, email3 string) string {
	return fmt.Sprintf(`
    resource "pagerduty_team" "test1" {
        name        = "%s"
    }
    resource "pagerduty_team" "test2" {
        name        = "%s"
    }

    resource "pagerduty_user" "test_wo_team" {
      name = "%s"
      email = "%s"
    }
    resource "pagerduty_user" "test_w_team1" {
      name = "%s"
      email = "%s"
    }
    resource "pagerduty_user" "test_w_team2" {
      name = "%s"
      email = "%s"
    }

    resource "pagerduty_team_membership" "test1" {
      team_id = pagerduty_team.test1.id
      user_id = pagerduty_user.test_w_team1.id
    }
    resource "pagerduty_team_membership" "test2" {
      depends_on = [pagerduty_team_membership.test1]
      team_id = pagerduty_team.test2.id
      user_id = pagerduty_user.test_w_team2.id
    }

    data "pagerduty_users" "test_all_users" {
      depends_on = [pagerduty_user.test_w_team1, pagerduty_user.test_wo_team]
    }

    data "pagerduty_users" "test_by_1_team" {
      depends_on = [pagerduty_team_membership.test1]
      team_ids = [pagerduty_team.test1.id]
    }
    data "pagerduty_users" "test_by_2_team" {
      depends_on = [pagerduty_team_membership.test2]
      team_ids = [pagerduty_team.test1.id, pagerduty_team.test2.id]
    }
`, teamname1, teamname2, username1, email1, username2, email2, username3, email3)
}
