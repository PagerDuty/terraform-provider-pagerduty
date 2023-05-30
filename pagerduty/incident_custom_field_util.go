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

func validateIncidentCustomFieldDataType() schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		dt := pagerduty.IncidentCustomFieldDataTypeFromString(v.(string))
		if !dt.IsKnown() {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Unknown data_type %v", v),
				AttributePath: p,
			})
		}
		return diags
	}
}

func validateIncidentCustomFieldFieldType() schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		dt := pagerduty.IncidentCustomFieldFieldTypeFromString(v.(string))
		if !dt.IsKnown() {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Unknown field_type %v", v),
				AttributePath: p,
			})
		}
		return diags
	}
}

func validateIncidentCustomFieldValue(value string, datatype pagerduty.IncidentCustomFieldDataType, multiValue bool, generateError func() error) error {
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

	parsedValue, err := convertIncidentCustomFieldValueForBuild(value, datatype, multiValue)
	if err != nil {
		return generateError()
	}

	switch datatype {
	case pagerduty.IncidentCustomFieldDataTypeString:
		return validatedConvertedIncidentCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.IncidentCustomFieldDataTypeInt:
		var i int64
		return validatedConvertedIncidentCustomFieldValue(reflect.TypeOf(i), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.IncidentCustomFieldDataTypeFloat:
		var f float64
		return validatedConvertedIncidentCustomFieldValue(reflect.TypeOf(f), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.IncidentCustomFieldDataTypeBool:
		var b bool
		return validatedConvertedIncidentCustomFieldValue(reflect.TypeOf(b), parsedValue, multiValue, noopValidator, generateError)
	case pagerduty.IncidentCustomFieldDataTypeDateTime:
		return validatedConvertedIncidentCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, datetimeValidator, generateError)
	case pagerduty.IncidentCustomFieldDataTypeUrl:
		return validatedConvertedIncidentCustomFieldValue(reflect.TypeOf(""), parsedValue, multiValue, urlValidator, generateError)
	default:
		return fmt.Errorf("unrecognized datatype: %v", datatype)
	}
}

func validatedConvertedIncidentCustomFieldValue(t reflect.Type, valueToValidate interface{}, expectSlice bool, valueValidator func(interface{}) error, errorGenerator func() error) error {
	if expectSlice {
		arr, ok := valueToValidate.([]interface{})
		if !ok {
			return errorGenerator()
		}
		for _, v := range arr {
			err := validatedConvertedIncidentCustomFieldValue(t, v, false, valueValidator, errorGenerator)
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

func convertIncidentCustomFieldValueForBuild(value string, datatype pagerduty.IncidentCustomFieldDataType, multiValue bool) (interface{}, error) {
	if multiValue {
		var v []interface{}
		err := json.Unmarshal([]byte(value), &v)
		if err != nil {
			return nil, err
		} else {
			if datatype == pagerduty.IncidentCustomFieldDataTypeInt {
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
		case pagerduty.IncidentCustomFieldDataTypeBool:
			return strconv.ParseBool(value)
		case pagerduty.IncidentCustomFieldDataTypeFloat:
			return strconv.ParseFloat(value, 64)
		case pagerduty.IncidentCustomFieldDataTypeInt:
			return strconv.ParseInt(value, 10, 64)
		default:
			return value, nil
		}
	}
}

func convertIncidentCustomFieldValueForFlatten(value interface{}, multiValue bool) (string, error) {
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
