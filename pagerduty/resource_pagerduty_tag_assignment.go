package pagerduty

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyTagAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyTagAssignmentCreate,
		Read:   resourcePagerDutyTagAssignmentRead,
		Delete: resourcePagerDutyTagAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyTagAssignmentImport,
		},
		Schema: map[string]*schema.Schema{
			"entity_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validateValueDiagFunc([]string{
					"users",
					"teams",
					"escalation_policies",
				}),
			},
			"entity_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func buildTagAssignmentStruct(d *schema.ResourceData) *pagerduty.TagAssignment {
	assignment := &pagerduty.TagAssignment{
		// Hard-coding "tag_reference" here because using the "tag" type doesn't allow users to delete the tags as
		// they receive no tag id from the PagerDuty API at this time
		Type:       "tag_reference",
		EntityType: d.Get("entity_type").(string),
		EntityID:   d.Get("entity_id").(string),
	}
	if attr, ok := d.GetOk("tag_id"); ok {
		assignment.TagID = attr.(string)
	}

	return assignment
}

func resourcePagerDutyTagAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	assignment := buildTagAssignmentStruct(d)
	assignments := &pagerduty.TagAssignments{
		Add: []*pagerduty.TagAssignment{assignment},
	}

	log.Printf("[INFO] Creating PagerDuty tag assignment with tagID %s for %s entity with ID %s", assignment.TagID, assignment.EntityType, assignment.EntityID)

	retryErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.Tags.Assign(assignment.EntityType, assignment.EntityID, assignments); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else {
			// create tag_assignment id using the entityID.tagID as PagerDuty API does not return one
			assignmentID := createAssignmentID(assignment.EntityID, assignment.TagID)
			d.SetId(assignmentID)
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}
	// give PagerDuty 2 seconds to save the assignment correctly
	time.Sleep(2 * time.Second)
	return resourcePagerDutyTagAssignmentRead(d, meta)

}

func resourcePagerDutyTagAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	assignment := buildTagAssignmentStruct(d)

	log.Printf("[INFO] Reading PagerDuty tag assignment with tagID %s for %s entity with ID %s", assignment.TagID, assignment.EntityType, assignment.EntityID)

	ok, err := isFoundTagAssignmentEntity(d.Get("entity_id").(string), d.Get("entity_type").(string), meta)
	if err != nil {
		return err
	}
	if !ok {
		d.SetId("")
		return nil
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		if tagResponse, _, err := client.Tags.ListTagsForEntity(assignment.EntityType, assignment.EntityID); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return resource.NonRetryableError(err)
			}

			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if tagResponse != nil {
			var foundTag *pagerduty.Tag

			// loop tags and find matching ID
			for _, tag := range tagResponse.Tags {
				if tag.ID == assignment.TagID {
					foundTag = tag
					break
				}
			}
			if foundTag == nil {
				d.SetId("")
				return nil
			}
		}
		return nil
	})
}

func resourcePagerDutyTagAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	assignment := buildTagAssignmentStruct(d)
	assignments := &pagerduty.TagAssignments{
		Remove: []*pagerduty.TagAssignment{assignment},
	}
	log.Printf("[INFO] Deleting PagerDuty tag assignment with tagID %s for entityID %s", assignment.TagID, assignment.EntityID)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if _, err := client.Tags.Assign(assignment.EntityType, assignment.EntityID, assignments); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else {
			d.SetId("")
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return nil
}

func resourcePagerDutyTagAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), ".")
	if len(ids) != 3 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_tag_assignment. Expecting an importation ID formed as '<entity_type>.<entity_id>.<tag_id>'")
	}
	entityType, entityID, tagID := ids[0], ids[1], ids[2]

	client, err := meta.(*Config).Client()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	// give PagerDuty 2 seconds to save the assignment correctly
	time.Sleep(2 * time.Second)
	tagResponse, _, err := client.Tags.ListTagsForEntity(entityType, entityID)

	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("error importing pagerduty_tag_assignment: %s", err.Error())
	}
	var foundTag *pagerduty.Tag
	// loop tags and find matching ID
	for _, tag := range tagResponse.Tags {
		if tag.ID == tagID {
			// create tag_assignment id using the entityID.tagID as PagerDuty API does not return one
			assignmentID := createAssignmentID(entityID, tagID)
			d.SetId(assignmentID)
			d.Set("entity_id", entityID)
			d.Set("entity_type", entityType)
			d.Set("tag_id", tagID)
			foundTag = tag
			break
		}
	}
	if foundTag == nil {
		d.SetId("")
		return []*schema.ResourceData{}, fmt.Errorf("error importing pagerduty_tag_assignment. Tag not found for entity")
	}

	return []*schema.ResourceData{d}, err
}

func isFoundTagAssignmentEntity(entityID, entityType string, meta interface{}) (bool, error) {
	var isFound bool
	client, err := meta.(*Config).Client()
	if err != nil {
		return isFound, err
	}

	fetchUser := func(id string) (*pagerduty.User, *pagerduty.Response, error) {
		return client.Users.Get(id, &pagerduty.GetUserOptions{})
	}
	fetchTeam := func(id string) (*pagerduty.Team, *pagerduty.Response, error) {
		return client.Teams.Get(id)
	}
	fetchEscalationPolicy := func(id string) (*pagerduty.EscalationPolicy, *pagerduty.Response, error) {
		return client.EscalationPolicies.Get(id, &pagerduty.GetEscalationPolicyOptions{})
	}
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		if entityType == "users" {
			_, _, err = fetchUser(entityID)
		}
		if entityType == "teams" {
			_, _, err = fetchTeam(entityID)
		}
		if entityType == "escalation_policies" {
			_, _, err = fetchEscalationPolicy(entityID)
		}

		if err != nil && isErrCode(err, http.StatusNotFound) {
			return nil
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		if isErrCode(err, http.StatusBadRequest) {
			return resource.NonRetryableError(err)
		}

		isFound = true

		return nil
	})
	if retryErr != nil {
		return isFound, retryErr
	}
	return isFound, nil
}

func createAssignmentID(entityID, tagID string) string {
	return fmt.Sprintf("%v.%v", entityID, tagID)
}
