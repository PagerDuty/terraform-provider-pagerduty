package pagerduty

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyVendor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyVendorRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use `name` instead. This attribute will be removed in a future version",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyVendorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty vendor")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListVendorsOptions{
		Query: searchName,
	}
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Vendors.List(o)
		if err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.Vendor

		for _, vendor := range resp.Vendors {
			if strings.EqualFold(vendor.Name, searchName) {
				found = vendor
				break
			}
		}

		// We didn't find an exact match, so let's fallback to partial matching.
		if found == nil {
			pr := regexp.MustCompile("(?i)" + searchName)
			for _, vendor := range resp.Vendors {
				if pr.MatchString(vendor.Name) {
					found = vendor
					break
				}
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any vendor with the name: %s", searchName),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("type", found.GenericServiceType)

		return nil
	})
}
