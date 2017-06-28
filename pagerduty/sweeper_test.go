package pagerduty

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/heimweh/go-pagerduty/pagerduty"
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
		Name: "pagerdutypagerduty_team_user",
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

func sweeperClient() (*pagerduty.Client, error) {
	if os.Getenv("PAGERDUTY_SWEEPER_TOKEN") == "" {
		return nil, fmt.Errorf("$PAGERDUTY_SWEEPER_TOKEN must be set")
	}

	config := &pagerduty.Config{
		Token: os.Getenv("PAGERDUTY_SWEEPER_TOKEN"),
		Debug: true,
	}

	return pagerduty.NewClient(config)
}
