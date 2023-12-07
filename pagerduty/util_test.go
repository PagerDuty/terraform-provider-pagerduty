package pagerduty

import (
	"strings"
	"testing"
)

func TestResourcePagerDutyParseColonCompoundID(t *testing.T) {
	resourceIDComponents := []string{"PABC12A", "PABC12B"}
	validCompoundResourceID := strings.Join(resourceIDComponents, ":")

	first, second, err := resourcePagerDutyParseColonCompoundID(validCompoundResourceID)

	if err != nil {
		t.Fatalf("%s: expected success while parsing invalid compound resource id: got %s", validCompoundResourceID, err)
	}

	for i, component := range []string{first, second} {
		expectedResourceIDComponent := resourceIDComponents[i]

		if expectedResourceIDComponent != component {
			t.Errorf(
				"%s: expected component %d of a valid compound resource ID to be %s: got %s",
				validCompoundResourceID,
				i+1,
				expectedResourceIDComponent,
				component,
			)
		}
	}
}

func TestResourcePagerDutyParseColonCompoundIDFailsForInvalidCompoundIDs(t *testing.T) {
	invalidCompoundResourceID := "PABC12APABC12B"

	_, _, err := resourcePagerDutyParseColonCompoundID(invalidCompoundResourceID)

	if err == nil {
		t.Fatalf("%s: expected errors while parsing invalid compound resource id: got success", invalidCompoundResourceID)
	}
}
