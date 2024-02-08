package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyAutomationActionsActionTeamAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAutomationActionsActionTeamAssociationCreate,
		Read:   resourcePagerDutyAutomationActionsActionTeamAssociationRead,
		Delete: resourcePagerDutyAutomationActionsActionTeamAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"action_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePagerDutyAutomationActionsActionTeamAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	actionID := d.Get("action_id").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[INFO] Creating PagerDuty AutomationActionsActionTeamAssociation %s:%s", d.Get("action_id").(string), d.Get("team_id").(string))

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if teamRef, _, err := client.AutomationActionsAction.AssociateToTeam(actionID, teamID); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			if isErrCode(err, 429) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		} else if teamRef != nil {
			d.SetId(fmt.Sprintf("%s:%s", actionID, teamID))
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return fetchPagerDutyAutomationActionsActionTeamAssociation(d, meta, handleNotFoundError)
}

func fetchPagerDutyAutomationActionsActionTeamAssociation(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	actionID, teamID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.AutomationActionsAction.GetAssociationToTeam(actionID, teamID)
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

		if resp.Team.ID != teamID {
			log.Printf("[WARN] Removing %s since the user: %s is not a member of: %s", d.Id(), actionID, teamID)
			d.SetId("")
			return nil
		}

		d.Set("action_id", actionID)
		d.Set("team_id", resp.Team.ID)

		return nil
	})
}

func resourcePagerDutyAutomationActionsActionTeamAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return fetchPagerDutyAutomationActionsActionTeamAssociation(d, meta, handleNotFoundError)
}

func resourcePagerDutyAutomationActionsActionTeamAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	actionID, teamID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty AutomationActionsActionTeamAssociation %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.AutomationActionsAction.DissociateToTeam(actionID, teamID); err != nil {
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
