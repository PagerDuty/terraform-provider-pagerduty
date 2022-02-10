package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyUserContactMethod() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyUserContactMethodCreate,
		Read:   resourcePagerDutyUserContactMethodRead,
		Update: resourcePagerDutyUserContactMethodUpdate,
		Delete: resourcePagerDutyUserContactMethodDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyUserContactMethodImport,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"email_contact_method",
					"phone_contact_method",
					"push_notification_contact_method",
					"sms_contact_method",
				}),
			},

			"send_short_email": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"country_code": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"blacklisted": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"label": {
				Type:     schema.TypeString,
				Required: true,
			},

			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildUserContactMethodStruct(d *schema.ResourceData) *pagerduty.ContactMethod {
	contactMethod := &pagerduty.ContactMethod{
		Type:    d.Get("type").(string),
		Label:   d.Get("label").(string),
		Address: d.Get("address").(string),
	}

	if v, ok := d.GetOk("send_short_email"); ok {
		contactMethod.SendShortEmail = v.(bool)
	}

	if v, ok := d.GetOk("country_code"); ok {
		contactMethod.CountryCode = v.(int)
	}

	if v, ok := d.GetOk("enabled"); ok {
		contactMethod.Enabled = v.(bool)
	}

	return contactMethod
}

func fetchPagerDutyUserContactMethod(d *schema.ResourceData, meta interface{}, handle404Errors bool) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID := d.Get("user_id").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err := client.Users.GetContactMethod(userID, d.Id())
		if checkErr := getErrorHandler(handle404Errors)(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		d.Set("address", resp.Address)
		d.Set("blacklisted", resp.BlackListed)
		d.Set("country_code", resp.CountryCode)
		d.Set("enabled", resp.Enabled)
		d.Set("label", resp.Label)
		d.Set("send_short_email", resp.SendShortEmail)
		d.Set("type", resp.Type)

		return nil
	})
}

func resourcePagerDutyUserContactMethodCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	userID := d.Get("user_id").(string)

	contactMethod := buildUserContactMethodStruct(d)

	resp, _, err := client.Users.CreateContactMethod(userID, contactMethod)
	if err != nil {
		return fmt.Errorf("error while creating contact method %s: %w", contactMethod.ID, err)
	}

	d.SetId(resp.ID)

	return fetchPagerDutyUserContactMethod(d, meta, false)
}

func resourcePagerDutyUserContactMethodRead(d *schema.ResourceData, meta interface{}) error {
	return fetchPagerDutyUserContactMethod(d, meta, true)
}

func resourcePagerDutyUserContactMethodUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	contactMethod := buildUserContactMethodStruct(d)

	log.Printf("[INFO] Updating PagerDuty user contact method %s", d.Id())

	userID := d.Get("user_id").(string)

	if _, _, err := client.Users.UpdateContactMethod(userID, d.Id(), contactMethod); err != nil {
		return fmt.Errorf("error while updating contact method %s: %w", contactMethod.ID, err)
	}

	return resourcePagerDutyUserContactMethodRead(d, meta)
}

func resourcePagerDutyUserContactMethodDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty user contact method %s", d.Id())

	userID := d.Get("user_id").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.Users.DeleteContactMethod(userID, d.Id()); err != nil {
			if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
				return checkErr.ReturnVal
			}
		}

		d.SetId("")

		return nil
	})
}

func resourcePagerDutyUserContactMethodImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	ids := strings.Split(d.Id(), ":")

	if len(ids) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_user_contact_method. Expecting an ID formed as '<user_id>:<contact_method_id>'")
	}
	uid, id := ids[0], ids[1]

	_, _, err = client.Users.GetContactMethod(uid, id)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.SetId(id)
	d.Set("user_id", uid)

	return []*schema.ResourceData{d}, nil
}
