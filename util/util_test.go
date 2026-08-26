package util

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func TestValidateTZValueDiagFunc(t *testing.T) {
	notValidTZ1 := "not a valid TZ"

	cases := []struct {
		given string
		want  diag.Diagnostics
		path  cty.Path
	}{
		{
			given: "America/Montevideo",
			want:  nil,
			path:  cty.Path{},
		},
		{
			given: "America/Indiana/Indianapolis",
			want:  nil,
			path:  cty.Path{},
		},
		{
			given: notValidTZ1,
			want:  diag.Diagnostics{diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("\"%s\" is a not valid input. Please refer to the list of allowed Time Zone values at https://developer.pagerduty.com/docs/1afe25e9c94cb-types#time-zone", notValidTZ1), Detail: "", AttributePath: cty.Path{cty.GetAttrStep{Name: "time_zone"}}}},
			path:  cty.Path{cty.GetAttrStep{Name: "time_zone"}},
		},
	}

	for _, c := range cases {
		got := ValidateTZValueDiagFunc(c.given, c.path)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("want %v; got %v", c.want, got)
		}
	}
}

func TestValidateTimeOfDay(t *testing.T) {
	cases := []struct {
		given   string
		wantErr bool
	}{
		// Valid HH:MM:SS values.
		{given: "08:00:00", wantErr: false},
		{given: "00:00:00", wantErr: false},
		{given: "23:59:59", wantErr: false},
		// Out-of-range components.
		{given: "24:00:00", wantErr: true},
		{given: "12:60:00", wantErr: true},
		// Junk-wrapped values that must be rejected once the regex is anchored.
		{given: "12:30:00xyz", wantErr: true},
		{given: "xyz08:00:00", wantErr: true},
		{given: "99:08:00:00", wantErr: true},
		{given: "", wantErr: true},
	}

	for _, c := range cases {
		_, errs := ValidateTimeOfDay(c.given, "start_time_of_day")
		gotErr := len(errs) > 0
		if gotErr != c.wantErr {
			t.Errorf("ValidateTimeOfDay(%q): wantErr=%v, got errs=%v", c.given, c.wantErr, errs)
		}
	}
}
