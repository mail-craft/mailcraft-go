package mailcraft

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is returned for any non-2xx response from the MailCraft API.
//
// It covers both error shapes the API returns: {"error": {"type", "message"}}
// for business-rule failures (plan limits, suppressed recipients, ...), and
// {"message", "errors": {"field": [...]}} for validation failures (422).
type APIError struct {
	// StatusCode is the HTTP status, e.g. 402 or 422.
	StatusCode int
	// Message explains what went wrong.
	Message string
	// Type is the business-rule error type, e.g. "plan_limit_reached". Empty for validation errors.
	Type string
	// Errors holds field validation errors for 422 responses.
	Errors map[string][]string
}

func (e *APIError) Error() string {
	if e.Type != "" {
		return fmt.Sprintf("mailcraft: %d %s: %s", e.StatusCode, e.Type, e.Message)
	}
	return fmt.Sprintf("mailcraft: %d: %s", e.StatusCode, e.Message)
}

func newAPIError(status int, raw []byte) *APIError {
	var body struct {
		Error *struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
		Message *string             `json:"message"`
		Errors  map[string][]string `json:"errors"`
	}

	if json.Unmarshal(raw, &body) == nil {
		if body.Error != nil {
			return &APIError{StatusCode: status, Message: body.Error.Message, Type: body.Error.Type}
		}
		if body.Message != nil {
			return &APIError{StatusCode: status, Message: *body.Message, Errors: body.Errors}
		}
	}

	return &APIError{StatusCode: status, Message: http.StatusText(status)}
}
