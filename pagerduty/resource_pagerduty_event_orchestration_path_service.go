package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

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

var eventOrchestrationPathServiceCatchAllActionsSchema = map[string]*schema.Schema{
	"suppress": {
		Type:     schema.TypeBool,
		Optional: true,
	},
	"suspend": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"priority": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"annotate": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"pagerduty_automation_actions": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"action_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
	"automation_actions": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
				"headers": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: eventOrchestrationAutomationActionObjectSchema,
					},
				},
				"parameters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: eventOrchestrationAutomationActionObjectSchema,
					},
				},
			},
		},
	},
	"severity": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validateEventOrchestrationPathSeverity(),
	},
	"event_action": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validateEventOrchestrationPathEventAction(),
	},
	"variables": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: eventOrchestrationPathVariablesSchema,
		},
	},
	"extractions": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: eventOrchestrationPathExtractionsSchema,
		},
	},
}

var eventOrchestrationPathServiceRuleActionsSchema = buildEventOrchestrationPathServiceRuleActionsSchema()

func buildEventOrchestrationPathServiceRuleActionsSchema() map[string]*schema.Schema {
	a := eventOrchestrationPathServiceCatchAllActionsSchema
	a["route_to"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return a
}

func resourcePagerDutyEventOrchestrationPathService() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePagerDutyEventOrchestrationPathServiceRead,
		Create: resourcePagerDutyEventOrchestrationPathServiceCreate,
		Update: resourcePagerDutyEventOrchestrationPathServiceUpdate,
		Delete: resourcePagerDutyEventOrchestrationPathServiceDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyEventOrchestrationPathServiceImport,
		},
		CustomizeDiff: checkExtractions,
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"label": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: PagerDutyEventOrchestrationPathConditions,
										},
									},
									"actions": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: eventOrchestrationPathServiceRuleActionsSchema,
										},
									},
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"catch_all": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: eventOrchestrationPathServiceRuleActionsSchema,
							},
						},
					},
				},
			},
		},
	}
}

func resourcePagerDutyEventOrchestrationPathServiceRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		id := d.Id()
		t := "service"
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type %s for orchestration: %s", t, id)

		if path, _, err := client.EventOrchestrationPaths.Get(d.Id(), t); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if path != nil {
			setEventOrchestrationPathServiceProps(d, path)
		}
		return nil
	})

}

func resourcePagerDutyEventOrchestrationPathServiceCreate(d *schema.ResourceData, meta interface{}) error {
	return resourcePagerDutyEventOrchestrationPathServiceUpdate(d, meta)
}

func resourcePagerDutyEventOrchestrationPathServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	payload := buildServicePathStruct(d)
	var servicePath *pagerduty.EventOrchestrationPath

	log.Printf("[INFO] Creating PagerDuty Event Orchestration Service Path: %s", payload.Parent.ID)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if path, _, err := client.EventOrchestrationPaths.Update(payload.Parent.ID, "service", payload); err != nil {
			return resource.RetryableError(err)
		} else if path != nil {
			d.SetId(path.Parent.ID)
			servicePath = path
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	setEventOrchestrationPathServiceProps(d, servicePath)

	return nil
}

func resourcePagerDutyEventOrchestrationPathServiceDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourcePagerDutyEventOrchestrationPathServiceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	id := d.Id()

	_, _, pErr := client.EventOrchestrationPaths.Get(id, "service")
	if pErr != nil {
		return []*schema.ResourceData{}, pErr
	}

	d.SetId(id)
	d.Set("service", id)

	return []*schema.ResourceData{d}, nil
}

func buildServicePathStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	return &pagerduty.EventOrchestrationPath{
		Parent: &pagerduty.EventOrchestrationPathReference{
			ID: d.Get("service").(string),
		},
		Sets:     expandServicePathSets(d.Get("sets")),
		CatchAll: expandServicePathCatchAll(d.Get("catch_all")),
	}
}

// TODO: see if we can reuse expand functions for all orch path sets.
// Maybe pass in the rule actions and catch-all rule actions expanding function?
func expandServicePathSets(v interface{}) []*pagerduty.EventOrchestrationPathSet {
	var sets []*pagerduty.EventOrchestrationPathSet

	for _, set := range v.([]interface{}) {
		s := set.(map[string]interface{})

		orchPathSet := &pagerduty.EventOrchestrationPathSet{
			ID:    s["id"].(string),
			Rules: expandServicePathRules(s["rules"].(interface{})),
		}

		sets = append(sets, orchPathSet)
	}

	return sets
}

func expandServicePathRules(v interface{}) []*pagerduty.EventOrchestrationPathRule {
	items := v.([]interface{})
	rules := []*pagerduty.EventOrchestrationPathRule{}

	for _, rule := range items {
		r := rule.(map[string]interface{})

		ruleInSet := &pagerduty.EventOrchestrationPathRule{
			ID:       r["id"].(string),
			Label:    r["label"].(string),
			Disabled: r["disabled"].(bool),
			// TODO: move conditions logic to util
			Conditions: expandRouterConditions(r["conditions"].(interface{})),
			Actions:    expandServicePathActions(r["actions"].([]interface{})),
		}

		rules = append(rules, ruleInSet)
	}
	return rules
}

func expandServicePathCatchAll(v interface{}) *pagerduty.EventOrchestrationPathCatchAll {
	var catchAll = new(pagerduty.EventOrchestrationPathCatchAll)

	for _, ca := range v.([]interface{}) {
		if ca != nil {
			am := ca.(map[string]interface{})
			catchAll.Actions = expandServicePathActions(am["actions"])
		}
	}

	return catchAll
}

func expandServicePathActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = &pagerduty.EventOrchestrationPathRuleActions{
		AutomationActions:          []*pagerduty.EventOrchestrationPathAutomationAction{},
		PagerdutyAutomationActions: []*pagerduty.EventOrchestrationPathPagerdutyAutomationAction{},
		Variables:                  []*pagerduty.EventOrchestrationPathActionVariables{},
		Extractions:                []*pagerduty.EventOrchestrationPathActionExtractions{},
	}

	for _, i := range v.([]interface{}) {
		if i == nil {
			continue
		}
		a := i.(map[string]interface{})

		actions.RouteTo = a["route_to"].(string)
		actions.Suppress = a["suppress"].(bool)
		actions.Suspend = intTypeToIntPtr(a["suspend"].(int))
		actions.Priority = a["priority"].(string)
		actions.Annotate = a["annotate"].(string)
		actions.Severity = a["severity"].(string)
		actions.EventAction = a["event_action"].(string)
		actions.PagerdutyAutomationActions = expandServicePathPagerDutyAutomationActions(a["pagerduty_automation_actions"])
		actions.AutomationActions = expandServicePathAutomationActions(a["automation_actions"])
		actions.Variables = expandEventOrchestrationPathVariables(a["variables"])
		actions.Extractions = expandEventOrchestrationPathExtractions(a["extractions"])
	}

	return actions
}

func expandServicePathPagerDutyAutomationActions(v interface{}) []*pagerduty.EventOrchestrationPathPagerdutyAutomationAction {
	result := []*pagerduty.EventOrchestrationPathPagerdutyAutomationAction{}

	for _, i := range v.([]interface{}) {
		a := i.(map[string]interface{})
		pdaa := &pagerduty.EventOrchestrationPathPagerdutyAutomationAction{
			ActionId: a["action_id"].(string),
		}

		result = append(result, pdaa)
	}

	return result
}

func expandServicePathAutomationActions(v interface{}) []*pagerduty.EventOrchestrationPathAutomationAction {
	result := []*pagerduty.EventOrchestrationPathAutomationAction{}

	for _, i := range v.([]interface{}) {
		a := i.(map[string]interface{})
		aa := &pagerduty.EventOrchestrationPathAutomationAction{
			Name:       a["name"].(string),
			Url:        a["url"].(string),
			AutoSend:   a["auto_send"].(bool),
			Headers:    expandEventOrchestrationAutomationActionObjects(a["headers"]),
			Parameters: expandEventOrchestrationAutomationActionObjects(a["parameters"]),
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

func setEventOrchestrationPathServiceProps(d *schema.ResourceData, p *pagerduty.EventOrchestrationPath) error {
	d.SetId(p.Parent.ID)
	d.Set("service", p.Parent.ID)
	// TODO: see if we can reuse expand functions for all orch path sets.
	// Maybe pass in the rule actions and catch-all rule actions expanding function?
	d.Set("sets", flattenServicePathSets(p.Sets))
	d.Set("catch_all", flattenServicePathCatchAll(p.CatchAll))
	return nil
}

func flattenServicePathSets(orchPathSets []*pagerduty.EventOrchestrationPathSet) []interface{} {
	var flattenedSets []interface{}

	for _, set := range orchPathSets {
		flattenedSet := map[string]interface{}{
			"id":    set.ID,
			"rules": flattenServicePathRules(set.Rules),
		}
		flattenedSets = append(flattenedSets, flattenedSet)
	}
	return flattenedSets
}

func flattenServicePathCatchAll(catchAll *pagerduty.EventOrchestrationPathCatchAll) []map[string]interface{} {
	var caMap []map[string]interface{}

	c := make(map[string]interface{})

	c["actions"] = flattenServicePathActions(catchAll.Actions)
	caMap = append(caMap, c)

	return caMap
}

func flattenServicePathRules(rules []*pagerduty.EventOrchestrationPathRule) []interface{} {
	var flattenedRules []interface{}

	for _, rule := range rules {
		flattenedRule := map[string]interface{}{
			"id":       rule.ID,
			"label":    rule.Label,
			"disabled": rule.Disabled,
			// TODO: move conditions logic to util
			"conditions": flattenRouterConditions(rule.Conditions),
			"actions":    flattenServicePathActions(rule.Actions),
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	return flattenedRules
}

func flattenServicePathActions(actions *pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	flattenedAction := map[string]interface{}{
		"route_to":     actions.RouteTo,
		"severity":     actions.Severity,
		"event_action": actions.EventAction,
		"suppress":     actions.Suppress,
		"suspend":      actions.Suspend,
		"priority":     actions.Priority,
		"annotate":     actions.Annotate,
	}

	if actions.Variables != nil {
		flattenedAction["variables"] = flattenEventOrchestrationPathVariables(actions.Variables)
	}
	if actions.Extractions != nil {
		flattenedAction["extractions"] = flattenEventOrchestrationPathExtractions(actions.Extractions)
	}
	if actions.PagerdutyAutomationActions != nil {
		flattenedAction["pagerduty_automation_actions"] = flattenServicePathPagerDutyAutomationActions(actions.PagerdutyAutomationActions)
	}
	if actions.AutomationActions != nil {
		flattenedAction["automation_actions"] = flattenServicePathAutomationActions(actions.AutomationActions)
	}

	actionsMap = append(actionsMap, flattenedAction)

	return actionsMap
}

func flattenServicePathPagerDutyAutomationActions(v []*pagerduty.EventOrchestrationPathPagerdutyAutomationAction) []interface{} {
	var result []interface{}

	for _, i := range v {
		pdaa := map[string]string{
			"action_id": i.ActionId,
		}

		result = append(result, pdaa)
	}

	return result
}

func flattenServicePathAutomationActions(v []*pagerduty.EventOrchestrationPathAutomationAction) []interface{} {
	var result []interface{}

	for _, i := range v {
		pdaa := map[string]interface{}{
			"name":       i.Name,
			"url":        i.Url,
			"auto_send":  i.AutoSend,
			"headers":    flattenServicePathAutomationActionObjects(i.Headers),
			"parameters": flattenServicePathAutomationActionObjects(i.Parameters),
		}

		result = append(result, pdaa)
	}

	return result
}

func flattenServicePathAutomationActionObjects(v []*pagerduty.EventOrchestrationPathAutomationActionObject) []interface{} {
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
