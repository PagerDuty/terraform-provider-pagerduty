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
		Read:   resourcePagerDutyEventOrchestrationPathRouterRead,
		Create: resourcePagerDutyEventOrchestrationPathRouterCreate,
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
			"parent": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"self": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
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

	log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "router", d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		path := buildRouterPathStruct(d)
		if routerPath, _, err := client.EventOrchestrationPaths.Get(path.Parent.ID, path.Type); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if routerPath != nil {
			d.SetId(path.Parent.ID)
			d.Set("type", routerPath.Type)
		}
		return nil
	})

}

func resourcePagerDutyEventOrchestrationPathRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PagerDuty Event Orchestration Path of type: %s for orchestration: %s", "router", d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		routerPathStruct := buildRouterPathStruct(d)
		if routerPath, _, err := client.EventOrchestrationPaths.Get(routerPathStruct.Parent.ID, "router"); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if routerPath != nil {
			d.SetId(routerPathStruct.Parent.ID)
			d.Set("type", routerPath.Type)
		}
		return nil
	})

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

	if attr, ok := d.GetOk("parent"); ok {
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
