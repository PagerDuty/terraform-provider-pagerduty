package optiontypes

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Int64Value struct {
	basetypes.Int64Value
}

func NewInt64Value(i int64) attr.Value { return Int64Value{basetypes.NewInt64Value(i)} }
func NewInt64Null() attr.Value         { return Int64Value{basetypes.NewInt64Null()} }
func NewInt64Unknown() attr.Value      { return Int64Value{basetypes.NewInt64Unknown()} }

func (i Int64Value) Type(ctx context.Context) attr.Type { return Int64Type{} }

func (i Int64Value) Equal(o attr.Value) bool {
	log.Printf("[cg] o=%#v i=%#v", o, i)
	if i.IsNull() || o.IsNull() {
		return true
	}
	if o2, ok := o.(Int64Value); ok {
		o = o2.Int64Value
	}
	return i.Int64Value.Equal(o)
}

type Int64Type struct {
	basetypes.Int64Type
}

func (t Int64Type) String() string { return "optiontypes.Int64Type" }

func (t Int64Type) Equal(o attr.Type) bool {
	if _, ok := o.(Int64Type); ok {
		return true
	}
	_, ok := o.(basetypes.Int64Type)
	return ok
}
