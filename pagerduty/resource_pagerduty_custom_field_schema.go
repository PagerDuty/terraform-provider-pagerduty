package pagerduty

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyCustomFieldSchema() *schema.Resource {
	depMsg := "The standalone custom field schema feature has been removed from PagerDuty's Public API."

	f := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		return diag.FromErr(errors.New(depMsg))
	}

	return &schema.Resource{
		ReadContext:   f,
		UpdateContext: f,
		DeleteContext: f,
		CreateContext: f,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
