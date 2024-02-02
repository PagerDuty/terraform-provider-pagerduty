package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEventOrchestrations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyEventOrchestrationsRead,

		Schema: map[string]*schema.Schema{
			"name_filter": {
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

	nameFilter := d.Get("name_filter").(string)

	var eoList []*pagerduty.EventOrchestration
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.EventOrchestrations.List()
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}

		re, err := regexp.Compile(nameFilter)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("invalid regexp for name_filter provided %s", nameFilter))
		}
		for _, orchestration := range resp.Orchestrations {
			if re.MatchString(orchestration.Name) {
				eoList = append(eoList, orchestration)
			}
		}
		if len(eoList) == 0 {
			return retry.NonRetryableError(fmt.Errorf("Unable to locate any Event Orchestration matching the expression: %s", nameFilter))
		}

		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	var orchestrations []*pagerduty.EventOrchestration
	for _, orchestration := range eoList {
		// Get orchestration matched by ID so we can set the integrations property
		// since the list endpoint does not return it
		retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
			orch, _, err := client.EventOrchestrations.Get(orchestration.ID)
			if err != nil {
				if isErrCode(err, http.StatusBadRequest) {
					return retry.NonRetryableError(err)
				}

				return retry.RetryableError(err)
			}
			orchestrations = append(orchestrations, orch)
			return nil
		})
		if retryErr != nil {
			time.Sleep(2 * time.Second)
			return retryErr
		}
	}

	d.SetId(id.UniqueId())
	d.Set("name_filter", nameFilter)
	d.Set("event_orchestrations", flattenPagerDutyEventOrchestrations(orchestrations))

	return nil
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
