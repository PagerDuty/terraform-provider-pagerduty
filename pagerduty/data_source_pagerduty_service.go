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

func dataSourcePagerDutyService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_resolve_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"acknowledgement_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alert_creation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"escalation_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"teams": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The set of teams associated with the service",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyServiceRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty service")

	searchName := d.Get("name").(string)

	lookupOffset := 0
	more := true
	o := &pagerduty.ListServicesOptions{
		Query: searchName,
		Limit: 100,
	}
	foundServices := make([]*pagerduty.Service, 0)

	retryErr := retry.Retry(5*time.Minute, func() *retry.RetryError {
		for more {
			o.Offset = lookupOffset
			resp, _, err := client.Services.List(o)
			if err != nil {
				if isErrCode(err, http.StatusBadRequest) {
					return retry.NonRetryableError(err)
				}

				return retry.RetryableError(err)
			}
			more = resp.More
			lookupOffset += resp.Limit

			foundServices = append(foundServices, resp.Services...)
		}

		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	var found *pagerduty.Service
	for _, service := range foundServices {
		if service.Name == searchName {
			found = service
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any service with the name: %s", searchName)
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
	d.Set("type", found.Type)
	d.Set("auto_resolve_timeout", found.AutoResolveTimeout)
	d.Set("acknowledgement_timeout", found.AcknowledgementTimeout)
	d.Set("alert_creation", found.AlertCreation)
	d.Set("description", found.Description)
	d.Set("teams", teams)
	d.Set("escalation_policy", found.EscalationPolicy.ID)

	return nil
}
