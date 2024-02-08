package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyEventRuleCreate,
		Read:   resourcePagerDutyEventRuleRead,
		Update: resourcePagerDutyEventRuleUpdate,
		Delete: resourcePagerDutyEventRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"action_json": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition_json": {
				Type:     schema.TypeString,
				Required: true,
			},
			"advanced_condition_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catch_all": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildEventRuleStruct(d *schema.ResourceData) *pagerduty.EventRule {
	eventRule := &pagerduty.EventRule{
		Actions:   expandString(d.Get("action_json").(string)),
		Condition: expandString(d.Get("condition_json").(string)),
	}

	if attr, ok := d.GetOk("advanced_condition_json"); ok {
		eventRule.AdvancedCondition = expandString(attr.(string))
	}

	if attr, ok := d.GetOk("catch_all"); ok {
		eventRule.CatchAll = attr.(bool)
	}

	return eventRule
}

func resourcePagerDutyEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	eventRule := buildEventRuleStruct(d)

	log.Printf("[INFO] Creating PagerDuty event rule: %s", eventRule.Condition)

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if eventRule, _, err := client.EventRules.Create(eventRule); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else if eventRule != nil {
			d.SetId(eventRule.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return resourcePagerDutyEventRuleRead(d, meta)
}

func resourcePagerDutyEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty event rule: %s", d.Id())

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.EventRules.List()
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
		}
		var foundRule *pagerduty.EventRule

		for _, rule := range resp.EventRules {
			log.Printf("[DEBUG] Resp rule.ID: %s", rule.ID)
			if rule.ID == d.Id() {
				foundRule = rule
				break
			}
		}
		// check if eventRule  not  found
		if foundRule == nil {
			d.SetId("")
			return nil
		}
		// if event rule is found set to ResourceData
		d.Set("action_json", flattenSlice(foundRule.Actions))
		d.Set("condition_json", flattenSlice(foundRule.Condition))
		if foundRule.AdvancedCondition != nil {
			d.Set("advanced_condition_json", flattenSlice(foundRule.AdvancedCondition))
		}
		d.Set("catch_all", foundRule.CatchAll)
		return nil
	})
}

func resourcePagerDutyEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	eventRule := buildEventRuleStruct(d)

	log.Printf("[INFO] Updating PagerDuty event rule: %s", d.Id())

	if _, _, err := client.EventRules.Update(d.Id(), eventRule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty event rule: %s", d.Id())

	if _, err := client.EventRules.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
