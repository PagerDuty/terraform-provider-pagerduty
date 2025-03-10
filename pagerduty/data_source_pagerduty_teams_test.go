package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyTeams_Basic(t *testing.T) {
	teamname1 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	teamname2 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	teamname3 := fmt.Sprintf("tf-team-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyTeamsConfig(teamname1, teamname2, teamname3),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyTeamsExists("data.pagerduty_teams.test_all_teams"),
					testAccDataSourcePagerDutyTeamsExists("data.pagerduty_teams.test_by_1_team"),
					resource.TestCheckResourceAttrSet(
						"data.pagerduty_users.test_all_teams", "teams.#"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.pagerduty_users.test_all_teams",
						"teams.*",
						map[string]string{
							"name": teamname1,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.pagerduty_users.test_all_teams",
						"teams.*",
						map[string]string{
							"name": teamname2,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.pagerduty_users.test_all_teams",
						"teams.*",
						map[string]string{
							"name": teamname3,
						}),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyTeamsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a team ID from PagerDuty")
		}

		return nil
	}
}

func testAccDataSourcePagerDutyTeamsConfig(teamname1, teamname2, teamname3 string) string {
	return fmt.Sprintf(`
    resource "pagerduty_team" "test1" {
        name        = "%s"
    }
    resource "pagerduty_team" "test2" {
        name        = "%s"
    }
	resource "pagerduty_team" "test3" {
        name        = "%s"
    }

    data "pagerduty_teams" "test_by_1_team" {
      depends_on = [pagerduty_team.test1, pagerduty_team.test2, pagerduty_team.test3]
      query = "%s"
    }

    data "pagerduty_teams" "test_all_teams" {
      depends_on = [pagerduty_team.test1, pagerduty_team.test2, pagerduty_team.test3]
    }
`, teamname1, teamname2, teamname3, teamname1)
}
