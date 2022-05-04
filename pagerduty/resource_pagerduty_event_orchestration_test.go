package pagerduty

import (
	"fmt"
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
	// TODO: delete all orchestrations created by the tests
	return nil
}

func TestAccPagerDutyEventOrchestration_Basic(t *testing.T) {
	orchestration := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyOrchestrationConfigNameOnly(orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyEventOrchestrationExists("pagerduty_event_orchestration.nameonly"),
					resource.TestCheckResourceAttr(
						"pagerduty_event_orchestration.nameonly", "name", orchestration),
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

func testAccCheckPagerDutyOrchestrationConfigNameOnly(n string) string {
	return fmt.Sprintf(`

resource "pagerduty_event_orchestration" "nameonly" {
	name = "%s"
}
`, n)
}
