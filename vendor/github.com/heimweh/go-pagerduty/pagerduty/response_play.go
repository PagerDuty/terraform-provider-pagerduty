package pagerduty

import (
	"fmt"
)

// ResponsePlayService handles the communication with response_plays
// related methods of the PagerDuty API.
type ResponsePlayService service

// ResponsePlay represents a response play.
type ResponsePlay struct {
	ID                 string                 `json:"id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Description        string                 `json:"description,omitempty"`
	Team               *TeamReference         `json:"team,omitempty"`
	Subscribers        *[]SubscriberReference `json:"subscribers,omitempty"`
	SubscribersMessage string                 `json:"subscribers_message"`
	Responders         *[]Responder           `json:"responders,omitempty"`
	RespondersMessage  string                 `json:"responders_message"`
	Runnability        string                 `json:"runnability,omitempty"`
	ConferenceNumber   string                 `json:"conference_number,omitempty"`
	ConferenceURL      string                 `json:"conference_url,omitempty"`
}

// Responder represents a responder within a response play object
type Responder struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// ResponsePlayPayload represents payload with a response play object
type ResponsePlayPayload struct {
	ResponsePlay *ResponsePlay `json:"response_play,omitempty"`
}

// ListResponsePlaysResponse represents a list response of response plays.
type ListResponsePlaysResponse struct {
	Total         int             `json:"total,omitempty"`
	ResponsePlays []*ResponsePlay `json:"response_plays,omitempty"`
	Offset        int             `json:"offset,omitempty"`
	More          bool            `json:"more,omitempty"`
	Limit         int             `json:"limit,omitempty"`
}

// List lists existing response_plays.
func (s *ResponsePlayService) List() (*ListResponsePlaysResponse, *Response, error) {
	u := "/response_plays"
	v := new(ListResponsePlaysResponse)

	responsePlays := make([]*ResponsePlay, 0)

	// Create a handler closure capable of parsing data from the response_plays endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListResponsePlaysResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		responsePlays = append(responsePlays, result.ResponsePlays...)

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
	v.ResponsePlays = responsePlays

	return v, nil, nil
}

// Create creates a new response play.
func (s *ResponsePlayService) Create(responsePlay *ResponsePlay) (*ResponsePlay, *Response, error) {
	u := "/response_plays"
	v := new(ResponsePlayPayload)
	p := &ResponsePlayPayload{ResponsePlay: responsePlay}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.ResponsePlay, resp, nil
}

// Get gets a new response play.
func (s *ResponsePlayService) Get(ID string) (*ResponsePlay, *Response, error) {
	u := fmt.Sprintf("/response_plays/%s", ID)
	v := new(ResponsePlayPayload)
	p := &ResponsePlayPayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.ResponsePlay, resp, nil
}

// Delete deletes an existing response_play.
func (s *ResponsePlayService) Delete(ID string) (*Response, error) {
	u := fmt.Sprintf("/response_plays/%s", ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Update updates an existing response_play.
func (s *ResponsePlayService) Update(ID string, responsePlay *ResponsePlay) (*ResponsePlay, *Response, error) {
	u := fmt.Sprintf("/response_plays/%s", ID)
	v := new(ResponsePlayPayload)
	p := ResponsePlayPayload{ResponsePlay: responsePlay}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.ResponsePlay, resp, nil
}
