package pagerduty

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestAccPagerDutyTagAssignment_User(t *testing.T) {
	tagLabel := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTagAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagAssignmentConfig(tagLabel, username, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTagAssignmentExists("pagerduty_tag_assignment.foo", "users"),
					resource.TestCheckResourceAttr(
						"pagerduty_tag.foo", "label", tagLabel),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "name", username),
					resource.TestCheckResourceAttr(
						"pagerduty_user.foo", "email", email),
				),
			},
		},
	})
}
func TestAccPagerDutyTagAssignment_Team(t *testing.T) {
	tagLabel := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTagAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagAssignmentTeamConfig(tagLabel, team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTagAssignmentExists("pagerduty_tag_assignment.foo", "teams"),
					resource.TestCheckResourceAttr(
						"pagerduty_tag.foo", "label", tagLabel),
					resource.TestCheckResourceAttr(
						"pagerduty_team.foo", "name", team),
				),
			},
		},
	})
}
func TestAccPagerDutyTagAssignment_EP(t *testing.T) {
	tagLabel := fmt.Sprintf("tf-%s", acctest.RandString(5))
	ep := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTagAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagAssignmentEPConfig(tagLabel, username, email, ep),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyTagAssignmentExists("pagerduty_tag_assignment.foo", "teams"),
					resource.TestCheckResourceAttr(
						"pagerduty_tag.foo", "label", tagLabel),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", ep),
				),
			},
		},
	})
}

func testAccCheckPagerDutyTagAssignmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_tag_assignment" {
			continue
		}
		ids := strings.Split(r.Primary.ID, ".")

		entityID, tagID := ids[0], ids[1]
		entityType := "users"

		response, _, err := client.Tags.ListTagsForEntity(entityType, entityID)
		if err != nil {
			// if there are no tags for the entity that's okay
			return nil
		}
		// find tag the test created
		for _, tag := range response.Tags {
			if tag.ID == tagID {
				return fmt.Errorf("Tag %s still exists and is connected to %s ID %s", tag.ID, entityType, entityID)
			}
		}
	}
	return nil
}

func testAccCheckPagerDutyTagAssignmentExists(n, entityType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag Assignment ID is set")
		}
		ids := strings.Split(rs.Primary.ID, ".")

		entityID, tagID := ids[0], ids[1]

		client := testAccProvider.Meta().(*pagerduty.Client)
		response, _, err := client.Tags.ListTagsForEntity(entityType, entityID)
		if err != nil {
			return err
		}
		// find tag the test created
		var isFound bool = false
		for _, tag := range response.Tags {
			if tag.ID == tagID {
				isFound = true
				break
			}
		}
		if !isFound {
			return fmt.Errorf("Tag %s is no longer connected to %s ID %s", tagID, entityType, entityID)
		}
		return nil
	}
}

func testAccCheckPagerDutyTagAssignmentConfig(tagLabel, username, email string) string {
	return fmt.Sprintf(`
resource "pagerduty_tag" "foo" {
	label = "%s"
}
resource "pagerduty_user" "foo" {
	name = "%s"
	email = "%s"
}
resource "pagerduty_tag_assignment" "foo" {
	entity_type = "users"
	entity_id = pagerduty_user.foo.id
	tag_id = pagerduty_tag.foo.id
}
`, tagLabel, username, email)
}

func testAccCheckPagerDutyTagAssignmentTeamConfig(tagLabel, team string) string {
	return fmt.Sprintf(`
resource "pagerduty_tag" "foo" {
	label = "%s"
}
resource "pagerduty_team" "foo" {
	name = "%s"
}
resource "pagerduty_tag_assignment" "foo" {
	entity_type = "teams"
	entity_id = pagerduty_team.foo.id
	tag_id = pagerduty_tag.foo.id
}
`, tagLabel, team)
}

func testAccCheckPagerDutyTagAssignmentEPConfig(tagLabel, username, email, ep string) string {
	return fmt.Sprintf(`
resource "pagerduty_tag" "foo" {
	label = "%s"
}

resource "pagerduty_user" "foo" {
	name        = "%s"
	email       = "%s"
}
  
resource "pagerduty_escalation_policy" "foo" {
	name        = "%s"
	description = "foo"
	num_loops   = 1
  
	rule {
	  escalation_delay_in_minutes = 10
  
	  target {
		type = "user_reference"
		id   = pagerduty_user.foo.id
	  }
	}
  }
resource "pagerduty_tag_assignment" "foo" {
	entity_type = "teams"
	entity_id = pagerduty_escalation_policy.foo.id
	tag_id = pagerduty_tag.foo.id
}
`, tagLabel, username, email, ep)
}
