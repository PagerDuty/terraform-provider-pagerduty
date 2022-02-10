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

func dataSourcePagerDutyUserContactMethod() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyUserContactMethodRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the contact method to find in the PagerDuty API",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the contact method",
			},

			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blacklisted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"country_code": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"device_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"send_short_email": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyUserContactMethodRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty user's contact method")

	userId := d.Get("user_id").(string)
	searchLabel := d.Get("label").(string)
	searchType := d.Get("type").(string)

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Users.ListContactMethods(userId)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		var found *pagerduty.ContactMethod

		for _, contactMethod := range resp.ContactMethods {
			if contactMethod.Label == searchLabel &&
				contactMethod.Type == searchType {
				found = contactMethod
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(fmt.Errorf("Unable to locate any contact methods with the label: %s", searchLabel))
		}

		d.SetId(found.ID)
		d.Set("address", found.Address)
		d.Set("blacklisted", found.BlackListed)
		d.Set("country_code", found.CountryCode)
		d.Set("device_type", found.DeviceType)
		d.Set("enabled", found.Enabled)
		d.Set("label", found.Label)
		d.Set("send_short_email", found.SendShortEmail)
		d.Set("type", found.Type)

		return nil
	}))
}
