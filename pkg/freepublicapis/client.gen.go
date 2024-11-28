// Package freepublicapis provides a client and primitives to interact with the HTTP API.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.
package freepublicapis

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/go-api-libs/api"
	"github.com/go-json-experiment/json"
)

var (
	baseURL = &url.URL{
		Host:   "www.freepublicapis.com",
		Path:   "/api",
		Scheme: "https",
	}

	jsonOpts = json.JoinOptions(
		json.WithMarshalers(json.MarshalFuncV2(jsonutil.URLMarshal)),
		json.WithUnmarshalers(json.UnmarshalFuncV2(jsonutil.URLUnmarshal)))
)

// Client conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The HTTP client to use for requests.
	cli *http.Client
}

// NewClient creates a new Client.
func NewClient() (*Client, error) {
	return &Client{cli: http.DefaultClient}, nil
}

// GetRandom defines an operation.
// GET /random
func (c *Client) GetRandom(ctx context.Context) (*SimpleAPIInfo, error) {
	return GetRandom[SimpleAPIInfo](ctx, c)
}

// GetRandom defines an operation.
// You can define a custom result to unmarshal the response into.
// GET /random
func GetRandom[R any](ctx context.Context, c *Client) (*R, error) {
	u := baseURL.JoinPath("/random")
	req := (&http.Request{
		Header:     http.Header{},
		Host:       u.Host,
		Method:     http.MethodGet,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        u,
	}).WithContext(ctx)

	rsp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusOK:
		// Returns simple information about a single API
		switch rsp.Header.Get("Content-Type") {
		case "application/json":
			var out R
			if err := json.UnmarshalRead(rsp.Body, &out, jsonOpts); err != nil {
				return nil, api.WrapDecodingError(rsp, err)
			}

			return &out, nil
		default:
			return nil, api.NewErrUnknownContentType(rsp)
		}
	default:
		return nil, api.NewErrUnknownStatusCode(rsp)
	}
}

// GetAPI defines an operation.
// GET /apis/{id}
func (c *Client) GetAPI(ctx context.Context, id int) (*SimpleAPIInfo, error) {
	return GetAPI[SimpleAPIInfo](ctx, c, id)
}

// GetAPI defines an operation.
// You can define a custom result to unmarshal the response into.
// GET /apis/{id}
func GetAPI[R any](ctx context.Context, c *Client, id int) (*R, error) {
	u := baseURL.JoinPath("apis", strconv.Itoa(id))
	req := (&http.Request{
		Header:     http.Header{},
		Host:       u.Host,
		Method:     http.MethodGet,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        u,
	}).WithContext(ctx)

	rsp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusOK:
		// Returns simple information about a single API
		switch rsp.Header.Get("Content-Type") {
		case "application/json":
			var out R
			if err := json.UnmarshalRead(rsp.Body, &out, jsonOpts); err != nil {
				return nil, api.WrapDecodingError(rsp, err)
			}

			return &out, nil
		default:
			return nil, api.NewErrUnknownContentType(rsp)
		}
	default:
		return nil, api.NewErrUnknownStatusCode(rsp)
	}
}

// ListApis defines an operation.
// GET /apis
func (c *Client) ListApis(ctx context.Context, params *ListApisParams) (APIInfos, error) {
	return ListApis[APIInfos](ctx, c, params)
}

// ListApis defines an operation.
// You can define a custom result to unmarshal the response into.
// GET /apis
func ListApis[R any](ctx context.Context, c *Client, params *ListApisParams) (R, error) {
	u := baseURL.JoinPath("/apis")

	if params != nil {
		q := make(url.Values, 2)

		if params.Limit != nil {
			q["limit"] = []string{strconv.Itoa(*params.Limit)}
		}

		if params.Sort != nil {
			q["sort"] = []string{*params.Sort}
		}

		u.RawQuery = q.Encode()
	}

	req := (&http.Request{
		Header:     http.Header{},
		Host:       u.Host,
		Method:     http.MethodGet,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        u,
	}).WithContext(ctx)

	var out R
	rsp, err := c.cli.Do(req)
	if err != nil {
		return out, err
	}
	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusOK:
		// Returns a list of API information
		switch rsp.Header.Get("Content-Type") {
		case "application/json":
			if err := json.UnmarshalRead(rsp.Body, &out, jsonOpts); err != nil {
				return out, api.WrapDecodingError(rsp, err)
			}

			return out, nil
		default:
			return out, api.NewErrUnknownContentType(rsp)
		}
	default:
		return out, api.NewErrUnknownStatusCode(rsp)
	}
}
