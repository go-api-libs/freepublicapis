package cassette

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions(
	jsontext.Multiline(true),
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(jsonutil.HTTPHeaderMarshal),
		json.MarshalToFunc(jsonutil.URLMarshal),
	)),
	json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(jsonutil.HTTPHeaderUnmarshal),
		json.UnmarshalFromFunc(urlUnmarshal),
	)),
)

// urlUnmarshal reads a URL as [jsonutil.URLUnmarshal] does, its query parameters sorted by name, so that recordings
// made before requests were normalized match those made after.
func urlUnmarshal(dec *jsontext.Decoder, u *url.URL) error {
	if err := jsonutil.URLUnmarshal(dec, u); err != nil {
		return err
	}

	*u = normalizeURL(*u)

	return nil
}

// InteractionsReadFile reads the interactions in the JSON file at path.
func InteractionsReadFile(path string) (Interactions, error) {
	return jsonutil.ReadFile[Interactions](path, jsonOpts)
}

// WriteFile writes the interactions as JSON to the file at path, masked as [Interactions.Mask] does; ias are left as
// they are.
func (ias Interactions) WriteFile(path string) error {
	masked := ias.clone()
	masked.Mask()

	return jsonutil.WriteFile(path, masked, jsonOpts)
}

// MarshalWrite writes the interactions as JSON to w, as they are: unlike [Interactions.WriteFile], it does not mask
// them.
func (ias Interactions) MarshalWrite(w io.Writer) error {
	return json.MarshalWrite(w, ias, jsonOpts)
}

var mu sync.Mutex

// AddInteraction adds ia to the interactions in the file at path, masked as [Interactions.Mask] does and with its bodies trimmed as
// [Interactions.TrimBodies] does with [MaxStringLen], writing the file right away; ia is left as it is. A call the
// file holds already, the same request answered with the same status and body, is not added again: a client retrying
// a call, or polling one, records it once.
func AddInteraction(path string, ia *Interaction) error {
	ia = ia.clone()
	ia.trimBodies(MaxStringLen)
	ia.maskWith(defaultMasker())

	mu.Lock()
	defer mu.Unlock()

	ias, err := InteractionsReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("making dir for interactions file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("reading interactions file: %w", err)
	}

	if slices.ContainsFunc(ias, ia.sameCall) {
		return nil
	}

	if err := append(ias, ia).WriteFile(path); err != nil {
		return fmt.Errorf("writing interactions file: %w", err)
	}

	return nil
}

// sameCall reports whether ia and other are the same call: the same request answered with the same status and body.
// Headers are left out, since a response's, such as its date, differ from call to call.
func (ia *Interaction) sameCall(other *Interaction) bool {
	return ia.Request.Method == other.Request.Method &&
		ia.Request.URL.String() == other.Request.URL.String() &&
		bytes.Equal(canonical(ia.Request.Body), canonical(other.Request.Body)) &&
		ia.Response.StatusCode == other.Response.StatusCode &&
		bytes.Equal(canonical(ia.Response.Body), canonical(other.Response.Body))
}

// canonical is b in canonical JSON: b itself if it is JSON, else the JSON string it is written as (see [Body]), so
// that a body compares alike before it is written and after it is read back.
func canonical(b Body) []byte {
	v := jsontext.Value(bytes.Clone(b))
	if !v.IsValid() {
		v, _ = jsontext.AppendQuote(nil, b)
	}

	if err := v.Canonicalize(); err != nil {
		return b
	}

	return v
}

// clone returns a deep copy of ias.
func (ias Interactions) clone() Interactions {
	out := make(Interactions, len(ias))
	for i, ia := range ias {
		out[i] = ia.clone()
	}

	return out
}

// clone returns a deep copy of ia.
func (ia *Interaction) clone() *Interaction {
	c := *ia
	c.Request.Headers = ia.Request.Headers.Clone()
	c.Request.Body = bytes.Clone(ia.Request.Body)
	c.Response.Headers = ia.Response.Headers.Clone()
	c.Response.Body = bytes.Clone(ia.Response.Body)

	return &c
}
