package pagerduty

import "fmt"

// AutomationActionsAction handles the communication with Automation Actions
// related methods of the PagerDuty API.
type AutomationActionsActionService service

type AutomationActionsAction struct {
	ID                   string                               `json:"id"`
	Name                 string                               `json:"name"`
	Description          *string                              `json:"description,omitempty"`
	ActionType           string                               `json:"action_type"`
	RunnerID             *string                              `json:"runner,omitempty"`
	ActionDataReference  AutomationActionsActionDataReference `json:"action_data_reference"`
	Services             []*ServiceReference                  `json:"services,omitempty"`
	Teams                []*TeamReference                     `json:"teams,omitempty"`
	Privileges           *AutomationActionsPrivileges         `json:"privileges,omitempty"`
	Type                 *string                              `json:"type,omitempty"`
	ActionClassification *string                              `json:"action_classification,omitempty"`
	RunnerType           *string                              `json:"runner_type,omitempty"`
	CreationTime         *string                              `json:"creation_time,omitempty"`
	ModifyTime           *string                              `json:"modify_time,omitempty"`
}

type AutomationActionsActionDataReference struct {
	ProcessAutomationJobId        *string `json:"process_automation_job_id,omitempty"`
	ProcessAutomationJobArguments *string `json:"process_automation_job_arguments,omitempty"`
	Script                        *string `json:"script,omitempty"`
	InvocationCommand             *string `json:"invocation_command,omitempty"`
}

type AutomationActionsActionPayload struct {
	Action *AutomationActionsAction `json:"action,omitempty"`
}

// Create creates a new action
func (s *AutomationActionsActionService) Create(action *AutomationActionsAction) (*AutomationActionsAction, *Response, error) {
	u := "/automation_actions/actions"
	v := new(AutomationActionsActionPayload)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	resp, err := s.client.newRequestDoOptions("POST", u, nil, &AutomationActionsActionPayload{Action: action}, &v, o)
	if err != nil {
		return nil, nil, err
	}

	return v.Action, resp, nil
}

// Get retrieves information about an action.
func (s *AutomationActionsActionService) Get(id string) (*AutomationActionsAction, *Response, error) {
	u := fmt.Sprintf("/automation_actions/actions/%s", id)
	v := new(AutomationActionsActionPayload)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	resp, err := s.client.newRequestDoOptions("GET", u, nil, nil, &v, o)
	if err != nil {
		return nil, nil, err
	}

	return v.Action, resp, nil
}

// Delete deletes an existing action.
func (s *AutomationActionsActionService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/automation_actions/actions/%s", id)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	return s.client.newRequestDoOptions("DELETE", u, nil, nil, nil, o)
}
