package pagerduty

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutySchedule_import(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	schedule := fmt.Sprintf("tf-%s", acctest.RandString(5))
	location := "Europe/Berlin"
	t.Setenv("PAGERDUTY_TIME_ZONE", location)
	start := testAccTimeNow().Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)
	rotationVirtualStart := testAccTimeNow().Add(24 * time.Hour).Round(1 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyScheduleConfig(username, email, schedule, location, start, rotationVirtualStart),
			},
			{
				ResourceName:      "pagerduty_schedule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyUserDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "pagerduty_user" {
			continue
		}

		ctx := context.Background()
		if _, err := testAccProvider.client.GetUserWithContext(ctx, r.Primary.ID, pagerduty.GetUserOptions{}); err == nil {
			return fmt.Errorf("User still exists")
		}
	}
	return nil
}
