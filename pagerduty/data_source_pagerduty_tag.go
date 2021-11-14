package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyTagRead,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The label of the tag to find in the PagerDuty API",
			},
		},
	}
}

func dataSourcePagerDutyTagRead(d *schema.ResourceData, meta interface{}) error {
	client, _ := meta.(*Config).Client()

	log.Printf("[INFO] Reading PagerDuty tag")

	searchTag := d.Get("label").(string)

	o := &pagerduty.ListTagsOptions{
		Query: searchTag,
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Tags.List(o)
		if err != nil {
			if isErrCode(err, 429) {
				// Delaying retry by 30s as recommended by PagerDuty
				// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
				time.Sleep(30 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		var found *pagerduty.Tag

		for _, tag := range resp.Tags {
			if tag.Label == searchTag {
				found = tag
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any tag with label: %s", searchTag),
			)
		}

		d.SetId(found.ID)
		d.Set("label", found.Label)

		return nil
	})
}
