package validate

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringHasNoSuffix struct {
	suffixes []string
}

var _ validator.String = (*stringHasNoSuffix)(nil)

func (v *stringHasNoSuffix) Description(context.Context) string {
	list := strings.Join(v.suffixes, ", ")
	return "Validates string does not end with any of these: " + list
}

func (v *stringHasNoSuffix) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *stringHasNoSuffix) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, suffix := range v.suffixes {
		if strings.HasSuffix(req.ConfigValue.ValueString(), suffix) {
			res.Diagnostics.AddError("Invalid Value", "string should not have suffix "+suffix)
		}
	}
}

func StringHasNoSuffix(suffixes ...string) validator.String {
	return &stringHasNoSuffix{suffixes: suffixes}
}
