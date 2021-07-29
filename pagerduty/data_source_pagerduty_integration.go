package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyIntegrationRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_type": {
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

func dataSourcePagerDutyIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty service")

	searchName := d.Get("service_name").(string)

	o := &pagerduty.ListServicesOptions{
		Query: searchName,
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Services.List(o)
		if err != nil {
			if isErrCode(err, 429) {
				// Delaying retry by 30s as recommended by PagerDuty
				// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
				time.Sleep(30 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
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

		integrationType := d.Get("integration_type").(string)
		for _, integration := range found.Integrations {
			if strings.EqualFold(integration.Summary, integrationType) {
				integrationDetails, _, err := client.Services.GetIntegration(found.ID, integration.ID, &pagerduty.GetIntegrationOptions{})
				if err != nil {
					if isErrCode(err, 429) {
						time.Sleep(30 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				d.SetId(integration.ID)
				d.Set("service_name", found.Name)
				d.Set("integration_key", integrationDetails.IntegrationKey)

				return nil
			}

		}
		return resource.NonRetryableError(
			fmt.Errorf("Unable to locate any integration of type %s on service %s", integrationType, searchName),
		)
	})
}
