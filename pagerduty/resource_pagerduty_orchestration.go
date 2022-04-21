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
	t := v.([]interface{})[0].(map[string]interface{})
	team := *pagerduty.OrchestrationObject{
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
		if orchestration, _, err := client.Rulesets.Create(orchestration); err != nil {
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

	return orchestration
}
