package pagerduty

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyAutomationActionsAction() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAutomationActionsActionCreate,
		Read:   resourcePagerDutyAutomationActionsActionRead,
		Delete: resourcePagerDutyAutomationActionsActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"action_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"script",
					"process_automation",
				}),
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"runner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"action_data_reference": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"process_automation_job_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true, // Requires creation of new resource while support for update is not implemented
						},
						"process_automation_job_arguments": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true, // Requires creation of new resource while support for update is not implemented
						},
						"script": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true, // Requires creation of new resource while support for update is not implemented
						},
						"invocation_command": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true, // Requires creation of new resource while support for update is not implemented
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action_classification": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateValueFunc([]string{
					"diagnostic",
					"remediation",
				}),
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"runner_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func buildAutomationActionsActionStruct(d *schema.ResourceData) (*pagerduty.AutomationActionsAction, error) {

	automationActionsAction := pagerduty.AutomationActionsAction{
		Name:       d.Get("name").(string),
		ActionType: d.Get("action_type").(string),
	}

	if automationActionsAction.ActionType != "process_automation" {
		return nil, errors.New("only actions of action_type process_automation can be created")
	}

	// The API does not allow new actions without a description, but legacy actions without a description exist
	if attr, ok := d.GetOk("description"); ok {
		val := attr.(string)
		automationActionsAction.Description = &val
	} else {
		return nil, errors.New("action description must be specified when creating an action")
	}

	if attr, ok := d.GetOk("action_data_reference"); ok {
		automationActionsAction.ActionDataReference = expandActionDataReference(attr)
	} else {
		return nil, errors.New("action_data_reference must be specified when creating an action")
	}

	return &automationActionsAction, nil
}

func expandActionDataReference(v interface{}) pagerduty.AutomationActionsActionDataReference {
	attr_map := v.([]interface{})[0].(map[string]interface{})
	process_automation_job_id := attr_map["process_automation_job_id"].(string)
	adr := pagerduty.AutomationActionsActionDataReference{
		ProcessAutomationJobId: &process_automation_job_id,
	}
	return adr
}

func resourcePagerDutyAutomationActionsActionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	automationActionsAction, err := buildAutomationActionsActionStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating PagerDuty AutomationActionsAction %s", automationActionsAction.Name)

	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
		if automationActionsAction, _, err := client.AutomationActionsAction.Create(automationActionsAction); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if automationActionsAction != nil {
			d.SetId(automationActionsAction.ID)
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return resourcePagerDutyAutomationActionsActionRead(d, meta)
}

func resourcePagerDutyAutomationActionsActionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty AutomationActionsAction %s", d.Id())

	return resource.Retry(30*time.Second, func() *resource.RetryError {
		if automationActionsAction, _, err := client.AutomationActionsAction.Get(d.Id()); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if automationActionsAction != nil {
			d.Set("name", automationActionsAction.Name)
			d.Set("type", automationActionsAction.Type)
			d.Set("action_type", automationActionsAction.ActionType)
			d.Set("creation_time", automationActionsAction.CreationTime)

			if automationActionsAction.Description != nil {
				d.Set("description", &automationActionsAction.Description)
			}

			f_adr := flattenActionDataReference(automationActionsAction.ActionDataReference)
			if err := d.Set("action_data_reference", f_adr); err != nil {
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
}

func flattenActionDataReference(adr pagerduty.AutomationActionsActionDataReference) []interface{} {
	adr_map := map[string]interface{}{}

	if adr.ProcessAutomationJobId != nil {
		adr_map["process_automation_job_id"] = *adr.ProcessAutomationJobId
	}

	if adr.ProcessAutomationJobArguments != nil {
		adr_map["process_automation_job_arguments"] = *adr.ProcessAutomationJobArguments
	}

	if adr.Script != nil {
		adr_map["script"] = *adr.Script
	}

	if adr.InvocationCommand != nil {
		adr_map["invocation_command"] = *adr.InvocationCommand
	}

	return []interface{}{adr_map}
}

func resourcePagerDutyAutomationActionsActionDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty AutomationActionsAction %s", d.Id())

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.AutomationActionsAction.Delete(d.Id()); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	d.SetId("")

	// giving the API time to catchup
	time.Sleep(time.Second)
	return nil
}
