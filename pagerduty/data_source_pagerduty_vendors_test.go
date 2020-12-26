package pagerduty

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePagerDutyVendors_Basic(t *testing.T) {
	dataSourceName := "data.pagerduty_vendors.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyVendorsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "names.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourcePagerDutyVendorsConfig = `
data "pagerduty_vendors" "foo" {
}
`
