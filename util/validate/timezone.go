package validate

import (
	"context"
	"fmt"
	"time"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type timezoneValidator struct {
	util.StringDescriber
}

func Timezone() validator.String {
	return &timezoneValidator{
		util.StringDescriber{Value: "checks time zone is supported by the machine's tzdata"},
	}
}

func (v timezoneValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() {
		return
	}
	value := req.ConfigValue.ValueString()
	_, err := time.LoadLocation(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path, fmt.Sprintf("Timezone %q is invalid", value), err.Error(),
		)
	}
}
