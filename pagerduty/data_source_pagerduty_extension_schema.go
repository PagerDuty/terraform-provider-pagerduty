package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyExtensionSchema() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyExtensionSchemaRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePagerDutyExtensionSchemaRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty Extension Schema")

	searchName := d.Get("name").(string)

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, _, err := client.ExtensionSchemas.List(&pagerduty.ListExtensionSchemasOptions{Query: searchName})
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)
			return retry.RetryableError(err)
		}

		var found *pagerduty.ExtensionSchema

		for _, schema := range resp.ExtensionSchemas {
			if strings.EqualFold(schema.Label, searchName) {
				found = schema
				break
			}
		}

		if found == nil {
			return retry.NonRetryableError(
				fmt.Errorf("Unable to locate any extension schema with the name: %s", searchName),
			)
		}

		d.SetId(found.ID)
		d.Set("name", found.Label)
		d.Set("type", found.Type)

		return nil
	})
}
