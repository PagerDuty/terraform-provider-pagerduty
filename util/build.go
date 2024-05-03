package util

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringToUintPointer(p path.Path, s types.String, diags *diag.Diagnostics) *uint {
	if s.IsNull() || s.IsUnknown() || s.ValueString() == "" || s.ValueString() == "null" {
		return nil
	}
	if val, err := strconv.Atoi(s.ValueString()); err == nil {
		uintvalue := uint(val)
		return &uintvalue
	} else {
		diags.AddError(fmt.Sprintf("Value for %q is not a valid number", p), err.Error())
	}
	return nil
}
