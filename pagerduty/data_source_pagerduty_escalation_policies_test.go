package pagerduty

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePagerDutyEscalationPolicies_Basic(t *testing.T) {
	dataSourceName := "data.pagerduty_escalation_policies.all"
	username := fmt.Sprintf("tf-%s", acctest.RandString(5))
	email := fmt.Sprintf("%s@foo.com", username)
	escalationPolicy1 := fmt.Sprintf("tf-%s1", acctest.RandString(5))
	escalationPolicy2 := fmt.Sprintf("tf-%s2", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePagerDutyEscalationPoliciesConfig(username, email, escalationPolicy1, escalationPolicy2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "names.#", "2"),
				),
			},
		},
	})
}

func testAccDataSourcePagerDutyEscalationPoliciesConfig(username, email, escalationPolicy1 string, escalationPolicy2 string) string {
	return fmt.Sprintf(`
resource "pagerduty_user" "test" {
  name  = "%s"
  email = "%s"
}

resource "pagerduty_escalation_policy" "test1" {
  name        = "%s"
  num_loops   = 2

  rule {
    escalation_delay_in_minutes = 10

    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

resource "pagerduty_escalation_policy" "test2" {
  name        = "%s"
  num_loops   = 2

  rule {
    escalation_delay_in_minutes = 15

    target {
      type = "user_reference"
      id   = pagerduty_user.test.id
    }
  }
}

data "pagerduty_escalation_policies" "all" {}
`, username, email, escalationPolicy1, escalationPolicy2)
}
