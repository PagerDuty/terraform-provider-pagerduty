package pagerduty

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyScheduleCreate,
		Read:   resourcePagerDutyScheduleRead,
		Update: resourcePagerDutyScheduleUpdate,
		Delete: resourcePagerDutyScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"layer": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"start": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" {
									return false
								}
								return true
							},
						},
						"end": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rotation_virtual_start": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" {
									return false
								}
								return true
							},
						},
						"rotation_turn_length_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"users": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"restriction": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_time_of_day": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_day_of_week": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"duration_seconds": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildScheduleStruct(d *schema.ResourceData) *pagerduty.Schedule {
	schedule := &pagerduty.Schedule{
		Name:           d.Get("name").(string),
		TimeZone:       d.Get("time_zone").(string),
		ScheduleLayers: expandScheduleLayers(d.Get("layer")),
	}

	if attr, ok := d.GetOk("description"); ok {
		schedule.Description = attr.(string)
	}

	return schedule
}

func resourcePagerDutyScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	schedule := buildScheduleStruct(d)

	log.Printf("[INFO] Creating PagerDuty schedule: %s", schedule.Name)

	schedule, _, err := client.Schedules.Create(schedule)
	if err != nil {
		return err
	}

	d.SetId(schedule.ID)

	return resourcePagerDutyScheduleRead(d, meta)
}

func resourcePagerDutyScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty schedule: %s", d.Id())

	schedule, _, err := client.Schedules.Get(d.Id(), &pagerduty.GetScheduleOptions{})
	if err != nil {
		return err
	}

	d.Set("name", schedule.Name)
	d.Set("time_zone", schedule.TimeZone)
	d.Set("description", schedule.Description)

	if err := d.Set("layer", flattenScheduleLayers(schedule.ScheduleLayers)); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	schedule := buildScheduleStruct(d)

	log.Printf("[INFO] Updating PagerDuty schedule: %s", d.Id())

	if _, _, err := client.Schedules.Update(d.Id(), schedule); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty schedule: %s", d.Id())

	if _, err := client.Schedules.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandScheduleLayers(v interface{}) []*pagerduty.ScheduleLayer {
	var scheduleLayers []*pagerduty.ScheduleLayer

	for _, sl := range v.([]interface{}) {
		rsl := sl.(map[string]interface{})
		scheduleLayer := &pagerduty.ScheduleLayer{
			ID:                        rsl["id"].(string),
			Name:                      rsl["name"].(string),
			Start:                     rsl["start"].(string),
			End:                       rsl["end"].(string),
			RotationVirtualStart:      rsl["rotation_virtual_start"].(string),
			RotationTurnLengthSeconds: rsl["rotation_turn_length_seconds"].(int),
		}

		for _, slu := range rsl["users"].([]interface{}) {
			user := &pagerduty.UserReferenceWrapper{
				User: &pagerduty.UserReference{
					ID:   slu.(string),
					Type: "user",
				},
			}
			scheduleLayer.Users = append(scheduleLayer.Users, user)
		}

		for _, slr := range rsl["restriction"].([]interface{}) {
			rslr := slr.(map[string]interface{})

			restriction := &pagerduty.Restriction{
				Type:            rslr["type"].(string),
				StartTimeOfDay:  rslr["start_time_of_day"].(string),
				StartDayOfWeek:  rslr["start_day_of_week"].(int),
				DurationSeconds: rslr["duration_seconds"].(int),
			}

			scheduleLayer.Restrictions = append(scheduleLayer.Restrictions, restriction)
		}

		scheduleLayers = append(scheduleLayers, scheduleLayer)
	}

	return scheduleLayers
}

func flattenScheduleLayers(v []*pagerduty.ScheduleLayer) []map[string]interface{} {
	var scheduleLayers []map[string]interface{}

	for _, sl := range v {
		scheduleLayer := map[string]interface{}{
			"id":                           sl.ID,
			"name":                         sl.Name,
			"end":                          sl.End,
			"start":                        sl.Start,
			"rotation_virtual_start":       sl.RotationVirtualStart,
			"rotation_turn_length_seconds": sl.RotationTurnLengthSeconds,
		}

		var users []string

		for _, slu := range sl.Users {
			users = append(users, slu.User.ID)
		}

		scheduleLayer["users"] = users

		var restrictions []map[string]interface{}

		for _, slr := range sl.Restrictions {
			restriction := map[string]interface{}{
				"duration_seconds":  slr.DurationSeconds,
				"start_time_of_day": slr.StartTimeOfDay,
				"type":              slr.Type,
			}

			if slr.StartDayOfWeek > 0 {
				restriction["start_day_of_week"] = slr.StartDayOfWeek
			}

			restrictions = append(restrictions, restriction)
		}

		scheduleLayer["restriction"] = restrictions

		scheduleLayers = append(scheduleLayers, scheduleLayer)
	}

	// Reverse the final result and return it
	resultReversed := make([]map[string]interface{}, 0, len(scheduleLayers))

	for i := len(scheduleLayers) - 1; i >= 0; i-- {
		resultReversed = append(resultReversed, scheduleLayers[i])
	}

	return resultReversed
}
