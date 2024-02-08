package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationPathUnrouted() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyEventOrchestrationPathUnroutedRead,
		CreateContext: resourcePagerDutyEventOrchestrationPathUnroutedCreate,
		UpdateContext: resourcePagerDutyEventOrchestrationPathUnroutedUpdate,
		DeleteContext: resourcePagerDutyEventOrchestrationPathUnroutedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePagerDutyEventOrchestrationPathUnroutedImport,
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

func resourcePagerDutyEventOrchestrationPathUnroutedRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "unrouted", d.Id())

		if unroutedPath, _, err := client.EventOrchestrationPaths.GetContext(ctx, d.Id(), "unrouted"); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
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

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diags
}

// EventOrchestrationPath cannot be created, use update to add / edit / remove rules and sets
func resourcePagerDutyEventOrchestrationPathUnroutedCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourcePagerDutyEventOrchestrationPathUnroutedUpdate(ctx, d, meta)
}

// EventOrchestrationPath cannot be deleted, use update to add / edit / remove rules and sets
func resourcePagerDutyEventOrchestrationPathUnroutedDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	// In order to delete an Unrouted Orchestration an empty orchestration path
	// config should be sent as an update.
	emptyPath := emptyOrchestrationPathStructBuilder("unrouted")
	orchestrationID := d.Get("event_orchestration").(string)

	log.Printf("[INFO] Deleting PagerDuty Unrouted Event Orchestration Path: %s", orchestrationID)

	retryErr := retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
		if _, _, err := client.EventOrchestrationPaths.UpdateContext(ctx, orchestrationID, "unrouted", emptyPath); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	d.SetId("")
	return nil
}

func resourcePagerDutyEventOrchestrationPathUnroutedUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	unroutedPath := buildUnroutedPathStructForUpdate(d)
	var warnings []*pagerduty.EventOrchestrationPathWarning

	log.Printf("[INFO] Updating PagerDuty EventOrchestrationPath of type: %s for orchestration: %s", "unrouted", unroutedPath.Parent.ID)

	retryErr := retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
		response, _, err := client.EventOrchestrationPaths.UpdateContext(ctx, unroutedPath.Parent.ID, "unrouted", unroutedPath)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}

		if response == nil {
			return retry.NonRetryableError(fmt.Errorf("no event orchestration unrouted found"))
		}

		d.SetId(unroutedPath.Parent.ID)
		d.Set("event_orchestration", unroutedPath.Parent.ID)
		warnings = response.Warnings

		if unroutedPath.Sets != nil {
			d.Set("set", flattenUnroutedSets(unroutedPath.Sets))
		}

		if response.OrchestrationPath.CatchAll != nil {
			d.Set("catch_all", flattenUnroutedCatchAll(response.OrchestrationPath.CatchAll))
		}

		return nil
	})

	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	return convertEventOrchestrationPathWarningsToDiagnostics(warnings, diags)
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
	actions := &pagerduty.EventOrchestrationPathRuleActions{
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
	catchAll := new(pagerduty.EventOrchestrationPathCatchAll)

	for _, ca := range v.([]interface{}) {
		if ca != nil {
			am := ca.(map[string]interface{})
			catchAll.Actions = expandUnroutedCatchAllActions(am["actions"])
		}
	}

	return catchAll
}

func expandUnroutedCatchAllActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	actions := new(pagerduty.EventOrchestrationPathRuleActions)
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

func resourcePagerDutyEventOrchestrationPathUnroutedImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	// given an orchestration ID import the unrouted orchestration path
	orchestrationID := d.Id()
	pathType := "unrouted"
	_, _, err = client.EventOrchestrationPaths.GetContext(ctx, orchestrationID, pathType)

	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(orchestrationID)
	d.Set("event_orchestration", orchestrationID)

	return []*schema.ResourceData{d}, nil
}
