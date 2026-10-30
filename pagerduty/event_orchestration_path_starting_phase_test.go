package pagerduty

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestEventOrchestrationPathStartingPhaseExpandFlatten(t *testing.T) {
	type pathFuncs struct {
		name    string
		schema  map[string]*schema.Schema
		parent  map[string]interface{}
		expand  func(interface{}) *pagerduty.EventOrchestrationPathRuleActions
		flatten func(*pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{}
	}

	paths := []pathFuncs{
		{
			name:    "global",
			schema:  resourcePagerDutyEventOrchestrationPathGlobal().Schema,
			parent:  map[string]interface{}{"event_orchestration": "E-ORC-1"},
			expand:  expandGlobalPathActions,
			flatten: flattenGlobalPathActions,
		},
		{
			name:    "service",
			schema:  resourcePagerDutyEventOrchestrationPathService().Schema,
			parent:  map[string]interface{}{"service": "PSVC123"},
			expand:  expandServicePathActions,
			flatten: flattenServicePathActions,
		},
	}

	for _, p := range paths {
		for _, phase := range []string{"triage", "active", ""} {
			t.Run(p.name+"/"+phaseLabel(phase), func(t *testing.T) {
				actions := map[string]interface{}{}
				if phase != "" {
					actions["starting_phase"] = phase
				}
				raw := map[string]interface{}{
					"set": []interface{}{
						map[string]interface{}{
							"id": "start",
							"rule": []interface{}{
								map[string]interface{}{"actions": []interface{}{actions}},
							},
						},
					},
					"catch_all": []interface{}{
						map[string]interface{}{"actions": []interface{}{actions}},
					},
				}
				for k, v := range p.parent {
					raw[k] = v
				}
				d := schema.TestResourceDataRaw(t, p.schema, raw)

				for _, attr := range []string{"set.0.rule.0.actions", "catch_all.0.actions"} {
					expanded := p.expand(d.Get(attr))
					if phase == "" {
						if expanded.StartingPhase != nil {
							t.Errorf("%s: expected StartingPhase to be nil when unset, got %q", attr, *expanded.StartingPhase)
						}
					} else if expanded.StartingPhase == nil || *expanded.StartingPhase != phase {
						t.Errorf("%s: expected StartingPhase %q, got %v", attr, phase, expanded.StartingPhase)
					}

					flattened := p.flatten(expanded)
					if got := flattened[0]["starting_phase"]; got != phase {
						t.Errorf("%s: expected flattened starting_phase %q, got %q", attr, phase, got)
					}
				}
			})
		}
	}
}

func TestEventOrchestrationPathStartingPhaseValidation(t *testing.T) {
	validate := validateEventOrchestrationPathStartingPhase()
	for _, v := range []string{"triage", "active"} {
		if diags := validate(v, nil); diags.HasError() {
			t.Errorf("expected %q to be valid, got %v", v, diags)
		}
	}
	for _, v := range []string{"", "Triage", "acknowledged", "resolved"} {
		if diags := validate(v, nil); !diags.HasError() {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}

func phaseLabel(p string) string {
	if p == "" {
		return "unset"
	}
	return p
}
