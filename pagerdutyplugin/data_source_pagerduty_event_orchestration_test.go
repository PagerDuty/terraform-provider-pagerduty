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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestration("pagerduty_event_orchestration.test", "data.pagerduty_event_orchestration.test"),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyEventOrchestration_WithIntegrations(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationWithIntegrationsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						a := s.RootModule().Resources["data.pagerduty_event_orchestration.test"].Primary.Attributes
						t.Log("[CG]", a)
						return nil
					},
					resource.TestCheckResourceAttr(
						"data.pagerduty_event_orchestration.test", "name", name),
					resource.TestCheckResourceAttr(
						"data.pagerduty_event_orchestration.test", "integration.0.label", "foo"),
					resource.TestCheckResourceAttr(
						"data.pagerduty_event_orchestration.test", "integration.0.parameters.0.type", "foo"),
					resource.TestCheckResourceAttr(
						"data.pagerduty_event_orchestration.test", "integration.1.label", "bar"),
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

func testAccDataSourcePagerDutyEventOrchestrationWithIntegrationsConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_event_orchestration" "test" {
  name = "%s"
}

resource "pagerduty_event_orchestration_integration" "foo" {
  label               = "foo"
  event_orchestration = pagerduty_event_orchestration.test.id
}

resource "pagerduty_event_orchestration_integration" "bar" {
  label               = "bar"
  event_orchestration = pagerduty_event_orchestration.test.id
}

data "pagerduty_event_orchestration" "test" {
  name = pagerduty_event_orchestration.test.name

  depends_on = [
    pagerduty_event_orchestration_integration.foo,
    pagerduty_event_orchestration_integration.bar
  ]
}
`, name)
}
