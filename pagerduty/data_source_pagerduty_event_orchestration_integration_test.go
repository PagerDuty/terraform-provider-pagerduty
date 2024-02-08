package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyEventOrchestrationIntegration_Basic(t *testing.T) {
	on := fmt.Sprintf("tf-orchestration-%s", acctest.RandString(5))
	lbl := fmt.Sprintf("tf-integration-%s", acctest.RandString(5))
	irn := "pagerduty_event_orchestration_integration.orch_int"
	n := "data.pagerduty_event_orchestration_integration.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// find by id
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationIntegrationByIdConfig(on, lbl),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationIntegration(irn, n),
				),
			},
			// find by id, ignore label
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationIntegrationByIdLabelConfig(on, lbl),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationIntegration(irn, n),
				),
			},
			// find by label
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationIntegrationByLabelConfig(on, lbl),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrationIntegration(irn, n),
				),
			},
			// id and label are both not set
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationIntegrationIdLabelNullConfig(on, lbl),
				ExpectError: regexp.MustCompile("Invalid Event Orchestration Integration data source configuration: ID and label cannot both be null"),
			},
			// bad event_orchestration
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationIntegrationBadOrchConfig(on, lbl),
				ExpectError: regexp.MustCompile("Unable to find an Integration with ID '(.+)' on PagerDuty Event Orchestration 'bad-orchestration-id'"),
			},
			// bad id
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationIntegrationBadIdConfig(on, lbl),
				ExpectError: regexp.MustCompile("Unable to find an Integration with ID 'bad-integration-id' on PagerDuty Event Orchestration '(.+)'"),
			},
			// bad label
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationIntegrationBadLabelConfig(on, lbl),
				ExpectError: regexp.MustCompile("Unable to find an Integration on Event Orchestration '(.+)' with label 'bad-integration-label'"),
			},
			// ambiguous label
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationIntegrationAmbiguousLabelConfig(on, lbl),
				ExpectError: regexp.MustCompile("Ambiguous Integration label: '" + lbl + "'. Found 2 Integrations with this label on Event Orchestration '(.+)'. Please use the Integration ID instead or make Integration labels unique within Event Orchestration."),
			},
		},
	})
}

func testAccDataSourcePagerDutyEventOrchestrationIntegration(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected the Event Orchestration Integration ID to be set")
		}

		testAtts := []string{"id", "label", "parameters.0.routing_key", "parameters.0.type"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the Event Orchestration Integration %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationBaseConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationByIdConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			id = pagerduty_event_orchestration_integration.orch_int.id
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationByIdLabelConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			id = pagerduty_event_orchestration_integration.orch_int.id
			label = "No such label"
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationByLabelConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%[1]s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%[2]s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%[2]s"
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationIdLabelNullConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationBadOrchConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = "bad-orchestration-id"
			id = pagerduty_event_orchestration_integration.orch_int.id
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationBadIdConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			id = "bad-integration-id"
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationBadLabelConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "bad-integration-label"
		}
		`, on, lbl)
}

func testAccDataSourcePagerDutyEventOrchestrationIntegrationAmbiguousLabelConfig(on, lbl string) string {
	return fmt.Sprintf(`
		resource "pagerduty_event_orchestration" "orch" {
			name = "%[1]s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%[2]s"
		}

		resource "pagerduty_event_orchestration_integration" "orch_int_2" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%[2]s"
		}

		data "pagerduty_event_orchestration_integration" "by_id" {
			event_orchestration = pagerduty_event_orchestration.orch.id
			label = "%[2]s"
		}
		`, on, lbl)
}
