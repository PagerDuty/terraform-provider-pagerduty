package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationPathRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyEventOrchestrationPathRouterRead,
		CreateContext: resourcePagerDutyEventOrchestrationPathRouterCreate,
		UpdateContext: resourcePagerDutyEventOrchestrationPathRouterUpdate,
		DeleteContext: resourcePagerDutyEventOrchestrationPathRouterDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyEventOrchestrationPathRouterImport,
		},
		Schema: map[string]*schema.Schema{
			"event_orchestration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, // Router can only have 'start' set
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
										MaxItems: 1, //there can only be one action for router
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"route_to": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: func(v interface{}, key string) (warns []string, errs []error) {
														value := v.(string)
														if value == "unrouted" {
															errs = append(errs, fmt.Errorf("route_to within a set's rule has to be a Service ID. Got: %q", v))
														}
														return
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
									"route_to": {
										Type:     schema.TypeString,
										Required: true,
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

func resourcePagerDutyEventOrchestrationPathRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type %s for orchestration: %s", "router", d.Id())

		if routerPath, _, err := client.EventOrchestrationPaths.Get(d.Id(), "router"); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if routerPath != nil {
			d.Set("event_orchestration", routerPath.Parent.ID)

			if routerPath.Sets != nil {
				d.Set("set", flattenSets(routerPath.Sets))
			}

			if routerPath.CatchAll != nil {
				d.Set("catch_all", flattenCatchAll(routerPath.CatchAll))
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
func resourcePagerDutyEventOrchestrationPathRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourcePagerDutyEventOrchestrationPathRouterUpdate(ctx, d, meta)
}

func resourcePagerDutyEventOrchestrationPathRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}

func resourcePagerDutyEventOrchestrationPathRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	routerPath := buildRouterPathStructForUpdate(d)
	var warnings []*pagerduty.EventOrchestrationPathWarning

	log.Printf("[INFO] Updating PagerDuty Event Orchestration Path of type %s for orchestration: %s", "router", routerPath.Parent.ID)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		response, _, err := client.EventOrchestrationPaths.Update(routerPath.Parent.ID, "router", routerPath)
		if err != nil {
			return resource.RetryableError(err)
		}
		if response == nil {
			return resource.NonRetryableError(fmt.Errorf("No Event Orchestration Router found."))
		}
		d.SetId(routerPath.Parent.ID)
		d.Set("event_orchestration", routerPath.Parent.ID)
		warnings = response.Warnings

		if routerPath.Sets != nil {
			d.Set("set", flattenSets(routerPath.Sets))
		}
		if response.OrchestrationPath.CatchAll != nil {
			d.Set("catch_all", flattenCatchAll(response.OrchestrationPath.CatchAll))
		}
		return nil
	})

	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	return convertEventOrchestrationPathWarningsToDiagnostics(warnings, diags)
}

func buildRouterPathStructForUpdate(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {

	orchPath := &pagerduty.EventOrchestrationPath{
		Parent: &pagerduty.EventOrchestrationPathReference{
			ID: d.Get("event_orchestration").(string),
		},
	}

	if attr, ok := d.GetOk("set"); ok {
		orchPath.Sets = expandSets(attr)
	}

	if attr, ok := d.GetOk("catch_all"); ok {
		orchPath.CatchAll = expandCatchAll(attr)
	}

	return orchPath
}

func expandSets(v interface{}) []*pagerduty.EventOrchestrationPathSet {
	var sets []*pagerduty.EventOrchestrationPathSet

	for _, set := range v.([]interface{}) {
		s := set.(map[string]interface{})

		orchPathSet := &pagerduty.EventOrchestrationPathSet{
			ID:    s["id"].(string),
			Rules: expandRules(s["rule"]),
		}

		sets = append(sets, orchPathSet)
	}

	return sets
}

func expandRules(v interface{}) []*pagerduty.EventOrchestrationPathRule {
	items := v.([]interface{})
	rules := []*pagerduty.EventOrchestrationPathRule{}

	for _, rule := range items {
		r := rule.(map[string]interface{})

		ruleInSet := &pagerduty.EventOrchestrationPathRule{
			ID:         r["id"].(string),
			Label:      r["label"].(string),
			Disabled:   r["disabled"].(bool),
			Conditions: expandEventOrchestrationPathConditions(r["condition"]),
			Actions:    expandRouterActions(r["actions"]),
		}

		rules = append(rules, ruleInSet)
	}
	return rules
}

func expandRouterActions(v interface{}) *pagerduty.EventOrchestrationPathRuleActions {
	var actions = new(pagerduty.EventOrchestrationPathRuleActions)
	for _, ai := range v.([]interface{}) {
		am := ai.(map[string]interface{})
		actions.RouteTo = am["route_to"].(string)
	}

	return actions
}

func expandCatchAll(v interface{}) *pagerduty.EventOrchestrationPathCatchAll {
	var catchAll = new(pagerduty.EventOrchestrationPathCatchAll)

	for _, ca := range v.([]interface{}) {
		am := ca.(map[string]interface{})
		catchAll.Actions = expandRouterActions(am["actions"])
	}

	return catchAll
}

func flattenSets(orchPathSets []*pagerduty.EventOrchestrationPathSet) []interface{} {
	var flattenedSets []interface{}
	for _, set := range orchPathSets {
		flattenedSet := map[string]interface{}{
			"id":   set.ID,
			"rule": flattenRules(set.Rules),
		}
		flattenedSets = append(flattenedSets, flattenedSet)
	}
	return flattenedSets
}

func flattenRules(rules []*pagerduty.EventOrchestrationPathRule) []interface{} {
	var flattenedRules []interface{}

	for _, rule := range rules {
		flattenedRule := map[string]interface{}{
			"id":        rule.ID,
			"label":     rule.Label,
			"disabled":  rule.Disabled,
			"condition": flattenEventOrchestrationPathConditions(rule.Conditions),
			"actions":   flattenRouterActions(rule.Actions),
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	return flattenedRules
}

func flattenRouterActions(actions *pagerduty.EventOrchestrationPathRuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	am := make(map[string]interface{})
	am["route_to"] = actions.RouteTo
	actionsMap = append(actionsMap, am)
	return actionsMap
}

func flattenCatchAll(catchAll *pagerduty.EventOrchestrationPathCatchAll) []map[string]interface{} {
	var caMap []map[string]interface{}

	c := make(map[string]interface{})

	c["actions"] = flattenRouterActions(catchAll.Actions)
	caMap = append(caMap, c)

	return caMap
}

func resourcePagerDutyEventOrchestrationPathRouterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	// given an orchestration ID import the router orchestration path
	orchestrationID := d.Id()
	pathType := "router"
	_, _, err = client.EventOrchestrationPaths.Get(orchestrationID, pathType)

	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(orchestrationID)
	d.Set("event_orchestration", orchestrationID)

	return []*schema.ResourceData{d}, nil
}
