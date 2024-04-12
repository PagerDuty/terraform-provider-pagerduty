package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyTeamMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyTeamMembershipCreate,
		Read:   resourcePagerDutyTeamMembershipRead,
		Update: resourcePagerDutyTeamMembershipUpdate,
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
				Default:  "manager",
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"observer",
					"responder",
					"manager",
				}),
			},
		},
	}
}

func maxRetries() int {
	return 4
}

func retryDelayMs() int {
	return 500
}

func calculateDelay(retryCount int) time.Duration {
	return time.Duration(retryCount*retryDelayMs()) * time.Millisecond
}

func fetchPagerDutyTeamMembershipWithRetries(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error, retryCount int, neededRole string) error {
	if retryCount >= maxRetries() {
		return nil
	}
	if err := fetchPagerDutyTeamMembership(d, meta, errCallback); err != nil {
		return err
	}
	fetchedRole, userId, teamId := d.Get("role").(string), d.Get("user_id"), d.Get("team_id")
	if strings.Compare(neededRole, fetchedRole) == 0 {
		return nil
	}
	log.Printf("[DEBUG] Warning role '%s' fetched from PD is different from the role '%s' from config for user: %s from team: %s, retrying...", fetchedRole, neededRole, userId, teamId)

	retryCount++
	time.Sleep(calculateDelay(retryCount))
	return fetchPagerDutyTeamMembershipWithRetries(d, meta, errCallback, retryCount, neededRole)
}

func fetchPagerDutyTeamMembership(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID, teamID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading user: %s from team: %s", userID, teamID)
	return retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.Teams.GetMembers(teamID, &pagerduty.GetMembersOptions{})
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

		for _, member := range resp.Members {
			if member.User.ID == userID {
				d.Set("user_id", userID)
				d.Set("team_id", teamID)
				d.Set("role", member.Role)

				return nil
			}
		}

		log.Printf("[WARN] Removing %s since the user: %s is not a member of: %s", d.Id(), userID, teamID)
		d.SetId("")

		return nil
	})
}

func resourcePagerDutyTeamMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)
	role := d.Get("role").(string)

	log.Printf("[DEBUG] Adding user: %s to team: %s with role: %s", userID, teamID, role)

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.Teams.AddUserWithRole(teamID, userID, role); err != nil {
			if isErrCode(err, 500) {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, teamID))

	return fetchPagerDutyTeamMembershipWithRetries(d, meta, genError, 0, d.Get("role").(string))
}

func resourcePagerDutyTeamMembershipRead(d *schema.ResourceData, meta interface{}) error {
	return fetchPagerDutyTeamMembership(d, meta, handleNotFoundError)
}

func resourcePagerDutyTeamMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)
	role := d.Get("role").(string)

	log.Printf("[DEBUG] Updating user: %s to team: %s with role: %s", userID, teamID, role)

	// To update existing membership resource, We can use the same API as creating a new membership.
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.Teams.AddUserWithRole(teamID, userID, role); err != nil {
			if isErrCode(err, 500) {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, teamID))

	return fetchPagerDutyTeamMembershipWithRetries(d, meta, genError, 0, d.Get("role").(string))
}

func resourcePagerDutyTeamMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID, teamID, err := resourcePagerDutyParseColonCompoundID(d.Id())
	if err != nil {
		return err
	}

	var isFoundErrRemovingUserFromTeam bool
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.Teams.RemoveUser(teamID, userID); err != nil {
			if isErrCode(err, 400) && strings.Contains(err.Error(), "User cannot be removed as they belong to an escalation policy on this team") {
				if !isFoundErrRemovingUserFromTeam {
					// Giving some time for the escalation policies to be removed during destroy operations.
					time.Sleep(2 * time.Second)
					isFoundErrRemovingUserFromTeam = true
					return retry.RetryableError(err)
				}

				return retry.NonRetryableError(err)
			}
			if isErrCode(err, 400) {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil && isFoundErrRemovingUserFromTeam {
		// Extract Escalation Policies associated to the team for which the userID is
		// a rule target.
		epsAssociatedToUser, err := extractEPsAssociatedToTeamAndUser(client, teamID, userID)
		if err != nil {
			return fmt.Errorf("%v; %w", retryErr, err)
		}

		if len(epsAssociatedToUser) > 0 {
			pdURLData, err := url.Parse(epsAssociatedToUser[0])
			if err != nil {
				return fmt.Errorf("%v; %w", retryErr, err)
			}

			accountSubdomain := strings.Split(pdURLData.Hostname(), ".")[0]
			var formatEPsList = func(eps []string) string {
				var formated []string
				for _, ep := range eps {
					formated = append(formated, fmt.Sprintf("\t* %s", ep))
				}
				return strings.Join(formated, "\n")
			}

			return fmt.Errorf(`User %[1]q can't be removed from Team %[2]q as they belong to an Escalation Policy on this team.
Please take only one of the following remediation measures in order to unblock the Team Membership removal:
  1. Remove the user from the following Escalation Policies:
%[4]s
  2. Remove the Escalation Policies from the Team https://%[3]s.pagerduty.com/teams/%[2]s

After completing one of the above given remediation options come back to continue with the destruction of Team Membership.`,
				userID,
				teamID,
				accountSubdomain,
				formatEPsList(epsAssociatedToUser),
			)
		}
	}
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	d.SetId("")

	return nil
}

func buildEPsIdsList(l []*pagerduty.EscalationPolicy) []string {
	eps := []string{}
	for _, o := range l {
		eps = append(eps, o.HTMLURL)
	}
	return unique(eps)
}

// extractEPsAssociatedToTeamAndUser returns the IDs of escalation policies
// associated to the specified team and for which the specified user is a rule
// target.
func extractEPsAssociatedToTeamAndUser(c *pagerduty.Client, teamID, userID string) ([]string, error) {
	var eps []*pagerduty.EscalationPolicy
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := c.EscalationPolicies.List(&pagerduty.ListEscalationPoliciesOptions{TeamIDs: []string{teamID}})
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
		}
		eps = resp.EscalationPolicies
		return nil
	})
	if retryErr != nil {
		return nil, retryErr
	}

	// filter all team escalation policies to only those for which the specified
	// user is a target.
	userEPs := []*pagerduty.EscalationPolicy{}
	for _, ep := range eps {
		for _, rule := range ep.EscalationRules {
			for _, target := range rule.Targets {
				if target.ID == userID {
					userEPs = append(userEPs, ep)
				}
			}
		}
	}

	epsAssociatedToTeamAndUser := buildEPsIdsList(userEPs)
	return epsAssociatedToTeamAndUser, nil
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
