package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyIncidentWorkflow(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	dataSourceName := fmt.Sprintf("data.pagerduty_incident_workflow.%s", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyIncidentWorkflowConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyIncidentWorkflowConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "input" {
  name = "%[1]s"
}

data "pagerduty_incident_workflow" "%[1]s" {
  depends_on = [
    pagerduty_incident_workflow.input
  ]
  name = "%[1]s"
}
`, name)
}

func TestAccDataSourcePagerDutyIncidentWorkflow_Missing(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePagerDutyIncidentWorkflowConfigBad(name),
				ExpectError: regexp.MustCompile(fmt.Sprintf("unable to locate any incident workflow with name: %s-incorrect", name)),
			},
		},
	})
}

func testAccDataSourcePagerDutyIncidentWorkflowConfigBad(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "input" {
  name = "%[1]s"
}

data "pagerduty_incident_workflow" "%[1]s" {
  name = "%[1]s-incorrect"
}
`, name)

}
