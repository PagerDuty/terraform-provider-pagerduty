package pagerduty

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyScheduleOverride() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyScheduleOverrideCreate,
		Read:   resourcePagerDutyScheduleOverrideRead,
		Delete: resourcePagerDutyScheduleOverrideDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func buildScheduleOverrideStruct(d *schema.ResourceData) (*pagerduty.Override, error) {
	override := &pagerduty.Override{
		User: &pagerduty.UserReference{
			ID:   d.Get("user").(string),
			Type: "user",
		},
		Start: d.Get("start").(string),
		End:   d.Get("end").(string),
	}
	return override, nil
}

func resourcePagerDutyScheduleOverrideCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	override, err := buildScheduleOverrideStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating override for PagerDuty schedule: %s", d.Get("schedule").(string))

	newOverride, _, err := client.Schedules.CreateOverride(d.Get("schedule").(string), override)
	if err != nil {
		return err
	}

	d.SetId(newOverride.ID)

	return resourcePagerDutyScheduleOverrideRead(d, meta)
}

func resourcePagerDutyScheduleOverrideRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty schedule override: %s", d.Id())

	listOverridesOptions := pagerduty.ListOverridesOptions{
		Since: d.Get("start").(string),
		Until: d.Get("end").(string),
	}

	overrides, _, err := client.Schedules.ListOverrides(d.Get("schedule").(string), &listOverridesOptions)
	if err != nil {
		return err
	}

	var matchingOverrides []*pagerduty.Override
	for _, o := range overrides.Overrides {
		if o.ID == d.Id() {
			matchingOverrides = append(matchingOverrides, o)
		}
	}
	if len(matchingOverrides) != 1 {
		err := errors.New(fmt.Sprintf("Could not find override: %s", d.Get("ID").(string)))
		return handleNotFoundError(err, d)
	}

	d.Set("user", matchingOverrides[0].User.ID)
	d.Set("start", matchingOverrides[0].Start)
	d.Set("end", matchingOverrides[0].End)

	return nil
}

func resourcePagerDutyScheduleOverrideDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty schedule override: %s", d.Id())

	_, err := client.Schedules.DeleteOverride(d.Get("schedule").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
