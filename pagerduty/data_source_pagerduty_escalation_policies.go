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

func dataSourcePagerDutyEscalationPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyEscalationPoliciesRead,

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

func dataSourcePagerDutyEscalationPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _ := meta.(*Config).Client()

	log.Printf("[INFO] Reading all PagerDuty escalation policies")

	o := &pagerduty.ListEscalationPoliciesOptions{}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EscalationPolicies.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
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
	}))
}
