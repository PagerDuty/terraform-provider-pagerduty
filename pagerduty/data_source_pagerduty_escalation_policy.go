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

func dataSourcePagerDutyEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyEscalationPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyEscalationPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty escalation policy")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListEscalationPoliciesOptions{
		Query: searchName,
	}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EscalationPolicies.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		var found *pagerduty.EscalationPolicy

		for _, policy := range resp.EscalationPolicies {
			if policy.Name == searchName {
				found = policy
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any escalation policy with the name: %s", searchName),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)

		return nil
	}))
}
