package pagerduty

import (
	"context"
)

// CustomFieldSchemaService handles the communication with field schema
// related methods of the PagerDuty API.
//
// Deprecated: No current replacement
type CustomFieldSchemaService service

// CustomFieldSchema represents a field schema.
//
// Deprecated: This struct should no longer be used
type CustomFieldSchema struct {
	ID                  string                                 `json:"id,omitempty"`
	Title               string                                 `json:"title,omitempty"`
	Type                string                                 `json:"type,omitempty"`
	Description         *string                                `json:"description,omitempty"`
	FieldConfigurations []*CustomFieldSchemaFieldConfiguration `json:"field_configurations,omitempty"`
}

// ListCustomFieldSchemaOptions represents options when retrieving a list of field schemas.
//
// Deprecated: This struct should no longer be used
type ListCustomFieldSchemaOptions struct {
	Offset int  `url:"offset,omitempty"`
	Limit  int  `url:"limit,omitempty"`
	Total  bool `url:"total,omitempty"`
}

// GetCustomFieldSchemaOptions represents options when retrieving a field schema
//
// Deprecated: This struct should no longer be used
type GetCustomFieldSchemaOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

// ListCustomFieldSchemaResponse represents a list response of field schemas
//
// Deprecated
type ListCustomFieldSchemaResponse struct {
	Total   int                  `json:"total,omitempty"`
	Schemas []*CustomFieldSchema `json:"schemas,omitempty"`
	Offset  int                  `json:"offset,omitempty"`
	More    bool                 `json:"more,omitempty"`
	Limit   int                  `json:"limit,omitempty"`
}

// Deprecated: This struct should no longer be used
type CustomFieldSchemaPayload struct {
	Schema *CustomFieldSchema `json:"schema,omitempty"`
}

// List lists existing field schemas. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) List(o *ListCustomFieldSchemaOptions) (*ListCustomFieldSchemaResponse, *Response, error) {
	return s.ListContext(context.Background(), o)
}

// ListContext lists existing field schemas. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) ListContext(_ context.Context, _ *ListCustomFieldSchemaOptions) (*ListCustomFieldSchemaResponse, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Get gets a field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) Get(id string, o *GetCustomFieldSchemaOptions) (*CustomFieldSchema, *Response, error) {
	return s.GetContext(context.Background(), id, o)
}

// GetContext gets a field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) GetContext(_ context.Context, _ string, _ *GetCustomFieldSchemaOptions) (*CustomFieldSchema, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Create creates a field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) Create(schema *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	return s.CreateContext(context.Background(), schema)
}

// CreateContext creates a field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) CreateContext(_ context.Context, _ *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Update updates a field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) Update(id string, schema *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	return s.UpdateContext(context.Background(), id, schema)
}

// UpdateContext updates a field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) UpdateContext(_ context.Context, _ string, _ *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	return nil, nil, customFieldDeprecationError()
}

// Delete removes an existing field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an existing field schema.
//
// Deprecated: No current replacement
func (s *CustomFieldSchemaService) DeleteContext(_ context.Context, _ string) (*Response, error) {
	return nil, customFieldDeprecationError()
}
