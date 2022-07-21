package pagerduty

import (
	"context"
	"fmt"
)

// CustomFieldService handles the communication with field related methods of the PagerDuty API.
type CustomFieldService service

// CustomField represents a custom field.
type CustomField struct {
	ID           string               `json:"id,omitempty"`
	Name         string               `json:"name,omitempty"`
	DisplayName  string               `json:"display_name,omitempty"`
	Type         string               `json:"type,omitempty"`
	Summary      string               `json:"summary,omitempty"`
	Self         string               `json:"self,omitempty"`
	DataType     CustomFieldDataType  `json:"datatype,omitempty"`
	Description  *string              `json:"description,omitempty"`
	MultiValue   bool                 `json:"multi_value"`
	FixedOptions bool                 `json:"fixed_options"`
	FieldOptions []*CustomFieldOption `json:"field_options,omitempty"`
}

// ListCustomFieldResponse represents a list response of fields
type ListCustomFieldResponse struct {
	Total  int            `json:"total,omitempty"`
	Fields []*CustomField `json:"fields,omitempty"`
	Offset int            `json:"offset,omitempty"`
	More   bool           `json:"more,omitempty"`
	Limit  int            `json:"limit,omitempty"`
}

// CustomFieldPayload represents payload with a field object
type CustomFieldPayload struct {
	Field *CustomField `json:"field,omitempty"`
}

// ListCustomFieldOptions represents options when retrieving a list of fields.
type ListCustomFieldOptions struct {
	Offset   int      `url:"offset,omitempty"`
	Limit    int      `url:"limit,omitempty"`
	Total    bool     `url:"total,omitempty"`
	Includes []string `url:"include,brackets,omitempty"`
}

type listCustomFieldOptionsGen struct {
	options *ListCustomFieldOptions
}

func (o *listCustomFieldOptionsGen) currentOffset() int {
	return o.options.Offset
}

func (o *listCustomFieldOptionsGen) changeOffset(i int) {
	o.options.Offset = i
}

func (o *listCustomFieldOptionsGen) buildStruct() interface{} {
	return o.options
}

// GetCustomFieldOptions represents options when retrieving a field.
type GetCustomFieldOptions struct {
	Includes []string `url:"include,brackets,omitempty"`
}

var customFieldsEarlyAccessHeader = RequestOptions{
	Type:  "header",
	Label: "X-EARLY-ACCESS",
	Value: "flex-service-early-access",
}

// List lists existing custom fields. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
func (s *CustomFieldService) List(o *ListCustomFieldOptions) (*ListCustomFieldResponse, *Response, error) {
	return s.ListContext(context.Background(), o)
}

// ListContext lists existing custom fields. If a non-zero Limit is passed as an option, only a single page of results will be
// returned. Otherwise, the entire list of fields will be returned.
func (s *CustomFieldService) ListContext(ctx context.Context, o *ListCustomFieldOptions) (*ListCustomFieldResponse, *Response, error) {
	u := "/customfields/fields"
	v := new(ListCustomFieldResponse)

	if o == nil {
		o = &ListCustomFieldOptions{}
	}

	if o.Limit != 0 {
		resp, err := s.client.newRequestDoContext(ctx, "GET", u, o, nil, &v)
		if err != nil {
			return nil, nil, err
		}

		return v, resp, nil
	} else {
		fields := make([]*CustomField, 0)

		// Create a handler closure capable of parsing data from the fields endpoint
		// and appending resultant response plays to the return slice.
		responseHandler := func(response *Response) (ListResp, *Response, error) {
			var result ListCustomFieldResponse

			if err := s.client.DecodeJSON(response, &result); err != nil {
				return ListResp{}, response, err
			}

			fields = append(fields, result.Fields...)

			// Return stats on the current page. Caller can use this information to
			// adjust for requesting additional pages.
			return ListResp{
				More:   result.More,
				Offset: result.Offset,
				Limit:  result.Limit,
			}, response, nil
		}
		err := s.client.newRequestPagedGetQueryDo(u, responseHandler, &listCustomFieldOptionsGen{
			options: o,
		}, customFieldsEarlyAccessHeader)
		if err != nil {
			return nil, nil, err
		}
		v.Fields = fields

		return v, nil, nil
	}
}

// Get gets a custom field.
func (s *CustomFieldService) Get(id string, o *GetCustomFieldOptions) (*CustomField, *Response, error) {
	return s.GetContext(context.Background(), id, o)
}

// GetContext gets a custom field.
func (s *CustomFieldService) GetContext(ctx context.Context, id string, o *GetCustomFieldOptions) (*CustomField, *Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s", id)
	v := new(CustomFieldPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "GET", u, o, nil, v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.Field, resp, nil
}

// Create creates a new custom field.
func (s *CustomFieldService) Create(field *CustomField) (*CustomField, *Response, error) {
	return s.CreateContext(context.Background(), field)
}

// CreateContext creates a new custom field.
func (s *CustomFieldService) CreateContext(ctx context.Context, field *CustomField) (*CustomField, *Response, error) {
	u := "/customfields/fields"
	v := new(CustomFieldPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "POST", u, nil, &CustomFieldPayload{Field: field}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.Field, resp, nil
}

// Delete removes an existing custom field.
func (s *CustomFieldService) Delete(id string) (*Response, error) {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an existing custom field.
func (s *CustomFieldService) DeleteContext(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s", id)
	return s.client.newRequestDoOptionsContext(ctx, "DELETE", u, nil, nil, nil, customFieldsEarlyAccessHeader)
}

// Update updates an existing custom field.
func (s *CustomFieldService) Update(id string, field *CustomField) (*CustomField, *Response, error) {
	return s.UpdateContext(context.Background(), id, field)
}

// UpdateContext updates an existing custom field.
func (s *CustomFieldService) UpdateContext(ctx context.Context, id string, field *CustomField) (*CustomField, *Response, error) {
	u := fmt.Sprintf("/customfields/fields/%s", id)
	v := new(CustomFieldPayload)

	resp, err := s.client.newRequestDoOptionsContext(ctx, "PUT", u, nil, &CustomFieldPayload{Field: field}, &v, customFieldsEarlyAccessHeader)
	if err != nil {
		return nil, nil, err
	}

	return v.Field, resp, nil
}
