package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationPathUnrouted() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePagerDutyEventOrchestrationPathUnroutedRead,
		Create: resourcePagerDutyEventOrchestrationPathUnroutedCreate,
		Update: resourcePagerDutyEventOrchestrationPathUnroutedUpdate,
		Delete: resourcePagerDutyEventOrchestrationPathUnroutedDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyEventOrchestrationPathUnroutedImport,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
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
				MinItems: 1, // An Unrouted Orchestration must contain at least a "start" set
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
										Required: true, // even if there are no actions, API returns actions as an empty list
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"route_to": {
													Type:     schema.TypeString,
													Optional: true, // If there is only start set we don't need route_to
													//TODO: validate func, The ID of a Set from this Unrouted Orchestration whose rules you also want to use with event that match this rule.
												},
												"severity": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validateValueFunc([]string{
														"info",
														"error",
														"warning",
														"critical",
													}),
												},
												"event_action": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validateValueFunc([]string{
														"trigger",
														"resolve",
													}),
												},
												"variables": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
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
														}}},
												"extractions": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"target": {
																Type:     schema.TypeString,
																Required: true,
															},
															"template": {
																Type:     schema.TypeString,
																Required: true,
															},
														}},
												},
											},
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
								Schema: map[string]*schema.Schema{
									"suppress": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateValueFunc([]string{
											"info",
											"error",
											"warning",
											"critical",
										}),
									},
									"event_action": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateValueFunc([]string{
											"trigger",
											"resolve",
										}),
									},
									"variables": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
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
											}}},
									"extractions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target": {
													Type:     schema.TypeString,
													Required: true,
												},
												"template": {
													Type:     schema.TypeString,
													Required: true,
												},
											}},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourcePagerDutyEventOrchestrationPathUnroutedRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {

		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "unrouted", d.Id())

		if unroutedPath, _, err := client.EventOrchestrationPaths.Get(d.Id(), "unrouted"); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if unroutedPath != nil {

			if unroutedPath.Sets != nil {
				d.Set("sets", flattenUnroutedSets(unroutedPath.Sets))
			}

			if unroutedPath.CatchAll != nil {
				d.Set("catch_all", flattenUnroutedCatchAll(unroutedPath.CatchAll))
			}
		}
		return nil
	})

}

// EventOrchestrationPath cannot be created, use update to add / edit / remove rules and sets
func resourcePagerDutyEventOrchestrationPathUnroutedCreate(d *schema.ResourceData, meta interface{}) error {
	return resourcePagerDutyEventOrchestrationPathUnroutedUpdate(d, meta)
}

// EventOrchestrationPath cannot be deleted, use update to add / edit / remove rules and sets
func resourcePagerDutyEventOrchestrationPathUnroutedDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourcePagerDutyEventOrchestrationPathUnroutedUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	updatePath := buildUnroutedPathStructForUpdate(d)

	log.Printf("[INFO] Updating PagerDuty EventOrchestrationPath of type: %s for orchestration: %s", "unrouted", updatePath.Parent.ID)

	return performUnroutedPathUpdate(d, updatePath, client)
}

func performUnroutedPathUpdate(d *schema.ResourceData, unroutedPath *pagerduty.EventOrchestrationPath, client *pagerduty.Client) error {
	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		updatedPath, _, err := client.EventOrchestrationPaths.Update(unroutedPath.Parent.ID, "unrouted", unroutedPath)
		if err != nil {
			return resource.RetryableError(err)
		}
		if updatedPath == nil {
			return resource.NonRetryableError(fmt.Errorf("No Event Orchestration Unrouted found."))
		}
		// set props
		d.SetId(unroutedPath.Parent.ID)
		if unroutedPath.Sets != nil {
			d.Set("sets", flattenUnroutedSets(unroutedPath.Sets))
		}
		if updatedPath.CatchAll != nil {
			d.Set("catch_all", flattenUnroutedCatchAll(updatedPath.CatchAll))
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return nil
}

func buildUnroutedPathParent(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	orchPath := &pagerduty.EventOrchestrationPath{}

	if attr, ok := d.GetOk("parent"); ok {
		orchPath.Parent = expandOrchestrationPathUnroutedParent(attr)
	}

	return orchPath
}

func buildUnroutedPathStructForUpdate(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {

	orchPath := buildUnroutedPathParent(d)

	if attr, ok := d.GetOk("parent"); ok {
		orchPath.Parent = expandOrchestrationPathUnroutedParent(attr)
	}

	if attr, ok := d.GetOk("sets"); ok {
		orchPath.Sets = expandUnroutedSets(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("catch_all"); ok {
		orchPath.CatchAll = expandUnroutedCatchAll(attr.([]interface{}))
	}

	return orchPath
}

func expandOrchestrationPathUnroutedParent(v interface{}) *pagerduty.EventOrchestrationPathReference {
	var parent *pagerduty.EventOrchestrationPathReference
	p := v.([]interface{})[0].(map[string]interface{})
	parent = &pagerduty.EventOrchestrationPathReference{
		ID: p["id"].(string),
	}

	return parent
}

func expandUnroutedSets(v interface{}) []*pagerduty.EventOrchestrationPathSet {
	var sets []*pagerduty.EventOrchestrationPathSet

	for _, set := range v.([]interface{}) {
		s := set.(map[string]interface{})

		orchPathSet := &pagerduty.EventOrchestrationPathSet{
			ID:    s["id"].(string),
			Rules: expandUnroutedRules(s["rules"]),
		}

		sets = append(sets, orchPathSet)
	}

	return sets
}

func expandUnroutedRules(v interface{}) []*pagerduty.EventOrchestrationPathRule {
	items := v.([]interface{})
	rules := []*pagerduty.EventOrchestrationPathRule{}

	for _, rule := range items {
		r := rule.(map[string]interface{})

		ruleInSet := &pagerduty.EventOrchestrationPathRule{
			ID:         r["id"].(string),
			Label:      r["label"].(string),
			Disabled:   r["disabled"].(bool),
			Conditions: expandUnroutedConditions(r["conditions"]),
			Actions:    expandUnroutedActions(r["actions"].([]interface{})),
		}

		rules = append(rules, ruleInSet)
	}
	return rules
}

func expandUnroutedActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = new(pagerduty.EventOrchestrationPathRuleActions)
	for _, ai := range v.([]interface{}) {
		if ai != nil {
			am := ai.(map[string]interface{})
			actions.RouteTo = am["route_to"].(string)
			actions.Severity = am["severity"].(string)
			actions.EventAction = am["event_action"].(string)
			actions.Variables = expandUnroutedActionsVariables(am["variables"])
			actions.Extractions = expandUnroutedActionsExtractions(am["extractions"])
		}
	}

	return actions
}

func expandUnroutedConditions(v interface{}) []*pagerduty.EventOrchestrationPathRuleCondition {
	items := v.([]interface{})
	conditions := []*pagerduty.EventOrchestrationPathRuleCondition{}

	for _, cond := range items {
		c := cond.(map[string]interface{})

		cx := &pagerduty.EventOrchestrationPathRuleCondition{
			Expression: c["expression"].(string),
		}

		conditions = append(conditions, cx)
	}

	return conditions
}

func expandUnroutedActionsExtractions(v interface{}) []*pagerduty.EventOrchestrationPathActionExtractions {
	var unroutedExtractions []*pagerduty.EventOrchestrationPathActionExtractions

	for _, eai := range v.([]interface{}) {
		ea := eai.(map[string]interface{})
		ext := &pagerduty.EventOrchestrationPathActionExtractions{
			Target:   ea["target"].(string),
			Template: ea["template"].(string),
		}
		unroutedExtractions = append(unroutedExtractions, ext)
	}
	return unroutedExtractions
}

func expandUnroutedActionsVariables(v interface{}) []*pagerduty.EventOrchestrationPathActionVariables {
	var unroutedVariables []*pagerduty.EventOrchestrationPathActionVariables

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})

		unroutedVar := &pagerduty.EventOrchestrationPathActionVariables{
			Name:  rer["name"].(string),
			Path:  rer["path"].(string),
			Type:  rer["type"].(string),
			Value: rer["value"].(string),
		}

		unroutedVariables = append(unroutedVariables, unroutedVar)
	}

	return unroutedVariables
}

func expandUnroutedCatchAll(v interface{}) *pagerduty.EventOrchestrationPathCatchAll {
	var catchAll = new(pagerduty.EventOrchestrationPathCatchAll)

	for _, ca := range v.([]interface{}) {
		if ca != nil {
			am := ca.(map[string]interface{})
			catchAll.Actions = expandUnroutedCatchAllActions(am["actions"])
		}
	}

	return catchAll
}

func expandUnroutedCatchAllActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = new(pagerduty.EventOrchestrationPathRuleActions)
	for _, ai := range v.([]interface{}) {
		if ai != nil {
			am := ai.(map[string]interface{})
			actions.Severity = am["severity"].(string)
			actions.EventAction = am["event_action"].(string)
			actions.Variables = expandUnroutedActionsVariables(am["variables"])
			actions.Extractions = expandUnroutedActionsExtractions(am["extractions"])
		}
	}

	return actions
}

func flattenUnroutedSets(orchPathSets []*pagerduty.EventOrchestrationPathSet) []interface{} {
	var flattenedSets []interface{}

	for _, set := range orchPathSets {
		flattenedSet := map[string]interface{}{
			"id":    set.ID,
			"rules": flattenUnroutedRules(set.Rules),
		}
		flattenedSets = append(flattenedSets, flattenedSet)
	}
	return flattenedSets
}

func flattenUnroutedRules(rules []*pagerduty.EventOrchestrationPathRule) []interface{} {
	var flattenedRules []interface{}

	for _, rule := range rules {
		flattenedRule := map[string]interface{}{
			"id":         rule.ID,
			"label":      rule.Label,
			"disabled":   rule.Disabled,
			"conditions": flattenUnroutedConditions(rule.Conditions),
			"actions":    flattenUnroutedActions(rule.Actions),
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	return flattenedRules
}

func flattenUnroutedConditions(conditions []*pagerduty.EventOrchestrationPathRuleCondition) []interface{} {
	var flattendConditions []interface{}

	for _, condition := range conditions {
		flattendCondition := map[string]interface{}{
			"expression": condition.Expression,
		}
		flattendConditions = append(flattendConditions, flattendCondition)
	}

	return flattendConditions
}

func flattenUnroutedActions(actions *pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	flattenedAction := map[string]interface{}{
		"route_to":     actions.RouteTo,
		"severity":     actions.Severity,
		"event_action": actions.EventAction,
	}

	if actions.Variables != nil {
		flattenedAction["variables"] = flattenUnroutedActionsVariables(actions.Variables)
	}
	if actions.Variables != nil {
		flattenedAction["extractions"] = flattenUnroutedActionsExtractions(actions.Extractions)
	}

	actionsMap = append(actionsMap, flattenedAction)

	return actionsMap
}

func flattenUnroutedActionsVariables(v []*pagerduty.EventOrchestrationPathActionVariables) []interface{} {
	var flatVariablesList []interface{}

	for _, s := range v {
		flatVariable := map[string]interface{}{
			"name":  s.Name,
			"path":  s.Path,
			"type":  s.Type,
			"value": s.Value,
		}
		flatVariablesList = append(flatVariablesList, flatVariable)
	}
	return flatVariablesList
}

func flattenUnroutedActionsExtractions(e []*pagerduty.EventOrchestrationPathActionExtractions) []interface{} {
	var flatExtractionsList []interface{}

	for _, s := range e {
		flatExtraction := map[string]interface{}{
			"target":   s.Target,
			"template": s.Template,
		}
		flatExtractionsList = append(flatExtractionsList, flatExtraction)
	}
	return flatExtractionsList
}

func flattenUnroutedCatchAll(catchAll *pagerduty.EventOrchestrationPathCatchAll) []map[string]interface{} {
	var caMap []map[string]interface{}

	c := make(map[string]interface{})

	c["actions"] = flattenUnroutedCatchAllActions(catchAll.Actions)
	caMap = append(caMap, c)

	return caMap
}

func flattenUnroutedCatchAllActions(actions *pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	flattenedAction := map[string]interface{}{
		"severity":     actions.Severity,
		"event_action": actions.EventAction,
		"suppress":     actions.Suppress, // By default suppress is set to "true" by API for unrouted
	}

	if actions.Variables != nil {
		flattenedAction["variables"] = flattenUnroutedActionsVariables(actions.Variables)
	}
	if actions.Variables != nil {
		flattenedAction["extractions"] = flattenUnroutedActionsExtractions(actions.Extractions)
	}

	actionsMap = append(actionsMap, flattenedAction)

	return actionsMap
}

func resourcePagerDutyEventOrchestrationPathUnroutedImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	// given an orchestration ID import the unrouted orchestration path
	orchestrationID := d.Id()
	pathType := "unrouted"
	_, _, err = client.EventOrchestrationPaths.Get(orchestrationID, pathType)

	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(orchestrationID)

	return []*schema.ResourceData{d}, nil
}
