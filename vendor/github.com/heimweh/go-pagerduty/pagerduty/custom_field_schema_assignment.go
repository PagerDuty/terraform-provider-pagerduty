package pagerduty

import (
	"context"
)

// CustomFieldSchemaAssignmentService handles the communication with field schema assignment
// related methods of the PagerDuty API.
//
// Deprecated: No current replacement
type CustomFieldSchemaAssignmentService service

// CustomFieldSchemaAssignment represents a field schema assignment.
//
// Deprecated: This struct should no longer be used
type CustomFieldSchemaAssignment struct {
	ID      string                      `json:"id,omitempty"`
	Type    string                      `json:"type,omitempty"`
	Service *ServiceReference           `json:"service,omitempty"`
	Schema  *CustomFieldSchemaReference `json:"schema,omitempty"`
}

// Deprecated: This struct should no longer be used
type CustomFieldSchemaAssignmentPayload struct {
	SchemaAssignment *CustomFieldSchemaAssignment `json:"schema_assignment,omitempty"`
}

// ListCustomFieldSchemaAssignmentsResponse represents a list response of resources assigned to field schemas.
//
// Deprecated: No current replacement
type ListCustomFieldSchemaAssignmentsResponse struct {
	Total             int                            `json:"total,omitempty"`
	SchemaAssignments []*CustomFieldSchemaAssignment `json:"schema_assignments,omitempty"`
	Offset            int                            `json:"offset,omitempty"`
	More              bool                           `json:"more,omitempty"`
	Limit             int                            `json:"limit,omitempty"`
}

// ListCustomFieldSchemaAssignmentsOptions represents options when retrieving a list of fields schema assignments.
//
// Deprecated: No current replacement
type ListCustomFieldSchemaAssignmentsOptions struct {
	Offset int  `url:"offset,omitempty"`
	Limit  int  `url:"limit,omitempty"`
	Total  bool `url:"total,omitempty"`
}

// Create creates a custom field schema assignment.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) Create(a *CustomFieldSchemaAssignment) (*CustomFieldSchemaAssignment, *Response, error) {
	return s.CreateContext(context.Background(), a)
}

// CreateContext creates a custom field schema assignment.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) CreateContext(_ context.Context, _ *CustomFieldSchemaAssignment) (*CustomFieldSchemaAssignment, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Delete removes a schema assignment
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes a schema assignment
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) DeleteContext(_ context.Context, _ string) (*Response, error) {
	return nil, customFieldDeprecationError()
}

// ListForSchema returns a list of assignments or the passed schema id
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) ListForSchema(schemaID string, o *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	return s.ListForSchemaContext(context.Background(), schemaID, o)
}

// ListForSchemaContext returns a list of assignments for the passed schema id
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) ListForSchemaContext(_ context.Context, _ string, _ *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// ListForService returns a list of schema assignments for the passed service id.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) ListForService(serviceID string, o *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	return s.ListForServiceContext(context.Background(), serviceID, o)
}

// ListForServiceContext returns a list of schema assignments for the passed service id.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaAssignmentService) ListForServiceContext(_ context.Context, _ string, _ *ListCustomFieldSchemaAssignmentsOptions) (*ListCustomFieldSchemaAssignmentsResponse, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}
