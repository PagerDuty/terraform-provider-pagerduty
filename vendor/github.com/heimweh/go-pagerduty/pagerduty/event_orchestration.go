package pagerduty

import (
	"fmt"
)

type EventOrchestrationService service

type EventOrchestration struct {
	ID           string                           `json:"id,omitempty"`
	Name         string                           `json:"name,omitempty"`
	Description  string                           `json:"description"`
	Team         *EventOrchestrationObject        `json:"team"`
	Routes       int                              `json:"routes,omitempty"`
	Integrations []*EventOrchestrationIntegration `json:"integrations,omitempty"`
}

type EventOrchestrationObject struct {
	Type string  `json:"type,omitempty"`
	ID   *string `json:"id"`
}

type EventOrchestrationPayload struct {
	Orchestration *EventOrchestration `json:"orchestration,omitempty"`
}

type ListEventOrchestrationsResponse struct {
	Total          int                   `json:"total,omitempty"`
	Offset         int                   `json:"offset,omitempty"`
	More           bool                  `json:"more,omitempty"`
	Limit          int                   `json:"limit,omitempty"`
	Orchestrations []*EventOrchestration `json:"orchestrations,omitempty"`
}

var eventOrchestrationBaseUrl = "/event_orchestrations"

func (s *EventOrchestrationService) List() (*ListEventOrchestrationsResponse, *Response, error) {
	v := new(ListEventOrchestrationsResponse)
	v.Total = 0

	orchestrations := make([]*EventOrchestration, 0)

	// Create a handler closure capable of parsing data from the event orchestrations endpoint
	// and appending resultant orchestrations to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListEventOrchestrationsResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		v.Total += result.Total
		v.Offset = result.Offset
		v.More = result.More
		v.Limit = result.Limit
		orchestrations = append(orchestrations, result.Orchestrations...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(eventOrchestrationBaseUrl, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.Orchestrations = orchestrations

	return v, nil, nil
}

func (s *EventOrchestrationService) Create(orchestration *EventOrchestration) (*EventOrchestration, *Response, error) {
	v := new(EventOrchestrationPayload)
	p := &EventOrchestrationPayload{Orchestration: orchestration}

	resp, err := s.client.newRequestDo("POST", eventOrchestrationBaseUrl, nil, p, v)

	if err != nil {
		return nil, nil, err
	}

	return v.Orchestration, resp, nil
}

func (s *EventOrchestrationService) Get(ID string) (*EventOrchestration, *Response, error) {
	u := fmt.Sprintf("%s/%s", eventOrchestrationBaseUrl, ID)
	v := new(EventOrchestrationPayload)
	p := &EventOrchestrationPayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Orchestration, resp, nil
}

func (s *EventOrchestrationService) Update(ID string, orchestration *EventOrchestration) (*EventOrchestration, *Response, error) {
	u := fmt.Sprintf("%s/%s", eventOrchestrationBaseUrl, ID)
	v := new(EventOrchestrationPayload)
	p := &EventOrchestrationPayload{Orchestration: orchestration}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Orchestration, resp, nil
}

func (s *EventOrchestrationService) Delete(ID string) (*Response, error) {
	u := fmt.Sprintf("%s/%s", eventOrchestrationBaseUrl, ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
