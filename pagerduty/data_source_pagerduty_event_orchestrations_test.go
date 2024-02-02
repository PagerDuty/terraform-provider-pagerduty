package pagerduty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourcePagerDutyEventOrchestrations_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	multipleMatchesName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	notMatchingName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePagerDutyEventOrchestrations("pagerduty_event_orchestration.test", "data.pagerduty_event_orchestrations.by_name"),
				),
			},
			{
				Config: testAccDataSourcePagerDutyEventOrchestrationsMultipleMatchesConfig(multipleMatchesName, notMatchingName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.pagerduty_event_orchestrations.by_name", "event_orchestrations.0.name", fmt.Sprintf("%s-matching-eo-name1", multipleMatchesName)),
					resource.TestCheckResourceAttr(
						"data.pagerduty_event_orchestrations.by_name", "event_orchestrations.1.name", fmt.Sprintf("%s-matching-eo-name2", multipleMatchesName)),
					resource.TestCheckNoResourceAttr(
						"data.pagerduty_event_orchestrations.by_name", "event_orchestrations.2"),
				),
			},
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationsNotFoundConfig(name),
				ExpectError: regexp.MustCompile("Unable to locate any Event Orchestration matching the expression"),
			},
			{
				Config:      testAccDataSourcePagerDutyEventOrchestrationsInvalidRegexConfig(),
				ExpectError: regexp.MustCompile("invalid regexp for name_filter provided"),
			},
		},
	})
}

func testAccDataSourcePagerDutyEventOrchestrations(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get an Event Orchestration ID from PagerDuty")
		}

		testAtts := []string{"id", "name", "integration"}

		for _, att := range testAtts {
			sub_att := fmt.Sprintf("event_orchestrations.0.%s", att)
			if a[sub_att] != srcA[att] {
				return fmt.Errorf("Expected the Event Orchestration %s to be: %s, but got: %s", att, srcA[att], a[sub_att])
			}
		}

		return nil
	}
}

func testAccDataSourcePagerDutyEventOrchestrationsConfig(name string) string {
	return fmt.Sprintf(`
resource "pagerduty_event_orchestration" "test" {
  name                    = "%s"
}

data "pagerduty_event_orchestrations" "by_name" {
  name_filter = pagerduty_event_orchestration.test.name
}
`, name)
}

func testAccDataSourcePagerDutyEventOrchestrationsMultipleMatchesConfig(matchingName, notMatchingName string) string {
	return fmt.Sprintf(`
resource "pagerduty_event_orchestration" "test1" {
  name                    = "%[1]s-matching-eo-name1"
}
resource "pagerduty_event_orchestration" "test2" {
  # this explicit dependecy is introduced to ensure the order of EO on the Data Source, because the test check relies on this order
  depends_on = [
    pagerduty_event_orchestration.test1,
  ]

  name                    = "%[1]s-matching-eo-name2"
}
resource "pagerduty_event_orchestration" "test3" {
  name                    = "%[2]s"
}

data "pagerduty_event_orchestrations" "by_name" {
  depends_on = [
    pagerduty_event_orchestration.test1,
    pagerduty_event_orchestration.test2,
    pagerduty_event_orchestration.test3,
  ]

  name_filter = "^%[1]s*"
}
`, matchingName, notMatchingName)
}

func testAccDataSourcePagerDutyEventOrchestrationsNotFoundConfig(name string) string {
	return fmt.Sprintf(`
data "pagerduty_event_orchestrations" "not_found" {
  name_filter = %q
}
`, name)
}

func testAccDataSourcePagerDutyEventOrchestrationsInvalidRegexConfig() string {
	return fmt.Sprintf(`
data "pagerduty_event_orchestrations" "invalid_regex" {
  name_filter = ")"
}
`)
}
