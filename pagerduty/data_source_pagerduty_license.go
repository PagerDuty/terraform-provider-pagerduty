package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

var licenseSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"summary": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"description": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"role_group": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"current_value": {
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	},
	"allocations_available": {
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	},
	"valid_roles": {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"self": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"html_url": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
}

func dataSourcePagerDutyLicense() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePagerDutyLicenseRead,
		Schema: licenseSchema,
	}
}

func dataSourcePagerDutyLicenseRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Fetching PagerDuty Licenses")

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		licenses, _, err := client.Licenses.List()
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.License

		for _, license := range licenses {
			if licenseIsMatch(license, d) {
				found = license
				break
			}
		}

		if found == nil {
			id, name, description := d.Get("id").(string), d.Get("name").(string), d.Get("description").(string)

			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any license with the configured id: %s, name: %s or description: %s", id, name, description))
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("description", found.Description)
		d.Set("type", found.Type)
		d.Set("summary", found.Summary)
		d.Set("role_group", found.RoleGroup)
		d.Set("allocations_available", found.AllocationsAvailable)
		d.Set("current_value", found.CurrentValue)
		d.Set("valid_roles", found.ValidRoles)
		d.Set("self", found.Self)
		d.Set("html_url", found.HTMLURL)

		return nil
	})
}

func licenseIsMatch(license *pagerduty.License, d *schema.ResourceData) bool {
	id, name, description := d.Get("id").(string), d.Get("name").(string), d.Get("description").(string)

	if license.ID == id {
		return true
	}
	if strings.Contains(license.Name, name) && strings.Contains(license.Description, description) {
		return true
	}
	return false
}
