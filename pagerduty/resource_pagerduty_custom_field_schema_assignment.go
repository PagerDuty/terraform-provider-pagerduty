package pagerduty

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyCustomFieldSchemaAssignment() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyCustomFieldSchemaAssignmentRead,
		CreateContext: resourcePagerDutyCustomFieldSchemaAssignmentCreate,
		DeleteContext: resourcePagerDutyCustomFieldSchemaAssignmentDelete,
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

func resourcePagerDutyCustomFieldSchemaAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceID := d.Get("service")
	resp, _, err := client.CustomFieldSchemaAssignments.ListForServiceContext(ctx, serviceID.(string), nil)
	if err != nil {
		return diag.FromErr(err)
	} else {
		for _, fsa := range resp.SchemaAssignments {
			if fsa.ID == d.Id() {
				return nil
			}
		}
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return nil
	}
}

func resourcePagerDutyCustomFieldSchemaAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Removing PagerDuty field schema assignment %s", d.Id())

	_, err = client.CustomFieldSchemaAssignments.DeleteContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildFieldSchemaAssignmentStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty field schema assignment %s -> %s", a.Service.ID, a.Schema.ID)

	createdAssignment, _, err := client.CustomFieldSchemaAssignments.CreateContext(ctx, a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldSchemaAssignment(d, createdAssignment)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenFieldSchemaAssignment(d *schema.ResourceData, a *pagerduty.CustomFieldSchemaAssignment) error {
	d.SetId(a.ID)
	d.Set("schema", a.Schema.ID)
	d.Set("service", a.Service.ID)
	return nil
}

func buildFieldSchemaAssignmentStruct(d *schema.ResourceData) (*pagerduty.CustomFieldSchemaAssignment, error) {
	a := pagerduty.CustomFieldSchemaAssignment{
		Schema: &pagerduty.CustomFieldSchemaReference{
			ID:   d.Get("schema").(string),
			Type: "schema_reference",
		},
		Service: &pagerduty.ServiceReference{
			ID:   d.Get("service").(string),
			Type: "service_reference",
		},
	}

	return &a, nil
}
