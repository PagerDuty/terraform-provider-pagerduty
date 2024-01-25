package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		escalationPolicy, _, err := client.EscalationPolicies.Create(escalationPolicy)
		if err != nil {
			if isErrCode(err, 429) {
				// Delaying retry by 30s as recommended by PagerDuty
				// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
				time.Sleep(30 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		d.SetId(escalationPolicy.ID)
		readErr = fetchEscalationPolicy(d, meta, genError)
		if readErr != nil {
			return resource.NonRetryableError(readErr)
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

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		escalationPolicy, _, err := client.EscalationPolicies.Get(d.Id(), o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			log.Printf("[WARN] Escalation Policy read error")
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}

			return nil
		}

		d.Set("name", escalationPolicy.Name)
		d.Set("description", escalationPolicy.Description)
		d.Set("num_loops", escalationPolicy.NumLoops)

		if err := d.Set("teams", flattenTeams(escalationPolicy.Teams)); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error setting teams: %s", err))
		}

		if err := d.Set("rule", flattenEscalationRules(escalationPolicy.EscalationRules)); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourcePagerDutyEscalationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	escalationPolicy := buildEscalationPolicyStruct(d)

	log.Printf("[INFO] Updating PagerDuty escalation policy: %s", d.Id())

	retryErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, _, err := client.EscalationPolicies.Update(d.Id(), escalationPolicy); err != nil {
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

func resourcePagerDutyEscalationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty escalation policy: %s", d.Id())

	// Retrying to give other resources (such as services) to delete
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.EscalationPolicies.Delete(d.Id()); err != nil {
			if isErrCode(err, 400) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
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

func flattenEscalationRules(v []*pagerduty.EscalationRule) []map[string]interface{} {
	var escalationRules []map[string]interface{}

	for _, er := range v {
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

		for _, ert := range er.Targets {
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
