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
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateValueFunc([]string{
					"observer",
					"responder",
					"manager",
				}),
				ForceNew: true,
			},
		},
	}
}
func resourcePagerDutyTeamMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)
	role := d.Get("role").(string)

	log.Printf("[DEBUG] Adding user: %s to team: %s with role: %s", userID, teamID, role)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.Teams.AddUserWithRole(teamID, userID, role); err != nil {
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

	d.Set("role", "")

	if isTeamMember(user, teamID) {
		resp, _, err := client.Teams.GetMembers(teamID, &pagerduty.GetMembersOptions{})
		if err != nil {
			return err
		}

		for _, member := range resp.Members {
			if member.User.ID == userID {
				d.Set("role", member.Role)
				break
			}
		}
	} else {
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
