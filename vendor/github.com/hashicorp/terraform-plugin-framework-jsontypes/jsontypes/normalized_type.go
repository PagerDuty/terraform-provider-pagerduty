// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*NormalizedType)(nil)
	_ xattr.TypeWithValidate  = (*NormalizedType)(nil)
)

// NormalizedType is an attribute type that represents a valid JSON string (RFC 7159). Semantic equality logic is defined for NormalizedType
// such that inconsequential differences between JSON strings are ignored (whitespace, property order, etc). If you need strict, byte-for-byte,
// string equality, consider using ExactType.
type NormalizedType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t NormalizedType) String() string {
	return "jsontypes.NormalizedType"
}

// ValueType returns the Value type.
func (t NormalizedType) ValueType(ctx context.Context) attr.Value {
	return Normalized{}
}

// Equal returns true if the given type is equivalent.
func (t NormalizedType) Equal(o attr.Type) bool {
	other, ok := o.(NormalizedType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// Validate implements type validation. This type requires the value provided to be a String value that is valid JSON format (RFC 7159).
func (t NormalizedType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"JSON Normalized Type Validation Error",
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
			"JSON Normalized Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)

		return diags
	}

	if ok := json.Valid([]byte(valueString)); !ok {
		diags.AddAttributeError(
			path,
			"Invalid JSON String Value",
			"A string value was provided that is not valid JSON string format (RFC 7159).\n\n"+
				"Given Value: "+valueString+"\n",
		)

		return diags
	}

	return diags
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t NormalizedType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return Normalized{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t NormalizedType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
