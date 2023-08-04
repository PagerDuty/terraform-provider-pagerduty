package pagerduty

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyCustomFieldSchemaAssignment() *schema.Resource {
	depMsg := "The custom field schema feature has been removed from PagerDuty's Public API."

	f := func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		return diag.FromErr(errors.New(depMsg))
	}

	return &schema.Resource{
		ReadContext:   f,
		CreateContext: f,
		DeleteContext: f,
		Schema: map[string]*schema.Schema{
			"schema": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
