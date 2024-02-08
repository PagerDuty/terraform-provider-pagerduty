package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPagerDutyUser_import(t *testing.T) {
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.test", username)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPagerDutyUserConfig(username, email),
			},

			{
				ResourceName:      "pagerduty_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
