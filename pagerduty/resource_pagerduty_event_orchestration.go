package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyOrchestration() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyOrchestrationCreate,
		Read: resourcePagerDutyOrchestrationRead,
		Update: resourcePagerDutyOrchestrationUpdate,
		Delete: resourcePagerDutyOrchestrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildOrchestrationStruct(d *schema.ResourceData) *pagerduty.Orchestration {
	orchestration := &pagerduty.Orchestration{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("description"); ok {
		orchestration.Description = attr.(string)
	}

	if attr, ok := d.GetOk("team"); ok {
		orchestration.Team = expandOrchestrationTeam(attr)
	}

	return orchestration
}

// TODO why is "team" a list?
func expandOrchestrationTeam(v interface{}) *pagerduty.OrchestrationObject {
	var team *pagerduty.OrchestrationObject
	t := v.([]interface{})[0].(map[string]interface{})
	team = &pagerduty.OrchestrationObject{
		ID: t["id"].(string),
	}

	return team
}

func resourcePagerDutyOrchestrationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orchestration := buildOrchestrationStruct(d)

	log.Printf("[INFO] Creating PagerDuty orchestration: %s", orchestration.Name)

	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
		if orchestration, _, err := client.Orchestrations.Create(orchestration); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if orchestration != nil {
			d.SetId(orchestration.ID)
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	setOrchestrationProps(d, orchestration)

	return nil
}

func resourcePagerDutyOrchestrationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		orch, _, err := client.Orchestrations.Get(d.Id())
		if err != nil {
			errResp := handleNotFoundError(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}
		
		setOrchestrationProps(d, orch)

		return nil
	})
}

func resourcePagerDutyOrchestrationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orchestration := buildOrchestrationStruct(d)

	log.Printf("[INFO] Updating PagerDuty orchestration: %s", d.Id())

	if _, _, err := client.Orchestrations.Update(d.Id(), orchestration); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyOrchestrationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty orchestration: %s", d.Id())
	if _, err := client.Orchestrations.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func setOrchestrationProps(d *schema.ResourceData, o *pagerduty.Orchestration) error {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	// TODO: set team, number of routes, integrations if exist
	return nil
}