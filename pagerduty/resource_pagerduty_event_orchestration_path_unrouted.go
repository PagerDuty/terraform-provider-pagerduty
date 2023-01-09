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
		CustomizeDiff: checkExtractions,
		Schema: map[string]*schema.Schema{
			"event_orchestration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1, // An Unrouted Orchestration must contain at least a "start" set
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule": {
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
									"condition": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: eventOrchestrationPathConditionsSchema,
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
												},
												"severity": {
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: validateEventOrchestrationPathSeverity(),
												},
												"event_action": {
													Type:             schema.TypeString,
													Optional:         true,
													ValidateDiagFunc: validateEventOrchestrationPathEventAction(),
												},
												"variable": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: eventOrchestrationPathVariablesSchema,
													},
												},
												"extraction": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: eventOrchestrationPathExtractionsSchema,
													},
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
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"info",
											"error",
											"warning",
											"critical",
										}),
									},
									"event_action": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"trigger",
											"resolve",
										}),
									},
									"variable": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: eventOrchestrationPathVariablesSchema,
										},
									},
									"extraction": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: eventOrchestrationPathExtractionsSchema,
										},
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
				d.Set("set", flattenUnroutedSets(unroutedPath.Sets))
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
			return resource.NonRetryableError(fmt.Errorf("no event orchestration unrouted found"))
		}
		d.SetId(unroutedPath.Parent.ID)
		d.Set("event_orchestration", unroutedPath.Parent.ID)
		if unroutedPath.Sets != nil {
			d.Set("set", flattenUnroutedSets(unroutedPath.Sets))
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

func buildUnroutedPathStructForUpdate(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {

	orchPath := &pagerduty.EventOrchestrationPath{
		Parent: &pagerduty.EventOrchestrationPathReference{
			ID: d.Get("event_orchestration").(string),
		},
	}

	if attr, ok := d.GetOk("set"); ok {
		orchPath.Sets = expandUnroutedSets(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("catch_all"); ok {
		orchPath.CatchAll = expandUnroutedCatchAll(attr.([]interface{}))
	}

	return orchPath
}

func expandUnroutedSets(v interface{}) []*pagerduty.EventOrchestrationPathSet {
	var sets []*pagerduty.EventOrchestrationPathSet

	for _, set := range v.([]interface{}) {
		s := set.(map[string]interface{})

		orchPathSet := &pagerduty.EventOrchestrationPathSet{
			ID:    s["id"].(string),
			Rules: expandUnroutedRules(s["rule"]),
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
			Conditions: expandEventOrchestrationPathConditions(r["condition"]),
			Actions:    expandUnroutedActions(r["actions"]),
		}

		rules = append(rules, ruleInSet)
	}
	return rules
}

func expandUnroutedActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = &pagerduty.EventOrchestrationPathRuleActions{
		Variables:   []*pagerduty.EventOrchestrationPathActionVariables{},
		Extractions: []*pagerduty.EventOrchestrationPathActionExtractions{},
	}

	for _, ai := range v.([]interface{}) {
		if ai != nil {
			am := ai.(map[string]interface{})
			actions.RouteTo = am["route_to"].(string)
			actions.Severity = am["severity"].(string)
			actions.EventAction = am["event_action"].(string)
			actions.Variables = expandEventOrchestrationPathVariables(am["variable"])
			actions.Extractions = expandEventOrchestrationPathExtractions(am["extraction"])
		}
	}

	return actions
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
			actions.Variables = expandEventOrchestrationPathVariables(am["variable"])
			actions.Extractions = expandEventOrchestrationPathExtractions(am["extraction"])
		}
	}

	return actions
}

func flattenUnroutedSets(orchPathSets []*pagerduty.EventOrchestrationPathSet) []interface{} {
	var flattenedSets []interface{}

	for _, set := range orchPathSets {
		flattenedSet := map[string]interface{}{
			"id":   set.ID,
			"rule": flattenUnroutedRules(set.Rules),
		}
		flattenedSets = append(flattenedSets, flattenedSet)
	}
	return flattenedSets
}

func flattenUnroutedRules(rules []*pagerduty.EventOrchestrationPathRule) []interface{} {
	var flattenedRules []interface{}

	for _, rule := range rules {
		flattenedRule := map[string]interface{}{
			"id":        rule.ID,
			"label":     rule.Label,
			"disabled":  rule.Disabled,
			"condition": flattenEventOrchestrationPathConditions(rule.Conditions),
			"actions":   flattenUnroutedActions(rule.Actions),
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	return flattenedRules
}

func flattenUnroutedActions(actions *pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	flattenedAction := map[string]interface{}{
		"route_to":     actions.RouteTo,
		"severity":     actions.Severity,
		"event_action": actions.EventAction,
	}

	if actions.Variables != nil {
		flattenedAction["variable"] = flattenEventOrchestrationPathVariables(actions.Variables)
	}
	if actions.Extractions != nil {
		flattenedAction["extraction"] = flattenEventOrchestrationPathExtractions(actions.Extractions)
	}

	actionsMap = append(actionsMap, flattenedAction)

	return actionsMap
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
		flattenedAction["variable"] = flattenEventOrchestrationPathVariables(actions.Variables)
	}
	if actions.Variables != nil {
		flattenedAction["extraction"] = flattenEventOrchestrationPathExtractions(actions.Extractions)
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
	d.Set("event_orchestration", orchestrationID)

	return []*schema.ResourceData{d}, nil
}
