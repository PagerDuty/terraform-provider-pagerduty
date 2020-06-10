package pagerduty

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty user")

	searchEmail := d.Get("email").(string)

	o := &pagerduty.ListUsersOptions{
		Query: searchEmail,
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Users.List(o)
		if err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		}

		var found *pagerduty.User

		for _, user := range resp.Users {
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
