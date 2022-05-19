package pagerduty

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

var PagerDutyEventOrchestrationPathParent = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"self": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var PagerDutyEventOrchestrationPathConditions = map[string]*schema.Schema{
	"expression": {
		Type:     schema.TypeString,
		Required: true,
	},
}

func checkExtractions(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
	sn := diff.Get("sets.#").(int)

	for si := 0; si < sn; si++ {
		rn := diff.Get(fmt.Sprintf("sets.%d.rules.#", si)).(int)
		for ri := 0; ri < rn; ri++ {
			res := checkExtractionAttributes(diff, fmt.Sprintf("sets.%d.rules.%d.actions.0.extractions", si, ri))
			if res != nil {
				return res
			}
		}
	}
	return checkExtractionAttributes(diff, "catch_all.0.actions.0.extractions")
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
