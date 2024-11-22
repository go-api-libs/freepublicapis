// Package freepublicapis_test provides tests for the client package.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.
package freepublicapis_test

import "net/http"

type testRoundTripper struct {
	rsp *http.Response
	err error
}

func (t *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.rsp, t.err
}
