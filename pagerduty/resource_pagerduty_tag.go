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

func resourcePagerDutyTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyTagCreate,
		ReadContext:   resourcePagerDutyTagRead,
		DeleteContext: resourcePagerDutyTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"summary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildTagStruct(d *schema.ResourceData) *pagerduty.Tag {
	tag := &pagerduty.Tag{
		Label: d.Get("label").(string),
		Type:  "tag",
	}

	if attr, ok := d.GetOk("summary"); ok {
		tag.Summary = attr.(string)
	}

	return tag
}

func resourcePagerDutyTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	tag := buildTagStruct(d)

	log.Printf("[INFO] Creating PagerDuty tag %s", tag.Label)

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if tag, _, err := client.Tags.Create(tag); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if tag != nil {
			d.SetId(tag.ID)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return resourcePagerDutyTagRead(ctx, d, meta)

}

func resourcePagerDutyTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty tag %s", d.Id())

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		tag, _, err := client.Tags.Get(d.Id())
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		if tag != nil {
			log.Printf("Tag Type: %v", tag.Type)
			d.Set("label", tag.Label)
			d.Set("summary", tag.Summary)
			d.Set("html_url", tag.HTMLURL)
		}
		return nil
	}))
}

func resourcePagerDutyTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty tag %s", d.Id())

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if _, err := client.Tags.Delete(d.Id()); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	d.SetId("")

	// giving the API time to catchup
	time.Sleep(time.Second)
	return nil
}
