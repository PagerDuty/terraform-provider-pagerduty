package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyBusinessServiceDependency() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyBusinessServiceDependencyAssociate,
		Read:   resourcePagerDutyBusinessServiceDependencyRead,
		Delete: resourcePagerDutyBusinessServiceDependencyDisassociate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"relationship": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"supporting_service": {
							Required: true,
							Type:     schema.TypeList,
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
								},
							},
						},
						"dependent_service": {
							Required: true,
							Type:     schema.TypeList,
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
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "service_dependency",
						},
					},
				},
			},
		},
	}
}

func buildBusinessServiceDependencyStruct(d *schema.ResourceData) (*pagerduty.ServiceRelationship, error) {
	var rel *pagerduty.ServiceRelationship

	rel = new(pagerduty.ServiceRelationship)

	for _, r := range d.Get("relationship").([]interface{}) {
		relmap := r.(map[string]interface{})
		rel.SupportingService = expandService(relmap["supporting_service"].(interface{}))
		rel.DependentService = expandService(relmap["dependent_service"].(interface{}))
	}
	if attr, ok := d.GetOk("type"); ok {
		rel.Type = attr.(string)
	}
	return rel, nil
}

func expandService(v interface{}) *pagerduty.ServiceObj {
	var so *pagerduty.ServiceObj
	so = new(pagerduty.ServiceObj)

	for _, s := range v.([]interface{}) {
		sm := s.(map[string]interface{})

		so.ID = sm["id"].(string)
		so.Type = sm["type"].(string)
	}

	return so
}
func resourcePagerDutyBusinessServiceDependencyAssociate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	serviceDependency, err := buildBusinessServiceDependencyStruct(d)
	if err != nil {
		return err
	}
	var r []*pagerduty.ServiceRelationship
	r = append(r, serviceDependency)

	dependencies := *&pagerduty.ListServiceRelationships{
		Relationships: r,
	}
	log.Printf("[INFO] Associating PagerDuty dependency %s", serviceDependency.ID)

	_, err = client.BusinessServices.AssociateServiceDependencies(&dependencies)
	if err != nil {
		return err
	}
	// API doesn't return a response, so we're creating the ID with the same pattern the API uses
	d.SetId(fmt.Sprintf("D-%s-%s", serviceDependency.DependentService.ID, serviceDependency.SupportingService.ID))

	return resourcePagerDutyBusinessServiceDependencyRead(d, meta)
}

func resourcePagerDutyBusinessServiceDependencyDisassociate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	splitID := strings.Split(d.Id(), "-")
	var bServeID string

	if len(splitID) > 1 {
		bServeID = splitID[1]
	} else {
		return fmt.Errorf("No business service id found in service relationship ID: %s", d.Id())
	}

	log.Printf("[INFO] Disassociating PagerDuty dependency %s", d.Id())

	// listServiceRelationships by calling get dependencies using the serviceDependency.DependentService.ID
	depResp, _, err := client.BusinessServices.GetDependencies(bServeID)
	if err != nil {
		return err
	}

	var foundRel *pagerduty.ServiceRelationship

	// // loop serviceRelationships until relationship.IDs match
	for _, rel := range depResp.Relationships {
		if rel.ID == d.Id() {
			foundRel = rel
			break
		}
	}
	// // check if relationship not found
	if foundRel == nil {
		d.SetId("")
		return nil
	}

	// // set matching Relationship to r
	var r []*pagerduty.ServiceRelationship

	r = append(r, foundRel)

	dependencies := *&pagerduty.ListServiceRelationships{
		Relationships: r,
	}
	_, err = client.BusinessServices.DisassociateServiceDependencies(&dependencies)
	if err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyBusinessServiceDependencyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	serviceDependency, err := buildBusinessServiceDependencyStruct(d)
	log.Printf("[INFO] Reading PagerDuty dependency %s", serviceDependency.ID)

	// Pausing to let the PD API sync.
	time.Sleep(2 * time.Second)
	dependencies, _, err := client.BusinessServices.GetDependencies(serviceDependency.DependentService.ID)

	var foundRel *pagerduty.ServiceRelationship

	if err != nil {
		return err
	}
	for _, rel := range dependencies.Relationships {
		if rel.ID == serviceDependency.ID {
			foundRel = rel
			break
		}
	}
	if foundRel != nil {
		d.Set("relationship", flattenRelationship(foundRel))
	}

	log.Printf("[DEBUGGIN] %s", d.Get("id"))

	return nil
}

func flattenRelationship(r *pagerduty.ServiceRelationship) map[string]interface{} {
	relationship := map[string]interface{}{
		"supporting_service": flattenService(r.SupportingService),
		"dependent_service":  flattenService(r.DependentService),
	}
	return relationship
}

func flattenService(s *pagerduty.ServiceObj) map[string]interface{} {
	service := map[string]interface{}{
		"id":   s.ID,
		"type": s.Type,
	}

	return service
}
