package validate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func AlternativesForPath(p path.Path, alt []attr.Value) *alternativesForPathValidator {
	return &alternativesForPathValidator{Path: p, Alternatives: alt}
}

type alternativesForPathValidator struct {
	Path         path.Path
	Alternatives []attr.Value
}

var _ validator.String = (*alternativesForPathValidator)(nil)

func (v *alternativesForPathValidator) Description(_ context.Context) string { return "" }
func (v *alternativesForPathValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *alternativesForPathValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
}
