package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyServiceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyServiceStatusRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_incident_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyServiceStatusRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty service status")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListServicesOptions{
		Query: searchName,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Services.List(o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.Service

		for _, service := range resp.Services {
			if service.Name == searchName {
				found = service
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any service with the name: %s", searchName),
			)
		}

		var teams []map[string]interface{}
		for _, team := range found.Teams {
			teams = append(teams, map[string]interface{}{
				"id":   team.ID,
				"name": team.Summary,
			})
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("last_incident_timestamp", found.LastIncidentTimestamp)
		d.Set("status", found.Status)

		return nil
	})
}
