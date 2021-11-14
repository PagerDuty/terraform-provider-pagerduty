package pagerduty

import "fmt"

// MaintenanceWindowService handles the communication with add-on related methods
// of the PagerDuty API.
type MaintenanceWindowService service

// MaintenanceWindow represents a PagerDuty maintenance window.
type MaintenanceWindow struct {
	CreatedBy      *UserReference      `json:"created_by,omitempty"`
	Description    string              `json:"description,omitempty"`
	EndTime        string              `json:"end_time,omitempty"`
	HTMLURL        string              `json:"html_url,omitempty"`
	ID             string              `json:"id,omitempty"`
	Self           string              `json:"self,omitempty"`
	SequenceNumber int                 `json:"sequence_number,omitempty"`
	Services       []*ServiceReference `json:"services,omitempty"`
	Src            string              `json:"src,omitempty"`
	StartTime      string              `json:"start_time,omitempty"`
	Summary        string              `json:"summary,omitempty"`
	Teams          []*TeamReference    `json:"teams,omitempty"`
	Type           string              `json:"type,omitempty"`
}

// ListMaintenanceWindowsOptions represents options when listing maintenance windows.
type ListMaintenanceWindowsOptions struct {
	Filter     string   `url:"filter,omitempty"`
	Include    []string `url:"include,omitempty,brackets"`
	Query      string   `url:"query,omitempty"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	TeamIDs    []string `url:"team_ids,omitempty,brackets"`
}

// ListMaintenanceWindowsResponse represents a list response of maintenance windows.
type ListMaintenanceWindowsResponse struct {
	Limit              int                  `json:"limit,omitempty"`
	MaintenanceWindows []*MaintenanceWindow `json:"maintenance_windows,omitempty"`
	More               bool                 `json:"more,omitempty"`
	Offset             int                  `json:"offset,omitempty"`
	Total              int                  `json:"total,omitempty"`
}

// MaintenanceWindowPayload represents a maintenance window.
type MaintenanceWindowPayload struct {
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty"`
}

// List lists existing maintenance windows.
func (s *MaintenanceWindowService) List(o *ListMaintenanceWindowsOptions) (*ListMaintenanceWindowsResponse, *Response, error) {
	u := "/maintenance_windows"
	v := new(ListMaintenanceWindowsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Create creates a new maintenancce window.
func (s *MaintenanceWindowService) Create(maintenanceWindow *MaintenanceWindow) (*MaintenanceWindow, *Response, error) {
	u := "/maintenance_windows"
	v := new(MaintenanceWindowPayload)

	resp, err := s.client.newRequestDo("POST", u, nil, &MaintenanceWindowPayload{MaintenanceWindow: maintenanceWindow}, v)
	if err != nil {
		return nil, nil, err
	}

	return v.MaintenanceWindow, resp, nil
}

// Delete removes an existing maintenance window.
func (s *MaintenanceWindowService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/maintenance_windows/%s", id)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Get retrieves information about a maintenance window.
func (s *MaintenanceWindowService) Get(id string) (*MaintenanceWindow, *Response, error) {
	u := fmt.Sprintf("/maintenance_windows/%s", id)
	v := new(MaintenanceWindowPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.MaintenanceWindow, resp, nil
}

// Update updates an existing maintenance window.
func (s *MaintenanceWindowService) Update(id string, maintenanceWindow *MaintenanceWindow) (*MaintenanceWindow, *Response, error) {
	u := fmt.Sprintf("/maintenance_windows/%s", id)
	v := new(MaintenanceWindowPayload)
	resp, err := s.client.newRequestDo("PUT", u, nil, &MaintenanceWindowPayload{MaintenanceWindow: maintenanceWindow}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.MaintenanceWindow, resp, nil
}
