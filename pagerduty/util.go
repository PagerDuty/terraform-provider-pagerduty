package pagerduty

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func timeToUTC(v string) (string, error) {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return "", err
	}

	return t.UTC().String(), nil
}

// validateRFC3339 validates that a date string has the correct RFC3339 layout
func validateRFC3339(v interface{}, k string) (we []string, errors []error) {
	value := v.(string)
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		errors = append(errors, fmt.Errorf("%s is not a valid format for argument: %s. Expected format: %s (RFC3339)", value, k, time.RFC3339))
	}

	return
}

// Validate a value against a set of possible values
func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, val := range values {
			if value == val {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%#v is an invalid value for argument %s. Must be one of %#v", value, k, values))
		}
		return
	}
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, string(v.(string)))
	}
	return vs
}
