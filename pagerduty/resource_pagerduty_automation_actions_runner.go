package pagerduty

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyAutomationActionsRunner() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAutomationActionsRunnerCreate,
		Read:   resourcePagerDutyAutomationActionsRunnerRead,
		Update: resourcePagerDutyAutomationActionsRunnerUpdate,
		Delete: resourcePagerDutyAutomationActionsRunnerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"runner_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"sidecar",
					"runbook",
				}),
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runbook_base_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runbook_api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_seen": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func buildAutomationActionsRunnerStruct(d *schema.ResourceData) (*pagerduty.AutomationActionsRunner, error) {
	automationActionsRunner := pagerduty.AutomationActionsRunner{
		Name:       d.Get("name").(string),
		RunnerType: d.Get("runner_type").(string),
	}

	if automationActionsRunner.RunnerType != "runbook" {
		return nil, errors.New("only runners of runner_type runbook can be created")
	}

	// The API does not allow new runners without a description, but legacy runners without a description exist
	if attr, ok := d.GetOk("description"); ok {
		val := attr.(string)
		automationActionsRunner.Description = &val
	} else {
		return nil, errors.New("runner description must be specified when creating a runbook runner")
	}

	if attr, ok := d.GetOk("runbook_base_uri"); ok {
		val := attr.(string)
		automationActionsRunner.RunbookBaseUri = &val
	} else {
		return nil, errors.New("runbook_base_uri must be specified when creating a runbook runner")
	}

	if attr, ok := d.GetOk("runbook_api_key"); ok {
		val := attr.(string)
		automationActionsRunner.RunbookApiKey = &val
	} else {
		return nil, errors.New("runbook_api_key must be specified when creating a runbook runner")
	}

	return &automationActionsRunner, nil
}

func resourcePagerDutyAutomationActionsRunnerCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	automationActionsRunner, err := buildAutomationActionsRunnerStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating PagerDuty AutomationActionsRunner %s", automationActionsRunner.Name)

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if automationActionsRunner, _, err := client.AutomationActionsRunner.Create(automationActionsRunner); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		} else if automationActionsRunner != nil {
			d.SetId(automationActionsRunner.ID)
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return resourcePagerDutyAutomationActionsRunnerRead(d, meta)
}

func resourcePagerDutyAutomationActionsRunnerRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty AutomationActionsRunner %s", d.Id())

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		if automationActionsRunner, _, err := client.AutomationActionsRunner.Get(d.Id()); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
		} else if automationActionsRunner != nil {
			d.Set("name", automationActionsRunner.Name)
			d.Set("type", automationActionsRunner.Type)
			d.Set("runner_type", automationActionsRunner.RunnerType)
			d.Set("creation_time", automationActionsRunner.CreationTime)

			if automationActionsRunner.Description != nil {
				d.Set("description", &automationActionsRunner.Description)
			}

			if automationActionsRunner.RunbookBaseUri != nil {
				d.Set("runbook_base_uri", &automationActionsRunner.RunbookBaseUri)
			}

			if automationActionsRunner.LastSeenTime != nil {
				d.Set("last_seen", &automationActionsRunner.LastSeenTime)
			}
		}
		return nil
	})
}

func resourcePagerDutyAutomationActionsRunnerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	automationActionsRunner, err := buildAutomationActionsRunnerStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating PagerDuty AutomationActionsRunner %s", d.Id())

	if _, _, err := client.AutomationActionsRunner.Update(d.Id(), automationActionsRunner); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyAutomationActionsRunnerDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty AutomationActionsRunner %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.AutomationActionsRunner.Delete(d.Id()); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
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
