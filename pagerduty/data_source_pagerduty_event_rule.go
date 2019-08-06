package pagerduty

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEventRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyEventRuleRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty event rule")

	searchId := d.Get("id").(string)

	resp, _, err := client.EventRules.List()
	if err != nil {
		return err
	}

	var found *pagerduty.EventRule

	for _, rule := range resp.EventRules {
		if rule.ID == searchId {
			found = rule
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any escalation policy with the name: %s", searchId)
	}

	d.SetId(found.ID)

	return nil
}
