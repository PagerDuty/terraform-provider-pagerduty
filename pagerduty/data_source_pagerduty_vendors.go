package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyVendors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyVendorsRead,

		Schema: map[string]*schema.Schema{
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePagerDutyVendorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading all PagerDuty vendors")

	o := &pagerduty.ListVendorsOptions{}
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Vendors.List(o)
		if err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		}

		var ids []string
		var names []string
		var types []string

		for _, vendor := range resp.Vendors {
			ids = append(ids, vendor.ID)
			names = append(names, vendor.Name)
			types = append(types, vendor.GenericServiceType)
		}

		d.SetId(fmt.Sprintf("%d", len(resp.Vendors)))
		d.Set("ids", ids)
		d.Set("names", names)
		d.Set("types", types)

		return nil
	})
}
