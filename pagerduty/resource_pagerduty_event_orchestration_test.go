package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_event_orchestration", &resource.Sweeper{
		Name: "pagerduty_event_orchestration",
		F:    testSweepEventOrchestration,
	})
}

func testSweepEventOrchestration(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.EventOrchestrations.List()
	if err != nil {
		return err
	}

	for _, orchestration := range resp.Orchestrations {
		if strings.HasPrefix(orchestration.Name, "tf-name-") {
			log.Printf("Destroying Event Orchestration %s (%s)", orchestration.Name, orchestration.ID)
			if _, err := client.EventOrchestrations.Delete(orchestration.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyEventOrchestration_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-description-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-team-%s", acctest.RandString(5))
	nameUpdated := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	descriptionUpdated := fmt.Sprintf("tf-description-%s", acctest.RandString(5))
	teamNameUpdated := fmt.Sprintf("tf-team-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationConfigNameOnly(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationExists("pagerduty_event_orchestration.nameonly"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.nameonly", "name", name),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationConfig(name, description, teamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationExists("pagerduty_event_orchestration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.foo", "name", name,
					),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.foo", "description", description,
					),
				),
			},
			{
				Config: testAccCheckPagerDutyEventOrchestrationConfigUpdated(nameUpdated, descriptionUpdated, teamNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationExists("pagerduty_event_orchestration.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.foo", "name", nameUpdated,
					),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.foo", "description", descriptionUpdated,
					),
				),
			},
		},
	})
}

func testAccCheckPagerDutyEventOrchestrationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_event_orchestration" {
			continue
		}
		if _, _, err := client.EventOrchestrations.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Event Orchestration still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyEventOrchestrationExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		orch, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if orch.Primary.ID == "" {
			return fmt.Errorf("No Event Orchestration ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.EventOrchestrations.Get(orch.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != orch.Primary.ID {
			return fmt.Errorf("Event Orchrestration not found: %v - %v", orch.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyEventOrchestrationConfig(name, description, team string) string {
	return fmt.Sprintf(`

resource "pagerduty_team" "foo" {
	name = "%s"
}
resource "pagerduty_event_orchestration" "foo" {
	name = "%s"
	description = "%s"
	team {
		id = pagerduty_team.foo.id
	}
}
`, team, name, description)
}

func testAccCheckPagerDutyEventOrchestrationConfigNameOnly(n string) string {
	return fmt.Sprintf(`

resource "pagerduty_event_orchestration" "nameonly" {
	name = "%s"
}
`, n)
}

func testAccCheckPagerDutyEventOrchestrationConfigUpdated(name, description, team string) string {
	return fmt.Sprintf(`

resource "pagerduty_team" "foo" {
	name = "%s"
}
resource "pagerduty_event_orchestration" "foo" {
	name = "%s"
	description = "%s"
	team {
		id = pagerduty_team.foo.id
	}
}
`, team, name, description)
}
