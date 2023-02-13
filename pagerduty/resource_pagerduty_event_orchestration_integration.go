package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyEventOrchestrationIntegrationCreate,
		ReadContext:   resourcePagerDutyEventOrchestrationIntegrationRead,
		UpdateContext: resourcePagerDutyEventOrchestrationIntegrationUpdate,
		DeleteContext: resourcePagerDutyEventOrchestrationIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyEventOrchestrationIntegrationImport,
		},
		CustomizeDiff: checkEventOrchestrationChange,
		Schema: map[string]*schema.Schema{
			map[string]*schema.Schema{
				"event_orchestration": {
					Type:     schema.TypeString,
					Required: true,
				},
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"label": {
					Type:     schema.TypeString,
					Required: true,
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
			}
		},
	}
}

func getEventOrchestrationIntegrationPayloadData(d *schema.ResourceData) (string, *pagerduty.EventOrchestrationIntegration) {
	orchestrationId := d.Get("event_orchestration").(string)
	
	integration := &pagerduty.EventOrchestrationIntegration{
		Label: d.Get("label").(string),
	}

	return orchestrationId, integration
}

func resourcePagerDutyEventOrchestrationIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	orchestrationId, payload := getEventOrchestrationIntegrationPayloadData(d)
	var integration *pagerduty.EventOrchestrationIntegration	

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		log.Printf("[INFO] Creating Integration '%s' for PagerDuty Event Orchestration: %s", payload.Label, orchestrationId)

		if response, _, err := client.EventOrchestrationIntegrations.Create(orchestrationId, payload); err != nil {
			return resource.RetryableError(err)
		} else if response != nil {
			d.SetId(response.Id)
			integration = response
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	setEventOrchestrationIntegrationProps(d, integration)

	return diags
}

func resourcePagerDutyEventOrchestrationIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	orchestrationId := d.Get("event_orchestration").(string)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {		
		log.Printf("[INFO] Reading Integration '%s' for PagerDuty Event Orchestration: %s", id, orchestrationId)

		if integration, _, err := client.EventOrchestrationIntegrations.Get(orchestrationId, id, t); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if integration != nil {
			setEventOrchestrationIntegrationProps(d, integration)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diags
}

// func resourcePagerDutyEventOrchestrationIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
// 	client, err := meta.(*Config).Client()
// 	if err != nil {
// 		return err
// 	}

// 	orchestration := buildEventOrchestrationStruct(d)

// 	log.Printf("[INFO] Updating PagerDuty Event Orchestration: %s", d.Id())

// 	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
// 		id := d.Id()
// 		orchestrationId := d.Get("event_orchestration").(string)
// 		log.Printf("[INFO] Updating Integration '%s' for PagerDuty Event Orchestration: %s", id, orchestrationId)

// 		if integration, _, err := client.EventOrchestrations.Update(d.Id(), orchestration); err != nil {
// 			if isErrCode(err, 400) || isErrCode(err, 429) {
// 				return resource.RetryableError(err)
// 			}
// 			return resource.NonRetryableError(err)
// 		}

// 		return nil
// 	})

// 	if retryErr != nil {
// 		return retryErr
// 	}

// 	return nil
// }
func resourcePagerDutyEventOrchestrationIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	orchestrationId, payload := getEventOrchestrationIntegrationPayloadData(d)
	var integration *pagerduty.EventOrchestrationIntegration

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		log.Printf("[INFO] Updating Integration '%s' for PagerDuty Event Orchestration: %s", id, orchestrationId)

		if response, _, err := client.EventOrchestrationIntegrations.Update(payload.Parent.ID, "global", payload); err != nil {
			return resource.RetryableError(err)
		} else if response != nil {
			d.SetId(response.OrchestrationPath.Parent.ID)
			globalPath = response.OrchestrationPath
			warnings = response.Warnings
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	setEventOrchestrationPathGlobalProps(d, globalPath)

	return diags
}

func resourcePagerDutyEventOrchestrationIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty Event Orchestration: %s", d.Id())
	if _, err := client.EventOrchestrations.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func flattenEventOrchestrationIntegrationParameters(p *pagerduty.EventOrchestrationIntegrationParameters) []interface{} {
	result := map[string]interface{}{
		"routing_key": p.RoutingKey,
		"type":        p.Type,
	}

	return []interface{}{result}
}

func setEventOrchestrationIntegrationProps(d *schema.ResourceData, i *pagerduty.EventOrchestrationIntegration) error {
	d.Set("id", i.ID)
	d.Set("label", i.Label)
	d.Set("parameters", flattenEventOrchestrationIntegrationParameters(i.Parameters))

	return nil
}
