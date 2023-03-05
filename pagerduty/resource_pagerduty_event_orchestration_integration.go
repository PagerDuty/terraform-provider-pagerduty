package pagerduty

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
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
			StateContext: resourcePagerDutyEventOrchestrationIntegrationImport,
		},
		Schema: map[string]*schema.Schema{
			"event_orchestration": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: addIntegrationMigrationWarning(),
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

func addIntegrationMigrationWarning() schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "Modifying the event_orchestration property of the 'pagerduty_event_orchestration_integration' resource will cause all future events sent with this integration's routing key to be evaluated against the new Event Orchestration.",
			AttributePath: p,
		})

		return diags
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
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	oid, payload := getEventOrchestrationIntegrationPayloadData(d)

	retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Creating Integration '%s' for PagerDuty Event Orchestration '%s'", payload.Label, oid)

		if integration, _, err := client.EventOrchestrationIntegrations.CreateContext(ctx, oid, payload); err != nil {
			if isErrCode(err, 400) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		} else if integration != nil {
			// Try reading an integration after creation, retry if not found:
			if _, _, readErr := client.EventOrchestrationIntegrations.GetContext(ctx, oid, integration.ID); readErr != nil {
				log.Printf("[WARN] Cannot locate Integration '%s' on PagerDuty Event Orchestration '%s'. Retrying creation...", integration.ID, oid)
				return resource.RetryableError(readErr)
			} else {
				d.SetId(integration.ID)
				setEventOrchestrationIntegrationProps(d, integration)
			}
		}
		return nil
	})

	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourcePagerDutyEventOrchestrationIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	oid := d.Get("event_orchestration").(string)

	retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Reading Integration '%s' for PagerDuty Event Orchestration: %s", id, oid)

		if integration, _, err := client.EventOrchestrationIntegrations.GetContext(ctx, oid, id); err != nil {
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

	return nil
}

func resourcePagerDutyEventOrchestrationIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()

	// Migrate integration if the event_orchestration property was modified
	if d.HasChange("event_orchestration") {
		o, n := d.GetChange("event_orchestration")
		sourceOrchId := o.(string)
		destinationOrchId := n.(string)

		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Migrating Event Orchestration Integration '%s': source - '%s', destination - '%s'", id, sourceOrchId, destinationOrchId)

			if _, _, err := client.EventOrchestrationIntegrations.MigrateFromOrchestrationContext(ctx, destinationOrchId, sourceOrchId, id); err != nil {
				if isErrCode(err, 400) {
					return resource.NonRetryableError(err)
				}
				return resource.RetryableError(err)
			}
			return nil
		})

		if retryErr != nil {
			time.Sleep(2 * time.Second)
			return diag.FromErr(retryErr)
		}
	}

	// Update migrated integration if the label property was modified
	if d.HasChange("label") {
		orchestrationId, payload := getEventOrchestrationIntegrationPayloadData(d)

		retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Updating Integration '%s' for PagerDuty Event Orchestration: %s", id, orchestrationId)

			if integration, _, err := client.EventOrchestrationIntegrations.UpdateContext(ctx, orchestrationId, id, payload); err != nil {
				if isErrCode(err, 400) {
					return resource.NonRetryableError(err)
				}
				return resource.RetryableError(err)
			} else if integration != nil {
				setEventOrchestrationIntegrationProps(d, integration)
			}
			return nil
		})

		if retryErr != nil {
			time.Sleep(2 * time.Second)
			return diag.FromErr(retryErr)
		}
	}

	return nil
}

func resourcePagerDutyEventOrchestrationIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	oid, _ := getEventOrchestrationIntegrationPayloadData(d)

	retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Deleting Integration '%s' for PagerDuty Event Orchestration: %s", id, oid)
		if _, err := client.EventOrchestrationIntegrations.DeleteContext(ctx, oid, id); err != nil {
			return resource.RetryableError(err)
		} else {
			// Try reading an integration after deletion, retry if still found:
			if integr, _, readErr := client.EventOrchestrationIntegrations.GetContext(ctx, oid, id); readErr == nil && integr != nil {
				log.Printf("[WARN] Integration '%s' still exists on PagerDuty Event Orchestration '%s'. Retrying deletion...", id, oid)
				return resource.RetryableError(fmt.Errorf("Integration '%s' still exists on PagerDuty Event Orchestration '%s'.", id, oid))
			}
		}
		return nil
	})

	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	d.SetId("")

	return nil
}

func resourcePagerDutyEventOrchestrationIntegrationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), ".")
	if len(ids) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_event_orchestration_integration. Expected import ID format: <orchestration_id>.<integration_id>")
	}
	oid, id := ids[0], ids[1]

	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	retryErr := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Reading Integration '%s' for PagerDuty Event Orchestration: %s", id, oid)

		if integration, _, err := client.EventOrchestrationIntegrations.GetContext(ctx, oid, id); err != nil {
			return resource.RetryableError(err)
		} else if integration != nil {
			d.SetId(id)
			d.Set("event_orchestration", oid)
			setEventOrchestrationIntegrationProps(d, integration)
		}
		return nil
	})

	if retryErr != nil {
		return []*schema.ResourceData{}, retryErr
	}

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
	d.Set("label", i.Label)
	d.Set("parameters", flattenEventOrchestrationIntegrationParameters(i.Parameters))

	return nil
}
