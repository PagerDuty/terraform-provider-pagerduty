package pagerduty

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyCustomFieldSchemaFieldConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyCustomFieldSchemaFieldConfigurationRead,
		UpdateContext: resourcePagerDutyCustomFieldSchemaFieldConfigurationUpdate,
		DeleteContext: resourcePagerDutyCustomFieldSchemaFieldConfigurationDelete,
		CreateContext: resourcePagerDutyCustomFieldSchemaFieldConfigurationCreate,
		// this function does not actually customize the diff but uses this hook
		// to validate the combination of required and the default_value prefixed
		// attributes
		CustomizeDiff: validateCustomFieldsForSchema,
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
				Type:             schema.TypeString,
				ValidateDiagFunc: validateCustomFieldDataTypeForFieldConfiguration(),
				Optional:         true,
			},
		},
	}
}

func resourcePagerDutyCustomFieldSchemaFieldConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	schemaID, fieldConfiguration, err := buildFieldSchemaFieldConfigurationStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty field configuration %s %s.", schemaID, fieldConfiguration.Field.ID)

	createdFieldConfiguration, _, err := client.CustomFieldSchemas.CreateFieldConfigurationContext(ctx, schemaID, fieldConfiguration)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldSchemaFieldConfiguration(d, schemaID, createdFieldConfiguration)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaFieldConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaID := d.Get("schema").(string)
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CustomFieldSchemas.DeleteFieldConfigurationContext(ctx, schemaID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaFieldConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	schemaID, fieldConfiguration, err := buildFieldSchemaFieldConfigurationStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty field schema field configuration %s", d.Id())

	updatedFieldConfiguration, _, err := client.CustomFieldSchemas.UpdateFieldConfigurationContext(ctx, schemaID, d.Id(), fieldConfiguration)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenFieldSchemaFieldConfiguration(d, schemaID, updatedFieldConfiguration)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyCustomFieldSchemaFieldConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaID := d.Get("schema").(string)
	log.Printf("[INFO] Reading PagerDuty field schema field configuration %s %s", schemaID, d.Id())
	err := fetchFieldSchemaFieldConfiguration(ctx, schemaID, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func fetchFieldSchemaFieldConfiguration(ctx context.Context, schemaID string, d *schema.ResourceData, meta interface{}, errorCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		fieldConfiguration, _, err := client.CustomFieldSchemas.GetFieldConfigurationContext(ctx, schemaID, d.Id(), nil)
		if err != nil {
			log.Printf("[WARN] Field Schema Field Configuration read error")
			errResp := errorCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenFieldSchemaFieldConfiguration(d, schemaID, fieldConfiguration); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil

	})
}

func flattenFieldSchemaFieldConfiguration(d *schema.ResourceData, schemaID string, fc *pagerduty.CustomFieldSchemaFieldConfiguration) error {
	if fc.DefaultValue != nil {
		v, err := convertCustomFieldValueForFlatten(fc.DefaultValue.Value, fc.DefaultValue.MultiValue)
		if err != nil {
			return err
		}
		d.Set("default_value_datatype", fc.DefaultValue.DataType.String())
		d.Set("default_value_multi_value", fc.DefaultValue.MultiValue)
		d.Set("default_value", v)
	}

	d.SetId(fc.ID)
	d.Set("schema", schemaID)
	d.Set("required", fc.Required)

	return nil
}

func buildFieldSchemaFieldConfigurationStruct(d *schema.ResourceData) (string, *pagerduty.CustomFieldSchemaFieldConfiguration, error) {
	schemaID := d.Get("schema").(string)
	fieldConfiguration := pagerduty.CustomFieldSchemaFieldConfiguration{
		Field: &pagerduty.CustomField{
			ID: d.Get("field").(string),
		},
	}

	required, hadRequired := d.GetOk("required")

	if hadRequired && required.(bool) {
		fieldConfiguration.Required = true

		dt := pagerduty.CustomFieldDataTypeFromString(d.Get("default_value_datatype").(string))
		mv := d.Get("default_value_multi_value").(bool)
		v, err := convertCustomFieldValueForBuild(d.Get("default_value").(string), dt, mv)
		if err == nil {
			fieldConfiguration.DefaultValue = &pagerduty.CustomFieldDefaultValue{
				DataType:   dt,
				MultiValue: mv,
				Value:      v,
			}
		} else {
			return schemaID, nil, err
		}
	}

	return schemaID, &fieldConfiguration, nil
}

func validateCustomFieldsForSchema(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	defaultValue, hadDefaultValue := d.GetOk("default_value")
	required, hadRequired := d.GetOk("required")
	if hadRequired && required.(bool) {
		if !hadDefaultValue || defaultValue == "" {
			return fmt.Errorf("required field without default value")
		}
	}

	if hadDefaultValue {
		defaultValueDataType, hadDefaultValueDataType := d.GetOk("default_value_datatype")
		if !hadDefaultValueDataType || defaultValueDataType == "" {
			return fmt.Errorf("required field without default value datatype")
		}
		dt := pagerduty.CustomFieldDataTypeFromString(defaultValueDataType.(string))
		if !dt.IsKnown() {
			// this is essentially impossible since default_value_datatype is validated
			return fmt.Errorf("unknown default value datatype: %s", defaultValueDataType)
		}
		defaultValueMultiValue, hadDefaultValueMultiValue := d.GetOk("default_value_multi_value")
		err := validateDefaultFieldValue(defaultValue.(string), dt, hadDefaultValueMultiValue && defaultValueMultiValue.(bool))
		if err != nil {
			return err
		}
	}
	return nil
}

func validateDefaultFieldValue(value string, datatype pagerduty.CustomFieldDataType, multiValue bool) error {
	generateError := func() error {
		if multiValue {
			return fmt.Errorf("invalid default value for datatype %v (multi-value): %v", datatype, value)
		} else {
			return fmt.Errorf("invalid default value for datatype %v: %v", datatype, value)
		}
	}

	return validateCustomFieldValue(value, datatype, multiValue, generateError)

}

func validateCustomFieldDataTypeForFieldConfiguration() schema.SchemaValidateDiagFunc {
	return func(v interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		dt := pagerduty.CustomFieldDataTypeFromString(v.(string))
		if !dt.IsKnown() {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Unknown datatype %v", v),
				AttributePath: p,
			})
		}
		return diags
	}
}
