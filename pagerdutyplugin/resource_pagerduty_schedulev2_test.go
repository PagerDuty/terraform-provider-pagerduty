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
					resource.TestCheckResourceAttr("pagerduty_schedulev2.test", "rotation.0.event.0.assignment_strategy.0.type", "user_assignment_strategy"),
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
        type = "user_assignment_strategy"

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
        type = "user_assignment_strategy"

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
        type = "user_assignment_strategy"

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
        type = "user_assignment_strategy"

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
