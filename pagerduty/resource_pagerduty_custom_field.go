package pagerduty

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyCustomField() *schema.Resource {
	depMsg := "The standalone custom field feature has been removed from PagerDuty's Public API. The incident_custom_field resource provides similar, but not identical, functionality."

	f := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		return diag.FromErr(errors.New(depMsg))
	}

	return &schema.Resource{
		DeprecationMessage: depMsg,
		ReadContext:        f,
		UpdateContext:      f,
		DeleteContext:      f,
		CreateContext:      f,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datatype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"multi_value": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"fixed_options": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}
