package pagerduty

import (
	"context"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The v3 API echoes back the time zone it was given, so flattening must keep it
// verbatim; rewriting it produces a perpetual plan diff against the config.
func TestFlattenScheduleV3TimeZone(t *testing.T) {
	for _, tz := range []string{"Etc/UTC", "America/New_York"} {
		var diags diag.Diagnostics
		schedule := &pagerduty.ScheduleV3{ID: "PSCHED1", Name: "test", TimeZone: tz}

		model := flattenScheduleV3(context.Background(), schedule, nil, types.ListNull(types.StringType), &diags)

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got := model.TimeZone.ValueString(); got != tz {
			t.Errorf("time_zone: want %q; got %q", tz, got)
		}
	}
}
