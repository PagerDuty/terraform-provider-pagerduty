package tztypes

import (
	"context"
	"fmt"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type StringValue struct {
	basetypes.StringValue
}

func NewStringNull() StringValue {
	return StringValue{StringValue: basetypes.NewStringNull()}
}

func NewStringValue(v string) StringValue {
	return StringValue{StringValue: basetypes.NewStringValue(v)}
}

func (s StringValue) Type(_ context.Context) attr.Type {
	return StringType{}
}

type StringType struct {
	basetypes.StringType
}

func (t StringType) String() string {
	return "tztypes.StringType"
}

func (t StringType) Equal(o attr.Type) bool {
	_, ok := o.(StringType)
	if ok {
		return true
	}

	return t.StringType.Equal(o)
}

func (t StringType) Validate(ctx context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var valueString string
	if err := in.As(&valueString); err != nil {
		diags.AddAttributeError(
			path,
			"Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	if !util.IsValidTZ(valueString) {
		diags.AddAttributeError(
			path,
			"Invalid String Value",
			"A string value was provided that is not a valid timezone.\n"+
				"Given Value: "+valueString,
		)
		return
	}

	return
}
