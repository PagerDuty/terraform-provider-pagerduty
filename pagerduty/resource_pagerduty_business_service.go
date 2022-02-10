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

func resourcePagerDutyBusinessService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyBusinessServiceCreate,
		ReadContext:   resourcePagerDutyBusinessServiceRead,
		UpdateContext: resourcePagerDutyBusinessServiceUpdate,
		DeleteContext: resourcePagerDutyBusinessServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"summary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:       schema.TypeString,
				Optional:   true,
				Default:    "business_service",
				Deprecated: "This will change to a computed resource in the next major release.",
				ValidateFunc: validateValueFunc([]string{
					"business_service",
				}),
			},
			"point_of_contact": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildBusinessServiceStruct(d *schema.ResourceData) (*pagerduty.BusinessService, error) {
	businessService := pagerduty.BusinessService{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("description"); ok {
		businessService.Description = attr.(string)
	}
	if attr, ok := d.GetOk("type"); ok {
		businessService.Type = attr.(string)
	}
	if attr, ok := d.GetOk("summary"); ok {
		businessService.Summary = attr.(string)
	}
	if attr, ok := d.GetOk("self"); ok {
		businessService.Self = attr.(string)
	}
	businessService.PointOfContact = d.Get("point_of_contact").(string)

	if attr, ok := d.GetOk("html_url"); ok {
		businessService.HTMLUrl = attr.(string)
	}
	if attr, ok := d.GetOk("team"); ok {
		businessService.Team = &pagerduty.BusinessServiceTeam{
			ID: attr.(string),
		}
	}

	return &businessService, nil
}

func resourcePagerDutyBusinessServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {

		businessService, err := buildBusinessServiceStruct(d)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		log.Printf("[INFO] Creating PagerDuty business service %s", businessService.Name)
		if businessService, _, err = client.BusinessServices.Create(businessService); err != nil {
			return resource.RetryableError(err)
		} else if businessService != nil {
			d.SetId(businessService.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	return resourcePagerDutyBusinessServiceRead(ctx, d, meta)
}

func resourcePagerDutyBusinessServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty business service %s", d.Id())

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		businessService, _, err := client.BusinessServices.Get(d.Id())
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		if businessService != nil {
			d.Set("name", businessService.Name)
			d.Set("html_url", businessService.HTMLUrl)
			d.Set("description", businessService.Description)
			d.Set("type", businessService.Type)
			d.Set("point_of_contact", businessService.PointOfContact)
			d.Set("summary", businessService.Summary)
			d.Set("self", businessService.Self)
			if businessService.Team != nil {
				d.Set("team", businessService.Team.ID)
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

func resourcePagerDutyBusinessServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	businessService, err := buildBusinessServiceStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] poc: %v", businessService.PointOfContact)
	log.Printf("[DEBUG] point_of_contact: %v", d.Get("point_of_contact"))

	log.Printf("[INFO] Updating PagerDuty business service %s", d.Id())

	if _, _, err := client.BusinessServices.Update(d.Id(), businessService); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyBusinessServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty business service %s", d.Id())

	if _, err := client.BusinessServices.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
