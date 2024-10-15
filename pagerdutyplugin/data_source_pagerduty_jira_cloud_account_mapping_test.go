package pagerduty

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePagerDutyJiraCloudAccountMapping_Basic(t *testing.T) {
	subdomain := os.Getenv("PAGERDUTY_SUBDOMAIN")
	if subdomain == "" {
		t.Skip("Missing env variable PAGERDUTY_SUBDOMAIN")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyJiraCloudAccountMappingConfig(subdomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pagerduty_jira_cloud_account_mapping.by_subdomain", "subdomain", subdomain),
					resource.TestCheckResourceAttrSet("data.pagerduty_jira_cloud_account_mapping.by_subdomain", "id"),
					resource.TestCheckResourceAttrSet("data.pagerduty_jira_cloud_account_mapping.by_subdomain", "base_url"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyJiraCloudAccountMappingConfig(subdomain string) string {
	return fmt.Sprintf(`
data "pagerduty_jira_cloud_account_mapping" "by_subdomain" {
  subdomain = "%s"
}
`, subdomain)
}
