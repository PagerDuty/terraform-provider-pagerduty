package rangetypes

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Int64Value struct {
	basetypes.Int64Value
	RangeType Int64Type
}

func NewInt64Null(t Int64Type) Int64Value {
	return Int64Value{Int64Value: basetypes.NewInt64Null(), RangeType: t}
}

func NewInt64Value(v int64, t Int64Type) Int64Value {
	return Int64Value{Int64Value: basetypes.NewInt64Value(v), RangeType: t}
}

func (s Int64Value) Type(_ context.Context) attr.Type {
	return s.RangeType
}

type Int64Type struct {
	basetypes.Int64Type
	Start int64
	End   int64
}

func (t Int64Type) String() string {
	return "rangetypes.Int64Type"
}

func (t Int64Type) Equal(o attr.Type) bool {
	if t2, ok := o.(Int64Type); ok {
		return t.Start == t2.Start && t.End == t2.End
	}
	return t.Int64Type.Equal(o)
}

func (t Int64Type) addTypeValidationError(err error, path path.Path, diags *diag.Diagnostics) {
	diags.AddAttributeError(
		path,
		"Type Validation Error",
		"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
			"Please report the following to the provider developer:\n\n"+err.Error(),
	)
}

func (t Int64Type) Validate(ctx context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return
	}

	if !in.Type().Is(tftypes.Number) {
		err := fmt.Errorf("expected Int64 value, received %T with value: %v", in, in)
		t.addTypeValidationError(err, path, &diags)
		return
	}

	if !in.IsKnown() || in.IsNull() {
		return
	}

	var valueFloat big.Float
	if err := in.As(&valueFloat); err != nil {
		t.addTypeValidationError(err, path, &diags)
		return
	}
	valueInt64, _ := valueFloat.Int64()

	if valueInt64 < t.Start || valueInt64 > int64(t.End) {
		diags.AddAttributeError(
			path,
			"Invalid Int64 Value",
			fmt.Sprintf("A value was provided that is not inside valid range (%v, %v).\n"+
				"Given Value: %v", t.Start, t.End, valueInt64),
		)
		return
	}

	return
}
