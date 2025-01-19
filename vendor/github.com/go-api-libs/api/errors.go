package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrUnknownStatusCode is returned when the response status code is not known,
	// i.e. it has not been defined in the OpenAPI specification.
	ErrUnknownStatusCode = errors.New("unknown status code")
	// ErrUnknownContentType is returned when the response content type is not known,
	// i.e. it has not been defined in the OpenAPI specification
	// for the status code of the response.
	ErrUnknownContentType = errors.New("unknown content type")
	// ErrStatusCode is returned when the status code indicates an error.
	ErrStatusCode = errors.New("status code indicates an error")
)

// Error is an error returned by the API client.
// Use errors.As to check if the underlying error is of a specific type.
type Error struct {
	// The HTTP response.
	Response *http.Response
	// The underlying error.
	Err error
	// Indicates that this is a custom error returned by the API.
	IsCustom bool
}

func NewErrUnknownStatusCode(rsp *http.Response) error {
	return &Error{Response: rsp, Err: ErrUnknownStatusCode}
}

func NewErrStatusCode(rsp *http.Response) error {
	return &Error{Response: rsp, Err: ErrStatusCode}
}

func NewErrUnknownContentType(rsp *http.Response) error {
	return &Error{Response: rsp, Err: ErrUnknownContentType}
}

// Error returns the error message.
func (e *Error) Error() string {
	switch e.Err {
	case ErrUnknownContentType:
		return fmt.Sprintf("%s: %v %q", e.Response.Status, e.Err,
			e.Response.Header.Get("Content-Type"))
	case ErrStatusCode:
		return fmt.Sprintf("got %s calling %s", e.Response.Status, e.Response.Request.URL)
	default:
		return fmt.Sprintf("%s: %v", e.Response.Status, e.Err)
	}
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error { return e.Err }

// StatusCode returns the HTTP status code of the response.
func (e Error) StatusCode() int { return e.Response.StatusCode }

// ContentType returns the HTTP content type.
func (e Error) ContentType() string { return e.Response.Header.Get("Content-Type") }

// DecodingError is an error returned when decoding the response body failed.
type DecodingError struct{ Err error }

// Unwrap returns the underlying error, e.g. a json.SemanticError.
func (e *DecodingError) Unwrap() error { return e.Err }

// Error returns the error message.
func (d DecodingError) Error() string {
	return fmt.Sprintf("decoding response body: %v", d.Err)
}

// WrapDecodingError wraps the error in a DecodingError,
// which is then wrapped into an api.Error containing the response.
// Use errors.As to check if the error returned by an API library is a DecodingError.
func WrapDecodingError(rsp *http.Response, err error) error {
	return &Error{Response: rsp, Err: &DecodingError{Err: err}}
}

// NewErrCustom wraps a custom error in an api.Error.
// Use this if the status code is unsuccessful and the API returns more information,
// e.g. as part of a JSON that can be parsed and transformed into an error.
func NewErrCustom(rsp *http.Response, err error) error {
	return &Error{Response: rsp, Err: err, IsCustom: true}
}
