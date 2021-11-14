package pagerduty

import "fmt"

// SlackConnectionService handles the communication with the integration slack
// related methods of the PagerDuty API.
type SlackConnectionService service

// SlackConnection represents a slack connection.
type SlackConnection struct {
	ID               string           `json:"id,omitempty"`
	SourceID         string           `json:"source_id,omitempty"`
	SourceName       string           `json:"source_name,omitempty"`
	SourceType       string           `json:"source_type,omitempty"`
	ChannelID        string           `json:"channel_id,omitempty"`
	ChannelName      string           `json:"channel_name,omitempty"`
	WorkspaceID      string           `json:"workspace_id,omitempty"`
	Config           ConnectionConfig `json:"config,omitempty"`
	NotificationType string           `json:"notification_type,omitempty"`
}

// ConnectionConfig represents a config object in a slack connection
type ConnectionConfig struct {
	Events     []string `json:"events,omitempty"`
	Priorities []string `json:"priorities"`
	Urgency    *string  `json:"urgency"`
}

// SlackConnectionPayload represents payload with a slack connect object
type SlackConnectionPayload struct {
	SlackConnection *SlackConnection `json:"slack_connection,omitempty"`
}

// ListSlackConnectionsResponse represents a list response of slack connections.
type ListSlackConnectionsResponse struct {
	Total            int                `json:"total,omitempty"`
	SlackConnections []*SlackConnection `json:"slack_connections,omitempty"`
	Offset           int                `json:"offset,omitempty"`
	More             bool               `json:"more,omitempty"`
	Limit            int                `json:"limit,omitempty"`
}

// List lists existing slack connections.
func (s *SlackConnectionService) List(workspaceID string) (*ListSlackConnectionsResponse, *Response, error) {
	u := fmt.Sprintf("/integration-slack/workspaces/%s/connections", workspaceID)
	v := new(ListSlackConnectionsResponse)

	slackConnections := make([]*SlackConnection, 0)

	// Create a handler closure capable of parsing data from the integration-slack connections endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListSlackConnectionsResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		slackConnections = append(slackConnections, result.SlackConnections...)

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
	v.SlackConnections = slackConnections

	return v, nil, nil
}

// Create creates a new slack connection.
func (s *SlackConnectionService) Create(workspaceID string, sconn *SlackConnection) (*SlackConnection, *Response, error) {
	u := fmt.Sprintf("/integration-slack/workspaces/%s/connections", workspaceID)
	v := new(SlackConnectionPayload)
	p := &SlackConnectionPayload{SlackConnection: sconn}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}
	// Slack Connection in Terraform Provider needs workspaceID set to the object
	v.SlackConnection.WorkspaceID = workspaceID

	return v.SlackConnection, resp, nil
}

// Get gets a slack connection.
func (s *SlackConnectionService) Get(workspaceID, ID string) (*SlackConnection, *Response, error) {
	u := fmt.Sprintf("/integration-slack/workspaces/%s/connections/%s", workspaceID, ID)
	v := new(SlackConnectionPayload)
	p := &SlackConnectionPayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.SlackConnection, resp, nil
}

// Delete deletes a slack connection.
func (s *SlackConnectionService) Delete(workspaceID, ID string) (*Response, error) {
	u := fmt.Sprintf("/integration-slack/workspaces/%s/connections/%s", workspaceID, ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Update updates a slack connection.
func (s *SlackConnectionService) Update(workspaceID, ID string, sconn *SlackConnection) (*SlackConnection, *Response, error) {
	u := fmt.Sprintf("/integration-slack/workspaces/%s/connections/%s", workspaceID, ID)
	v := new(SlackConnectionPayload)
	p := SlackConnectionPayload{SlackConnection: sconn}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.SlackConnection, resp, nil
}
