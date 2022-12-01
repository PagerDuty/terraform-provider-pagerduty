package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyIncidentWorkflow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyIncidentWorkflowRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyIncidentWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty incident workflow")

	searchName := d.Get("name").(string)

	err = resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.IncidentWorkflows.ListContext(ctx, &pagerduty.ListIncidentWorkflowOptions{})
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.IncidentWorkflow

		for _, iw := range resp.IncidentWorkflows {
			if iw.Name == searchName {
				found = iw
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("unable to locate any incident workflow with name: %s", searchName),
			)
		}

		err = flattenIncidentWorkflow(d, found, false, nil)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}
	return nil

}
