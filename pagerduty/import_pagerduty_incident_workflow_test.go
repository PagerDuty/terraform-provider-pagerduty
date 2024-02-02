package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyIncidentWorkflow_import(t *testing.T) {
	workflowName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowConfigNoSteps(workflowName),
			},

			{
				ResourceName:      "pagerduty_incident_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyIncidentWorkflowConfigNoSteps(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "test" {
  name = "%s"
  description = "some description"
}
`, name)
}
