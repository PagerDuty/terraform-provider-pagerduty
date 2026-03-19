package pagerduty

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyScheduleV2_Basic(t *testing.T) {
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
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyScheduleV2("pagerduty_schedulev2.test", "data.pagerduty_schedulev2.by_name"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyScheduleV2(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		if r == nil {
			return fmt.Errorf("data source %s not found in state", n)
		}
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a v3 schedule ID from PagerDuty")
		}

		for _, att := range []string{"id", "name"} {
			if a[att] != srcA[att] {
				return fmt.Errorf("expected the v3 schedule %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyScheduleV2Config(username, email, scheduleName, effectiveSince, startTime, endTime string) string {
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

data "pagerduty_schedulev2" "by_name" {
  name = pagerduty_schedulev2.test.name
}
`, username, email, scheduleName, startTime, endTime, effectiveSince)
}
