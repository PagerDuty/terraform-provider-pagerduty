package planmodify

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseNullForRemovedWithState sets plan to null if the list has an state, but
// the configuration is now null
func UseNullForRemovedWithState() planmodifier.List {
	return useNullForRemovedWithStateModifier{}
}

type useNullForRemovedWithStateModifier struct{}

func (m useNullForRemovedWithStateModifier) Description(_ context.Context) string {
	return "Removes the value if the list has an state, but the configuration changes to null"
}

func (m useNullForRemovedWithStateModifier) MarkdownDescription(_ context.Context) string {
	return "Removes the value if the list has an state, but the configuration changes to null"
}

func (m useNullForRemovedWithStateModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known or an unknown configuration value.
	if !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = req.ConfigValue
}
