package validate

import (
	"context"
	"strings"
	"unicode"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type validateIsAllowedString struct {
	validateFn func(s string) bool
	util.StringDescriber
}

func (v validateIsAllowedString) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if ok := v.validateFn(req.ConfigValue.ValueString()); !ok {
		resp.Diagnostics.AddError(v.Value, "")
	}
}

func IsAllowedString(mode util.StringContentValidationMode) validator.String {
	switch mode {
	case util.NoNonPrintableChars:
		return validateIsAllowedString{
			func(s string) bool {
				for _, char := range s {
					if !unicode.IsPrint(char) {
						return false
					}
				}
				return s != "" && !strings.HasSuffix(s, " ")
			},
			util.StringDescriber{Value: "Name can not be blank, nor contain non-printable characters. Trailing white spaces are not allowed either."},
		}
	default:
		return validateIsAllowedString{
			func(s string) bool { return false },
			util.StringDescriber{Value: "Invalid mode while using func IsAllowedStringValidator(mode StringContentValidationMode)"},
		}
	}
}
