package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyService() *schema.Resource {
	return &schema.Resource{
		Create:        resourcePagerDutyServiceCreate,
		Read:          resourcePagerDutyServiceRead,
		Update:        resourcePagerDutyServiceUpdate,
		Delete:        resourcePagerDutyServiceDelete,
		CustomizeDiff: customizePagerDutyServiceDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIsAllowedString(NoNonPrintableChars),
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"alert_creation": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					// Once migrated, alert_creation arguments previously defined as create_incidents would have been reported diffs for all matching services. As this is no longer configurable, opt to suppress this diff.
					return true
				},
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"create_alerts_and_incidents",
					"create_incidents",
				}),
			},
			"alert_grouping": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"time",
					"intelligent",
					"rules",
				}),
				Deprecated:    "Use `alert_grouping_parameters.type`",
				ConflictsWith: []string{"alert_grouping_parameters"},
			},
			"alert_grouping_timeout": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Use `alert_grouping_parameters.config.timeout`",
				ConflictsWith: []string{"alert_grouping_parameters"},
			},
			"alert_grouping_parameters": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				Deprecated:    "Use a resource `pagerduty_alert_grouping_setting` instead.\nFollow the migration guide at https://registry.terraform.io/providers/PagerDuty/pagerduty/latest/docs/resources/alert_grouping_setting#migration-from-alert_grouping_parameters",
				ConflictsWith: []string{"alert_grouping", "alert_grouping_timeout"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateDiagFunc: validateValueDiagFunc([]string{
								"time",
								"intelligent",
								"content_based",
							}),
						},
						"config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"fields": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"aggregate": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"all",
											"any",
										}),
									},
									"time_window": {
										Type:             schema.TypeInt,
										Optional:         true,
										Computed:         true,
										ValidateDiagFunc: validateTimeWindow,
									},
								},
							},
						},
					},
				},
			},
			"auto_pause_notifications_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{120, 180, 300, 600, 900})),
						},
					},
				},
			},
			"auto_resolve_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "14400",
			},
			"last_incident_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acknowledgement_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1800",
			},
			"escalation_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"incident_urgency_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"urgency": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"during_support_hours": {
							Type:     schema.TypeList,
							MaxItems: 1,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"urgency": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"outside_support_hours": {
							Type:     schema.TypeList,
							MaxItems: 1,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"urgency": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"support_hours": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				MinItems: 1,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_zone": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: util.ValidateTZValueDiagFunc,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"days_of_week": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 7,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"scheduled_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to_urgency": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"at": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"response_play": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func customizePagerDutyServiceDiff(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
	in := diff.Get("incident_urgency_rule.#").(int)
	for i := 0; i <= in; i++ {
		t := diff.Get(fmt.Sprintf("incident_urgency_rule.%d.type", i)).(string)
		if t == "use_support_hours" && diff.Get(fmt.Sprintf("incident_urgency_rule.%d.urgency", i)).(string) != "" {
			return fmt.Errorf("general urgency cannot be set for a use_support_hours incident urgency rule type")
		}
	}

	incidentUrgencyRuleType := diff.Get("incident_urgency_rule.0.type").(string)
	if incidentUrgencyRuleType == "use_support_hours" {
		if diff.Get("support_hours.#").(int) != 1 {
			return fmt.Errorf("when using type = use_support_hours in incident_urgency_rule you must specify exactly one (otherwise optional) support_hours block")
		}
	}

	// Due to alert_grouping_parameters.type = null is a valid configuration
	// for disabling Service's Alert Grouping configuration and having an
	// empty alert_grouping_parameters.config block is also valid, API ignore
	// this input fields, and turns out that API response for Service
	// configuration doesn't bring a representation of this HCL, which leads
	// to a permadiff, described in
	// https://github.com/PagerDuty/terraform-provider-pagerduty/issues/700
	//
	// So, bellow is the formated representation alert_grouping_parameters
	// value when this permadiff appears and must be ignored.
	ignoreThisAlertGroupingParamsConfigDiff := `[]interface {}{map[string]interface {}{"config":[]interface {}{interface {}(nil)}, "type":""}}`
	if agpdiff, ok := diff.Get("alert_grouping_parameters").([]interface{}); ok && diff.NewValueKnown("alert_grouping_parameters") && fmt.Sprintf("%#v", agpdiff) == ignoreThisAlertGroupingParamsConfigDiff {
		diff.Clear("alert_grouping_parameters")
	}

	if agpType, ok := diff.Get("alert_grouping_parameters.0.type").(string); ok {
		agppath := "alert_grouping_parameters.0.config.0."
		timeoutVal := diff.Get(agppath + "timeout").(int)
		aggregateVal := diff.Get(agppath + "aggregate").(string)
		fieldsVal := diff.Get(agppath + "fields").([]interface{})
		timeWindowVal := diff.Get(agppath + "time_window").(int)
		hasChangeAgpType := diff.HasChange("alert_grouping_parameters")

		if agpType == "content_based" && (aggregateVal == "" || len(fieldsVal) == 0) {
			return fmt.Errorf("When using Alert grouping parameters configuration of type \"content_based\" is in use, attributes \"aggregate\" and \"fields\" are required")
		}
		if timeWindowVal == 86400 && agpType != "content_based" {
			return fmt.Errorf("Alert grouping parameters configuration attribute \"time_window\" with a value of 86400 is only supported by \"content-based\" type Alert Grouping")
		}
		if (aggregateVal != "" || len(fieldsVal) > 0) && (agpType != "" && hasChangeAgpType && agpType != "content_based") {
			return fmt.Errorf("Alert grouping parameters configuration attributes \"aggregate\" and \"fields\" are only supported by \"content_based\" type Alert Grouping")
		}
		if timeoutVal > 0 && (agpType != "" && hasChangeAgpType && agpType != "time") {
			return fmt.Errorf("Alert grouping parameters configuration attribute \"timeout\" is only supported by \"time\" type Alert Grouping")
		}
		if (timeWindowVal > 300) && (agpType != "" && hasChangeAgpType && (agpType != "intelligent" && agpType != "content_based")) {
			return fmt.Errorf("Alert grouping parameters configuration attribute \"time_window\" is only supported by \"intelligent\" and \"content-based\" type Alert Grouping")
		}
	}

	return nil
}

func validateTimeWindow(v interface{}, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	tw := v.(int)
	if (tw < 300 || tw > 3600) && tw != 86400 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Alert grouping time window value must be between 300 and 3600 or exactly 86400(86400 is supported only for content-based alert grouping), current setting is %d", tw),
			AttributePath: p,
		})
	}
	return diags
}

func buildServiceStruct(d *schema.ResourceData) (*pagerduty.Service, error) {
	service := pagerduty.Service{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("description"); ok {
		service.Description = attr.(string)
	}

	if attr, ok := d.GetOk("auto_resolve_timeout"); ok {
		if attr.(string) != "null" {
			if val, err := strconv.Atoi(attr.(string)); err == nil {
				service.AutoResolveTimeout = &val
			} else {
				return nil, err
			}
		}
	}

	if attr, ok := d.GetOk("acknowledgement_timeout"); ok {
		if attr.(string) != "null" {
			if val, err := strconv.Atoi(attr.(string)); err == nil {
				service.AcknowledgementTimeout = &val
			} else {
				return nil, err
			}
		}
	}

	if attr, ok := d.GetOk("alert_creation"); ok {
		service.AlertCreation = attr.(string)
	}

	if attr, ok := d.GetOk("alert_grouping"); ok {
		ag := attr.(string)
		if ag != "rules" {
			service.AlertGrouping = &ag
		}
	}

	if attr, ok := d.GetOk("alert_grouping_parameters"); ok {
		service.AlertGroupingParameters = expandAlertGroupingParameters(attr)
	} else {
		// Clear AlertGroupingParameters as it takes precedence over AlertGrouping and AlertGroupingTimeout which are apparently deprecated (that's not explicitly documented in the API)
		service.AlertGroupingParameters = nil
	}

	if attr, ok := d.GetOk("alert_grouping_timeout"); ok {
		if attr.(string) != "null" {
			if val, err := strconv.Atoi(attr.(string)); err == nil {
				service.AlertGroupingTimeout = &val
			} else {
				return nil, err
			}
		}
	}
	if attr, ok := d.GetOk("auto_pause_notifications_parameters"); ok {
		service.AutoPauseNotificationsParameters = expandAutoPauseNotificationsParameters(attr)
	}

	if attr, ok := d.GetOk("escalation_policy"); ok {
		service.EscalationPolicy = &pagerduty.EscalationPolicyReference{
			ID:   attr.(string),
			Type: "escalation_policy_reference",
		}
	}

	if attr, ok := d.GetOk("incident_urgency_rule"); ok {
		service.IncidentUrgencyRule = expandIncidentUrgencyRule(attr)
		if service.IncidentUrgencyRule.Type == "use_support_hours" {
			service.ScheduledActions = make([]*pagerduty.ScheduledAction, 1)
		}
	}

	if attr, ok := d.GetOk("scheduled_actions"); ok {
		service.ScheduledActions = expandScheduledActions(attr)
	}

	if attr, ok := d.GetOk("support_hours"); ok {
		service.SupportHours = expandSupportHours(attr)
	}

	if attr, ok := d.GetOk("response_play"); ok {
		if attr.(string) != "null" {
			service.ResponsePlay = &pagerduty.ResponsePlayReference{
				ID:   attr.(string),
				Type: "response_play_reference",
			}
		}
	}
	return &service, nil
}

func fetchService(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		service, _, err := client.Services.Get(d.Id(), &pagerduty.GetServiceOptions{
			Includes: []string{"auto_pause_notifications_parameters"},
		})
		if err != nil {
			log.Printf("[WARN] Service read error")
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			if errResp != nil {
				return retry.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenService(d, service); err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
}

func resourcePagerDutyServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	service, err := buildServiceStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating PagerDuty service %s", service.Name)

	service, _, err = client.Services.Create(service)
	if err != nil {
		return err
	}

	d.SetId(service.ID)

	// We wait for internal subsystem to sync. Otherwise fields like
	// alert_grouping_parameters will return empty.
	time.Sleep(500 * time.Millisecond)

	return fetchService(d, meta, genError)
}

func resourcePagerDutyServiceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading PagerDuty service %s", d.Id())
	return fetchService(d, meta, handleNotFoundError)
}

func resourcePagerDutyServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	service, err := buildServiceStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating PagerDuty service %s", d.Id())

	updatedService, _, err := client.Services.Update(d.Id(), service)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	return flattenService(d, updatedService)
}

func resourcePagerDutyServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty service %s", d.Id())

	if _, err := client.Services.Delete(d.Id()); err != nil {
		return handleNotFoundError(err, d)
	}

	d.SetId("")

	// giving the API time to catchup
	time.Sleep(time.Second)
	return nil
}

func flattenService(d *schema.ResourceData, service *pagerduty.Service) error {
	d.Set("name", service.Name)
	d.Set("type", service.Type)
	d.Set("html_url", service.HTMLURL)
	d.Set("status", service.Status)
	d.Set("created_at", service.CreatedAt)
	d.Set("escalation_policy", service.EscalationPolicy.ID)
	d.Set("description", service.Description)
	if service.AutoResolveTimeout == nil {
		d.Set("auto_resolve_timeout", "null")
	} else {
		d.Set("auto_resolve_timeout", strconv.Itoa(*service.AutoResolveTimeout))
	}
	d.Set("last_incident_timestamp", service.LastIncidentTimestamp)
	if service.AcknowledgementTimeout == nil {
		d.Set("acknowledgement_timeout", "null")
	} else {
		d.Set("acknowledgement_timeout", strconv.Itoa(*service.AcknowledgementTimeout))
	}
	d.Set("alert_creation", service.AlertCreation)
	if service.AlertGrouping != nil && *service.AlertGrouping != "" {
		d.Set("alert_grouping", *service.AlertGrouping)
	}
	if service.AlertGroupingTimeout == nil {
		d.Set("alert_grouping_timeout", "null")
	} else {
		d.Set("alert_grouping_timeout", strconv.Itoa(*service.AlertGroupingTimeout))
	}

	_, hasGrouping := d.GetOk("alert_grouping")
	_, hasGroupingParams := d.GetOk("alert_grouping_parameters")
	if service.AlertGroupingParameters != nil && (!hasGrouping && hasGroupingParams) {
		if err := d.Set("alert_grouping_parameters", flattenAlertGroupingParameters(service.AlertGroupingParameters)); err != nil {
			return err
		}
	}

	if service.AutoPauseNotificationsParameters != nil {
		if err := d.Set("auto_pause_notifications_parameters", flattenAutoPauseNotificationsParameters(service.AutoPauseNotificationsParameters)); err != nil {
			return err
		}
	}

	if service.IncidentUrgencyRule != nil {
		if err := d.Set("incident_urgency_rule", flattenIncidentUrgencyRule(service.IncidentUrgencyRule)); err != nil {
			return err
		}
	}

	if service.SupportHours != nil {
		if err := d.Set("support_hours", flattenSupportHours(service.SupportHours)); err != nil {
			return err
		}
	}

	if service.ScheduledActions != nil {
		if err := d.Set("scheduled_actions", flattenScheduledActions(service.ScheduledActions)); err != nil {
			return err
		}
	}
	if service.ResponsePlay != nil {
		d.Set("response_play", service.ResponsePlay.ID)
	}
	return nil
}

func expandAlertGroupingParameters(v interface{}) *pagerduty.AlertGroupingParameters {
	alertGroupingParameters := &pagerduty.AlertGroupingParameters{
		Config: &pagerduty.AlertGroupingConfig{},
	}
	// First We capture a possible nil value for the interface to avoid the a
	// panic
	ragp, ok := v.([]interface{})
	if !ok || isNilFunc(ragp[0]) {
		return nil
	}
	ragpVal := ragp[0].(map[string]interface{})
	groupingType := ""
	if ragpVal["type"].(string) != "" {
		groupingType = ragpVal["type"].(string)
		alertGroupingParameters.Type = &groupingType
	}

	alertGroupingParameters.Config = expandAlertGroupingConfig(groupingType, ragpVal["config"])
	if groupingType == "content_based" && alertGroupingParameters.Config != nil {
		alertGroupingParameters.Config.Timeout = nil
	}
	return alertGroupingParameters
}

func expandAutoPauseNotificationsParameters(v interface{}) *pagerduty.AutoPauseNotificationsParameters {
	autoPauseNotificationsParameters := &pagerduty.AutoPauseNotificationsParameters{}
	riur := make(map[string]interface{})

	data, ok := v.([]interface{})
	if ok && len(data) > 0 && !isNilFunc(data[0]) {
		riur = data[0].(map[string]interface{})
	} else {
		return autoPauseNotificationsParameters
	}

	autoPauseNotificationsParameters.Enabled = riur["enabled"].(bool)
	if autoPauseNotificationsParameters.Enabled {
		timeout := riur["timeout"].(int)
		autoPauseNotificationsParameters.Timeout = &timeout
	}
	return autoPauseNotificationsParameters
}

func expandAlertGroupingConfig(groupingType string, v interface{}) *pagerduty.AlertGroupingConfig {
	alertGroupingConfig := &pagerduty.AlertGroupingConfig{}
	rconfig := v.([]interface{})
	if len(rconfig) == 0 || rconfig[0] == nil {
		return nil
	}
	config := rconfig[0].(map[string]interface{})

	if groupingType == "time" {
		if val, ok := config["timeout"]; ok {
			to := val.(int)
			alertGroupingConfig.Timeout = &to
		}
	}

	if groupingType == "intelligent" || groupingType == "content_based" {
		if val, ok := config["time_window"]; ok {
			to := val.(int)
			alertGroupingConfig.TimeWindow = &to
		}
	}

	if groupingType == "content_based" {
		alertGroupingConfig.Fields = []string{}
		if val, ok := config["fields"]; ok {
			for _, field := range val.([]interface{}) {
				alertGroupingConfig.Fields = append(alertGroupingConfig.Fields, field.(string))
			}
		}
		if val, ok := config["aggregate"]; ok {
			agg := val.(string)
			alertGroupingConfig.Aggregate = &agg
		}
	}

	return alertGroupingConfig
}

func flattenAlertGroupingParameters(v *pagerduty.AlertGroupingParameters) interface{} {
	alertGroupingParameters := map[string]interface{}{}

	if v.Config == nil && v.Type == nil {
		return []interface{}{alertGroupingParameters}
	} else {
		alertGroupingParameters = map[string]interface{}{"type": "", "config": []map[string]interface{}{{"aggregate": nil, "fields": nil, "timeout": nil, "time_window": nil}}}
	}

	if v.Type != nil {
		alertGroupingParameters["type"] = v.Type
	}

	if v.Config != nil {
		alertGroupingParameters["config"] = flattenAlertGroupingConfig(v.Config)
	}

	return []interface{}{alertGroupingParameters}
}

func flattenAlertGroupingConfig(v *pagerduty.AlertGroupingConfig) interface{} {
	alertGroupingConfig := map[string]interface{}{
		"aggregate":   v.Aggregate,
		"fields":      v.Fields,
		"timeout":     v.Timeout,
		"time_window": v.TimeWindow,
	}

	return []interface{}{alertGroupingConfig}
}

func flattenAutoPauseNotificationsParameters(v *pagerduty.AutoPauseNotificationsParameters) []interface{} {
	autoPauseNotificationsParameters := map[string]interface{}{
		"enabled": v.Enabled,
	}
	if v.Enabled {
		autoPauseNotificationsParameters["timeout"] = v.Timeout
	}
	if !v.Enabled && v.Timeout == nil {
		autoPauseNotificationsParameters["timeout"] = 120
	}

	return []interface{}{autoPauseNotificationsParameters}
}

func expandIncidentUrgencyRule(v interface{}) *pagerduty.IncidentUrgencyRule {
	incidentUrgencyRule := &pagerduty.IncidentUrgencyRule{}
	riur := make(map[string]interface{})

	data, ok := v.([]interface{})
	if ok && len(data) > 0 && !isNilFunc(data[0]) {
		riur = data[0].(map[string]interface{})
	} else {
		return incidentUrgencyRule
	}

	incidentUrgencyRule.Type = riur["type"].(string)

	if val, ok := riur["urgency"]; ok {
		incidentUrgencyRule.Urgency = val.(string)
	}

	if val, ok := riur["during_support_hours"]; ok {
		if len(val.([]interface{})) > 0 {
			incidentUrgencyRule.DuringSupportHours = expandIncidentUrgencyType(val)
		}
	}

	if val, ok := riur["outside_support_hours"]; ok {
		if len(val.([]interface{})) > 0 {
			incidentUrgencyRule.OutsideSupportHours = expandIncidentUrgencyType(val)
		}
	}

	return incidentUrgencyRule
}

func flattenIncidentUrgencyRule(v *pagerduty.IncidentUrgencyRule) []interface{} {
	incidentUrgencyRule := map[string]interface{}{
		"type":    v.Type,
		"urgency": v.Urgency,
	}

	if v.DuringSupportHours != nil {
		incidentUrgencyRule["during_support_hours"] = flattenIncidentUrgencyType(v.DuringSupportHours)
	}

	if v.OutsideSupportHours != nil {
		incidentUrgencyRule["outside_support_hours"] = flattenIncidentUrgencyType(v.OutsideSupportHours)
	}

	return []interface{}{incidentUrgencyRule}
}

func expandIncidentUrgencyType(v interface{}) *pagerduty.IncidentUrgencyType {
	incidentUrgencyType := &pagerduty.IncidentUrgencyType{}
	riut := make(map[string]interface{})

	data, ok := v.([]interface{})
	if ok && len(data) > 0 && !isNilFunc(data[0]) {
		riut = data[0].(map[string]interface{})
	} else {
		return incidentUrgencyType
	}

	if v, ok := riut["type"]; ok {
		incidentUrgencyType.Type = v.(string)
	}

	if v, ok := riut["urgency"]; ok {
		incidentUrgencyType.Urgency = v.(string)
	}

	return incidentUrgencyType
}

func flattenIncidentUrgencyType(v *pagerduty.IncidentUrgencyType) []interface{} {
	incidentUrgencyType := map[string]interface{}{
		"type":    v.Type,
		"urgency": v.Urgency,
	}
	return []interface{}{incidentUrgencyType}
}

func expandSupportHours(v interface{}) *pagerduty.SupportHours {
	supportHours := &pagerduty.SupportHours{}

	rsh := make(map[string]interface{})

	data, ok := v.([]interface{})
	if ok && len(data) > 0 && !isNilFunc(data[0]) {
		rsh = data[0].(map[string]interface{})
	} else {
		return supportHours
	}

	if v, ok := rsh["type"]; ok {
		supportHours.Type = v.(string)
	}

	if v, ok := rsh["time_zone"]; ok {
		supportHours.TimeZone = v.(string)
	}

	if v, ok := rsh["start_time"]; ok {
		supportHours.StartTime = v.(string)
	}

	if v, ok := rsh["end_time"]; ok {
		supportHours.EndTime = v.(string)
	}

	if v, ok := rsh["days_of_week"]; ok {
		var daysOfWeek []int

		for _, dof := range v.([]interface{}) {
			daysOfWeek = append(daysOfWeek, dof.(int))
		}

		supportHours.DaysOfWeek = daysOfWeek
	}

	return supportHours
}

func flattenSupportHours(v *pagerduty.SupportHours) []interface{} {
	supportHours := map[string]interface{}{}

	if v.Type != "" {
		supportHours["type"] = v.Type
	}

	if v.TimeZone != "" {
		supportHours["time_zone"] = v.TimeZone
	}

	if v.StartTime != "" {
		supportHours["start_time"] = v.StartTime
	}

	if v.EndTime != "" {
		supportHours["end_time"] = v.EndTime
	}

	if len(v.DaysOfWeek) > 0 {
		supportHours["days_of_week"] = v.DaysOfWeek
	}

	return []interface{}{supportHours}
}

func expandScheduledActions(v interface{}) []*pagerduty.ScheduledAction {
	var scheduledActions []*pagerduty.ScheduledAction

	for _, sa := range v.([]interface{}) {
		rsa := sa.(map[string]interface{})

		scheduledAction := &pagerduty.ScheduledAction{
			Type:      rsa["type"].(string),
			ToUrgency: rsa["to_urgency"].(string),
			At:        expandScheduledActionAt(rsa["at"]),
		}

		scheduledActions = append(scheduledActions, scheduledAction)
	}

	return scheduledActions
}

func flattenScheduledActions(v []*pagerduty.ScheduledAction) []interface{} {
	var scheduledActions []interface{}

	for _, sa := range v {
		scheduledAction := map[string]interface{}{
			"type":       sa.Type,
			"to_urgency": sa.ToUrgency,
			"at":         flattenScheduledActionAt(sa.At),
		}
		scheduledActions = append(scheduledActions, scheduledAction)
	}

	return scheduledActions
}

func expandScheduledActionAt(v interface{}) *pagerduty.At {
	rat := make(map[string]interface{})

	data, ok := v.([]interface{})
	if ok && len(data) > 0 && !isNilFunc(data[0]) {
		rat = data[0].(map[string]interface{})
	} else {
		return nil
	}

	return &pagerduty.At{
		Type: rat["type"].(string),
		Name: rat["name"].(string),
	}
}

func flattenScheduledActionAt(v *pagerduty.At) []interface{} {
	at := map[string]interface{}{"type": v.Type, "name": v.Name}
	return []interface{}{at}
}
