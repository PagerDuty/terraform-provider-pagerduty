package pagerduty

import (
	"log"

	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyExtensionCreate,
		Read:   resourcePagerDutyExtensionRead,
		Update: resourcePagerDutyExtensionUpdate,
		Delete: resourcePagerDutyExtensionDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyExtensionImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"endpoint_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extension_objects": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"extension_schema": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func buildExtensionStruct(d *schema.ResourceData) *pagerduty.Extension {
	Extension := &pagerduty.Extension{
		Name:        d.Get("name").(string),
		Type:        "extension",
		EndpointURL: d.Get("endpoint_url").(string),
		ExtensionSchema: &pagerduty.ExtensionSchemaReference{
			Type: "extension_schema_reference",
			ID:   d.Get("extension_schema").(string),
		},
		ExtensionObjects: expandServices(d.Get("extension_objects")),
	}

	return Extension
}

func resourcePagerDutyExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	Extension := buildExtensionStruct(d)

	log.Printf("[INFO] Creating PagerDuty extension %s", Extension.Name)

	Extension, _, err := client.Extensions.Create(Extension)
	if err != nil {
		return err
	}

	d.SetId(Extension.ID)

	return nil
}

func resourcePagerDutyExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty extension %s", d.Id())

	extension, _, err := client.Extensions.Get(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("name", extension.Name)

	return nil
}

func resourcePagerDutyExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	extension := buildExtensionStruct(d)

	log.Printf("[INFO] Updating PagerDuty extension %s", d.Id())

	if _, _, err := client.Extensions.Update(d.Id(), extension); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty extension %s", d.Id())

	if _, err := client.Extensions.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourcePagerDutyExtensionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*pagerduty.Client)

	extension, _, err := client.Extensions.Get(d.Id())

	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_extension. Expecting an importation ID for extension.")
	}

	d.Set("endpoint_url", extension.EndpointURL)
	d.Set("extension_objects", []string{extension.ExtensionObjects[0].ID})
	d.Set("extension_schema", extension.ExtensionSchema.ID)

	return []*schema.ResourceData{d}, err
}
