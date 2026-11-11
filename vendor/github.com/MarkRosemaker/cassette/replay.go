package cassette

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"
)

// ErrNoInteraction is returned, wrapped, for a request no recorded interaction answers.
var ErrNoInteraction = errors.New("no recorded interaction")

// Recorder is an [http.RoundTripper] that answers requests from recorded interactions, and, with [RecordMisses],
// sends and records those it has none for. Make one with [Replay]. It is safe for concurrent use.
type Recorder struct {
	mu sync.Mutex

	recorded Interactions // what it answers from, scaffolds left out
	used     []bool       // whether each of recorded has answered
	next     int          // the index of the next of recorded to answer, in order
	calls    Interactions // what it answered or recorded, in call order

	inOrder bool
	match   Matcher
	send    http.RoundTripper // where a request without a recording goes; nil: it fails
}

// Option configures a [Recorder].
type Option func(*Recorder)

// Replay returns a [Recorder] that answers each request from the first of ias it matches (see [MatchRequest]) that
// has not answered yet, or, once all have, the last of them again. A request none of ias matches fails with an error
// wrapping [ErrNoInteraction]; nothing reaches the network. Interactions without a response, scaffolds, are left out.
func Replay(ias Interactions, opts ...Option) *Recorder {
	r := &Recorder{match: MatchRequest()}

	for _, ia := range ias {
		if ia.Response.StatusCode > 0 {
			r.recorded = append(r.recorded, ia)
		}
	}

	r.used = make([]bool, len(r.recorded))

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// RecordMisses makes a [Recorder] send a request it has no recording for through rt, [http.DefaultTransport] if
// nil, and record it, instead of failing. With no recordings to start from, it records every request.
func RecordMisses(rt http.RoundTripper) Option {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return func(r *Recorder) { r.send = rt }
}

// InOrder makes a [Recorder] answer with its recordings in their order, each once: a request that does not match the
// next one fails, saying why.
func InOrder() Option { return func(r *Recorder) { r.inOrder = true } }

// WithMatcher makes a [Recorder] match requests to recordings with m instead of [MatchRequest]().
func WithMatcher(m Matcher) Option { return func(r *Recorder) { r.match = m } }

// RoundTrip implements [http.RoundTripper].
func (r *Recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	r.mu.Lock()

	ia, err := r.answer(req)
	if ia != nil {
		r.calls = append(r.calls, ia)
	}

	r.mu.Unlock()

	switch {
	case err != nil:
		return nil, err
	case ia != nil:
		return ia.Response.httpResponse(req), nil
	case r.send != nil:
		return r.record(req)
	default:
		return nil, fmt.Errorf("%s %s: %w", req.Method, req.URL, ErrNoInteraction)
	}
}

// answer is the recording that answers req, if any; err is set where none can, in order.
func (r *Recorder) answer(req *http.Request) (*Interaction, error) {
	if r.inOrder {
		if r.next >= len(r.recorded) {
			if r.send != nil {
				return nil, nil
			}

			return nil, fmt.Errorf("unexpected request %s %s: %w", req.Method, req.URL, ErrNoInteraction)
		}

		if err := r.match(req, r.recorded[r.next].Request); err != nil {
			return nil, fmt.Errorf("interaction #%d: %w", r.next, err)
		}

		r.used[r.next] = true
		r.next++

		return r.recorded[r.next-1], nil
	}

	last := -1

	for i, ia := range r.recorded {
		if r.match(req, ia.Request) != nil {
			continue
		}

		if !r.used[i] {
			r.used[i] = true
			return ia, nil
		}

		last = i
	}

	if last == -1 {
		return nil, nil
	}

	return r.recorded[last], nil
}

// record sends req and records what it is answered with.
func (r *Recorder) record(req *http.Request) (*http.Response, error) {
	rec, err := NewRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := r.send.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	recResp, err := NewResponse(resp)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.calls = append(r.calls, &Interaction{Request: rec, Response: recResp})
	r.mu.Unlock()

	return resp, nil
}

// Interactions are the interactions the recorder answered with or recorded, in the order of the calls.
func (r *Recorder) Interactions() Interactions {
	r.mu.Lock()
	defer r.mu.Unlock()

	return slices.Clone(r.calls)
}

// Unused are the recordings that have not answered any request.
func (r *Recorder) Unused() Interactions {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out Interactions

	for i, ia := range r.recorded {
		if !r.used[i] {
			out = append(out, ia)
		}
	}

	return out
}

// httpResponse is the response rsp records, to req.
func (rsp Response) httpResponse(req *http.Request) *http.Response {
	return &http.Response{
		Status:        fmt.Sprintf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode)),
		StatusCode:    rsp.StatusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        rsp.Headers.Clone(),
		Body:          io.NopCloser(bytes.NewReader(rsp.Body)),
		ContentLength: int64(len(rsp.Body)),
		Request:       req,
	}
}

// Matcher reports why the request req does not match the recorded request rec, or nil if it does.
type Matcher func(req *http.Request, rec Request) error

// MatchOption configures [MatchRequest].
type MatchOption func(*matcher)

type matcher struct {
	headers      bool
	ignore       map[string]bool
	optionalAuth func(*http.Request) bool
}

// MatchRequest returns the [Matcher] a [Recorder] uses unless told otherwise: a request matches a recording with the
// same method, the same URL, its query parameters in any order, and the same body, compared as
// JSON where both are, so that the order of members does not matter. Headers are not compared, unless
// [CompareHeaders] says so.
func MatchRequest(opts ...MatchOption) Matcher {
	m := &matcher{}
	for _, opt := range opts {
		opt(m)
	}

	return m.match
}

// CompareHeaders makes [MatchRequest] compare headers too, all but those named, exactly. Authorization is compared
// by its scheme alone, since a recording masks the credential after it, and a Content-Type recorded without a body
// is left out: it says nothing of a request that has none.
func CompareHeaders(ignore ...string) MatchOption {
	return func(m *matcher) {
		m.headers = true
		m.ignore = map[string]bool{}

		for _, name := range ignore {
			m.ignore[http.CanonicalHeaderKey(name)] = true
		}
	}
}

// OptionalAuthorization makes [MatchRequest] match a recording without Authorization to a request that has one where
// optional reports true for it: a call whose credentials are optional may be sent with them though it was recorded
// without.
func OptionalAuthorization(optional func(*http.Request) bool) MatchOption {
	return func(m *matcher) { m.optionalAuth = optional }
}

func (m *matcher) match(req *http.Request, rec Request) error {
	got, err := NewRequest(req)
	if err != nil {
		return err
	}

	if got.Method != rec.Method {
		return fmt.Errorf("got method %s, want %s", got.Method, rec.Method)
	}

	if want := normalizeURL(rec.URL); got.URL.String() != want.String() {
		return fmt.Errorf("got URL %s, want %s", &got.URL, &want)
	}

	if !bytes.Equal(canonical(got.Body), canonical(rec.Body)) {
		return fmt.Errorf("got body %s, want %s", got.Body, rec.Body)
	}

	if m.headers {
		return m.matchHeaders(req, got.Headers.Clone(), rec)
	}

	return nil
}

// matchHeaders compares the headers of req, got, with those rec was recorded with, as [CompareHeaders] says.
func (m *matcher) matchHeaders(req *http.Request, got http.Header, rec Request) error {
	want := rec.Headers.Clone()
	if want == nil {
		want = http.Header{}
	}

	if got == nil {
		got = http.Header{}
	}

	if len(rec.Body) == 0 {
		want.Del("Content-Type")
	}

	gotScheme, _, _ := strings.Cut(got.Get("Authorization"), " ")
	wantScheme, _, _ := strings.Cut(want.Get("Authorization"), " ")

	if wantScheme == "" && m.optionalAuth != nil && m.optionalAuth(req) {
		wantScheme = gotScheme
	}

	if gotScheme != wantScheme {
		return fmt.Errorf("got Authorization scheme %q, want %q", gotScheme, wantScheme)
	}

	for name := range m.ignore {
		got.Del(name)
		want.Del(name)
	}

	got.Del("Authorization")
	want.Del("Authorization")

	if !maps.EqualFunc(got, want, slices.Equal) {
		return fmt.Errorf("got headers %v, want %v", got, want)
	}

	return nil
}
