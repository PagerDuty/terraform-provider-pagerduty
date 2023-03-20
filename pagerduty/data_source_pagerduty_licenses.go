package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyLicenses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyLicensesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"licenses": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourcePagerDutyLicensesRead(d *schema.ResourceData, meta interface{}) error {
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

		d.SetId(d.Get("name").(string))
		d.Set("licenses", flattenLicenses(licenses))
		return nil
	})
}

func flattenLicenses(licenses []*pagerduty.License) []map[string]interface{} {
	updated := make([]map[string]interface{}, len(licenses))
	for i, license := range licenses {
		updated[i] = flattenLicense(license)
	}

	return updated
}

func flattenLicense(l *pagerduty.License) map[string]interface{} {
	var license = map[string]interface{}{
		"id":                    l.ID,
		"type":                  l.Type,
		"name":                  l.Name,
		"description":           l.Description,
		"summary":               l.Summary,
		"role_group":            l.RoleGroup,
		"allocations_available": l.AllocationsAvailable,
		"current_value":         l.CurrentValue,
		"self":                  l.Self,
		"html_url":              l.HTMLURL,
		"valid_roles":           l.ValidRoles,
	}

	return license
}
