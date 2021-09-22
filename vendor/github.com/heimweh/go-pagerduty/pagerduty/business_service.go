package pagerduty

import "fmt"

// BusinessServiceService handles the communication with business service
// related methods of the PagerDuty API.
type BusinessServiceService service

// BusinessService represents a business service.
type BusinessService struct {
	ID             string               `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Type           string               `json:"type,omitempty"`
	Summary        string               `json:"summary,omitempty"`
	Self           string               `json:"self,omitempty"`
	PointOfContact string               `json:"point_of_contact,omitempty"`
	HTMLUrl        string               `json:"html_url,omitempty"`
	Description    string               `json:"description,omitempty"`
	Team           *BusinessServiceTeam `json:"team,omitempty"`
}

// BusinessServiceTeam represents a team object in a business service
type BusinessServiceTeam struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Self string `json:"self,omitempty"`
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

// List lists existing business services.
func (s *BusinessServiceService) List() (*ListBusinessServicesResponse, *Response, error) {
	u := "/business_services"
	v := new(ListBusinessServicesResponse)

	businessServices := make([]*BusinessService, 0)

	// Create a handler closure capable of parsing data from the business_services endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListBusinessServicesResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		businessServices = append(businessServices, result.BusinessServices...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(u, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.BusinessServices = businessServices

	return v, nil, nil
}

// Create creates a new business service.
func (s *BusinessServiceService) Create(bservice *BusinessService) (*BusinessService, *Response, error) {
	u := "/business_services"
	v := new(BusinessServicePayload)
	p := &BusinessServicePayload{BusinessService: bservice}

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
func (s *BusinessServiceService) Update(ID string, bserv *BusinessService) (*BusinessService, *Response, error) {
	u := fmt.Sprintf("/business_services/%s", ID)
	v := new(BusinessServicePayload)
	p := BusinessServicePayload{BusinessService: bserv}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.BusinessService, resp, nil
}
