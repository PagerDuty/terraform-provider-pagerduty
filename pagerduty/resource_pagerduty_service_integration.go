package pagerduty

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

const (
	errEmailIntegrationMustHaveEmail = "integration_email attribute must be set for an integration type generic_email_inbound_integration"
)

func resourcePagerDutyServiceIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyServiceIntegrationCreate,
		ReadContext:   resourcePagerDutyServiceIntegrationRead,
		UpdateContext: resourcePagerDutyServiceIntegrationUpdate,
		DeleteContext: resourcePagerDutyServiceIntegrationDelete,
		CustomizeDiff: func(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
			t := diff.Get("type").(string)
			if t == "generic_email_inbound_integration" && diff.Get("integration_email").(string) == "" && diff.NewValueKnown("integration_email") {
				return errors.New(errEmailIntegrationMustHaveEmail)
			}
			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourcePagerDutyServiceIntegrationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"vendor"},
				ValidateFunc: validateValueFunc([]string{
					"aws_cloudwatch_inbound_integration",
					"cloudkick_inbound_integration",
					"event_transformer_api_inbound_integration",
					"events_api_v2_inbound_integration",
					"generic_email_inbound_integration",
					"generic_events_api_inbound_integration",
					"keynote_inbound_integration",
					"nagios_inbound_integration",
					"pingdom_inbound_integration",
					"sql_monitor_inbound_integration",
				}),
			},
			"vendor": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"type"},
				Computed:      true,
			},
			"integration_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"integration_email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildServiceIntegrationStruct(d *schema.ResourceData) (*pagerduty.Integration, error) {
	serviceIntegration := &pagerduty.Integration{
		Name: d.Get("name").(string),
		Type: "service_integration",
		Service: &pagerduty.ServiceReference{
			Type: "service",
			ID:   d.Get("service").(string),
		},
	}

	if attr, ok := d.GetOk("integration_key"); ok {
		serviceIntegration.IntegrationKey = attr.(string)
	}

	if attr, ok := d.GetOk("integration_email"); ok {
		serviceIntegration.IntegrationEmail = attr.(string)
	}

	if attr, ok := d.GetOk("type"); ok {
		serviceIntegration.Type = attr.(string)
	}

	if attr, ok := d.GetOk("vendor"); ok {
		serviceIntegration.Vendor = &pagerduty.VendorReference{
			ID:   attr.(string),
			Type: "vendor",
		}
	}

	if serviceIntegration.Type == "generic_email_inbound_integration" && serviceIntegration.IntegrationEmail == "" {
		return nil, errors.New(errEmailIntegrationMustHaveEmail)
	}

	return serviceIntegration, nil
}

func fetchPagerDutyServiceIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}, handle404Errors bool) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	service := d.Get("service").(string)

	o := &pagerduty.GetIntegrationOptions{}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		serviceIntegration, _, err := client.Services.GetIntegration(service, d.Id(), o)
		if checkErr := getErrorHandler(handle404Errors)(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		d.Set("name", serviceIntegration.Name)
		d.Set("type", serviceIntegration.Type)

		if serviceIntegration.Service != nil {
			d.Set("service", serviceIntegration.Service.ID)
		}

		if serviceIntegration.Vendor != nil {
			d.Set("vendor", serviceIntegration.Vendor.ID)
		}

		if serviceIntegration.IntegrationKey != "" {
			d.Set("integration_key", serviceIntegration.IntegrationKey)
		}

		if serviceIntegration.IntegrationEmail != "" {
			d.Set("integration_email", serviceIntegration.IntegrationEmail)
		}

		if serviceIntegration.HTMLURL != "" {
			d.Set("html_url", serviceIntegration.HTMLURL)
		}

		return nil
	}))
}

func resourcePagerDutyServiceIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceIntegration, err := buildServiceIntegrationStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty service integration %s", serviceIntegration.Name)

	service := d.Get("service").(string)

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if serviceIntegration, _, err := client.Services.CreateIntegration(service, serviceIntegration); err != nil {
			if isErrCode(err, 400) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if serviceIntegration != nil {
			d.SetId(serviceIntegration.ID)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return fetchPagerDutyServiceIntegration(ctx, d, meta, false)
}

func resourcePagerDutyServiceIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty service integration %s", d.Id())
	return fetchPagerDutyServiceIntegration(ctx, d, meta, true)
}

func resourcePagerDutyServiceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceIntegration, err := buildServiceIntegrationStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	service := d.Get("service").(string)

	log.Printf("[INFO] Updating PagerDuty service integration %s", d.Id())

	if _, _, err := client.Services.UpdateIntegration(service, d.Id(), serviceIntegration); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyServiceIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	service := d.Get("service").(string)

	log.Printf("[INFO] Removing PagerDuty service integration %s", d.Id())

	if _, err := client.Services.DeleteIntegration(service, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourcePagerDutyServiceIntegrationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	ids := strings.Split(d.Id(), ".")

	if len(ids) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_service_integration. Expecting an importation ID formed as '<service_id>.<integration_id>'")
	}
	sid, id := ids[0], ids[1]

	_, _, err = client.Services.GetIntegration(sid, id, nil)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	// These are set because an import also calls Read behind the scenes
	d.SetId(id)
	d.Set("service", sid)

	return []*schema.ResourceData{d}, nil
}
