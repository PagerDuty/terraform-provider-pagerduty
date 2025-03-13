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

func dataSourcePagerDutyTeams() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyTeamsRead,

		Schema: map[string]*schema.Schema{
			"query": {
				Description: "Filters the result, showing only the records whose name matches the query.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"teams": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of teams whose name matches the query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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

func dataSourcePagerDutyTeamsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty teams")

	query := d.Get("query").(string)

	o := &pagerduty.ListTeamsOptions{}
	if query != "" {
		o.Query = query
	}

	var pdTeams = make([]*pagerduty.Team, 0, 25)
	more := true
	offset := 0

	for more {
		log.Printf("[DEBUG] Getting PagerDuty teams at offset %d", offset)
		retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, _, err := client.Teams.List(o)
			if err != nil {
				if isErrCode(err, http.StatusBadRequest) {
					return retry.NonRetryableError(err)
				}

				// Delaying retry by 30s as recommended by PagerDuty
				// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
				time.Sleep(30 * time.Second)
				return retry.RetryableError(err)
			}

			pdTeams = append(pdTeams, resp.Teams...)
			more = resp.More
			offset += resp.Limit
			o.Offset = offset

			return nil
		})
	}

	var teams []map[string]interface{}
	for _, team := range pdTeams {
		teams = append(teams, map[string]interface{}{
			"id":          team.ID,
			"name":        team.Name,
			"summary":     team.Summary,
			"description": team.Description,
		})
	}

	// Since this data doesn't have a unique ID, this force this data to be
	// refreshed in every Terraform apply
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("teams", teams)

	return nil
}
