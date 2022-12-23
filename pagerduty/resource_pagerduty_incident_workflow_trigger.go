package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyIncidentWorkflowTrigger() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyIncidentWorkflowTriggerRead,
		UpdateContext: resourcePagerDutyIncidentWorkflowTriggerUpdate,
		DeleteContext: resourcePagerDutyIncidentWorkflowTriggerDelete,
		CreateContext: resourcePagerDutyIncidentWorkflowTriggerCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: validateIncidentWorkflowTrigger,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"manual",
					"conditional",
				}),
			},
			"workflow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subscribed_to_all_services": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePagerDutyIncidentWorkflowTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	iwt, err := buildIncidentWorkflowTriggerStruct(d, true)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty incident workflow trigger %s for %s.", iwt.Type, iwt.Workflow.ID)

	createdWorkflowTrigger, _, err := client.IncidentWorkflowTriggers.CreateContext(ctx, iwt)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentWorkflowTrigger(d, createdWorkflowTrigger)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty incident workflow trigger %s", d.Id())
	err := fetchIncidentWorkflowTrigger(ctx, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	iwt, err := buildIncidentWorkflowTriggerStruct(d, false)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty incident workflow trigger %s", d.Id())

	updatedWorkflowTrigger, _, err := client.IncidentWorkflowTriggers.UpdateContext(ctx, d.Id(), iwt)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentWorkflowTrigger(d, updatedWorkflowTrigger)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.IncidentWorkflowTriggers.DeleteContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func validateIncidentWorkflowTrigger(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	triggerType := d.Get("type").(string)
	_, hadCondition := d.GetOk("condition")
	if triggerType == "manual" && hadCondition {
		return fmt.Errorf("when trigger type manual is used, condition must not be specified")
	}
	if triggerType == "conditional" && !hadCondition {
		return fmt.Errorf("when trigger type conditional is used, condition must be specified")
	}

	s, hadServices := d.GetOk("services")
	all := d.Get("subscribed_to_all_services").(bool)
	if all && hadServices && len(s.([]interface{})) > 0 {
		return fmt.Errorf("when subscribed_to_all_services is true, services must either be not defined or empty")
	}

	return nil
}

func fetchIncidentWorkflowTrigger(ctx context.Context, d *schema.ResourceData, meta interface{}, errorCallback func(err error, d *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		iwt, _, err := client.IncidentWorkflowTriggers.GetContext(ctx, d.Id())
		if err != nil {
			log.Printf("[WARN] Incident workflow trigger read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenIncidentWorkflowTrigger(d, iwt); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil

	})
}

func flattenIncidentWorkflowTrigger(d *schema.ResourceData, t *pagerduty.IncidentWorkflowTrigger) error {
	d.SetId(t.ID)
	d.Set("type", t.TriggerType.String())
	d.Set("workflow", t.Workflow.ID)
	d.Set("services", flattenIncidentWorkflowEnabledServices(t.Services))
	d.Set("subscribed_to_all_services", t.SubscribedToAllServices)
	if t.Condition != nil {
		d.Set("condition", t.Condition)
	}

	return nil
}

func flattenIncidentWorkflowEnabledServices(s []*pagerduty.ServiceReference) []string {
	services := make([]string, len(s))
	for i, v := range s {
		services[i] = v.ID
	}
	return services
}

func buildIncidentWorkflowTriggerStruct(d *schema.ResourceData, forUpdate bool) (*pagerduty.IncidentWorkflowTrigger, error) {
	iwt := pagerduty.IncidentWorkflowTrigger{
		SubscribedToAllServices: d.Get("subscribed_to_all_services").(bool),
	}

	if forUpdate {
		iwt.Workflow = &pagerduty.IncidentWorkflow{
			ID: d.Get("workflow").(string),
		}
		iwt.TriggerType = pagerduty.IncidentWorkflowTriggerTypeFromString(d.Get("type").(string))
	}

	if services, ok := d.GetOk("services"); ok {
		iwt.Services = buildIncidentWorkflowTriggerServices(services)
	}

	if condition, ok := d.GetOk("condition"); ok {
		condStr := condition.(string)
		iwt.Condition = &condStr
	}

	return &iwt, nil
}

func buildIncidentWorkflowTriggerServices(s interface{}) []*pagerduty.ServiceReference {
	services := s.([]interface{})
	newServices := make([]*pagerduty.ServiceReference, len(services))
	for i, v := range services {
		newServices[i] = &pagerduty.ServiceReference{
			ID: v.(string),
		}
	}
	return newServices
}
