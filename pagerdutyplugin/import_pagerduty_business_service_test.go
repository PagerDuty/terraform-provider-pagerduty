package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyBusinessService_import(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	desc := fmt.Sprintf("tf-%s", acctest.RandString(5))
	poc := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyBusinessServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyBusinessServiceConfig(name, desc, poc),
			},
			{
				ResourceName:      "pagerduty_business_service.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
