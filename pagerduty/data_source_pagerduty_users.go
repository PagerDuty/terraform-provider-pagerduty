package pagerduty

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyUsersRead,

		Schema: map[string]*schema.Schema{
			"team": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePagerDutyUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	searchTeam := []string{d.Get("team").(string)}
	log.Printf("[DEBUG] *** Getting users for team: %s", searchTeam)

	o := &pagerduty.ListUsersOptions{
		TeamIDs: searchTeam,
	}

	resp, _, err := client.Users.List(o)
	if err != nil {
		log.Println("[ERROR] *** We've errored: " + err.Error())
		return err
	}

	filteredUsers := resp.Users[:]
	if len(filteredUsers) < 1 {
		return fmt.Errorf("Unable to locate any user for Team: %s", searchTeam)
	}

	var userids []string

	//var found *pagerduty.User
	for _, user := range filteredUsers {
		userids = append(userids, user.ID)
	}

	log.Printf("[DEBUG] *** Setting the user ids: %d", len(userids))
	d.SetId(searchTeam[0])

	if err := d.Set("users", userids); err != nil {
		return fmt.Errorf("Error setting users: %s", err)
	}

	return nil
}
