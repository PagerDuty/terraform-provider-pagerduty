package pagerduty

import (
	"context"
	"fmt"
)

// IncidentWorkflowService handles the communication with incident workflow
// related methods of the PagerDuty API.
type IncidentWorkflowService service

// IncidentWorkflow represents an incident workflow.
type IncidentWorkflow struct {
	ID          string                  `json:"id,omitempty"`
	Type        string                  `json:"type,omitempty"`
	Name        string                  `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Self        string                  `json:"self,omitempty"`
	Steps       []*IncidentWorkflowStep `json:"steps,omitempty"`
	Team        *TeamReference          `json:"team,omitempty"`
}

// IncidentWorkflowStep represents a step in an incident workflow.
type IncidentWorkflowStep struct {
	ID            string                               `json:"id,omitempty"`
	Type          string                               `json:"type,omitempty"`
	Name          string                               `json:"name,omitempty"`
	Description   *string                              `json:"description,omitempty"`
	Configuration *IncidentWorkflowActionConfiguration `json:"action_configuration,omitempty"`
}

// IncidentWorkflowActionConfiguration represents the configuration for an incident workflow action
type IncidentWorkflowActionConfiguration struct {
	ActionID          string                                    `json:"action_id,omitempty"`
	Description       *string                                   `json:"description,omitempty"`
	Inputs            []*IncidentWorkflowActionInput            `json:"inputs,omitempty"`
	InlineStepsInputs []*IncidentWorkflowActionInlineStepsInput `json:"inline_steps_inputs,omitempty"`
}

// IncidentWorkflowActionInput represents the configuration for an incident workflow action input with a serialized string as the value
type IncidentWorkflowActionInput struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// IncidentWorkflowActionInlineStepsInput represents the configuration for an incident workflow action input with a series of inlined steps as the value
type IncidentWorkflowActionInlineStepsInput struct {
	Name  string                                       `json:"name,omitempty"`
	Value *IncidentWorkflowActionInlineStepsInputValue `json:"value,omitempty"`
}

// IncidentWorkflowActionInlineStepsInputValue represents the value for an inline_steps_input input
type IncidentWorkflowActionInlineStepsInputValue struct {
	Steps []*IncidentWorkflowActionInlineStep `json:"steps,omitempty"`
}

// IncidentWorkflowActionInlineStep represents a single step within an inline_steps_input input's value
type IncidentWorkflowActionInlineStep struct {
	Name          string                               `json:"name,omitempty"`
	Configuration *IncidentWorkflowActionConfiguration `json:"action_configuration,omitempty"`
}

// ListIncidentWorkflowResponse represents a list response of incident workflows.
type ListIncidentWorkflowResponse struct {
	Total             int                 `json:"total,omitempty"`
	IncidentWorkflows []*IncidentWorkflow `json:"incident_workflows,omitempty"`
	Offset            int                 `json:"offset,omitempty"`
	More              bool                `json:"more,omitempty"`
	Limit             int                 `json:"limit,omitempty"`
}

// IncidentWorkflowPayload represents payload with an incident workflow object.
type IncidentWorkflowPayload struct {
	IncidentWorkflow *IncidentWorkflow `json:"incident_workflow,omitempty"`
}

// ListIncidentWorkflowOptions represents options when retrieving a list of incident workflows.
type ListIncidentWorkflowOptions struct {
	Offset   int      `url:"offset,omitempty"`
	Limit    int      `url:"limit,omitempty"`
	Total    bool     `url:"total,omitempty"`
	Includes []string `url:"include,brackets,omitempty"`
}

type listIncidentWorkflowOptionsGen struct {
	options *ListIncidentWorkflowOptions
}

func (o *listIncidentWorkflowOptionsGen) currentOffset() int {
	return o.options.Offset
}

func (o *listIncidentWorkflowOptionsGen) changeOffset(i int) {
	o.options.Offset = i
}

func (o *listIncidentWorkflowOptionsGen) buildStruct() interface{} {
	return o.options
}

// List lists existing incident workflows. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of incident workflows will be returned.
func (s *IncidentWorkflowService) List(o *ListIncidentWorkflowOptions) (*ListIncidentWorkflowResponse, *Response, error) {
	return s.ListContext(context.Background(), o)
}

// ListContext lists existing incident workflows. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of incident workflows will be returned.
func (s *IncidentWorkflowService) ListContext(ctx context.Context, o *ListIncidentWorkflowOptions) (*ListIncidentWorkflowResponse, *Response, error) {
	u := "/incident_workflows"
	v := new(ListIncidentWorkflowResponse)

	if o == nil {
		o = &ListIncidentWorkflowOptions{}
	}

	if o.Limit != 0 {
		resp, err := s.client.newRequestDoContext(ctx, "GET", u, o, nil, &v)
		if err != nil {
			return nil, nil, err
		}

		return v, resp, nil
	} else {
		workflows := make([]*IncidentWorkflow, 0)

		// Create a handler closure capable of parsing data from the workflows endpoint
		// and appending resultant response plays to the return slice.
		responseHandler := func(response *Response) (ListResp, *Response, error) {
			var result ListIncidentWorkflowResponse

			if err := s.client.DecodeJSON(response, &result); err != nil {
				return ListResp{}, response, err
			}

			workflows = append(workflows, result.IncidentWorkflows...)

			// Return stats on the current page. Caller can use this information to
			// adjust for requesting additional pages.
			return ListResp{
				More:   result.More,
				Offset: result.Offset,
				Limit:  result.Limit,
			}, response, nil
		}
		err := s.client.newRequestPagedGetQueryDoContext(ctx, u, responseHandler, &listIncidentWorkflowOptionsGen{
			options: o,
		})
		if err != nil {
			return nil, nil, err
		}
		v.IncidentWorkflows = workflows

		return v, nil, nil
	}
}

// Get gets an incident workflow.
func (s *IncidentWorkflowService) Get(id string) (*IncidentWorkflow, *Response, error) {
	return s.GetContext(context.Background(), id)
}

// GetContext gets an incident workflow.
func (s *IncidentWorkflowService) GetContext(ctx context.Context, id string) (*IncidentWorkflow, *Response, error) {
	u := fmt.Sprintf("/incident_workflows/%s", id)
	v := new(IncidentWorkflowPayload)

	resp, err := s.client.newRequestDoContext(ctx, "GET", u, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v.IncidentWorkflow, resp, nil
}

// Create creates a new incident workflow.
func (s *IncidentWorkflowService) Create(iw *IncidentWorkflow) (*IncidentWorkflow, *Response, error) {
	return s.CreateContext(context.Background(), iw)
}

// CreateContext creates a new incident workflow.
func (s *IncidentWorkflowService) CreateContext(ctx context.Context, iw *IncidentWorkflow) (*IncidentWorkflow, *Response, error) {
	u := "/incident_workflows"
	v := new(IncidentWorkflowPayload)

	resp, err := s.client.newRequestDoContext(ctx, "POST", u, nil, &IncidentWorkflowPayload{IncidentWorkflow: iw}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.IncidentWorkflow, resp, nil
}

// Delete removes an existing incident workflow.
func (s *IncidentWorkflowService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an existing incident workflow.
func (s *IncidentWorkflowService) DeleteContext(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("/incident_workflows/%s", id)
	return s.client.newRequestDoContext(ctx, "DELETE", u, nil, nil, nil)
}

// Update updates an existing incident workflow.
func (s *IncidentWorkflowService) Update(id string, iw *IncidentWorkflow) (*IncidentWorkflow, *Response, error) {
	return s.UpdateContext(context.Background(), id, iw)
}

// UpdateContext updates an existing incident workflow.
func (s *IncidentWorkflowService) UpdateContext(ctx context.Context, id string, iw *IncidentWorkflow) (*IncidentWorkflow, *Response, error) {
	u := fmt.Sprintf("/incident_workflows/%s", id)
	v := new(IncidentWorkflowPayload)

	resp, err := s.client.newRequestDoContext(ctx, "PUT", u, nil, &IncidentWorkflowPayload{IncidentWorkflow: iw}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.IncidentWorkflow, resp, nil
}
