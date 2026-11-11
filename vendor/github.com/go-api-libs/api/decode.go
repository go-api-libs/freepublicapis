package api

import (
	"encoding/json/v2"
	"errors"
	"net/http"

	"github.com/MarkRosemaker/cassette"
)

// DecodeLenient is a debug-capable client's second try at a response body it could not decode as opts say, err
// being why. Decoding it the first time is the client's own, so that a client not in debug mode does nothing more
// than decode.
//
// In debug mode, ia is the call the client recorded: DecodeLenient adds it to the interactions in debugFile, for
// openapi-enrich to learn from, and decodes the body ia recorded into out again, with members the specification does
// not know allowed, so that the response is still of use while they wait to be added. Without debug mode, ia is nil,
// and it returns err as a decoding error.
func DecodeLenient[T any](rsp *http.Response, err error, ia *cassette.Interaction, debugFile string, out *T,
	opts json.Options,
) error {
	if ia == nil {
		return WrapDecodingError(rsp, err)
	}

	if kerr := cassette.AddInteraction(debugFile, ia); kerr != nil {
		return errors.Join(WrapDecodingError(rsp, err), kerr)
	}

	*out = *new(T)
	if err := json.Unmarshal(ia.Response.Body, out, opts, json.RejectUnknownMembers(false)); err != nil {
		return WrapDecodingError(rsp, err)
	}

	return nil
}
