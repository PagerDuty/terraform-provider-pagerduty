package optiontypes

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ListValue struct {
	basetypes.ListValue
}

func (l ListValue) Equal(o attr.Value) bool {
	other, ok := o.(ListValue)
	if !ok {
		if l.IsNull() || other.IsNull() {
			return true
		}
		o = other.ListValue
	}

	return l.ListValue.Equal(o)
}

func NewListNull(t attr.Type) ListValue {
	return ListValue{basetypes.NewListNull(t)}
}

func NewListUnknown(t attr.Type) ListValue {
	return ListValue{basetypes.NewListUnknown(t)}
}

func NewListValue(t attr.Type, elements []attr.Value) (ListValue, diag.Diagnostics) {
	l, diags := basetypes.NewListValue(t, elements)
	return ListValue{l}, diags
}

type ListType struct {
	basetypes.ListType
}

func (l ListType) Equal(o attr.Type) bool {
	if _, ok := o.(ListType); ok {
		return true
	}
	_, ok := o.(basetypes.ListType)
	return ok
}

func (l ListType) String() string {
	return "optiontypes.ListType[" + l.ElementType().String() + "]"
}
