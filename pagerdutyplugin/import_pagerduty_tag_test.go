package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyTag_import(t *testing.T) {
	tag := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPagerDutyTagDestroy,
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
