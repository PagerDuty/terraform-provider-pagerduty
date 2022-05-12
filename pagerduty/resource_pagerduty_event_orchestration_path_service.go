package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
	"log"
	"time"
)

var eventOrchestrationAutomationActionObject = map[string]*schema.Schema{
	"key": {
		Type:     schema.TypeString,
		Required: true,
	},
	"value": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var eventOrchestrationPathServiceCatchAllActions = map[string]*schema.Schema{
	// suppress
	// suspend
	// priority
	// annotate
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
						Schema: eventOrchestrationAutomationActionObject,
					},
				},
				"parameters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: eventOrchestrationAutomationActionObject,
					},
				},
			},
		},
	},
	// severity
	// event_action
	// variables
	// extractions
}

var eventOrchestrationPathServiceRuleActions = buildEventOrchestrationPathServiceRuleActions()

func buildEventOrchestrationPathServiceRuleActions() map[string]*schema.Schema {
	a := eventOrchestrationPathServiceCatchAllActions
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
			State: schema.ImportStatePassthrough, //TODO: resourcePagerDutyEventOrchestrationPathServiceImport
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: PagerDutyEventOrchestrationPathParent,
				},
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
											Schema: eventOrchestrationPathServiceRuleActions,
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
			// "catch_all": {
			// 	Type:     schema.TypeList,
			// 	MaxItems: 1,
			// 	Required: true, // TODO: figure out how to make this optional with a default value
			// 	Elem: &schema.Resource{
			// 		Schema: eventOrchestrationPathServiceCatchAllActions,
			// 	},
			// },
		},
	}
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

func resourcePagerDutyEventOrchestrationPathServiceDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func buildServicePathStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	return &pagerduty.EventOrchestrationPath{
		// TODO: use shared method
		Parent: &pagerduty.EventOrchestrationPathReference{
			ID: d.Get("parent.0.id").(string),
		},
		Sets: expandServicePathSets(d.Get("sets")),
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

func expandServicePathActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = new(pagerduty.EventOrchestrationPathRuleActions)

	for _, i := range v.([]interface{}) {
		if i == nil {
			continue
		}
		a := i.(map[string]interface{})
		// TODO:
		actions.RouteTo = a["route_to"].(string)
		// suppress
		// suspend
		// priority
		// annotate
		actions.PagerdutyAutomationActions = expandServicePathPagerDutyAutomationActions(a["pagerduty_automation_actions"])
		actions.AutomationActions = expandServicePathAutomationActions(a["automation_actions"])
		// severity
		// event_action
		// variables
		// extractions
	}

	return actions
}

func expandServicePathPagerDutyAutomationActions(v interface{}) []*pagerduty.EventOrchestrationPathPagerdutyAutomationAction {
	var result []*pagerduty.EventOrchestrationPathPagerdutyAutomationAction

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
	var result []*pagerduty.EventOrchestrationPathAutomationAction

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
	var result []*pagerduty.EventOrchestrationPathAutomationActionObject

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
	d.Set("type", p.Type)
	d.Set("parent", flattenServicePathParent(p.Parent))
	// TODO: see if we can reuse expand functions for all orch path sets.
	// Maybe pass in the rule actions and catch-all rule actions expanding function?
	d.Set("sets", flattenServicePathSets(p.Sets))
	return nil
}

func flattenServicePathParent(p *pagerduty.EventOrchestrationPathReference) []interface{} {
	var parent = map[string]interface{}{
		"id":   p.ID,
		"type": p.Type,
		"self": p.Self,
	}

	return []interface{}{parent}
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

	am := make(map[string]interface{})
	am["route_to"] = actions.RouteTo
	am["pagerduty_automation_actions"] = flattenServicePathPagerDutyAutomationActions(actions.PagerdutyAutomationActions)
	am["automation_actions"] = flattenServicePathAutomationActions(actions.AutomationActions)
	actionsMap = append(actionsMap, am)

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
