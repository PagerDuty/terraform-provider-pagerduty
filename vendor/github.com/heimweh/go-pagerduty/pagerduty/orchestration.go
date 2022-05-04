package pagerduty

import (
	"fmt"
)

type OrchestrationService service

type Orchestration struct {
	ID          string               `json:"id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Team        *OrchestrationObject `json:"team,omitempty"`
	// TODO: add Integrations, Routes, Updater, Creator + expand tests to verify these props
}

type OrchestrationObject struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type OrchestrationPayload struct {
	Orchestration *Orchestration `json:"orchestration,omitempty"`
}

var orchestrationBaseUrl = "/event_orchestrations"

func (s *OrchestrationService) Create(orchestration *Orchestration) (*Orchestration, *Response, error) {
	v := new(OrchestrationPayload)
	p := &OrchestrationPayload{Orchestration: orchestration}

	resp, err := s.client.newRequestDo("POST", orchestrationBaseUrl, nil, p, v)

	if err != nil {
		return nil, nil, err
	}

	return v.Orchestration, resp, nil
}

func (s *OrchestrationService) Get(ID string) (*Orchestration, *Response, error) {
	u := fmt.Sprintf("%s/%s", orchestrationBaseUrl, ID)
	v := new(OrchestrationPayload)
	p := &OrchestrationPayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Orchestration, resp, nil
}

func (s *OrchestrationService) Update(ID string, orchestration *Orchestration) (*Orchestration, *Response, error) {
	u := fmt.Sprintf("%s/%s", orchestrationBaseUrl, ID)
	v := new(OrchestrationPayload)
	p := &OrchestrationPayload{Orchestration: orchestration}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Orchestration, resp, nil
}

func (s *OrchestrationService) Delete(ID string) (*Response, error) {
	u := fmt.Sprintf("%s/%s", orchestrationBaseUrl, ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}
