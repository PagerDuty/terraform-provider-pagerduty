package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyCustomFieldOption() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyCustomFieldOptionRead,
		UpdateContext: resourcePagerDutyCustomFieldOptionUpdate,
		DeleteContext: resourcePagerDutyCustomFieldOptionDelete,
		CreateContext: resourcePagerDutyCustomFieldOptionCreate,
		// this function does not actually customize the diff but uses this hook
		// to validate the combination of datatype and value.
		CustomizeDiff: validateCustomFieldOptionValue,
		Schema: map[string]*schema.Schema{
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datatype": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateCustomFieldDataType(),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func validateCustomFieldOptionValue(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	datatype := pagerduty.CustomFieldDataTypeFromString(diff.Get("datatype").(string))
	value := diff.Get("value").(string)

	generateError := func() error {
		return fmt.Errorf("invalid value for datatype %v: %v", datatype, value)
	}

	return validateCustomFieldValue(value, datatype, false, generateError)
}

func resourcePagerDutyCustomFieldOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	fieldID, fieldOption, err := buildFieldOptionStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty field option %s: %s", fieldID, fieldOption.Data.Value)

	createdFieldOption, _, err := client.CustomFields.CreateFieldOptionContext(ctx, fieldID, fieldOption)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldOption(d, fieldID, createdFieldOption)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fieldID := d.Get("field").(string)
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CustomFields.DeleteFieldOptionContext(ctx, fieldID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	fieldID, fieldOption, err := buildFieldOptionStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty field Option %s:%s", fieldID, d.Id())

	updatedFieldOption, _, err := client.CustomFields.UpdateFieldOptionContext(ctx, fieldID, d.Id(), fieldOption)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldOption(d, fieldID, updatedFieldOption)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fieldID := d.Get("field").(string)
	log.Printf("[INFO] Reading PagerDuty field option %s:%s", fieldID, d.Id())
	err := fetchFieldOption(ctx, fieldID, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenFieldOption(d *schema.ResourceData, fieldID string, fieldOption *pagerduty.CustomFieldOption) error {
	value, err := convertCustomFieldValueForFlatten(fieldOption.Data.Value, false)
	if err != nil {
		return err
	}

	d.SetId(fieldOption.ID)
	d.Set("field", fieldID)
	d.Set("datatype", fieldOption.Data.DataType.String())
	d.Set("value", value)
	return nil
}

func buildFieldOptionStruct(d *schema.ResourceData) (string, *pagerduty.CustomFieldOption, error) {
	fieldID := d.Get("field").(string)

	dt := pagerduty.CustomFieldDataTypeFromString(d.Get("datatype").(string))
	v, err := convertCustomFieldValueForBuild(d.Get("value").(string), dt, false)
	if err != nil {
		return fieldID, nil, err
	}

	fieldOption := pagerduty.CustomFieldOption{
		Data: &pagerduty.CustomFieldOptionData{
			DataType: dt,
			Value:    v,
		},
	}

	return fieldID, &fieldOption, nil
}

func fetchFieldOption(ctx context.Context, fieldID string, d *schema.ResourceData, meta interface{}, errorCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		fieldOption, _, err := client.CustomFields.GetFieldOptionContext(ctx, fieldID, d.Id())
		if err != nil {
			log.Printf("[WARN] Field option read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenFieldOption(d, fieldID, fieldOption); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil

	})
}
