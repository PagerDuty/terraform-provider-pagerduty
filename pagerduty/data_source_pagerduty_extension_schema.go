package pagerduty

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty Extension Schema")

	searchName := d.Get("name").(string)

	resp, _, err := client.ExtensionSchemas.List(&pagerduty.ListExtensionSchemasOptions{Query: searchName})
	if err != nil {
		return err
	}

	var found *pagerduty.ExtensionSchema

	for _, schema := range resp.ExtensionSchemas {
		if strings.EqualFold(schema.Label, searchName) {
			found = schema
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any extension schema with the name: %s", searchName)
	}

	d.SetId(found.ID)
	d.Set("name", found.Label)
	d.Set("type", found.Type)

	return nil
}
