package pagerduty

// DEPRECATED
// Please don't add functions to this file and use the 'util' module instead.

import (
	"time"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func timeToUTC(v string) (time.Time, error) {
	return util.TimeToUTC(v)
}

func validateRFC3339(v interface{}, k string) ([]string, []error) {
	return util.ValidateRFC3339(v, k)
}

func genErrorTimeFormatRFC339(value, k string) error {
	return util.GenErrorTimeFormatRFC339(value, k)
}

func suppressRFC3339Diff(k, oldTime, newTime string, d *schema.ResourceData) bool {
	return util.SuppressRFC3339Diff(k, oldTime, newTime, d)
}

func suppressScheduleLayerStartDiff(k, oldTime, newTime string, d *schema.ResourceData) bool {
	return util.SuppressScheduleLayerStartDiff(k, oldTime, newTime, d)
}

// func parseRFC3339Time(k, oldTime, newTime string) (time.Time, time.Time, error) {
// 	return util.ParseRFC3339Time(k, oldTime, newTime)
// }

func suppressLeadTrailSpaceDiff(k, old, new string, d *schema.ResourceData) bool {
	return util.SuppressLeadTrailSpaceDiff(k, old, new, d)
}

func suppressCaseDiff(k, old, new string, d *schema.ResourceData) bool {
	return util.SuppressCaseDiff(k, old, new, d)
}

func validateValueDiagFunc(values []string) schema.SchemaValidateDiagFunc {
	return util.ValidateValueDiagFunc(values)
}

type StringContentValidationMode = util.StringContentValidationMode

const (
	NoContentValidation               = util.NoContentValidation
	NoNonPrintableChars               = util.NoNonPrintableChars
	NoNonPrintableCharsOrSpecialChars = util.NoNonPrintableCharsOrSpecialChars
)

func validateIsAllowedString(mode StringContentValidationMode) schema.SchemaValidateDiagFunc {
	return util.ValidateIsAllowedString(mode)
}

func expandStringList(configured []interface{}) []string {
	return util.ExpandStringList(configured)
}

func expandString(v string) []interface{} {
	return util.ExpandString(v)
}

func flattenSlice(v []interface{}) interface{} {
	return util.FlattenSlice(v)
}

func stringTypeToStringPtr(v string) *string {
	return util.StringTypeToStringPtr(v)
}

func stringPtrToStringType(v *string) string {
	return util.StringPtrToStringType(v)
}

func intTypeToIntPtr(v int) *int {
	return util.IntTypeToIntPtr(v)
}

func renderRoundedPercentage(p float64) string {
	return util.RenderRoundedPercentage(p)
}

func isNilFunc(i interface{}) bool {
	return util.IsNilFunc(i)
}

func unique(s []string) []string {
	return util.Unique(s)
}

func resourcePagerDutyParseColonCompoundID(id string) (string, string, error) {
	return util.ResourcePagerDutyParseColonCompoundID(id)
}
