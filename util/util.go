package util

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TimeToUTC(v string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}, err
	}

	return t.UTC(), nil
}

// TimeNowInLoc returns the current time in the given location.
// If an error occurs when trying to load the location, we just return the
// current local time.
func TimeNowInLoc(name string) time.Time {
	loc, err := time.LoadLocation(name)
	now := time.Now()
	if err != nil {
		log.Printf("[WARN] Failed to load location: %s", err)
		return now
	}
	return now.In(loc)
}

// ValidateRFC3339 validates that a date string has the correct RFC3339 layout
func ValidateRFC3339(v interface{}, k string) (we []string, errors []error) {
	value := v.(string)
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		errors = append(errors, GenErrorTimeFormatRFC339(value, k))
	}
	if t.Second() > 0 {
		errors = append(errors, fmt.Errorf("please set the time %s to a full minute, e.g. 11:23:00, not 11:23:05", value))
	}

	return
}

func GenErrorTimeFormatRFC339(value, k string) error {
	return fmt.Errorf("%s is not a valid format for argument: %s. Expected format: %s (RFC3339)", value, k, time.RFC3339)
}

func SuppressRFC3339Diff(k, oldTime, newTime string, d *schema.ResourceData) bool {
	oldT, newT, err := ParseRFC3339Time(k, oldTime, newTime)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	return oldT.Equal(newT)
}

// issue: https://github.com/PagerDuty/terraform-provider-pagerduty/issues/200
// The start value of schedule layer can't be set to a time in the past. So if
// the value passed in is before the current time then PagerDuty will set the
// start to the current time. Thus, we do not need to show diff if both newT and
// oldT is in the past, as it will not bring
// any real changes to the schedule layer.
func SuppressScheduleLayerStartDiff(k, oldTime, newTime string, d *schema.ResourceData) bool {
	oldT, newT, err := ParseRFC3339Time(k, oldTime, newTime)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	return oldT.Equal(newT) || (newT.Before(time.Now()) && oldT.Before(time.Now()))
}

func ParseRFC3339Time(k, oldTime, newTime string) (time.Time, time.Time, error) {
	var t time.Time
	oldT, err := time.Parse(time.RFC3339, oldTime)
	if err != nil {
		return t, t, fmt.Errorf("[ERROR] Failed to parse %q (old %q). Expected format: %s (RFC3339)", oldTime, k, time.RFC3339)
	}

	newT, err := time.Parse(time.RFC3339, newTime)
	if err != nil {
		return t, t, fmt.Errorf("[ERROR] Failed to parse %q (new %q). Expected format: %s (RFC3339)", oldTime, k, time.RFC3339)
	}

	return oldT, newT, nil
}

func SuppressLeadTrailSpaceDiff(k, prev, next string, d *schema.ResourceData) bool {
	trimmedInput := strings.TrimSpace(next)
	repeatedSpaceMatcher := regexp.MustCompile(`\s+`)
	return prev == repeatedSpaceMatcher.ReplaceAllLiteralString(trimmedInput, " ")
}

func SuppressCaseDiff(k, prev, next string, d *schema.ResourceData) bool {
	return prev == strings.ToLower(next)
}

// Validate a value against a set of possible values
func ValidateValueDiagFunc(values []string) schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value := v.(string)
		valid := false
		for _, val := range values {
			if value == val {
				valid = true
				break
			}
		}

		if !valid {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("%#v is an invalid value. Must be one of %#v", value, values),
				AttributePath: p,
			})
		}
		return diags
	}
}

type StringContentValidationMode int64

const (
	NoContentValidation StringContentValidationMode = iota
	NoNonPrintableChars
	NoNonPrintableCharsOrSpecialChars
)

// ValidateIsAllowedString will always validate if string provided is not empty,
// neither has trailing white spaces. Additionally the string content validation
// will be done based on the `mode` set.
//
//	mode: NoContentValidation | NoNonPrintableChars | NoNonPrintableCharsOrSpecialChars
func ValidateIsAllowedString(mode StringContentValidationMode) schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		fillDiags := func() {
			summary := "Name can not be blank. Trailing white spaces are not allowed either."
			switch mode {
			case NoNonPrintableChars:
				summary = "Name can not be blank, nor contain non-printable characters. Trailing white spaces are not allowed either."
			case NoNonPrintableCharsOrSpecialChars:
				summary = "Name can not be blank, nor contain the characters '\\', '/', '&', '<', '>', or any non-printable characters. Trailing white spaces are not allowed either."
			}
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
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

// ExpandStringList takes the result of flatmap.Expand for an array of strings
// and returns a []string
func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, string(v.(string)))
	}
	return vs
}

func ExpandString(v string) []interface{} {
	var obj []interface{}
	if err := json.Unmarshal([]byte(v), &obj); err != nil {
		log.Printf("[ERROR] Could not unmarshal field %s: %v", v, err)
		return nil
	}

	return obj
}

func FlattenSlice(v []interface{}) interface{} {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("[ERROR] Could not marshal field %s: %v", v, err)
		return nil
	}
	return string(b)
}

// StringTypeToStringPtr is a helper that returns a pointer to
// the string value passed in or nil if the string is empty.
func StringTypeToStringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// StringPtrToStringType is a helper that returns the string value passed in
// or an empty string if the given pointer is nil.
func StringPtrToStringType(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func IntTypeToIntPtr(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

// RenderRoundedPercentage is a helper function to render percentanges
// represented as float64 numbers, by its round with two decimals string
// representation.
func RenderRoundedPercentage(p float64) string {
	return fmt.Sprintf("%.2f", math.Round(p*100))
}

// IsNilFunc is a helper which verifies if an empty interface expecting a
// nullable value indeed has a `nil` type assigned or it's just empty.
func IsNilFunc(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// Unique will remove duplicates from a strings slice.
func Unique(s []string) []string {
	result := []string{}
	uniqueVals := make(map[string]bool)
	for _, v := range s {
		if _, ok := uniqueVals[v]; !ok {
			uniqueVals[v] = true
			result = append(result, v)
		}
	}
	return result
}

func ResourcePagerDutyParseColonCompoundID(id string) (string, string, error) {
	parts := strings.Split(id, ":")

	if len(parts) < 2 {
		return "", "", fmt.Errorf("%s: expected colon compound ID to have at least two components", id)
	}

	return parts[0], parts[1], nil
}

func ValidateTZValueDiagFunc(v interface{}, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := v.(string)
	valid := false

	foundAt := sort.SearchStrings(validTZ, value)
	if foundAt < len(validTZ) && validTZ[foundAt] == value {
		valid = true
	}

	if !valid {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("%q is a not valid input. Please refer to the list of allowed Time Zone values at https://developer.pagerduty.com/docs/1afe25e9c94cb-types#time-zone", value),
			AttributePath: p,
		})
	}

	return diags
}

// validTZ at the moment there not an API to fetch this values, so hardcoding
// them here
var validTZ []string = []string{
	"Africa/Algiers",
	"Africa/Cairo",
	"Africa/Casablanca",
	"Africa/Harare",
	"Africa/Johannesburg",
	"Africa/Monrovia",
	"Africa/Nairobi",
	"America/Argentina/Buenos_Aires",
	"America/Bogota",
	"America/Caracas",
	"America/Chicago",
	"America/Chihuahua",
	"America/Denver",
	"America/Godthab",
	"America/Guatemala",
	"America/Guyana",
	"America/Halifax",
	"America/Indiana/Indianapolis",
	"America/Juneau",
	"America/La_Paz",
	"America/Lima",
	"America/Lima",
	"America/Los_Angeles",
	"America/Mazatlan",
	"America/Mexico_City",
	"America/Mexico_City",
	"America/Monterrey",
	"America/Montevideo",
	"America/New_York",
	"America/Phoenix",
	"America/Puerto_Rico",
	"America/Regina",
	"America/Santiago",
	"America/Sao_Paulo",
	"America/St_Johns",
	"America/Tijuana",
	"Asia/Almaty",
	"Asia/Baghdad",
	"Asia/Baku",
	"Asia/Bangkok",
	"Asia/Bangkok",
	"Asia/Chongqing",
	"Asia/Colombo",
	"Asia/Dhaka",
	"Asia/Dhaka",
	"Asia/Hong_Kong",
	"Asia/Irkutsk",
	"Asia/Jakarta",
	"Asia/Jerusalem",
	"Asia/Kabul",
	"Asia/Kamchatka",
	"Asia/Karachi",
	"Asia/Karachi",
	"Asia/Kathmandu",
	"Asia/Kolkata",
	"Asia/Kolkata",
	"Asia/Kolkata",
	"Asia/Kolkata",
	"Asia/Krasnoyarsk",
	"Asia/Kuala_Lumpur",
	"Asia/Kuwait",
	"Asia/Magadan",
	"Asia/Muscat",
	"Asia/Muscat",
	"Asia/Novosibirsk",
	"Asia/Rangoon",
	"Asia/Riyadh",
	"Asia/Seoul",
	"Asia/Shanghai",
	"Asia/Singapore",
	"Asia/Srednekolymsk",
	"Asia/Taipei",
	"Asia/Tashkent",
	"Asia/Tbilisi",
	"Asia/Tehran",
	"Asia/Tokyo",
	"Asia/Tokyo",
	"Asia/Tokyo",
	"Asia/Ulaanbaatar",
	"Asia/Urumqi",
	"Asia/Vladivostok",
	"Asia/Yakutsk",
	"Asia/Yekaterinburg",
	"Asia/Yerevan",
	"Atlantic/Azores",
	"Atlantic/Cape_Verde",
	"Atlantic/South_Georgia",
	"Australia/Adelaide",
	"Australia/Brisbane",
	"Australia/Darwin",
	"Australia/Hobart",
	"Australia/Melbourne",
	"Australia/Melbourne",
	"Australia/Perth",
	"Australia/Sydney",
	"Etc/GMT+12",
	"Etc/UTC",
	"Europe/Amsterdam",
	"Europe/Athens",
	"Europe/Belgrade",
	"Europe/Berlin",
	"Europe/Bratislava",
	"Europe/Brussels",
	"Europe/Bucharest",
	"Europe/Budapest",
	"Europe/Copenhagen",
	"Europe/Dublin",
	"Europe/Helsinki",
	"Europe/Istanbul",
	"Europe/Kaliningrad",
	"Europe/Kiev",
	"Europe/Lisbon",
	"Europe/Ljubljana",
	"Europe/London",
	"Europe/London",
	"Europe/Madrid",
	"Europe/Minsk",
	"Europe/Moscow",
	"Europe/Moscow",
	"Europe/Paris",
	"Europe/Prague",
	"Europe/Riga",
	"Europe/Rome",
	"Europe/Samara",
	"Europe/Sarajevo",
	"Europe/Skopje",
	"Europe/Sofia",
	"Europe/Stockholm",
	"Europe/Tallinn",
	"Europe/Vienna",
	"Europe/Vilnius",
	"Europe/Volgograd",
	"Europe/Warsaw",
	"Europe/Zagreb",
	"Europe/Zurich",
	"Europe/Zurich",
	"Pacific/Apia",
	"Pacific/Auckland",
	"Pacific/Auckland",
	"Pacific/Chatham",
	"Pacific/Fakaofo",
	"Pacific/Fiji",
	"Pacific/Guadalcanal",
	"Pacific/Guam",
	"Pacific/Honolulu",
	"Pacific/Majuro",
	"Pacific/Midway",
	"Pacific/Noumea",
	"Pacific/Pago_Pago",
	"Pacific/Port_Moresby",
	"Pacific/Tongatapu",
}

// CheckJSONEqual returns a function that can be used as in input for
// `resource.TestCheckResourceAttrWith`, it compares two json strings are
// equivalent in data.
func CheckJSONEqual(expected string) resource.CheckResourceAttrWithFunc {
	return resource.CheckResourceAttrWithFunc(func(value string) error {
		var exp interface{}
		if err := json.Unmarshal([]byte(expected), &exp); err != nil {
			return err
		}

		var got interface{}
		if err := json.Unmarshal([]byte(value), &got); err != nil {
			return err
		}

		if !reflect.DeepEqual(exp, got) {
			return fmt.Errorf(`Received value "%v", but expected "%v"`, got, exp)
		}

		return nil
	})
}

// Returns a pair of lists with additions and removals necessary to make set
// `from` turn into set `to`.
func CalculateDiff(from, to []string) (additions, deletions []string) {
	setA := make(map[string]struct{})
	for _, a := range from {
		setA[a] = struct{}{}
	}

	setB := make(map[string]struct{})
	for _, b := range to {
		setB[b] = struct{}{}
	}

	for b := range setB {
		if _, found := setA[b]; !found {
			additions = append(additions, b)
		}
	}

	for a := range setA {
		if _, found := setB[a]; !found {
			deletions = append(deletions, a)
		}
	}

	return
}

var UserAgentAppend string
