package pagerduty

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyUserContactMethod() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyUserContactMethodCreate,
		Read:   resourcePagerDutyUserContactMethodRead,
		Update: resourcePagerDutyUserContactMethodUpdate,
		Delete: resourcePagerDutyUserContactMethodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			"contact_method_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourcePagerDutyUserContactMethodParseID(id string) (string, string) {
	// userID, contactMethodID
	parts := strings.Split(id, ":")
	return parts[0], parts[1]
}

func resourcePagerDutyUserContactMethodCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID := d.Get("user_id").(string)

	contactMethod := buildUserContactMethodStruct(d)

	resp, _, err := client.Users.CreateContactMethod(userID, contactMethod)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, resp.ID))

	return resourcePagerDutyUserContactMethodRead(d, meta)
}

func resourcePagerDutyUserContactMethodRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID, cmID := resourcePagerDutyUserContactMethodParseID(d.Id())

	resp, _, err := client.Users.GetContactMethod(userID, cmID)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("address", resp.Address)
	d.Set("blacklisted", resp.BlackListed)
	d.Set("contact_method_id", resp.ID)
	d.Set("country_code", resp.CountryCode)
	d.Set("enabled", resp.Enabled)
	d.Set("label", resp.Label)
	d.Set("send_short_email", resp.SendShortEmail)
	d.Set("type", resp.Type)
	d.Set("user_id", userID)

	return nil
}

func resourcePagerDutyUserContactMethodUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	contactMethod := buildUserContactMethodStruct(d)

	log.Printf("[INFO] Updating PagerDuty user contact method %s", d.Id())

	userID, cmID := resourcePagerDutyUserContactMethodParseID(d.Id())

	if _, _, err := client.Users.UpdateContactMethod(userID, cmID, contactMethod); err != nil {
		return err
	}

	return resourcePagerDutyUserContactMethodRead(d, meta)
}

func resourcePagerDutyUserContactMethodDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty user contact method %s", d.Id())

	userID, cmID := resourcePagerDutyUserContactMethodParseID(d.Id())

	if _, err := client.Users.DeleteContactMethod(userID, cmID); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
