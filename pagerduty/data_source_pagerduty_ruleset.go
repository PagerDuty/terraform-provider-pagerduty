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

func dataSourcePagerDutyRuleset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyRulesetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routing_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePagerDutyRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty ruleset")

	searchName := d.Get("name").(string)

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Rulesets.List()
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		var found *pagerduty.Ruleset

		for _, ruleset := range resp.Rulesets {
			if ruleset.Name == searchName {
				found = ruleset
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any ruleset with the name: %s", searchName),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("routing_keys", found.RoutingKeys)

		return nil
	}))
}
