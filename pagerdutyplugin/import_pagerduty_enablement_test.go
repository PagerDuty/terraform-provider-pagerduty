package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyEnablement_Service_Import(t *testing.T) {
	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)
	escalationPolicy := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementServiceConfig(serviceName, username, email, escalationPolicy, "false"),
			},
			{
				ResourceName:      "pagerduty_enablement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPagerDutyEnablement_EventOrchestration_Import(t *testing.T) {
	orchestrationName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyEnablementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyEnablementEventOrchestrationConfig(orchestrationName, "false"),
			},
			{
				ResourceName:      "pagerduty_enablement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
