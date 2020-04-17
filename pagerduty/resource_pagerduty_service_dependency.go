package pagerduty

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyServiceDependency() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyServiceDependencyAssociate,
		Read:   resourcePagerDutyServiceDependencyRead,
		Delete: resourcePagerDutyServiceDependencyDisassociate,
		Importer: &schema.ResourceImporter{
			State: resourcePagerDutyServiceDependencyImport,
		},
		Schema: map[string]*schema.Schema{
			"dependency": {
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

func buildServiceDependencyStruct(d *schema.ResourceData) (*pagerduty.ServiceDependency, error) {
	var rel *pagerduty.ServiceDependency
	rel = new(pagerduty.ServiceDependency)

	log.Printf("[BUGIN] buildServiceDepStruct-dependency: %v", d.Get("dependency"))
	for _, r := range d.Get("dependency").([]interface{}) {
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
func resourcePagerDutyServiceDependencyAssociate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	serviceDependency, err := buildServiceDependencyStruct(d)
	if err != nil {
		return err
	}
	var r []*pagerduty.ServiceDependency
	r = append(r, serviceDependency)

	input := pagerduty.ListServiceDependencies{
		Relationships: r,
	}
	log.Printf("[INFO] Associating PagerDuty dependency %s", serviceDependency.ID)

	var dependencies *pagerduty.ListServiceDependencies
	retryErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		if dependencies, _, err = client.ServiceDependencies.AssociateServiceDependencies(&input); err != nil {
			if isErrCode(err, 404) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		} else {
			for _, r := range dependencies.Relationships {
				d.SetId(r.ID)
			}
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}
	return resourcePagerDutyServiceDependencyRead(d, meta)
}

func resourcePagerDutyServiceDependencyDisassociate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	dependency, err := buildServiceDependencyStruct(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Disassociating PagerDuty dependency %s", dependency.DependentService.ID)

	// listServiceRelationships by calling get dependencies using the serviceDependency.DependentService.ID
	depResp, _, err := client.ServiceDependencies.GetBusinessServiceDependencies(dependency.DependentService.ID)
	if err != nil {
		return err
	}

	var foundDep *pagerduty.ServiceDependency

	// loop serviceRelationships until relationship.IDs match
	for _, rel := range depResp.Relationships {
		if rel.ID == d.Id() {
			foundDep = rel
			break
		}
	}
	// check if relationship not found
	if foundDep == nil {
		d.SetId("")
		return nil
	}
	// convertType is needed because the PagerDuty API returns the 'reference' values in responses but wants the other
	// values in requests
	foundDep.SupportingService.Type = convertType(foundDep.SupportingService.Type)
	foundDep.DependentService.Type = convertType(foundDep.DependentService.Type)

	// set matching Dependency to r
	var r []*pagerduty.ServiceDependency

	r = append(r, foundDep)

	input := pagerduty.ListServiceDependencies{
		Relationships: r,
	}
	_, _, err = client.ServiceDependencies.DisassociateServiceDependencies(&input)
	if err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyServiceDependencyRead(d *schema.ResourceData, meta interface{}) error {
	serviceDependency, err := buildServiceDependencyStruct(d)
	log.Printf("[INFO] Reading PagerDuty dependency %s", serviceDependency.ID)

	if err = findDependencySetState(d.Id(), serviceDependency.DependentService.ID, d, meta); err != nil {
		return err
	}

	log.Printf("[DEBUGGIN] %s", d.Get("dependency"))

	return nil
}

func flattenRelationship(r *pagerduty.ServiceDependency) map[string]interface{} {
	relationship := map[string]interface{}{
		"supporting_service": flattenService(r.SupportingService),
		"dependent_service":  flattenService(r.DependentService),
	}
	log.Printf("[flattenRelationship] relationship: %v", relationship)

	return relationship
}

func flattenService(s *pagerduty.ServiceObj) map[string]interface{} {
	service := map[string]interface{}{
		"id":   s.ID,
		"type": convertType(s.Type),
	}
	return service
}

// convertType is needed because the PagerDuty API returns the 'reference' values in responses but wants the other
// values in requests
func convertType(s string) string {
	switch s {
	case "business_service_reference":
		s = "business_service"
	case "technical_service_reference":
		s = "service"
	}
	return s
}

func findDependencySetState(depID, busServiceID string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	// Pausing to let the PD API sync.
	time.Sleep(2 * time.Second)
	dependencies, _, err := client.ServiceDependencies.GetBusinessServiceDependencies(busServiceID)

	if err != nil {
		return err
	}
	for _, rel := range dependencies.Relationships {
		if rel.ID == depID {
			d.SetId(rel.ID)
			d.Set("dependency", flattenRelationship(rel))
			log.Printf("[findDependencySetState] dependency: %v", d.Get("dependency"))
			break
		}
	}

	return nil
}

func resourcePagerDutyServiceDependencyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), ".")

	if len(ids) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_service_dependency. Expecting an importation ID formed as '<business_service_id>.<service_dependency_id>'")
	}
	sid, id := ids[0], ids[1]

	if err := findDependencySetState(id, sid, d, meta); err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
