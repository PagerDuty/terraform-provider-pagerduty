package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyAddon() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyAddonCreate,
		ReadContext:   resourcePagerDutyAddonRead,
		UpdateContext: resourcePagerDutyAddonUpdate,
		DeleteContext: resourcePagerDutyAddonDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func fetchPagerDutyAddon(ctx context.Context, d *schema.ResourceData, meta interface{}, handle404Errors bool) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		addon, _, err := client.Addons.Get(d.Id())
		if checkErr := getErrorHandler(handle404Errors)(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		d.Set("name", addon.Name)
		d.Set("src", addon.Src)

		return nil
	}))
}

func resourcePagerDutyAddonCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	addon := buildAddonStruct(d)

	log.Printf("[INFO] Creating PagerDuty add-on %s", addon.Name)

	addon, _, err = client.Addons.Install(addon)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(addon.ID)
	// Retrying on creates incase of eventual consistency on creation
	return fetchPagerDutyAddon(ctx, d, meta, false)
}

func resourcePagerDutyAddonRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading PagerDuty add-on %s", d.Id())
	return fetchPagerDutyAddon(ctx, d, meta, true)
}

func resourcePagerDutyAddonUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	addon := buildAddonStruct(d)

	log.Printf("[INFO] Updating PagerDuty add-on %s", d.Id())

	if _, _, err := client.Addons.Update(d.Id(), addon); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePagerDutyAddonDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty add-on %s", d.Id())

	if _, err := client.Addons.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
