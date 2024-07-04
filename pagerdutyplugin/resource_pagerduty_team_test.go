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
	resource.AddTestSweepers("pagerduty_team", &resource.Sweeper{
		Name: "pagerduty_team",
		F:    testSweepTeam,
		Dependencies: []string{
			"pagerduty_escalation_policy",
			"pagerduty_service",
			"pagerduty_schedule",
		},
	})
}

func testSweepTeam(region string) error {
	ctx := context.Background()
	response, err := testAccProvider.client.ListTeamsWithContext(ctx, pagerduty.ListTeamOptions{})
	if err != nil {
		return err
	}

	for _, team := range response.Teams {
		if strings.HasPrefix(team.Name, "test") || strings.HasPrefix(team.Name, "tf-") {
			log.Printf("Destroying team %s (%s)", team.Name, team.ID)
			if err := testAccProvider.client.DeleteTeamWithContext(ctx, team.ID); err != nil {
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamConfig(team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamExists("pagerduty_team.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", team),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "description", "foo"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_team.foo", "html_url"),
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
			// Validating that externally removed teams are detected and planed for
			// re-creation
			{
				Config: testAccCheckPagerDutyTeamConfigUpdated(teamUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroyTeam("pagerduty_team.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccCheckPagerDutyTeamConfigUpdated(teamUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"pagerduty_team.foo", "id"),
				),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPagerDutyTeam_DefaultRole(t *testing.T) {
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	defaultRole := "manager"
	defaultRoleUpdated := "none"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamDefaultRoleConfig(team, defaultRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamExists("pagerduty_team.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", team),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "default_role", defaultRole),
					resource.TestCheckResourceAttrSet(
						"pagerduty_team.foo", "html_url"),
				),
			},
			{
				Config: testAccCheckPagerDutyTeamDefaultRoleConfig(team, defaultRoleUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "default_role", defaultRoleUpdated),
				),
			},
		},
	})
}

func TestAccPagerDutyTeam_Parent(t *testing.T) {
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	parent := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTeamWithParentConfig(team, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTeamExists("pagerduty_team.foo"),
					testAccCheckPagerDutyTeamExists("pagerduty_team.parent"),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", team),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "description", "foo"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_team.foo", "html_url"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_team.foo", "parent"),
					resource.TestCheckResourceAttr(
						"pagerduty_team.parent", "name", parent),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTeamDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_team" {
			continue
		}

		if _, err := testAccProvider.client.GetTeamWithContext(ctx, r.Primary.ID); err == nil {
			return fmt.Errorf("Team still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyTeamExists(_ string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			ctx := context.Background()
			if _, err := testAccProvider.client.GetTeamWithContext(ctx, r.Primary.ID); err != nil {
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

func testAccCheckPagerDutyTeamDefaultRoleConfig(team, defaultRole string) string {
	return fmt.Sprintf(`

resource "pagerduty_team" "foo" {
  name         = "%s"
  description  = "foo"
  default_role = "%s"
}
`, team, defaultRole)
}

func testAccCheckPagerDutyTeamWithParentConfig(team, parent string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "parent" {
	name        = "%s"
	description = "parent"
}	
resource "pagerduty_team" "foo" {
	name        = "%s"
	description = "foo"
	parent = pagerduty_team.parent.id
}`, parent, team)
}

func testAccExternallyDestroyTeam(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Team ID is set")
		}

		ctx := context.Background()
		err := testAccProvider.client.DeleteTeamWithContext(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}
