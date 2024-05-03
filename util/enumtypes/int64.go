package enumtypes

import (
	"context"
	"fmt"
	"math/big"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Int64Value struct {
	basetypes.Int64Value
	EnumType Int64Type
}

func NewInt64Null(t Int64Type) Int64Value {
	return Int64Value{Int64Value: basetypes.NewInt64Null(), EnumType: t}
}

func NewInt64Value(v int64, t Int64Type) Int64Value {
	return Int64Value{Int64Value: basetypes.NewInt64Value(v), EnumType: t}
}

func (s Int64Value) Type(_ context.Context) attr.Type {
	return s.EnumType
}

type Int64Type struct {
	basetypes.Int64Type
	OneOf []int64
}

func (t Int64Type) Int64() string {
	return "enumtypes.Int64Type"
}

func (t Int64Type) Equal(o attr.Type) bool {
	if t2, ok := o.(Int64Type); ok {
		return slices.Equal(t.OneOf, t2.OneOf)
	}
	return t.Int64Type.Equal(o)
}

func (t Int64Type) Validate(ctx context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return
	}

	if !in.Type().Is(tftypes.Number) {
		err := fmt.Errorf("expected Int64 value, received %T with value: %v", in, in)
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

	var valueFloat big.Float
	if err := in.As(&valueFloat); err != nil {
		diags.AddAttributeError(
			path,
			"Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}
	valueInt64, _ := valueFloat.Int64()

	found := false
	for _, v := range t.OneOf {
		if v == valueInt64 {
			found = true
			break
		}
	}

	if !found {
		diags.AddAttributeError(
			path,
			"Invalid Int64 Value",
			fmt.Sprintf(
				"A string value was provided that is not valid.\n"+
					"Given Value: %v\n"+
					"Expecting One Of: %v",
				valueInt64,
				t.OneOf,
			),
		)
		return
	}

	return
}
