package pagerduty

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcePagerDutyVendor_Basic(t *testing.T) {
	dataSourceName := "data.pagerduty_vendor.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyVendorConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "PZQ6AUS"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Amazon CloudWatch"),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyVendor_ExactMatch(t *testing.T) {
	dataSourceName := "data.pagerduty_vendor.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPagerDutyScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyExactMatchConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "PKG4M95"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "Sentry"),
				),
			},
		},
	})
}

const testAccDataSourcePagerDutyVendorConfig = `
data "pagerduty_vendor" "foo" {
  name = "cloudwatch"
}
`

const testAccDataSourcePagerDutyExactMatchConfig = `
data "pagerduty_vendor" "foo" {
  name = "sentry"
}
`
