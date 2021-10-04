package pagerduty

import "fmt"

// TagService handles the communication with tag
// related methods of the PagerDuty API.
type TagService service

// Tag represents a tag.
type Tag struct {
	Label   string `json:"label,omitempty"`
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// ListTagsOptions represents options when listing tags.
type ListTagsOptions struct {
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Total  int    `url:"total,omitempty"`
	Query  string `url:"query,omitempty"`
}

// ListTagsResponse represents a list response of tags.
type ListTagsResponse struct {
	Limit  int    `json:"limit,omitempty"`
	More   bool   `json:"more,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Total  int    `json:"total,omitempty"`
	Tags   []*Tag `json:"tags,omitempty"`
}

// TagPayload represents payload with a tag object
type TagPayload struct {
	Tag *Tag `json:"tag,omitempty"`
}

// TagAssignments represent an object for adding and removing tag to entity assignments.
type TagAssignments struct {
	Add    []*TagAssignment `json:"add,omitempty"`
	Remove []*TagAssignment `json:"remove,omitempty"`
}

// TagAssignment represents a single tag assignment to an entity
type TagAssignment struct {
	Type       string `json:"type"`
	TagID      string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
}

// List lists existing tags.
func (s *TagService) List(o *ListTagsOptions) (*ListTagsResponse, *Response, error) {
	u := "/tags"
	v := new(ListTagsResponse)

	tags := make([]*Tag, 0)

	// Create a handler closure capable of parsing data from the response_plays endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListTagsResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		tags = append(tags, result.Tags...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(u, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.Tags = tags

	return v, nil, nil
}

// List Tags for a given Entity.
func (s *TagService) ListTagsForEntity(e, eid string) (*ListTagsResponse, *Response, error) {
	u := fmt.Sprintf("/%s/%s/tags", e, eid)
	v := new(ListTagsResponse)

	tags := make([]*Tag, 0)

	// Create a handler closure capable of parsing data from the response_plays endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListTagsResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		tags = append(tags, result.Tags...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(u, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.Tags = tags

	return v, nil, nil
}

// Create creates a new tag.
func (s *TagService) Create(tag *Tag) (*Tag, *Response, error) {
	u := "/tags"
	v := new(TagPayload)
	p := &TagPayload{Tag: tag}

	resp, err := s.client.newRequestDo("POST", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Tag, resp, nil
}

// Delete removes an existing tag.
func (s *TagService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/tags/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about a tag.
func (s *TagService) Get(id string) (*Tag, *Response, error) {
	u := fmt.Sprintf("/tags/%s", id)
	v := new(TagPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Tag, resp, nil
}

// Assign adds and removes tag assignments with entities
func (s *TagService) Assign(e, eid string, a *TagAssignments) (*Response, error) {
	u := fmt.Sprintf("/%s/%s/change_tags", e, eid)

	resp, err := s.client.newRequestDo("POST", u, nil, a, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
