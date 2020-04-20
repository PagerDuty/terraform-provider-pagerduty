package pagerduty

import "fmt"

// BusinessServiceService handles the communication with business service
// related methods of the PagerDuty API.
type BusinessServiceService service

// BusinessService represents a business service.
type BusinessService struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Type           string `json:"type,omitempty"`
	Summary        string `json:"summary,omitempty"`
	Self           string `json:"self,omitempty"`
	PointOfContact string `json:"point_of_contact,omitempty"`
	HTMLUrl        string `json:"html_url,omitempty"`
	Description    string `json:"description,omitempty"`
}

// BusinessServicePayload represents payload with a business service object
type BusinessServicePayload struct {
	BusinessService *BusinessService `json:"business_service,omitempty"`
}

// ListBusinessServicesResponse represents a list response of business services.
type ListBusinessServicesResponse struct {
	Total            int                `json:"total,omitempty"`
	BusinessServices []*BusinessService `json:"business_services,omitempty"`
	Offset           int                `json:"offset,omitempty"`
	More             bool               `json:"more,omitempty"`
	Limit            int                `json:"limit,omitempty"`
}

// ServiceRelationship represents a relationship between a business and technical service
type ServiceRelationship struct {
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

// ListServiceRelationships represents a list of relationships for a business service
type ListServiceRelationships struct {
	Relationships []*ServiceRelationship `json:"relationships,omitempty"`
}

// List lists existing business services.
func (s *BusinessServiceService) List() (*ListBusinessServicesResponse, *Response, error) {
	u := "/business_services"
	v := new(ListBusinessServicesResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new business service.
func (s *BusinessServiceService) Create(ruleset *BusinessService) (*BusinessService, *Response, error) {
	u := "/business_services"
	v := new(BusinessServicePayload)
	p := &BusinessServicePayload{BusinessService: ruleset}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.BusinessService, resp, nil
}

// Get gets a business service.
func (s *BusinessServiceService) Get(ID string) (*BusinessService, *Response, error) {
	u := fmt.Sprintf("/business_services/%s", ID)
	v := new(BusinessServicePayload)
	p := &BusinessServicePayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.BusinessService, resp, nil
}

// Delete deletes a business service.
func (s *BusinessServiceService) Delete(ID string) (*Response, error) {
	u := fmt.Sprintf("/business_services/%s", ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Update updates a business service.
func (s *BusinessServiceService) Update(ID string, ruleset *BusinessService) (*BusinessService, *Response, error) {
	u := fmt.Sprintf("/business_services/%s", ID)
	v := new(BusinessServicePayload)
	p := BusinessServicePayload{BusinessService: ruleset}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.BusinessService, resp, nil
}

// AssociateServiceDependencies Create new dependencies between two services
func (s *BusinessServiceService) AssociateServiceDependencies(relationships *ListServiceRelationships) (*Response, error) {
	u := "/service_dependencies/associate"

	return s.client.newRequestDo("POST", u, nil, relationships, nil)
}

// DisassociateServiceDependencies Disassociate dependencies between two services.
func (s *BusinessServiceService) DisassociateServiceDependencies(relationships *ListServiceRelationships) (*Response, error) {
	u := "/service_dependencies/disassociate"

	return s.client.newRequestDo("POST", u, nil, relationships, nil)
}

// GetDependencies gets all immediate dependencies of a business service.
func (s *BusinessServiceService) GetDependencies(businessServiceID string) (*ListServiceRelationships, *Response, error) {
	u := fmt.Sprintf("/service_dependencies/business_services/%s", businessServiceID)
	v := new(ListServiceRelationships)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
