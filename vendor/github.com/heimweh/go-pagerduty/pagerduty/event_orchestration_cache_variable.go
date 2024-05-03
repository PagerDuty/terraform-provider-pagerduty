package pagerduty

import (
	"context"
	"fmt"
)

type EventOrchestrationCacheVariableService service

type EventOrchestrationCacheVariableCondition struct {
	// A PCL string: https://developer.pagerduty.com/docs/ZG9jOjM1NTE0MDc0-pcl-overview
	Expression string `json:"expression,omitempty"`
}

// Configuration for a cache variable changes depending on the type:
//   - if `Type` is `recent_value`; then use `Regex` and `Source`
//   - if `Type` is `trigger_event_count`; then use `TTLSeconds`
type EventOrchestrationCacheVariableConfiguration struct {
	Type       string `json:"type,omitempty"`
	Regex      string `json:"regex,omitempty"`
	Source     string `json:"source,omitempty"`
	TTLSeconds int    `json:"ttl_seconds,omitempty"`
}

type EventOrchestrationCacheVariable struct {
	ID            string                                        `json:"id,omitempty"`
	Name          string                                        `json:"name,omitempty"`
	Disabled      bool                                          `json:"disabled"`
	Conditions    []*EventOrchestrationCacheVariableCondition   `json:"conditions"`
	Configuration *EventOrchestrationCacheVariableConfiguration `json:"configuration,omitempty"`
	CreatedAt     string                                        `json:"created_at,omitempty"`
	CreatedBy     *UserReference                                `json:"created_by,omitempty"`
	UpdatedAt     string                                        `json:"updated_at,omitempty"`
	UpdatedBy     *UserReference                                `json:"updated_by,omitempty"`
}

type EventOrchestrationCacheVariablePayload struct {
	CacheVariable *EventOrchestrationCacheVariable `json:"cache_variable,omitempty"`
}

type ListEventOrchestrationCacheVariablesResponse struct {
	Total          int                                `json:"total,omitempty"`
	CacheVariables []*EventOrchestrationCacheVariable `json:"cache_variables,omitempty"`
}

const CacheVariableTypeGlobal string = "global"
const CacheVariableTypeService string = "service"

func buildEventOrchestrationCacheVariableUrl(cacheVariableType string, orchestrationId string, cacheVariableId string) string {
	if cacheVariableType == CacheVariableTypeService {
		return fmt.Sprintf("%s/services/%s/cache_variables/%s", eventOrchestrationBaseUrl, orchestrationId, cacheVariableId)
	}

	return fmt.Sprintf("%s/%s/cache_variables/%s", eventOrchestrationBaseUrl, orchestrationId, cacheVariableId)
}

func (s *EventOrchestrationCacheVariableService) List(ctx context.Context, cacheVariableType string, orchestrationId string) (*ListEventOrchestrationCacheVariablesResponse, *Response, error) {
	u := buildEventOrchestrationCacheVariableUrl(cacheVariableType, orchestrationId, "")
	v := new(ListEventOrchestrationCacheVariablesResponse)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, v)

	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (s *EventOrchestrationCacheVariableService) Create(ctx context.Context, cacheVariableType string, orchestrationId string, cacheVariable *EventOrchestrationCacheVariable) (*EventOrchestrationCacheVariable, *Response, error) {
	u := buildEventOrchestrationCacheVariableUrl(cacheVariableType, orchestrationId, "")
	v := new(EventOrchestrationCacheVariablePayload)
	p := &EventOrchestrationCacheVariablePayload{CacheVariable: cacheVariable}

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, p, v)

	if err != nil {
		return nil, nil, err
	}

	return v.CacheVariable, resp, nil
}

func (s *EventOrchestrationCacheVariableService) Get(ctx context.Context, cacheVariableType string, orchestrationId string, cacheVariableId string) (*EventOrchestrationCacheVariable, *Response, error) {
	u := buildEventOrchestrationCacheVariableUrl(cacheVariableType, orchestrationId, cacheVariableId)
	v := new(EventOrchestrationCacheVariablePayload)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.CacheVariable, resp, nil
}

func (s *EventOrchestrationCacheVariableService) Update(ctx context.Context, cacheVariableType string, orchestrationId string, cacheVariableId string, cacheVariable *EventOrchestrationCacheVariable) (*EventOrchestrationCacheVariable, *Response, error) {
	u := buildEventOrchestrationCacheVariableUrl(cacheVariableType, orchestrationId, cacheVariableId)
	v := new(EventOrchestrationCacheVariablePayload)
	p := &EventOrchestrationCacheVariablePayload{CacheVariable: cacheVariable}

	resp, err := s.client.newRequestDoContext(ctx, "PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.CacheVariable, resp, nil
}

func (s *EventOrchestrationCacheVariableService) Delete(ctx context.Context, cacheVariableType string, orchestrationId string, cacheVariableId string) (*Response, error) {
	u := buildEventOrchestrationCacheVariableUrl(cacheVariableType, orchestrationId, cacheVariableId)
	return s.client.newRequestDoContext(ctx, "DELETE", u, nil, nil, nil)
}
