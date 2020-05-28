package pagerduty

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyRuleset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyRulesetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyRulesetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty ruleset")

	searchName := d.Get("name").(string)

	resp, _, err := client.Rulesets.List()
	if err != nil {
		return err
	}

	var found *pagerduty.Ruleset

	for _, ruleset := range resp.Rulesets {
		if ruleset.Name == searchName {
			found = ruleset
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any ruleset with the name: %s", searchName)
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)

	return nil
}
