package pagerduty

import (
	"context"
	"fmt"
)

type EventOrchestrationIntegrationService service

type EventOrchestrationIntegrationParameters struct {
	RoutingKey string `json:"routing_key,omitempty"`
	Type       string `json:"type,omitempty"`
}

type EventOrchestrationIntegration struct {
	ID         string                                   `json:"id,omitempty"`
	Label      string                                   `json:"label,omitempty"`
	Parameters *EventOrchestrationIntegrationParameters `json:"parameters,omitempty"`
}

type EventOrchestrationIntegrationPayload struct {
	Integration *EventOrchestrationIntegration `json:"integration,omitempty"`
}

type EventOrchestrationIntegrationMigrationPayload struct {
	SourceType    string `json:"source_type,omitempty"`
	SourceId      string `json:"source_id,omitempty"`
	IntegrationId string `json:"integration_id,omitempty"`
}

type ListEventOrchestrationIntegrationsResponse struct {
	Total        int                              `json:"total,omitempty"`
	Integrations []*EventOrchestrationIntegration `json:"integrations,omitempty"`
}

func buildEventOrchestrationIntegrationUrl(orchestrationId string, lastUrlSegment string) string {
	url := fmt.Sprintf("%s/%s/integrations", eventOrchestrationBaseUrl, orchestrationId)

	if len(lastUrlSegment) > 0 {
		url = fmt.Sprintf("%s/%s", url, lastUrlSegment)
	}

	return url
}

func (s *EventOrchestrationIntegrationService) ListContext(ctx context.Context, orchestrationId string) (*ListEventOrchestrationIntegrationsResponse, *Response, error) {
	u := buildEventOrchestrationIntegrationUrl(orchestrationId, "")
	v := new(ListEventOrchestrationIntegrationsResponse)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, v)

	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (s *EventOrchestrationIntegrationService) CreateContext(ctx context.Context, orchestrationId string, integration *EventOrchestrationIntegration) (*EventOrchestrationIntegration, *Response, error) {
	u := buildEventOrchestrationIntegrationUrl(orchestrationId, "")
	v := new(EventOrchestrationIntegrationPayload)
	p := &EventOrchestrationIntegrationPayload{Integration: integration}

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, p, v)

	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

func (s *EventOrchestrationIntegrationService) GetContext(ctx context.Context, orchestrationId string, id string) (*EventOrchestrationIntegration, *Response, error) {
	u := buildEventOrchestrationIntegrationUrl(orchestrationId, id)
	v := new(EventOrchestrationIntegrationPayload)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

func (s *EventOrchestrationIntegrationService) UpdateContext(ctx context.Context, orchestrationId string, id string, integration *EventOrchestrationIntegration) (*EventOrchestrationIntegration, *Response, error) {
	u := buildEventOrchestrationIntegrationUrl(orchestrationId, id)
	v := new(EventOrchestrationIntegrationPayload)
	p := &EventOrchestrationIntegrationPayload{Integration: integration}

	resp, err := s.client.newRequestDoContext(ctx, "PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Integration, resp, nil
}

func (s *EventOrchestrationIntegrationService) DeleteContext(ctx context.Context, orchestrationId string, id string) (*Response, error) {
	u := buildEventOrchestrationIntegrationUrl(orchestrationId, id)
	return s.client.newRequestDoContext(ctx, "DELETE", u, nil, nil, nil)
}

func (s *EventOrchestrationIntegrationService) MigrateFromOrchestrationContext(ctx context.Context, destinationOrchestrationId string, sourceOrchestrationId string, id string) (*ListEventOrchestrationIntegrationsResponse, *Response, error) {
	u := buildEventOrchestrationIntegrationUrl(destinationOrchestrationId, "migration")
	v := new(ListEventOrchestrationIntegrationsResponse)
	p := &EventOrchestrationIntegrationMigrationPayload{
		SourceType:    "orchestration",
		SourceId:      sourceOrchestrationId,
		IntegrationId: id,
	}

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, p, v)

	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
