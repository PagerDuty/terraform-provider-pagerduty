package pagerduty

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		// CustomizeDiff: checkEventOrchestrationChange, // Cannot return diags, only error
		Schema: map[string]*schema.Schema{
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

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		log.Printf("[INFO] Creating Integration '%s' for PagerDuty Event Orchestration: %s", payload.Label, orchestrationId)

		if integration, _, err := client.EventOrchestrationIntegrations.Create(orchestrationId, payload); err != nil {
			return resource.RetryableError(err)
		} else if integration != nil {
			d.SetId(integration.ID)
			setEventOrchestrationIntegrationProps(d, integration)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diags
}

func resourcePagerDutyEventOrchestrationIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	// Testing warnings on `terraform plan`:
	if d.HasChange("event_orchestration") {
		diag := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  ">>> test warning summary",
			Detail:   ">>> test warning detail",
		}
		diags = append(diags, diag)
	}

	id := d.Id()
	orchestrationId := d.Get("event_orchestration").(string)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		log.Printf("[INFO] Reading Integration '%s' for PagerDuty Event Orchestration: %s", id, orchestrationId)

		if integration, _, err := client.EventOrchestrationIntegrations.Get(orchestrationId, id); err != nil {
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

func resourcePagerDutyEventOrchestrationIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	orchestrationId, payload := getEventOrchestrationIntegrationPayloadData(d)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		log.Printf("[INFO] Updating Integration '%s' for PagerDuty Event Orchestration: %s", id, orchestrationId)

		if integration, _, err := client.EventOrchestrationIntegrations.Update(orchestrationId, id, payload); err != nil {
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

func resourcePagerDutyEventOrchestrationIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	orchestrationId, _ := getEventOrchestrationIntegrationPayloadData(d)

	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if _, err := client.EventOrchestrationIntegrations.Delete(orchestrationId, id); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	d.SetId("")

	return diags
}

func resourcePagerDutyEventOrchestrationIntegrationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	id := d.Id()
	orchestrationId, _ := getEventOrchestrationIntegrationPayloadData(d)
	_, _, err = client.EventOrchestrationIntegrations.Get(orchestrationId, id)

	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(id)
	d.Set("event_orchestration", orchestrationId)

	return []*schema.ResourceData{d}, nil
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
