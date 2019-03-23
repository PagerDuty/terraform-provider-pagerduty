package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func init() {
	resource.AddTestSweepers("pagerduty_schedule_override", &resource.Sweeper{
		Name:         "pagerduty_schedule_override",
		F:            testSweepScheduleOverride,
		Dependencies: []string{},
	})
}

func testSweepScheduleOverride(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client, err := config.Client()
	if err != nil {
		return err
	}

	respUser, _, err := client.Users.List(&pagerduty.ListUsersOptions{})
	if err != nil {
		return err
	}

	for _, user := range respUser.Users {
		if strings.HasPrefix(user.Name, "tf-") {
			log.Printf("Destroying user %s (%s)", user.Name, user.ID)
			if _, err := client.Users.Delete(user.ID); err != nil {
				return err
			}
		}
	}

	resp, _, err := client.Schedules.List(&pagerduty.ListSchedulesOptions{})
	if err != nil {
		return err
	}

	for _, schedule := range resp.Schedules {
		if strings.HasPrefix(schedule.Name, "tf-") {
			log.Printf("Destroying schedule %s (%s)", schedule.Name, schedule.ID)
			if _, err := client.Schedules.Delete(schedule.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccPagerDutyScheduleOverride(t *testing.T) {
	usernameA := fmt.Sprintf("tf-%s", acctest.RandString(5))
	emailA := fmt.Sprintf("tf-%s@foo.com", acctest.RandString(5))
	usernameB := fmt.Sprintf("tf-%s", acctest.RandString(6))
	emailB := fmt.Sprintf("tf-%s@foo.com", acctest.RandString(6))

	scheduleName := "tf-test-override-schedule"

	overrideStart := "2021-04-03T21:00:00-04:00"
	overrideEnd := "2021-04-04T23:00:00-04:00"
	overrideLaterStart := "2021-04-04T21:00:00-04:00"
	overrideLaterEnd := "2021-04-05T23:00:00-04:00"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleOverrideDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyOverrideSimpleConfig(
					usernameA,
					emailA,
					usernameB,
					emailB,
					scheduleName,
					overrideStart,
					overrideEnd,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleOverrideExists("pagerduty_schedule_override.foo_override"),
					resource.TestCheckResourceAttr("pagerduty_user.userA", "name", usernameA),
					resource.TestCheckResourceAttr("pagerduty_user.userB", "name", usernameB),
					resource.TestCheckResourceAttrSet("pagerduty_schedule_override.foo_override", "user"),
					resource.TestCheckResourceAttrSet("pagerduty_schedule_override.foo_override", "start"),
					resource.TestCheckResourceAttrSet("pagerduty_schedule_override.foo_override", "end"),
					resource.TestCheckResourceAttrSet("pagerduty_schedule_override.foo_override", "schedule"),
				),
			},
			{
				Config: testAccCheckPagerDutyOverrideSimpleConfig(
					usernameA,
					emailA,
					usernameB,
					emailB,
					scheduleName,
					overrideLaterStart,
					overrideLaterEnd,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleOverrideExists("pagerduty_schedule_override.foo_override"),
				),
			},
			{
				Config: testAccCheckPagerDutyScheduleBaseConfig(
					usernameA,
					emailA,
					usernameB,
					emailB,
					scheduleName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPagerDutyScheduleNowExistsWithoutOverride(
						fmt.Sprintf("pagerduty_schedule.foo"),
						overrideStart,
						overrideLaterEnd,
					),
				),
			},
		},
	})
}

func testAccCheckPagerDutyScheduleOverrideDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pagerduty.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_schedule_override" {
			continue
		}

		listOverridesOptions := pagerduty.ListOverridesOptions{
			Since: r.Primary.Attributes["start"],
			Until: r.Primary.Attributes["end"],
		}

		overrides, _, err := client.Schedules.ListOverrides(
			r.Primary.Attributes["schedule"],
			&listOverridesOptions,
		)
		if err != nil {
			return err
		}

		if len(overrides.Overrides) != 0 {
			return fmt.Errorf("Override still exists")
		}
	}
	return nil
}

func testAccCheckPagerDutyScheduleNowExistsWithoutOverride(
	n,
	overrideStart,
	overrideEnd string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No schedule ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)

		listOverridesOptions := pagerduty.ListOverridesOptions{
			Since: overrideStart,
			Until: overrideEnd,
		}

		overrides, _, err := client.Schedules.ListOverrides(rs.Primary.ID, &listOverridesOptions)
		if err != nil {
			return err
		}

		if len(overrides.Overrides) != 0 {
			return fmt.Errorf("Expected to find no overrides; found %d", len(overrides.Overrides))
		}

		return nil
	}
}

func testAccCheckPagerDutyScheduleOverrideExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No override ID is set")
		}

		client := testAccProvider.Meta().(*pagerduty.Client)

		listOverridesOptions := pagerduty.ListOverridesOptions{
			Since: rs.Primary.Attributes["start"],
			Until: rs.Primary.Attributes["end"],
		}
		overrides, _, err := client.Schedules.ListOverrides(
			rs.Primary.Attributes["schedule"],
			&listOverridesOptions,
		)
		if err != nil {
			return err
		}

		if len(overrides.Overrides) != 1 {
			return fmt.Errorf("Expected to find a single override; found %d", len(overrides.Overrides))
		}

		if overrides.Overrides[0].ID != rs.Primary.ID {
			return fmt.Errorf("Override not found: %v - %v", rs.Primary.ID, overrides.Overrides[0].ID)
		}

		return nil
	}
}

func testAccCheckPagerDutyScheduleBaseConfig(
	usernameA,
	emailA,
	usernameB,
	emailB,
	scheduleName string,
) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "userA" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_user" "userB" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_schedule" "foo" {
  name      = "%s"
  time_zone = "America/New_York"

  layer {
    name                         = "Night Shift"
    start                        = "2019-04-01T20:00:00-04:00"
    rotation_virtual_start       = "2019-04-07T20:00:00-04:00"
    rotation_turn_length_seconds = 86400
    users                        = ["${pagerduty_user.userA.id}"]
  }
}
`, usernameA, emailA, usernameB, emailB, scheduleName)
}

func testAccCheckPagerDutyOverrideSimpleConfig(
	usernameA,
	emailA,
	usernameB,
	emailB,
	scheduleName,
	overrideStart,
	overrideEnd string,
) string {
	return testAccCheckPagerDutyScheduleBaseConfig(
		usernameA, emailA, usernameB, emailB, scheduleName,
	) +
		fmt.Sprintf(`
resource "pagerduty_schedule_override" "foo_override" {
  user     = "${pagerduty_user.userB.id}"
  start    = "%s"
  end      = "%s"
  schedule = "${pagerduty_schedule.foo.id}"
}
`, overrideStart, overrideEnd)
}

// TODO:
// create a schedule, observe who is on call, create an override, observe again
// test that no overrides exist in a given window
// test creation of two overrides within a given window
// test the overlapping overrides scenario
// handle situation when we try to delete an override from the past
