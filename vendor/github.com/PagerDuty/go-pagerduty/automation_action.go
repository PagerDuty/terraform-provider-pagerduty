package pagerduty

import (
	"context"

	"github.com/google/go-querystring/query"
)

type ActionDataReference struct {
	Script                        *string `json:"script,omitempty"`
	InvocationCommand             *string `json:"invocation_command,omitempty"`
	ProcessAutomationJobID        *string `json:"process_automation_job_id,omitempty"`
	ProcessAutomationJobArguments *string `json:"process_automation_job_arguments,omitempty"`
	ProcessAutomationNodeFilter   *string `json:"process_automation_node_filter,omitempty"`
}

type Priviledges struct {
	Permissions []string `json:"permissions"`
}

type AutomationAction struct {
	APIObject
	Name                 string               `json:"name,omitempty"`
	Description          string               `json:"description,omitempty"`
	ActionType           string               `json:"action_type,omitempty"`
	ActionClassification *string              `json:"action_classification,omitempty"`
	Runner               string               `json:"runner,omitempty"`
	RunnerType           string               `json:"runner_type,omitempty"`
	Services             []APIObject          `json:"services,omitempty"`
	Teams                []APIObject          `json:"teams,omitempty"`
	ActionDataReference  *ActionDataReference `json:"action_data_reference,omitempty"`
	Priviledges          *Priviledges         `json:"priviledges,omitempty"`
	Metadata             interface{}          `json:"metadata,omitempty"`
	CreationTime         string               `json:"creation_time,omitempty"`
	ModifyTime           string               `json:"modify_time,omitempty"`
	LastRun              string               `json:"last_run,omitempty"`
	LastRunBy            *APIObject           `json:"last_run_by,omitempty"`
}

type AutomationActionResponse struct {
	Action AutomationAction `json:"action,omitempty"`
}

type CreateAutomationActionOptions AutomationActionResponse

func (c *Client) CreateAutomationActionWithContext(ctx context.Context, o CreateAutomationActionOptions) (*AutomationAction, error) {
	resp, err := c.post(ctx, "/automation_actions/actions", o, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Action, nil
}

func (c *Client) GetAutomationActionWithContext(ctx context.Context, id string) (*AutomationAction, error) {
	resp, err := c.get(ctx, "/automation_actions/actions/"+id, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Action, nil
}

func (c *Client) DeleteAutomationActionWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/automation_actions/actions/"+id)
	return err
}

type UpdateAutomationActionOptions AutomationActionResponse

func (c *Client) UpdateAutomationActionWithContext(ctx context.Context, o UpdateAutomationActionOptions) (*AutomationAction, error) {
	resp, err := c.put(ctx, "/automation_actions/actions/"+o.Action.ID, o, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Action, nil
}

type AssociateAutomationActionServiceOptions struct {
	Service APIReference `json:"service,omitempty"`
}

type AutomationActionServiceResponse struct {
	Service APIObject `json:"service,omitempty"`
}

func (c *Client) AssociateAutomationActionServiceWithContext(ctx context.Context, id string, o AssociateAutomationActionServiceOptions) (*AutomationActionServiceResponse, error) {
	resp, err := c.post(ctx, "/automation_actions/actions/"+id+"/services", o, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionServiceResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DisassociateAutomationActionServiceWithContext(ctx context.Context, actionID, serviceID string) error {
	_, err := c.delete(ctx, "/automation_actions/actions/"+actionID+"/services/"+serviceID)
	return err
}

func (c *Client) GetAutomationActionServiceWithContext(ctx context.Context, actionID, serviceID string) (*AutomationActionServiceResponse, error) {
	resp, err := c.get(ctx, "/automation_actions/actions/"+actionID+"/services/"+serviceID, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionServiceResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type AssociateAutomationActionTeamOptions struct {
	Team APIReference `json:"team,omitempty"`
}

type AutomationActionTeamResponse struct {
	Team APIObject `json:"team,omitempty"`
}

func (c *Client) AssociateAutomationActionTeamWithContext(ctx context.Context, id string, o AssociateAutomationActionTeamOptions) (*AutomationActionTeamResponse, error) {
	resp, err := c.post(ctx, "/automation_actions/actions/"+id+"/teams", o, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionTeamResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DisassociateAutomationActionTeamWithContext(ctx context.Context, actionID, teamID string) error {
	_, err := c.delete(ctx, "/automation_actions/actions/"+actionID+"/teams/"+teamID)
	return err
}

func (c *Client) GetAutomationActionTeamWithContext(ctx context.Context, actionID, teamID string) (*AutomationActionTeamResponse, error) {
	resp, err := c.get(ctx, "/automation_actions/actions/"+actionID+"/teams/"+teamID, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionTeamResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type AutomationActionsRunnerActions struct {
	Actions []APIObject `json:"actions,omitempty"`
	More    bool        `json:"more,omitempty"`
}

type AutomationActionsRunner struct {
	APIObject
	Secret            string                          `json:"secret,omitempty"`
	RunnerType        string                          `json:"runner_type,omitempty"`
	Name              string                          `json:"name,omitempty"`
	Description       string                          `json:"description,omitempty"`
	LastSeen          string                          `json:"last_seen,omitempty"`
	Status            string                          `json:"status,omitempty"`
	CreationTime      string                          `json:"creation_time,omitempty"`
	RunbookBaseURI    string                          `json:"runbook_base_uri,omitempty"`
	RunbookAPIKey     string                          `json:"runbook_api_key,omitempty"`
	Teams             []APIObject                     `json:"teams,omitempty"`
	Priviledges       *Priviledges                    `json:"priviledges,omitempty"`
	AssociatedActions *AutomationActionsRunnerActions `json:"associated_actions,omitempty"`
	Metadata          interface{}                     `json:"metadata,omitempty"`
}

type AutomationActionsRunnerResponse struct {
	Runner AutomationActionsRunner `json:"runner,omitempty"`
}

func (c *Client) CreateAutomationActionsRunnerWithContext(ctx context.Context, a AutomationActionsRunner) (*AutomationActionsRunner, error) {
	d := map[string]AutomationActionsRunner{
		"runner": a,
	}

	resp, err := c.post(ctx, "/automation_actions/runners/", d, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionsRunnerResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Runner, nil
}

type ListAutomationActionsRunnersOptions struct {
	Cursor   string   `url:"include,omitempty"`
	Limit    uint     `url:"limit,omitempty"`
	Name     string   `url:"name,omitempty"`
	Includes []string `url:"include,brackets,omitempty"`
}

type ListAutomationActionsRunnersResponse struct {
	Runners     []AutomationActionsRunner `json:"runners,omitempty"`
	Priviledges *Priviledges              `json:"priviledges,omitempty"`
	cursor
	// NextCursor  string                    `json:"next_cursor,omitempty"`
	// Limit       uint                      `json:"limit,omitempty"`
	cursorHandler
}

func (c *Client) ListAutomationActionsRunnersWithContext(ctx context.Context, o ListAutomationActionsRunnersOptions) (*ListAutomationActionsRunnersResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, "/automation_actions/runners?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var response ListAutomationActionsRunnersResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetAutomationActionsRunnerWithContext(ctx context.Context, id string) (*AutomationActionsRunner, error) {
	resp, err := c.get(ctx, "/automation_actions/runners/"+id, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionsRunnerResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Runner, nil
}

func (c *Client) UpdateAutomationActionsRunnerWithContext(ctx context.Context, a AutomationActionsRunner) (*AutomationActionsRunner, error) {
	d := map[string]AutomationActionsRunner{
		"runner": a,
	}

	resp, err := c.put(ctx, "/automation_actions/runners/"+a.ID, d, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionsRunnerResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response.Runner, nil
}

func (c *Client) DeleteAutomationActionsRunnerWithContext(ctx context.Context, id string) error {
	_, err := c.delete(ctx, "/automation_actions/runners/"+id)
	return err
}

type AssociateAutomationActionsRunnerTeamOptions struct {
	Team APIReference `json:"team,omitempty"`
}

type AutomationActionsRunnerTeamResponse struct {
	Team APIObject `json:"team,omitempty"`
}

func (c *Client) AssociateAutomationActionsRunnerTeamWithContext(ctx context.Context, id string, o AssociateAutomationActionsRunnerTeamOptions) (*AutomationActionsRunnerTeamResponse, error) {
	resp, err := c.post(ctx, "/automation_actions/runners/"+id+"/teams", o, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionsRunnerTeamResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DisassociateAutomationActionsRunnerTeamWithContext(ctx context.Context, actionID, teamID string) error {
	_, err := c.delete(ctx, "/automation_actions/runners/"+actionID+"/teams/"+teamID)
	return err
}

func (c *Client) GetAutomationActionsRunnerTeamWithContext(ctx context.Context, actionID, teamID string) (*AutomationActionsRunnerTeamResponse, error) {
	resp, err := c.get(ctx, "/automation_actions/runners/"+actionID+"/teams/"+teamID, nil)
	if err != nil {
		return nil, err
	}

	var response AutomationActionsRunnerTeamResponse
	if err := c.decodeJSON(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
