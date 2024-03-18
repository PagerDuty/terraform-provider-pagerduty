package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// Deprecated: Migrated to pagerdutyplugin.resourceBusinessService. Kept for testing purposes.
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
				Type:       schema.TypeString,
				Optional:   true,
				Default:    "business_service",
				Deprecated: "This will change to a computed resource in the next major release.",
				ValidateDiagFunc: validateValueDiagFunc([]string{
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

func resourcePagerDutyBusinessServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	retryErr := retry.Retry(5*time.Minute, func() *retry.RetryError {
		businessService, err := buildBusinessServiceStruct(d)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		log.Printf("[INFO] Creating PagerDuty business service %s", businessService.Name)
		if businessService, _, err = client.BusinessServices.Create(businessService); err != nil {
			return retry.RetryableError(err)
		} else if businessService != nil {
			d.SetId(businessService.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return resourcePagerDutyBusinessServiceRead(d, meta)
}

func resourcePagerDutyBusinessServiceRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty business service %s", d.Id())

	retryErr := retry.Retry(5*time.Minute, func() *retry.RetryError {
		if businessService, _, err := client.BusinessServices.Get(d.Id()); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			if err := handleNotFoundError(err, d); err == nil {
				return nil
			}

			return retry.RetryableError(err)
		} else if businessService != nil {
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
		return retryErr
	}

	return nil
}

func resourcePagerDutyBusinessServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	businessService, err := buildBusinessServiceStruct(d)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] poc: %v", businessService.PointOfContact)
	log.Printf("[DEBUG] point_of_contact: %v", d.Get("point_of_contact"))

	log.Printf("[INFO] Updating PagerDuty business service %s", d.Id())

	if _, _, err := client.BusinessServices.Update(d.Id(), businessService); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyBusinessServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty business service %s", d.Id())

	if _, err := client.BusinessServices.Delete(d.Id()); err != nil {
		return handleNotFoundError(err, d)
	}

	d.SetId("")

	return nil
}
