package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyService() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyServiceCreate,
		Read:   resourcePagerDutyServiceRead,
		Update: resourcePagerDutyServiceUpdate,
		Delete: resourcePagerDutyServiceDelete,
		CustomizeDiff: func(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
			in := diff.Get("incident_urgency_rule.#").(int)
			for i := 0; i <= in; i++ {
				t := diff.Get(fmt.Sprintf("incident_urgency_rule.%d.type", i)).(string)
				if t == "use_support_hours" && diff.Get(fmt.Sprintf("incident_urgency_rule.%d.urgency", i)).(string) != "" {
					return fmt.Errorf("general urgency cannot be set for a use_support_hours incident urgency rule type")
				}
			}
			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateCantBeBlankOrNotPrintableChars,
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
				Default:  "create_incidents",
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								_, err := time.LoadLocation(val.(string))
								if err != nil {
									errs = append(errs, err)
								}
								return
							},
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
		service.AlertGrouping = &ag
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

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		service, _, err := client.Services.Get(d.Id(), &pagerduty.GetServiceOptions{
			Includes: []string{"auto_pause_notifications_parameters"},
		})
		if err != nil {
			log.Printf("[WARN] Service read error")
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		if err := flattenService(d, service); err != nil {
			return resource.NonRetryableError(err)
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
		return err
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
		return err
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
	if service.AlertGroupingParameters != nil {
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
	pre := v.([]interface{})[0]
	if isNilFunc(pre) {
		return nil
	}
	riur := pre.(map[string]interface{})
	groupingType := ""
	if riur["type"].(string) != "" {
		groupingType = riur["type"].(string)
		alertGroupingParameters.Type = &groupingType
	}

	// For Intelligent grouping type, config is null
	alertGroupingParameters.Config = nil
	if groupingType == "content_based" || groupingType == "time" {
		alertGroupingParameters.Config = expandAlertGroupingConfig(riur["config"], groupingType)
	}
	return alertGroupingParameters
}

func expandAutoPauseNotificationsParameters(v interface{}) *pagerduty.AutoPauseNotificationsParameters {
	riur := v.([]interface{})[0].(map[string]interface{})
	autoPauseNotificationsParameters := &pagerduty.AutoPauseNotificationsParameters{
		Enabled: riur["enabled"].(bool),
	}
	if autoPauseNotificationsParameters.Enabled {
		timeout := riur["timeout"].(int)
		autoPauseNotificationsParameters.Timeout = &timeout
	}
	return autoPauseNotificationsParameters
}

func expandAlertGroupingConfig(v interface{}, groupingType string) *pagerduty.AlertGroupingConfig {
	alertGroupingConfig := &pagerduty.AlertGroupingConfig{}
	if len(v.([]interface{})) == 0 || v.([]interface{})[0] == nil {
		return nil
	}
	riur := v.([]interface{})[0].(map[string]interface{})

	alertGroupingConfig.Fields = []string{}
	if val, ok := riur["fields"]; ok {
		for _, field := range val.([]interface{}) {
			alertGroupingConfig.Fields = append(alertGroupingConfig.Fields, field.(string))
		}
	}
	if val, ok := riur["aggregate"]; ok {
		agg := val.(string)
		alertGroupingConfig.Aggregate = &agg
	}
	if val, ok := riur["timeout"]; ok && groupingType == "time" {
		to := val.(int)
		alertGroupingConfig.Timeout = &to
	}
	return alertGroupingConfig
}
func flattenAlertGroupingParameters(v *pagerduty.AlertGroupingParameters) interface{} {
	alertGroupingParameters := map[string]interface{}{}

	if v.Config == nil && v.Type == nil {
		return []interface{}{alertGroupingParameters}
	} else {
		alertGroupingParameters = map[string]interface{}{"type": "", "config": []map[string]interface{}{{"aggregate": nil, "fields": nil, "timeout": nil}}}
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
		"aggregate": v.Aggregate,
		"fields":    v.Fields,
		"timeout":   v.Timeout,
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

	return []interface{}{autoPauseNotificationsParameters}
}

func expandIncidentUrgencyRule(v interface{}) *pagerduty.IncidentUrgencyRule {
	riur := v.([]interface{})[0].(map[string]interface{})
	incidentUrgencyRule := &pagerduty.IncidentUrgencyRule{
		Type: riur["type"].(string),
	}

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
	riut := v.([]interface{})[0].(map[string]interface{})
	incidentUrgencyType := &pagerduty.IncidentUrgencyType{}

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

	rsh := v.([]interface{})[0].(map[string]interface{})

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
	rat := v.([]interface{})[0].(map[string]interface{})
	return &pagerduty.At{
		Type: rat["type"].(string),
		Name: rat["name"].(string),
	}
}

func flattenScheduledActionAt(v *pagerduty.At) []interface{} {
	at := map[string]interface{}{"type": v.Type, "name": v.Name}
	return []interface{}{at}
}
