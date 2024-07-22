package pagerduty

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

var eventOrchestrationPathGlobalCatchAllActionsSchema = map[string]*schema.Schema{
	"drop_event": {
		Type:     schema.TypeBool,
		Optional: true,
	},
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
	"automation_action": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: eventOrchestrationAutomationActionSchema,
		},
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
	"incident_custom_field_update": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: eventOrchestrationIncidentCustomFieldsObjectSchema,
		},
	},
	"escalation_policy": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

var eventOrchestrationPathGlobalRuleActionsSchema = buildEventOrchestrationPathGlobalRuleActionsSchema()

func buildEventOrchestrationPathGlobalRuleActionsSchema() map[string]*schema.Schema {
	a := eventOrchestrationPathGlobalCatchAllActionsSchema
	a["route_to"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return a
}

func resourcePagerDutyEventOrchestrationPathGlobal() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyEventOrchestrationPathGlobalRead,
		CreateContext: resourcePagerDutyEventOrchestrationPathGlobalCreate,
		UpdateContext: resourcePagerDutyEventOrchestrationPathGlobalUpdate,
		DeleteContext: resourcePagerDutyEventOrchestrationPathGlobalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePagerDutyEventOrchestrationPathGlobalImport,
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
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: eventOrchestrationPathGlobalRuleActionsSchema,
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
								Schema: eventOrchestrationPathGlobalCatchAllActionsSchema,
							},
						},
					},
				},
			},
		},
	}
}

func resourcePagerDutyEventOrchestrationPathGlobalRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		id := d.Id()
		t := "global"
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type %s for orchestration: %s", t, id)

		if path, _, err := client.EventOrchestrationPaths.GetContext(ctx, d.Id(), t); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
		} else if path != nil {
			setEventOrchestrationPathGlobalProps(d, path)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diags
}

func resourcePagerDutyEventOrchestrationPathGlobalCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourcePagerDutyEventOrchestrationPathGlobalUpdate(ctx, d, meta)
}

func resourcePagerDutyEventOrchestrationPathGlobalUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	payload := buildGlobalPathStruct(d)
	var globalPath *pagerduty.EventOrchestrationPath
	var warnings []*pagerduty.EventOrchestrationPathWarning

	log.Printf("[INFO] Creating PagerDuty Event Orchestration Global Path: %s", payload.Parent.ID)

	retryErr := retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
		if response, _, err := client.EventOrchestrationPaths.UpdateContext(ctx, payload.Parent.ID, "global", payload); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else if response != nil {
			d.SetId(response.OrchestrationPath.Parent.ID)
			globalPath = response.OrchestrationPath
			warnings = response.Warnings
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	setEventOrchestrationPathGlobalProps(d, globalPath)

	return convertEventOrchestrationPathWarningsToDiagnostics(warnings, diags)
}

func resourcePagerDutyEventOrchestrationPathGlobalDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	// In order to delete a Global Orchestration an empty orchestration path
	// config should be sent as an update.
	emptyPath := emptyOrchestrationPathStructBuilder("global_orchestration")
	orchestrationID := d.Get("event_orchestration").(string)

	log.Printf("[INFO] Deleting PagerDuty Global Event Orchestration Path: %s", orchestrationID)

	retryErr := retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
		if _, _, err := client.EventOrchestrationPaths.UpdateContext(ctx, orchestrationID, "global", emptyPath); err != nil {
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

func resourcePagerDutyEventOrchestrationPathGlobalImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	orchestrationID := d.Id()
	_, _, err = client.EventOrchestrationPaths.GetContext(ctx, orchestrationID, "global")

	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(orchestrationID)
	d.Set("event_orchestration", orchestrationID)

	return []*schema.ResourceData{d}, nil
}

func buildGlobalPathStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	return &pagerduty.EventOrchestrationPath{
		Parent: &pagerduty.EventOrchestrationPathReference{
			ID: d.Get("event_orchestration").(string),
		},
		Sets:     expandGlobalPathSets(d.Get("set")),
		CatchAll: expandGlobalPathCatchAll(d.Get("catch_all")),
	}
}

func expandGlobalPathSets(v interface{}) []*pagerduty.EventOrchestrationPathSet {
	var sets []*pagerduty.EventOrchestrationPathSet

	for _, set := range v.([]interface{}) {
		s := set.(map[string]interface{})

		orchPathSet := &pagerduty.EventOrchestrationPathSet{
			ID:    s["id"].(string),
			Rules: expandGlobalPathRules(s["rule"].(interface{})),
		}

		sets = append(sets, orchPathSet)
	}

	return sets
}

func expandGlobalPathRules(v interface{}) []*pagerduty.EventOrchestrationPathRule {
	items := v.([]interface{})
	rules := []*pagerduty.EventOrchestrationPathRule{}

	for _, rule := range items {
		r := rule.(map[string]interface{})

		ruleInSet := &pagerduty.EventOrchestrationPathRule{
			ID:         r["id"].(string),
			Label:      r["label"].(string),
			Disabled:   r["disabled"].(bool),
			Conditions: expandEventOrchestrationPathConditions(r["condition"]),
			Actions:    expandGlobalPathActions(r["actions"]),
		}

		rules = append(rules, ruleInSet)
	}
	return rules
}

func expandGlobalPathCatchAll(v interface{}) *pagerduty.EventOrchestrationPathCatchAll {
	catchAll := new(pagerduty.EventOrchestrationPathCatchAll)

	for _, ca := range v.([]interface{}) {
		if ca != nil {
			am := ca.(map[string]interface{})
			catchAll.Actions = expandGlobalPathActions(am["actions"])
		}
	}

	return catchAll
}

func expandGlobalPathActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = &pagerduty.EventOrchestrationPathRuleActions{
		AutomationActions:          []*pagerduty.EventOrchestrationPathAutomationAction{},
		Variables:                  []*pagerduty.EventOrchestrationPathActionVariables{},
		Extractions:                []*pagerduty.EventOrchestrationPathActionExtractions{},
		IncidentCustomFieldUpdates: []*pagerduty.EventOrchestrationPathIncidentCustomFieldUpdate{},
	}

	for _, i := range v.([]interface{}) {
		if i == nil {
			continue
		}
		a := i.(map[string]interface{})

		actions.DropEvent = a["drop_event"].(bool)
		actions.RouteTo = a["route_to"].(string)
		actions.Suppress = a["suppress"].(bool)
		actions.Suspend = intTypeToIntPtr(a["suspend"].(int))
		actions.Priority = a["priority"].(string)
		actions.EscalationPolicy = stringTypeToStringPtr(a["escalation_policy"].(string))
		actions.Annotate = a["annotate"].(string)
		actions.Severity = a["severity"].(string)
		actions.EventAction = a["event_action"].(string)
		actions.AutomationActions = expandEventOrchestrationPathAutomationActions(a["automation_action"])
		actions.Variables = expandEventOrchestrationPathVariables(a["variable"])
		actions.Extractions = expandEventOrchestrationPathExtractions(a["extraction"])
		actions.IncidentCustomFieldUpdates = expandEventOrchestrationPathIncidentCustomFields(a["incident_custom_field_update"])
	}

	return actions
}

func setEventOrchestrationPathGlobalProps(d *schema.ResourceData, p *pagerduty.EventOrchestrationPath) error {
	d.SetId(p.Parent.ID)
	d.Set("event_orchestration", p.Parent.ID)
	d.Set("set", flattenGlobalPathSets(p.Sets))
	d.Set("catch_all", flattenGlobalPathCatchAll(p.CatchAll))
	return nil
}

func flattenGlobalPathSets(orchPathSets []*pagerduty.EventOrchestrationPathSet) []interface{} {
	var flattenedSets []interface{}

	for _, set := range orchPathSets {
		flattenedSet := map[string]interface{}{
			"id":   set.ID,
			"rule": flattenGlobalPathRules(set.Rules),
		}
		flattenedSets = append(flattenedSets, flattenedSet)
	}
	return flattenedSets
}

func flattenGlobalPathCatchAll(catchAll *pagerduty.EventOrchestrationPathCatchAll) []map[string]interface{} {
	var caMap []map[string]interface{}

	c := make(map[string]interface{})

	c["actions"] = flattenGlobalPathActions(catchAll.Actions)
	caMap = append(caMap, c)

	return caMap
}

func flattenGlobalPathRules(rules []*pagerduty.EventOrchestrationPathRule) []interface{} {
	var flattenedRules []interface{}

	for _, rule := range rules {
		flattenedRule := map[string]interface{}{
			"id":        rule.ID,
			"label":     rule.Label,
			"disabled":  rule.Disabled,
			"condition": flattenEventOrchestrationPathConditions(rule.Conditions),
			"actions":   flattenGlobalPathActions(rule.Actions),
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	return flattenedRules
}

func flattenGlobalPathActions(actions *pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	flattenedAction := map[string]interface{}{
		"drop_event":        actions.DropEvent,
		"route_to":          actions.RouteTo,
		"severity":          actions.Severity,
		"event_action":      actions.EventAction,
		"suppress":          actions.Suppress,
		"suspend":           actions.Suspend,
		"priority":          actions.Priority,
		"annotate":          actions.Annotate,
		"escalation_policy": stringPtrToStringType(actions.EscalationPolicy),
	}

	if actions.Variables != nil {
		flattenedAction["variable"] = flattenEventOrchestrationPathVariables(actions.Variables)
	}
	if actions.Extractions != nil {
		flattenedAction["extraction"] = flattenEventOrchestrationPathExtractions(actions.Extractions)
	}
	if actions.AutomationActions != nil {
		flattenedAction["automation_action"] = flattenEventOrchestrationAutomationActions(actions.AutomationActions)
	}
	if actions.IncidentCustomFieldUpdates != nil {
		flattenedAction["incident_custom_field_update"] = flattenEventOrchestrationIncidentCustomFieldUpdates(actions.IncidentCustomFieldUpdates)
	}

	actionsMap = append(actionsMap, flattenedAction)

	return actionsMap
}
