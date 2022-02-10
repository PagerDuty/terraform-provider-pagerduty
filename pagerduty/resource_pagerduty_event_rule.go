package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyEventRuleCreate,
		ReadContext:   resourcePagerDutyEventRuleRead,
		UpdateContext: resourcePagerDutyEventRuleUpdate,
		DeleteContext: resourcePagerDutyEventRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourcePagerDutyEventRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	eventRule := buildEventRuleStruct(d)

	log.Printf("[INFO] Creating PagerDuty event rule: %s", eventRule.Condition)

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if eventRule, _, err := client.EventRules.Create(eventRule); err != nil {
			return resource.RetryableError(err)
		} else if eventRule != nil {
			d.SetId(eventRule.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	return resourcePagerDutyEventRuleRead(ctx, d, meta)
}

func resourcePagerDutyEventRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty event rule: %s", d.Id())

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EventRules.List()
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
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
	}))
}
func resourcePagerDutyEventRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	eventRule := buildEventRuleStruct(d)

	log.Printf("[INFO] Updating PagerDuty event rule: %s", d.Id())

	if _, _, err := client.EventRules.Update(d.Id(), eventRule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyEventRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty event rule: %s", d.Id())

	if _, err := client.EventRules.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
