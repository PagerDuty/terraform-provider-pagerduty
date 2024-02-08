package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyEventOrchestration_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestration("pagerduty_event_orchestration.test", "data.pagerduty_event_orchestration.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyEventOrchestration(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get an Event Orchestration ID from PagerDuty")
		}

		testAtts := []string{"id", "name", "integration"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the Event Orchestration %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEventOrchestrationConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_event_orchestration" "test" {
  name                    = "%s"
}

data "pagerduty_event_orchestration" "by_name" {
  name = pagerduty_event_orchestration.test.name
}
`, name)
}
