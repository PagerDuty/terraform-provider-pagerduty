package pagerduty

import (
	"context"
	"fmt"
)

// CustomFieldSchemaAssignmentService handles the communication with field schema assignment
// related methods of the PagerDuty API.
type CustomFieldSchemaAssignmentService service

// CustomFieldSchemaAssignment represents a field schema assignment.
type CustomFieldSchemaAssignment struct {
	ID      string                      `json:"id,omitempty"`
	Type    string                      `json:"type,omitempty"`
	Service *ServiceReference           `json:"service,omitempty"`
	Schema  *CustomFieldSchemaReference `json:"schema,omitempty"`
}

type CustomFieldSchemaAssignmentPayload struct {
	SchemaAssignment *CustomFieldSchemaAssignment `json:"schema_assignment,omitempty"`
}

// ListCustomFieldSchemaAssignmentsResponse represents a list response of resources assigned to field schemas.
type ListCustomFieldSchemaAssignmentsResponse struct {
	Total             int                            `json:"total,omitempty"`
	SchemaAssignments []*CustomFieldSchemaAssignment `json:"schema_assignments,omitempty"`
	Offset            int                            `json:"offset,omitempty"`
	More              bool                           `json:"more,omitempty"`
	Limit             int                            `json:"limit,omitempty"`
}

// ListCustomFieldSchemaAssignmentsOptions represents options when retrieving a list of fields schema assignments.
type ListCustomFieldSchemaAssignmentsOptions struct {
	Offset int  `url:"offset,omitempty"`
	Limit  int  `url:"limit,omitempty"`
	Total  bool `url:"total,omitempty"`
}

// fullListCustomFieldSchemaAssignmentsOptions represents options when retrieving a list of schema assignments
type fullListCustomFieldSchemaAssignmentsOptions struct {
	SchemaID  string `url:"schema_id,omitempty"`
	ServiceID string `url:"service_id,omitempty"`
	Offset    int    `url:"offset,omitempty"`
	Limit     int    `url:"limit,omitempty"`
	Total     bool   `url:"total,omitempty"`
}

type fullListCustomFieldSchemaAssignmentsOptionsGen struct {
	options *fullListCustomFieldSchemaAssignmentsOptions
}

func (o *fullListCustomFieldSchemaAssignmentsOptionsGen) currentOffset() int {
	return o.options.Offset
}

func (o *fullListCustomFieldSchemaAssignmentsOptionsGen) changeOffset(i int) {
	o.options.Offset = i
}

func (o *fullListCustomFieldSchemaAssignmentsOptionsGen) buildStruct() interface{} {
	return o.options
}

// Create creates a custom field schema assignment.
func (s *CustomFieldSchemaAssignmentService) Create(a *CustomFieldSchemaAssignment) (*CustomFieldSchemaAssignment, *Response, error) {
	return s.CreateContext(context.Background(), a)
}

// CreateContext creates a custom field schema assignment.
func (s *CustomFieldSchemaAssignmentService) CreateContext(ctx context.Context, a *CustomFieldSchemaAssignment) (*CustomFieldSchemaAssignment, *Response, error) {
	u := "/customfields/schema_assignments"
	v := new(CustomFieldSchemaAssignmentPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "POST", u, nil, &CustomFieldSchemaAssignmentPayload{SchemaAssignment: a}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.SchemaAssignment, resp, nil
}

// Delete removes a schema assignment
func (s *CustomFieldSchemaAssignmentService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes a schema assignment
func (s *CustomFieldSchemaAssignmentService) DeleteContext(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("/customfields/schema_assignments/%s", id)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "DELETE", u, nil, nil, nil, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListForSchema returns a list of assignments or the passed schema id
func (s *CustomFieldSchemaAssignmentService) ListForSchema(schemaID string, o *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	return s.ListForSchemaContext(context.Background(), schemaID, o)
}

// ListForSchemaContext returns a list of assignments for the passed schema id
func (s *CustomFieldSchemaAssignmentService) ListForSchemaContext(ctx context.Context, schemaID string, o *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	fo := fullListCustomFieldSchemaAssignmentsOptions{
		SchemaID: schemaID,
	}
	if o != nil {
		fo.Limit = o.Limit
		fo.Offset = o.Offset
		fo.Total = o.Total
	}

	return s.listContext(ctx, &fo)
}

// ListForService returns a list of schema assignments for the passed service id.
func (s *CustomFieldSchemaAssignmentService) ListForService(serviceID string, o *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	return s.ListForServiceContext(context.Background(), serviceID, o)
}

// ListForServiceContext returns a list of schema assignments for the passed service id.
func (s *CustomFieldSchemaAssignmentService) ListForServiceContext(ctx context.Context, serviceID string, o *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	fo := fullListCustomFieldSchemaAssignmentsOptions{
		ServiceID: serviceID,
	}
	if o != nil {
		fo.Limit = o.Limit
		fo.Offset = o.Offset
		fo.Total = o.Total
	}

	return s.listContext(ctx, &fo)
}

// listContext lists existing custom fields. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
func (s *CustomFieldSchemaAssignmentService) listContext(ctx context.Context, o *fullListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	u := "/customfields/schema_assignments"
	v := new(ListCustomFieldSchemaAssignmentsResponse)

	if o == nil {
		o = &fullListCustomFieldSchemaAssignmentsOptions{}
	}

	if o.Limit != 0 {
		resp, err := s.client.newRequestDoContext(ctx, "GET", u, o, nil, &v)
		if err != nil {
			return nil, nil, err
		}

		return v, resp, nil
	} else {
		assignments := make([]*CustomFieldSchemaAssignment, 0)

		// Create a handler closure capable of parsing data from the fields endpoint
		// and appending resultant response plays to the return slice.
		responseHandler := func(response *Response) (ListResp, *Response, error) {
			var result ListCustomFieldSchemaAssignmentsResponse

			if err := s.client.DecodeJSON(response, &result); err != nil {
				return ListResp{}, response, err
			}

			assignments = append(assignments, result.SchemaAssignments...)

			// Return stats on the current page. Caller can use this information to
			// adjust for requesting additional pages.
			return ListResp{
				More:   result.More,
				Offset: result.Offset,
				Limit:  result.Limit,
			}, response, nil
		}
		err := s.client.newRequestPagedGetQueryDo(u, responseHandler, &fullListCustomFieldSchemaAssignmentsOptionsGen{
			options: o,
		}, customFieldsEarlyAccessHeader)
		if err != nil {
			return nil, nil, err
		}
		v.SchemaAssignments = assignments

		return v, nil, nil
	}
}
