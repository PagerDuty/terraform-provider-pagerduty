package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("pagerduty_orchestration", &resource.Sweeper{
		Name: "pagerduty_orchestration",
		F:    testSweepOrchestration,
	})
}

func testSweepOrchestration(region string) error {
	// TODO: delete all orchestrations created by the tests
	return nil
}

func TestAccPagerDutyOrchestration_Basic(t *testing.T) {
	orchestration := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyOrchestrationConfigNameOnly(orchestration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyOrchestrationExists("pagerduty_orchestration.nameonly"),
					resource.TestCheckResourceAttr(
						"pagerduty_orchestration.nameonly", "name", orchestration),
				),
			},
		},
	})
}

func testAccCheckPagerDutyOrchestrationDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_orchestration" {
			continue
		}
		if _, _, err := client.Orchestrations.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("Orchestration still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyOrchestrationExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		orch, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}
		if orch.Primary.ID == "" {
			return fmt.Errorf("No Orchestration ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		found, _, err := client.Orchestrations.Get(orch.Primary.ID)
		if err != nil {
			return err
		}
		if found.ID != orch.Primary.ID {
			return fmt.Errorf("Orchrestration not found: %v - %v", orch.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyOrchestrationConfigNameOnly(n string) string {
	return fmt.Sprintf(`

resource "pagerduty_orchestration" "nameonly" {
	name = "%s"
}
`, n)
}