// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*Normalized)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*Normalized)(nil)
)

// Normalized represents a valid JSON string (RFC 7159). Semantic equality logic is defined for Normalized
// such that inconsequential differences between JSON strings are ignored (whitespace, property order, etc). If you
// need strict, byte-for-byte, string equality, consider using ExactType.
type Normalized struct {
	basetypes.StringValue
}

// Type returns a NormalizedType.
func (v Normalized) Type(_ context.Context) attr.Type {
	return NormalizedType{}
}

// Equal returns true if the given value is equivalent.
func (v Normalized) Equal(o attr.Value) bool {
	other, ok := o.(Normalized)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true if the given JSON string value is semantically equal to the current JSON string value. When compared,
// these JSON string values are "normalized" by marshalling them to empty Go structs. This prevents Terraform data consistency errors and
// resource drift due to inconsequential differences in the JSON strings (whitespace, property order, etc).
func (v Normalized) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(Normalized)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	result, err := jsonEqual(newValue.ValueString(), v.ValueString())

	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected error occurred while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return false, diags
	}

	return result, diags
}

func jsonEqual(s1, s2 string) (bool, error) {
	s1, err := normalizeJSONString(s1)
	if err != nil {
		return false, err
	}

	s2, err = normalizeJSONString(s2)
	if err != nil {
		return false, err
	}

	return s1 == s2, nil
}

func normalizeJSONString(jsonStr string) (string, error) {
	dec := json.NewDecoder(strings.NewReader(jsonStr))

	// This ensures the JSON decoder will not parse JSON numbers into Go's float64 type; avoiding Go
	// normalizing the JSON number representation or imposing limits on numeric range. See the unit test cases
	// of StringSemanticEquals for examples.
	dec.UseNumber()

	var temp interface{}
	if err := dec.Decode(&temp); err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(&temp)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// Unmarshal calls (encoding/json).Unmarshal with the Normalized StringValue and `target` input. A null or unknown value will produce an error diagnostic.
// See encoding/json docs for more on usage: https://pkg.go.dev/encoding/json#Unmarshal
func (v Normalized) Unmarshal(target any) diag.Diagnostics {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("Normalized JSON Unmarshal Error", "json string value is null"))
		return diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("Normalized JSON Unmarshal Error", "json string value is unknown"))
		return diags
	}

	err := json.Unmarshal([]byte(v.ValueString()), target)
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("Normalized JSON Unmarshal Error", err.Error()))
	}

	return diags
}

// NewNormalizedNull creates a Normalized with a null value. Determine whether the value is null via IsNull method.
func NewNormalizedNull() Normalized {
	return Normalized{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewNormalizedUnknown creates a Normalized with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewNormalizedUnknown() Normalized {
	return Normalized{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewNormalizedValue creates a Normalized with a known value. Access the value via ValueString method.
func NewNormalizedValue(value string) Normalized {
	return Normalized{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewNormalizedPointerValue creates a Normalized with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewNormalizedPointerValue(value *string) Normalized {
	return Normalized{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
