package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyMaintenanceWindowCreate,
		ReadContext:   resourcePagerDutyMaintenanceWindowRead,
		UpdateContext: resourcePagerDutyMaintenanceWindowUpdate,
		DeleteContext: resourcePagerDutyMaintenanceWindowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateRFC3339,
				DiffSuppressFunc: suppressRFC3339Diff,
			},
			"end_time": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateRFC3339,
				DiffSuppressFunc: suppressRFC3339Diff,
			},

			"services": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
		},
	}
}

func buildMaintenanceWindowStruct(d *schema.ResourceData) *pagerduty.MaintenanceWindow {
	window := &pagerduty.MaintenanceWindow{
		StartTime: d.Get("start_time").(string),
		EndTime:   d.Get("end_time").(string),
		Services:  expandServices(d.Get("services").(*schema.Set)),
	}

	if v, ok := d.GetOk("description"); ok {
		window.Description = v.(string)
	}

	return window
}

func resourcePagerDutyMaintenanceWindowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	window := buildMaintenanceWindowStruct(d)

	log.Printf("[INFO] Creating PagerDuty maintenance window")

	window, _, err = client.MaintenanceWindows.Create(window)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(window.ID)

	return nil
}

func resourcePagerDutyMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty maintenance window %s", d.Id())

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		window, _, err := client.MaintenanceWindows.Get(d.Id())
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		d.Set("description", window.Description)
		d.Set("start_time", window.StartTime)
		d.Set("end_time", window.EndTime)

		if err := d.Set("services", flattenServices(window.Services)); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	}))
}

func resourcePagerDutyMaintenanceWindowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	window := buildMaintenanceWindowStruct(d)

	log.Printf("[INFO] Updating PagerDuty maintenance window %s", d.Id())

	if _, _, err := client.MaintenanceWindows.Update(d.Id(), window); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyMaintenanceWindowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty maintenance window %s", d.Id())

	if _, err := client.MaintenanceWindows.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func expandServices(v *schema.Set) []*pagerduty.ServiceReference {
	var services []*pagerduty.ServiceReference

	for _, srv := range v.List() {
		service := &pagerduty.ServiceReference{
			Type: "service_reference",
			ID:   srv.(string),
		}
		services = append(services, service)
	}

	return services
}

func flattenServices(v []*pagerduty.ServiceReference) *schema.Set {
	var services []interface{}

	for _, srv := range v {
		services = append(services, srv.ID)
	}

	return schema.NewSet(schema.HashString, services)
}
