package pagerduty

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyUserNotificationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyUserNotificationRuleCreate,
		Read:   resourcePagerDutyUserNotificationRuleRead,
		Update: resourcePagerDutyUserNotificationRuleUpdate,
		Delete: resourcePagerDutyUserNotificationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"contact_method_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"contact_method_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"start_delay_in_minutes": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"urgency": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"high",
					"low",
				}),
			},
		},
	}
}

func buildUserNotificationRuleStruct(d *schema.ResourceData) *pagerduty.NotificationRule {
	rule := &pagerduty.NotificationRule{
		Type:                "assignment_notification_rule",
		StartDelayInMinutes: d.Get("start_delay_in_minutes").(int),
		Urgency:             d.Get("urgency").(string),
		ContactMethod: &pagerduty.ContactMethodReference{
			ID:   d.Get("contact_method_id").(string),
			Type: d.Get("contact_method_type").(string),
		},
	}

	return rule
}

func resourcePagerDutyUserNotificationRuleParseID(id string) (string, string) {
	// userID, ruleID
	parts := strings.Split(id, ":")
	return parts[0], parts[1]
}

func resourcePagerDutyUserNotificationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID := d.Get("user_id").(string)

	rule := buildUserNotificationRuleStruct(d)

	resp, _, err := client.Users.CreateNotificationRule(userID, rule)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", userID, resp.ID))

	return resourcePagerDutyUserNotificationRuleUpdate(d, meta)
}

func resourcePagerDutyUserNotificationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	userID, ruleID := resourcePagerDutyUserNotificationRuleParseID(d.Id())

	resp, _, err := client.Users.GetNotificationRule(userID, ruleID)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("user_id", userID)
	d.Set("contact_method_id", resp.ContactMethod.ID)
	d.Set("contact_method_type", resp.ContactMethod.Type)
	d.Set("urgency", resp.Urgency)
	d.Set("start_delay_in_minutes", resp.StartDelayInMinutes)

	return nil
}

func resourcePagerDutyUserNotificationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	rule := buildUserNotificationRuleStruct(d)

	log.Printf("[INFO] Updating PagerDuty user notification rule: %s", d.Id())

	userID, ruleID := resourcePagerDutyUserNotificationRuleParseID(d.Id())

	if _, _, err := client.Users.UpdateNotificationRule(userID, ruleID, rule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyUserNotificationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty user notification rule: %s", d.Id())

	userID, ruleID := resourcePagerDutyUserNotificationRuleParseID(d.Id())

	if _, err := client.Users.DeleteNotificationRule(userID, ruleID); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
