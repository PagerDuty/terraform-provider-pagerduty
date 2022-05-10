package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestration() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyEventOrchestrationCreate,
		Read:   resourcePagerDutyEventOrchestrationRead,
		Update: resourcePagerDutyEventOrchestrationUpdate,
		Delete: resourcePagerDutyEventOrchestrationDelete,
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
			"routes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"integrations": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true, // Tests keep failing if "Optional: true" is not provided
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"routing_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
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

func buildEventOrchestrationStruct(d *schema.ResourceData) *pagerduty.EventOrchestration {
	orchestration := &pagerduty.EventOrchestration{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("description"); ok {
		orchestration.Description = attr.(string)
	}

	if attr, ok := d.GetOk("team"); ok {
		orchestration.Team = expandOrchestrationTeam(attr)
	} else {
		var tId *string
		orchestration.Team = &pagerduty.EventOrchestrationObject{
			ID: tId,
		}
	}

	return orchestration
}

func expandOrchestrationTeam(v interface{}) *pagerduty.EventOrchestrationObject {
	var team *pagerduty.EventOrchestrationObject
	t := v.([]interface{})[0].(map[string]interface{})
	team = &pagerduty.EventOrchestrationObject{
		ID: stringTypeToStringPtr(t["id"].(string)),
	}

	return team
}

func resourcePagerDutyEventOrchestrationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	payload := buildEventOrchestrationStruct(d)
	var orchestration *pagerduty.EventOrchestration

	log.Printf("[INFO] Creating PagerDuty Event Orchestration: %s", payload.Name)

	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
		if orch, _, err := client.EventOrchestrations.Create(payload); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if orch != nil {
			d.SetId(orch.ID)
			orchestration = orch
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	setEventOrchestrationProps(d, orchestration)

	return nil
}

func resourcePagerDutyEventOrchestrationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		orch, _, err := client.EventOrchestrations.Get(d.Id())
		if err != nil {
			errResp := handleNotFoundError(err, d)
			if errResp != nil {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(errResp)
			}

			return nil
		}

		setEventOrchestrationProps(d, orch)

		return nil
	})
}

func resourcePagerDutyEventOrchestrationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orchestration := buildEventOrchestrationStruct(d)

	log.Printf("[INFO] Updating PagerDuty Event Orchestration: %s", d.Id())

	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
		if _, _, err := client.EventOrchestrations.Update(d.Id(), orchestration); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	return nil
}

func resourcePagerDutyEventOrchestrationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PagerDuty Event Orchestration: %s", d.Id())
	if _, err := client.EventOrchestrations.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func flattenEventOrchestrationTeam(v *pagerduty.EventOrchestrationObject) []interface{} {
	team := map[string]interface{}{
		"id": v.ID,
	}

	return []interface{}{team}
}

func flattenEventOrchestrationIntegrations(eoi []*pagerduty.EventOrchestrationIntegration) []interface{} {
	var result []interface{}

	for _, i := range eoi {
		integration := map[string]interface{}{
			"id":         i.ID,
			"parameters": flattenEventOrchestrationIntegrationParameters(i.Parameters),
		}
		result = append(result, integration)
	}
	return result
}

func flattenEventOrchestrationIntegrationParameters(p *pagerduty.EventOrchestrationIntegrationParameters) []interface{} {
	result := map[string]interface{}{
		"routing_key": p.RoutingKey,
		"type":        p.Type,
	}

	return []interface{}{result}
}

func setEventOrchestrationProps(d *schema.ResourceData, o *pagerduty.EventOrchestration) error {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("routes", o.Routes)

	if o.Team != nil {
		d.Set("team", flattenEventOrchestrationTeam(o.Team))
	}

	if len(o.Integrations) > 0 {
		d.Set("integrations", flattenEventOrchestrationIntegrations(o.Integrations))
	}

	return nil
}
