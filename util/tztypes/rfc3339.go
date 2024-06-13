package tztypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type RFC3339Value struct {
	basetypes.StringValue
}

func NewRFC3339Null() RFC3339Value {
	return RFC3339Value{StringValue: basetypes.NewStringNull()}
}

func NewRFC3339Value(v string) RFC3339Value {
	return RFC3339Value{StringValue: basetypes.NewStringValue(v)}
}

func (s RFC3339Value) Type(_ context.Context) attr.Type {
	return RFC3339Type{}
}

type RFC3339Type struct {
	basetypes.StringType
}

func (t RFC3339Type) String() string {
	return "tztypes.RFC3339Type"
}

func (t RFC3339Type) Equal(o attr.Type) bool {
	_, ok := o.(RFC3339Type)
	if ok {
		return true
	}

	return t.StringType.Equal(o)
}

func (t RFC3339Type) Validate(ctx context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
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

	if _, err := time.Parse(time.RFC3339, valueString); err != nil {
		diags.AddAttributeError(
			path,
			"Invalid String Value",
			"A string value was provided that is not a valid RFC3339 time.\n"+
				"Given Value: "+valueString,
		)
		return
	}

	return
}
