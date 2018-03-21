package pagerduty

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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
				Required: false,
			},
		},
	}
}

func dataSourcePagerDutyTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty team")

	searchTeam := d.Get("name").(string)

	resp, _, err := client.Teams.List(&pagerduty.ListTeamsOptions{})
	if err != nil {
		return err
	}

	var found *pagerduty.Team

	for _, team := range resp.Teams {
		if team.Name == searchTeam {
			found = team
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any team with name: %s", searchTeam)
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)
	d.Set("description", found.Description)

	return nil
}
