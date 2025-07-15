package planmodify

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type UseStateForUnknownIfFunc func(context.Context, planmodifier.StringRequest) bool

// UseStateForUnknownIf returns a plan modifier that copies a known prior state
// value into the planned value. Use this when it is known that an unconfigured
// value will remain the same after a resource update, except for special cases
// where state gets invalidated by some condition.
func UseStateForUnknownIf(fn UseStateForUnknownIfFunc) planmodifier.String {
	return useStateForUnknownIfModifier{fn: fn}
}

// useStateForUnknownIfModifier implements the plan modifier.
type useStateForUnknownIfModifier struct {
	fn UseStateForUnknownIfFunc
}

// Description returns a human-readable description of the plan modifier.
func (m useStateForUnknownIfModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change, while not invalidated."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useStateForUnknownIfModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change, while not invalidated."
}

// PlanModifyString implements the plan modification logic.
func (m useStateForUnknownIfModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	if m.fn(ctx, req) {
		resp.PlanValue = req.StateValue
	}
}
