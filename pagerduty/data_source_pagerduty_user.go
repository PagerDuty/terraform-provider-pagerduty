package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyUserRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePagerDutyUserRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty user")

	searchEmail := d.Get("email").(string)

	o := &pagerduty.ListUsersOptions{
		Query: searchEmail,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.Users.ListAll(o)
		if err != nil {
			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.FullUser

		for _, user := range resp {
			if user.Email == searchEmail {
				found = user
				break
			}
		}

		if found == nil {
			return resource.NonRetryableError(
				fmt.Errorf("Unable to locate any user with the email: %s", searchEmail),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Name)
		d.Set("email", found.Email)

		return nil
	})
}
