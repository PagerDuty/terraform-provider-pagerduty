package pagerduty

import (
	"encoding/json"
	"log"

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
			"advanced_conditions_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameters": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
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
							Optional: true,
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

	if attr, ok := d.GetOk("advanced_conditions_json"); ok {
		rule.AdvancedConditions = expandString(attr.(string))
	}
	if attr, ok := d.GetOk("action"); ok {
		rule.Actions = expandActions(attr.([]interface{}))
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
		con := &pagerduty.RuleConditions{
			Operator:          vm["operator"].(string),
			RuleSubconditions: expandSubConditions(vm["subconditions"].([]interface{})),
		}
		conditions = con
	}

	return conditions
}

func expandActions(v interface{}) []*pagerduty.RuleAction {
	var actions []*pagerduty.RuleAction

	for _, ai := range v.([]interface{}) {
		am := ai.(map[string]interface{})
		act := &pagerduty.RuleAction{
			Action:     am["action"].(string),
			Parameters: expandActionParameters(am["parameters"].([]interface{})),
		}
		actions = append(actions, act)
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

func expandActionParameters(v interface{}) map[string]string {
	parms := make(map[string]string)

	for _, pi := range v.([]interface{}) {
		pm := pi.(map[string]interface{})

		for key, value := range pm {
			if value != nil {
				parms[key] = value.(string)
			}
		}
	}
	return parms
}

func flattenConditions(conditions *pagerduty.RuleConditions) interface{} {
	b, err := json.Marshal(conditions)
	if err != nil {
		log.Printf("[ERROR] Could not conditions field: %v", err)
		return nil
	}
	return string(b)
}

func flattenActions(actions []*pagerduty.RuleAction) interface{} {
	b, err := json.Marshal(actions)
	if err != nil {
		log.Printf("[ERROR] Could not flatten actions field: %v", err)
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

	rule, _, err := client.Rulesets.GetRule(rulesetID, d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("conditions", flattenConditions(rule.Conditions))

	if rule.Actions != nil {
		d.Set("action", flattenActions(rule.Actions))
	}
	d.Set("position", rule.Position)
	d.Set("disabled", rule.Disabled)

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
