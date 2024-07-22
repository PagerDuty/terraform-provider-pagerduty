package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
			StateContext: resourcePagerDutyEventOrchestrationPathRouterImport,
		},
		CustomizeDiff: checkDynamicRoutingRule,
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
										MaxItems: 1, // there can only be one action for router
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dynamic_route_to": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"lookup_by": {
																Type:     schema.TypeString,
																Required: true,
															},
															"regex": {
																Type:     schema.TypeString,
																Required: true,
															},
															"source": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
												"route_to": {
													Type:     schema.TypeString,
													Optional: true,
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

func checkDynamicRoutingRule(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
	rNum := diff.Get("set.0.rule.#").(int)
	draIdxs := []int{}
	errorMsgs := []string{}

	for ri := 0; ri < rNum; ri++ {
		dra := diff.Get(fmt.Sprintf("set.0.rule.%d.actions.0.dynamic_route_to", ri))
		hasDra := isNonEmptyList(dra)
		if !hasDra {
			continue
		}
		draIdxs = append(draIdxs, ri)
	}
	// 1. Only the first rule of the first ("start") set can have the Dynamic Routing action:
	if len(draIdxs) > 1 {
		idxs := []string{}
		for _, idx := range draIdxs {
			idxs = append(idxs, fmt.Sprintf("%d", idx))
		}
		errorMsgs = append(errorMsgs, fmt.Sprintf("A Router can have at most one Dynamic Routing rule; Rules with the dynamic_route_to action found at indexes: %s", strings.Join(idxs, ", ")))
	}
	// 2. The Dynamic Routing action can only be used in the first rule of the first set:
	if len(draIdxs) > 0 && draIdxs[0] != 0 {
		errorMsgs = append(errorMsgs, fmt.Sprintf("The Dynamic Routing rule must be the first rule in a Router"))
	}
	// 3. If the Dynamic Routing rule is the first rule of the first set,
	// validate its configuration. It cannot have any conditions or the `route_to` action:
	if len(draIdxs) == 1 && draIdxs[0] == 0 {
		condNum := diff.Get("set.0.rule.0.condition.#").(int)
		// diff.NewValueKnown(str) will return false if the value is based on interpolation that was unavailable at diff time,
		// which may be the case for the `route_to` action when it references a pagerduty_service resource.
		// Source: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/helper/schema#ResourceDiff.NewValueKnown
		routeToValueKnown := diff.NewValueKnown("set.0.rule.0.actions.0.route_to")
		routeTo := diff.Get("set.0.rule.0.actions.0.route_to").(string)
		if condNum > 0 {
			errorMsgs = append(errorMsgs, fmt.Sprintf("Dynamic Routing rules cannot have conditions"))
		}
		if !routeToValueKnown || routeToValueKnown && routeTo != "" {
			errorMsgs = append(errorMsgs, fmt.Sprintf("Dynamic Routing rules cannot have the `route_to` action"))
		}
	}

	if len(errorMsgs) > 0 {
		return fmt.Errorf("Invalid Dynamic Routing rule configuration:\n- %s", strings.Join(errorMsgs, "\n- "))
	}
	return nil
}

func resourcePagerDutyEventOrchestrationPathRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type %s for orchestration: %s", "router", d.Id())

		if routerPath, _, err := client.EventOrchestrationPaths.GetContext(ctx, d.Id(), "router"); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
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
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	// In order to delete an Orchestration Router an empty orchestration path
	// config should be sent as an update.
	emptyPath := emptyOrchestrationPathStructBuilder("router")
	routerID := d.Get("event_orchestration").(string)

	log.Printf("[INFO] Deleting PagerDuty Event Orchestration Router Path: %s", routerID)

	retryErr := retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
		if _, _, err := client.EventOrchestrationPaths.UpdateContext(ctx, routerID, "router", emptyPath); err != nil {
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

func resourcePagerDutyEventOrchestrationPathRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	routerPath := buildRouterPathStructForUpdate(d)
	var warnings []*pagerduty.EventOrchestrationPathWarning

	log.Printf("[INFO] Updating PagerDuty Event Orchestration Path of type %s for orchestration: %s", "router", routerPath.Parent.ID)

	retryErr := retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
		response, _, err := client.EventOrchestrationPaths.UpdateContext(ctx, routerPath.Parent.ID, "router", routerPath)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}
		if response == nil {
			return retry.NonRetryableError(fmt.Errorf("No Event Orchestration Router found."))
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
	actions := new(pagerduty.EventOrchestrationPathRuleActions)
	for _, ai := range v.([]interface{}) {
		am := ai.(map[string]interface{})
		dra := am["dynamic_route_to"]
		if isNonEmptyList(dra) {
			actions.DynamicRouteTo = expandRouterDynamicRouteToAction(dra)
		} else {
			actions.RouteTo = am["route_to"].(string)
		}
	}

	return actions
}

func expandCatchAll(v interface{}) *pagerduty.EventOrchestrationPathCatchAll {
	catchAll := new(pagerduty.EventOrchestrationPathCatchAll)

	for _, ca := range v.([]interface{}) {
		am := ca.(map[string]interface{})
		catchAll.Actions = expandRouterActions(am["actions"])
	}

	return catchAll
}

func expandRouterDynamicRouteToAction(v interface{}) *pagerduty.EventOrchestrationPathDynamicRouteTo {
	dr := new(pagerduty.EventOrchestrationPathDynamicRouteTo)
	for _, i := range v.([]interface{}) {
		dra := i.(map[string]interface{})
		dr.LookupBy = dra["lookup_by"].(string)
		dr.Regex = dra["regex"].(string)
		dr.Source = dra["source"].(string)
	}
	return dr
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
	if actions.DynamicRouteTo != nil {
		am["dynamic_route_to"] = flattenRouterDynamicRouteToAction(actions.DynamicRouteTo)
	} else {
		am["route_to"] = actions.RouteTo
	}
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

func flattenRouterDynamicRouteToAction(dra *pagerduty.EventOrchestrationPathDynamicRouteTo) []map[string]interface{} {
	var dr []map[string]interface{}

	dr = append(dr, map[string]interface{}{
		"lookup_by": dra.LookupBy,
		"regex":     dra.Regex,
		"source":    dra.Source,
	})

	return dr
}

func resourcePagerDutyEventOrchestrationPathRouterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	// given an orchestration ID import the router orchestration path
	orchestrationID := d.Id()
	pathType := "router"
	_, _, err = client.EventOrchestrationPaths.GetContext(ctx, orchestrationID, pathType)

	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(orchestrationID)
	d.Set("event_orchestration", orchestrationID)

	return []*schema.ResourceData{d}, nil
}
