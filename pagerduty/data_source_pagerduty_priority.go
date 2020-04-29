package pagerduty

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyPriority() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyPriorityRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the priority to find in the PagerDuty API",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyPriorityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty priority")

	searchTeam := d.Get("name").(string)

	resp, _, err := client.Priorities.List()
	if err != nil {
		return err
	}

	var found *pagerduty.Priority

	for _, priority := range resp.Priorities {
		if strings.EqualFold(priority.Name, searchTeam) {
			found = priority
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any priority with name: %s", searchTeam)
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)
	d.Set("description", found.Description)

	return nil
}
