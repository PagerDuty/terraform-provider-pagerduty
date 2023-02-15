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

func dataSourcePagerDutyFieldSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyFieldSchemaRead,
		Schema: map[string]*schema.Schema{
			"title": {
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

func dataSourcePagerDutyFieldSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty data source")

	search := d.Get("title").(string)

	err = resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.CustomFieldSchemas.ListContext(ctx, nil)
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.CustomFieldSchema

		for _, fs := range resp.Schemas {
			if fs.Title == search {
				found = fs
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("unable to locate any field schema with title: %s", search),
			)
		}

		err = flattenFieldSchema(d, found)
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
