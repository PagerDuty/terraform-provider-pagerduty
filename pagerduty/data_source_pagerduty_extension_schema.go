package pagerduty

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func dataSourcePagerDutyExtensionSchema() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePagerDutyExtensionSchemaRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use `name` instead. This attribute will be removed in a future version",
			},
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

	resp, _, err := client.ExtensionSchemas.List()
	if err != nil {
		return err
	}

	var found *pagerduty.ExtensionSchema = findExtensionSchema(resp.ExtensionSchemas, searchName)

	if found == nil {
		return fmt.Errorf("Unable to locate any extension schema with the name: %s", searchName)
	}

	d.SetId(found.ID)
	d.Set("name", found.Label)
	d.Set("type", found.Type)

	return nil
}

func findExtensionSchema(extensionSchemas []*pagerduty.ExtensionSchema, searchName string) *pagerduty.ExtensionSchema {
	r := regexp.MustCompile("(?i)" + searchName)
	var closeMatch *pagerduty.ExtensionSchema
	for _, extensionSchema := range extensionSchemas {
		if searchName == extensionSchema.Label {
			return extensionSchema
		}
		if r.MatchString(extensionSchema.Label) {
			closeMatch = extensionSchema
		}
	}

	if closeMatch != nil {
		return closeMatch
	}
	return nil
}
