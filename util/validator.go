package util

import (
	"strings"
	"unicode"

	"github.com/hashicorp/go-cty/cty"
	v2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValidateIsAllowedString will always validate if string provided is not empty,
// neither has trailing white spaces. Additionally the string content validation
// will be done based on the `mode` set.
//
//	mode: NoContentValidation | NoNonPrintableChars | NoNonPrintableCharsOrSpecialChars
func ReValidateIsAllowedString(mode StringContentValidationMode) schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) v2diag.Diagnostics {
		var diags v2diag.Diagnostics

		fillDiags := func() {
			summary := "Name can not be blank. Trailing white spaces are not allowed either."
			switch mode {
			case NoNonPrintableChars:
				summary = "Name can not be blank, nor contain non-printable characters. Trailing white spaces are not allowed either."
			case NoNonPrintableCharsOrSpecialChars:
				summary = "Name can not be blank, nor contain the characters '\\', '/', '&', '<', '>', or any non-printable characters. Trailing white spaces are not allowed either."
			}
			diags = append(diags, v2diag.Diagnostic{
				Severity:      v2diag.Error,
				Summary:       summary,
				AttributePath: p,
			})
		}

		value := v.(string)
		if value == "" {
			fillDiags()
			return diags
		}

		for _, char := range value {
			if (mode == NoNonPrintableChars || mode == NoNonPrintableCharsOrSpecialChars) && !unicode.IsPrint(char) {
				fillDiags()
				return diags
			}
			if mode == NoNonPrintableCharsOrSpecialChars {
				switch char {
				case '\\', '/', '&', '<', '>':
					fillDiags()
					return diags
				}
			}
		}

		if strings.HasSuffix(value, " ") {
			fillDiags()
			return diags
		}

		return diags
	}
}
