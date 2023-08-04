package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyServiceIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyServiceIntegrationRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_summary": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "examples 'Amazon CloudWatch', 'New Relic",
			},

			"integration_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourcePagerDutyServiceIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty service")

	searchName := d.Get("service_name").(string)

	o := &pagerduty.ListServicesOptions{
		Query: searchName,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Services.List(o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			return handleError(err)
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
				fmt.Errorf("unable to locate any service with the name: %s", searchName),
			)
		}

		integrationSummary := d.Get("integration_summary").(string)
		for _, integration := range found.Integrations {
			if strings.EqualFold(integration.Summary, integrationSummary) {
				integrationDetails, _, err := client.Services.GetIntegration(found.ID, integration.ID, &pagerduty.GetIntegrationOptions{})
				if err != nil {
					return handleError(err)
				}
				d.SetId(integration.ID)
				d.Set("service_name", found.Name)
				d.Set("integration_key", integrationDetails.IntegrationKey)

				return nil
			}

		}
		return resource.NonRetryableError(
			fmt.Errorf("unable to locate any integration of type %s on service %s", integrationSummary, searchName),
		)
	})
}

func handleError(err error) *resource.RetryError {
	time.Sleep(30 * time.Second)
	return resource.RetryableError(err)
}
