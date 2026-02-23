package validate

import (
	"context"
	"fmt"
	"sort"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type validTimeZone struct{}

var _ validator.String = (*validTimeZone)(nil)

func (v *validTimeZone) Description(context.Context) string {
	return "Validates that the value is a supported IANA time zone."
}

func (v *validTimeZone) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *validTimeZone) ValidateString(_ context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	foundAt := sort.SearchStrings(util.ValidTZ, value)
	if foundAt >= len(util.ValidTZ) || util.ValidTZ[foundAt] != value {
		res.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Time Zone",
			fmt.Sprintf("%q is not a valid time zone. Please refer to the list of allowed values at https://developer.pagerduty.com/docs/1afe25e9c94cb-types#time-zone", value),
		)
	}
}

// ValidTimeZone returns a Framework validator that checks the value is a
// supported IANA time zone accepted by the PagerDuty API.
func ValidTimeZone() validator.String {
	return &validTimeZone{}
}
