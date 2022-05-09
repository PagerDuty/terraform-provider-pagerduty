package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationPathRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyEventOrchestrationPathRouterCreate,
		Read:   resourcePagerDutyEventOrchestrationPathRouterRead,
		Update: resourcePagerDutyEventOrchestrationPathRouterUpdate,
		Delete: resourcePagerDutyEventOrchestrationPathRouterUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //TODO: resourcePagerDutyEventOrchestrationPathImport
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "parent": {
			// 	Type:     schema.TypeList,
			// 	Required: true,
			// 	MaxItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"type": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"self": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// },
			// "sets": {
			// 	Type:     schema.TypeList,
			// 	Required: true, //TODO: is it always going to have a set?
			// 	MaxItems: 1,    // Router can only have 'start' set
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"rules": {
			// 				Type:     schema.TypeList,
			// 				Required: true, // even if there are no rules, API returns rules as an empty list
			// 				MaxItems: 1000, // TODO: do we need this?! Router allows a max of 1000 rules
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"id": {
			// 							Type:     schema.TypeString,
			// 							Optional: true, // If the start set has no rules, empty list is returned by API for rules. TODO: there is a validation on id
			// 						},
			// 						"label": {
			// 							Type:     schema.TypeString,
			// 							Optional: true,
			// 						},
			// 						"conditions": {
			// 							Type:     schema.TypeList,
			// 							Optional: true,
			// 							Elem: &schema.Resource{
			// 								Schema: map[string]*schema.Schema{
			// 									"expression": {
			// 										Type:     schema.TypeString,
			// 										Required: true,
			// 									},
			// 								},
			// 							},
			// 						},
			// 						"actions": {
			// 							Type:     schema.TypeList,
			// 							Optional: true,
			// 							MaxItems: 1, //there can only be one action for router
			// 							Elem: &schema.Resource{
			// 								Schema: map[string]*schema.Schema{
			// 									"route_to": {
			// 										Type:     schema.TypeString,
			// 										Required: true,
			// 										//TODO: validate func, cannot be unrouted, should be some serviceID
			// 									},
			// 								},
			// 							},
			// 						},
			// 						"disabled": {
			// 							Type:     schema.TypeBool,
			// 							Optional: true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			// "catch_all": {
			// 	Type:     schema.TypeList,
			// 	Optional: true, //if not supplied, API creates it
			// 	MaxItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"actions": {
			// 				Type:     schema.TypeList,
			// 				Optional: true, //if not provided, API defaults to unrouted
			// 				MaxItems: 1,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"route_to": {
			// 							Type:     schema.TypeString,
			// 							Optional: true, //if not provided, API defaults to unrouted
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

// func buildEventOrchestrationRouterStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
// 	orchestrationPath := &pagerduty.EventOrchestrationPath{
// 		Type: d.Get("type").(string),
// 		//Self: d.Get("self").(string),
// 	}

// 	// if attr, ok := d.GetOk("description"); ok {
// 	// 	orchestration.Description = attr.(string)
// 	// }

// 	// if attr, ok := d.GetOk("team"); ok {
// 	// 	orchestration.Team = expandOrchestrationTeam(attr)
// 	// }
// 	return orchestrationPath
// }

func resourcePagerDutyEventOrchestrationPathRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	// TODO: figure out a way to get the Type differently
	log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "router", d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		if routerPath, _, err := client.EventOrchestrationPaths.Get(d.Id(), "router"); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if routerPath != nil {
			d.Set("type", routerPath.Type)
		}
		return nil
	})

}

//TODO: As a temporary fix to get rid of "Inconistency issue - root resource created but not present", made the create to have same logic as read.
func resourcePagerDutyEventOrchestrationPathRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
	routerPath := buildRouterPathStruct(d)
	// TODO: figure out a way to get the Type differently
	log.Printf("[INFO] Updating PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "router", d.Id())

	return performRouterPathUpdate(d.Id(), routerPath, client)
}

func performRouterPathUpdate(orchestrationID string, routerPath *pagerduty.EventOrchestrationPath, client *pagerduty.Client) error {
	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if _, _, err := client.EventOrchestrationPaths.Update(orchestrationID, "router", routerPath); err != nil {
			return resource.RetryableError(err)
		}
		//TODO: figure out rule ordering
		// else if rule.Position != nil && *updatedRouterPath.Position != *rule.Position && rule.CatchAll != true {
		// 	log.Printf("[INFO] PagerDuty ruleset rule %s position %d needs to be %d", updatedRouterPath.ID, *updatedRouterPath.Position, *rule.Position)
		// 	return resource.RetryableError(fmt.Errorf("Error updating ruleset rule %s position %d needs to be %d", updatedRouterPath.ID, *updatedRouterPath.Position, *rule.Position))
		// }
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return nil
}

func resourcePagerDutyEventOrchestrationPathRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePagerDutyEventOrchestrationPathRouterDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildRouterPathStruct(d *schema.ResourceData) *pagerduty.EventOrchestrationPath {
	orchPath := &pagerduty.EventOrchestrationPath{
		Type: d.Get("type").(string),
	}

	if attr, ok := d.GetOk("Parent"); ok {
		orchPath.Parent = expandOrchestrationPathParent(attr)
	}

	return orchPath
}

func expandOrchestrationPathParent(v interface{}) *pagerduty.EventOrchestrationPathReference {
	var parent *pagerduty.EventOrchestrationPathReference
	p := v.([]interface{})[0].(map[string]interface{})
	parent = &pagerduty.EventOrchestrationPathReference{
		ID:   p["id"].(string),
		Type: p["type"].(string),
		Self: p["self"].(string),
	}

	return parent
}

// func buildEventOrchestrationStruct(d *schema.ResourceData) *pagerduty.EventOrchestration {
// 	orchestration := &pagerduty.EventOrchestration{
// 		Name: d.Get("name").(string),
// 	}

// 	if attr, ok := d.GetOk("description"); ok {
// 		orchestration.Description = attr.(string)
// 	}

// 	if attr, ok := d.GetOk("team"); ok {
// 		orchestration.Team = expandOrchestrationTeam(attr)
// 	}

// 	return orchestration
// }

// func expandOrchestrationTeam(v interface{}) *pagerduty.EventOrchestrationObject {
// 	var team *pagerduty.EventOrchestrationObject
// 	t := v.([]interface{})[0].(map[string]interface{})
// 	team = &pagerduty.EventOrchestrationObject{
// 		ID: t["id"].(string),
// 	}

// 	return team
// }

// func resourcePagerDutyEventOrchestrationCreate(d *schema.ResourceData, meta interface{}) error {
// 	client, err := meta.(*Config).Client()
// 	if err != nil {
// 		return err
// 	}

// 	payload := buildEventOrchestrationStruct(d)
// 	var orchestration *pagerduty.EventOrchestration

// 	log.Printf("[INFO] Creating PagerDuty Event Orchestration: %s", payload.Name)

// 	retryErr := resource.Retry(10*time.Second, func() *resource.RetryError {
// 		if orch, _, err := client.EventOrchestrations.Create(payload); err != nil {
// 			if isErrCode(err, 400) || isErrCode(err, 429) {
// 				return resource.RetryableError(err)
// 			}

// 			return resource.NonRetryableError(err)
// 		} else if orch != nil {
// 			d.SetId(orch.ID)
// 			orchestration = orch
// 		}
// 		return nil
// 	})

// 	if retryErr != nil {
// 		return retryErr
// 	}

// 	setEventOrchestrationProps(d, orchestration)

// 	return nil
// }

// func resourcePagerDutyEventOrchestrationRead(d *schema.ResourceData, meta interface{}) error {
// 	client, err := meta.(*Config).Client()
// 	if err != nil {
// 		return err
// 	}

// 	return resource.Retry(2*time.Minute, func() *resource.RetryError {
// 		orch, _, err := client.EventOrchestrations.Get(d.Id())
// 		if err != nil {
// 			errResp := handleNotFoundError(err, d)
// 			if errResp != nil {
// 				time.Sleep(2 * time.Second)
// 				return resource.RetryableError(errResp)
// 			}

// 			return nil
// 		}

// 		setEventOrchestrationProps(d, orch)

// 		return nil
// 	})
// }

// func resourcePagerDutyEventOrchestrationUpdate(d *schema.ResourceData, meta interface{}) error {
// 	client, err := meta.(*Config).Client()
// 	if err != nil {
// 		return err
// 	}

// 	orchestration := buildEventOrchestrationStruct(d)

// 	log.Printf("[INFO] Updating PagerDuty Event Orchestration: %s", d.Id())

// 	if _, _, err := client.EventOrchestrations.Update(d.Id(), orchestration); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourcePagerDutyEventOrchestrationDelete(d *schema.ResourceData, meta interface{}) error {
// 	client, err := meta.(*Config).Client()
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("[INFO] Deleting PagerDuty Event Orchestration: %s", d.Id())
// 	if _, err := client.EventOrchestrations.Delete(d.Id()); err != nil {
// 		return err
// 	}

// 	d.SetId("")

// 	return nil
// }

// func flattenEventOrchestrationTeam(v *pagerduty.EventOrchestrationObject) []interface{} {
// 	team := map[string]interface{}{
// 		"id": v.ID,
// 	}

// 	return []interface{}{team}
// }

// func flattenEventOrchestrationIntegrations(eoi []*pagerduty.EventOrchestrationIntegration) []interface{} {
// 	var result []interface{}

// 	for _, i := range eoi {
// 		integration := map[string]interface{}{
// 			"id":   i.ID,
// 			"parameters": flattenEventOrchestrationIntegrationParameters(i.Parameters),
// 		}
// 		result = append(result, integration)
// 	}
// 	return result
// }

// func flattenEventOrchestrationIntegrationParameters(p *pagerduty.EventOrchestrationIntegrationParameters) []interface{} {
// 	result := map[string]interface{}{
// 		"routing_key": p.RoutingKey,
// 		"type": p.Type,
// 	}

// 	return []interface{}{result}
// }

// func setEventOrchestrationProps(d *schema.ResourceData, o *pagerduty.EventOrchestration) error {
// 	d.Set("name", o.Name)
// 	d.Set("description", o.Description)
// 	d.Set("routes", o.Routes)

// 	if o.Team != nil {
// 		d.Set("team", flattenEventOrchestrationTeam(o.Team))
// 	}

//   if len(o.Integrations) > 0 {
// 		d.Set("integrations", flattenEventOrchestrationIntegrations(o.Integrations))
// 	}

// 	return nil
// }
