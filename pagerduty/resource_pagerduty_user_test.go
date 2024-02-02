package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_user", &resource.Sweeper{
		Name: "pagerduty_user",
		F:    testSweepUser,
		Dependencies: []string{
			"pagerduty_team",
			"pagerduty_schedule",
			"pagerduty_escalation_policy",
		},
	})
}

func testSweepUser(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.Users.List(&pagerduty.ListUsersOptions{})
	if err != nil {
		return err
	}

	for _, user := range resp.Users {
		if strings.HasPrefix(user.Name, "test") || strings.HasPrefix(user.Name, "tf") {
			log.Printf("Destroying user %s (%s)", user.Name, user.ID)
			if _, err := client.Users.Delete(user.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyUser_Basic(t *testing.T) {
	username := fmt.Sprintf("tf %s", acctest.RandString(5))
	usernameSpaces := " " + strings.ReplaceAll(username, " ", "  ") + " "
	usernameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", strings.ReplaceAll(username, " ", ""))
	emailUpdated := fmt.Sprintf("%s@foo.test", usernameUpdated)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserConfig(usernameSpaces, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", strings.ToLower(email)),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "color", "green"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "role", "user"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "job_title", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "description", "foo"),
					resource.TestCheckResourceAttrSet(
						"pagerduty_user.foo", "html_url"),
				),
			},
			{
				Config: testAccCheckPagerDutyUserConfigUpdated(usernameUpdated, emailUpdated, "observer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", usernameUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", emailUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "color", "red"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "role", "observer"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "job_title", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "description", "bar"),
				),
			},
			{
				Config: testAccCheckPagerDutyUserConfigUpdated(usernameUpdated, emailUpdated, "read_only_user"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "role", "read_only_user"),
				),
			},
			{
				Config: testAccCheckPagerDutyUserConfigUpdated(usernameUpdated, emailUpdated, "restricted_access"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "role", "restricted_access"),
				),
			},
		},
	})
}

func TestAccPagerDutyUserWithTeams_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	team1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team2 := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserWithTeamsConfig(team1, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", email),
				),
			},
			{
				Config: testAccCheckPagerDutyUserWithTeamsConfigUpdated(team1, team2, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", email),
				),
			},
			{
				Config: testAccCheckPagerDutyUserWithNoTeamsConfig(team1, team2, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExists("pagerduty_user.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", email),
				),
			},
		},
	})
}

func TestAccPagerDutyUserWithLicenses_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	licensesName := "test"

	// Licenses are required, but an account may only have 1 license. Therefore,
	// this index value `i` is set to 0 so that any integration tests on accounts
	// with only 1 license will still pass. Since an account is not guaranteed to
	// have multiple licenses, there is no test for changing licenses.
	i := "0"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserWithLicensesConfig(username, email, licensesName, i),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyUserExistsWithLicense("pagerduty_user.foo", fmt.Sprintf("data.pagerduty_licenses.%s", licensesName), i),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", email),
				),
			},
		},
	})
}

func testAccCheckPagerDutyUserDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_user" {
			continue
		}

		if _, _, err := client.Users.Get(r.Primary.ID, &pagerduty.GetUserOptions{}); err == nil {
			return fmt.Errorf("User still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No user ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.Users.Get(rs.Primary.ID, &pagerduty.GetUserOptions{})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("User not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

// This tests that the assigned license from the plan using index `i` matches
// the fetched users assigned license.
func testAccCheckPagerDutyUserExistsWithLicense(userResource, licenseData, i string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[userResource]
		if !ok {
			return fmt.Errorf("Not found: %s", userResource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No user ID is set")
		}

		dataR, ok := s.RootModule().Resources[licenseData]
		if !ok {
			return fmt.Errorf("Not found: %s", licenseData)
		}
		dataA := dataR.Primary.Attributes

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, err := client.Users.GetWithLicense(rs.Primary.ID, &pagerduty.GetUserOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("User not found: %v - %v", rs.Primary.ID, found)
		}

		licenseAttr := fmt.Sprintf("licenses.%s.id", i)
		licenseID, ok := dataA[licenseAttr]
		if !ok {
			return fmt.Errorf("Could not find %v in data.pagerduty_licenses", licenseAttr)
		}
		if licenseID != found.License.ID {
			return fmt.Errorf("User's assigned license %s does not match the configured license %s", found.License.ID, licenseID)
		}

		return nil
	}
}

func testAccCheckPagerDutyUserConfig(username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "green"
  role        = "user"
  job_title   = "foo"
  description = "foo"
  time_zone   = "Europe/Berlin"
}`, username, email)
}

func testAccCheckPagerDutyUserConfigUpdated(username, email, role string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
  color       = "red"
  role        = "%s"
  job_title   = "bar"
  description = "bar"
  time_zone   = "Europe/Dublin"
}`, username, email, role)
}

func testAccCheckPagerDutyUserWithTeamsConfig(team, username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name = "%s"
}

resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
}
`, team, username, email)
}

// The license is dynamically assigned based on the fetched licenses using
// the provided index `i`. The role is dynamically assigned based on the
// fetched licenses valid_roles. However, "owner" role is unique and so it is
// excluded with `local.invalid_roles`.
func testAccCheckPagerDutyUserWithLicensesConfig(username, email, licensesName, i string) string {
	return fmt.Sprintf(`
locals {
	invalid_roles = ["owner"]
}

data "pagerduty_licenses" "%s" {}

resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
	license = data.pagerduty_licenses.test.licenses[%s].id
	role = tolist(setsubtract(data.pagerduty_licenses.test.licenses[%s].valid_roles, local.invalid_roles))[0]
}
`, licensesName, username, email, i, i)
}

func testAccCheckPagerDutyUserWithTeamsConfigUpdated(team1, team2, username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name = "%s"
}

resource "pagerduty_team" "bar" {
  name = "%s"
}

resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team_membership" "foo" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.foo.id
}

resource "pagerduty_team_membership" "bar" {
  user_id = pagerduty_user.foo.id
  team_id = pagerduty_team.bar.id
}
`, team1, team2, username, email)
}

func testAccCheckPagerDutyUserWithNoTeamsConfig(team1, team2, username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_team" "foo" {
  name = "%s"
}

resource "pagerduty_team" "bar" {
  name = "%s"
}

resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}
`, team1, team2, username, email)
}
