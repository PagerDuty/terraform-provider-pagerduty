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
				Type:     schema.TypeSet,
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
		ExtensionObjects: expandServiceObjects(d.Get("extension_objects")),
	}

	return Extension
}

func resourcePagerDutyExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	extension := buildExtensionStruct(d)

	log.Printf("[INFO] Creating PagerDuty extension %s", extension.Name)

	extension, _, err := client.Extensions.Create(extension)
	if err != nil {
		return err
	}

	d.SetId(extension.ID)

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
	d.Set("summary", extension.Summary)
	d.Set("endpoint_url", extension.EndpointURL)
	d.Set("extension_objects", extension.ExtensionObjects)
	d.Set("extension_schema", extension.ExtensionSchema)

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
		if perr, ok := err.(*pagerduty.Error); ok && perr.Code == 5001 {
			log.Printf("[WARN] Extension (%s) not found, removing from state", d.Id())
			return nil
		}
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

func expandServiceObjects(v interface{}) []*pagerduty.ServiceReference {
	var services []*pagerduty.ServiceReference

	for _, srv := range v.(*schema.Set).List() {
		service := &pagerduty.ServiceReference{
			Type: "service_reference",
			ID:   srv.(string),
		}
		services = append(services, service)
	}

	return services
}
