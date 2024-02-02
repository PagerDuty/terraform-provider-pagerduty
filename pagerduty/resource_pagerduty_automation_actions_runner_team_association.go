package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePagerDutyAutomationActionsRunnerTeamAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAutomationActionsRunnerTeamAssociationCreate,
		Read:   resourcePagerDutyAutomationActionsRunnerTeamAssociationRead,
		Delete: resourcePagerDutyAutomationActionsRunnerTeamAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"runner_id": {
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

func resourcePagerDutyAutomationActionsRunnerTeamAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	runnerID := d.Get("runner_id").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[INFO] Creating PagerDuty AutomationActionsRunnerTeamAssociation %s:%s", d.Get("runner_id").(string), d.Get("team_id").(string))

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if teamRef, _, err := client.AutomationActionsRunner.AssociateToTeam(runnerID, teamID); err != nil {
			if isErrCode(err, 429) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		} else if teamRef != nil {
			d.SetId(fmt.Sprintf("%s:%s", runnerID, teamID))
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return fetchPagerDutyAutomationActionsRunnerTeamAssociation(d, meta, handleNotFoundError)
}

func fetchPagerDutyAutomationActionsRunnerTeamAssociation(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	runnerID, teamID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.AutomationActionsRunner.GetAssociationToTeam(runnerID, teamID)
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
			log.Printf("[WARN] Removing %s since the user: %s is not a member of: %s", d.Id(), runnerID, teamID)
			d.SetId("")
			return nil
		}

		d.Set("runner_id", runnerID)
		d.Set("team_id", resp.Team.ID)

		return nil
	})
}

func resourcePagerDutyAutomationActionsRunnerTeamAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return fetchPagerDutyAutomationActionsRunnerTeamAssociation(d, meta, handleNotFoundError)
}

func resourcePagerDutyAutomationActionsRunnerTeamAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	runnerID, teamID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty AutomationActionsRunnerTeamAssociation %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.AutomationActionsRunner.DissociateFromTeam(runnerID, teamID); err != nil {
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
