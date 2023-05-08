package pagerduty

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePagerDutyFieldSchema() *schema.Resource {
	depMsg := "The custom field schema feature has been removed from PagerDuty's Public API."

	return &schema.Resource{
		DeprecationMessage: depMsg,
		ReadContext: func(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
			return diag.FromErr(errors.New(depMsg))
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
