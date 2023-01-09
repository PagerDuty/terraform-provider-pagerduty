package pagerduty

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_incident_workflows", &resource.Sweeper{
		Name:         "pagerduty_incident_workflows",
		Dependencies: []string{"pagerduty_incident_workflow_triggers"},
		F:            testSweepIncidentWorkflow,
	})
}

func testSweepIncidentWorkflow(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	workflowsResp, _, err := client.IncidentWorkflows.List(&pagerduty.ListIncidentWorkflowOptions{})
	if err != nil {
		return err
	}

	for _, iw := range workflowsResp.IncidentWorkflows {
		if strings.HasPrefix(iw.Name, "tf-") {
			log.Printf("Destroying incident workflow %s (%s)", iw.Name, iw.ID)
			if _, err := client.IncidentWorkflows.Delete(iw.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyIncidentWorkflow_Basic(t *testing.T) {
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
				Config: testAccCheckPagerDutyIncidentWorkflowConfig(workflowName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowExists("pagerduty_incident_workflow.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "name", workflowName),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.#", "2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.value", "second update"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowConfigUpdate(workflowName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowExists("pagerduty_incident_workflow.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "name", workflowName),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "description", "some description"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.#", "2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.value", "second update updated"),
				),
			},
		},
	})
}

func TestAccPagerDutyIncidentWorkflow_Team(t *testing.T) {
	workflowName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowConfigWithTeam(workflowName, teamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowExists("pagerduty_incident_workflow.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "name", workflowName),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttrSet("pagerduty_incident_workflow.test", "team"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentWorkflowConfigWithTeam(name, team string) string {
	return fmt.Sprintf(`
%s

resource "pagerduty_incident_workflow" "test" {
  name = "%s"
  team = pagerduty_team.foo.id
}
`, testAccCheckPagerDutyTeamConfig(team), name)
}

func testAccCheckPagerDutyIncidentWorkflowConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "test" {
  name = "%s"
  step {
    name           = "Example Step"
    action         = "pagerduty.com:incident-workflows:send-status-update:1"
    input {
      name = "Message"
      value = "first update"
    }
  }
  step {
    name          = "Another Step"
    action        = "pagerduty.com:incident-workflows:send-status-update:1"
    input {
      name = "Message"
      value = "second update"
    }
  }
}
`, name)
}

func testAccCheckPagerDutyIncidentWorkflowConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "test" {
  name = "%s"
  description = "some description"
  step {
    name           = "Example Step"
    action         = "pagerduty.com:incident-workflows:send-status-update:1"
    input {
      name = "Message"
      value = "first update"
    }
  }
  step {
    name          = "Another Step"
    action        = "pagerduty.com:incident-workflows:send-status-update:1"
    input {
      name = "Message"
      value = "second update updated"
    }
  }
}
`, name)
}

func testAccCheckPagerDutyIncidentWorkflowDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_incident_workflow" {
			continue
		}

		if _, _, err := client.IncidentWorkflows.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("incident workflow still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyIncidentWorkflowExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no incident workflow ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.IncidentWorkflows.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("incident workflow not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func TestFlattenIncidentWorkflowStepsOneGenerated(t *testing.T) {

	n := &pagerduty.IncidentWorkflow{
		Steps: []*pagerduty.IncidentWorkflowStep{
			{
				ID:   "abc-123",
				Name: "something",
				Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
					Inputs: []*pagerduty.IncidentWorkflowActionInput{
						{
							Name:  "test1-value",
							Value: "test1-value",
						},
						{
							Name:  "test2-value",
							Value: "test2-value",
						},
						{
							Name:  "test3-value",
							Value: "test3-value",
						},
					},
				},
			},
		},
	}
	o := map[string][]string{
		"abc-123": {"test1-value", "test2-value"},
	}
	r := flattenIncidentWorkflowSteps(n, o)
	l := r[0]["input"].(*[]interface{})
	if len(*l) != 3 {
		t.Errorf("flattened step had wrong number of inputs. want 2 got %v", len(*l))
	}
	for i, v := range *l {
		if i < 2 {
			if _, hadGen := v.(map[string]interface{})["generated"]; hadGen {
				t.Errorf("was not expecting input %v to be generated", i)
			}
		} else {
			if gen, hadGen := v.(map[string]interface{})["generated"]; !hadGen || !(gen.(bool)) {
				t.Errorf("was expecting input %v to be generated", i)
			}
		}
	}
}

func testAccPreCheckIncidentWorkflows(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_INCIDENT_WORKFLOWS"); v == "" {
		t.Skip("PAGERDUTY_ACC_INCIDENT_WORKFLOWS not set. Skipping Incident Workflows-related test")
	}
}
