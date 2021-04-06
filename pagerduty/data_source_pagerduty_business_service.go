package pagerduty

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyBusinessService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyBusinessServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyBusinessServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty business service")

	searchName := d.Get("name").(string)

	resp, _, err := client.BusinessServices.List()
	if err != nil {
		return err
	}

	var found *pagerduty.BusinessService

	for _, businessService := range resp.BusinessServices {
		if businessService.Name == searchName {
			found = businessService
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any business service with the name: %s", searchName)
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)

	return nil
}
