package pagerduty

import (
	"reflect"
	"testing"
)

func TestJiraCloudCustomFieldValue(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want interface{}
	}{
		{"numeric id stays a string", "42", "42"},
		{"boolean text stays a string", "true", "true"},
		{"null text stays a string", "null", "null"},
		{"plain text", "pagerduty, incident", "pagerduty, incident"},
		{"quoted text is not unwrapped", `"42"`, `"42"`},
		{"empty", "", ""},
		{"object is decoded", `{"id":"10000","displayName":"Sec Level 1"}`, map[string]interface{}{"id": "10000", "displayName": "Sec Level 1"}},
		{"array is decoded", `[{"id":"1"},{"id":"2"}]`, []interface{}{map[string]interface{}{"id": "1"}, map[string]interface{}{"id": "2"}}},
		{"malformed object text stays a string", `{"id":`, `{"id":`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := jiraCloudCustomFieldValue(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("jiraCloudCustomFieldValue(%q) = %#v (%T), want %#v (%T)", c.in, got, got, c.want, c.want)
			}
		})
	}
}
