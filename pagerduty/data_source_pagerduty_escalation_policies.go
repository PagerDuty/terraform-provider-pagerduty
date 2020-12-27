package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEscalationPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyEscalationPoliciesRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePagerDutyEscalationPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading all PagerDuty escalation policies")

	o := &pagerduty.ListEscalationPoliciesOptions{}

	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EscalationPolicies.List(o)
		if err != nil {
			time.Sleep(15 * time.Second)
			return resource.RetryableError(err)
		}

		var ids []string
		var names []string

		for _, ep := range resp.EscalationPolicies {
			ids = append(ids, ep.ID)
			names = append(names, ep.Name)
		}

		d.SetId(fmt.Sprintf("%d", len(resp.EscalationPolicies)))
		d.Set("ids", ids)
		d.Set("names", names)

		return nil
	})
}
