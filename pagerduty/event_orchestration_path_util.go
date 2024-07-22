package pagerduty

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		c := cond.(map[string]interface{})

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
		rer := er.(map[string]interface{})

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
		ea := eai.(map[string]interface{})
		ext := &pagerduty.EventOrchestrationPathIncidentCustomFieldUpdate{
			ID:    ea["id"].(string),
			Value: ea["value"].(string),
		}
		res = append(res, ext)
	}
	return res
}

func expandEventOrchestrationPathExtractions(v interface{}) []*pagerduty.EventOrchestrationPathActionExtractions {
	res := []*pagerduty.EventOrchestrationPathActionExtractions{}

	for _, eai := range v.([]interface{}) {
		ea := eai.(map[string]interface{})
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
		a := i.(map[string]interface{})
		aa := &pagerduty.EventOrchestrationPathAutomationAction{
			Name:       a["name"].(string),
			Url:        a["url"].(string),
			AutoSend:   a["auto_send"].(bool),
			Headers:    expandEventOrchestrationAutomationActionObjects(a["header"]),
			Parameters: expandEventOrchestrationAutomationActionObjects(a["parameter"]),
		}

		result = append(result, aa)
	}

	return result
}

func expandEventOrchestrationAutomationActionObjects(v interface{}) []*pagerduty.EventOrchestrationPathAutomationActionObject {
	result := []*pagerduty.EventOrchestrationPathAutomationActionObject{}

	for _, i := range v.([]interface{}) {
		o := i.(map[string]interface{})
		obj := &pagerduty.EventOrchestrationPathAutomationActionObject{
			Key:   o["key"].(string),
			Value: o["value"].(string),
		}

		result = append(result, obj)
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

func flattenEventOrchestrationAutomationActions(v []*pagerduty.EventOrchestrationPathAutomationAction) []interface{} {
	var result []interface{}

	for _, i := range v {
		pdaa := map[string]interface{}{
			"name":      i.Name,
			"url":       i.Url,
			"auto_send": i.AutoSend,
			"header":    flattenEventOrchestrationAutomationActionObjects(i.Headers),
			"parameter": flattenEventOrchestrationAutomationActionObjects(i.Parameters),
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
