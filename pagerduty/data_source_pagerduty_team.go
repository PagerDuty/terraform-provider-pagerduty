package pagerduty

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePagerDutyTeamRead,

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
			"parent": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePagerDutyTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty team")

	searchTeam := d.Get("name").(string)

	o := &pagerduty.ListTeamsOptions{
		Query: searchTeam,
	}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Teams.List(o)
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
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
		d.Set("parent", found.Parent)

		return nil
	}))
}
