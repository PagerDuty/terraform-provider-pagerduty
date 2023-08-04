package pagerduty

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyCustomFieldSchemaFieldConfiguration() *schema.Resource {
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
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_value_multi_value": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"default_value_datatype": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
