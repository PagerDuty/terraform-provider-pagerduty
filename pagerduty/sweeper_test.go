package pagerduty

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	resource.AddTestSweepers("pagerduty_service", &resource.Sweeper{
		Name: "pagerduty_service",
		F:    testSweepService,
	})

	resource.AddTestSweepers("pagerduty_escalation_policy", &resource.Sweeper{
		Name: "pagerduty_escalation_policy",
		F:    testSweepEscalationPolicy,
	})

	resource.AddTestSweepers("pagerduty_user", &resource.Sweeper{
		Name: "pagerduty_user",
		F:    testSweepUser,
	})

	resource.AddTestSweepers("pagerduty_team", &resource.Sweeper{
		Name: "pagerduty_team",
		F:    testSweepTeam,
	})

	resource.AddTestSweepers("pagerduty_schedule", &resource.Sweeper{
		Name: "pagerduty_schedule",
		F:    testSweepSchedule,
	})

	resource.AddTestSweepers("pagerduty_addon", &resource.Sweeper{
		Name: "pagerduty_addon",
		F:    testSweepAddon,
	})
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedConfigForRegion returns a common config setup needed for the sweeper
// functions for a given region
func sharedConfigForRegion(region string) (*Config, error) {
	if os.Getenv("PAGERDUTY_TOKEN") == "" {
		return nil, fmt.Errorf("$PAGERDUTY_TOKEN must be set")
	}

	config := &Config{
		Token: os.Getenv("PAGERDUTY_TOKEN"),
	}

	return config, nil
}
