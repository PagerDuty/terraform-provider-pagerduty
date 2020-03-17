package pagerduty

import (
	"log"

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
	log.Printf("[INFO] Associating PagerDuty dependency between business service %s and service %s", serviceDependency.SupportingService.ID, serviceDependency.DependentService.ID)

	_, err = client.BusinessServices.AssociateServiceDependencies(&dependencies)
	if err != nil {
		return err
	}

	return resourcePagerDutyBusinessServiceDependencyRead(d, meta)
}

func resourcePagerDutyBusinessServiceDependencyDisassociate(d *schema.ResourceData, meta interface{}) error {
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
	log.Printf("[INFO] Disassociating PagerDuty dependency between business service %s and service %s", serviceDependency.SupportingService.ID, serviceDependency.DependentService.ID)

	_, err = client.BusinessServices.DisassociateServiceDependencies(&dependencies)
	if err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyBusinessServiceDependencyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	serviceDependency, err := buildBusinessServiceDependencyStruct(d)
	// log.Printf("[INFO] Reading PagerDuty dependency between business service %s and service %s", serviceDependency.SupportingService.ID, serviceDependency.DependentService.ID)

	dependencies, _, err := client.BusinessServices.GetDependencies(serviceDependency.SupportingService.ID)
	var foundDep *pagerduty.ServiceRelationship

	if err != nil {
		return err
	}
	for _, rel := range dependencies.Relationships {
		if rel.DependentService.ID == serviceDependency.DependentService.ID {
			// log.Printf("[DEBUG] FoundDep.SupportingService: %s", rel.SupportingService.ID)
			foundDep = rel
			break
		}
	}
	if foundDep != nil {
		d.Set("supporting_service", flattenService(foundDep.SupportingService))
		d.Set("dependent_service", flattenService(foundDep.DependentService))
	}
	return nil
}

func flattenService(s *pagerduty.ServiceObj) map[string]interface{} {
	service := map[string]interface{}{
		"id":   s.ID,
		"type": s.Type,
	}
	return service
}
