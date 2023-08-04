package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyEscalationPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyEscalationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty escalation policy")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListEscalationPoliciesOptions{
		Query: searchName,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EscalationPolicies.List(o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
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
	})
}
