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

func dataSourcePagerDutyIncidentCustomField() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyIncidentCustomFieldRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"field_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyIncidentCustomFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty data source")

	searchName := d.Get("name").(string)

	err = resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.IncidentCustomFields.ListContext(ctx, nil)
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.IncidentCustomField

		for _, field := range resp.Fields {
			if field.Name == searchName {
				found = field
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("unable to locate any field with name: %s", searchName),
			)
		}

		err = flattenIncidentCustomField(d, found)
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
