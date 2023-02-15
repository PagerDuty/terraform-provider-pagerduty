package pagerduty

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyCustomFieldSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyCustomFieldSchemaRead,
		UpdateContext: resourcePagerDutyCustomFieldSchemaUpdate,
		DeleteContext: resourcePagerDutyCustomFieldSchemaDelete,
		CreateContext: resourcePagerDutyCustomFieldSchemaCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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

func resourcePagerDutyCustomFieldSchemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	fieldSchema, err := buildFieldSchemaStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty field schema %s.", fieldSchema.Title)

	createdFieldSchema, _, err := client.CustomFieldSchemas.CreateContext(ctx, fieldSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldSchema(d, createdFieldSchema)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CustomFieldSchemas.DeleteContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	fieldSchema, err := buildFieldSchemaStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty field schema %s", d.Id())

	updatedFieldSchema, _, err := client.CustomFieldSchemas.UpdateContext(ctx, d.Id(), fieldSchema)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldSchema(d, updatedFieldSchema)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty field schema %s", d.Id())
	err := fetchFieldSchema(ctx, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func fetchFieldSchema(ctx context.Context, d *schema.ResourceData, meta interface{}, errorCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		fieldSchema, _, err := client.CustomFieldSchemas.GetContext(ctx, d.Id(), nil)
		if err != nil {
			log.Printf("[WARN] Field Schema read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenFieldSchema(d, fieldSchema); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil

	})
}

func flattenFieldSchema(d *schema.ResourceData, fs *pagerduty.CustomFieldSchema) error {
	d.SetId(fs.ID)
	d.Set("title", fs.Title)
	if fs.Description != nil {
		d.Set("description", *(fs.Description))
	}
	return nil
}

func buildFieldSchemaStruct(d *schema.ResourceData) (*pagerduty.CustomFieldSchema, error) {
	fieldSchema := pagerduty.CustomFieldSchema{
		Title: d.Get("title").(string),
	}
	if desc, ok := d.GetOk("description"); ok {
		str := desc.(string)
		fieldSchema.Description = &str
	}
	return &fieldSchema, nil
}
