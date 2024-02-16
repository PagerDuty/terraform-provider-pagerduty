package pagerduty

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccPagerDutyIncidentWorkflow_InlineInputs(t *testing.T) {
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
				Config: testAccCheckPagerDutyIncidentWorkflowInlineInputsConfig(workflowName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowExists("pagerduty_incident_workflow.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "name", workflowName),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.#", "3"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.name", "Step 1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.action", "pagerduty.com:incident-workflows:send-status-update:1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.name", "Message"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.value", "first update"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.generated", "false"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.name", "Step 2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.action", "pagerduty.com:logic:incident-workflows-loop-until:2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.#", "3"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.name", "Condition"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.value", "incident.status matches 'resolved'"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.1.name", "Delay between loops"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.1.value", "10"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.1.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.2.name", "Maximum loops"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.2.value", "20"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.2.generated", "true"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.name", "Actions"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.0.name", "Step 2a"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.0.action", "pagerduty.com:incident-workflows:send-status-update:1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.0.input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.0.input.0.name", "Message"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.0.input.0.value", "Loop update"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.inline_steps_input.0.step.0.input.0.generated", "false"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.name", "Step 3"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.action", "pagerduty.com:incident-workflows:add-conference-bridge:5"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.#", "3"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.0.name", "Conference Number"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.0.value", "+1 415-555-1212,,,,1234#,"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.0.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.1.name", "Conference URL"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.1.value", "https://www.testconferenceurl.com/"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.1.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.2.name", "Overwrite existing conference bridge"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.2.value", "No"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.2.input.2.generated", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowInlineInputsConfigUpdate(workflowName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowExists("pagerduty_incident_workflow.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "name", workflowName),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow.test", "description", "some description"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.#", "2"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.name", "Step 1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.action", "pagerduty.com:logic:incident-workflows-loop-until:2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.#", "3"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.name", "Condition"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.value", "incident.status matches 'resolved'"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.0.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.1.name", "Delay between loops"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.1.value", "10"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.1.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.2.name", "Maximum loops"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.2.value", "20"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.input.2.generated", "true"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.name", "Actions"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.#", "2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.0.name", "Step 1a"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.0.action", "pagerduty.com:incident-workflows:send-status-update:1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.0.input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.0.input.0.name", "Message"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.0.input.0.value", "Loop update 1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.0.input.0.generated", "false"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.1.name", "Step 1b"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.1.action", "pagerduty.com:incident-workflows:send-status-update:1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.1.input.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.1.input.0.name", "Message"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.1.input.0.value", "Loop update 2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.0.inline_steps_input.0.step.1.input.0.generated", "false"),

					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.name", "Step 2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.action", "pagerduty.com:incident-workflows:add-conference-bridge:5"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.#", "2"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.name", "Conference Number"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.value", "+1 415-555-1212,,,,1234#,"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.0.generated", "false"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.1.name", "Overwrite existing conference bridge"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.1.value", "Yes"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow.test", "step.1.input.1.generated", "false"),
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

func testAccCheckPagerDutyIncidentWorkflowInlineInputsConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "test" {
  name = "%s"
  step {
    name           = "Step 1"
    action         = "pagerduty.com:incident-workflows:send-status-update:1"
    input {
      name = "Message"
      value = "first update"
    }
  }
  step {
    name           = "Step 2"
    action         = "pagerduty.com:logic:incident-workflows-loop-until:2"
    input {
      name = "Condition"
      value = "incident.status matches 'resolved'"
    }
    input {
      name = "Delay between loops"
      value = "10"
    }
    inline_steps_input {
      name = "Actions"
      step {
        name = "Step 2a"
        action = "pagerduty.com:incident-workflows:send-status-update:1"
        input {
          name = "Message"
          value = "Loop update"
        }
      }
    }
  }
  step {
    name          = "Step 3"
    action        = "pagerduty.com:incident-workflows:add-conference-bridge:5"
    input {
      name = "Conference URL"
      value = "https://www.testconferenceurl.com/"
    }
    input {
      name = "Conference Number"
      value = "+1 415-555-1212,,,,1234#,"
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

func testAccCheckPagerDutyIncidentWorkflowInlineInputsConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_incident_workflow" "test" {
  name = "%s"
  description = "some description"
  step {
    name           = "Step 1"
    action         = "pagerduty.com:logic:incident-workflows-loop-until:2"
    input {
      name = "Condition"
      value = "incident.status matches 'resolved'"
    }
    input {
      name = "Delay between loops"
      value = "10"
    }
    inline_steps_input {
      name = "Actions"
      step {
        name = "Step 1a"
        action = "pagerduty.com:incident-workflows:send-status-update:1"
        input {
          name = "Message"
          value = "Loop update 1"
        }
      }
      step {
        name = "Step 1b"
        action = "pagerduty.com:incident-workflows:send-status-update:1"
        input {
          name = "Message"
          value = "Loop update 2"
        }
      }
    }
  }
  step {
    name          = "Step 2"
    action        = "pagerduty.com:incident-workflows:add-conference-bridge:5"
    input {
      name = "Conference Number"
      value = "+1 415-555-1212,,,,1234#,"
    }
    input {
      name = "Overwrite existing conference bridge"
      value = "Yes"
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
	iw := &pagerduty.IncidentWorkflow{
		Steps: []*pagerduty.IncidentWorkflowStep{
			{
				ID:   "abc-123",
				Name: "step1",
				Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
					ActionID: "step1-action-id",
					Inputs: []*pagerduty.IncidentWorkflowActionInput{
						{
							Name:  "step1-input1",
							Value: "test1-value",
						},
						{
							Name:  "step1-input2",
							Value: "test2-value",
						},
						{
							Name:  "step1-input3",
							Value: "test3-value",
						},
					},
					InlineStepsInputs: []*pagerduty.IncidentWorkflowActionInlineStepsInput{
						{
							Name: "step1-inlineinput1",
							Value: &pagerduty.IncidentWorkflowActionInlineStepsInputValue{
								Steps: []*pagerduty.IncidentWorkflowActionInlineStep{
									{
										Name: "step1a",
										Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
											ActionID: "step1a-action-id",
											Inputs: []*pagerduty.IncidentWorkflowActionInput{
												{
													Name:  "step1a-input1",
													Value: "inlineval1",
												},
												{
													Name:  "step1a-input2",
													Value: "inlineval2",
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "step1-inlineinput2",
							Value: &pagerduty.IncidentWorkflowActionInlineStepsInputValue{
								Steps: []*pagerduty.IncidentWorkflowActionInlineStep{
									{
										Name: "step1b",
										Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
											ActionID: "step1b-action-id",
											Inputs: []*pagerduty.IncidentWorkflowActionInput{
												{
													Name:  "step1b-input1",
													Value: "inlineval3",
												},
												{
													Name:  "step1b-input2",
													Value: "inlineval4",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	specifiedSteps := []*SpecifiedStep{
		{
			// "step1-input1" is generated
			SpecifiedInputNames: []string{"step1-input2", "step1-input3"},
			SpecifiedInlineInputs: map[string][]*SpecifiedStep{
				"step1-inlineinput1": {
					{
						// "step1a-input1" is generated
						SpecifiedInputNames:   []string{"step1a-input2"},
						SpecifiedInlineInputs: map[string][]*SpecifiedStep{},
					},
				},
				"step1-inlineinput2": {
					{
						// "step1b-input2" is generated
						SpecifiedInputNames:   []string{"step1b-input1"},
						SpecifiedInlineInputs: map[string][]*SpecifiedStep{},
					},
				},
			},
		},
	}
	flattenedSteps := flattenIncidentWorkflowSteps(iw, specifiedSteps, false)
	step1Inputs := flattenedSteps[0]["input"].(*[]interface{})
	if len(*step1Inputs) != 3 {
		t.Errorf("flattened step1 had wrong number of inputs. want 3 got %v", len(*step1Inputs))
	}
	for i, v := range *step1Inputs {
		if i == 0 {
			if generated, hadGen := v.(map[string]interface{})["generated"]; !hadGen || !generated.(bool) {
				t.Errorf("was not expecting step1 input %v to have generated=false", i)
			}
		} else {
			if _, hadGen := v.(map[string]interface{})["generated"]; hadGen {
				t.Errorf("was not expecting step1 input %v to have generated key set", i)
			}
		}
	}

	step1aInlineStepInputs := flattenedSteps[0]["inline_steps_input"].(*[]interface{})
	step1aInlineStepInputs1 := (*step1aInlineStepInputs)[0].(map[string]interface{})
	step1aInlineStepInputs1Steps := step1aInlineStepInputs1["step"].(*[]interface{})
	step1aInputs := (*step1aInlineStepInputs1Steps)[0].(map[string]interface{})["input"].(*[]interface{})
	if len(*step1aInputs) != 2 {
		t.Errorf("flattened step1a had wrong number of inputs. want 2 got %v", len(*step1aInputs))
	}
	for i, v := range *step1aInputs {
		if i == 0 {
			if generated, hadGen := v.(map[string]interface{})["generated"]; !hadGen || !generated.(bool) {
				t.Errorf("was not expecting step1a input %v to have generated=false", i)
			}
		} else {
			if _, hadGen := v.(map[string]interface{})["generated"]; hadGen {
				t.Errorf("was not expecting step1a input %v to have generated key set", i)
			}
		}
	}

	step1bInlineStepInputs := flattenedSteps[0]["inline_steps_input"].(*[]interface{})
	step1bInlineStepInputs2 := (*step1bInlineStepInputs)[1].(map[string]interface{})
	step1bInlineStepInputs2Steps := step1bInlineStepInputs2["step"].(*[]interface{})
	step1bInputs := (*step1bInlineStepInputs2Steps)[0].(map[string]interface{})["input"].(*[]interface{})
	if len(*step1bInputs) != 2 {
		t.Errorf("flattened step1b had wrong number of inputs. want 2 got %v", len(*step1aInputs))
	}
	for i, v := range *step1bInputs {
		if i == 0 {
			if _, hadGen := v.(map[string]interface{})["generated"]; hadGen {
				t.Errorf("was not expecting step1b input %v to have generated key set", i)
			}
		} else {
			if generated, hadGen := v.(map[string]interface{})["generated"]; !hadGen || !generated.(bool) {
				t.Errorf("was not expecting step1b input %v to have generated=false", i)
			}
		}
	}
}

func TestFlattenIncidentWorkflowStepsWithoutSpecifiedSteps(t *testing.T) {
	iw := &pagerduty.IncidentWorkflow{
		Steps: []*pagerduty.IncidentWorkflowStep{
			{
				ID:   "abc-123",
				Name: "step1",
				Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
					ActionID: "step1-action-id",
					Inputs: []*pagerduty.IncidentWorkflowActionInput{
						{
							Name:  "step1-input1",
							Value: "test1-value",
						},
					},
					InlineStepsInputs: []*pagerduty.IncidentWorkflowActionInlineStepsInput{
						{
							Name: "step1-inlineinput1",
							Value: &pagerduty.IncidentWorkflowActionInlineStepsInputValue{
								Steps: []*pagerduty.IncidentWorkflowActionInlineStep{
									{
										Name: "step1a",
										Configuration: &pagerduty.IncidentWorkflowActionConfiguration{
											ActionID: "step1a-action-id",
											Inputs: []*pagerduty.IncidentWorkflowActionInput{
												{
													Name:  "step1a-input1",
													Value: "inlineval1",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	flattenedSteps := flattenIncidentWorkflowSteps(iw, nil, true)
	step1Inputs := flattenedSteps[0]["input"].(*[]interface{})
	if len(*step1Inputs) != 1 {
		t.Errorf("flattened step1 had wrong number of inputs. want 1 got %v", len(*step1Inputs))
	}
	for i, v := range *step1Inputs {
		if _, hadGen := v.(map[string]interface{})["generated"]; hadGen {
			t.Errorf("was not expecting step1a input %v to have generated key set", i)
		}
	}

	step1aInlineStepInputs := flattenedSteps[0]["inline_steps_input"].(*[]interface{})
	step1aInlineStepInputs1 := (*step1aInlineStepInputs)[0].(map[string]interface{})
	step1aInlineStepInputs1Steps := step1aInlineStepInputs1["step"].(*[]interface{})
	step1aInputs := (*step1aInlineStepInputs1Steps)[0].(map[string]interface{})["input"].(*[]interface{})
	if len(*step1aInputs) != 1 {
		t.Errorf("flattened step1a had wrong number of inputs. want 1 got %v", len(*step1aInputs))
	}
	for i, v := range *step1aInputs {
		if _, hadGen := v.(map[string]interface{})["generated"]; hadGen {
			t.Errorf("was not expecting step1a input %v to have generated key set", i)
		}
	}
}

func testAccPreCheckIncidentWorkflows(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_INCIDENT_WORKFLOWS"); v == "" {
		t.Skip("PAGERDUTY_ACC_INCIDENT_WORKFLOWS not set. Skipping Incident Workflows-related test")
	}
}
