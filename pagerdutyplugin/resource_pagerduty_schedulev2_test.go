package pagerduty

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyScheduleV2_Basic(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	// Use times far enough in the future to avoid effective_since adjustment
	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "name", scheduleName),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "description", "Managed by Terraform"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "id"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.#", "1"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.0.id"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.#", "1"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.0.event.0.id"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.name", "Weekly On-Call"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.recurrence.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.type", "rotating_member_assignment_strategy"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.member.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.member.0.type", "user_member"),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleV2_Update(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	scheduleNameUpdated := fmt.Sprintf("tf-%s-updated", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "name", scheduleName),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "time_zone", "America/New_York"),
				),
			},
			{
				Config: testAccPagerDutyScheduleV2ConfigUpdated(username, email, scheduleNameUpdated, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "name", scheduleNameUpdated),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "time_zone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "description", "Updated by Terraform"),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleV2_MultipleRotations(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2MultipleRotationsConfig(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.#", "2"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.0.id"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.1.id"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.name", "Primary Rotation"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.1.event.0.name", "Secondary Rotation"),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleV2_Import(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
				),
			},
			{
				ResourceName:      "pagerduty_schedulev2.test",
				ImportState:       true,
				ImportStateVerify: true,
				// effective_since may be adjusted by the API for past times.
				// start_time/end_time are UTC-normalized by the API on import (no prior state to compare against).
				ImportStateVerifyIgnore: []string{
					"rotation.0.event.0.effective_since",
					"rotation.0.event.0.start_time",
					"rotation.0.event.0.end_time",
					// shifts_per_member is not set in the config so state preserves null,
					// but import reads the API's stored value (1). One normalizing apply
					// after import brings them back into sync.
					"rotation.0.event.0.assignment_strategy.0.shifts_per_member",
				},
			},
		},
	})
}

// --- Helper functions ---

func testAccCheckPagerDutyScheduleV2Destroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_schedulev2" {
			continue
		}
		ctx := context.Background()
		if _, err := testAccProvider.client.GetScheduleV3(ctx, r.Primary.ID); err == nil {
			return fmt.Errorf("v3 schedule still exists: %s", r.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPagerDutyScheduleV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if r.Primary.ID == "" {
			return fmt.Errorf("no ID set for %s", n)
		}
		ctx := context.Background()
		if _, err := testAccProvider.client.GetScheduleV3(ctx, r.Primary.ID); err != nil {
			return fmt.Errorf("error fetching v3 schedule %s: %s", r.Primary.ID, err)
		}
		return nil
	}
}

// --- Config functions ---

func testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Managed by Terraform"

  rotation {
    event {
      name            = "Weekly On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince)
}

func testAccPagerDutyScheduleV2ConfigUpdated(username, email, scheduleName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/Los_Angeles"
  description = "Updated by Terraform"

  rotation {
    event {
      name            = "Weekly On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince)
}

func TestAccPagerDutyScheduleV2_PastEffectiveSince(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	// Use a date 48 hours in the past. The API normalizes past effective_since
	// values to the current server time. The provider must preserve the configured
	// value in state to avoid a Framework "inconsistent result after apply" error.
	effectiveSince := time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					// State must hold the configured value, not the API-normalized current time.
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.effective_since", effectiveSince),
				),
			},
			{
				// Re-plan with the same config. Fails if there is a perpetual diff
				// caused by the read-side normalization not preserving the state value.
				Config:             testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccPagerDutyScheduleV2_ActiveEventFieldChange reproduces issue #1123:
// when an event becomes active (effective_since in the past) and the user
// changes a shift-defining field (start_time/end_time), the provider's
// Update path must DELETE+CREATE the underlying event because the API
// rejects PUTs that mutate those fields on active events. The new event
// gets a new server-assigned ID. Without a plan modifier that marks the
// ID as unknown for this case, Terraform fails the apply with
// "Provider produced inconsistent result after apply".
func TestAccPagerDutyScheduleV2_ActiveEventFieldChange(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	// Past effective_since makes the event active immediately. The API
	// normalizes the value server-side; the provider preserves the configured
	// value in state, so subsequent updates still see an active event.
	effectiveSince := time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"
	startTimeUpdated := time.Now().UTC().Add(48*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTimeUpdated := time.Now().UTC().Add(48*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.0.event.0.id"),
				),
			},
			{
				// Shift the start/end by 24h on an already-active event.
				// Pre-fix this apply errors with "inconsistent result after apply"
				// because the new event ID differs from state.
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTimeUpdated, endTimeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.0.event.0.id"),
				),
			},
			{
				// Re-plan with the same config: must be empty.
				Config:             testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTimeUpdated, endTimeUpdated),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccPagerDutyScheduleV2_EffectiveUntil(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	effectiveUntil := time.Now().UTC().Add(30 * 24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2EffectiveUntilConfig(username, email, scheduleName, effectiveSince, effectiveUntil, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.effective_until", effectiveUntil),
				),
			},
			{
				// Remove effective_until — the event becomes indefinite.
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckNoResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.effective_until"),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleV2_UpdateEvents(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"
	startTimeUpdated := time.Now().UTC().Add(48*time.Hour).Format("2006-01-02") + "T08:00:00Z"
	endTimeUpdated := time.Now().UTC().Add(48*time.Hour).Format("2006-01-02") + "T20:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2UpdateEventsConfig(username, email, scheduleName, effectiveSince, startTime, endTime, "Weekly On-Call", "RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.name", "Weekly On-Call"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.recurrence.0", "RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"),
				),
			},
			{
				Config: testAccPagerDutyScheduleV2UpdateEventsConfig(username, email, scheduleName, effectiveSince, startTimeUpdated, endTimeUpdated, "Daily On-Call", "RRULE:FREQ=DAILY"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.name", "Daily On-Call"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.recurrence.0", "RRULE:FREQ=DAILY"),
				),
			},
		},
	})
}

func TestAccPagerDutyScheduleV2_RotationCountChange(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.#", "1"),
				),
			},
			{
				// Add a second rotation, exercising the new-rotation creation branch.
				Config: testAccPagerDutyScheduleV2MultipleRotationsConfig(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.#", "2"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.0.id"),
					resource.TestCheckResourceAttrSet("pagerduty_schedulev2.test", "rotation.1.id"),
				),
			},
			{
				// Remove the second rotation, exercising the rotation deletion branch.
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.#", "1"),
				),
			},
		},
	})
}

func testAccPagerDutyScheduleV2MultipleRotationsConfig(username, email, scheduleName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Multi-rotation schedule"

  rotation {
    event {
      name            = "Primary Rotation"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }

  rotation {
    event {
      name            = "Secondary Rotation"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=TU,TH"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince, startTime, endTime, effectiveSince)
}

func testAccPagerDutyScheduleV2EffectiveUntilConfig(username, email, scheduleName, effectiveSince, effectiveUntil, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Managed by Terraform"

  rotation {
    event {
      name            = "Weekly On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      effective_until = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince, effectiveUntil)
}

func testAccPagerDutyScheduleV2UpdateEventsConfig(username, email, scheduleName, effectiveSince, startTime, endTime, eventName, recurrence string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Managed by Terraform"

  rotation {
    event {
      name            = "%s"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["%s"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, eventName, startTime, endTime, effectiveSince, recurrence)
}

func TestAccPagerDutyScheduleV2_EveryMemberStrategy(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2EveryMemberStrategyConfig(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "name", scheduleName),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.type", "every_member_assignment_strategy"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.member.#", "1"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.member.0.type", "user_member"),
				),
			},
		},
	})
}

func testAccPagerDutyScheduleV2EveryMemberStrategyConfig(username, email, scheduleName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Every-member strategy test"

  rotation {
    event {
      name            = "All Hands On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "every_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince)
}

// TestAccPagerDutyScheduleV2_WithEscalationPolicy covers the bug reported in
// GitHub issue #1105 where any state refresh on a schedule that is referenced
// by an escalation policy caused a JSON unmarshal panic:
// "cannot unmarshal object into Go struct field
//
//	ScheduleV3.schedule.escalation_policies of type string"
func TestAccPagerDutyScheduleV2_WithEscalationPolicy(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	escalationPolicyName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				// Create the schedule and an escalation policy that references it.
				// On the subsequent plan/refresh the API will return escalation_policies
				// as objects — this step validates the fix for issue #1105.
				Config: testAccPagerDutyScheduleV2WithEscalationPolicyConfig(username, email, scheduleName, escalationPolicyName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "name", scheduleName),
				),
			},
			{
				// Second plan/refresh — triggers the exact read path that caused
				// the unmarshal panic. If the fix is correct, this must produce no error.
				Config:   testAccPagerDutyScheduleV2WithEscalationPolicyConfig(username, email, scheduleName, escalationPolicyName, effectiveSince, startTime, endTime),
				PlanOnly: true,
			},
		},
	})
}

func testAccPagerDutyScheduleV2WithEscalationPolicyConfig(username, email, scheduleName, escalationPolicyName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Managed by Terraform"

  rotation {
    event {
      name            = "Weekly On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}

resource "pagerduty_escalation_policy" "test" {
  name        = "%s"
  description = "Managed by Terraform"
  num_loops   = 1

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "schedule_reference"
      id   = pagerduty_schedulev2.test.id
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince, escalationPolicyName)
}

// TestAccPagerDutyScheduleV2_MultipleTeamsOrderStable reproduces the team
// reorder perpetual diff (same class as issue #1001): the v3 API returns a
// schedule's associated teams in an unstable order, so a `teams` list with more
// than one entry would show a spurious reorder on every plan. The provider must
// preserve the prior-state ordering on refresh so re-planning yields no diff.
func TestAccPagerDutyScheduleV2_MultipleTeamsOrderStable(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	team1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2MultipleTeamsConfig(username, email, team1, team2, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "teams.#", "2"),
				),
			},
			{
				// Re-plan with the same config. Fails if the unordered API team
				// response produces a perpetual reorder diff.
				Config:             testAccPagerDutyScheduleV2MultipleTeamsConfig(username, email, team1, team2, scheduleName, effectiveSince, startTime, endTime),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccPagerDutyScheduleV2_OffsetTimestamps reproduces the "inconsistent
// result after apply" failure on start_time/end_time supplied with a timezone
// offset (e.g. "+02:00"). The v3 API normalizes those to UTC ("Z"), so the
// provider must preserve the configured strings in state when they represent the
// same instant — both at apply time (consistency check) and on refresh (no diff).
func TestAccPagerDutyScheduleV2_OffsetTimestamps(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	day := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	startTime := day + "T09:00:00+02:00"
	endTime := day + "T17:00:00+02:00"
	effectiveSince := day + "T09:00:00+02:00"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					// State must hold the configured offset strings, not the API's UTC form.
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.start_time", startTime),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.end_time", endTime),
				),
			},
			{
				// Re-plan with the same config. Fails if the UTC-normalized read
				// value is not reconciled against the configured offset value.
				Config:             testAccPagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccPagerDutyScheduleV2_OffsetEffectiveUntil guards the optional
// effective_until field against the same UTC-normalization problem as
// start_time/end_time: a timezone-offset value (e.g. "+02:00") would otherwise
// fail with "Provider produced inconsistent result after apply" on create and a
// perpetual diff on refresh.
func TestAccPagerDutyScheduleV2_OffsetEffectiveUntil(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	loc := time.FixedZone("UTC+2", 2*3600)
	day := time.Now().Add(24 * time.Hour)
	startTime := time.Date(day.Year(), day.Month(), day.Day(), 9, 0, 0, 0, loc).Format(time.RFC3339)
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 17, 0, 0, 0, loc).Format(time.RFC3339)
	effectiveSince := time.Date(day.Year(), day.Month(), day.Day(), 9, 0, 0, 0, loc).Format(time.RFC3339)
	effectiveUntil := time.Now().Add(30 * 24 * time.Hour).In(loc).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2EffectiveUntilConfig(username, email, scheduleName, effectiveSince, effectiveUntil, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					// State must hold the configured offset string, not the API's UTC form.
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.effective_until", effectiveUntil),
				),
			},
			{
				Config:             testAccPagerDutyScheduleV2EffectiveUntilConfig(username, email, scheduleName, effectiveSince, effectiveUntil, startTime, endTime),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccPagerDutyScheduleV2_TimeRepresentationChangeStable guards the
// reconcile/equality path for active events. When state holds the API's UTC form
// (e.g. after import) and the config expresses the same instants with a timezone
// offset, the event is semantically unchanged: it must NOT be DELETE+CREATEd
// (its server-assigned ID must stay the same) and the apply must not error with
// "inconsistent result after apply". The event is made active (past
// effective_since) so a spurious recreate would destroy a live shift.
func TestAccPagerDutyScheduleV2_TimeRepresentationChangeStable(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	// Same instants expressed two ways: UTC "Z" form and an equivalent "+02:00"
	// offset form. effective_since is in the past so the event is active.
	loc := time.FixedZone("UTC+2", 2*3600)
	esInstant := time.Now().UTC().Add(-48 * time.Hour)
	day := time.Now().UTC().Add(24 * time.Hour)
	startInstant := time.Date(day.Year(), day.Month(), day.Day(), 9, 0, 0, 0, time.UTC)
	endInstant := time.Date(day.Year(), day.Month(), day.Day(), 17, 0, 0, 0, time.UTC)

	esZ := esInstant.Format(time.RFC3339)
	stZ := startInstant.Format(time.RFC3339)
	etZ := endInstant.Format(time.RFC3339)

	esOff := esInstant.In(loc).Format(time.RFC3339)
	stOff := startInstant.In(loc).Format(time.RFC3339)
	etOff := endInstant.In(loc).Format(time.RFC3339)

	var eventID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, esZ, stZ, etZ),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					testAccStoreScheduleV2EventID("pagerduty_schedulev2.test", &eventID),
				),
			},
			{
				// Same instants, now in offset form. The event must survive in place.
				Config: testAccPagerDutyScheduleV2Config(username, email, scheduleName, esOff, stOff, etOff),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					// ID unchanged proves the active event was not recreated.
					testAccCheckScheduleV2EventIDStable("pagerduty_schedulev2.test", &eventID),
					// State adopts the configured offset representation.
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.start_time", stOff),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.effective_since", esOff),
				),
			},
			{
				Config:             testAccPagerDutyScheduleV2Config(username, email, scheduleName, esOff, stOff, etOff),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccPagerDutyScheduleV2_NoDescription covers GitHub issue #1151: creating a
// schedule whose config omits the Optional+Computed description must resolve it
// from the API response instead of leaving it unknown, which fails the apply
// with "Provider returned invalid result object after apply" and taints the
// resource.
func TestAccPagerDutyScheduleV2_NoDescription(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_V3"); v == "" {
		t.Skip("PAGERDUTY_ACC_SCHEDULE_V3 must be set to run v3 schedule acceptance tests")
	}
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	scheduleName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	effectiveSince := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	startTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T09:00:00Z"
	endTime := time.Now().UTC().Add(24*time.Hour).Format("2006-01-02") + "T17:00:00Z"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyScheduleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPagerDutyScheduleV2NoDescriptionConfig(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleV2Exists("pagerduty_schedulev2.test"),
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "description", ""),
				),
			},
			{
				Config:             testAccPagerDutyScheduleV2NoDescriptionConfig(username, email, scheduleName, effectiveSince, startTime, endTime),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// testAccStoreScheduleV2EventID records the first event's server-assigned ID so a
// later step can assert it did not change.
func testAccStoreScheduleV2EventID(n string, dst *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		id := r.Primary.Attributes["rotation.0.event.0.id"]
		if id == "" {
			return fmt.Errorf("no event id set for %s", n)
		}
		*dst = id
		return nil
	}
}

// testAccCheckScheduleV2EventIDStable fails if the first event's ID differs from
// the previously stored value, which would indicate a DELETE+CREATE recreate.
func testAccCheckScheduleV2EventIDStable(n string, want *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		got := r.Primary.Attributes["rotation.0.event.0.id"]
		if got != *want {
			return fmt.Errorf("event was recreated: id was %q, now %q", *want, got)
		}
		return nil
	}
}

func testAccPagerDutyScheduleV2MultipleTeamsConfig(username, email, team1, team2, scheduleName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_team" "test1" {
  name = "%s"
}

resource "pagerduty_team" "test2" {
  name = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name        = "%s"
  time_zone   = "America/New_York"
  description = "Managed by Terraform"

  teams = [
    pagerduty_team.test1.id,
    pagerduty_team.test2.id,
  ]

  rotation {
    event {
      name            = "Weekly On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, team1, team2, scheduleName, startTime, endTime, effectiveSince)
}

func testAccPagerDutyScheduleV2NoDescriptionConfig(username, email, scheduleName, effectiveSince, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedulev2" "test" {
  name      = "%s"
  time_zone = "America/New_York"

  rotation {
    event {
      name            = "Weekly On-Call"
      start_time      = "%s"
      end_time        = "%s"
      effective_since = "%s"
      recurrence      = ["RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"]

      assignment_strategy {
        type = "rotating_member_assignment_strategy"

        member {
          type    = "user_member"
          user_id = pagerduty_user.test.id
        }
      }
    }
  }
}
`, username, email, scheduleName, startTime, endTime, effectiveSince)
}
