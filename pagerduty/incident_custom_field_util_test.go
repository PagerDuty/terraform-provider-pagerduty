package pagerduty

import (
	"reflect"
	"testing"

	"github.com/heimweh/go-pagerduty/pagerduty"
)

func TestPagerDutyIncidentCustomField_ConvertValueForBuild(t *testing.T) {
	v, _ := convertIncidentCustomFieldValueForBuild("5", pagerduty.IncidentCustomFieldDataTypeInt, false)
	if v != int64(5) {
		t.Errorf("Unexpected parse int value")
	}

	v, _ = convertIncidentCustomFieldValueForBuild("[5, 6]", pagerduty.IncidentCustomFieldDataTypeInt, true)
	if !reflect.DeepEqual(v, []interface{}{int64(5), int64(6)}) {
		t.Errorf("Unexpected parse []int value")
	}

	v, _ = convertIncidentCustomFieldValueForBuild("5.4", pagerduty.IncidentCustomFieldDataTypeFloat, false)
	if v != 5.4 {
		t.Errorf("Unexpected parse float value")
	}

	v, _ = convertIncidentCustomFieldValueForBuild("[5.4, 6.7]", pagerduty.IncidentCustomFieldDataTypeFloat, true)
	if !reflect.DeepEqual(v, []interface{}{5.4, 6.7}) {
		t.Errorf("Unexpected parse []float value")
	}

	v, _ = convertIncidentCustomFieldValueForBuild("false", pagerduty.IncidentCustomFieldDataTypeBool, false)
	if v != false {
		t.Errorf("Unexpected parse bool value")
	}

	v, _ = convertIncidentCustomFieldValueForBuild(`["foo","bar"]`, pagerduty.IncidentCustomFieldDataTypeString, true)
	if !reflect.DeepEqual(v, []interface{}{"foo", "bar"}) {
		t.Errorf("Unexpected parse []string value")
	}
}

func TestPagerDutyIncidentCustomField_ConvertDefaultValueForFlatten(t *testing.T) {
	v, _ := convertIncidentCustomFieldValueForFlatten(5, false)
	if v != "5" {
		t.Errorf("Unexpected flatten int value")
	}

	v, _ = convertIncidentCustomFieldValueForFlatten([]int{5, 6}, true)
	if v != "[5,6]" {
		t.Errorf("Unexpected flatten []int value")
	}

	v, _ = convertIncidentCustomFieldValueForFlatten(5.4, false)
	if v != "5.4" {
		t.Errorf("Unexpected flatten float value")
	}

	v, _ = convertIncidentCustomFieldValueForFlatten([]float64{5.4, 6.7}, true)
	if v != "[5.4,6.7]" {
		t.Errorf("Unexpected flatten []float value")
	}

	v, _ = convertIncidentCustomFieldValueForFlatten(false, false)
	if v != "false" {
		t.Errorf("Unexpected flatten bool value")
	}

	v, _ = convertIncidentCustomFieldValueForFlatten([]string{"foo", "bar"}, true)
	if v != `["foo","bar"]` {
		t.Errorf("Unexpected flatten []string value")
	}
}
