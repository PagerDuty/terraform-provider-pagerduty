package pagerduty

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEventOrchestrations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyEventOrchestrationsRead,

		Schema: map[string]*schema.Schema{
			"search": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_orchestrations": {
				Type:     schema.TypeList,
				Computed: true,
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
						"integration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
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
				},
			},
		},
	}
}

func dataSourcePagerDutyEventOrchestrationsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty Event Orchestrations")

	searchName := d.Get("search").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.EventOrchestrations.List()
		if err != nil {
			return resource.RetryableError(err)
		}

		var orchestrations []*pagerduty.EventOrchestration
		re := regexp.MustCompile(searchName)
		for _, orchestration := range resp.Orchestrations {
			if re.MatchString(orchestration.Name) {
				// Get orchestration matched by ID so we can set the integrations property
				// since the list endpoint does not return it
				orch, _, err := client.EventOrchestrations.Get(orchestration.ID)
				if err != nil {
					return resource.RetryableError(err)
				}
				orchestrations = append(orchestrations, orch)
			}
		}

		if len(orchestrations) == 0 {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any Event Orchestration with the name: %s", searchName),
			)
		}

		d.SetId(resource.UniqueId())
		d.Set("search", searchName)
		d.Set("event_orchestrations", flattenPagerDutyEventOrchestrations(orchestrations))

		return nil
	})
}

func flattenPagerDutyEventOrchestrations(orchestrations []*pagerduty.EventOrchestration) []interface{} {
	var result []interface{}

	for _, o := range orchestrations {
		orchestration := map[string]interface{}{
			"id":          o.ID,
			"name":        o.Name,
			"integration": flattenEventOrchestrationIntegrations(o.Integrations),
		}
		result = append(result, orchestration)
	}
	return result
}
