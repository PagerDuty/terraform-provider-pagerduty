package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyRuleset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyRulesetRead,

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

func dataSourcePagerDutyRulesetRead(d *schema.ResourceData, meta interface{}) error {
	client, _ := meta.(*Config).Client()

	log.Printf("[INFO] Reading PagerDuty ruleset")

	searchName := d.Get("name").(string)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Rulesets.List()
		if checkErr := handleGenericErrors(err, d); checkErr != nil {
			return checkErr
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
	})
}
