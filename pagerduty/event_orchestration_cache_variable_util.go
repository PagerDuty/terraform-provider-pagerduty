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

var resourceEventOrchestrationCacheVariableConditionSchema = map[string]*schema.Schema{
	"expression": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var resourceEventOrchestrationCacheVariableConfigurationSchema = map[string]*schema.Schema{
	"type": {
		Type:     schema.TypeString,
		Required: true,
		ValidateDiagFunc: validateValueDiagFunc([]string{
			"recent_value",
			"trigger_event_count",
		}),
	},
	"regex": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"source": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"ttl_seconds": {
		Type:     schema.TypeInt,
		Optional: true,
	},
}

var dataSourceEventOrchestrationCacheVariableConditionSchema = map[string]*schema.Schema{
	"expression": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var dataSourceEventOrchestrationCacheVariableConfigurationSchema = map[string]*schema.Schema{
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"regex": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"source": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ttl_seconds": {
		Type:     schema.TypeInt,
		Computed: true,
	},
}

func checkConfiguration(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
	t := diff.Get("configuration.0.type").(string)
	s := diff.Get("configuration.0.source").(string)
	r := diff.Get("configuration.0.regex").(string)
	ts := diff.Get("configuration.0.ttl_seconds").(int)

	if t == "recent_value" && (r == "" || s == "") {
		return fmt.Errorf("Invalid configuration: regex and source cannot be null when type is recent_value")
	}
	if t == "trigger_event_count" && ts == 0 {
		return fmt.Errorf("Invalid configuration: ttl_seconds cannot be null when type is trigger_event_count")
	}
	if (r != "" || s != "") && ts != 0 {
		return fmt.Errorf("Invalid configuration: ttl_seconds cannot be used in conjuction with regex and source")
	}
	return nil
}

func getIdentifier(cacheVariableType string) string {
	switch cacheVariableType {
	case pagerduty.CacheVariableTypeGlobal:
		return "event_orchestration"
	case pagerduty.CacheVariableTypeService:
		return "service"
	}
	return ""
}

func setEventOrchestrationCacheVariableProps(d *schema.ResourceData, cv *pagerduty.EventOrchestrationCacheVariable) error {
	d.Set("name", cv.Name)
	d.Set("disabled", cv.Disabled)
	d.Set("configuration", flattenEventOrchestrationCacheVariableConfiguration(cv.Configuration))
	d.Set("condition", flattenEventOrchestrationCacheVariableConditions(cv.Conditions))

	return nil
}

func getEventOrchestrationCacheVariablePayloadData(d *schema.ResourceData, cacheVariableType string) (string, *pagerduty.EventOrchestrationCacheVariable) {
	orchestrationId := d.Get(getIdentifier(cacheVariableType)).(string)

	cacheVariable := &pagerduty.EventOrchestrationCacheVariable{
		Name:          d.Get("name").(string),
		Conditions:    expandEventOrchestrationCacheVariableConditions(d.Get("condition")),
		Configuration: expandEventOrchestrationCacheVariableConfiguration(d.Get("configuration")),
		Disabled:      d.Get("disabled").(bool),
	}

	return orchestrationId, cacheVariable
}

func fetchPagerDutyEventOrchestrationCacheVariable(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string, oid string, id string) (*pagerduty.EventOrchestrationCacheVariable, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return nil, err
	}

	if cacheVariable, _, err := client.EventOrchestrationCacheVariables.Get(ctx, cacheVariableType, oid, id); err != nil {
		return nil, err
	} else if cacheVariable != nil {
		d.SetId(id)
		d.Set(getIdentifier(cacheVariableType), oid)
		setEventOrchestrationCacheVariableProps(d, cacheVariable)
		return cacheVariable, nil
	}

	return nil, fmt.Errorf("Reading Cache Variable '%s' for PagerDuty Event Orchestration '%s' returned `nil`.", id, oid)
}

/*
****

	Schema expand and flatten helpers

****
*/
func expandEventOrchestrationCacheVariableConditions(v interface{}) []*pagerduty.EventOrchestrationCacheVariableCondition {
	conditions := []*pagerduty.EventOrchestrationCacheVariableCondition{}

	for _, cond := range v.([]interface{}) {
		c := cond.(map[string]interface{})

		cx := &pagerduty.EventOrchestrationCacheVariableCondition{
			Expression: c["expression"].(string),
		}

		conditions = append(conditions, cx)
	}

	return conditions
}

func expandEventOrchestrationCacheVariableConfiguration(v interface{}) *pagerduty.EventOrchestrationCacheVariableConfiguration {
	conf := &pagerduty.EventOrchestrationCacheVariableConfiguration{}

	for _, i := range v.([]interface{}) {
		c := i.(map[string]interface{})

		conf.Type = c["type"].(string)
		conf.Regex = c["regex"].(string)
		conf.Source = c["source"].(string)
		conf.TTLSeconds = c["ttl_seconds"].(int)
	}

	return conf
}

func flattenEventOrchestrationCacheVariableConditions(conds []*pagerduty.EventOrchestrationCacheVariableCondition) []interface{} {
	var flattenedConds []interface{}

	for _, cond := range conds {
		flattenedCond := map[string]interface{}{
			"expression": cond.Expression,
		}
		flattenedConds = append(flattenedConds, flattenedCond)
	}

	return flattenedConds
}

func flattenEventOrchestrationCacheVariableConfiguration(conf *pagerduty.EventOrchestrationCacheVariableConfiguration) []interface{} {
	result := map[string]interface{}{
		"type":        conf.Type,
		"regex":       conf.Regex,
		"source":      conf.Source,
		"ttl_seconds": conf.TTLSeconds,
	}

	return []interface{}{result}
}

/*
****

	Resource contexts

****
*/
func resourceEventOrchestrationCacheVariableImport(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string) ([]*schema.ResourceData, error) {
	oid, id, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	parent_identifier := "orchestration_id"

	if cacheVariableType == pagerduty.CacheVariableTypeService {
		parent_identifier = "service_id"
	}

	if oid == "" || id == "" {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing cache variable. Expected import ID format: <%s>:<cache_variable_id>", parent_identifier)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading Cache Variable '%s' for PagerDuty Event Orchestration: %s", id, oid)

		if _, err := fetchPagerDutyEventOrchestrationCacheVariable(ctx, d, meta, cacheVariableType, oid, id); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}
		return nil
	})

	if retryErr != nil {
		return []*schema.ResourceData{}, retryErr
	}

	return []*schema.ResourceData{d}, nil
}

func resourceEventOrchestrationCacheVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	oid, payload := getEventOrchestrationCacheVariablePayloadData(d, cacheVariableType)

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Creating Cache Variable '%s' for PagerDuty Event Orchestration '%s'", payload.Name, oid)

		if cacheVariable, _, err := client.EventOrchestrationCacheVariables.Create(ctx, cacheVariableType, oid, payload); err != nil {
			if isErrCode(err, http.StatusBadRequest) || isErrCode(err, http.StatusNotFound) || isErrCode(err, http.StatusForbidden) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		} else if cacheVariable != nil {
			// Try reading an cache variable after creation, retry if not found:
			if _, readErr := fetchPagerDutyEventOrchestrationCacheVariable(ctx, d, meta, cacheVariableType, oid, cacheVariable.ID); readErr != nil {
				log.Printf("[WARN] Cannot locate Cache Variable '%s' on PagerDuty Event Orchestration '%s'. Retrying creation...", cacheVariable.ID, oid)
				return retry.RetryableError(readErr)
			}
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourceEventOrchestrationCacheVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string) diag.Diagnostics {
	id := d.Id()
	oid := d.Get(getIdentifier(cacheVariableType)).(string)

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading Cache Variable '%s' for PagerDuty Event Orchestration: %s", id, oid)

		if _, err := fetchPagerDutyEventOrchestrationCacheVariable(ctx, d, meta, cacheVariableType, oid, id); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			} else if isErrCode(err, http.StatusNotFound) {
				log.Printf("[WARN] Removing %s because it's gone", d.Id())
				d.SetId("")
				return nil
			}

			return retry.RetryableError(err)
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourceEventOrchestrationCacheVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	oid, payload := getEventOrchestrationCacheVariablePayloadData(d, cacheVariableType)

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Updating Cache Variable '%s' for PagerDuty Event Orchestration: %s", id, oid)
		if cacheVariable, _, err := client.EventOrchestrationCacheVariables.Update(ctx, cacheVariableType, oid, id, payload); err != nil {
			if isErrCode(err, http.StatusBadRequest) || isErrCode(err, http.StatusNotFound) || isErrCode(err, http.StatusForbidden) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else if cacheVariable != nil {
			d.SetId(cacheVariable.ID)
			setEventOrchestrationCacheVariableProps(d, cacheVariable)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourceEventOrchestrationCacheVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	oid, _ := getEventOrchestrationCacheVariablePayloadData(d, cacheVariableType)

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Deleting Cache Variable '%s' for PagerDuty Event Orchestration: %s", id, oid)
		if _, err := client.EventOrchestrationCacheVariables.Delete(ctx, cacheVariableType, oid, id); err != nil {
			if isErrCode(err, http.StatusBadRequest) || isErrCode(err, http.StatusNotFound) || isErrCode(err, http.StatusForbidden) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else {
			// Try reading an cache variable after deletion, retry if still found:
			if cacheVariable, _, readErr := client.EventOrchestrationCacheVariables.Get(ctx, cacheVariableType, oid, id); readErr == nil && cacheVariable != nil {
				log.Printf("[WARN] Cache Variable '%s' still exists on PagerDuty Event Orchestration '%s'. Retrying deletion...", id, oid)
				return retry.RetryableError(fmt.Errorf("Cache Variable '%s' still exists on PagerDuty Event Orchestration '%s'.", id, oid))
			}
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	d.SetId("")

	return nil
}

/*
****

	Data Source context and helpers

****
*/

func dataSourceEventOrchestrationCacheVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cacheVariableType string) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return diag.FromErr(fmt.Errorf("Invalid Event Orchestration Cache Variable data source configuration: ID and name cannot both be null"))
	}

	oid := d.Get(getIdentifier(cacheVariableType)).(string)

	if id != "" && name != "" {
		diag := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Event Orchestration Cache Variable data source has both the ID and name attributes configured. Using ID '%s' to read data.", id),
		}
		diags = append(diags, diag)
	}

	if id != "" {
		return getEventOrchestrationCacheVariableById(ctx, d, meta, diags, cacheVariableType, oid, id)
	}

	return getEventOrchestrationCacheVariableByName(ctx, d, meta, diags, cacheVariableType, oid, name)
}

func getEventOrchestrationCacheVariableById(ctx context.Context, d *schema.ResourceData, meta interface{}, diags diag.Diagnostics, cacheVariableType string, oid, id string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading Cache Variable data source by ID '%s' for PagerDuty Event Orchestration '%s'", id, oid)

		if cacheVariable, _, err := client.EventOrchestrationCacheVariables.Get(ctx, cacheVariableType, oid, id); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else if cacheVariable != nil {
			d.SetId(cacheVariable.ID)
			setEventOrchestrationCacheVariableProps(d, cacheVariable)
		}
		return nil
	})

	if retryErr != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to find a Cache Variable with ID '%s' on PagerDuty Event Orchestration '%s'", id, oid),
		}
		return append(diags, diag)
	}

	return diags
}

func getEventOrchestrationCacheVariableByName(ctx context.Context, d *schema.ResourceData, meta interface{}, diags diag.Diagnostics, cacheVariableType string, oid, name string) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		log.Printf("[INFO] Reading Cache Variable data source by name '%s' for PagerDuty Event Orchestration '%s'", name, oid)

		resp, _, err := client.EventOrchestrationCacheVariables.List(ctx, cacheVariableType, oid)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}

		var matches []*pagerduty.EventOrchestrationCacheVariable

		for _, i := range resp.CacheVariables {
			if i.Name == name {
				matches = append(matches, i)
			}
		}

		count := len(matches)

		if count == 0 {
			return retry.NonRetryableError(
				fmt.Errorf("Unable to find a Cache Variable on Event Orchestration '%s' with name '%s'", oid, name),
			)
		}

		// This case should theoretically be impossible since Cache Variables must have
		// unique names per Event Orchestration
		if count > 1 {
			return retry.NonRetryableError(
				fmt.Errorf("Ambiguous Cache Variable name: '%s'. Found %v Cache Variables with this name on Event Orchestration '%s'. Please use the Cache Variable ID instead or make Cache Variable names unique within Event Orchestration.", name, count, oid),
			)
		}

		found := matches[0]
		d.SetId(found.ID)
		setEventOrchestrationCacheVariableProps(d, found)

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return diags
}
