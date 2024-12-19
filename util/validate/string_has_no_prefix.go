package validate

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringHasNoPrefix struct {
	prefixes []string
}

var _ validator.String = (*stringHasNoPrefix)(nil)

func (v *stringHasNoPrefix) Description(context.Context) string {
	list := strings.Join(v.prefixes, ", ")
	return "Validates string does not start with any of these: " + list
}

func (v *stringHasNoPrefix) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *stringHasNoPrefix) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, prefix := range v.prefixes {
		if strings.HasPrefix(req.ConfigValue.ValueString(), prefix) {
			res.Diagnostics.AddError("Invalid Value", "string should not have prefix "+prefix)
		}
	}
}

func StringHasNoPrefix(prefixes ...string) validator.String {
	return &stringHasNoPrefix{prefixes: prefixes}
}
