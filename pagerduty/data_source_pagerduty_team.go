package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team to find in the PagerDuty API",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty team")

	searchTeam := d.Get("name").(string)

	o := &pagerduty.ListTeamsOptions{
		Query: searchTeam,
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Teams.List(o)
		if err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.Team

		for _, team := range resp.Teams {
			if team.Name == searchTeam {
				found = team
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any team with name: %s", searchTeam),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("description", found.Description)

		return nil
	})
}
