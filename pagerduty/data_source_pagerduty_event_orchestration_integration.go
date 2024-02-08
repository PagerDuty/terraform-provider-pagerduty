package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyEventOrchestrationIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyEventOrchestrationIntegrationRead,
		Schema: map[string]*schema.Schema{
			"event_orchestration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourcePagerDutyEventOrchestrationIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Get("id").(string)
	lbl := d.Get("label").(string)

	if id == "" && lbl == "" {
		return diag.FromErr(fmt.Errorf("Invalid Event Orchestration Integration data source configuration: ID and label cannot both be null"))
	}

	oid := d.Get("event_orchestration").(string)

	if id != "" && lbl != "" {
		diag := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Event Orchestration Integration data source has both the ID and label attributes configured. Using ID '%s' to read data.", id),
		}
		diags = append(diags, diag)
	}

	if id != "" {
		return getEventOrchestrationIntegrationById(ctx, d, meta, diags, oid, id)
	}

	return getEventOrchestrationIntegrationByLabel(ctx, d, meta, diags, oid, lbl)
}

func getEventOrchestrationIntegrationById(ctx context.Context, d *schema.ResourceData, meta interface{}, diags diag.Diagnostics, oid, id string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading Integration data source by ID '%s' for PagerDuty Event Orchestration '%s'", id, oid)

		if integration, _, err := client.EventOrchestrationIntegrations.GetContext(ctx, oid, id); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		} else if integration != nil {
			d.SetId(integration.ID)
			setEventOrchestrationIntegrationProps(d, integration)
		}
		return nil
	})

	if retryErr != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to find an Integration with ID '%s' on PagerDuty Event Orchestration '%s'", id, oid),
		}
		return append(diags, diag)
	}

	return diags
}

func getEventOrchestrationIntegrationByLabel(ctx context.Context, d *schema.ResourceData, meta interface{}, diags diag.Diagnostics, oid, lbl string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading Integration data source by label '%s' for PagerDuty Event Orchestration '%s'", lbl, oid)

		resp, _, err := client.EventOrchestrationIntegrations.ListContext(ctx, oid)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		var matches []*pagerduty.EventOrchestrationIntegration

		for _, i := range resp.Integrations {
			if i.Label == lbl {
				matches = append(matches, i)
			}
		}

		count := len(matches)

		if count == 0 {
			return retry.NonRetryableError(
				fmt.Errorf("Unable to find an Integration on Event Orchestration '%s' with label '%s'", oid, lbl),
			)
		}

		if count > 1 {
			return retry.NonRetryableError(
				fmt.Errorf("Ambiguous Integration label: '%s'. Found %v Integrations with this label on Event Orchestration '%s'. Please use the Integration ID instead or make Integration labels unique within Event Orchestration.", lbl, count, oid),
			)
		}

		found := matches[0]
		d.SetId(found.ID)
		setEventOrchestrationIntegrationProps(d, found)

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diags
}
