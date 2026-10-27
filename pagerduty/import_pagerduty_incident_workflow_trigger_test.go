package pagerduty

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyIncidentWorkflowTrigger_import(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	workflow := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualSingleService(username, email, escalationPolicy, service, workflow),
			},

			{
				ResourceName:      "pagerduty_incident_workflow_trigger.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccPagerDutyIncidentWorkflowTrigger_import_incidentType references a
// pre-existing incident type via PAGERDUTY_ACC_INCIDENT_TYPE_ID rather than
// creating one: incident types cannot be deleted through the API.
func TestAccPagerDutyIncidentWorkflowTrigger_import_incidentType(t *testing.T) {
	workflow := fmt.Sprintf("tf-%s", acctest.RandString(5))
	incidentTypeID := os.Getenv("PAGERDUTY_ACC_INCIDENT_TYPE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
			testAccPreCheckIncidentType(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigIncidentType(workflow, incidentTypeID),
			},

			{
				ResourceName:      "pagerduty_incident_workflow_trigger.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
