package pagerduty

import (
	"context"
	"fmt"
)

// IncidentWorkflowTriggerService handles the communication with incident workflow
// trigger related methods of the PagerDuty API.
type IncidentWorkflowTriggerService service

// IncidentWorkflowTrigger represents an incident workflow.
type IncidentWorkflowTrigger struct {
	ID                      string                              `json:"id,omitempty"`
	Type                    string                              `json:"type,omitempty"`
	TriggerType             IncidentWorkflowTriggerType         `json:"trigger_type,omitempty"`
	Workflow                *IncidentWorkflow                   `json:"workflow,omitempty"`
	Services                []*ServiceReference                 `json:"services,omitempty"`
	Condition               *string                             `json:"condition,omitempty"`
	SubscribedToAllServices bool                                `json:"is_subscribed_to_all_services,omitempty"`
	Permissions             *IncidentWorkflowTriggerPermissions `json:"permissions,omitempty"`
}

type IncidentWorkflowTriggerPermissions struct {
	Restricted bool   `json:"restricted"`
	TeamID     string `json:"team_id,omitempty"`
}

// ListIncidentWorkflowTriggerResponse represents a list response of incident workflow triggers.
type ListIncidentWorkflowTriggerResponse struct {
	Triggers      []*IncidentWorkflowTrigger `json:"triggers,omitempty"`
	NextPageToken string                     `json:"next_page_token,omitempty"`
	Limit         int                        `json:"limit,omitempty"`
}

// IncidentWorkflowTriggerPayload represents payload with an incident workflow trigger object.
type IncidentWorkflowTriggerPayload struct {
	Trigger *IncidentWorkflowTrigger `json:"trigger,omitempty"`
}

// ListIncidentWorkflowTriggerOptions represents options when retrieving a list of incident workflow triggers.
type ListIncidentWorkflowTriggerOptions struct {
	IncidentID  string                      `url:"incident_id,omitempty"`
	WorkflowID  string                      `url:"workflow_id,omitempty"`
	ServiceID   string                      `url:"service_id,omitempty"`
	TriggerType IncidentWorkflowTriggerType `url:"trigger_type,omitempty"`
	Limit       int                         `url:"limit,omitempty"`
	PageToken   string                      `url:"page_token,omitempty"`
}

type listIncidentWorkflowTriggerOptionsGen struct {
	options *ListIncidentWorkflowTriggerOptions
}

func (o *listIncidentWorkflowTriggerOptionsGen) currentCursor() string {
	return o.options.PageToken
}

func (o *listIncidentWorkflowTriggerOptionsGen) changeCursor(s string) {
	o.options.PageToken = s
}

func (o *listIncidentWorkflowTriggerOptionsGen) buildStruct() interface{} {
	return o.options
}

// List lists existing incident workflow triggers. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of incident workflows will be returned.
func (s *IncidentWorkflowTriggerService) List(o *ListIncidentWorkflowTriggerOptions) (*ListIncidentWorkflowTriggerResponse, *Response, error) {
	return s.ListContext(context.Background(), o)
}

// ListContext lists existing incident workflow triggers. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of incident workflows will be returned.
func (s *IncidentWorkflowTriggerService) ListContext(ctx context.Context, o *ListIncidentWorkflowTriggerOptions) (*ListIncidentWorkflowTriggerResponse, *Response, error) {
	u := "/incident_workflows/triggers"
	v := new(ListIncidentWorkflowTriggerResponse)

	if o == nil {
		o = &ListIncidentWorkflowTriggerOptions{}
	}

	if o.Limit != 0 {
		resp, err := s.client.newRequestDoContext(ctx, "GET", u, o, nil, &v)
		if err != nil {
			return nil, nil, err
		}

		return v, resp, nil
	} else {
		triggers := make([]*IncidentWorkflowTrigger, 0)

		// Create a handler closure capable of parsing data from the workflows endpoint
		// and appending resultant response plays to the return slice.
		responseHandler := func(response *Response) (CursorListResp, *Response, error) {
			var result ListIncidentWorkflowTriggerResponse

			if err := s.client.DecodeJSON(response, &result); err != nil {
				return CursorListResp{}, response, err
			}

			triggers = append(triggers, result.Triggers...)

			// Return stats on the current page. Caller can use this information to
			// adjust for requesting additional pages.
			return CursorListResp{
				Limit:      result.Limit,
				NextCursor: result.NextPageToken,
			}, response, nil
		}
		err := s.client.newRequestCursorPagedGetQueryDoContext(ctx, u, responseHandler, &listIncidentWorkflowTriggerOptionsGen{
			options: o,
		})
		if err != nil {
			return nil, nil, err
		}
		v.Triggers = triggers

		return v, nil, nil
	}
}

// Get gets an incident workflow trigger.
func (s *IncidentWorkflowTriggerService) Get(id string) (*IncidentWorkflowTrigger, *Response, error) {
	return s.GetContext(context.Background(), id)
}

// GetContext gets an incident workflow trigger.
func (s *IncidentWorkflowTriggerService) GetContext(ctx context.Context, id string) (*IncidentWorkflowTrigger, *Response, error) {
	u := fmt.Sprintf("/incident_workflows/triggers/%s", id)
	v := new(IncidentWorkflowTriggerPayload)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Trigger, resp, nil
}

// Create creates a new incident workflow trigger.
func (s *IncidentWorkflowTriggerService) Create(t *IncidentWorkflowTrigger) (*IncidentWorkflowTrigger, *Response, error) {
	return s.CreateContext(context.Background(), t)
}

// CreateContext creates a new incident workflow trigger.
func (s *IncidentWorkflowTriggerService) CreateContext(ctx context.Context, t *IncidentWorkflowTrigger) (*IncidentWorkflowTrigger, *Response, error) {
	u := "/incident_workflows/triggers"
	v := new(IncidentWorkflowTriggerPayload)

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, &t, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Trigger, resp, nil
}

// Delete removes an existing incident workflow trigger.
func (s *IncidentWorkflowTriggerService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an existing incident workflow trigger.
func (s *IncidentWorkflowTriggerService) DeleteContext(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("/incident_workflows/triggers/%s", id)
	return s.client.newRequestDoContext(ctx, "DELETE", u, nil, nil, nil)
}

// Update updates an existing incident workflow trigger.
func (s *IncidentWorkflowTriggerService) Update(id string, t *IncidentWorkflowTrigger) (*IncidentWorkflowTrigger, *Response, error) {
	return s.UpdateContext(context.Background(), id, t)
}

// UpdateContext updates an existing incident workflow trigger.
func (s *IncidentWorkflowTriggerService) UpdateContext(ctx context.Context, id string, t *IncidentWorkflowTrigger) (*IncidentWorkflowTrigger, *Response, error) {
	u := fmt.Sprintf("/incident_workflows/triggers/%s", id)
	v := new(IncidentWorkflowTriggerPayload)

	resp, err := s.client.newRequestDoContext(ctx, "PUT", u, nil, &t, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Trigger, resp, nil
}
