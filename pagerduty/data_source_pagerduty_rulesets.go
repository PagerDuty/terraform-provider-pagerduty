package pagerduty

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyRulesets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyRulesetsRead,

		Schema: map[string]*schema.Schema{
			"search": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rulesets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routing_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourcePagerDutyRulesetsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty rulesets")

	searchName := d.Get("search").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Rulesets.List()
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var rulesets []*pagerduty.Ruleset
		re := regexp.MustCompile(searchName)
		for _, ruleset := range resp.Rulesets {
			if re.MatchString(ruleset.Name) {
				rulesets = append(rulesets, ruleset)
			}
		}

		if len(rulesets) == 0 {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any ruleset with the name: %s", searchName),
			)
		}

		d.SetId(resource.UniqueId())
		d.Set("search", searchName)
		d.Set("rulesets", flattenPagerDutyRulesets(rulesets))

		return nil
	})
}

func flattenPagerDutyRulesets(rulesets []*pagerduty.Ruleset) []interface{} {
	var result []interface{}

	for _, i := range rulesets {
		ruleset := map[string]interface{}{
			"id":           i.ID,
			"name":         i.Name,
			"routing_keys": i.RoutingKeys,
		}
		result = append(result, ruleset)
	}
	return result
}
