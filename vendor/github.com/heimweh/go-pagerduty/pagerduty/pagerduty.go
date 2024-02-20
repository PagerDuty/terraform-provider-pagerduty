package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/heimweh/go-pagerduty/persistentconfig"
	"github.com/spf13/afero"
)

const (
	defaultBaseURL                    = "https://api.pagerduty.com"
	defaultAppOauthTokenGenerationURL = "https://identity.pagerduty.com/oauth/token"
	defaultUserAgent                  = "heimweh/go-pagerduty(terraform)"
	defaultRegion                     = "us"
	jitterPercent                     = 0.3
)

// AuthTokenType is an enum of available tokens types
// authenticating calls
type AuthTokenType int64

const (
	AuthTokenTypeAPIToken AuthTokenType = iota
	AuthTokenTypeScopedOauthToken
	AuthTokenTypeUseAppCredentials
)

func (d AuthTokenType) String() string {
	return authTokenTypeToStringMapping[d]
}

var authTokenTypeToStringMapping = map[AuthTokenType]string{
	AuthTokenTypeAPIToken:          "api_token",
	AuthTokenTypeScopedOauthToken:  "scoped_oauth_token",
	AuthTokenTypeUseAppCredentials: "use_app_credentials",
}

type service struct {
	client *Client
}

// Config represents the configuration for a PagerDuty client
type Config struct {
	BaseURL                   string
	HTTPClient                *http.Client
	Token                     string
	UserAgent                 string
	Debug                     bool
	APIAuthTokenType          *AuthTokenType
	AppOauthScopedTokenParams *persistentconfig.AppOauthScopedTokenParams
	clientPersistentConfig    *persistentconfig.ClientPersistentConfig
}

// Client manages the communication with the PagerDuty API
type Client struct {
	baseURL                          *url.URL
	client                           *http.Client
	Config                           *Config
	Abilities                        *AbilityService
	Addons                           *AddonService
	EscalationPolicies               *EscalationPolicyService
	Extensions                       *ExtensionService
	MaintenanceWindows               *MaintenanceWindowService
	Rulesets                         *RulesetService
	EventOrchestrations              *EventOrchestrationService
	EventOrchestrationPaths          *EventOrchestrationPathService
	EventOrchestrationIntegrations   *EventOrchestrationIntegrationService
	EventOrchestrationCacheVariables *EventOrchestrationCacheVariableService
	Schedules                        *ScheduleService
	Services                         *ServicesService
	Teams                            *TeamService
	ExtensionSchemas                 *ExtensionSchemaService
	Users                            *UserService
	Licenses                         *LicenseService
	Vendors                          *VendorService
	EventRules                       *EventRuleService
	BusinessServices                 *BusinessServiceService
	ServiceDependencies              *ServiceDependencyService
	Priorities                       *PriorityService
	ResponsePlays                    *ResponsePlayService
	SlackConnections                 *SlackConnectionService
	Tags                             *TagService
	WebhookSubscriptions             *WebhookSubscriptionService
	BusinessServiceSubscribers       *BusinessServiceSubscriberService
	OnCall                           *OnCallService
	AutomationActionsRunner          *AutomationActionsRunnerService
	AutomationActionsAction          *AutomationActionsActionService
	Incidents                        *IncidentService
	IncidentWorkflows                *IncidentWorkflowService
	IncidentWorkflowTriggers         *IncidentWorkflowTriggerService
	CustomFields                     *CustomFieldService
	CustomFieldSchemas               *CustomFieldSchemaService
	CustomFieldSchemaAssignments     *CustomFieldSchemaAssignmentService
	IncidentCustomFields             *IncidentCustomFieldService
}

// Response is a wrapper around http.Response
type Response struct {
	Response  *http.Response
	BodyBytes []byte
}

// RequestOptions is an object to setting options for HTTP requests
type RequestOptions struct {
	Type  string
	Label string
	Value string
}

// NewClient returns a new PagerDuty API client.
func NewClient(config *Config) (*Client, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultUserAgent
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}

	if config.APIAuthTokenType == nil {
		defaultTokenType := AuthTokenTypeAPIToken
		config.APIAuthTokenType = &defaultTokenType
	}

	if *config.APIAuthTokenType == AuthTokenTypeUseAppCredentials {
		clientPersistentConfig := persistentconfig.ClientPersistentConfig{
			Fs: afero.NewOsFs(), // Using host file system
		}
		if err := clientPersistentConfig.Load(); err != nil {
			return nil, err
		}
		config.AppOauthScopedTokenParams.Token = clientPersistentConfig.Token
		config.clientPersistentConfig = &clientPersistentConfig
	}

	c := &Client{
		baseURL: baseURL,
		client:  config.HTTPClient,
		Config:  config,
	}

	c.Abilities = &AbilityService{c}
	c.Addons = &AddonService{c}
	c.EscalationPolicies = &EscalationPolicyService{c}
	c.MaintenanceWindows = &MaintenanceWindowService{c}
	c.Rulesets = &RulesetService{c}
	c.EventOrchestrations = &EventOrchestrationService{c}
	c.EventOrchestrationIntegrations = &EventOrchestrationIntegrationService{c}
	c.EventOrchestrationCacheVariables = &EventOrchestrationCacheVariableService{c}
	c.EventOrchestrationPaths = &EventOrchestrationPathService{c}
	c.Schedules = &ScheduleService{c}
	c.Services = &ServicesService{c}
	c.Teams = &TeamService{c}
	c.Users = &UserService{c}
	c.Licenses = &LicenseService{c}
	c.Vendors = &VendorService{c}
	c.Extensions = &ExtensionService{c}
	c.ExtensionSchemas = &ExtensionSchemaService{c}
	c.EventRules = &EventRuleService{c}
	c.BusinessServices = &BusinessServiceService{c}
	c.ServiceDependencies = &ServiceDependencyService{c}
	c.Priorities = &PriorityService{c}
	c.ResponsePlays = &ResponsePlayService{c}
	c.SlackConnections = &SlackConnectionService{c}
	c.Tags = &TagService{c}
	c.WebhookSubscriptions = &WebhookSubscriptionService{c}
	c.BusinessServiceSubscribers = &BusinessServiceSubscriberService{c}
	c.OnCall = &OnCallService{c}
	c.AutomationActionsRunner = &AutomationActionsRunnerService{c}
	c.AutomationActionsAction = &AutomationActionsActionService{c}
	c.Incidents = &IncidentService{c}
	c.IncidentWorkflows = &IncidentWorkflowService{c}
	c.IncidentWorkflowTriggers = &IncidentWorkflowTriggerService{c}
	c.CustomFields = &CustomFieldService{c}
	c.CustomFieldSchemas = &CustomFieldSchemaService{c}
	c.CustomFieldSchemaAssignments = &CustomFieldSchemaAssignmentService{c}
	c.IncidentCustomFields = &IncidentCustomFieldService{c}

	InitCache(c)
	PopulateCache()

	return c, nil
}

func (c *Client) newRequest(method, url string, body interface{}, options ...RequestOptions) (*http.Request, error) {
	return c.newRequestContext(context.Background(), method, url, body, options...)
}

func (c *Client) newRequestContext(ctx context.Context, method, url string, body interface{}, options ...RequestOptions) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	if c.Config.Debug {
		log.Printf("[DEBUG] PagerDuty - Preparing %s request to %s with body: %s", method, url, buf)
	}

	u := c.baseURL.String() + url

	req, err := http.NewRequestWithContext(ctx, method, u, buf)
	if err != nil {
		return nil, err
	}

	if len(options) > 0 {
		for _, o := range options {
			if o.Type == "header" {
				req.Header.Add(o.Label, o.Value)
			}
		}
	}
	req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.Config.UserAgent)

	// Defaults to API Token Authorization header configuration
	authHeader := fmt.Sprintf("Token token=%s", c.Config.Token)
	if *c.Config.APIAuthTokenType == AuthTokenTypeUseAppCredentials || *c.Config.APIAuthTokenType == AuthTokenTypeScopedOauthToken {
		log.Printf("[INFO] Pagerduty - Using Scoped Oauth")
		authHeader = fmt.Sprintf("Bearer %s", c.Config.AppOauthScopedTokenParams.Token)
	}
	req.Header.Add("Authorization", authHeader)

	return req, nil
}

type scopedOauthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *Client) generateScopedOauthAccessToken() error {
	aotp := c.Config.AppOauthScopedTokenParams
	region := aotp.Region
	if region == "" {
		log.Printf("[INFO] Pagerduty - Using default region %q", defaultRegion)
		region = defaultRegion
	}
	subdomain := aotp.PDSubDomain
	u := defaultAppOauthTokenGenerationURL
	scopes := strings.Join(availableOauthScopes(), " ")

	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", aotp.ClientID)
	data.Add("client_secret", aotp.ClientSecret)
	data.Add("scope", fmt.Sprintf("as_account-%s.%s %s", region, subdomain, scopes))
	encodedData := data.Encode()
	payload := strings.NewReader(encodedData)

	req, err := http.NewRequest("POST", u, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", c.Config.UserAgent)

	internalClient := &http.Client{}

	v := new(scopedOauthResponse)
	resp, err := internalClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		return fmt.Errorf("with status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return err
	}

	// clientPersistentConfig := persistentconfig.ClientPersistentConfig{}
	// err = clientPersistentConfig.Load()
	// if err != nil {
	// 	return err
	// }
	c.Config.clientPersistentConfig.SetCredential("token", v.AccessToken)
	c.Config.AppOauthScopedTokenParams.Token = v.AccessToken

	return nil
}

func (c *Client) newRequestDo(method, url string, qryOptions, body, v interface{}) (*Response, error) {
	return c.newRequestDoContext(context.Background(), method, url, qryOptions, body, v)
}

func (c *Client) newRequestDoContext(ctx context.Context, method, url string, qryOptions, body, v interface{}) (*Response, error) {
	if qryOptions != nil {
		values, err := query.Values(qryOptions)
		if err != nil {
			return nil, err
		}

		if v := values.Encode(); v != "" {
			url = fmt.Sprintf("%s?%s", url, v)
		}
	}
	req, err := c.newRequestContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, v)
	if err != nil {
		if respErr, ok := err.(*Error); ok && respErr.needToRetry {
			return c.newRequestDoContext(ctx, method, url, nil, body, v)
		}

		return nil, err
	}

	return resp, nil
}

func (c *Client) newRequestDoOptions(method, url string, qryOptions, body, v interface{}, reqOptions ...RequestOptions) (*Response, error) {
	return c.newRequestDoOptionsContext(context.Background(), method, url, qryOptions, body, v, reqOptions...)
}

func (c *Client) newRequestDoOptionsContext(ctx context.Context, method, url string, qryOptions, body, v interface{}, reqOptions ...RequestOptions) (*Response, error) {
	if qryOptions != nil {
		values, err := query.Values(qryOptions)
		if err != nil {
			return nil, err
		}

		if v := values.Encode(); v != "" {
			url = fmt.Sprintf("%s?%s", url, v)
		}
	}
	req, err := c.newRequestContext(ctx, method, url, body, reqOptions...)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req, v)
	if err != nil {
		if respErr, ok := err.(*Error); ok && respErr.needToRetry {
			return c.newRequestDoOptionsContext(ctx, method, url, nil, body, v)
		}

		return nil, err
	}

	return resp, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	sLogger := newSecureLogger()
	sLogger.LogReq(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	sLogger.LogRes(resp)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := &Response{
		Response:  resp,
		BodyBytes: bodyBytes,
	}

	if err := c.checkResponse(response); err != nil {
		return response, err
	}

	if v != nil {
		if err := c.DecodeJSON(response, v); err != nil {
			return response, err
		}
	}

	return response, nil
}

// ListResp represents a list response from the PagerDuty API
type ListResp struct {
	Offset int  `json:"offset,omitempty"`
	Limit  int  `json:"limit,omitempty"`
	More   bool `json:"more,omitempty"`
	Total  int  `json:"total,omitempty"`
}

// responseHandler is capable of parsing a response. At a minimum it must
// extract the page information for the current page. It can also execute
// additional necessary handling; for example, if a closure, it has access
// to the scope in which it was defined, and can be used to append data to
// a specific slice. The responseHandler is responsible for closing the response.
type responseHandler func(response *Response) (ListResp, *Response, error)

// CursorListResp represents a cursor-paginated list response from the PagerDuty API
type CursorListResp struct {
	NextCursor string
	Limit      int
}

// cursorResponseHandler is capable of parsing a response. At a minimum it must
// extract the page information for the current page. It can also execute
// additional necessary handling; for example, if a closure, it has access
// to the scope in which it was defined, and can be used to append data to
// a specific slice. The responseHandler is responsible for closing the response.
type cursorResponseHandler func(response *Response) (CursorListResp, *Response, error)

// offsetQueryOptionsGen enables updating the offset across multiple
// pages of a list response while retaining query parameters other than the
// offset. For list-style functions where the underlying method supports
// other query parameters, this interface should be implemented to update
// the offset when fetching subsequent pages of results.
type offsetQueryOptionsGen interface {
	currentOffset() int
	changeOffset(int)
	buildStruct() interface{}
}

type simpleOffsetQueryOptionsGen struct {
	Offset int `url:"offset,omitempty"`
}

func (o *simpleOffsetQueryOptionsGen) currentOffset() int {
	return o.Offset
}

func (o *simpleOffsetQueryOptionsGen) changeOffset(i int) {
	o.Offset = i
}

func (o *simpleOffsetQueryOptionsGen) buildStruct() interface{} {
	return o
}

func (c *Client) newRequestPagedGetDo(basePath string, handler responseHandler, reqOptions ...RequestOptions) error {
	return c.newRequestPagedGetQueryDo(basePath, handler, &simpleOffsetQueryOptionsGen{}, reqOptions...)
}

func (c *Client) newRequestPagedGetQueryDo(basePath string, handler responseHandler, qryOptions offsetQueryOptionsGen, reqOptions ...RequestOptions) error {
	return c.newRequestPagedGetQueryDoContext(context.Background(), basePath, handler, qryOptions, reqOptions...)
}

func (c *Client) newRequestPagedGetQueryDoContext(ctx context.Context, basePath string, handler responseHandler, qryOptions offsetQueryOptionsGen, reqOptions ...RequestOptions) error {
	// Indicates whether there are still additional pages associated with request.
	var stillMore bool

	// Offset to set for the next page request.
	var nextOffset int

	// While there are more pages, keep adjusting the offset to get all results.
	for stillMore, nextOffset = true, qryOptions.currentOffset(); stillMore; {
		qryOptions.changeOffset(nextOffset)
		response, err := c.newRequestDoOptionsContext(ctx, "GET", basePath, qryOptions.buildStruct(), nil, nil, reqOptions...)
		if err != nil {
			return err
		}

		// Call handler to extract page information and execute additional necessary handling.
		pageInfo, _, err := handler(response)
		if err != nil {
			return err
		}

		// Bump the offset as necessary and set whether more results exist.
		nextOffset = pageInfo.Offset + pageInfo.Limit
		stillMore = pageInfo.More
	}

	return nil
}

// cursorQueryOptionsGen enables updating the cursor across multiple
// pages of a list response while retaining query parameters other than the
// cursor. For list-style functions where the underlying method supports
// other query parameters, this interface should be implemented to update
// the cursor when fetching subsequent pages of results.
type cursorQueryOptionsGen interface {
	currentCursor() string
	changeCursor(string)
	buildStruct() interface{}
}

func (c *Client) newRequestCursorPagedGetQueryDo(basePath string, handler cursorResponseHandler, qryOptions cursorQueryOptionsGen, reqOptions ...RequestOptions) error {
	return c.newRequestCursorPagedGetQueryDoContext(context.Background(), basePath, handler, qryOptions, reqOptions...)
}

func (c *Client) newRequestCursorPagedGetQueryDoContext(ctx context.Context, basePath string, handler cursorResponseHandler, qryOptions cursorQueryOptionsGen, reqOptions ...RequestOptions) error {
	// Indicates whether there are still additional pages associated with request.
	var stillMore bool

	// Cursor to set for the next page request.
	var nextCursor string

	// While there are more pages, keep adjusting the offset to get all results.
	for stillMore, nextCursor = true, qryOptions.currentCursor(); stillMore; {
		qryOptions.changeCursor(nextCursor)
		response, err := c.newRequestDoOptionsContext(ctx, "GET", basePath, qryOptions.buildStruct(), nil, nil, reqOptions...)
		if err != nil {
			return err
		}

		// Call handler to extract page information and execute additional necessary handling.
		pageInfo, _, err := handler(response)
		if err != nil {
			return err
		}

		// Bump the offset as necessary and set whether more results exist.
		nextCursor = pageInfo.NextCursor
		stillMore = nextCursor != ""
	}

	return nil
}

// ValidateAuth validates a token against the PagerDuty API
func (c *Client) ValidateAuth() error {
	_, _, err := c.Abilities.List()
	return err
}

// DecodeJSON decodes json body to given interface
func (c *Client) DecodeJSON(res *Response, v interface{}) error {
	return json.Unmarshal(res.BodyBytes, v)
}

func (c *Client) checkResponse(res *Response) error {
	if res.Response.StatusCode >= 200 && res.Response.StatusCode <= 299 {
		return nil
	}

	return c.decodeErrorResponse(res)
}

func (c *Client) decodeErrorResponse(res *Response) error {
	// Try to decode error response or fallback with standard error
	v := &errorResponse{Error: &Error{ErrorResponse: res}}
	err := c.DecodeJSON(res, v)

	if handledError := handleRatelimitError(res, v); handledError != nil {
		return handledError
	}

	if handledError := c.handleScopedOAuthError(res, v); handledError != nil {
		return handledError
	}

	if err != nil {
		return fmt.Errorf("%s API call to %s failed: %v", res.Response.Request.Method, res.Response.Request.URL.String(), res.Response.Status)
	}
	log.Printf("[INFO] v.Error %+v", v.Error)

	return v.Error
}

func (c *Client) handleScopedOAuthError(res *Response, v *errorResponse) error {
	isUsingScopedAPITokenFromCredentials := *c.Config.APIAuthTokenType == AuthTokenTypeUseAppCredentials
	isOauthScopeMissing := isUsingScopedAPITokenFromCredentials && res.Response.StatusCode == http.StatusForbidden
	needNewOauthScopedAccessToken := isUsingScopedAPITokenFromCredentials && res.Response.StatusCode == http.StatusUnauthorized
	if isOauthScopeMissing {
		return fmt.Errorf("%s API call to %s failed because %s API scope is required", res.Response.Request.Method, res.Response.Request.URL.String(), v.Error.RequiredScopes)
	}
	if needNewOauthScopedAccessToken {
		err := c.generateScopedOauthAccessToken()
		if err != nil {
			return fmt.Errorf("API call to obtain a new Scoped Oauth Access Token failed: %v", err)
		}
		v.Error.needToRetry = true
		return v.Error
	}

	return nil
}

// handleRatelimitError will handle rate limit errors from responses with http
// code 429. Delaying retry based on ratelimit-reset recommended by PagerDuty
// https://developer.pagerduty.com/docs/72d3b724589e3-rest-api-rate-limits#reaching-the-limit
func handleRatelimitError(res *Response, v *errorResponse) error {
	var markErrorAsRetryable = func(waitFor time.Duration) error {
		reqMethod := res.Response.Request.Method
		reqEndpoint := res.Response.Request.URL
		log.Printf(
			"[INFO] Rate limit hit, throttling by %v seconds until next retry to %s: %s",
			strconv.FormatFloat(waitFor.Seconds(), 'f', 1, 64),
			strings.ToUpper(reqMethod),
			reqEndpoint)
		time.Sleep(waitFor)
		v.Error.needToRetry = true
		return v.Error
	}

	ratelimitReset := res.Response.Header.Get("ratelimit-reset")

	if res.Response.StatusCode == http.StatusTooManyRequests && ratelimitReset == "" {
		baseDelay := 5 * time.Second
		jitter := 1 + (jitterPercent * rand.Float64())
		waitFor := time.Duration(float64(baseDelay) * jitter)

		return markErrorAsRetryable(waitFor)
	}

	if res.Response.StatusCode == http.StatusTooManyRequests && ratelimitReset != "" {
		headerWaitSeconds, err := strconv.ParseInt(ratelimitReset, 10, 0)
		if err == nil {
			baseDelay := 500 * time.Millisecond
			headerWait := time.Duration(headerWaitSeconds) * time.Second
			jitter := 1 + (jitterPercent * rand.Float64())
			extraWait := time.Duration(float64(baseDelay) * jitter)

			waitFor := headerWait + extraWait

			return markErrorAsRetryable(waitFor)
		}
	}

	return nil
}

func availableOauthScopes() []string {
	return []string{
		"abilities.read",
		"addons.read",
		"addons.write",
		"analytics.read",
		"audit_records.read",
		"change_events.read",
		"change_events.write",
		"custom_fields.read",
		"custom_fields.write",
		"escalation_policies.read",
		"escalation_policies.write",
		"event_orchestrations.read",
		"event_orchestrations.write",
		"event_rules.read",
		"event_rules.write",
		"extension_schemas.read",
		"extensions.read",
		"extensions.write",
		"incident_workflows.read",
		"incident_workflows.write",
		"incident_workflows:instances.write",
		"incidents.read",
		"incidents.write",
		"licenses.read",
		"notifications.read",
		"oncalls.read",
		"priorities.read",
		"response_plays.read",
		"response_plays.write",
		"schedules.read",
		"schedules.write",
		"services.read",
		"services.write",
		"standards.read",
		"standards.write",
		"status_dashboards.read",
		"subscribers.read",
		"subscribers.write",
		"tags.read",
		"tags.write",
		"teams.read",
		"teams.write",
		"templates.read",
		"templates.write",
		"users.read",
		"users.write",
		"users:contact_methods.read",
		"users:contact_methods.write",
		"users:sessions.read",
		"users:sessions.write",
		"vendors.read",
	}
}
