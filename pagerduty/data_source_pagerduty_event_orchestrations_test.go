package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePagerDutyEventOrchestrations_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrations("pagerduty_event_orchestration.test", "data.pagerduty_event_orchestrations.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyEventOrchestrations(src, n string) resource.TestCheckFunc {
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
			sub_att := fmt.Sprintf("event_orchestrations.0.%s", att)
			if a[sub_att] != srcA[att] {
				return fmt.Errorf("Expected the Event Orchestration %s to be: %s, but got: %s", att, srcA[att], a[sub_att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEventOrchestrationsConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_event_orchestration" "test" {
  name                    = "%s"
}

data "pagerduty_event_orchestrations" "by_name" {
  search = pagerduty_event_orchestration.test.name
}
`, name)
}
