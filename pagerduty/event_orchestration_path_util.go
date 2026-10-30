package pagerduty

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

var eventOrchestrationPathConditionsSchema = map[string]*schema.Schema{
	"expression": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var eventOrchestrationPathVariablesSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"path": {
		Type:     schema.TypeString,
		Required: true,
	},
	"type": {
		Type:     schema.TypeString,
		Required: true,
	},
	"value": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var eventOrchestrationPathExtractionsSchema = map[string]*schema.Schema{
	"regex": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"source": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"target": {
		Type:     schema.TypeString,
		Required: true,
	},
	"template": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

var eventOrchestrationAutomationActionObjectSchema = map[string]*schema.Schema{
	"key": {
		Type:     schema.TypeString,
		Required: true,
	},
	"value": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var eventOrchestrationIncidentCustomFieldsObjectSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"value": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var eventOrchestrationIncidentTypeObjectSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var eventOrchestrationAutomationActionSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"url": {
		Type:     schema.TypeString,
		Required: true,
	},
	"auto_send": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"header": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: eventOrchestrationAutomationActionObjectSchema,
		},
	},
	"parameter": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: eventOrchestrationAutomationActionObjectSchema,
		},
	},
	"trigger_types": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
}

func invalidExtractionRegexTemplateNilConfig() string {
	return `
		extraction {
			target = "event.summary"
		}`
}

func invalidExtractionRegexTemplateValConfig() string {
	return `
		extraction {
			regex = ".*"
			template = "hi"
			target = "event.summary"
		}`
}

func invalidExtractionRegexNilSourceConfig() string {
	return `
		extraction {
			regex = ".*"
			target = "event.summary"
		}`
}

func validateEventOrchestrationPathSeverity() schema.SchemaValidateDiagFunc {
	return validateValueDiagFunc([]string{
		"info",
		"error",
		"warning",
		"critical",
	})
}

func validateEventOrchestrationPathEventAction() schema.SchemaValidateDiagFunc {
	return validateValueDiagFunc([]string{
		"trigger",
		"resolve",
	})
}

func checkExtractions(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
	sn := diff.Get("set.#").(int)

	for si := 0; si < sn; si++ {
		rn := diff.Get(fmt.Sprintf("set.%d.rule.#", si)).(int)
		for ri := 0; ri < rn; ri++ {
			res := checkExtractionAttributes(diff, fmt.Sprintf("set.%d.rule.%d.actions.0.extraction", si, ri))
			if res != nil {
				return res
			}
		}
	}
	return checkExtractionAttributes(diff, "catch_all.0.actions.0.extraction")
}

func checkExtractionAttributes(diff *schema.ResourceDiff, loc string) error {
	num := diff.Get(fmt.Sprintf("%s.#", loc)).(int)
	for i := 0; i < num; i++ {
		prefix := fmt.Sprintf("%s.%d", loc, i)
		r := diff.Get(fmt.Sprintf("%s.regex", prefix)).(string)
		t := diff.Get(fmt.Sprintf("%s.template", prefix)).(string)

		if r == "" && t == "" {
			return fmt.Errorf("Invalid configuration in %s: regex and template cannot both be null", prefix)
		}
		if r != "" && t != "" {
			return fmt.Errorf("Invalid configuration in %s: regex and template cannot both have values", prefix)
		}

		s := diff.Get(fmt.Sprintf("%s.source", prefix)).(string)
		if r != "" && s == "" {
			return fmt.Errorf("Invalid configuration in %s: source can't be blank", prefix)
		}
	}
	return nil
}

func expandEventOrchestrationPathConditions(v interface{}) []*pagerduty.EventOrchestrationPathRuleCondition {
	conditions := []*pagerduty.EventOrchestrationPathRuleCondition{}

	for _, cond := range v.([]interface{}) {
		if cond == nil {
			continue
		}
		c, ok := cond.(map[string]interface{})
		if !ok || c == nil {
			continue
		}

		cx := &pagerduty.EventOrchestrationPathRuleCondition{
			Expression: c["expression"].(string),
		}

		conditions = append(conditions, cx)
	}

	return conditions
}

func flattenEventOrchestrationPathConditions(conditions []*pagerduty.EventOrchestrationPathRuleCondition) []interface{} {
	var flattendConditions []interface{}

	for _, condition := range conditions {
		flattendCondition := map[string]interface{}{
			"expression": condition.Expression,
		}
		flattendConditions = append(flattendConditions, flattendCondition)
	}

	return flattendConditions
}

func expandEventOrchestrationPathVariables(v interface{}) []*pagerduty.EventOrchestrationPathActionVariables {
	res := []*pagerduty.EventOrchestrationPathActionVariables{}

	for _, er := range v.([]interface{}) {
		rer, ok := er.(map[string]interface{})
		if !ok || rer == nil {
			continue
		}

		av := &pagerduty.EventOrchestrationPathActionVariables{
			Name:  rer["name"].(string),
			Path:  rer["path"].(string),
			Type:  rer["type"].(string),
			Value: rer["value"].(string),
		}

		res = append(res, av)
	}

	return res
}

func flattenEventOrchestrationPathVariables(v []*pagerduty.EventOrchestrationPathActionVariables) []interface{} {
	var res []interface{}

	for _, s := range v {
		fv := map[string]interface{}{
			"name":  s.Name,
			"path":  s.Path,
			"type":  s.Type,
			"value": s.Value,
		}
		res = append(res, fv)
	}
	return res
}

func expandEventOrchestrationPathIncidentCustomFields(v interface{}) []*pagerduty.EventOrchestrationPathIncidentCustomFieldUpdate {
	res := []*pagerduty.EventOrchestrationPathIncidentCustomFieldUpdate{}

	for _, eai := range v.([]interface{}) {
		ea, ok := eai.(map[string]interface{})
		if !ok || ea == nil {
			continue
		}
		ext := &pagerduty.EventOrchestrationPathIncidentCustomFieldUpdate{
			ID:    ea["id"].(string),
			Value: ea["value"].(string),
		}
		res = append(res, ext)
	}
	return res
}

func expandEventOrchestrationPathIncidentType(v interface{}) *pagerduty.EventOrchestrationPathIncidentType {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	it := l[0].(map[string]interface{})
	return &pagerduty.EventOrchestrationPathIncidentType{
		ID:   it["id"].(string),
		Name: it["name"].(string),
	}
}

func expandEventOrchestrationPathExtractions(v interface{}) []*pagerduty.EventOrchestrationPathActionExtractions {
	res := []*pagerduty.EventOrchestrationPathActionExtractions{}

	for _, eai := range v.([]interface{}) {
		ea, ok := eai.(map[string]interface{})
		if !ok || ea == nil {
			continue
		}
		ext := &pagerduty.EventOrchestrationPathActionExtractions{
			Target:   ea["target"].(string),
			Regex:    ea["regex"].(string),
			Template: ea["template"].(string),
			Source:   ea["source"].(string),
		}
		res = append(res, ext)
	}
	return res
}

func flattenEventOrchestrationPathExtractions(e []*pagerduty.EventOrchestrationPathActionExtractions) []interface{} {
	var res []interface{}

	for _, s := range e {
		e := map[string]interface{}{
			"target":   s.Target,
			"regex":    s.Regex,
			"template": s.Template,
			"source":   s.Source,
		}
		res = append(res, e)
	}
	return res
}

func expandEventOrchestrationPathAutomationActions(v interface{}) []*pagerduty.EventOrchestrationPathAutomationAction {
	result := []*pagerduty.EventOrchestrationPathAutomationAction{}

	for _, i := range v.([]interface{}) {
		a, ok := i.(map[string]interface{})
		if !ok || a == nil {
			continue
		}
		aa := &pagerduty.EventOrchestrationPathAutomationAction{
			Name:         a["name"].(string),
			Url:          a["url"].(string),
			AutoSend:     a["auto_send"].(bool),
			Headers:      expandEventOrchestrationAutomationActionObjects(a["header"]),
			Parameters:   expandEventOrchestrationAutomationActionObjects(a["parameter"]),
			TriggerTypes: expandEventOrchestrationAutomationTriggerTypes(a["trigger_types"]),
		}

		result = append(result, aa)
	}

	return result
}

func expandEventOrchestrationAutomationActionObjects(v interface{}) []*pagerduty.EventOrchestrationPathAutomationActionObject {
	result := []*pagerduty.EventOrchestrationPathAutomationActionObject{}

	for _, i := range v.([]interface{}) {
		o, ok := i.(map[string]interface{})
		if !ok || o == nil {
			continue
		}
		obj := &pagerduty.EventOrchestrationPathAutomationActionObject{
			Key:   o["key"].(string),
			Value: o["value"].(string),
		}

		result = append(result, obj)
	}

	return result
}

func expandEventOrchestrationAutomationTriggerTypes(v interface{}) []string {
	if v == nil {
		return nil
	}

	var result []string

	for _, i := range v.([]interface{}) {
		result = append(result, i.(string))
	}

	return result
}

func flattenEventOrchestrationIncidentCustomFieldUpdates(v []*pagerduty.EventOrchestrationPathIncidentCustomFieldUpdate) []interface{} {
	var result []interface{}

	for _, i := range v {
		custom_field := map[string]string{
			"id":    i.ID,
			"value": i.Value,
		}

		result = append(result, custom_field)
	}

	return result
}

func flattenEventOrchestrationPathIncidentType(v *pagerduty.EventOrchestrationPathIncidentType) []interface{} {
	if v == nil {
		return nil
	}
	return []interface{}{
		map[string]string{
			"id":   v.ID,
			"name": v.Name,
		},
	}
}

func flattenEventOrchestrationAutomationActions(v []*pagerduty.EventOrchestrationPathAutomationAction) []interface{} {
	var result []interface{}

	for _, i := range v {
		pdaa := map[string]interface{}{
			"name":          i.Name,
			"url":           i.Url,
			"auto_send":     i.AutoSend,
			"header":        flattenEventOrchestrationAutomationActionObjects(i.Headers),
			"parameter":     flattenEventOrchestrationAutomationActionObjects(i.Parameters),
			"trigger_types": i.TriggerTypes,
		}

		result = append(result, pdaa)
	}

	return result
}

func flattenEventOrchestrationAutomationActionObjects(v []*pagerduty.EventOrchestrationPathAutomationActionObject) []interface{} {
	var result []interface{}

	for _, i := range v {
		pdaa := map[string]interface{}{
			"key":   i.Key,
			"value": i.Value,
		}

		result = append(result, pdaa)
	}

	return result
}

func convertEventOrchestrationPathWarningsToDiagnostics(warnings []*pagerduty.EventOrchestrationPathWarning, diags diag.Diagnostics) diag.Diagnostics {
	if warnings == nil {
		return diags
	}

	for _, warning := range warnings {
		diag := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  warning.Message,
			Detail:   fmt.Sprintf("Feature: %s\nFeature Type: %s\nRule ID: %s\nWarning Type: %s", warning.Feature, warning.FeatureType, warning.RuleId, warning.WarningType),
		}
		diags = append(diags, diag)
	}

	return diags
}

// checkExistingOrchestrationPathConfig fetches the current path config for
// pathType ("service", "global", "router", "unrouted") and returns an error if
// non-trivial configuration already exists that would be silently overwritten.
// resourceType is the Terraform resource name used in the import hint.
func checkExistingOrchestrationPathConfig(ctx context.Context, client *pagerduty.Client, orchID, pathType, resourceType string) error {
	var existingPath *pagerduty.EventOrchestrationPath

	retryErr := retry.RetryContext(ctx, 5*time.Second, func() *retry.RetryError {
		path, _, err := client.EventOrchestrationPaths.GetContext(ctx, orchID, pathType)
		if err != nil {
			if isErrCode(err, http.StatusForbidden) {
				// For service paths the API returns 403 (not 404) when the parent
				// service has been deleted; treat it as "no existing config".
				// For other path types a 403 is a real permission error.
				if pathType == "service" {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		existingPath = path
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	if existingPath == nil {
		return nil
	}

	hasNonTrivialConfig := false

	if pathType == "router" {
		// Router is trivial when catch_all routes to "unrouted" and there are no rules.
		if existingPath.CatchAll != nil && existingPath.CatchAll.Actions != nil {
			rt := existingPath.CatchAll.Actions.RouteTo
			if rt != "" && rt != "unrouted" {
				hasNonTrivialConfig = true
			}
		}
		if !hasNonTrivialConfig && len(existingPath.Sets) > 0 && len(existingPath.Sets[0].Rules) > 0 {
			hasNonTrivialConfig = true
		}
	} else {
		if existingPath.CatchAll != nil && existingPath.CatchAll.Actions != nil {
			a := existingPath.CatchAll.Actions
			// Mirror every field of EventOrchestrationPathRuleActions; update here when
			// new action fields are added to that struct in go-pagerduty.
			// DynamicRouteTo is router-only and intentionally omitted (never set on
			// non-router catch_all by the API).
			// The unrouted path always has suppress=true set by the API; exclude it
			// from the trivial check so a fresh orchestration is never falsely blocked.
			suppressNonTrivial := a.Suppress && pathType != "unrouted"
			if suppressNonTrivial || a.DropEvent || a.Priority != "" || a.Severity != "" ||
				a.EventAction != "" || a.Annotate != "" || a.RouteTo != "" ||
				a.Suspend != nil || a.EscalationPolicy != nil ||
				len(a.Variables) > 0 || len(a.Extractions) > 0 ||
				len(a.AutomationActions) > 0 || len(a.PagerdutyAutomationActions) > 0 ||
				len(a.IncidentCustomFieldUpdates) > 0 {
				hasNonTrivialConfig = true
			}
		}
		if !hasNonTrivialConfig {
			if len(existingPath.Sets) >= 2 || (len(existingPath.Sets) > 0 && len(existingPath.Sets[0].Rules) > 0) {
				hasNonTrivialConfig = true
			}
		}
	}

	if hasNonTrivialConfig {
		return fmt.Errorf(
			"the %s orchestration (ID: %s) has existing configuration that might be overwritten; "+
				"please import this resource before creating it using: terraform import %s.<resource_name> %s",
			pathType, orchID, resourceType, orchID,
		)
	}

	return nil
}

func emptyOrchestrationPathStructBuilder(pathType string) *pagerduty.EventOrchestrationPath {
	commonEmptyOrchestrationPath := func() *pagerduty.EventOrchestrationPath {
		return &pagerduty.EventOrchestrationPath{
			CatchAll: &pagerduty.EventOrchestrationPathCatchAll{
				Actions: nil,
			},
			Sets: []*pagerduty.EventOrchestrationPathSet{
				{
					ID:    "start",
					Rules: []*pagerduty.EventOrchestrationPathRule{},
				},
			},
		}
	}
	routerEmptyOrchestrationPath := func() *pagerduty.EventOrchestrationPath {
		return &pagerduty.EventOrchestrationPath{
			CatchAll: &pagerduty.EventOrchestrationPathCatchAll{
				Actions: &pagerduty.EventOrchestrationPathRuleActions{
					RouteTo: "unrouted",
				},
			},
			Sets: []*pagerduty.EventOrchestrationPathSet{
				{
					ID:    "start",
					Rules: []*pagerduty.EventOrchestrationPathRule{},
				},
			},
		}
	}

	if pathType == "router" {
		return routerEmptyOrchestrationPath()
	}

	return commonEmptyOrchestrationPath()
}

func isNonEmptyList(arg interface{}) bool {
	return !isNilFunc(arg) && len(arg.([]interface{})) > 0
}
