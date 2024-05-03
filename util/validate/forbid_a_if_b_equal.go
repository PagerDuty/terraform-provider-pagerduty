package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ForbidAIfBEqual raises an error if path `a` is not null when path `b` is
// equal to expected value `exp`.
func ForbidAIfBEqual(a, b path.Path, expected attr.Value) resource.ConfigValidator {
	return &forbidIfEqual{
		dst: a,
		src: b,
		exp: expected,
	}
}

// ForbidAIfBEqual raises an error if path `a` is not null when path `b` is
// equal to expected value `exp`. Raising message `msg` when invalid.
func ForbidAIfBEqualWithMessage(a, b path.Path, expected attr.Value, message string) resource.ConfigValidator {
	return &forbidIfEqual{
		dst: a,
		src: b,
		exp: expected,
		msg: message,
	}
}

type forbidIfEqual struct {
	dst path.Path
	src path.Path
	exp attr.Value
	msg string
}

func (v *forbidIfEqual) Description(ctx context.Context) string         { return "" }
func (v *forbidIfEqual) MarkdownDescription(ctx context.Context) string { return "" }

func (v *forbidIfEqual) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var src attr.Value
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, v.src, &src)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if src.IsNull() || src.IsUnknown() {
		return
	}

	if src.Equal(v.exp) {
		var dst attr.Value
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, v.dst, &dst)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if !dst.IsNull() {
			detail := v.msg
			if detail == "" {
				detail = fmt.Sprintf("When the value of %s equals %s, field %s cannot have a value", v.src, v.exp, v.dst)
			}
			resp.Diagnostics.AddAttributeError(
				v.dst,
				fmt.Sprintf("Forbidden %s", v.dst),
				detail,
			)
			return
		}
	}
}
