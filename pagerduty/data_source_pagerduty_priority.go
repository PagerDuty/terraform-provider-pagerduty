package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// Deprecated: Migrated to pagerdutyplugin.dataSourcePriority. Kept for testing purposes.
func dataSourcePagerDutyPriority() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyPriorityRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the priority to find in the PagerDuty API",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyPriorityRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty priority")

	searchTeam := d.Get("name").(string)

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, _, err := client.Priorities.List()
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		var found *pagerduty.Priority

		for _, priority := range resp.Priorities {
			if strings.EqualFold(priority.Name, searchTeam) {
				found = priority
				break
			}
		}

		if found == nil {
			return retry.NonRetryableError(
				fmt.Errorf("Unable to locate any priority with name: %s", searchTeam),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("description", found.Description)

		return nil
	})
}
