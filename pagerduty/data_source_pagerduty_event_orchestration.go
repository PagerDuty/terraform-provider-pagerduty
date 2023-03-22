package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEventOrchestration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyEventOrchestrationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true, // Tests keep failing if "Optional: true" is not provided
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"routing_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePagerDutyEventOrchestrationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty Event Orchestration")

	searchName := d.Get("name").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EventOrchestrations.List()
		if err != nil {
			return resource.RetryableError(err)
		}

		var found *pagerduty.EventOrchestration

		for _, orchestration := range resp.Orchestrations {
			if orchestration.Name == searchName {
				found = orchestration
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any Event Orchestration with the name: %s", searchName),
			)
		}

		// Get the found orchestration by ID so we can set the integrations property
		// since the list ndpoint does not return it
		orch, _, err := client.EventOrchestrations.Get(found.ID)
		if err != nil {
			return resource.RetryableError(err)
		}

		d.SetId(orch.ID)
		d.Set("name", orch.Name)

		if len(orch.Integrations) > 0 {
			d.Set("integration", flattenEventOrchestrationIntegrations(orch.Integrations))
		}

		return nil
	})
}
