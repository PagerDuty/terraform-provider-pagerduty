package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyMaintenanceWindowCreate,
		Read:   resourcePagerDutyMaintenanceWindowRead,
		Update: resourcePagerDutyMaintenanceWindowUpdate,
		Delete: resourcePagerDutyMaintenanceWindowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourcePagerDutyMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	window := buildMaintenanceWindowStruct(d)

	log.Printf("[INFO] Creating PagerDuty maintenance window")

	window, _, err = client.MaintenanceWindows.Create(window)
	if err != nil {
		return err
	}

	d.SetId(window.ID)

	return nil
}

func resourcePagerDutyMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty maintenance window %s", d.Id())

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		window, _, err := client.MaintenanceWindows.Get(d.Id())
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := handleNotFoundError(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}

			return nil
		}

		d.Set("description", window.Description)
		d.Set("start_time", window.StartTime)
		d.Set("end_time", window.EndTime)

		if err := d.Set("services", flattenServices(window.Services)); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})
}

func resourcePagerDutyMaintenanceWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	window := buildMaintenanceWindowStruct(d)

	log.Printf("[INFO] Updating PagerDuty maintenance window %s", d.Id())

	if _, _, err := client.MaintenanceWindows.Update(d.Id(), window); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyMaintenanceWindowDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty maintenance window %s", d.Id())

	if _, err := client.MaintenanceWindows.Delete(d.Id()); err != nil {
		// 405: The maintenance window can't be deleted because it has already ended. This can be considered deleted
		// from terraform's perspective.
		if !isErrCode(err, 405) {
			return err
		}
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
