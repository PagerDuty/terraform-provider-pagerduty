package pagerduty

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyTeamMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyTeamMembersRead,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team to find via the PagerDuty API",
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The set of team memberships associated with the team",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePagerDutyTeamMembersRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	teamID := d.Get("team_id").(string)

	log.Printf("[INFO] Reading PagerDuty team members of %s", teamID)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Teams.GetMembers(teamID, &pagerduty.GetMembersOptions{})
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(err)
		}

		var mems []map[string]interface{}
		for _, member := range resp.Members {
			mems = append(mems, map[string]interface{}{
				"id":      member.User.ID,
				"type":    member.User.Type,
				"summary": member.User.Summary,
				"role":    member.Role,
			})
		}

		// Since this data doesn't have an unique ID, this forces the data to be
		// refreshed with each Terraform apply.
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

		d.Set("members", mems)
		d.Set("team_id", teamID)

		return nil
	})
}
