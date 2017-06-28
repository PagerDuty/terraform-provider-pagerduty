package pagerduty

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func testSweepTeam(region string) error {
	client, err := sweeperClient()
	if err != nil {
		return err
	}

	resp, _, err := client.Teams.List(&pagerduty.ListTeamsOptions{})
	if err != nil {
		return err
	}

	for _, team := range resp.Teams {
		if strings.HasPrefix(team.Name, "test") || strings.HasPrefix(team.Name, "tf-") {
			if _, err := client.Teams.Delete(team.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyTeam_Basic(t *testing.T) {
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamConfig(team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamExists("pagerduty_team.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", team),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "description", "foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyTeamConfigUpdated(teamUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamExists("pagerduty_team.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", teamUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "description", "bar"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_team" {
			continue
		}

		if _, _, err := client.Teams.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Team still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyTeamExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*pagerduty.Client)
		for _, r := range s.RootModule().Resources {
			if _, _, err := client.Teams.Get(r.Primary.ID); err != nil {
				return fmt.Errorf("Received an error retrieving team %s ID: %s", err, r.Primary.ID)
			}
		}
		return nil
	}
}

func testAccCheckPagerDutyTeamConfig(team string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "foo"
}`, team)
}

func testAccCheckPagerDutyTeamConfigUpdated(team string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name        = "%s"
  description = "bar"
}`, team)
}
