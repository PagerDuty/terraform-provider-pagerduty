package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyAutomationActionsActionServiceAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAutomationActionsActionServiceAssociationCreate,
		Read:   resourcePagerDutyAutomationActionsActionServiceAssociationRead,
		Delete: resourcePagerDutyAutomationActionsActionServiceAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"action_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePagerDutyAutomationActionsActionServiceAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	actionID := d.Get("action_id").(string)
	serviceID := d.Get("service_id").(string)

	log.Printf("[INFO] Creating PagerDuty AutomationActionsActionServiceAssociation %s:%s", d.Get("action_id").(string), d.Get("service_id").(string))

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if serviceRef, _, err := client.AutomationActionsAction.AssociateToService(actionID, serviceID); err != nil {
			if isErrCode(err, 429) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		} else if serviceRef != nil {
			d.SetId(fmt.Sprintf("%s:%s", actionID, serviceID))
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return fetchPagerDutyAutomationActionsActionServiceAssociation(d, meta, handleNotFoundError)
}

func fetchPagerDutyAutomationActionsActionServiceAssociation(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	actionID, serviceID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.AutomationActionsAction.GetAssociationToService(actionID, serviceID)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}

			return nil
		}

		if resp.Service.ID != serviceID {
			log.Printf("[WARN] Removing %s since the service: %s is not associated to the action: %s", d.Id(), actionID, serviceID)
			d.SetId("")
			return nil
		}

		d.Set("action_id", actionID)
		d.Set("service_id", resp.Service.ID)

		return nil
	})
}

func resourcePagerDutyAutomationActionsActionServiceAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return fetchPagerDutyAutomationActionsActionServiceAssociation(d, meta, handleNotFoundError)
}

func resourcePagerDutyAutomationActionsActionServiceAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	actionID, serviceID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty AutomationActionsActionServiceAssociation %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.AutomationActionsAction.DissociateFromService(actionID, serviceID); err != nil {
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
