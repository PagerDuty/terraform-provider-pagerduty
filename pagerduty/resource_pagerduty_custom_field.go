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

func resourcePagerDutyCustomField() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyCustomFieldRead,
		UpdateContext: resourcePagerDutyCustomFieldUpdate,
		DeleteContext: resourcePagerDutyCustomFieldDelete,
		CreateContext: resourcePagerDutyCustomFieldCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateCustomFieldDataType(),
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

func resourcePagerDutyCustomFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := buildFieldStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty field %s", field.Name)

	createdField, _, err := client.CustomFields.CreateContext(ctx, field)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenCustomField(d, createdField)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CustomFields.DeleteContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := buildFieldStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty field %s", d.Id())

	updatedField, _, err := client.CustomFields.UpdateContext(ctx, d.Id(), field)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenCustomField(d, updatedField)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty field %s", d.Id())
	err := fetchField(ctx, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func fetchField(ctx context.Context, d *schema.ResourceData, meta interface{}, errorCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		field, _, err := client.CustomFields.GetContext(ctx, d.Id(), nil)
		if err != nil {
			log.Printf("[WARN] Field read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenCustomField(d, field); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil

	})
}

func flattenCustomField(d *schema.ResourceData, field *pagerduty.CustomField) error {
	d.SetId(field.ID)
	d.Set("name", field.Name)
	if field.Description != nil {
		d.Set("description", *(field.Description))
	}
	d.Set("display_name", field.DisplayName)
	d.Set("datatype", field.DataType.String())
	d.Set("multi_value", field.MultiValue)
	d.Set("fixed_options", field.FixedOptions)
	return nil
}

func buildFieldStruct(d *schema.ResourceData) (*pagerduty.CustomField, error) {
	field := pagerduty.CustomField{
		Name:         d.Get("name").(string),
		DisplayName:  d.Get("display_name").(string),
		DataType:     pagerduty.CustomFieldDataTypeFromString(d.Get("datatype").(string)),
		MultiValue:   d.Get("multi_value").(bool),
		FixedOptions: d.Get("fixed_options").(bool),
	}
	if desc, ok := d.GetOk("description"); ok {
		str := desc.(string)
		field.Description = &str
	}
	return &field, nil
}
