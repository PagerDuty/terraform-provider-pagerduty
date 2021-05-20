package pagerduty

import "fmt"

// ServiceDependencyService handles the communication with service dependency
// related methods of the PagerDuty API.
type ServiceDependencyService service

// ServiceDependency represents a relationship between a business and technical service
type ServiceDependency struct {
	ID                string      `json:"id,omitempty"`
	Type              string      `json:"type,omitempty"`
	SupportingService *ServiceObj `json:"supporting_service,omitempty"`
	DependentService  *ServiceObj `json:"dependent_service,omitempty"`
}

// ServiceObj represents a service object in service relationship
type ServiceObj struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// ListServiceDependencies represents a list of dependencies for a service
type ListServiceDependencies struct {
	Relationships []*ServiceDependency `json:"relationships,omitempty"`
}

// AssociateServiceDependencies Create new dependencies between two services
func (s *ServiceDependencyService) AssociateServiceDependencies(dependencies *ListServiceDependencies) (*ListServiceDependencies, *Response, error) {
	u := "/service_dependencies/associate"
	v := new(ListServiceDependencies)

	resp, err := s.client.newRequestDo("POST", u, nil, dependencies, &v)

	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// DisassociateServiceDependencies Disassociate dependencies between two services.
func (s *ServiceDependencyService) DisassociateServiceDependencies(dependencies *ListServiceDependencies) (*ListServiceDependencies, *Response, error) {
	u := "/service_dependencies/disassociate"
	v := new(ListServiceDependencies)

	resp, err := s.client.newRequestDo("POST", u, nil, dependencies, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// GetServiceDependenciesForType gets all immediate dependencies of a dependent service.
func (s *ServiceDependencyService) GetServiceDependenciesForType(serviceID, serviceType string) (*ListServiceDependencies, *Response, error) {
	if serviceType == "business_service" || serviceType == "business_service_reference" {
		return s.getBusinessServiceDependencies(serviceID)
	} else if serviceType == "service" || serviceType == "technical_service_reference" {
		return s.getTechnicalServiceDependencies(serviceID)
	}
	// return a not found error
	return nil, nil, fmt.Errorf("dependent Service type of %s not found", serviceType)
}

// getBusinessServiceDependencies gets all immediate dependencies of a business service.
func (s *ServiceDependencyService) getBusinessServiceDependencies(businessServiceID string) (*ListServiceDependencies, *Response, error) {
	u := fmt.Sprintf("/service_dependencies/business_services/%s", businessServiceID)
	v := new(ListServiceDependencies)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// getTechnicalServiceDependencies gets all immediate dependencies of a technical service.
func (s *ServiceDependencyService) getTechnicalServiceDependencies(serviceID string) (*ListServiceDependencies, *Response, error) {
	u := fmt.Sprintf("/service_dependencies/technical_services/%s", serviceID)
	v := new(ListServiceDependencies)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
