package pagerduty

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyTag() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyTagCreate,
		Read:   resourcePagerDutyTagRead,
		Delete: resourcePagerDutyTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourcePagerDutyTagCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	tag := buildTagStruct(d)

	log.Printf("[INFO] Creating PagerDuty tag %s", tag.Label)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
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
		return retryErr
	}

	return resourcePagerDutyTagRead(d, meta)

}

func resourcePagerDutyTagRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty tag %s", d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		tag, _, err := client.Tags.Get(d.Id())
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			errResp := handleNotFoundError(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}
		if tag != nil {
			log.Printf("Tag Type: %v", tag.Type)
			d.Set("label", tag.Label)
			d.Set("summary", tag.Summary)
			d.Set("html_url", tag.HTMLURL)
		}
		return nil
	})
}

func resourcePagerDutyTagDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty tag %s", d.Id())

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.Tags.Delete(d.Id()); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	d.SetId("")

	// giving the API time to catchup
	time.Sleep(time.Second)
	return nil
}
