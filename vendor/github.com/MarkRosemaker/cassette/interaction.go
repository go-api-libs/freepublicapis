package cassette

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// Interactions represents a collection of interactions.
type Interactions []*Interaction

// Interaction represents a single observed HTTP request/response pair.
type Interaction struct {
	Request  Request  `json:"request"`
	Response Response `json:"response,omitzero"`
}

// Request represents an observed HTTP request.
type Request struct {
	Method string `json:"method"`
	// URL is the URL requested, its query parameters sorted by name. A scaffold, a request yet to
	// be made, has none.
	URL     url.URL     `json:"url"`
	Headers http.Header `json:"header,omitempty"`
	Body    Body        `json:"body,omitzero"`
	// BodyOmitted is set for a body that was not text, see [IsText], which is not recorded.
	BodyOmitted bool `json:"bodyOmitted,omitzero"`
}

// NewRequest creates a new [Request] out of an [*http.Request].
// If the request has a body of text, it is drained and restored; any other is left as it is, and recorded as omitted.
func NewRequest(req *http.Request) (Request, error) {
	r := Request{
		Method:  req.Method,
		URL:     normalizeURL(*req.URL),
		Headers: req.Header.Clone(),
	}

	if req.Body == nil || req.Body == http.NoBody {
		return r, nil
	}

	if !IsText(req.Header.Get("Content-Type")) {
		r.BodyOmitted = true
		return r, nil
	}

	// Drain and restore the request body.
	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		return r, err
	}

	r.Body = body
	req.Body = io.NopCloser(bytes.NewReader(body))

	return r, nil
}

// Create creates a corresponding [*http.Request].
func (r Request) Create(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL.String(), func() io.Reader {
		if r.Body == nil {
			return nil
		}

		return bytes.NewReader(r.Body)
	}())
	if err != nil {
		return nil, err
	}

	if r.Headers != nil {
		req.Header = r.Headers.Clone()
	}

	return req, nil
}

// IsScaffold reports whether r is a request yet to be made: one without a URL.
func (r Request) IsScaffold() bool { return r.URL == (url.URL{}) }

// normalizeURL returns u with its query parameters stably sorted by name, so that two requests that differ only in
// their order are recorded alike. Each parameter is kept as it was written: re-escaping it, as [url.Values.Encode]
// does, would turn a delimiter such as the comma of ids=1,2 into %2C.
func normalizeURL(u url.URL) url.URL {
	if u.RawQuery == "" {
		return u
	}

	params := strings.Split(u.RawQuery, "&")
	slices.SortStableFunc(params, func(a, b string) int { return cmp.Compare(queryName(a), queryName(b)) })
	u.RawQuery = strings.Join(params, "&")

	return u
}

// queryName is the unescaped name of the query parameter param, written as name=value.
func queryName(param string) string {
	name, _, _ := strings.Cut(param, "=")
	if unescaped, err := url.QueryUnescape(name); err == nil {
		return unescaped
	}

	return name
}

// Response represents an observed HTTP response.
type Response struct {
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"header,omitempty"`
	Body       Body        `json:"body,omitzero"`
	// BodyOmitted is set for a body that was not text, see [IsText], which is not recorded.
	BodyOmitted bool `json:"bodyOmitted,omitzero"`
}

// NewResponse creates a new [Response] out of an [*http.Response].
// If the response has a body of text, it is drained and restored; any other, such as a zip, is left to stream as it
// is, and recorded as omitted.
func NewResponse(resp *http.Response) (Response, error) {
	r := Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
	}

	if !IsText(resp.Header.Get("Content-Type")) {
		r.BodyOmitted = resp.Body != nil && resp.Body != http.NoBody
		return r, nil
	}

	// Drain and restore the response body.
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return r, err
	}

	r.Body = body
	resp.Body = io.NopCloser(bytes.NewReader(body))

	return r, nil
}

// Body is a recorded body: written as the JSON it is, or as a string if it is no JSON.
type Body []byte

// UnmarshalJSONFrom implements [json.UnmarshalerFrom]: a body of text, such as XML or HTML, written as a JSON string
// is read back as the text it is, as it was sent, not as the JSON string it was written as.
func (r *Request) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	type plain Request

	if err := json.UnmarshalDecode(dec, (*plain)(r)); err != nil {
		return err
	}

	r.Body = r.Body.asSent(r.Headers.Get("Content-Type"))

	return nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom], reading a body of text as [Request.UnmarshalJSONFrom] does.
func (rsp *Response) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	type plain Response

	if err := json.UnmarshalDecode(dec, (*plain)(rsp)); err != nil {
		return err
	}

	rsp.Body = rsp.Body.asSent(rsp.Headers.Get("Content-Type"))

	return nil
}

// asSent is b as it was sent with contentType: a JSON string, which is how a body of text other than JSON is written,
// is the text it holds. A body of JSON, or of no content type said, is the JSON it is.
func (b Body) asSent(contentType string) Body {
	if contentType == "" || isJSON(contentType) || jsontext.Value(b).Kind() != '"' {
		return b
	}

	text, err := jsontext.AppendUnquote(nil, []byte(b))
	if err != nil {
		return b
	}

	return text
}

// IsZero reports whether there is no body, which is left out when written: an empty JSON value, such as null or {},
// is a body all the same.
func (b Body) IsZero() bool { return len(b) == 0 }

// MarshalJSONTo implements [json.MarshalerTo].
func (b Body) MarshalJSONTo(enc *jsontext.Encoder) error {
	if jsontext.Value(b).IsValid() {
		return enc.WriteValue(jsontext.Value(b))
	}

	return enc.WriteToken(jsontext.String(string(b)))
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
func (b *Body) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	*b = Body(val.Clone())

	return nil
}
