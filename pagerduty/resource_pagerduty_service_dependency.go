package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

// This mutex was introduced in order to avoid race conditions in the backend
// due to parallel calls to `/service_dependencies/associate` endpoint.
var dependencyAssociationMutex sync.Mutex

func resourcePagerDutyServiceDependency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyServiceDependencyAssociate,
		ReadContext:   resourcePagerDutyServiceDependencyRead,
		DeleteContext: resourcePagerDutyServiceDependencyDisassociate,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePagerDutyServiceDependencyImport,
		},
		Schema: map[string]*schema.Schema{
			"dependency": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
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
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"business_service",
											"service",
										}),
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
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateDiagFunc: validateValueDiagFunc([]string{
											"business_service",
											"business_service_reference",
											"service",
											"technical_service_reference",
										}),
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func buildServiceDependencyStruct(d *schema.ResourceData) (*pagerduty.ServiceDependency, error) {
	rel := new(pagerduty.ServiceDependency)
	rel.ID = d.Id()

	for _, r := range d.Get("dependency").([]interface{}) {
		relmap := r.(map[string]interface{})
		rel.SupportingService = expandService(relmap["supporting_service"])
		rel.DependentService = expandService(relmap["dependent_service"])
	}

	if rel.SupportingService == nil {
		return nil, fmt.Errorf("dependent service not found for dependency: %v", d.Id())
	}
	if rel.DependentService == nil {
		return nil, fmt.Errorf("supporting service not found for dependency: %v", d.Id())
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

func resourcePagerDutyServiceDependencyAssociate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceDependency, err := buildServiceDependencyStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}
	var r []*pagerduty.ServiceDependency
	r = append(r, serviceDependency)

	input := pagerduty.ListServiceDependencies{
		Relationships: r,
	}
	log.Printf("[INFO] Associating PagerDuty dependency %s", serviceDependency.ID)

	var dependencies *pagerduty.ListServiceDependencies
	retryErr := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		// Lock the mutex to ensure only one API call to
		// `service_dependencies/associate` is done at a time
		dependencyAssociationMutex.Lock()
		dependencies, _, err = client.ServiceDependencies.AssociateServiceDependencies(&input)
		dependencyAssociationMutex.Unlock()

		if err != nil {
			if isErrCode(err, 404) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		} else {
			for _, r := range dependencies.Relationships {
				if err := d.Set("dependency", flattenRelationship(r)); err != nil {
					return retry.NonRetryableError(err)
				}
				d.SetId(r.ID)
			}
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	return nil
}

func resourcePagerDutyServiceDependencyDisassociate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	dependency, err := buildServiceDependencyStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Disassociating PagerDuty dependency %s", dependency.DependentService.ID)

	var foundDep *pagerduty.ServiceDependency

	// listServiceRelationships by calling get dependencies using the serviceDependency.DependentService.ID
	retryErr := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		if dependencies, _, err := client.ServiceDependencies.GetServiceDependenciesForType(dependency.DependentService.ID, dependency.DependentService.Type); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			if err = handleNotFoundError(err, d); err == nil {
				return nil
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)

			return retry.RetryableError(err)
		} else if dependencies != nil {
			for _, rel := range dependencies.Relationships {
				if rel.ID == d.Id() {
					foundDep = rel
					break
				}
			}
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(5 * time.Second)
		return diag.FromErr(retryErr)
	}
	// If the dependency is not found, then chances are it had been deleted
	// outside Terraform or be part of a stale state. So it's needed to be cleared
	// from the state.
	if foundDep == nil {
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
	retryErr = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		if _, _, err = client.ServiceDependencies.DisassociateServiceDependencies(&input); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}

			if err = handleNotFoundError(err, d); err == nil {
				return nil
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)

			return retry.RetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(5 * time.Second)
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourcePagerDutyServiceDependencyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	serviceDependency, err := buildServiceDependencyStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Reading PagerDuty dependency %s", serviceDependency.ID)

	if err = findDependencySetState(ctx, d.Id(), serviceDependency.DependentService.ID, serviceDependency.DependentService.Type, d, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenRelationship(r *pagerduty.ServiceDependency) []map[string]interface{} {
	var rels []map[string]interface{}

	relationship := map[string]interface{}{
		"supporting_service": flattenServiceReference(r.SupportingService),
		"dependent_service":  flattenServiceReference(r.DependentService),
	}
	rels = append(rels, relationship)

	return rels
}

func flattenServiceReference(s *pagerduty.ServiceObj) []map[string]interface{} {
	var servs []map[string]interface{}

	service := map[string]interface{}{
		"id":   s.ID,
		"type": convertType(s.Type),
	}
	servs = append(servs, service)
	return servs
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

func findDependencySetState(ctx context.Context, depID, serviceID, serviceType string, d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	// Pausing to let the PD API sync.
	time.Sleep(1 * time.Second)
	retryErr := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		if dependencies, _, err := client.ServiceDependencies.GetServiceDependenciesForType(serviceID, serviceType); err != nil {
			if isErrCode(err, http.StatusBadRequest) {
				return retry.NonRetryableError(err)
			}
			if err = handleNotFoundError(err, d); err == nil {
				return nil
			}

			// Delaying retry by 30s as recommended by PagerDuty
			// https://developer.pagerduty.com/docs/rest-api-v2/rate-limiting/#what-are-possible-workarounds-to-the-events-api-rate-limit
			time.Sleep(30 * time.Second)

			return retry.RetryableError(err)
		} else if dependencies != nil {
			depFound := false
			for _, rel := range dependencies.Relationships {
				if rel.ID == depID {
					if err := d.Set("dependency", flattenRelationship(rel)); err != nil {
						return retry.NonRetryableError(err)
					}
					d.SetId(rel.ID)
					depFound = true
					break
				}
			}
			if !depFound {
				d.SetId("")
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

func resourcePagerDutyServiceDependencyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), ".")

	if len(ids) != 3 {
		return []*schema.ResourceData{}, fmt.Errorf("Error importing pagerduty_service_dependency. Expecting an importation ID formed as '<supporting_service_id>.<supporting_service_type>.<service_dependency_id>'")
	}
	sid, st, id := ids[0], ids[1], ids[2]

	if err := findDependencySetState(ctx, id, sid, st, d, meta); err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
