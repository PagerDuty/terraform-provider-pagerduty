package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyScheduleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty schedule")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListSchedulesOptions{
		Query: searchName,
	}

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, _, err := client.Schedules.List(o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		var found *pagerduty.Schedule

		for _, schedule := range resp.Schedules {
			if schedule.Name == searchName {
				found = schedule
				break
			}
		}

		if found == nil {
			return retry.NonRetryableError(
				fmt.Errorf("Unable to locate any schedule with the name: %s", searchName),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)

		return nil
	})
}
