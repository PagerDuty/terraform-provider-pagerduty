package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyAddon() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyAddonCreate,
		Read:   resourcePagerDutyAddonRead,
		Update: resourcePagerDutyAddonUpdate,
		Delete: resourcePagerDutyAddonDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"src": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildAddonStruct(d *schema.ResourceData) *pagerduty.Addon {
	addon := &pagerduty.Addon{
		Name: d.Get("name").(string),
		Src:  d.Get("src").(string),
		Type: "full_page_addon",
	}

	return addon
}

func fetchPagerDutyAddon(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		addon, _, err := client.Addons.Get(d.Id())
		if err != nil {
			log.Printf("[WARN] Service read error")
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(errResp)
			}

			return nil
		}

		d.Set("name", addon.Name)
		d.Set("src", addon.Src)

		return nil
	})
}

func resourcePagerDutyAddonCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	addon := buildAddonStruct(d)

	log.Printf("[INFO] Creating PagerDuty add-on %s", addon.Name)

	addon, _, err = client.Addons.Install(addon)
	if err != nil {
		return err
	}

	d.SetId(addon.ID)
	// Retrying on creates incase of eventual consistency on creation
	return fetchPagerDutyAddon(d, meta, genError)
}

func resourcePagerDutyAddonRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading PagerDuty add-on %s", d.Id())
	return fetchPagerDutyAddon(d, meta, handleNotFoundError)
}

func resourcePagerDutyAddonUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	addon := buildAddonStruct(d)

	log.Printf("[INFO] Updating PagerDuty add-on %s", d.Id())

	if _, _, err := client.Addons.Update(d.Id(), addon); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyAddonDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty add-on %s", d.Id())

	if _, err := client.Addons.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
