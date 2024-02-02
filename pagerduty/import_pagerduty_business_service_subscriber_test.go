package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPagerDutyBusinessServiceSubscriber_import(t *testing.T) {
	businessServiceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	team := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyBusinessServiceSubscriberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceSubscriberTeamConfig(businessServiceName, team),
			},
			{
				ResourceName:      "pagerduty_business_service_subscriber.foo",
				ImportStateIdFunc: testAccCheckPagerDutyBusinessServiceSubscriberID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPagerDutyBusinessServiceSubscriberID(s *terraform.State) (string, error) {
	return fmt.Sprintf("%v.%v.%v", s.RootModule().Resources["pagerduty_business_service.foo"].Primary.ID, "team", s.RootModule().Resources["pagerduty_team.foo"].Primary.ID), nil
}
