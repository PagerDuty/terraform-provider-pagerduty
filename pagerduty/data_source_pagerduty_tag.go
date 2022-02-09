package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
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

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Tags.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr != nil {
			return checkErr
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
