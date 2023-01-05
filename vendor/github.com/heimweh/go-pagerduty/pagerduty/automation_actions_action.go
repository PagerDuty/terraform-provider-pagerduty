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

type AutomationActionsActionTeamAssociationPayload struct {
	Team *TeamReference `json:"team,omitempty"`
}

var automationActionsActionBaseUrl = "/automation_actions/actions"

// Create creates a new action
func (s *AutomationActionsActionService) Create(action *AutomationActionsAction) (*AutomationActionsAction, *Response, error) {
	u := automationActionsActionBaseUrl
	v := new(AutomationActionsActionPayload)

	resp, err := s.client.newRequestDoOptions("POST", u, nil, &AutomationActionsActionPayload{Action: action}, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Action, resp, nil
}

// Get retrieves information about an action.
func (s *AutomationActionsActionService) Get(id string) (*AutomationActionsAction, *Response, error) {
	u := fmt.Sprintf("%s/%s", automationActionsActionBaseUrl, id)
	v := new(AutomationActionsActionPayload)

	resp, err := s.client.newRequestDoOptions("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Action, resp, nil
}

// Update an existing action
func (s *AutomationActionsActionService) Update(ID string, action *AutomationActionsAction) (*AutomationActionsAction, *Response, error) {
	u := fmt.Sprintf("%s/%s", automationActionsActionBaseUrl, ID)
	v := new(AutomationActionsActionPayload)
	p := &AutomationActionsActionPayload{Action: action}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.Action, resp, nil
}

// Delete deletes an existing action.
func (s *AutomationActionsActionService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("%s/%s", automationActionsActionBaseUrl, id)

	return s.client.newRequestDoOptions("DELETE", u, nil, nil, nil)
}

// Associate an Automation Action with a team
func (s *AutomationActionsActionService) AssociateToTeam(actionID, teamID string) (*AutomationActionsActionTeamAssociationPayload, *Response, error) {
	u := fmt.Sprintf("%s/%s/teams", automationActionsActionBaseUrl, actionID)
	v := new(AutomationActionsActionTeamAssociationPayload)
	p := &AutomationActionsActionTeamAssociationPayload{
		Team: &TeamReference{ID: teamID, Type: "team_reference"},
	}

	resp, err := s.client.newRequestDoOptions("POST", u, nil, p, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Dissociate an Automation Action with a team
func (s *AutomationActionsActionService) DissociateToTeam(actionID, teamID string) (*Response, error) {
	u := fmt.Sprintf("%s/%s/teams/%s", automationActionsActionBaseUrl, actionID, teamID)

	return s.client.newRequestDoOptions("DELETE", u, nil, nil, nil)
}

// Gets the details of an Automation Action / team relation
func (s *AutomationActionsActionService) GetAssociationToTeam(actionID, teamID string) (*AutomationActionsActionTeamAssociationPayload, *Response, error) {
	u := fmt.Sprintf("%s/%s/teams/%s", automationActionsActionBaseUrl, actionID, teamID)
	v := new(AutomationActionsActionTeamAssociationPayload)

	resp, err := s.client.newRequestDoOptions("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
