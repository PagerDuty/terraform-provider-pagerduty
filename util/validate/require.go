package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Require checks a path is not null.
func Require(p path.Path) resource.ConfigValidator {
	return &requirePath{Path: p}
}

type requirePath struct {
	path.Path
}

func (v *requirePath) Description(ctx context.Context) string {
	return "Forces item to be present if its parent is present"
}

func (v *requirePath) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *requirePath) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var parent attr.Value
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, v.Path.ParentPath(), &parent)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if parent.IsNull() {
		return
	}

	var src attr.Value
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, v.Path, &src)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if src.IsNull() {
		resp.Diagnostics.AddAttributeError(
			v.Path,
			fmt.Sprintf("Required %s", v.Path),
			fmt.Sprintf("Field %s must have an explicit value", v.Path),
		)
		return
	}
}
