// Package freepublicapis_test provides tests for the client package.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.
package freepublicapis_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-api-libs/api"
	"github.com/go-api-libs/freepublicapis/pkg/freepublicapis"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func pointerTo[T any](v T) *T {
	return &v
}

type testRoundTripper struct {
	rsp *http.Response
	err error
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.rsp, t.err
}

func TestClient_Error(t *testing.T) {
	ctx := context.Background()

	c, err := freepublicapis.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Do", func(t *testing.T) {
		testErr := errors.New("test error")
		http.DefaultClient.Transport = &testRoundTripper{err: testErr}

		if _, err := c.GetRandom(ctx); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, testErr) {
			t.Fatalf("want: %v, got: %v", testErr, err)
		}

		if _, err := c.GetAPI(ctx, 275); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, testErr) {
			t.Fatalf("want: %v, got: %v", testErr, err)
		}

		if _, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
			Limit: 10,
			Sort: "best",
		}); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, testErr) {
			t.Fatalf("want: %v, got: %v", testErr, err)
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		errDecode := &api.DecodingError{}

		t.Run("GetRandom", func(t *testing.T) {
			// unknown status code
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{StatusCode: http.StatusTeapot}}

			if _, err := c.GetRandom(ctx); err == nil {
				t.Fatal("expected error")
			} else if !errors.Is(err, api.ErrUnknownStatusCode) {
				t.Fatalf("want: %v, got: %v", api.ErrUnknownStatusCode, err)
			}

			// unknown content type
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{
				Header:     http.Header{"Content-Type": []string{"foo"}},
				StatusCode: http.StatusOK,
			}}

			if _, err := c.GetRandom(ctx); err == nil {
				t.Fatal("expected error")
			} else if !errors.Is(err, api.ErrUnknownContentType) {
				t.Fatalf("want: %v, got: %v", api.ErrUnknownContentType, err)
			}

			// decoding error for known content type "application/json"
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("{")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				StatusCode: http.StatusOK,
			}}

			if _, err := c.GetRandom(ctx); err == nil {
				t.Fatal("expected error")
			} else if !errors.As(err, &errDecode) {
				t.Fatalf("want: %v, got: %v", errDecode, err)
			}
		})

		t.Run("GetAPI", func(t *testing.T) {
			// unknown status code
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{StatusCode: http.StatusTeapot}}

			if _, err := c.GetAPI(ctx, 275); err == nil {
				t.Fatal("expected error")
			} else if !errors.Is(err, api.ErrUnknownStatusCode) {
				t.Fatalf("want: %v, got: %v", api.ErrUnknownStatusCode, err)
			}

			// unknown content type
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{
				Header:     http.Header{"Content-Type": []string{"foo"}},
				StatusCode: http.StatusOK,
			}}

			if _, err := c.GetAPI(ctx, 275); err == nil {
				t.Fatal("expected error")
			} else if !errors.Is(err, api.ErrUnknownContentType) {
				t.Fatalf("want: %v, got: %v", api.ErrUnknownContentType, err)
			}

			// decoding error for known content type "application/json"
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("{")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				StatusCode: http.StatusOK,
			}}

			if _, err := c.GetAPI(ctx, 275); err == nil {
				t.Fatal("expected error")
			} else if !errors.As(err, &errDecode) {
				t.Fatalf("want: %v, got: %v", errDecode, err)
			}
		})

		t.Run("ListApis", func(t *testing.T) {
			// unknown status code
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{StatusCode: http.StatusTeapot}}

			if _, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 10,
				Sort:  "best",
			}); err == nil {
				t.Fatal("expected error")
			} else if !errors.Is(err, api.ErrUnknownStatusCode) {
				t.Fatalf("want: %v, got: %v", api.ErrUnknownStatusCode, err)
			}

			// unknown content type
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{
				Header:     http.Header{"Content-Type": []string{"foo"}},
				StatusCode: http.StatusOK,
			}}

			if _, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 10,
				Sort:  "best",
			}); err == nil {
				t.Fatal("expected error")
			} else if !errors.Is(err, api.ErrUnknownContentType) {
				t.Fatalf("want: %v, got: %v", api.ErrUnknownContentType, err)
			}

			// decoding error for known content type "application/json"
			http.DefaultClient.Transport = &testRoundTripper{rsp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("{")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				StatusCode: http.StatusOK,
			}}

			if _, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 10,
				Sort:  "best",
			}); err == nil {
				t.Fatal("expected error")
			} else if !errors.As(err, &errDecode) {
				t.Fatalf("want: %v, got: %v", errDecode, err)
			}
		})
	})
}

func replay(t *testing.T, cassette string) {
	t.Helper()

	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       cassette,
		Mode:               recorder.ModeReplayOnly,
		RealTransport:      http.DefaultTransport,
		SkipRequestLatency: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = r.Stop()
	})

	http.DefaultClient.Transport = r
}

func TestClient_VCR(t *testing.T) {
	ctx := context.Background()

	c, err := freepublicapis.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("2024-11-24", func(t *testing.T) {
		replay(t, "../../internal/probe/vcr/2024-11-24")

		{
			res, err := c.GetRandom(ctx)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.GetAPI(ctx, 275)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}
	})

	t.Run("2024-11-28", func(t *testing.T) {
		replay(t, "../../internal/probe/vcr/2024-11-28")

		{
			res, err := c.GetRandom(ctx)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.GetAPI(ctx, 275)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 10,
				Sort:  "best",
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, nil)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}
	})

	t.Run("2024-11-29", func(t *testing.T) {
		replay(t, "../../internal/probe/vcr/2024-11-29")

		{
			res, err := c.GetRandom(ctx)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.GetAPI(ctx, 275)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, nil)
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 100,
				Sort:  "best",
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 101,
				Sort:  "best",
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 301,
				Sort:  "best",
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 1000,
				Sort:  "best",
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}

		{
			res, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
				Limit: 10,
				Sort:  "best",
			})
			if err != nil {
				t.Fatal(err)
			} else if res == nil {
				t.Fatal("result is nil")
			}
		}
	})
}
