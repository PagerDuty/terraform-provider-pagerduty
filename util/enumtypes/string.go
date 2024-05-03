package enumtypes

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type StringValue struct {
	basetypes.StringValue
	EnumType StringType
}

func NewStringNull(t StringType) StringValue {
	return StringValue{StringValue: basetypes.NewStringNull(), EnumType: t}
}

func NewStringValue(v string, t StringType) StringValue {
	return StringValue{StringValue: basetypes.NewStringValue(v), EnumType: t}
}

func (s StringValue) Type(_ context.Context) attr.Type {
	return s.EnumType
}

type StringType struct {
	basetypes.StringType
	OneOf []string
}

func (t StringType) String() string {
	return "enumtypes.StringType"
}

func (t StringType) Equal(o attr.Type) bool {
	if t2, ok := o.(StringType); ok {
		return slices.Equal(t.OneOf, t2.OneOf)
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

	found := false
	for _, v := range t.OneOf {
		if v == valueString {
			found = true
			break
		}
	}

	if !found {
		diags.AddAttributeError(
			path,
			"Invalid String Value",
			"A string value was provided that is not valid.\n"+
				"Given Value: "+valueString+"\n"+
				"Expecting One Of: "+strings.Join(t.OneOf, ", "),
		)
		return
	}

	return
}
