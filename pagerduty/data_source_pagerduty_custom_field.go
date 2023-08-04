package pagerduty

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePagerDutyField() *schema.Resource {
	depMsg := "The standalone custom field feature has been removed from PagerDuty's Public API. The incident_custom_field data source provides similar, but not identical, functionality."

	return &schema.Resource{
		DeprecationMessage: depMsg,
		ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
			return diag.FromErr(errors.New(depMsg))
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datatype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multi_value": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"fixed_options": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
