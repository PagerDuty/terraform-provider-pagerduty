package pagerduty

import (
	"errors"
	"fmt"
)

var (
	// ErrNoToken is returned by NewClient if a user
	// passed an empty/missing token.
	ErrNoToken = errors.New("an empty token was provided")

	// ErrAuthFailure is returned by NewClient if a user
	// passed an invalid token and failed validation against the PagerDuty API.
	ErrAuthFailure = errors.New("failed to authenticate using the provided token")
)

type errorResponse struct {
	Error *Error `json:"error"`
}

// Error represents an error response from the PagerDuty API.
type Error struct {
	ErrorResponse  *Response
	Code           int         `json:"code,omitempty"`
	Errors         interface{} `json:"errors,omitempty"`
	Message        string      `json:"message,omitempty"`
	RequiredScopes string      `json:"required_scopes,omitempty"`
	TokenScopes    string      `json:"token_scopes,omitempty"`
	needToRetry    bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s API call to %s failed %v. Code: %d, Errors: %v, Message: %s", e.ErrorResponse.Response.Request.Method, e.ErrorResponse.Response.Request.URL.String(), e.ErrorResponse.Response.Status, e.Code, e.Errors, e.Message)
}
