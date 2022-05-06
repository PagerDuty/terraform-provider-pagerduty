package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyEventOrchestration_import(t *testing.T) {
	name := fmt.Sprintf("tf-name-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-description-%s", acctest.RandString(5))
	teamName := fmt.Sprintf("tf-team-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyEventOrchestrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEventOrchestrationConfig(name, description, teamName),
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
				ResourceName:      "pagerduty_event_orchestration.nameonly",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
