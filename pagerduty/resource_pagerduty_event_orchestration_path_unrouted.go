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
		Delete: resourcePagerDutyEventOrchestrationPathUnroutedUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //TODO: resourcePagerDutyEventOrchestrationPathImport
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
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"self": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"self": {
				Type:     schema.TypeString,
				Optional: true,
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
							Required: true, // even if there are no rules, API returns rules as an empty list
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true, // If the start set has no rules, empty list is returned by API for rules.
										// TODO: there is a validation on id
									},
									"label": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expression": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"actions": {
										Type:     schema.TypeList,
										Required: true, // even if there are no actions, API returns actions as an empty list
										MaxItems: 1,    //TODO check if this is valid
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
				Required: true, //if not supplied, API creates it
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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

// func buildEventOrchestrationUnroutedStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
// 	orchestrationPath := &pagerduty.EventOrchestrationPath{
// 		Type: d.Get("type").(string),
// 		//Self: d.Get("self").(string),
// 	}

// 	// if attr, ok := d.GetOk("description"); ok {
// 	// 	orchestration.Description = attr.(string)
// 	// }

// 	// if attr, ok := d.GetOk("team"); ok {
// 	// 	orchestration.Team = expandOrchestrationTeam(attr)
// 	// }
// 	return orchestrationPath
// }

func resourcePagerDutyEventOrchestrationPathUnroutedRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	// TODO: Check migration to RetryContext func
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		path := buildUnroutedPathStruct(d)
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "unrouted", path.Parent.ID)

		if unroutedPath, _, err := client.EventOrchestrationPaths.Get(path.Parent.ID, path.Type); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if unroutedPath != nil {
			d.SetId(path.Parent.ID)
			// d.Set("type", unroutedPath.Type)
			// d.Set("self", path.Parent.Self+"/"+unroutedPath.Type)
			if unroutedPath.Sets != nil {
				d.Set("sets", flattenSets(unroutedPath.Sets))
			}
		}
		return nil
	})

}

func resourcePagerDutyEventOrchestrationPathUnroutedCreate(d *schema.ResourceData, meta interface{}) error {
	// client, err := meta.(*Config).Client()
	// if err != nil {
	// 	return err
	// }

	// return resource.Retry(2*time.Minute, func() *resource.RetryError {
	// 	unroutedPathStruct := buildUnroutedPathStruct(d)
	// 	log.Printf("[INFO] Reading PagerDuty EventOrchestrationPath of type: %s for orchestration: %s", "unrouted", unroutedPathStruct.Parent.ID)

	// 	if unroutedPath, _, err := client.EventOrchestrationPaths.Get(unroutedPathStruct.Parent.ID, unroutedPathStruct.Type); err != nil {
	// 		time.Sleep(2 * time.Second)
	// 		return resource.RetryableError(err)
	// 	} else if unroutedPath != nil {
	// 		d.SetId(unroutedPathStruct.Parent.ID)
	// 		d.Set("type", unroutedPath.Type)
	// 	}
	// 	return nil
	// })
	// return resourcePagerDutyEventOrchestrationPathUnroutedRead(d, meta)
	return resourcePagerDutyEventOrchestrationPathUnroutedUpdate(d, meta)
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

func resourcePagerDutyEventOrchestrationPathUnroutedDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
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
			d.Set("sets", flattenSets(unroutedPath.Sets))
		}

		//TODO: figure out rule ordering
		// else if rule.Position != nil && *updatedUnroutedPath.Position != *rule.Position && rule.CatchAll != true {
		// 	log.Printf("[INFO] PagerDuty ruleset rule %s position %d needs to be %d", updatedUnroutedPath.ID, *updatedUnroutedPath.Position, *rule.Position)
		// 	return resource.RetryableError(fmt.Errorf("Error updating ruleset rule %s position %d needs to be %d", updatedUnroutedPath.ID, *updatedUnroutedPath.Position, *rule.Position))
		// }
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return nil
}

func buildUnroutedPathStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	orchPath := &pagerduty.EventOrchestrationPath{
		Type: d.Get("type").(string),
	}

	if attr, ok := d.GetOk("parent"); ok {
		orchPath.Parent = expandOrchestrationPathParent(attr)
	}

	return orchPath
}

func buildUnroutedPathStructForUpdate(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {

	// get the path-parent
	orchPath := &pagerduty.EventOrchestrationPath{}

	if attr, ok := d.GetOk("parent"); ok {
		orchPath.Parent = expandOrchestrationPathParent(attr)
	}

	// build other props
	// if attr, ok := d.GetOk("self"); ok {
	// 	orchPath.Self = attr.(string)
	// } else {
	// 	orchPath.Self = orchPath.Parent.Self + "/" + orchPath.Type
	// }

	if attr, ok := d.GetOk("sets"); ok {
		orchPath.Sets = expandSets(attr.([]interface{}))
	}

	return orchPath
}

func expandOrchestrationPathParent(v interface{}) *pagerduty.EventOrchestrationPathReference {
	var parent *pagerduty.EventOrchestrationPathReference
	p := v.([]interface{})[0].(map[string]interface{})
	parent = &pagerduty.EventOrchestrationPathReference{
		ID:   p["id"].(string),
		Type: p["type"].(string),
		Self: p["self"].(string),
	}

	return parent
}

func expandSets(v interface{}) []*pagerduty.EventOrchestrationPathSet {
	var sets []*pagerduty.EventOrchestrationPathSet

	for _, set := range v.([]interface{}) {
		s := set.(map[string]interface{})

		orchPathSet := &pagerduty.EventOrchestrationPathSet{
			ID:    s["id"].(string),
			Rules: expandRules(s["rules"]),
		}

		sets = append(sets, orchPathSet)
	}

	return sets
}

func expandRules(v interface{}) []*pagerduty.EventOrchestrationPathRule {
	var rules []*pagerduty.EventOrchestrationPathRule

	for _, rule := range v.([]interface{}) {
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
		// TODO: check if this is the right way to check on actions
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
	var conditions []*pagerduty.EventOrchestrationPathRuleCondition

	for _, cond := range v.([]interface{}) {
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

func flattenSets(orchPathSets []*pagerduty.EventOrchestrationPathSet) []interface{} {
	var flattenedSets []interface{}

	for _, set := range orchPathSets {
		flattenedSet := map[string]interface{}{
			"id":    set.ID,
			"rules": flattenRules(set.Rules),
		}
		flattenedSets = append(flattenedSets, flattenedSet)
	}
	return flattenedSets
}

func flattenRules(rules []*pagerduty.EventOrchestrationPathRule) []interface{} {
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
