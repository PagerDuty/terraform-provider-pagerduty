package pagerduty

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyUsersRead,

		Schema: map[string]*schema.Schema{
			"team_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of users who are members of the team",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePagerDutyUsersRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty users")

	pre := d.Get("team_ids").([]interface{})
	var teamIds []string
	for _, ti := range pre {
		teamIds = append(teamIds, ti.(string))
	}

	o := &pagerduty.ListUsersOptions{
		TeamIDs: teamIds,
	}

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.Users.ListAll(o)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		var users []map[string]interface{}
		for _, user := range resp {
			users = append(users, map[string]interface{}{
				"id":          user.ID,
				"name":        user.Name,
				"email":       user.Email,
				"role":        user.Role,
				"job_title":   user.JobTitle,
				"time_zone":   user.TimeZone,
				"description": user.Description,
			})
		}

		// Since this data doesn't have an unique ID, this force this data to be
		// refreshed in every Terraform apply
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.Set("users", users)

		return nil
	})
}
