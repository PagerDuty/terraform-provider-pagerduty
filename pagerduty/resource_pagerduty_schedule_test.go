package pagerduty

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_schedule", &resource.Sweeper{
		Name: "pagerduty_schedule",
		F:    testSweepSchedule,
	})
}

func testSweepSchedule(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	resp, _, err := client.Schedules.List(&pagerduty.ListSchedulesOptions{})
	if err != nil {
		return err
	}

	for _, schedule := range resp.Schedules {
		if strings.HasPrefix(schedule.Name, "test") || strings.HasPrefix(schedule.Name, "tf-") {
			log.Printf("Destroying schedule %s (%s)", schedule.Name, schedule.ID)
			if _, err := client.Schedules.Delete(schedule.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutySchedule_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	scheduleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	startWrongFormated := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC1123)
	startNotRounded := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Add(5 * time.Second).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rendered_coverage_percentage", "0.00"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "final_schedule.0.rendered_coverage_percentage", "0.00"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleConfigUpdated(username, email, scheduleUpdated, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", scheduleUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),
				),
			},
			{
				Config:      testAccCheckPagerDutyScheduleConfigRestrictionType(username, email, schedule, location, start, rotationVirtualStart),
				ExpectError: regexp.MustCompile("start_day_of_week must only be set for a weekly_restriction schedule restriction type"),
			},
			// Validating that a Weekly Restriction with no Start Day of Week set
			// returns a format error.
			{
				Config:      testAccCheckPagerDutyScheduleConfigRestrictionTypeWeeklyWithoutStartDayOfWeekSet(username, email, schedule, location, start, rotationVirtualStart),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("start_day_of_week must be set for a weekly_restriction schedule restriction type"),
			},
			// Validating that wrong formatted values for "start" attribute return a
			// format error.
			{
				Config:      testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, startWrongFormated, rotationVirtualStart),
				ExpectError: regexp.MustCompile("is not a valid format for argument:"),
			},
			// Validating that dates not minute rounded for "start" attribute are
			// acepted.
			{
				Config: testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, startNotRounded, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", startNotRounded),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleWithTeams_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	scheduleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	teamUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsConfig(username, email, schedule, location, start, rotationVirtualStart, team),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "teams.#", "1"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsConfigUpdated(username, email, scheduleUpdated, location, start, rotationVirtualStart, teamUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", scheduleUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "teams.#", "1"),
				),
			},
		},
	})
}

func TestAccPagerDutySchedule_BasicWithExternalDestroyHandling(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
				),
			},
			// Validating that externally removed schedule are detected and planed for
			// re-creation
			{
				Config: testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccExternallyDestroySchedule("pagerduty_schedule.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPagerDutyScheduleWithTeams_EscalationPolicyDependant(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy := fmt.Sprintf("ts-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantConfigUpdated(username, email, team, escalationPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleNoExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleWithTeams_EscalationPolicyDependantWithOneLayer(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy1 := fmt.Sprintf("ts-%s", acctest.RandString(5))
	escalationPolicy2 := fmt.Sprintf("ts-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckScheduleUsedByEPWithOneLayer(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOneLayerConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
				),
			},
			// Validating that deleting a Schedule used by an Escalation Policy with
			// one configured layer prompts the expected error.
			{
				Config:      testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantConfigUpdated(username, email, team, escalationPolicy1),
				ExpectError: regexp.MustCompile("It is not possible to continue with the destruction of the Schedule \".*\", because it is being used by the Escalation Policy \".*\" which has only one layer configured"),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithMultipleLayersUsingTheSameScheduleAsTargetConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
				),
			},
			// Validating that deleting a Schedule used by an Escalation Policy with
			// multiple configured layer but each layer has configured only the
			// Schedule try to be deleted
			{
				Config:      testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithMultipleLayersUsingTheSameScheduleAsTargetConfigUpdated(username, email, team, escalationPolicy1),
				ExpectError: regexp.MustCompile("It is not possible to continue with the destruction of the Schedule \".*\", because it is being used by the Escalation Policy \".*\" which has only one layer configured"),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantMultipleWithOneLayerConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, escalationPolicy2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
				),
			},
			// Validation that deleting a Schedule used by multiple Escalation
			// Policies with one configured layer prompts the expected error.
			{
				Config:      testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantMultipleWithOneLayerConfigUpdated(username, email, team, escalationPolicy1, escalationPolicy2),
				ExpectError: regexp.MustCompile("It is not possible to continue with the destruction of the Schedule \".*\", because it is being used by multiple Escalation Policies which have only one layer configured."),
			},
		},
	})
}

func TestAccPagerDutyScheduleWithTeams_EscalationPolicyDependantWithOpenIncidents(t *testing.T) {
	service1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicy1 := fmt.Sprintf("ts-%s", acctest.RandString(5))
	escalationPolicy2 := fmt.Sprintf("ts-%s", acctest.RandString(5))
	incident_id := ""
	p_incident_id := &incident_id
	unrelated_incident_id := ""
	p_unrelated_incident_id := &unrelated_incident_id

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOpenIncidentConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, service1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
					testAccCheckPagerDutyScheduleOpenIncidentOnService(p_incident_id, "pagerduty_service.foo", "pagerduty_escalation_policy.foo"),
				),
			},
			{
				Config:      testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOpenIncidentConfigUpdated(username, email, team, escalationPolicy1, service1),
				ExpectError: regexp.MustCompile("Before destroying Schedule \".*\" You must first resolve or reassign the following incidents related with Escalation Policies using this Schedule"),
			},
			{
				// Extra intermediate step with the original plan for resolving the
				// outstanding incident and retrying the schedule destroy after that.
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOpenIncidentConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, service1),
				Check: resource.ComposeTestCheckFunc(
					testAccPagerDutyScheduleResolveIncident(p_incident_id, "pagerduty_escalation_policy.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOpenIncidentConfigUpdated(username, email, team, escalationPolicy1, service1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithUnrelatedOpenIncidentConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, escalationPolicy2, service1, service2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.bar", "name", escalationPolicy2),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service1),
					resource.TestCheckResourceAttr(
						"pagerduty_service.bar", "name", service2),
					testAccCheckPagerDutyScheduleOpenIncidentOnService(p_unrelated_incident_id, "pagerduty_service.bar", "pagerduty_escalation_policy.bar"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithUnrelatedOpenIncidentConfigUpdated(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, escalationPolicy2, service1, service2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.bar", "name", escalationPolicy2),
					testAccPagerDutyScheduleResolveIncident(p_unrelated_incident_id, "pagerduty_escalation_policy.bar"),
				),
			},
		},
	})
}

func TestAccPagerDutySchedule_EscalationPolicyDependantWithOpenIncidents(t *testing.T) {
	service1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	service2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	schedule2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	escalationPolicy1 := fmt.Sprintf("ts-%s", acctest.RandString(5))
	escalationPolicy2 := fmt.Sprintf("ts-%s", acctest.RandString(5))
	incident_id := ""
	p_incident_id := &incident_id
	unrelated_incident_id := ""
	p_unrelated_incident_id := &unrelated_incident_id

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleEscalationPolicyDependantWithOpenIncidentConfig(username, email, schedule1, location, start, rotationVirtualStart, escalationPolicy1, service1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule1),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
					testAccCheckPagerDutyScheduleOpenIncidentOnService(p_incident_id, "pagerduty_service.foo", "pagerduty_escalation_policy.foo"),
				),
			},
			{
				Config:      testAccCheckPagerDutyScheduleEscalationPolicyDependantWithOpenIncidentConfigUpdated(username, email, escalationPolicy1, service1),
				ExpectError: regexp.MustCompile("Before destroying Schedule \".*\" You must first resolve or reassign the following incidents related with Escalation Policies using this Schedule"),
			},
			{
				// Extra intermediate step with the original plan for resolving the
				// outstanding incident and retrying the schedule destroy after that.
				Config: testAccCheckPagerDutyScheduleEscalationPolicyDependantWithOpenIncidentConfig(username, email, schedule1, location, start, rotationVirtualStart, escalationPolicy1, service1),
				Check: resource.ComposeTestCheckFunc(
					testAccPagerDutyScheduleResolveIncident(p_incident_id, "pagerduty_escalation_policy.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleEscalationPolicyDependantWithOpenIncidentConfigUpdated(username, email, escalationPolicy1, service1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleEscalationPolicyDependantWithUnrelatedOpenIncidentConfig(username, email, schedule1, schedule2, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule1),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.bar", "name", schedule2),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.bar", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.bar", "name", escalationPolicy2),
					resource.TestCheckResourceAttr(
						"pagerduty_service.foo", "name", service1),
					resource.TestCheckResourceAttr(
						"pagerduty_service.bar", "name", service2),
					testAccCheckPagerDutyScheduleOpenIncidentOnService(p_unrelated_incident_id, "pagerduty_service.bar", "pagerduty_escalation_policy.bar"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleEscalationPolicyDependantWithUnrelatedOpenIncidentConfigUpdated(username, email, schedule1, schedule2, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.bar", "name", schedule2),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.bar", "description", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.foo", "name", escalationPolicy1),
					resource.TestCheckResourceAttr(
						"pagerduty_escalation_policy.bar", "name", escalationPolicy2),
					testAccPagerDutyScheduleResolveIncident(p_unrelated_incident_id, "pagerduty_escalation_policy.bar"),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleOverflow_Basic(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	scheduleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "America/New_York"
	start := timeNowInLoc(location).Add(30 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(30 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleOverflowConfig(username, email, schedule, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleOverflowConfigUpdated(username, email, scheduleUpdated, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
				),
			},
		},
	})
}

func TestAccPagerDutySchedule_BasicWeek(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	scheduleUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "Australia/Melbourne"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleConfigWeek(username, email, schedule, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.restriction.0.start_day_of_week", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleConfigWeekUpdated(username, email, scheduleUpdated, location, start, rotationVirtualStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", scheduleUpdated),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.restriction.0.start_day_of_week", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),
				),
			},
		},
	})
}

func TestAccPagerDutySchedule_Multi(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "Europe/Berlin"
	start := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	end := timeNowInLoc(location).Add(72 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := timeNowInLoc(location).Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleConfigMulti(username, email, schedule, location, start, rotationVirtualStart, end),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),

					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "3"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.name", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.restriction.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.restriction.0.duration_seconds", "32101"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.restriction.0.start_time_of_day", "08:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_turn_length_seconds", "86400"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.users.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.0.rotation_virtual_start", rotationVirtualStart),

					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.name", "bar"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.restriction.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.restriction.0.duration_seconds", "32101"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.restriction.0.start_time_of_day", "08:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.restriction.0.start_day_of_week", "5"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.rotation_turn_length_seconds", "86400"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.users.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.end", end),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.1.rotation_virtual_start", rotationVirtualStart),

					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.name", "foobar"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.restriction.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.restriction.0.duration_seconds", "32101"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.restriction.0.start_time_of_day", "08:00:00"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.restriction.0.start_day_of_week", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.rotation_turn_length_seconds", "86400"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.users.#", "1"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.start", start),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.2.rotation_virtual_start", rotationVirtualStart),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleConfigMultiUpdated(username, email, schedule, location, start, rotationVirtualStart, end),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleExists("pagerduty_schedule.foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "name", schedule),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "description", "foo"),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "time_zone", location),
					resource.TestCheckResourceAttr(
						"pagerduty_schedule.foo", "layer.#", "2"),
				),
			},
		},
	})
}

func testAccCheckPagerDutyScheduleDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_schedule" {
			continue
		}

		if _, _, err := client.Schedules.Get(r.Primary.ID, &pagerduty.GetScheduleOptions{}); err == nil {
			return fmt.Errorf("Schedule still exists")
		}

	}
	return nil
}

func testAccCheckPagerDutyScheduleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schedule ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.Schedules.Get(rs.Primary.ID, &pagerduty.GetScheduleOptions{})
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Schedule not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckPagerDutyScheduleNoExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return nil
		}
		if rs != nil && rs.Primary.ID == "" {
			return nil
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		found, _, err := client.Schedules.Get(rs.Primary.ID, &pagerduty.GetScheduleOptions{})
		if err != nil {
			return err
		}

		if found.ID == rs.Primary.ID {
			return fmt.Errorf("Schedule still exists: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccExternallyDestroySchedule(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schedule ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()
		_, err := client.Schedules.Delete(rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleConfigRestrictionType(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
	  start_day_of_week = 5
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleConfigRestrictionTypeWeeklyWithoutStartDayOfWeekSet(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "weekly_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleConfigUpdated(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone = "%s"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleOverflowConfig(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name      = "%s"
  overflow  = true
  time_zone = "%s"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleOverflowConfigUpdated(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name      = "%s"
  overflow  = false
  time_zone = "%s"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleConfigWeek(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "weekly_restriction"
      start_time_of_day = "08:00:00"
			start_day_of_week = 1
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleConfigWeekUpdated(username, email, schedule, location, start, rotationVirtualStart string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone = "%s"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

		restriction {
      type              = "weekly_restriction"
      start_time_of_day = "08:00:00"
			start_day_of_week = 5
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleConfigMulti(username, email, schedule, location, start, rotationVirtualStart, end string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
	start                        = "%[5]v"
	end = null
    rotation_virtual_start       = "%[6]v"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }

  layer {
    name                         = "bar"
	start                        = "%[5]v"
	end							 = "%[7]v"
    rotation_virtual_start       = "%[6]v"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "weekly_restriction"
      start_time_of_day = "08:00:00"
			start_day_of_week = 5
      duration_seconds  = 32101
    }
  }

  layer {
    name                         = "foobar"
	start                        = "%[5]v"
	end = null
    rotation_virtual_start       = "%[6]v"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "weekly_restriction"
      start_time_of_day = "08:00:00"
			start_day_of_week = 1
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart, end)
}

func testAccCheckPagerDutyScheduleConfigMultiUpdated(username, email, schedule, location, start, rotationVirtualStart, end string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name        = "%s"
  email       = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
	start                        = "%[5]v"
	end = null
    rotation_virtual_start       = "%[6]v"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }

  layer {
    name                         = "bar"
	start                        = "%[5]v"
	end							 = "%[7]v"
    rotation_virtual_start       = "%[6]v"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "weekly_restriction"
      start_time_of_day = "08:00:00"
			start_day_of_week = 5
      duration_seconds  = 32101
    }
  }
}
`, username, email, schedule, location, start, rotationVirtualStart, end)
}

func testAccCheckPagerDutyScheduleWithTeamsConfig(username, email, schedule, location, start, rotationVirtualStart, team string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, team, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleWithTeamsConfigUpdated(username, email, schedule, location, start, rotationVirtualStart, team string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "bar"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "Managed by Terraform"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}
`, username, email, team, schedule, location, start, rotationVirtualStart)
}

func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
  
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy)
}
func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantConfigUpdated(username, email, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "bar"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, username, email, team, escalationPolicy)
}
func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOneLayerConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy)
}
func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithMultipleLayersUsingTheSameScheduleAsTargetConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy)
}

func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithMultipleLayersUsingTheSameScheduleAsTargetConfigUpdated(username, email, team, escalationPolicy string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "bar"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, username, email, team, escalationPolicy)
}
func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantMultipleWithOneLayerConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, escaltionPolicy2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
resource "pagerduty_escalation_policy" "bar" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy1, escaltionPolicy2)
}
func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantMultipleWithOneLayerConfigUpdated(username, email, team, escalationPolicy1, escaltionPolicy2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
resource "pagerduty_escalation_policy" "bar" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
`, username, email, team, escalationPolicy1, escaltionPolicy2)
}

func testAccCheckPagerDutyScheduleOpenIncidentOnService(p *string, sn, epn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[sn]
		if !ok {
			return fmt.Errorf("Not found service: %s", sn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		rep, ok := s.RootModule().Resources[epn]
		if !ok {
			return fmt.Errorf("Not found escalation policy: %s", epn)
		}

		if rep.Primary.ID == "" {
			return fmt.Errorf("No Escalation Policy ID is set")
		}

		client, _ := testAccProvider.Meta().(*Config).Client()

		incident := &pagerduty.Incident{
			Type:  "incident",
			Title: fmt.Sprintf("tf-%s", acctest.RandString(5)),
			Service: &pagerduty.ServiceReference{
				ID:   rs.Primary.ID,
				Type: "service_reference",
			},
			EscalationPolicy: &pagerduty.EscalationPolicyReference{
				ID:   rep.Primary.ID,
				Type: "escalation_policy_reference",
			},
		}
		resp, _, err := client.Incidents.Create(incident)
		if err != nil {
			return err
		}

		*p = resp.ID

		return nil
	}
}

func testAccPagerDutyScheduleResolveIncident(p *string, epn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, _ := testAccProvider.Meta().(*Config).Client()

		incident, _, err := client.Incidents.Get(*p)
		if err != nil {
			return err
		}

		// marking incident as resolved
		incident.Status = "resolved"
		_, _, err = client.Incidents.ManageIncidents([]*pagerduty.Incident{
			incident,
		}, &pagerduty.ManageIncidentsOptions{})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOpenIncidentConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
  
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy, service)
}

func testAccCheckPagerDutyScheduleEscalationPolicyDependantWithOpenIncidentConfig(username, email, schedule, location, start, rotationVirtualStart, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name = "%s"

  time_zone   = "%s"
  description = "foo"

  layer {
    name                         = "foo"
    start                        = "%s"
    rotation_virtual_start       = "%s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
  
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}
resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
`, username, email, schedule, location, start, rotationVirtualStart, escalationPolicy, service)
}

func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithUnrelatedOpenIncidentConfig(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, escalationPolicy2, service1, service2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[2]s"
}

resource "pagerduty_team" "foo" {
	name = "%[3]s"
	description = "fighters"
}

resource "pagerduty_schedule" "foo" {
  name = "%[4]s"

  time_zone   = "%[5]s"
  description = "foo"

  teams = [pagerduty_team.foo.id]

  layer {
    name                         = "foo"
    start                        = "%[6]s"
    rotation_virtual_start       = "%[7]s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[8]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
  
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "%[9]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_service" "foo" {
	name                    = "%[10]s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service" "bar" {
	name                    = "%[11]s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2)
}

func testAccCheckPagerDutyScheduleEscalationPolicyDependantWithUnrelatedOpenIncidentConfig(username, email, schedule1, schedule2, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[2]s"
}

resource "pagerduty_schedule" "foo" {
  name = "%[3]s"

  time_zone   = "%[5]s"
  description = "foo"

  layer {
    name                         = "foo"
    start                        = "%[6]s"
    rotation_virtual_start       = "%[7]s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_schedule" "bar" {
  name = "%[4]s"

  time_zone   = "%[5]s"
  description = "bar"

  layer {
    name                         = "bar"
    start                        = "%[6]s"
    rotation_virtual_start       = "%[7]s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[8]s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
  
  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.foo.id
    }
  }
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "%[9]s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.bar.id
    }
  }
}

resource "pagerduty_service" "foo" {
	name                    = "%[10]s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service" "bar" {
	name                    = "%[11]s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
`, username, email, schedule1, schedule2, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2)
}

func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithUnrelatedOpenIncidentConfigUpdated(username, email, schedule, location, start, rotationVirtualStart, team, escalationPolicy1, escalationPolicy2, service1, service2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[2]s"
}

resource "pagerduty_team" "foo" {
	name = "%[3]s"
	description = "fighters"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[8]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "%[9]s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_service" "foo" {
	name                    = "%[10]s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service" "bar" {
	name                    = "%[11]s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
`, username, email, team, schedule, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2)
}

func testAccCheckPagerDutyScheduleEscalationPolicyDependantWithUnrelatedOpenIncidentConfigUpdated(username, email, schedule1, schedule2, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%[1]s"
  email = "%[2]s"
}

resource "pagerduty_schedule" "bar" {
  name = "%[4]s"

  time_zone   = "%[5]s"
  description = "bar"

  layer {
    name                         = "bar"
    start                        = "%[6]s"
    rotation_virtual_start       = "%[7]s"
    rotation_turn_length_seconds = 86400
    users                        = [pagerduty_user.foo.id]

    restriction {
      type              = "daily_restriction"
      start_time_of_day = "08:00:00"
      duration_seconds  = 32101
    }
  }
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%[8]s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}

resource "pagerduty_escalation_policy" "bar" {
  name      = "%[9]s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
    target {
      type = "schedule_reference"
      id   = pagerduty_schedule.bar.id
    }
  }
}

resource "pagerduty_service" "foo" {
	name                    = "%[10]s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
resource "pagerduty_service" "bar" {
	name                    = "%[11]s"
	description             = "bar"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.bar.id
	alert_creation          = "create_incidents"
}
`, username, email, schedule1, schedule2, location, start, rotationVirtualStart, escalationPolicy1, escalationPolicy2, service1, service2)
}

func testAccCheckPagerDutyScheduleWithTeamsEscalationPolicyDependantWithOpenIncidentConfigUpdated(username, email, team, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "foo" {
	name = "%s"
	description = "bar"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2
  teams     = [pagerduty_team.foo.id]

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
`, username, email, team, escalationPolicy, service)
}

func testAccCheckPagerDutyScheduleEscalationPolicyDependantWithOpenIncidentConfigUpdated(username, email, escalationPolicy, service string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "foo" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_escalation_policy" "foo" {
  name      = "%s"
  num_loops = 2

  rule {
    escalation_delay_in_minutes = 10
    target {
      type = "user_reference"
      id   = pagerduty_user.foo.id
    }
  }
}
resource "pagerduty_service" "foo" {
	name                    = "%s"
	description             = "foo"
	auto_resolve_timeout    = 1800
	acknowledgement_timeout = 1800
	escalation_policy       = pagerduty_escalation_policy.foo.id
	alert_creation          = "create_incidents"
}
`, username, email, escalationPolicy, service)
}

func testAccPreCheckScheduleUsedByEPWithOneLayer(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_USED_BY_EP_W_1_LAYER"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_USED_BY_EP_W_1_LAYER not set. Skipping Schedule related test")
	}
}
