package pagerduty

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_incident_workflow_triggers", &resource.Sweeper{
		Name: "pagerduty_incident_workflow_triggers",
		F:    testSweepIncidentWorkflowTrigger,
	})
}

func testSweepIncidentWorkflowTrigger(region string) error {
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
			triggersResp, _, err := client.IncidentWorkflowTriggers.List(&pagerduty.ListIncidentWorkflowTriggerOptions{WorkflowID: iw.ID})
			if err != nil {
				return err
			}

			for _, t := range triggersResp.Triggers {
				log.Printf("Destroying incident workflow trigger %s", t.ID)
				if _, err := client.IncidentWorkflowTriggers.Delete(t.ID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func TestAccPagerDutyIncidentWorkflowTrigger_BadType(t *testing.T) {
	config := `
resource "pagerduty_incident_workflow_trigger" "my_first_workflow_trigger" {
  type             = "dummy"
  workflow         = "ignored"
  subscribed_to_all_services = true
}
`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`"dummy" is an invalid value. Must be one of \[]string{"manual", "conditional"}`),
			},
		},
	})
}

func TestAccPagerDutyIncidentWorkflowTrigger_ConditionWithManualType(t *testing.T) {
	config := `
resource "pagerduty_incident_workflow_trigger" "my_first_workflow_trigger" {
  type             = "manual"
  workflow         = "ignored"
  condition        = "something"
  subscribed_to_all_services = true
}
`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("when trigger type manual is used, condition must not be specified"),
			},
		},
	})
}

func TestAccPagerDutyIncidentWorkflowTrigger_SubscribedToAllWithInvalidServices(t *testing.T) {
	config := `
resource "pagerduty_incident_workflow_trigger" "my_first_workflow_trigger" {
  type       = "conditional"
  workflow   = "ignored"
  condition  = "something"
  subscribed_to_all_services = true
  services = ["abc-123"]
}
`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("when subscribed_to_all_services is true, services must either be not defined or empty"),
			},
		},
	})
}

func TestAccPagerDutyIncidentWorkflowTrigger_BasicManual(t *testing.T) {
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
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualSingleService(username, email, escalationPolicy, service, workflow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "manual"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualSingleService(username, email, escalationPolicy, service, workflow string) string {
	return fmt.Sprintf(`
%s

%s

resource "pagerduty_incident_workflow_trigger" "test" {
  type       = "manual"
  workflow   = pagerduty_incident_workflow.test.id
  services   = [pagerduty_service.foo.id]
  subscribed_to_all_services = false
}
`, testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service), testAccCheckPagerDutyIncidentWorkflowConfig(workflow))
}

func TestAccPagerDutyIncidentWorkflowTrigger_BasicConditionalAllServices(t *testing.T) {
	workflow := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", ""),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, "incident.priority matches 'P1'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", "incident.priority matches 'P1'"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, "incident.priority matches 'P2'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", "incident.priority matches 'P2'"),
				),
			},
		},
	})
}

func TestAccPagerDutyIncidentWorkflowTrigger_ManualWithTeamPermissions(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service := fmt.Sprintf("tf-%s", acctest.RandString(5))
	workflow := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamIDTFRef := "pagerduty_team.foo.id"
	emptyCondition := ""
	dummyCondition := "event.summary matches 'foo'"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissions(username, email, escalationPolicy, service, teamName, workflow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "manual"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "permissions.0.restricted", "false"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissionsUpdated(username, email, escalationPolicy, service, teamName, workflow, "manual", emptyCondition, "true", teamIDTFRef),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "manual"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "permissions.0.restricted", "true"),
					testAccCheckPagerDutyIncidentWorkflowTriggerCheckPermissionsTeamId("pagerduty_incident_workflow_trigger.test", "pagerduty_team.foo"),
				),
			},
			// Check input validation conditions for permissions configuration
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissionsUpdated(username, email, escalationPolicy, service, teamName, workflow, "conditional", dummyCondition, "true", teamIDTFRef),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("restricted can only be true when trigger type is manual"),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissionsUpdated(username, email, escalationPolicy, service, teamName, workflow, "manual", emptyCondition, "false", teamIDTFRef),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("team_id not allowed when restricted is false"),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissionsUpdated(username, email, escalationPolicy, service, teamName, workflow, "manual", emptyCondition, "true", `""`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
				),
				ExpectError: regexp.MustCompile("team_id must be specified when restricted is true"),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissions(username, email, escalationPolicy, service, workflow, team string) string {
	return fmt.Sprintf(`
%s

%s

resource "pagerduty_team" "foo" {
  name = %q
}

resource "pagerduty_incident_workflow_trigger" "test" {
  type                       = "manual"
  workflow                   = pagerduty_incident_workflow.test.id
  services                   = [pagerduty_service.foo.id]
  subscribed_to_all_services = false
}
`, testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service), testAccCheckPagerDutyIncidentWorkflowConfig(workflow), team)
}

func testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualWithPermissionsUpdated(username, email, escalationPolicy, service, workflow, team, triggerType, condition, isRestricted, teamId string) string {
	return fmt.Sprintf(`
%s

%s

resource "pagerduty_team" "foo" {
  name = "%s"
}

resource "pagerduty_incident_workflow_trigger" "test" {
  type                       = "%s"
  condition                  = "%s"
  workflow                   = pagerduty_incident_workflow.test.id
  services                   = [pagerduty_service.foo.id]
  subscribed_to_all_services = false
  permissions {
    restricted = %s
    team_id    = %s
  }
}
`, testAccCheckPagerDutyServiceConfig(username, email, escalationPolicy, service), testAccCheckPagerDutyIncidentWorkflowConfig(workflow), team, triggerType, condition, isRestricted, teamId)
}

func testAccCheckPagerDutyIncidentWorkflowTriggerCheckPermissionsTeamId(iwtName, teamName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsIWT, ok := s.RootModule().Resources[iwtName]
		if !ok {
			return fmt.Errorf("not found: %s", iwtName)
		}
		if rsIWT.Primary.ID == "" {
			return fmt.Errorf("no incident workflow trigger ID is set")
		}

		rsTeam, ok := s.RootModule().Resources[teamName]
		if !ok {
			return fmt.Errorf("not found: %s", teamName)
		}
		if rsTeam.Primary.ID == "" {
			return fmt.Errorf("no team ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.IncidentWorkflowTriggers.Get(rsIWT.Primary.ID)
		if err != nil {
			return err
		}

		if found.Permissions.TeamID != rsTeam.Primary.ID {
			return fmt.Errorf("incident workflow trigger team restriction wanted %q, but got %q", rsTeam.Primary.ID, found.Permissions.TeamID)
		}

		return nil
	}
}

func TestAccPagerDutyIncidentWorkflowTrigger_ChangeTypeCausesReplace(t *testing.T) {
	workflow := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", ""),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, "incident.priority matches 'P1'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", "incident.priority matches 'P1'"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualAllServices(workflow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "manual"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, condition string) string {
	return fmt.Sprintf(`
%s

resource "pagerduty_incident_workflow_trigger" "test" {
  type       = "conditional"
  workflow   = pagerduty_incident_workflow.test.id
  services   = []
  condition  = "%s"
  subscribed_to_all_services = true
}
`, testAccCheckPagerDutyIncidentWorkflowConfig(workflow), condition)
}

func testAccCheckPagerDutyIncidentWorkflowTriggerConfigManualAllServices(workflow string) string {
	return fmt.Sprintf(`
%s

resource "pagerduty_incident_workflow_trigger" "test" {
  type       = "manual"
  workflow   = pagerduty_incident_workflow.test.id
  services   = []
  subscribed_to_all_services = true
}
`, testAccCheckPagerDutyIncidentWorkflowConfig(workflow))
}

func TestAccPagerDutyIncidentWorkflowTrigger_CannotChangeType(t *testing.T) {
	workflow := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckIncidentWorkflows(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPagerDutyIncidentWorkflowTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", ""),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, "incident.priority matches 'P1'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", "incident.priority matches 'P1'"),
					resource.TestCheckResourceAttr("pagerduty_incident_workflow_trigger.test", "subscribed_to_all_services", "true"),
				),
			},
			{
				Config: testAccCheckPagerDutyIncidentWorkflowTriggerConfigConditionalAllServices(workflow, "incident.priority matches 'P2'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyIncidentWorkflowTriggerExists("pagerduty_incident_workflow_trigger.test"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "type", "conditional"),
					resource.TestCheckResourceAttr(
						"pagerduty_incident_workflow_trigger.test", "condition", "incident.priority matches 'P2'"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyIncidentWorkflowTriggerDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_incident_workflow_trigger" {
			continue
		}

		if _, _, err := client.IncidentWorkflowTriggers.Get(r.Primary.ID); err == nil {
			return fmt.Errorf("incident workflow trigger still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyIncidentWorkflowTriggerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no incident workflow trigger ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.IncidentWorkflowTriggers.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("incident workflow trigger not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
