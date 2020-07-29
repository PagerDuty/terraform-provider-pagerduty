package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyResponsePlay() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyRulesetRuleCreate,
		Read:   resourcePagerDutyRulesetRuleRead,
		Update: resourcePagerDutyRulesetRuleUpdate,
		Delete: resourcePagerDutyRulesetRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyRulesetRuleImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscribers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					ID: schema.TypeString,
					Type: schema.TypeString,
				},
			},
			"subscribers_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"responders": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					ID: schema.TypeString,
					Type: schema.TypeString,
				},
			},
			"responders_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runnability": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"conference_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"conference_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildResponsePlayStruct(d *schema.ResourceData) *pagerduty.ResponsePlay {
	responsePlay := &pagerduty.ResponsePlay{
		Name: d.Get("name").(string),
	}
	if attr, ok := d.GetOk("type"); ok {
		responsePlay.Type = attr.(string)
	}
	if attr, ok := d.GetOk("description"); ok {
		responsePlay.Description = attr.(string)
	}
	if attr, ok := d.GetOk("team"); ok {
		responsePlay.Team = &pagerduty.TeamReference{
			ID:   attr.(string),
			Type: "team",
		}
	}
	if attr, ok := d.GetOk("subscribers"); ok {
		responsePlay.Subscribers = expandSubscribers(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("subscribers_message"); ok {
		responsePlay.SubscribersMessage = attr.(string)
	}

	// todo: responders

	if attr, ok := d.GetOk("responders_message"); ok {
		responsePlay.RespondersMessage = attr.(string)
	}

	if attr, ok := d.GetOk("runnability"); ok {
		responsePlay.Runnability = attr.(string)
	}

	if attr, ok := d.GetOk("conference_number"); ok {
		responsePlay.ConferenceNumber = attr.(string)
	}

	if attr, ok := d.GetOk("conference_url"); ok {
		responsePlay.ConferenceURL = attr.(string)
	}

	return responsePlay
}

func expandSubscribers(v interface{}) []*pagerduty.SubscriberReference {
	var subscribers []*pagerduty.RuleConditions

	for _, si := range v.([]interface{}) {
		sm := &pagerduty.SubscriberReference{
			ID: si.(string),
			Type: ""
			
			(map[string]interface{})
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

func flattenConditions(conditions *pagerduty.RuleConditions) []map[string]interface{} {
	var cons []map[string]interface{}

	con := map[string]interface{}{
		"operator":      conditions.Operator,
		"subconditions": flattenSubconditions(conditions.RuleSubconditions),
	}
	cons = append(cons, con)

	return cons
}

func flattenSubconditions(subconditions []*pagerduty.RuleSubcondition) []interface{} {
	var flattenedSubConditions []interface{}

	for _, sc := range subconditions {
		flattenedSubCon := map[string]interface{}{
			"operator":  sc.Operator,
			"parameter": flattenSubconditionParameters(sc.Parameters),
		}
		flattenedSubConditions = append(flattenedSubConditions, flattenedSubCon)
	}
	return flattenedSubConditions
}

func flattenSubconditionParameters(p *pagerduty.ConditionParameter) []interface{} {

	flattenedParams := map[string]interface{}{
		"path":  p.Path,
		"value": p.Value,
	}

	return []interface{}{flattenedParams}
}

func flattenActions(actions *pagerduty.RuleActions) []map[string]interface{} {
	var actionsMap []map[string]interface{}

	am := make(map[string]interface{})

	if actions.Suppress != nil {
		am["suppress"] = flattenSuppress(actions.Suppress)
	}
	if actions.Severity != nil {
		am["severity"] = flattenActionParameter(actions.Severity)
	}
	if actions.Route != nil {
		am["route"] = flattenActionParameter(actions.Route)
	}
	if actions.Priority != nil {
		am["priority"] = flattenActionParameter(actions.Priority)
	}
	if actions.Annotate != nil {
		am["annotate"] = flattenActionParameter(actions.Annotate)
	}
	if actions.EventAction != nil {
		am["event_action"] = flattenActionParameter(actions.EventAction)
	}
	if actions.Extractions != nil {
		am["extractions"] = flattenExtractions(actions.Extractions)
	}
	actionsMap = append(actionsMap, am)

	return actionsMap
}

func flattenSuppress(s *pagerduty.RuleActionSuppress) []interface{} {

	sup := map[string]interface{}{
		"value":                 s.Value,
		"threshold_value":       s.ThresholdValue,
		"threshold_time_unit":   s.ThresholdTimeUnit,
		"threshold_time_amount": s.ThresholdTimeAmount,
	}
	return []interface{}{sup}
}
func flattenActionParameter(ap *pagerduty.RuleActionParameter) []interface{} {

	param := map[string]interface{}{
		"value": ap.Value,
	}
	return []interface{}{param}
}

func flattenExtractions(rae []*pagerduty.RuleActionExtraction) []interface{} {
	var flatExtractList []interface{}

	for _, ex := range rae {
		flatExtract := map[string]interface{}{
			"target": ex.Target,
			"source": ex.Source,
			"regex":  ex.Regex,
		}
		flatExtractList = append(flatExtractList, flatExtract)
	}
	return flatExtractList
}

func flattenTimeFrame(timeframe *pagerduty.RuleTimeFrame) []map[string]interface{} {
	var tfMap []map[string]interface{}

	tm := make(map[string]interface{})

	if timeframe.ScheduledWeekly != nil {
		tm["scheduled_weekly"] = flattenScheduledWeekly(timeframe.ScheduledWeekly)
	}
	if timeframe.ActiveBetween != nil {
		tm["active_between"] = flattenActiveBetween(timeframe.ActiveBetween)
	}
	tfMap = append(tfMap, tm)

	return tfMap
}

func flattenScheduledWeekly(s *pagerduty.ScheduledWeekly) []interface{} {

	fsw := map[string]interface{}{
		"timezone":   s.Timezone,
		"start_time": s.StartTime,
		"duration":   s.Duration,
		"weekdays":   s.Weekdays,
	}
	return []interface{}{fsw}
}

func flattenActiveBetween(ab *pagerduty.ActiveBetween) []interface{} {

	fab := map[string]interface{}{
		"start_time": ab.StartTime,
		"end_time":   ab.EndTime,
	}
	return []interface{}{fab}
}

func resourcePagerDutyRulesetRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	rule := buildRulesetRuleStruct(d)

	log.Printf("[INFO] Creating PagerDuty ruleset rule for ruleset: %s", rule.Ruleset.ID)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if rule, _, err := client.Rulesets.CreateRule(rule.Ruleset.ID, rule); err != nil {
			return resource.RetryableError(err)
		} else if rule != nil {
			d.SetId(rule.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return resourcePagerDutyRulesetRuleRead(d, meta)
}

func resourcePagerDutyRulesetRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty ruleset rule: %s", d.Id())
	rulesetID := d.Get("ruleset").(string)

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		if rule, _, err := client.Rulesets.GetRule(rulesetID, d.Id()); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if rule != nil {
			if rule.Conditions != nil {
				d.Set("conditions", flattenConditions(rule.Conditions))
			}
			if rule.Actions != nil {
				d.Set("actions", flattenActions(rule.Actions))
			}
			if rule.TimeFrame != nil {
				d.Set("time_frame", flattenTimeFrame(rule.TimeFrame))
			}
			d.Set("position", rule.Position)
			d.Set("disabled", rule.Disabled)
			d.Set("ruleset", rulesetID)
		}
		return nil
	})
}

func resourcePagerDutyRulesetRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	rule := buildRulesetRuleStruct(d)

	log.Printf("[INFO] Updating PagerDuty ruleset rule: %s", d.Id())
	rulesetID := d.Get("ruleset").(string)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if _, _, err := client.Rulesets.UpdateRule(rulesetID, d.Id(), rule); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return nil
}

func resourcePagerDutyRulesetRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty ruleset rule: %s", d.Id())
	rulesetID := d.Get("ruleset").(string)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if _, err := client.Rulesets.DeleteRule(rulesetID, d.Id()); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	d.SetId("")

	return nil
}

func resourcePagerDutyRulesetRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*pagerduty.Client)

	ids := strings.Split(d.Id(), ".")

	if len(ids) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_ruleset_rule. Expecting an importation ID formed as '<ruleset_id>.<ruleset_rule_id>'")
	}
	rulesetID, ruleID := ids[0], ids[1]

	_, _, err := client.Rulesets.GetRule(rulesetID, ruleID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(ruleID)
	d.Set("ruleset", rulesetID)

	return []*schema.ResourceData{d}, nil
}
