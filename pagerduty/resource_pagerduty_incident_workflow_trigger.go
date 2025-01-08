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
	if triggerType == "manual" && hadCondition {
		return fmt.Errorf("when trigger type manual is used, condition must not be specified")
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
	d.Set("subscribed_to_all_services", t.SubscribedToAllServices)
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
	iwt := pagerduty.IncidentWorkflowTrigger{
		SubscribedToAllServices: d.Get("subscribed_to_all_services").(bool),
	}

	if forUpdate {
		iwt.Workflow = &pagerduty.IncidentWorkflow{
			ID: d.Get("workflow").(string),
		}
		iwt.TriggerType = pagerduty.IncidentWorkflowTriggerTypeFromString(d.Get("type").(string))
	}

	if services, ok := d.GetOk("services"); ok {
		iwt.Services = buildIncidentWorkflowTriggerServices(services)
	}

	if condition, ok := d.GetOk("condition"); ok {
		condStr := condition.(string)
		iwt.Condition = &condStr
	}

	if permissions, ok := d.GetOk("permissions"); ok {
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
