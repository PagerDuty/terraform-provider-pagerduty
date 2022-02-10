package pagerduty

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"regexp"
	"strings"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyVendor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyVendorRead,

		Schema: map[string]*schema.Schema{
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

func dataSourcePagerDutyVendorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty vendor")

	searchName := d.Get("name").(string)

	o := &pagerduty.ListVendorsOptions{
		Query: searchName,
	}
	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Vendors.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
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
	}))
}
