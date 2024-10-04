package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyAlertGroupingSetting_Time(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyAlertGroupingSettingTimeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyAlertGroupingSetting("pagerduty_alert_grouping_setting.test", "data.pagerduty_alert_grouping_setting.by_name"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "name", name),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "type", "time"),
					resource.TestCheckResourceAttrSet("data.pagerduty_alert_grouping_setting.by_name", "description"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "services.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourcePagerDutyAlertGroupingSetting_ContentBased(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyAlertGroupingSettingContentBasedConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyAlertGroupingSetting("pagerduty_alert_grouping_setting.test", "data.pagerduty_alert_grouping_setting.by_name"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "name", name),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "type", "content_based"),
					resource.TestCheckResourceAttrSet("data.pagerduty_alert_grouping_setting.by_name", "description"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "services.#", "1"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "config.time_window", "300"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "config.aggregate", "any"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "config.fields.#", "1"),
					resource.TestCheckResourceAttr("data.pagerduty_alert_grouping_setting.by_name", "config.fields.0", "summary"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyAlertGroupingSetting(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a alert grouping setting ID from PagerDuty")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the alert grouping setting %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyAlertGroupingSettingTimeConfig(name string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "test" {
	name = "Default"
}

resource "pagerduty_service" "test" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.test.id
}

resource "pagerduty_alert_grouping_setting" "test" {
	name = "%[1]s"
	type = "time"
	services = [pagerduty_service.test.id]
	config {}
}

data "pagerduty_alert_grouping_setting" "by_name" {
	name = pagerduty_alert_grouping_setting.test.name
}
`, name)
}

func testAccDataSourcePagerDutyAlertGroupingSettingContentBasedConfig(name string) string {
	return fmt.Sprintf(`
data "pagerduty_escalation_policy" "test" {
	name = "Default"
}

resource "pagerduty_service" "test" {
	name = "%s"
	escalation_policy = data.pagerduty_escalation_policy.test.id
}

resource "pagerduty_alert_grouping_setting" "test" {
	name = "%[1]s"
	type = "content_based"
	services = [pagerduty_service.test.id]
	config {
		aggregate = "any"
		fields = ["summary"]
	}
}

data "pagerduty_alert_grouping_setting" "by_name" {
	name = pagerduty_alert_grouping_setting.test.name
}
`, name)
}
