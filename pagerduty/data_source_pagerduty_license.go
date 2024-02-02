package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		licenses, _, err := client.Licenses.List()
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		id, name, description := d.Get("id").(string), d.Get("name").(string), d.Get("description").(string)
		found := findBestMatchLicense(licenses, id, name, description)

		if found == nil {
			ids := licensesToStringOfIds(licenses)
			return retry.NonRetryableError(
				fmt.Errorf("Unable to locate any license with ids in [%s] with the configured id: '%s', name: '%s' or description: '%s'", ids, id, name, description))
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

func licensesToStringOfIds(licenses []*pagerduty.License) string {
	ids := make([]string, len(licenses))
	for i, v := range licenses {
		ids[i] = v.ID
	}
	return strings.Join(ids, ", ")
}

func findBestMatchLicense(licenses []*pagerduty.License, id, name, description string) *pagerduty.License {
	var found *pagerduty.License
	for _, license := range licenses {
		if licenseIsExactMatch(license, id, name, description) {
			found = license
			break
		}
	}

	// If there is no exact match for a license, check for substring matches
	// This allows customers to use a term such as "Full User", which is included
	// in the names of all licenses that support creating full users. However,
	// if id is set then it must match with licenseIsExactMatch
	if id == "" && found == nil {
		for _, license := range licenses {
			if licenseContainsMatch(license, name, description) {
				found = license
				break
			}
		}
	}

	return found
}

func licenseIsExactMatch(license *pagerduty.License, id, name, description string) bool {
	if id != "" {
		return license.ID == id && matchesOrIsUnset(license.Name, name) && matchesOrIsUnset(license.Description, description)
	}
	return license.Name == name && license.Description == description
}

func matchesOrIsUnset(licenseAttr, config string) bool {
	return config == "" || config == licenseAttr
}

func licenseContainsMatch(license *pagerduty.License, name, description string) bool {
	return strings.Contains(license.Name, name) && strings.Contains(license.Description, description)
}
