package pagerduty

import (
	"reflect"
	"testing"

	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestPagerDutyCustomField_ConvertValueForBuild(t *testing.T) {
	v, _ := convertCustomFieldValueForBuild("5", pagerduty.CustomFieldDataTypeInt, false)
	if v != int64(5) {
		t.Errorf("Unexpected parse int value")
	}

	v, _ = convertCustomFieldValueForBuild("[5, 6]", pagerduty.CustomFieldDataTypeInt, true)
	if !reflect.DeepEqual(v, []interface{}{int64(5), int64(6)}) {
		t.Errorf("Unexpected parse []int value")
	}

	v, _ = convertCustomFieldValueForBuild("5.4", pagerduty.CustomFieldDataTypeFloat, false)
	if v != 5.4 {
		t.Errorf("Unexpected parse float value")
	}

	v, _ = convertCustomFieldValueForBuild("[5.4, 6.7]", pagerduty.CustomFieldDataTypeFloat, true)
	if !reflect.DeepEqual(v, []interface{}{5.4, 6.7}) {
		t.Errorf("Unexpected parse []float value")
	}

	v, _ = convertCustomFieldValueForBuild("false", pagerduty.CustomFieldDataTypeBool, false)
	if v != false {
		t.Errorf("Unexpected parse bool value")
	}

	v, _ = convertCustomFieldValueForBuild(`["foo","bar"]`, pagerduty.CustomFieldDataTypeString, true)
	if !reflect.DeepEqual(v, []interface{}{"foo", "bar"}) {
		t.Errorf("Unexpected parse []string value")
	}
}

func TestPagerDutyCustomField_ConvertDefaultValueForFlatten(t *testing.T) {
	v, _ := convertCustomFieldValueForFlatten(5, false)
	if v != "5" {
		t.Errorf("Unexpected flatten int value")
	}

	v, _ = convertCustomFieldValueForFlatten([]int{5, 6}, true)
	if v != "[5,6]" {
		t.Errorf("Unexpected flatten []int value")
	}

	v, _ = convertCustomFieldValueForFlatten(5.4, false)
	if v != "5.4" {
		t.Errorf("Unexpected flatten float value")
	}

	v, _ = convertCustomFieldValueForFlatten([]float64{5.4, 6.7}, true)
	if v != "[5.4,6.7]" {
		t.Errorf("Unexpected flatten []float value")
	}

	v, _ = convertCustomFieldValueForFlatten(false, false)
	if v != "false" {
		t.Errorf("Unexpected flatten bool value")
	}

	v, _ = convertCustomFieldValueForFlatten([]string{"foo", "bar"}, true)
	if v != `["foo","bar"]` {
		t.Errorf("Unexpected flatten []string value")
	}
}
