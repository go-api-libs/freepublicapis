package cassette

import (
	"bytes"
	"encoding/json/jsontext"
	"mime"
	"strings"
	"unicode/utf8"
)

// MaxStringLen is the most bytes a JSON string in a recorded body keeps: a longer one, such as a whole image in
// base64, is cut to it and ends in "…". It is long enough for every URL seen so far, presigned ones included.
const MaxStringLen = 2048

// IsText reports whether a body of the media type contentType is text, and so worth recording: JSON, XML, a form or
// text/*. Any other, such as a zip, a PDF, an image or a video, is recorded only as omitted. A body without a media
// type is treated as text.
func IsText(contentType string) bool {
	if contentType == "" {
		return true
	}

	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return true
	}

	switch {
	case strings.HasPrefix(mt, "text/"),
		mt == "application/json", strings.HasSuffix(mt, "+json"), strings.Contains(mt, "/json"),
		mt == "application/xml", strings.HasSuffix(mt, "+xml"),
		mt == "application/x-www-form-urlencoded", mt == "application/javascript":
		return true
	default:
		return false
	}
}

// TrimBodies omits every recorded body that is not text, see [IsText], and cuts every JSON string in the others
// longer than maxLen bytes, see [MaxStringLen].
func (ias Interactions) TrimBodies(maxLen int) {
	for i := range ias {
		ias[i].trimBodies(maxLen)
	}
}

func (ia *Interaction) trimBodies(maxLen int) {
	ia.Request.Body, ia.Request.BodyOmitted = trimBody(ia.Request.Body, ia.Request.BodyOmitted, ia.Request.Headers.Get("Content-Type"), maxLen)
	ia.Response.Body, ia.Response.BodyOmitted = trimBody(ia.Response.Body, ia.Response.BodyOmitted, ia.Response.Headers.Get("Content-Type"), maxLen)
}

// trimBody is b without what [Interactions.TrimBodies] leaves out, and whether b is omitted.
func trimBody(b Body, omitted bool, contentType string, maxLen int) (Body, bool) {
	if len(b) == 0 {
		return b, omitted
	}

	if !IsText(contentType) {
		return nil, true
	}

	return shortenStrings(b, maxLen), omitted
}

// shortenStrings cuts every string in the JSON value b longer than maxLen bytes, its member names aside. A body that
// is not valid JSON is returned as it is.
func shortenStrings(b Body, maxLen int) Body {
	if maxLen <= 0 || len(b) <= maxLen || !jsontext.Value(b).IsValid() {
		return b
	}

	buf := &bytes.Buffer{}
	dec := jsontext.NewDecoder(bytes.NewReader(b))
	enc := jsontext.NewEncoder(buf)

	if err := shortenValue(dec, enc, maxLen); err != nil {
		return b
	}

	// the encoder ends the value with a newline that the body did not have
	return Body(bytes.TrimRight(buf.Bytes(), "\n"))
}

func shortenValue(dec *jsontext.Decoder, enc *jsontext.Encoder, maxLen int) error {
	switch dec.PeekKind() {
	case '{':
		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		if err := enc.WriteToken(jsontext.BeginObject); err != nil {
			return err
		}

		for dec.PeekKind() != '}' {
			name, err := dec.ReadToken()
			if err != nil {
				return err
			}

			if err := enc.WriteToken(name); err != nil {
				return err
			}

			if err := shortenValue(dec, enc, maxLen); err != nil {
				return err
			}
		}

		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		return enc.WriteToken(jsontext.EndObject)
	case '[':
		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		if err := enc.WriteToken(jsontext.BeginArray); err != nil {
			return err
		}

		for dec.PeekKind() != ']' {
			if err := shortenValue(dec, enc, maxLen); err != nil {
				return err
			}
		}

		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		return enc.WriteToken(jsontext.EndArray)
	default:
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}

		if s := tok.String(); tok.Kind() == jsontext.KindString && len(s) > maxLen {
			cut := maxLen
			for cut > 0 && !utf8.RuneStart(s[cut]) {
				cut--
			}

			tok = jsontext.String(s[:cut] + "…")
		}

		return enc.WriteToken(tok)
	}
}
