package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePagerDutyAutomationActionsRunner() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyAutomationActionsRunnerRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"runner_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyAutomationActionsRunnerRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty automation actions runner")

	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		runner, _, err := client.AutomationActionsRunner.Get(d.Get("id").(string))
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		d.SetId(runner.ID)
		d.Set("name", runner.Name)
		d.Set("runner_type", runner.RunnerType)
		d.Set("description", runner.Description)

		return nil
	})
}
