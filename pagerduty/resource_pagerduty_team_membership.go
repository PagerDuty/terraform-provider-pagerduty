package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyTeamMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyTeamMembershipCreate,
		Read:   resourcePagerDutyTeamMembershipRead,
		Delete: resourcePagerDutyTeamMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
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
func resourcePagerDutyTeamMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[DEBUG] Adding user: %s to team: %s", userID, teamID)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.Teams.AddUser(teamID, userID); err != nil {
			if isErrCode(err, 500) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, teamID))

	return resourcePagerDutyTeamMembershipRead(d, meta)
}

func resourcePagerDutyTeamMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID, teamID := resourcePagerDutyTeamMembershipParseID(d.Id())

	log.Printf("[DEBUG] Reading user: %s from team: %s", userID, teamID)

	user, _, err := client.Users.Get(userID, &pagerduty.GetUserOptions{})
	if err != nil {
		return err
	}

	if !isTeamMember(user, teamID) {
		log.Printf("[WARN] Removing %s since the user: %s is not a member of: %s", d.Id(), userID, teamID)
		d.SetId("")
	}

	d.Set("user_id", userID)
	d.Set("team_id", teamID)

	return nil
}

func resourcePagerDutyTeamMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID, teamID := resourcePagerDutyTeamMembershipParseID(d.Id())

	log.Printf("[DEBUG] Removing user: %s from team: %s", userID, teamID)

	if _, err := client.Teams.RemoveUser(teamID, userID); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourcePagerDutyTeamMembershipParseID(id string) (string, string) {
	parts := strings.Split(id, ":")
	return parts[0], parts[1]
}

func isTeamMember(user *pagerduty.User, teamID string) bool {
	var found bool

	for _, team := range user.Teams {
		if teamID == team.ID {
			found = true
		}
	}

	return found
}
