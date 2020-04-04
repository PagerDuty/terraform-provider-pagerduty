package pagerduty

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyRulesetRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyRulesetRuleCreate,
		Read:   resourcePagerDutyRulesetRuleRead,
		Update: resourcePagerDutyRulesetRuleUpdate,
		Delete: resourcePagerDutyRulesetRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ruleset": {
				Type:     schema.TypeString,
				Required: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subconditions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"time_frame": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheduled_weekly": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"start_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"duration": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"weekdays": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"active_between": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"end_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"suppress": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold_value": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"threshold_time_unit": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateValueFunc([]string{
											"minutes",
											"seconds",
											"hours",
										}),
									},
									"threshold_time_amount": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"severity": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateValueFunc([]string{
											"info",
											"error",
											"warning",
											"critical",
										}),
									},
								},
							},
						},
						"route": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"priority": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"annotate": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"event_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateValueFunc([]string{
											"trigger",
											"resolve",
										}),
									},
								},
							},
						},
						"extractions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"regex": {
										Type:     schema.TypeString,
										Optional: true,
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

func buildRulesetRuleStruct(d *schema.ResourceData) *pagerduty.RulesetRule {
	rule := &pagerduty.RulesetRule{
		Ruleset: &pagerduty.RulesetReference{
			Type: "ruleset",
			ID:   d.Get("ruleset").(string),
		},
		Conditions: expandConditions(d.Get("conditions").([]interface{})),
	}

	if attr, ok := d.GetOk("actions"); ok {
		rule.Actions = expandActions(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("time_frame"); ok {
		rule.TimeFrame = expandTimeFrame(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("position"); ok {
		rule.Position = attr.(int)
	}
	if attr, ok := d.GetOk("disabled"); ok {
		rule.Disabled = attr.(bool)
	}

	return rule
}

func expandConditions(v interface{}) *pagerduty.RuleConditions {
	var conditions *pagerduty.RuleConditions

	for _, vi := range v.([]interface{}) {
		vm := vi.(map[string]interface{})
		conditions = &pagerduty.RuleConditions{
			Operator:          vm["operator"].(string),
			RuleSubconditions: expandSubConditions(vm["subconditions"].([]interface{})),
		}
	}

	return conditions
}

func expandTimeFrame(v interface{}) *pagerduty.RuleTimeFrame {
	var tFrame = new(pagerduty.RuleTimeFrame)

	for _, tfi := range v.([]interface{}) {
		tfm := tfi.(map[string]interface{})

		if tfm["scheduled_weekly"] != nil {
			tFrame.ScheduledWeekly = expandScheduledWeekly(tfm["scheduled_weekly"].(interface{}))
		}
		if tfm["active_between"] != nil {
			tFrame.ActiveBetween = expandActiveBetween(tfm["active_between"].(interface{}))
		}
	}

	return tFrame
}

func expandScheduledWeekly(v interface{}) *pagerduty.ScheduledWeekly {
	var sw *pagerduty.ScheduledWeekly

	for _, swi := range v.([]interface{}) {
		swm := swi.(map[string]interface{})

		sw = &pagerduty.ScheduledWeekly{
			Timezone:  swm["timezone"].(string),
			StartTime: swm["start_time"].(int),
			Duration:  swm["duration"].(int),
			Weekdays:  convertToIntArray(swm["weekdays"].([]interface{})),
		}
	}
	return sw
}

func convertToIntArray(v []interface{}) []int {
	ints := make([]int, len(v))

	for i := range v {
		ints[i] = v[i].(int)
	}
	return ints
}

func expandActiveBetween(v interface{}) *pagerduty.ActiveBetween {
	var ab *pagerduty.ActiveBetween

	for _, abi := range v.([]interface{}) {
		abm := abi.(map[string]interface{})

		ab = &pagerduty.ActiveBetween{
			StartTime: abm["start_time"].(int),
			EndTime:   abm["end_time"].(int),
		}
	}

	return ab
}

func expandActions(v interface{}) *pagerduty.RuleActions {
	var actions = new(pagerduty.RuleActions)

	for _, ai := range v.([]interface{}) {
		am := ai.(map[string]interface{})
		log.Printf("[DEBUG] Severity: %v", am)

		if am["suppress"] != nil {
			actions.Suppress = expandSuppress(am["suppress"].(interface{}))
		}
		if am["extractions"] != nil {
			actions.Extractions = expandExtractions(am["extractions"].(interface{}))
		}
		if am["severity"] != nil {
			actions.Severity = expandActionParameters(am["severity"].(interface{}))
		}
		if am["route"] != nil {
			actions.Route = expandActionParameters(am["route"].(interface{}))
		}
		if am["priority"] != nil {
			actions.Priority = expandActionParameters(am["priority"].(interface{}))
		}
		if am["event_action"] != nil {
			actions.EventAction = expandActionParameters(am["event_action"].(interface{}))
		}
		if am["annotate"] != nil {
			actions.Annotate = expandActionParameters(am["annotate"].(interface{}))
		}
		if am["event_actions"] != nil {
			actions.Annotate = expandActionParameters(am["event_actions"].(interface{}))
		}
	}
	return actions
}

func expandSubConditions(v interface{}) []*pagerduty.RuleSubcondition {
	var sc []*pagerduty.RuleSubcondition

	for _, sci := range v.([]interface{}) {
		scm := sci.(map[string]interface{})
		scon := &pagerduty.RuleSubcondition{
			Operator:   scm["operator"].(string),
			Parameters: expandSubConditionParameters(scm["parameter"].([]interface{})),
		}
		sc = append(sc, scon)
	}
	return sc
}
func expandSubConditionParameters(v interface{}) *pagerduty.ConditionParameter {
	var parms *pagerduty.ConditionParameter

	for _, pi := range v.([]interface{}) {
		pm := pi.(map[string]interface{})
		cp := &pagerduty.ConditionParameter{
			Path:  pm["path"].(string),
			Value: pm["value"].(string),
		}
		parms = cp
	}
	return parms
}

func expandActionParameters(v interface{}) *pagerduty.RuleActionParameter {
	var rap *pagerduty.RuleActionParameter

	for _, pi := range v.([]interface{}) {
		pm := pi.(map[string]interface{})
		if pm["value"] != nil {
			rap = &pagerduty.RuleActionParameter{
				Value: pm["value"].(string),
			}
		}
	}
	return rap
}

func expandSuppress(v interface{}) *pagerduty.RuleActionSuppress {
	var ras *pagerduty.RuleActionSuppress

	for _, sai := range v.([]interface{}) {
		sa := sai.(map[string]interface{})
		ras = &pagerduty.RuleActionSuppress{
			Value:               sa["value"].(bool),
			ThresholdValue:      sa["threshold_value"].(int),
			ThresholdTimeUnit:   sa["threshold_time_unit"].(string),
			ThresholdTimeAmount: sa["threshold_time_amount"].(int),
		}
	}
	return ras
}

func expandExtractions(v interface{}) []*pagerduty.RuleActionExtraction {
	var rae []*pagerduty.RuleActionExtraction

	for _, eai := range v.([]interface{}) {
		ea := eai.(map[string]interface{})
		ext := &pagerduty.RuleActionExtraction{
			Target: ea["target"].(string),
			Source: ea["source"].(string),
			Regex:  ea["regex"].(string),
		}
		rae = append(rae, ext)
	}
	return rae
}

func flattenConditions(conditions *pagerduty.RuleConditions) interface{} {
	b, err := json.Marshal(conditions)
	if err != nil {
		log.Printf("[ERROR] Could not conditions field: %v", err)
		return nil
	}
	return string(b)
}

func flattenActions(actions *pagerduty.RuleActions) interface{} {
	b, err := json.Marshal(actions)
	if err != nil {
		log.Printf("[ERROR] Could not flatten actions field: %v", err)
		return nil
	}
	return string(b)
}

func flattenTimeFrame(timeframe *pagerduty.RuleTimeFrame) interface{} {
	b, err := json.Marshal(timeframe)
	if err != nil {
		log.Printf("[ERROR] Could not flatten ruleset rule time frame field: %v", err)
		return nil
	}
	return string(b)
}

func resourcePagerDutyRulesetRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	rule := buildRulesetRuleStruct(d)

	log.Printf("[INFO] Creating PagerDuty ruleset rule for ruleset: %s", rule.Ruleset.ID)

	rule, _, err := client.Rulesets.CreateRule(rule.Ruleset.ID, rule)
	if err != nil {
		return err
	}

	d.SetId(rule.ID)

	return resourcePagerDutyRulesetRuleRead(d, meta)
}

func resourcePagerDutyRulesetRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty ruleset rule: %s", d.Id())
	rulesetID := d.Get("ruleset").(string)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if rule, _, err := client.Rulesets.GetRule(rulesetID, d.Id()); err != nil {
			if isErrCode(err, 500) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		} else if rule != nil {
			d.Set("conditions", flattenConditions(rule.Conditions))

			if rule.Actions != nil {
				d.Set("actions", flattenActions(rule.Actions))
			}
			if rule.TimeFrame != nil {
				d.Set("time_frame", flattenTimeFrame(rule.TimeFrame))
			}
			d.Set("position", rule.Position)
			d.Set("disabled", rule.Disabled)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return nil
}

func resourcePagerDutyRulesetRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	rule := buildRulesetRuleStruct(d)

	log.Printf("[INFO] Updating PagerDuty ruleset rule: %s", d.Id())
	rulesetID := d.Get("ruleset").(string)

	if _, _, err := client.Rulesets.UpdateRule(rulesetID, d.Id(), rule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyRulesetRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty ruleset rule: %s", d.Id())
	rulesetID := d.Get("ruleset").(string)

	if _, err := client.Rulesets.DeleteRule(rulesetID, d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
