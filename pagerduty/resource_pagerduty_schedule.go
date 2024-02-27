package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutySchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyScheduleCreate,
		Read:   resourcePagerDutyScheduleRead,
		Update: resourcePagerDutyScheduleUpdate,
		Delete: resourcePagerDutyScheduleDelete,
		CustomizeDiff: func(context context.Context, diff *schema.ResourceDiff, i interface{}) error {
			ln := diff.Get("layer.#").(int)
			for li := 0; li <= ln; li++ {
				rn := diff.Get(fmt.Sprintf("layer.%d.restriction.#", li)).(int)
				for ri := 0; ri <= rn; ri++ {
					t := diff.Get(fmt.Sprintf("layer.%d.restriction.%d.type", li, ri)).(string)
					isStartDayOfWeekSetWhenDailyRestrictionType := t == "daily_restriction" && diff.Get(fmt.Sprintf("layer.%d.restriction.%d.start_day_of_week", li, ri)).(int) != 0
					if isStartDayOfWeekSetWhenDailyRestrictionType {
						return fmt.Errorf("start_day_of_week must only be set for a weekly_restriction schedule restriction type")
					}
					isStartDayOfWeekNotSetWhenWeeklyRestrictionType := t == "weekly_restriction" && diff.Get(fmt.Sprintf("layer.%d.restriction.%d.start_day_of_week", li, ri)).(int) == 0
					if isStartDayOfWeekNotSetWhenWeeklyRestrictionType {
						return fmt.Errorf("start_day_of_week must be set for a weekly_restriction schedule restriction type")
					}
					ds := diff.Get(fmt.Sprintf("layer.%d.restriction.%d.duration_seconds", li, ri)).(int)
					if t == "daily_restriction" && ds >= 3600*24 {
						return fmt.Errorf("duration_seconds for a daily_restriction schedule restriction type must be shorter than a day")
					}
				}
			}
			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"time_zone": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: util.ValidateTZValueDiagFunc,
			},

			"overflow": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},

			"layer": {
				Type:     schema.TypeList,
				Required: true,
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
							Required: true,
							ValidateFunc: func(v interface{}, k string) ([]string, []error) {
								var errors []error
								value := v.(string)
								_, err := time.Parse(time.RFC3339, value)
								if err != nil {
									errors = append(errors, genErrorTimeFormatRFC339(value, k))
								}

								return nil, errors
							},
							DiffSuppressFunc: suppressScheduleLayerStartDiff,
						},

						"end": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validateRFC3339,
							DiffSuppressFunc: suppressRFC3339Diff,
						},

						"rotation_virtual_start": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validateRFC3339,
							DiffSuppressFunc: suppressRFC3339Diff,
						},

						"rotation_turn_length_seconds": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3600, 365*24*3600),
						},

						"users": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"rendered_coverage_percentage": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"restriction": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"daily_restriction",
											"weekly_restriction",
										}),
									},

									"start_time_of_day": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]`), "must be of 00:00:00 format"),
									},

									"start_day_of_week": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 7),
									},

									"duration_seconds": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 7*24*3600-1),
									},
								},
							},
						},
					},
				},
			},
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"final_schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rendered_coverage_percentage": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildScheduleStruct(d *schema.ResourceData) (*pagerduty.Schedule, error) {
	layers, err := expandScheduleLayers(d.Get("layer"))
	if err != nil {
		return nil, err
	}

	schedule := &pagerduty.Schedule{
		Name:           d.Get("name").(string),
		TimeZone:       d.Get("time_zone").(string),
		ScheduleLayers: layers,
	}

	if attr, ok := d.GetOk("description"); ok {
		schedule.Description = attr.(string)
	}

	if attr, ok := d.GetOk("teams"); ok {
		schedule.Teams = expandSchedTeams(attr.([]interface{}))
	}

	return schedule, nil
}

func resourcePagerDutyScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	schedule, err := buildScheduleStruct(d)
	if err != nil {
		return err
	}

	o := &pagerduty.CreateScheduleOptions{}

	if v, ok := d.GetOk("overflow"); ok {
		o.Overflow = v.(bool)
	}

	log.Printf("[INFO] Creating PagerDuty schedule: %s", schedule.Name)

	schedule, _, err = client.Schedules.Create(schedule, o)
	if err != nil {
		return err
	}

	d.SetId(schedule.ID)

	return resourcePagerDutyScheduleRead(d, meta)
}

func resourcePagerDutyScheduleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading PagerDuty schedule: %s", d.Id())
	return fetchSchedule(d, meta, handleNotFoundError)
}

func fetchSchedule(d *schema.ResourceData, meta interface{}, errCallback func(error, *schema.ResourceData) error) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		schedule, _, err := client.Schedules.Get(d.Id(), &pagerduty.GetScheduleOptions{})
		if err != nil {
			log.Printf("[WARN] Schedule read error")
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			errResp := errCallback(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}
			return nil
		}
		if schedule != nil {
			d.Set("name", schedule.Name)
			d.Set("time_zone", schedule.TimeZone)
			d.Set("description", schedule.Description)

			layers, err := flattenScheduleLayers(schedule.ScheduleLayers)
			if err != nil {
				return retry.NonRetryableError(err)
			}

			if err := d.Set("layer", layers); err != nil {
				return retry.NonRetryableError(err)
			}
			if err := d.Set("teams", flattenShedTeams(schedule.Teams)); err != nil {
				return retry.NonRetryableError(fmt.Errorf("error setting teams: %s", err))
			}
			if err := d.Set("final_schedule", flattenScheFinalSchedule(schedule.FinalSchedule)); err != nil {
				return retry.NonRetryableError(fmt.Errorf("error setting final_schedule: %s", err))
			}

		}
		return nil
	})

	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return nil
}

func resourcePagerDutyScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	schedule, err := buildScheduleStruct(d)
	if err != nil {
		return err
	}

	opts := &pagerduty.UpdateScheduleOptions{}

	if v, ok := d.GetOk("overflow"); ok {
		opts.Overflow = v.(bool)
	}

	// A schedule layer can never be removed but it can be ended.
	// Here we determine which layer has been removed from the configuration
	// and we mark it as ended. This is to avoid diff issues.

	if d.HasChange("layer") {
		oraw, nraw := d.GetChange("layer")

		osl, err := expandScheduleLayers(oraw.([]interface{}))
		if err != nil {
			return err
		}

		nsl, err := expandScheduleLayers(nraw.([]interface{}))
		if err != nil {
			return err
		}

		// Checks to see if new schedule layers (nsl) include all old schedule layers (osl)
		for _, o := range osl {
			found := false
			for _, n := range nsl {
				// layer is found in both nsl and osl
				if o.ID == n.ID {
					found = true
				}
			}

			// If layer is not found in new schedule layers (nsl) set end value for layer
			if !found {
				end, err := timeToUTC(time.Now().Format(time.RFC3339))
				if err != nil {
					return err
				}
				endStr := end.String()
				o.End = &endStr
				schedule.ScheduleLayers = append(schedule.ScheduleLayers, o)
			}
		}
	}

	log.Printf("[INFO] Updating PagerDuty schedule: %s", d.Id())

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, _, err := client.Schedules.Update(d.Id(), schedule, opts); err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return nil
}

func resourcePagerDutyScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
	scheduleId := d.Id()

	log.Printf("[INFO] Starting deletion process of Schedule %s", scheduleId)
	var scheduleData *pagerduty.Schedule
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, _, err := client.Schedules.Get(scheduleId, &pagerduty.GetScheduleOptions{})
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
		}
		scheduleData = resp
		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	log.Printf("[INFO] Listing Escalation Policies that use schedule : %s", scheduleId)
	// Extracting Escalation Policies that use this Schedule
	epsUsingThisSchedule, err := extractEPsUsingASchedule(client, scheduleData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty schedule: %s", scheduleId)
	// Retrying to give other resources (such as escalation policies) to delete
	retryErr = retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, err := client.Schedules.Delete(scheduleId); err != nil {
			if !isErrCode(err, 400) {
				return retry.RetryableError(err)
			}
			isErrorScheduleUsedByEP := func(e *pagerduty.Error) bool {
				return strings.Compare(fmt.Sprintf("%v", e.Errors), "[Schedule can't be deleted if it's being used by escalation policies]") == 0
			}
			isErrorScheduleWOpenIncidents := func(e *pagerduty.Error) bool {
				return strings.Compare(fmt.Sprintf("%v", e.Errors), "[Schedule can't be deleted if it's being used by an escalation policy snapshot with open incidents]") == 0
			}

			// Handling of specific http 400 errors from API call DELETE /schedules
			e, ok := err.(*pagerduty.Error)
			if !ok || !isErrorScheduleUsedByEP(e) && !isErrorScheduleWOpenIncidents(e) {
				return retry.NonRetryableError(err)
			}

			var workaroundErr error
			// An Schedule with open incidents related can't be remove till those
			// incidents have been resolved.
			linksToIncidentsOpen, workaroundErr := listIncidentsOpenedRelatedToSchedule(client, scheduleData, epsUsingThisSchedule)
			if workaroundErr != nil {
				err = fmt.Errorf("%v; %w", err, workaroundErr)
				return retry.NonRetryableError(err)
			}

			hasToShowIncidentRemediationMessage := len(linksToIncidentsOpen) > 0
			if hasToShowIncidentRemediationMessage {
				var urlLinksMessage string
				for _, incident := range linksToIncidentsOpen {
					urlLinksMessage = fmt.Sprintf("%s\n%s", urlLinksMessage, incident)
				}
				return retry.NonRetryableError(fmt.Errorf("Before destroying Schedule %q You must first resolve or reassign the following incidents related with Escalation Policies using this Schedule... %s", scheduleId, urlLinksMessage))
			}

			// Returning at this point because the open incident (s) blocking the
			// deletion of the Schedule can't be tracked.
			if isErrorScheduleWOpenIncidents(e) && !hasToShowIncidentRemediationMessage {
				return retry.NonRetryableError(e)
			}

			epsDataUsingThisSchedule, errFetchingFullEPs := fetchEPsDataUsingASchedule(epsUsingThisSchedule, client)
			if errFetchingFullEPs != nil {
				err = fmt.Errorf("%v; %w", err, errFetchingFullEPs)
				return retry.RetryableError(err)
			}

			errBlockingBecauseOfEPs := detectUseOfScheduleByEPsWithOneLayer(scheduleId, epsDataUsingThisSchedule)
			if errBlockingBecauseOfEPs != nil {
				return retry.NonRetryableError(errBlockingBecauseOfEPs)
			}

			// Workaround for Schedule being used by escalation policies error
			log.Printf("[INFO] Dissociating Escalation Policies that use the Schedule: %s", scheduleId)
			workaroundErr = dissociateScheduleFromEPs(client, scheduleId, epsDataUsingThisSchedule)
			if workaroundErr != nil {
				err = fmt.Errorf("%v; %w", err, workaroundErr)
			}
			return retry.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	d.SetId("")

	return nil
}

func expandScheduleLayers(v interface{}) ([]*pagerduty.ScheduleLayer, error) {
	var scheduleLayers []*pagerduty.ScheduleLayer

	for _, sl := range v.([]interface{}) {
		rsl := sl.(map[string]interface{})

		// This is a temporary fix to prevent getting back the wrong rotation_virtual_start time.
		// The background here is that if a user specifies a rotation_virtual_start time to be:
		// "2017-09-01T10:00:00+02:00" the API returns back "2017-09-01T12:00:00+02:00".
		// With this fix in place, we get the correct rotation_virtual_start time, thus
		// eliminating the diff issues we've been seeing in the past.
		// This has been confirmed working by PagerDuty support.
		rvs, err := timeToUTC(rsl["rotation_virtual_start"].(string))
		if err != nil {
			return nil, err
		}

		// The type of layer.*.end is schema.TypeString. If the end is an empty string, it means the layer does not end.
		// A client should send a payload including `"end": null` to unset the end of layer.
		scheduleLayer := &pagerduty.ScheduleLayer{
			ID:                        rsl["id"].(string),
			Name:                      rsl["name"].(string),
			Start:                     rsl["start"].(string),
			End:                       stringTypeToStringPtr(rsl["end"].(string)),
			RotationVirtualStart:      rvs.String(),
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

	return scheduleLayers, nil
}

func flattenScheduleLayers(v []*pagerduty.ScheduleLayer) ([]map[string]interface{}, error) {
	var scheduleLayers []map[string]interface{}

	for _, sl := range v {
		// A schedule layer can never be removed but it can be ended.
		// Here we check each layer and if it has been ended we don't read it back
		// because it's not relevant anymore.
		endStr := stringPtrToStringType(sl.End)
		if endStr != "" {
			end, err := timeToUTC(endStr)
			if err != nil {
				return nil, err
			}

			if time.Now().UTC().After(end) {
				continue
			}
		}
		scheduleLayer := map[string]interface{}{
			"id":                           sl.ID,
			"name":                         sl.Name,
			"end":                          endStr,
			"start":                        sl.Start,
			"rotation_virtual_start":       sl.RotationVirtualStart,
			"rotation_turn_length_seconds": sl.RotationTurnLengthSeconds,
			"rendered_coverage_percentage": renderRoundedPercentage(sl.RenderedCoveragePercentage),
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

	return resultReversed, nil
}

// the expandShedTeams and flattenSchedTeams are based on the expandTeams and flattenTeams functions in the user
// resource. added these functions here for maintainability
func expandSchedTeams(v interface{}) []*pagerduty.TeamReference {
	var teams []*pagerduty.TeamReference

	for _, t := range v.([]interface{}) {
		team := &pagerduty.TeamReference{
			ID:   t.(string),
			Type: "team_reference",
		}
		teams = append(teams, team)
	}

	return teams
}

func flattenShedTeams(teams []*pagerduty.TeamReference) []string {
	res := make([]string, len(teams))
	for i, t := range teams {
		res[i] = t.ID
	}

	return res
}

func flattenScheFinalSchedule(finalSche *pagerduty.SubSchedule) []map[string]interface{} {
	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["name"] = finalSche.Name
	elem["rendered_coverage_percentage"] = renderRoundedPercentage(finalSche.RenderedCoveragePercentage)
	res = append(res, elem)

	return res
}

func listIncidentsOpenedRelatedToSchedule(c *pagerduty.Client, schedule *pagerduty.Schedule, epIDs []string) ([]string, error) {
	var incidents []*pagerduty.Incident
	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		var err error
		options := &pagerduty.ListIncidentsOptions{
			DateRange: "all",
			Statuses:  []string{"triggered", "acknowledged"},
			Limit:     100,
		}
		if len(schedule.Users) > 0 {
			for _, u := range schedule.Users {
				options.UserIDs = append(options.UserIDs, u.ID)
			}
		}

		incidents, err = c.Incidents.ListAll(options)
		if err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return retry.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		return nil, retryErr
	}

	filterIncidentsByEPs := func(incidents []*pagerduty.Incident, eps []string) []*pagerduty.Incident {
		var r []*pagerduty.Incident

		matchIndex := make(map[string]bool)
		for _, ep := range eps {
			matchIndex[ep] = true
		}
		for _, inc := range incidents {
			if matchIndex[inc.EscalationPolicy.ID] {
				r = append(r, inc)
			}
		}
		return r
	}
	incidents = filterIncidentsByEPs(incidents, epIDs)

	var linksToIncidents []string
	for _, inc := range incidents {
		linksToIncidents = append(linksToIncidents, inc.HTMLURL)
	}
	return linksToIncidents, nil
}

func extractEPsUsingASchedule(c *pagerduty.Client, schedule *pagerduty.Schedule) ([]string, error) {
	eps := []string{}
	for _, ep := range schedule.EscalationPolicies {
		eps = append(eps, ep.ID)
	}
	return eps, nil
}

func dissociateScheduleFromEPs(c *pagerduty.Client, scheduleID string, eps []*pagerduty.EscalationPolicy) error {
	for _, ep := range eps {
		errorMessage := fmt.Sprintf("Error while trying to dissociate Schedule %q from Escalation Policy %q", scheduleID, ep.ID)
		err := removeScheduleFromEP(c, scheduleID, ep)
		if err != nil {
			return fmt.Errorf("%w; %s", err, errorMessage)
		}
	}

	return nil
}

func removeScheduleFromEP(c *pagerduty.Client, scheduleID string, ep *pagerduty.EscalationPolicy) error {
	needsToUpdate := false
	epr := ep.EscalationRules
	// If the Escalation Policy using this Schedule has only one layer then this
	// workaround isn't applicable.
	if len(epr) < 2 {
		return nil
	}

	for ri, r := range epr {
		for index, target := range r.Targets {
			isScheduleConfiguredInEscalationRule := target.Type == "schedule_reference" && target.ID == scheduleID
			if !isScheduleConfiguredInEscalationRule {
				continue
			}

			if len(r.Targets) > 1 {
				// Removing Schedule as a configured Target from the Escalation Rules
				// slice.
				r.Targets = append(r.Targets[:index], r.Targets[index+1:]...)
			} else {
				// Removing Escalation Rules that will end up having no target configured.
				isLastRule := ri == len(epr)-1
				if isLastRule {
					epr = epr[:ri]
				} else {
					epr = append(epr[:ri], epr[ri+1:]...)
				}
			}
			needsToUpdate = true
		}
	}
	if !needsToUpdate {
		return nil
	}
	ep.EscalationRules = epr

	retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, _, err := c.EscalationPolicies.Update(ep.ID, ep)
		if err != nil {
			if !isErrCode(err, 404) {
				return retry.RetryableError(err)
			}
		}
		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	return nil
}

func detectUseOfScheduleByEPsWithOneLayer(scheduleId string, eps []*pagerduty.EscalationPolicy) error {
	epsFound := []*pagerduty.EscalationPolicy{}
	for _, ep := range eps {
		epHasNoLayers := len(ep.EscalationRules) == 0
		if epHasNoLayers {
			continue
		}

		epHasOneLayer := len(ep.EscalationRules) == 1 && len(ep.EscalationRules[0].Targets) == 1
		epHasMultipleLayersButAllTargetThisSchedule := func() bool {
			var meetCondition bool
			if len(ep.EscalationRules) == 1 {
				return meetCondition
			}
			meetConditionMapping := make(map[int]bool)
			for epli, epLayer := range ep.EscalationRules {
				meetConditionMapping[epli] = false
				isTargetingThisSchedule := epLayer.Targets[0].Type == "schedule_reference" && epLayer.Targets[0].ID == scheduleId
				if len(epLayer.Targets) == 1 && isTargetingThisSchedule {
					meetConditionMapping[epli] = true
				}
			}
			for _, mc := range meetConditionMapping {
				if !mc {
					meetCondition = false
					break
				}
				meetCondition = true
			}

			return meetCondition
		}

		if !epHasOneLayer && !epHasMultipleLayersButAllTargetThisSchedule() {
			continue
		}
		epsFound = append(epsFound, ep)
	}

	if len(epsFound) == 0 {
		return nil
	}

	tfState, err := getTFStateSnapshot()
	if err != nil {
		return err
	}

	epsNames := []string{}
	for _, ep := range epsFound {
		epState := tfState.GetResourceStateById(ep.ID)

		// To cover the case when the Schedule is used by an Escalation Policy which
		// is not being managed by the same TF config which is managing this Schedule.
		if epState == nil {
			return fmt.Errorf("It is not possible to continue with the destruction of the Schedule %q, because it is being used by Escalation Policy %q which has only one layer configured. Nevertheless, the mentioned Escalation Policy is not managed by this Terraform configuration. So in order to unblock this resource destruction, We suggest you to first make the appropiate changes on the Escalation Policy %s and come back for retrying.", scheduleId, ep.ID, ep.HTMLURL)
		}
		epsNames = append(epsNames, epState.Name)
	}

	displayError := fmt.Errorf(`It is not possible to continue with the destruction of the Schedule %q, because it is being used by the Escalation Policy %[2]q which has only one layer configured. Therefore in order to unblock this resource destruction, We suggest you to first execute "terraform apply (or destroy, please act accordingly) -target=pagerduty_escalation_policy.%[2]s"`, scheduleId, epsNames[0])
	if len(epsNames) > 1 {
		var epsListMessage string
		for _, ep := range epsNames {
			epsListMessage = fmt.Sprintf("%s\n%s", epsListMessage, ep)
		}
		displayError = fmt.Errorf(`It is not possible to continue with the destruction of the Schedule %q, because it is being used by multiple Escalation Policies which have only one layer configured. Therefore in order to unblock this resource destruction, We suggest you to first execute "terraform apply (or destroy, please act accordingly) -target=pagerduty_escalation_policy.<Escalation Policy Name here>". e.g: "terraform apply -target=pagerduty_escalation_policy.example". Replacing the example name with the following Escalation Policies which are blocking the deletion of the Schedule...%s`, scheduleId, epsListMessage)
	}

	return displayError
}

func fetchEPsDataUsingASchedule(eps []string, c *pagerduty.Client) ([]*pagerduty.EscalationPolicy, error) {
	fullEPs := []*pagerduty.EscalationPolicy{}
	for _, epID := range eps {
		retryErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
			ep, _, err := c.EscalationPolicies.Get(epID, &pagerduty.GetEscalationPolicyOptions{})
			if err != nil {
				if isErrCode(err, http.StatusBadRequest) {
					return retry.NonRetryableError(err)
				}

				return retry.RetryableError(err)
			}
			fullEPs = append(fullEPs, ep)
			return nil
		})
		if retryErr != nil {
			return fullEPs, retryErr
		}
	}

	return fullEPs, nil
}
