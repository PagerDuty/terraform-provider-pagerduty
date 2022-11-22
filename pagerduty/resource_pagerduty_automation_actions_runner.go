package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyAutomationActionsRunner() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAutomationActionsRunnerCreate,
		Read:   resourcePagerDutyAutomationActionsRunnerRead,
		Delete: resourcePagerDutyAutomationActionsRunnerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"runner_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"sidecar",
				}),
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
				ForceNew: true, // Requires creation of new resource while support for update is not implemented
			},
		},
	}
}

func buildAutomationActionsRunnerStruct(d *schema.ResourceData) *pagerduty.AutomationActionsRunner {
	automationActionsRunner := &pagerduty.AutomationActionsRunner{
		ID:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		RunnerType:  d.Get("runner_type").(string),
		Description: d.Get("description").(string),
	}

	return automationActionsRunner
}

func resourcePagerDutyAutomationActionsRunnerCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	automationActionsRunner := buildAutomationActionsRunnerStruct(d)

	log.Printf("[INFO] Creating PagerDuty AutomationActionsRunner %s", automationActionsRunner.Name)

	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
		if automationActionsRunner, _, err := client.AutomationActionsRunner.Create(automationActionsRunner); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
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

	return resource.Retry(30*time.Second, func() *resource.RetryError {
		if automationActionsRunner, _, err := client.AutomationActionsRunner.Get(d.Id()); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if automationActionsRunner != nil {
			log.Printf("AutomationActionsRunner Name: %v", automationActionsRunner.Name)
			log.Printf("AutomationActionsRunner RunnerType: %v", automationActionsRunner.RunnerType)
			d.Set("id", automationActionsRunner.ID)
			d.Set("name", automationActionsRunner.Name)
			d.Set("runner_type", automationActionsRunner.RunnerType)
			d.Set("description", automationActionsRunner.Description)
		}
		return nil
	})
}

func resourcePagerDutyAutomationActionsRunnerDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty AutomationActionsRunner %s", d.Id())

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.AutomationActionsRunner.Delete(d.Id()); err != nil {
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
