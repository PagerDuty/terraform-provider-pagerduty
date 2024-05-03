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
		dst:      a,
		src:      b,
		expected: expected,
	}
}

type requireIfEqual struct {
	dst      path.Path
	src      path.Path
	expected attr.Value
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

	if src.Equal(v.expected) {
		var dst attr.Value
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, v.dst, &dst)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if dst.IsNull() || dst.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				v.dst,
				fmt.Sprintf("Required %s", v.dst),
				fmt.Sprintf("When the value of %s equals %s, field %s must have an explicit value", v.src, v.expected, v.dst),
			)
			return
		}
	}
}
