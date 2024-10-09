package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyEscalationPolicyCreate,
		Read:   resourcePagerDutyEscalationPolicyRead,
		Update: resourcePagerDutyEscalationPolicyUpdate,
		Delete: resourcePagerDutyEscalationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIsAllowedString(NoNonPrintableCharsOrSpecialChars),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"num_loops": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 9),
			},
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 1,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"escalation_delay_in_minutes": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"escalation_rule_assignment_strategy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"assign_to_everyone",
											"round_robin",
										}),
									},
								},
							},
						},
						"target": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "user_reference",
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"user_reference",
											"schedule_reference",
										}),
									},
									"id": {
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

func buildEscalationPolicyStruct(d *schema.ResourceData) *pagerduty.EscalationPolicy {
	escalationPolicy := &pagerduty.EscalationPolicy{
		Name:            d.Get("name").(string),
		EscalationRules: expandEscalationRules(d.Get("rule").([]interface{})),
	}

	if attr, ok := d.GetOk("description"); ok {
		escalationPolicy.Description = attr.(string)
	}

	loops := d.Get("num_loops").(int)
	escalationPolicy.NumLoops = &loops

	if attr, ok := d.GetOk("teams"); ok {
		escalationPolicy.Teams = expandTeams(attr.([]interface{}))
	}

	return escalationPolicy
}

func resourcePagerDutyEscalationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
	var readErr error

	escalationPolicy := buildEscalationPolicyStruct(d)

	log.Printf("[INFO] Creating PagerDuty escalation policy: %s", escalationPolicy.Name)

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		escalationPolicy, _, err := client.EscalationPolicies.Create(escalationPolicy)
		if err != nil {
			if isErrCode(err, 429) {
				// Delaying retry by 30s as recommended by PagerDuty
				// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
				time.Sleep(30 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		d.SetId(escalationPolicy.ID)
		readErr = fetchEscalationPolicy(d, meta, genError)
		if readErr != nil {
			return retry.NonRetryableError(readErr)
		}
		return nil
	})
}

func resourcePagerDutyEscalationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading PagerDuty escalation policy: %s", d.Id())
	return fetchEscalationPolicy(d, meta, handleNotFoundError)
}

func fetchEscalationPolicy(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	o := &pagerduty.GetEscalationPolicyOptions{Includes: []string{"escalation_rule_assignment_strategies"}}

	var escalationPolicyFirstAttempt *pagerduty.EscalationPolicy

	escalationPolicyFirstAttempt, _, err = client.EscalationPolicies.Get(d.Id(), o)
	if err != nil && isErrCode(err, http.StatusForbidden) || isMalformedForbiddenError(err) {
		// Removing the inclusion of escalation_rule_assignment_strategies for
		// accounts wihtout the required entitlements.
		o = nil
	}

	if err == nil && escalationPolicyFirstAttempt != nil {
		return setResourceEPProps(d, escalationPolicyFirstAttempt)
	}

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		escalationPolicy, _, err := client.EscalationPolicies.Get(d.Id(), o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) || isErrCode(err, http.StatusForbidden) || isMalformedForbiddenError(err) {
				return retry.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			log.Printf("[WARN] Escalation Policy read error")
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}

			return nil
		}

		err = setResourceEPProps(d, escalationPolicy)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})
}

func setResourceEPProps(d *schema.ResourceData, escalationPolicy *pagerduty.EscalationPolicy) error {
	d.Set("name", escalationPolicy.Name)
	d.Set("description", escalationPolicy.Description)
	d.Set("num_loops", escalationPolicy.NumLoops)

	if err := d.Set("teams", flattenTeams(escalationPolicy.Teams)); err != nil {
		return fmt.Errorf("error setting teams: %s", err)
	}

	if err := d.Set("rule", flattenEscalationRules(escalationPolicy.EscalationRules, d)); err != nil {
		return err
	}
	return nil
}

func resourcePagerDutyEscalationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	escalationPolicy := buildEscalationPolicyStruct(d)

	log.Printf("[INFO] Updating PagerDuty escalation policy: %s", d.Id())

	_, _, err = client.EscalationPolicies.Update(d.Id(), escalationPolicy)
	if err == nil {
		return nil
	}

	if isErrCode(err, http.StatusForbidden) || isMalformedForbiddenError(err) {
		// Removing the inclusion of escalation_rule_assignment_strategies for
		// accounts wihtout the required entitlements.
		for idx, er := range escalationPolicy.EscalationRules {
			if er.EscalationRuleAssignmentStrategy == nil {
				continue
			}

			if er.EscalationRuleAssignmentStrategy.Type == "round_robin" {
				return fmt.Errorf("Round Robin Scheduling is available for accounts on the following pricing plans: Business, Digital Operations (legacy) and Enterprise for Incident Management. Therefore, set the escalation_rule_assignment_strategy to 'assign_to_everyone' for the escalation rule at index %d", idx)
			}

			er.EscalationRuleAssignmentStrategy = nil
			escalationPolicy.EscalationRules[idx] = er
		}
	}

	retryErr := retry.Retry(5*time.Minute, func() *retry.RetryError {
		if _, _, err := client.EscalationPolicies.Update(d.Id(), escalationPolicy); err != nil {
			if isErrCode(err, http.StatusBadRequest) || isErrCode(err, http.StatusForbidden) || isMalformedForbiddenError(err) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return nil
}

func resourcePagerDutyEscalationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty escalation policy: %s", d.Id())

	// Retrying to give other resources (such as services) to delete
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.EscalationPolicies.Delete(d.Id()); err != nil {
			if isErrCode(err, 400) {
				return retry.RetryableError(err)
			}

			err = handleNotFoundError(err, d)
			if err != nil {
				return retry.NonRetryableError(err)
			}

			return nil
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	d.SetId("")

	// giving the API time to catchup
	time.Sleep(time.Second)
	return nil
}

func expandEscalationRules(v interface{}) []*pagerduty.EscalationRule {
	var escalationRules []*pagerduty.EscalationRule

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		escalationRule := &pagerduty.EscalationRule{
			EscalationDelayInMinutes:         rer["escalation_delay_in_minutes"].(int),
			EscalationRuleAssignmentStrategy: expandEscalationRuleAssignmentStrategy(rer["escalation_rule_assignment_strategy"]),
		}

		for _, ert := range rer["target"].([]interface{}) {
			rert := ert.(map[string]interface{})
			escalationRuleTarget := &pagerduty.EscalationTargetReference{
				ID:   rert["id"].(string),
				Type: rert["type"].(string),
			}

			escalationRule.Targets = append(escalationRule.Targets, escalationRuleTarget)
		}

		escalationRules = append(escalationRules, escalationRule)
	}

	return escalationRules
}

func expandEscalationRuleAssignmentStrategy(v interface{}) *pagerduty.EscalationRuleAssignmentStrategy {
	log.Printf("expandEscalationRuleAssignmentStrategy_v is %+v", v)
	escalationRuleAssignmentStrategy := &pagerduty.EscalationRuleAssignmentStrategy{}
	pre := v.([]interface{})
	if len(pre) == 0 || isNilFunc(pre[0]) {
		return nil
	}

	eras := pre[0].(map[string]interface{})
	teras := eras["type"].(string)
	log.Printf("expandEscalationRuleAssignmentStrategy_teras is %#v", teras)
	escalationRuleAssignmentStrategy.Type = teras
	return escalationRuleAssignmentStrategy
}

func flattenEscalationRules(v []*pagerduty.EscalationRule, d *schema.ResourceData) []map[string]interface{} {
	var escalationRules []map[string]interface{}

	for i, er := range v {
		escalationRule := map[string]interface{}{
			"id":                          er.ID,
			"escalation_delay_in_minutes": er.EscalationDelayInMinutes,
		}

		var escalationRuleAssignmentStrategy []map[string]interface{}
		if er.EscalationRuleAssignmentStrategy != nil {
			eras := map[string]interface{}{"type": er.EscalationRuleAssignmentStrategy.Type}
			escalationRuleAssignmentStrategy = append(escalationRuleAssignmentStrategy, eras)
		}
		escalationRule["escalation_rule_assignment_strategy"] = escalationRuleAssignmentStrategy

		var targets []map[string]interface{}
		addedTargets := map[string]struct{}{}

		// Append targets in same orden as plan, then mark them as added
		if d != nil {
			targetsPlan := d.Get(fmt.Sprintf("rule.%d.target", i)).([]any)
			for _, tpValue := range targetsPlan {
				var ert *pagerduty.EscalationTargetReference
				for _, t := range er.Targets {
					if t.ID == tpValue.(map[string]any)["id"].(string) {
						ert = t
						break
					}
				}
				if ert == nil {
					continue
				}
				escalationRuleTarget := map[string]interface{}{"id": ert.ID, "type": ert.Type}
				targets = append(targets, escalationRuleTarget)
				addedTargets[ert.ID] = struct{}{}
			}
		}

		// Append targets not present in plan
		for _, ert := range er.Targets {
			if _, found := addedTargets[ert.ID]; found {
				continue
			}
			escalationRuleTarget := map[string]interface{}{"id": ert.ID, "type": ert.Type}
			targets = append(targets, escalationRuleTarget)
		}

		escalationRule["target"] = targets
		escalationRules = append(escalationRules, escalationRule)
	}

	return escalationRules
}

func expandTeams(v interface{}) []*pagerduty.TeamReference {
	var teams []*pagerduty.TeamReference

	for _, t := range v.([]interface{}) {
		team := &pagerduty.TeamReference{
			ID:   t.(string),
			Type: "team_reference",
		}
		teams = append(teams, team)
	}

	return teams
}

func flattenTeams(teams []*pagerduty.TeamReference) []string {
	res := make([]string, len(teams))
	for i, t := range teams {
		res[i] = t.ID
	}

	return res
}
