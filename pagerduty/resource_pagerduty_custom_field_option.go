package pagerduty

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyCustomFieldOption() *schema.Resource {
	depMsg := "The standalone custom field feature has been removed from PagerDuty's Public API. The incident_custom_field_option resource provides similar, but not identical, functionality."

	f := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		return diag.FromErr(errors.New(depMsg))
	}

	return &schema.Resource{
		ReadContext:   f,
		UpdateContext: f,
		DeleteContext: f,
		CreateContext: f,
		Schema: map[string]*schema.Schema{
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datatype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
