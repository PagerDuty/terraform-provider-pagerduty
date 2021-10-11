package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPagerDutyTag_import(t *testing.T) {
	tag := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyTagConfig(tag),
			},

			{
				ResourceName:      "pagerduty_tag.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
