package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyIncidentCustomFieldOption() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyIncidentCustomFieldOptionRead,
		UpdateContext: resourcePagerDutyIncidentCustomFieldOptionUpdate,
		DeleteContext: resourcePagerDutyIncidentCustomFieldOptionDelete,
		CreateContext: resourcePagerDutyIncidentCustomFieldOptionCreate,
		// this function does not actually customize the diff but uses this hook
		// to validate the combination of datatype and value.
		CustomizeDiff: validateIncidentCustomFieldOptionValue,
		Schema: map[string]*schema.Schema{
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateValueDiagFunc([]string{pagerduty.IncidentCustomFieldDataTypeString.String()}),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func validateIncidentCustomFieldOptionValue(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	datatype := pagerduty.IncidentCustomFieldDataTypeFromString(diff.Get("data_type").(string))
	value := diff.Get("value").(string)

	generateError := func() error {
		return fmt.Errorf("invalid value for data_type %v: %v", datatype, value)
	}

	return validateIncidentCustomFieldValue(value, datatype, false, generateError)
}

func resourcePagerDutyIncidentCustomFieldOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	fieldID, fieldOption, err := buildFieldOptionStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty incident custom field option %s: %s", fieldID, fieldOption.Data.Value)

	createdFieldOption, _, err := client.IncidentCustomFields.CreateFieldOptionContext(ctx, fieldID, fieldOption)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldOption(d, fieldID, createdFieldOption)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentCustomFieldOptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fieldID := d.Get("field").(string)
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.IncidentCustomFields.DeleteFieldOptionContext(ctx, fieldID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentCustomFieldOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	fieldID, fieldOption, err := buildFieldOptionStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty incident custom field Option %s:%s", fieldID, d.Id())

	updatedFieldOption, _, err := client.IncidentCustomFields.UpdateFieldOptionContext(ctx, fieldID, d.Id(), fieldOption)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldOption(d, fieldID, updatedFieldOption)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentCustomFieldOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fieldID := d.Get("field").(string)
	log.Printf("[INFO] Reading PagerDuty incident custom field option %s:%s", fieldID, d.Id())
	err := fetchFieldOption(ctx, fieldID, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenFieldOption(d *schema.ResourceData, fieldID string, fieldOption *pagerduty.IncidentCustomFieldOption) error {
	value, err := convertIncidentCustomFieldValueForFlatten(fieldOption.Data.Value, false)
	if err != nil {
		return err
	}

	d.SetId(fieldOption.ID)
	d.Set("field", fieldID)
	d.Set("data_type", fieldOption.Data.DataType.String())
	d.Set("value", value)
	return nil
}

func buildFieldOptionStruct(d *schema.ResourceData) (string, *pagerduty.IncidentCustomFieldOption, error) {
	fieldID := d.Get("field").(string)

	dt := pagerduty.IncidentCustomFieldDataTypeFromString(d.Get("data_type").(string))
	v, err := convertIncidentCustomFieldValueForBuild(d.Get("value").(string), dt, false)
	if err != nil {
		return fieldID, nil, err
	}

	fieldOption := pagerduty.IncidentCustomFieldOption{
		Data: &pagerduty.IncidentCustomFieldOptionData{
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

	return retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		fieldOption, _, err := client.IncidentCustomFields.GetFieldOptionContext(ctx, fieldID, d.Id())
		if err != nil {
			log.Printf("[WARN] Field option read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				if isErrCode(err, http.StatusBadRequest) {
					return retry.NonRetryableError(err)
				}

				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenFieldOption(d, fieldID, fieldOption); err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
}
