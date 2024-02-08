package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"runner_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_seen": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"runbook_base_uri": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		runner, _, err := client.AutomationActionsRunner.Get(d.Get("id").(string))
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		d.SetId(runner.ID)
		d.Set("name", runner.Name)
		d.Set("type", runner.Type)
		d.Set("runner_type", runner.RunnerType)
		d.Set("creation_time", runner.CreationTime)

		if runner.Description != nil {
			d.Set("description", &runner.Description)
		}

		if runner.RunbookBaseUri != nil {
			d.Set("runbook_base_uri", &runner.RunbookBaseUri)
		}

		if runner.LastSeenTime != nil {
			d.Set("last_seen", &runner.LastSeenTime)
		}

		return nil
	})
}
