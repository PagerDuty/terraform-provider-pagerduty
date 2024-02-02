package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyEventOrchestration_import(t *testing.T) {
	name := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-description-%s", acctest.RandString(5))
	team1 := fmt.Sprintf("tf-team1-%s", acctest.RandString(5))
	team2 := fmt.Sprintf("tf-team2-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationConfig(name, description, team1, team2),
			},
			{
				ResourceName:      "pagerduty_event_orchestration.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyEventOrchestrationNameOnly_import(t *testing.T) {
	name := fmt.Sprintf("tf-name-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationConfigNameOnly(name),
			},

			{
				ResourceName:      "pagerduty_event_orchestration.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
