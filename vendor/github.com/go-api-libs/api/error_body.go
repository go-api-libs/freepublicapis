package api

import "net/http"

// maxBodyErrorLen bounds what Error renders. These bodies are whole HTML pages
// often enough that returning one unabridged would be unusable in a log line.
const maxBodyErrorLen = 512

// ErrorBody is an error carrying a response body the specification declares no
// schema for, so nothing could be decoded from it.
type ErrorBody struct {
	// The undecoded response body.
	Body []byte
}

// Error returns the error message: the body, truncated if it is long.
func (e *ErrorBody) Error() string {
	switch {
	case len(e.Body) == 0:
		return "empty response body"
	case len(e.Body) > maxBodyErrorLen:
		return string(e.Body[:maxBodyErrorLen]) + "..."
	default:
		return string(e.Body)
	}
}

// NewErrBody wraps an undecoded response body in an api.Error.
// Use this when the status code is unsuccessful and the API returns a body the
// specification gives no schema for, leaving nothing to decode it into.
func NewErrBody(rsp *http.Response, body []byte) error {
	return &Error{Response: rsp, Err: &ErrorBody{Body: body}, IsCustom: true}
}
