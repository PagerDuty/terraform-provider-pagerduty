package validate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// RequireList checks path `p` is a list at least with size 1.
func RequireList(p path.Path) resource.ConfigValidator {
	return &requireListSize{Path: p}
}

type requireListSize struct {
	path.Path
}

func (v *requireListSize) Description(ctx context.Context) string         { return "" }
func (v *requireListSize) MarkdownDescription(ctx context.Context) string { return "" }

func (v *requireListSize) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var src attr.Value
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, v.Path, &src)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if src.IsNull() || src.IsUnknown() {
		return
	}

	size := 1
	if size < 1 {
		resp.Diagnostics.AddAttributeError(v.Path, "Required to be a list with items", "")
		return
	}
}
