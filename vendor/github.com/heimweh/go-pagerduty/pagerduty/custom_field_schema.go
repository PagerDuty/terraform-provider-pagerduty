package pagerduty

import (
	"context"
	"fmt"
)

// CustomFieldSchemaService handles the communication with field schema
// related methods of the PagerDuty API.
type CustomFieldSchemaService service

// CustomFieldSchema represents a field schema.
type CustomFieldSchema struct {
	ID                  string                                 `json:"id,omitempty"`
	Title               string                                 `json:"title,omitempty"`
	Type                string                                 `json:"type,omitempty"`
	Description         *string                                `json:"description,omitempty"`
	FieldConfigurations []*CustomFieldSchemaFieldConfiguration `json:"field_configurations,omitempty"`
}

// ListCustomFieldSchemaOptions represents options when retrieving a list of field schemas.
type ListCustomFieldSchemaOptions struct {
	Offset int  `url:"offset,omitempty"`
	Limit  int  `url:"limit,omitempty"`
	Total  bool `url:"total,omitempty"`
}

// GetCustomFieldSchemaOptions represents options when retrieving a field schema
type GetCustomFieldSchemaOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

type listCustomFieldSchemaOptionsGen struct {
	options *ListCustomFieldSchemaOptions
}

func (o *listCustomFieldSchemaOptionsGen) currentOffset() int {
	return o.options.Offset
}

func (o *listCustomFieldSchemaOptionsGen) changeOffset(i int) {
	o.options.Offset = i
}

func (o *listCustomFieldSchemaOptionsGen) buildStruct() interface{} {
	return o.options
}

// ListCustomFieldSchemaResponse represents a list response of field schemas
type ListCustomFieldSchemaResponse struct {
	Total   int                  `json:"total,omitempty"`
	Schemas []*CustomFieldSchema `json:"schemas,omitempty"`
	Offset  int                  `json:"offset,omitempty"`
	More    bool                 `json:"more,omitempty"`
	Limit   int                  `json:"limit,omitempty"`
}

type CustomFieldSchemaPayload struct {
	Schema *CustomFieldSchema `json:"schema,omitempty"`
}

// List lists existing field schemas. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
func (s *CustomFieldSchemaService) List(o *ListCustomFieldSchemaOptions) (*ListCustomFieldSchemaResponse, *Response, error) {
	return s.ListContext(context.Background(), o)
}

// ListContext lists existing field schemas. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
func (s *CustomFieldSchemaService) ListContext(ctx context.Context, o *ListCustomFieldSchemaOptions) (*ListCustomFieldSchemaResponse, *Response, error) {
	u := "/customfields/schemas"
	v := new(ListCustomFieldSchemaResponse)

	if o == nil {
		o = &ListCustomFieldSchemaOptions{}
	}

	if o.Limit != 0 {
		resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, o, nil, &v, customFieldsEarlyAccessHeader)
		if err != nil {
			return nil, nil, err
		}

		return v, resp, nil
	} else {
		schemas := make([]*CustomFieldSchema, 0)

		// Create a handler closure capable of parsing data from the fields endpoint
		// and appending resultant response plays to the return slice.
		responseHandler := func(response *Response) (ListResp, *Response, error) {
			var result ListCustomFieldSchemaResponse

			if err := s.client.DecodeJSON(response, &result); err != nil {
				return ListResp{}, response, err
			}

			schemas = append(schemas, result.Schemas...)

			// Return stats on the current page. Caller can use this information to
			// adjust for requesting additional pages.
			return ListResp{
				More:   result.More,
				Offset: result.Offset,
				Limit:  result.Limit,
			}, response, nil
		}
		err := s.client.newRequestPagedGetQueryDoContext(ctx, u, responseHandler, &listCustomFieldSchemaOptionsGen{
			options: o,
		}, customFieldsEarlyAccessHeader)
		if err != nil {
			return nil, nil, err
		}
		v.Schemas = schemas

		return v, nil, nil
	}
}

// Get gets a field schema.
func (s *CustomFieldSchemaService) Get(id string, o *GetCustomFieldSchemaOptions) (*CustomFieldSchema, *Response, error) {
	return s.GetContext(context.Background(), id, o)
}

// GetContext gets a field schema.
func (s *CustomFieldSchemaService) GetContext(ctx context.Context, id string, o *GetCustomFieldSchemaOptions) (*CustomFieldSchema, *Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s", id)
	v := new(CustomFieldSchemaPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, o, nil, v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.Schema, resp, nil
}

// Create creates a field schema.
func (s *CustomFieldSchemaService) Create(schema *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	return s.CreateContext(context.Background(), schema)
}

// CreateContext creates a field schema.
func (s *CustomFieldSchemaService) CreateContext(ctx context.Context, schema *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	u := "/customfields/schemas"
	v := new(CustomFieldSchemaPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "POST", u, nil, &CustomFieldSchemaPayload{Schema: schema}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.Schema, resp, nil
}

// Update updates a field schema.
func (s *CustomFieldSchemaService) Update(id string, schema *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	return s.UpdateContext(context.Background(), id, schema)
}

// UpdateContext updates a field schema.
func (s *CustomFieldSchemaService) UpdateContext(ctx context.Context, id string, schema *CustomFieldSchema) (*CustomFieldSchema, *Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s", id)
	v := new(CustomFieldSchemaPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "PUT", u, nil, &CustomFieldSchemaPayload{Schema: schema}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.Schema, resp, nil
}

// Delete removes an existing field schema.
func (s *CustomFieldSchemaService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an existing field schema.
func (s *CustomFieldSchemaService) DeleteContext(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("/customfields/schemas/%s", id)
	return s.client.newRequestDoOptionsContext(ctx, "DELETE", u, nil, nil, nil, customFieldsEarlyAccessHeader)
}
