package pagerduty

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyVendors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyVendorsRead,

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

func dataSourcePagerDutyVendorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _ := meta.(*Config).Client()

	log.Printf("[INFO] Reading all PagerDuty vendors")

	o := &pagerduty.ListVendorsOptions{}
	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Vendors.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
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
	}))
}
