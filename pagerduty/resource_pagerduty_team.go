package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyTeamCreate,
		ReadContext:   resourcePagerDutyTeamRead,
		UpdateContext: resourcePagerDutyTeamUpdate,
		DeleteContext: resourcePagerDutyTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
	return team
}

func resourcePagerDutyTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	team := buildTeamStruct(d)

	log.Printf("[INFO] Creating PagerDuty team %s", team.Name)

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if team, _, err := client.Teams.Create(team); err != nil {
			return resource.RetryableError(err)
		} else if team != nil {
			d.SetId(team.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	return resourcePagerDutyTeamRead(ctx, d, meta)

}

func resourcePagerDutyTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty team %s", d.Id())

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		team, _, err := client.Teams.Get(d.Id())
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		if team != nil {
			d.Set("name", team.Name)
			d.Set("description", team.Description)
			d.Set("html_url", team.HTMLURL)
		}
		return nil
	}))
}

func resourcePagerDutyTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	team := buildTeamStruct(d)

	log.Printf("[INFO] Updating PagerDuty team %s", d.Id())

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if _, _, err := client.Teams.Update(d.Id(), team); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	return resourcePagerDutyTeamRead(ctx, d, meta)
}

func resourcePagerDutyTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty team %s", d.Id())

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if _, err := client.Teams.Delete(d.Id()); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	d.SetId("")

	// giving the API time to catchup
	time.Sleep(time.Second)
	return nil
}
