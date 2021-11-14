package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.pagerduty.com"
)

type service struct {
	client *Client
}

// Config represents the configuration for a PagerDuty client
type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	UserAgent  string
	Debug      bool
}

// Client manages the communication with the PagerDuty API
type Client struct {
	baseURL             *url.URL
	client              *http.Client
	Config              *Config
	Abilities           *AbilityService
	Addons              *AddonService
	EscalationPolicies  *EscalationPolicyService
	Extensions          *ExtensionService
	MaintenanceWindows  *MaintenanceWindowService
	Rulesets            *RulesetService
	Schedules           *ScheduleService
	Services            *ServicesService
	Teams               *TeamService
	ExtensionSchemas    *ExtensionSchemaService
	Users               *UserService
	Vendors             *VendorService
	EventRules          *EventRuleService
	BusinessServices    *BusinessServiceService
	ServiceDependencies *ServiceDependencyService
	Priorities          *PriorityService
	ResponsePlays       *ResponsePlayService
	SlackConnections    *SlackConnectionService
	Tags                *TagService
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

	config.UserAgent = "heimweh/go-pagerduty(terraform)"

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
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
	c.Schedules = &ScheduleService{c}
	c.Services = &ServicesService{c}
	c.Teams = &TeamService{c}
	c.Users = &UserService{c}
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

	InitCache(c)
	PopulateCache()

	return c, nil
}

func (c *Client) newRequest(method, url string, body interface{}, options ...RequestOptions) (*http.Request, error) {
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

	req, err := http.NewRequest(method, u, buf)
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
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", c.Config.Token))
	req.Header.Add("Content-Type", "application/json")

	if c.Config.UserAgent != "" {
		req.Header.Add("User-Agent", c.Config.UserAgent)
	}
	return req, nil
}

func (c *Client) newRequestDo(method, url string, qryOptions, body, v interface{}) (*Response, error) {
	if qryOptions != nil {
		values, err := query.Values(qryOptions)
		if err != nil {
			return nil, err
		}

		if v := values.Encode(); v != "" {
			url = fmt.Sprintf("%s?%s", url, v)
		}
	}
	req, err := c.newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

func (c *Client) newRequestDoOptions(method, url string, qryOptions, body, v interface{}, reqOptions ...RequestOptions) (*Response, error) {
	if qryOptions != nil {
		values, err := query.Values(qryOptions)
		if err != nil {
			return nil, err
		}

		if v := values.Encode(); v != "" {
			url = fmt.Sprintf("%s?%s", url, v)
		}
	}
	req, err := c.newRequest(method, url, body, reqOptions...)
	if err != nil {
		return nil, err
	}

	return c.do(req, v)
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
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

func (c *Client) newRequestPagedGetDo(basePath string, handler responseHandler, reqOptions ...RequestOptions) error {
	// Indicates whether there are still additional pages associated with request.
	var stillMore bool

	// Offset to set for the next page request.
	var nextOffset int

	// While there are more pages, keep adjusting the offset to get all results.
	for stillMore, nextOffset = true, 0; stillMore; {
		response, err := c.newRequestDoOptions("GET", fmt.Sprintf("%s?offset=%d", basePath, nextOffset), nil, nil, nil, reqOptions...)
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
	if err := c.DecodeJSON(res, v); err != nil {
		return fmt.Errorf("%s API call to %s failed: %v", res.Response.Request.Method, res.Response.Request.URL.String(), res.Response.Status)
	}

	return v.Error
}
