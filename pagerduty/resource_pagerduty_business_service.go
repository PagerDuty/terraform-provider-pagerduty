package pagerduty

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyBusinessService() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyBusinessServiceCreate,
		Read:   resourcePagerDutyBusinessServiceRead,
		Update: resourcePagerDutyBusinessServiceUpdate,
		Delete: resourcePagerDutyBusinessServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "business_service",
				ValidateFunc: validateValueFunc([]string{
					"business_service",
					"business_service_reference",
				}),
			},
			"point_of_contact": {
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
	if attr, ok := d.GetOk("point_of_contact"); ok {
		businessService.PointOfContact = attr.(string)
	}
	if attr, ok := d.GetOk("html_url"); ok {
		businessService.HTMLUrl = attr.(string)
	}

	return &businessService, nil
}

func resourcePagerDutyBusinessServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	businessService, err := buildBusinessServiceStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating PagerDuty business service %s", businessService.Name)

	businessService, _, err = client.BusinessServices.Create(businessService)
	if err != nil {
		return err
	}

	d.SetId(businessService.ID)

	return resourcePagerDutyBusinessServiceRead(d, meta)
}

func resourcePagerDutyBusinessServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty business service %s", d.Id())

	businessService, _, err := client.BusinessServices.Get(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("name", businessService.Name)
	d.Set("html_url", businessService.HTMLUrl)
	d.Set("description", businessService.Description)
	d.Set("type", businessService.Type)
	d.Set("point_of_contact", businessService.PointOfContact)
	d.Set("summary", businessService.Summary)
	d.Set("self", businessService.Self)

	return nil
}

func resourcePagerDutyBusinessServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	businessService, err := buildBusinessServiceStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating PagerDuty business service %s", d.Id())

	if _, _, err := client.BusinessServices.Update(d.Id(), businessService); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyBusinessServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty business service %s", d.Id())

	if _, err := client.BusinessServices.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
