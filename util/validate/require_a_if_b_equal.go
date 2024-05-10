package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// RequireAIfBEqual checks path `a` is not null when path `b` is equal to `expected`.
func RequireAIfBEqual(a, b path.Path, expected attr.Value) resource.ConfigValidator {
	return &requireIfEqual{
		dst: a,
		src: b,
		exp: expected,
	}
}

// RequireAIfBEqualWithMessage checks path `a` is not null when path `b` is
// equal to `expected`. Raises error message `msg` when `a` is null.
func RequireAIfBEqualWithMessage(a, b path.Path, expected attr.Value, msg string) resource.ConfigValidator {
	return &requireIfEqual{
		dst: a,
		src: b,
		exp: expected,
		msg: msg,
	}
}

type requireIfEqual struct {
	dst path.Path
	src path.Path
	exp attr.Value
	msg string
}

func (v *requireIfEqual) Description(ctx context.Context) string         { return "" }
func (v *requireIfEqual) MarkdownDescription(ctx context.Context) string { return "" }

func (v *requireIfEqual) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
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
		if dst.IsNull() {
			detail := v.msg
			if detail == "" {
				detail = fmt.Sprintf("When the value of %s equals %s, field %s must have an explicit value", v.src, v.exp, v.dst)
			}
			resp.Diagnostics.AddAttributeError(v.dst, fmt.Sprintf("Required %s", v.dst), detail)
			return
		}
	}
}
