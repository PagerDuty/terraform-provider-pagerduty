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

func dataSourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyScheduleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty schedule")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListSchedulesOptions{
		Query: searchName,
	}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Schedules.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		var found *pagerduty.Schedule

		for _, schedule := range resp.Schedules {
			if schedule.Name == searchName {
				found = schedule
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any schedule with the name: %s", searchName),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)

		return nil
	}))
}
