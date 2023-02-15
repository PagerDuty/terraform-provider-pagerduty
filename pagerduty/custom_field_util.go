package pagerduty

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func validateCustomFieldDataType() schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		dt := pagerduty.CustomFieldDataTypeFromString(v.(string))
		if !dt.IsKnown() {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Unknown datatype %v", v),
				AttributePath: p,
			})
		}
		if !dt.IsAllowedOnField() {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Datatype %v is not allowed on fields", v),
				AttributePath: p,
			})
		}
		return diags
	}
}

func validateCustomFieldValue(value string, datatype pagerduty.CustomFieldDataType, multiValue bool, generateError func() error) error {
	noopValidator := func(v interface{}) error {
		return nil
	}
	datetimeValidator := func(v interface{}) error {
		_, err := time.Parse(time.RFC3339, v.(string))
		return err
	}
	urlValidator := func(v interface{}) error {
		u, err := url.Parse(v.(string))
		if err != nil {
			return err
		}
		if !u.IsAbs() {
			return fmt.Errorf(`parsed url default value "%v" is not an absolute url`, v)
		}
		return nil
	}

	parsedValue, err := convertCustomFieldValueForBuild(value, datatype, multiValue)
	if err != nil {
		return generateError()
	}

	switch datatype {
	case pagerduty.CustomFieldDataTypeString:
		return validatedConvertedCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.CustomFieldDataTypeFieldOption:
		return validatedConvertedCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.CustomFieldDataTypeInt:
		var i int64
		return validatedConvertedCustomFieldValue(reflect.TypeOf(i), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.CustomFieldDataTypeFloat:
		var f float64
		return validatedConvertedCustomFieldValue(reflect.TypeOf(f), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.CustomFieldDataTypeBool:
		var b bool
		return validatedConvertedCustomFieldValue(reflect.TypeOf(b), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.CustomFieldDataTypeDateTime:
		return validatedConvertedCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, datetimeValidator, generateError)
	case pagerduty.CustomFieldDataTypeUrl:
		return validatedConvertedCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, urlValidator, generateError)
	default:
		return fmt.Errorf("unrecognized datatype: %v", datatype)
	}
}

func validatedConvertedCustomFieldValue(t reflect.Type, valueToValidate interface{}, expectSlice bool, valueValidator func(interface{}) error, errorGenerator func() error) error {
	if expectSlice {
		arr, ok := valueToValidate.([]interface{})
		if !ok {
			return errorGenerator()
		}
		for _, v := range arr {
			err := validatedConvertedCustomFieldValue(t, v, false, valueValidator, errorGenerator)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		if reflect.TypeOf(valueToValidate) != t {
			return errorGenerator()
		}
		return valueValidator(valueToValidate)
	}
}

func convertCustomFieldValueForBuild(value string, datatype pagerduty.CustomFieldDataType, multiValue bool) (interface{}, error) {
	if multiValue {
		var v []interface{}
		err := json.Unmarshal([]byte(value), &v)
		if err != nil {
			return nil, err
		} else {
			if datatype == pagerduty.CustomFieldDataTypeInt {
				var iv []interface{}
				for _, ev := range v {
					fv, ok := ev.(float64)
					if !ok {
						return nil, fmt.Errorf("value %v not parseable as a number", ev)
					}
					iv = append(iv, int64(math.Round(fv)))
				}
				v = iv
			}
			return v, nil
		}
	} else {
		switch datatype {
		case pagerduty.CustomFieldDataTypeBool:
			return strconv.ParseBool(value)
		case pagerduty.CustomFieldDataTypeFloat:
			return strconv.ParseFloat(value, 64)
		case pagerduty.CustomFieldDataTypeInt:
			return strconv.ParseInt(value, 10, 64)
		default:
			return value, nil
		}
	}
}

func convertCustomFieldValueForFlatten(value interface{}, multiValue bool) (string, error) {
	if multiValue {
		b, err := json.Marshal(value)
		if err != nil {
			return "", err
		} else {
			return string(b), nil
		}
	} else {
		return fmt.Sprintf("%v", value), nil
	}
}
