package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationPathRouter() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePagerDutyEventOrchestrationPathRouterRead,
		Create: resourcePagerDutyEventOrchestrationPathRouterCreate,
		Update: resourcePagerDutyEventOrchestrationPathRouterUpdate,
		Delete: resourcePagerDutyEventOrchestrationPathRouterDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyEventOrchestrationPathRouterImport,
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
				MaxItems: 1, // Router can only have 'start' set
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

func resourcePagerDutyEventOrchestrationPathRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type %s for orchestration: %s", "router", d.Id())

		if routerPath, _, err := client.EventOrchestrationPaths.Get(d.Id(), "router"); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if routerPath != nil {
			d.Set("type", routerPath.Type)

			if routerPath.Sets != nil {
				d.Set("sets", flattenSets(routerPath.Sets))
			}

			if routerPath.CatchAll != nil {
				d.Set("catch_all", flattenCatchAll(routerPath.CatchAll))
			}
		}
		return nil
	})

}

// EventOrchestrationPath cannot be created, use update to add / edit / remove rules and sets
func resourcePagerDutyEventOrchestrationPathRouterCreate(d *schema.ResourceData, meta interface{}) error {
	return resourcePagerDutyEventOrchestrationPathRouterUpdate(d, meta)
}

func resourcePagerDutyEventOrchestrationPathRouterDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourcePagerDutyEventOrchestrationPathRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	updatePath := buildRouterPathStructForUpdate(d)

	log.Printf("[INFO] Updating PagerDuty Event Orchestration Path of type %s for orchestration: %s", "router", updatePath.Parent.ID)

	return performRouterPathUpdate(d, updatePath, client)
}

func performRouterPathUpdate(d *schema.ResourceData, routerPath *pagerduty.EventOrchestrationPath, client *pagerduty.Client) error {
	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		updatedPath, _, err := client.EventOrchestrationPaths.Update(routerPath.Parent.ID, "router", routerPath)
		if err != nil {
			return resource.RetryableError(err)
		}
		if updatedPath == nil {
			return resource.NonRetryableError(fmt.Errorf("No Event Orchestration Router found."))
		}
		// set props
		d.SetId(routerPath.Parent.ID)
		d.Set("type", updatedPath.Type)

		if routerPath.Sets != nil {
			d.Set("sets", flattenSets(routerPath.Sets))
		}
		if updatedPath.CatchAll != nil {
			d.Set("catch_all", flattenCatchAll(updatedPath.CatchAll))
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return nil
}

func buildRouterPathParent(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	orchPath := &pagerduty.EventOrchestrationPath{}

	if attr, ok := d.GetOk("parent"); ok {
		orchPath.Parent = expandOrchestrationPathParent(attr)
	}

	return orchPath
}

func buildRouterPathStructForUpdate(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {

	orchPath := buildRouterPathParent(d)

	if attr, ok := d.GetOk("parent"); ok {
		orchPath.Parent = expandOrchestrationPathParent(attr)
	}

	if attr, ok := d.GetOk("sets"); ok {
		orchPath.Sets = expandSets(attr)
	}

	if attr, ok := d.GetOk("catch_all"); ok {
		orchPath.CatchAll = expandCatchAll(attr)
	}

	return orchPath
}

func expandOrchestrationPathParent(v interface{}) *pagerduty.EventOrchestrationPathReference {
	var parent *pagerduty.EventOrchestrationPathReference
	p := v.([]interface{})[0].(map[string]interface{})
	parent = &pagerduty.EventOrchestrationPathReference{
		ID: p["id"].(string),
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
	items := v.([]interface{})
	rules := []*pagerduty.EventOrchestrationPathRule{}

	for _, rule := range items {
		r := rule.(map[string]interface{})

		ruleInSet := &pagerduty.EventOrchestrationPathRule{
			ID:         r["id"].(string),
			Label:      r["label"].(string),
			Disabled:   r["disabled"].(bool),
			Conditions: expandRouterConditions(r["conditions"].(interface{})),
			Actions:    expandRouterActions(r["actions"].([]interface{})),
		}

		rules = append(rules, ruleInSet)
	}
	return rules
}

func expandRouterConditions(v interface{}) []*pagerduty.EventOrchestrationPathRuleCondition {
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
			"conditions": flattenRouterConditions(rule.Conditions),
			"actions":    flattenRouterActions(rule.Actions),
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	return flattenedRules
}

func flattenRouterConditions(conditions []*pagerduty.EventOrchestrationPathRuleCondition) []interface{} {
	var flattendConditions []interface{}

	for _, condition := range conditions {
		flattendCondition := map[string]interface{}{
			"expression": condition.Expression,
		}
		flattendConditions = append(flattendConditions, flattendCondition)
	}

	return flattendConditions
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

	return []*schema.ResourceData{d}, nil
}
