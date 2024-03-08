package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// Deprecated: Migrated to pagerdutyplugin.resourceTeam. Kept for testing purposes.
func resourcePagerDutyTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyTeamCreate,
		Read:   resourcePagerDutyTeamRead,
		Update: resourcePagerDutyTeamUpdate,
		Delete: resourcePagerDutyTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildTeamStruct(d *schema.ResourceData) *pagerduty.Team {
	team := &pagerduty.Team{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("description"); ok {
		team.Description = attr.(string)
	}
	if attr, ok := d.GetOk("parent"); ok {
		team.Parent = &pagerduty.TeamReference{
			ID:   attr.(string),
			Type: "team_reference",
		}
	}
	if attr, ok := d.GetOk("default_role"); ok {
		team.DefaultRole = attr.(string)
	}
	return team
}

func resourcePagerDutyTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	team := buildTeamStruct(d)

	log.Printf("[INFO] Creating PagerDuty team %s", team.Name)

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if team, _, err := client.Teams.Create(team); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			return retry.RetryableError(err)
		} else if team != nil {
			d.SetId(team.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return resourcePagerDutyTeamRead(d, meta)
}

func resourcePagerDutyTeamRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty team %s", d.Id())

	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		if team, _, err := client.Teams.Get(d.Id()); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := handleNotFoundError(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}
		} else if team != nil {
			d.Set("name", team.Name)
			d.Set("description", team.Description)
			d.Set("html_url", team.HTMLURL)
			d.Set("default_role", team.DefaultRole)
		}
		return nil
	})
}

func resourcePagerDutyTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	team := buildTeamStruct(d)

	log.Printf("[INFO] Updating PagerDuty team %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, _, err := client.Teams.Update(d.Id(), team); err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return resourcePagerDutyTeamRead(d, meta)
}

func resourcePagerDutyTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty team %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.Teams.Delete(d.Id()); err != nil {
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
