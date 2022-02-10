package pagerduty

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyRuleset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyRulesetCreate,
		ReadContext:   resourcePagerDutyRulesetRead,
		UpdateContext: resourcePagerDutyRulesetUpdate,
		DeleteContext: resourcePagerDutyRulesetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"routing_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildRulesetStruct(d *schema.ResourceData) *pagerduty.Ruleset {
	ruleset := &pagerduty.Ruleset{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("team"); ok {
		ruleset.Team = expandTeam(attr)
	}

	if attr, ok := d.GetOk("routing_keys"); ok {
		ruleset.RoutingKeys = expandKeys(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("type"); ok {
		ruleset.Type = attr.(string)
	}

	return ruleset
}

func expandKeys(v []interface{}) []string {
	keys := make([]string, len(v))

	for i, k := range v {
		keys[i] = fmt.Sprintf("%v", k)
	}

	return keys
}

func expandTeam(v interface{}) *pagerduty.RulesetObject {
	var team *pagerduty.RulesetObject
	t := v.([]interface{})[0].(map[string]interface{})
	team = &pagerduty.RulesetObject{
		ID: t["id"].(string),
	}

	return team
}

func flattenTeam(v *pagerduty.RulesetObject) []interface{} {
	team := map[string]interface{}{
		"id": v.ID,
	}

	return []interface{}{team}
}

func fetchPagerDutyRuleset(ctx context.Context, d *schema.ResourceData, meta interface{}, handle404Errors bool) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		ruleset, _, err := client.Rulesets.Get(d.Id())
		if checkErr := getErrorHandler(handle404Errors)(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		d.Set("name", ruleset.Name)
		d.Set("type", ruleset.Type)

		// if ruleset is found set to ResourceData
		if ruleset.Team != nil {
			d.Set("team", flattenTeam(ruleset.Team))
		}
		d.Set("routing_keys", ruleset.RoutingKeys)

		return nil
	}))
}

func resourcePagerDutyRulesetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	ruleset := buildRulesetStruct(d)

	log.Printf("[INFO] Creating PagerDuty ruleset: %s", ruleset.Name)

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if ruleset, _, err := client.Rulesets.Create(ruleset); err != nil {
			if isErrCode(err, 400) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if ruleset != nil {
			d.SetId(ruleset.ID)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}
	return fetchPagerDutyRuleset(ctx, d, meta, false)
}

func resourcePagerDutyRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty ruleset: %s", d.Id())
	return fetchPagerDutyRuleset(ctx, d, meta, true)

}
func resourcePagerDutyRulesetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	ruleset := buildRulesetStruct(d)

	log.Printf("[INFO] Updating PagerDuty ruleset: %s", d.Id())

	if _, _, err := client.Rulesets.Update(d.Id(), ruleset); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyRulesetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty ruleset: %s", d.Id())

	if _, err := client.Rulesets.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
