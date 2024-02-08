package pagerduty

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyIncidentCustomField() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyIncidentCustomFieldRead,
		UpdateContext: resourcePagerDutyIncidentCustomFieldUpdate,
		DeleteContext: resourcePagerDutyIncidentCustomFieldDelete,
		CreateContext: resourcePagerDutyIncidentCustomFieldCreate,
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
			"data_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIncidentCustomFieldDataType(),
			},
			"field_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIncidentCustomFieldFieldType(),
			},
			"default_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePagerDutyIncidentCustomFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := buildFieldStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty incident custom field %s", field.Name)

	createdField, _, err := client.IncidentCustomFields.CreateContext(ctx, field)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentCustomField(d, createdField)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentCustomFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.IncidentCustomFields.DeleteContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentCustomFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := buildFieldStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty incident custom field %s", d.Id())

	updatedField, _, err := client.IncidentCustomFields.UpdateContext(ctx, d.Id(), field)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenIncidentCustomField(d, updatedField)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentCustomFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		field, _, err := client.IncidentCustomFields.GetContext(ctx, d.Id(), nil)
		if err != nil {
			log.Printf("[WARN] Incident custom field read error")
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenIncidentCustomField(d, field); err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
}

func flattenIncidentCustomField(d *schema.ResourceData, field *pagerduty.IncidentCustomField) error {
	d.SetId(field.ID)
	d.Set("name", field.Name)
	if field.Description != nil {
		d.Set("description", *(field.Description))
	}
	d.Set("display_name", field.DisplayName)
	d.Set("data_type", field.DataType.String())
	d.Set("field_type", field.FieldType.String())

	if field.DefaultValue != nil {
		v, err := convertIncidentCustomFieldValueForFlatten(field.DefaultValue, field.FieldType.IsMultiValue())
		if err != nil {
			return err
		}
		d.Set("default_value", v)
	}
	return nil
}

func buildFieldStruct(d *schema.ResourceData) (*pagerduty.IncidentCustomField, error) {
	field := pagerduty.IncidentCustomField{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
		DataType:    pagerduty.IncidentCustomFieldDataTypeFromString(d.Get("data_type").(string)),
		FieldType:   pagerduty.IncidentCustomFieldFieldTypeFromString(d.Get("field_type").(string)),
	}
	if desc, ok := d.GetOk("description"); ok {
		str := desc.(string)
		field.Description = &str
	}
	if df, ok := d.GetOk("default_value"); ok {
		field.DefaultValue = df
	}
	return &field, nil
}
