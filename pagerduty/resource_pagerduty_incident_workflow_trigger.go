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

func resourcePagerDutyIncidentWorkflowTrigger() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePagerDutyIncidentWorkflowTriggerRead,
		UpdateContext: resourcePagerDutyIncidentWorkflowTriggerUpdate,
		DeleteContext: resourcePagerDutyIncidentWorkflowTriggerDelete,
		CreateContext: resourcePagerDutyIncidentWorkflowTriggerCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: validateIncidentWorkflowTrigger,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"manual",
					"conditional",
					"incident_type",
				}),
			},
			"workflow": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subscribed_to_all_services": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"incident_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restricted": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourcePagerDutyIncidentWorkflowTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	iwt, err := buildIncidentWorkflowTriggerStruct(d, true)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating PagerDuty incident workflow trigger %s for %s.", iwt.Type, iwt.Workflow.ID)

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		createdWorkflowTrigger, _, err := client.IncidentWorkflowTriggers.CreateContext(ctx, iwt)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		}

		err = flattenIncidentWorkflowTrigger(d, createdWorkflowTrigger)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyIncidentWorkflowTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty incident workflow trigger %s", d.Id())
	err := fetchIncidentWorkflowTrigger(ctx, d, meta, handleNotFoundError)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePagerDutyIncidentWorkflowTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	iwt, err := buildIncidentWorkflowTriggerStruct(d, false)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty incident workflow trigger %s", d.Id())

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		updatedWorkflowTrigger, _, err := client.IncidentWorkflowTriggers.UpdateContext(ctx, d.Id(), iwt)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}

		err = flattenIncidentWorkflowTrigger(d, updatedWorkflowTrigger)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyIncidentWorkflowTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.IncidentWorkflowTriggers.DeleteContext(ctx, d.Id())
	if err != nil {
		if isErrCode(err, http.StatusNotFound) {
			return diag.FromErr(handleNotFoundError(err, d))
		}
		return diag.FromErr(err)
	}
	return nil
}

func validateIncidentWorkflowTrigger(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	triggerType := d.Get("type").(string)
	_, hadCondition := d.GetOk("condition")
	if triggerType != "conditional" && hadCondition {
		return fmt.Errorf("when trigger type %s is used, condition must not be specified", triggerType)
	}

	// pagerduty_incident_workflow_trigger.permissions input validation
	permissionRestricted := d.Get("permissions.0.restricted").(bool)
	permissionTeamID := d.Get("permissions.0.team_id").(string)
	if triggerType != "manual" && permissionRestricted {
		return fmt.Errorf("restricted can only be true when trigger type is manual")
	}
	if !permissionRestricted && permissionTeamID != "" {
		return fmt.Errorf("team_id not allowed when restricted is false")
	}

	s, hadServices := d.GetOk("services")
	all := d.Get("subscribed_to_all_services").(bool)
	if all && hadServices && len(s.([]interface{})) > 0 {
		return fmt.Errorf("when subscribed_to_all_services is true, services must either be not defined or empty")
	}

	// The API rejects is_subscribed_to_all_services outright when trigger_type
	// is incident_type (it isn't a service-scoped trigger type), so require it
	// be false in config rather than silently never sending a true value.
	if triggerType == "incident_type" && all {
		return fmt.Errorf("subscribed_to_all_services must be false when trigger type is incident_type")
	}

	// The API requires incident_types to contain at least one element whenever
	// it is sent (it rejects an explicit empty list with "must contain at least
	// 1 items"), so catch an empty list at plan time instead of a 400 at apply.
	if triggerType == "incident_type" && len(d.Get("incident_types").([]interface{})) == 0 {
		return fmt.Errorf("incident_types must contain at least one item when trigger type is incident_type")
	}

	return nil
}

func fetchIncidentWorkflowTrigger(ctx context.Context, d *schema.ResourceData, meta interface{}, errorCallback func(err error, d *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		iwt, _, err := client.IncidentWorkflowTriggers.GetContext(ctx, d.Id())
		if err != nil {
			log.Printf("[WARN] Incident workflow trigger read error")
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

		if err := flattenIncidentWorkflowTrigger(d, iwt); err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
}

func flattenIncidentWorkflowTrigger(d *schema.ResourceData, t *pagerduty.IncidentWorkflowTrigger) error {
	d.SetId(t.ID)
	d.Set("type", t.TriggerType.String())
	if t.Workflow != nil {
		d.Set("workflow", t.Workflow.ID)
	}
	d.Set("services", flattenIncidentWorkflowEnabledServices(t.Services))
	if t.TriggerType == pagerduty.IncidentWorkflowTriggerTypeIncidentType {
		// The API returns is_subscribed_to_all_services as true for incident_type
		// triggers even though it's never sent on create/update (validated to
		// always be false for this type); trusting the response here would
		// produce a permanent plan diff against the only value config is ever
		// allowed to hold for this type.
		d.Set("subscribed_to_all_services", false)
	} else {
		d.Set("subscribed_to_all_services", t.SubscribedToAllServices)
	}
	if t.Condition != nil {
		d.Set("condition", t.Condition)
	}
	if t.Permissions != nil {
		d.Set("permissions", []map[string]interface{}{
			{
				"restricted": t.Permissions.Restricted,
				"team_id":    t.Permissions.TeamID,
			},
		})
	}
	if len(t.IncidentTypes) > 0 {
		d.Set("incident_types", t.IncidentTypes)
	}

	return nil
}

func flattenIncidentWorkflowEnabledServices(s []*pagerduty.ServiceReference) []string {
	services := make([]string, len(s))
	for i, v := range s {
		services[i] = v.ID
	}
	return services
}

func buildIncidentWorkflowTriggerStruct(d *schema.ResourceData, forUpdate bool) (*pagerduty.IncidentWorkflowTrigger, error) {
	triggerType := d.Get("type").(string)

	// The API rejects is_subscribed_to_all_services outright for incident_type
	// triggers. validateIncidentWorkflowTrigger requires it be false in config
	// for that type, and the vendored client's omitempty tag on this plain bool
	// field already omits a false value from the request, so this needs no
	// special-casing here.
	iwt := pagerduty.IncidentWorkflowTrigger{
		SubscribedToAllServices: d.Get("subscribed_to_all_services").(bool),
	}

	if forUpdate {
		iwt.Workflow = &pagerduty.IncidentWorkflow{
			ID: d.Get("workflow").(string),
		}
		iwt.TriggerType = pagerduty.IncidentWorkflowTriggerTypeFromString(triggerType)
	}

	if services, ok := d.GetOk("services"); ok {
		iwt.Services = buildIncidentWorkflowTriggerServices(services)
	}

	// IncidentTypes must be set unconditionally, outside the forUpdate block, so
	// updates to the list are sent to the API (unlike Workflow/TriggerType, which
	// are immutable and only sent on Create). Use d.Get rather than d.GetOk: GetOk
	// treats an empty list as "no value", which would silently drop a shrink of
	// the list instead of sending it to the API; validateIncidentWorkflowTrigger
	// rejects an empty list for incident_type at plan time (the API requires at
	// least one element), so build never needs to send an empty one. Gate on
	// triggerType so manual/conditional triggers never send this field at all.
	if triggerType == "incident_type" {
		iwt.IncidentTypes = buildIncidentWorkflowTriggerIncidentTypes(d.Get("incident_types"))
	}

	// Special handling for condition to support empty string conditions
	// GetOk won't return true for empty strings, but we need to set them
	// for conditional triggers that execute on incident creation
	if condition, ok := d.GetOk("condition"); triggerType == "conditional" || ok {
		condStr := condition.(string)
		iwt.Condition = &condStr
	}

	// The API rejects the permissions key outright for incident_type triggers.
	// permissions is Optional+Computed, so once a computed default is written
	// to state by a prior apply, d.GetOk("permissions") sees it as present on
	// every subsequent update even for a trigger whose config never mentions
	// it; gate on triggerType rather than trusting GetOk alone.
	if permissions, ok := d.GetOk("permissions"); ok && triggerType != "incident_type" {
		p, err := expandIncidentWorkflowTriggerPermissions(permissions)
		if err != nil {
			return nil, err
		}
		iwt.Permissions = p
	}

	return &iwt, nil
}

func buildIncidentWorkflowTriggerServices(s interface{}) []*pagerduty.ServiceReference {
	services := s.([]interface{})
	newServices := make([]*pagerduty.ServiceReference, len(services))
	for i, v := range services {
		newServices[i] = &pagerduty.ServiceReference{
			ID: v.(string),
		}
	}
	return newServices
}

func buildIncidentWorkflowTriggerIncidentTypes(v interface{}) []string {
	incidentTypes := v.([]interface{})
	newIncidentTypes := make([]string, len(incidentTypes))
	for i, v := range incidentTypes {
		newIncidentTypes[i] = v.(string)
	}
	return newIncidentTypes
}

func expandIncidentWorkflowTriggerPermissions(v interface{}) (*pagerduty.IncidentWorkflowTriggerPermissions, error) {
	var permissions *pagerduty.IncidentWorkflowTriggerPermissions

	permissionsData, ok := v.([]interface{})
	if ok && len(permissionsData) > 0 {
		p := permissionsData[0].(map[string]interface{})

		// Unfortunately this validatation can't be made during diff checking, since
		// Diff Customization doesn't support computed/"known after apply" values
		// like team_id in this case. Based on
		// https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/customizing-differences
		// because of this, it will only be returned during the apply phase.
		if p["restricted"].(bool) && p["team_id"].(string) == "" {
			return nil, fmt.Errorf("team_id must be specified when restricted is true")
		}

		permissions = &pagerduty.IncidentWorkflowTriggerPermissions{
			Restricted: p["restricted"].(bool),
			TeamID:     p["team_id"].(string),
		}
	}

	return permissions, nil
}
