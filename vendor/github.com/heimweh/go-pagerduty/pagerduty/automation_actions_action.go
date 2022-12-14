package pagerduty

import "fmt"

// AutomationActionsAction handles the communication with schedule
// related methods of the PagerDuty API.
type AutomationActionsActionService service

// Associate an Automation Action with a team
func (s *AutomationActionsActionService) AssociateToTeam(actionID, teamID string) (*TeamReference, *Response, error) {
	u := fmt.Sprintf("/automation_actions/actions/%s/teams", actionID)
	v := new(struct {
		Team *TeamReference `json:"team,omitempty"`
	})
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	resp, err := s.client.newRequestDoOptions("POST", u, nil, &TeamReference{ID: teamID, Type: "team_reference"}, &v, o)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}

// Dissociate an Automation Action with a team
func (s *AutomationActionsActionService) DissociateToTeam(actionID, teamID string) (*Response, error) {
	u := fmt.Sprintf("/automation_actions/actions/%s/teams/%s", actionID, teamID)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	return s.client.newRequestDoOptions("DELETE", u, nil, nil, nil, o)
}

// Gets the details of an Automation Action / team relation
func (s *AutomationActionsActionService) GetAssociationToTeam(actionID, teamID string) (*TeamReference, *Response, error) {
	u := fmt.Sprintf("/automation_actions/actions/%s/teams/%s", actionID, teamID)
	v := new(struct {
		Team *TeamReference `json:"team,omitempty"`
	})
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	resp, err := s.client.newRequestDoOptions("GET", u, nil, nil, &v, o)
	if err != nil {
		return nil, nil, err
	}

	return v.Team, resp, nil
}
