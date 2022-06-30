package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyResponsePlay() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyResponsePlayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the response play.",
			},
			"from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A valid PagerDuty user's email is needed for the response play API ",
			},
		},
	}
}

func dataSourcePagerDutyResponsePlayRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty Response Play")

	searchResponsePlay := d.Get("name").(string)

	o := &pagerduty.ListResponsePlayOptions{
		From: d.Get("from").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.ResponsePlays.List(o)
		if err != nil {
			if isErrCode(err, 429) {
				// Delaying retry by 30s as recommended by PagerDuty
				// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
				time.Sleep(30 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		var found *pagerduty.ResponsePlay

		for _, play := range resp.ResponsePlays {
			if play.Name == searchResponsePlay {
				found = play
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("unable to locate any response play with name: %s", searchResponsePlay),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)

		return nil
	})
}
