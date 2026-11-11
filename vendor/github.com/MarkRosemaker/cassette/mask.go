package cassette

import (
	"bytes"
	"encoding/json/jsontext"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// masker configures which recorded values [Interactions.Mask] redacts.
//
// Redaction is shape-preserving: a masked e-mail address is replaced by another
// e-mail address, a masked UUID by another UUID, a masked number by a number.
// Schemas inferred from a masked recording therefore keep the types and formats
// they would have had from the original, so masking a cassette does not change
// the specification it produces — only the examples in it.
//
// The placeholders describe one fictional person: John Doe, johndoe,
// john.doe@example.com. A reader who meets one of them recognises the rest.
type masker struct {
	// HeaderKeys are header names whose values are redacted, in requests and
	// responses alike. Matching is case-insensitive.
	HeaderKeys []string

	// BodyKeys are JSON object keys whose values are redacted wherever they
	// occur in a request or response body. Matching is case-insensitive and
	// exact: "id" does not match "ownerId". When a matching key holds an object
	// or an array, the whole subtree below it is redacted.
	BodyKeys []string

	// IDKeys are keys holding identifiers, redacted only when the value looks
	// machine-generated: a UUID, or a run of at least 16 characters that has no
	// separators and contains a digit.
	//
	// Use this for keys far too common to redact outright. "id" holding an
	// account UUID or a numeric account ID is caught, while "id" holding a
	// structural name — Notion labels its title property "title" — is left alone.
	IDKeys []string

	// NameKeys are keys holding a person's name, replaced with "John Doe" when
	// the value actually reads as one: two to four capitalised words of letters.
	//
	// The shape check matters as much as the key. Recorded data is full of
	// capitalised phrases that are not people — company names, taxonomy labels,
	// page titles — and replacing those would corrupt the fixture while
	// protecting nobody.
	NameKeys []string

	// UsernameKeys are keys holding a handle, replaced with "johndoe" in the
	// same case as the original: "friedamuster" becomes "johndoe",
	// "FriedaMuster" becomes "JohnDoe", "@friedamuster" becomes "@johndoe".
	//
	// Casing is preserved because an API often stores the same handle several
	// ways — Habitica keeps both username and lowerCaseUsername — and a
	// replacement that flattened them would make the fixture describe something
	// the API never returns.
	UsernameKeys []string

	// Values are literal strings redacted wherever they occur in a body,
	// whatever key holds them. Use this for identifiers that live under keys too
	// generic to redact wholesale — an account UUID appearing under "id",
	// "ownerId" and "value" is caught by listing the UUID here, whereas listing
	// those keys in BodyKeys would destroy unrelated data.
	//
	// Prefer the key-based fields where they suffice: a value listed here has to
	// be written down somewhere, and a masker configured in tracked source would
	// publish the very string it is meant to hide.
	Values []string

	// KeepEmails leaves e-mail addresses in place. By default any string that is
	// an e-mail address is redacted wherever it occurs, whatever key holds it,
	// because an address is personal data no matter where it was recorded and
	// the key holding it is often too generic to list.
	KeepEmails bool
}

func (m masker) rules() rules {
	return rules{
		headers:   lowerSet(m.HeaderKeys),
		keys:      lowerSet(m.BodyKeys),
		ids:       lowerSet(m.IDKeys),
		names:     lowerSet(m.NameKeys),
		usernames: lowerSet(m.UsernameKeys),
		values:    m.Values,
		emails:    !m.KeepEmails,
	}
}

// defaultMasker returns the configuration used by [Interactions.Mask]: the
// headers and body keys that carry credentials in most APIs.
//
// It knows nothing about which identifiers are yours. A recording that carries
// an account ID under a generic key such as "id" needs that key listed in
// IDKeys, or the value itself in Values.
//
// X-Client is not masked: an API such as Habitica uses it to name the tool
// calling it, and a replayed call must send it as recorded.
func defaultMasker() masker {
	return masker{
		HeaderKeys: []string{
			"Authorization", "Proxy-Authorization",
			"Cookie", "Set-Cookie",
			"X-Api-Key", "X-Api-User", "X-Api-Token",
			"X-Auth-Token", "X-Csrf-Token",
			"X-Rd-Token", "X-Session-Id",
		},
		BodyKeys: []string{
			"access_token", "accessToken",
			"api_key", "apiKey",
			"client_secret", "clientSecret",
			"credentials", "hardware_info",
			"mac", "password", "passwd",
			"private_key", "privateKey",
			"refresh_token", "refreshToken",
			"secret", "session", "sessionId", "session_id",
			"signature", "token", "userId", "user_id", "uuid",
		},
	}
}

// rules is the compiled form of a [masker]: key sets rather than slices.
type rules struct {
	headers, keys, ids, names, usernames map[string]bool
	values                               []string
	emails                               bool
}

// scope records which kinds of redaction an enclosing key has switched on. It
// travels down the tree, so a matching key redacts everything below it.
type scope struct {
	all, id, name, username bool
}

// maskWith redacts sensitive values in place according to m.
func (ia *Interaction) maskWith(m masker) {
	ia.mask(m.rules())
}

// Mask redacts credentials in place: in headers, in bodies and in query values.
func (ias Interactions) Mask() { ias.maskWith(defaultMasker()) }

// maskWith redacts sensitive values in place according to m.
func (ias Interactions) maskWith(m masker) {
	r := m.rules()
	for _, ia := range ias {
		ia.mask(r)
	}
}

func (ia *Interaction) mask(r rules) {
	maskHeaders(ia.Request.Headers, r.headers)
	maskHeaders(ia.Response.Headers, r.headers)

	ia.Request.URL.RawQuery = maskQuery(ia.Request.URL.RawQuery, r.keys)

	ia.Request.Body = maskBody(ia.Request.Body, r)
	ia.Response.Body = maskBody(ia.Response.Body, r)
}

func lowerSet(keys []string) map[string]bool {
	set := make(map[string]bool, len(keys))
	for _, k := range keys {
		set[strings.ToLower(k)] = true
	}

	return set
}

func maskHeaders(h map[string][]string, keys map[string]bool) {
	for k, vals := range h {
		if !keys[strings.ToLower(k)] {
			continue
		}

		masked := make([]string, len(vals))
		for i, v := range vals {
			masked[i] = maskString(v)
		}

		h[k] = masked
	}
}

// maskQuery masks the values of the query parameters raw holds whose names are among keys, as it does the members of a
// body; every other parameter is kept as it was written.
func maskQuery(raw string, keys map[string]bool) string {
	if raw == "" {
		return raw
	}

	params := strings.Split(raw, "&")
	for i, param := range params {
		name, value, ok := strings.Cut(param, "=")
		if !ok || !keys[strings.ToLower(queryName(param))] {
			continue
		}

		if unescaped, err := url.QueryUnescape(value); err == nil {
			value = unescaped
		}

		// the asterisks of a masked value need no escaping, and read better without
		params[i] = name + "=" + strings.ReplaceAll(url.QueryEscape(maskString(value)), "%2A", "*")
	}

	return strings.Join(params, "&")
}

// maskBody rewrites b according to r. A body that is not valid JSON is returned
// unchanged: masking must never corrupt a recording.
func maskBody(b Body, r rules) Body {
	if len(bytes.TrimSpace(b)) == 0 {
		return b
	}

	buf := &bytes.Buffer{}
	dec := jsontext.NewDecoder(bytes.NewReader(b))
	enc := jsontext.NewEncoder(buf)

	if err := maskValue(dec, enc, r, scope{}); err != nil {
		return b
	}

	// The encoder terminates the value with a newline; the original body has no
	// such terminator, and adding one would show up as a diff in every recording.
	return Body(bytes.TrimRight(buf.Bytes(), "\n"))
}

// maskValue copies one JSON value from dec to enc, redacting according to r and
// whatever sc an enclosing key has already switched on.
func maskValue(dec *jsontext.Decoder, enc *jsontext.Encoder, r rules, sc scope) error {
	switch dec.PeekKind() {
	case '{':
		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		if err := enc.WriteToken(jsontext.BeginObject); err != nil {
			return err
		}

		for dec.PeekKind() != '}' {
			key, err := dec.ReadToken()
			if err != nil {
				return err
			}

			name := key.String()
			if err := enc.WriteToken(jsontext.String(name)); err != nil {
				return err
			}

			k := strings.ToLower(name)
			if err := maskValue(dec, enc, r, scope{
				all:      sc.all || r.keys[k],
				id:       sc.id || r.ids[k],
				name:     sc.name || r.names[k],
				username: sc.username || r.usernames[k],
			}); err != nil {
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
			if err := maskValue(dec, enc, r, sc); err != nil {
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

		return enc.WriteToken(maskToken(tok, r, sc))
	}
}

func maskToken(tok jsontext.Token, r rules, sc scope) jsontext.Token {
	switch tok.Kind() {
	case jsontext.KindString:
		s := tok.String()

		switch {
		case sc.all:
			return jsontext.String(maskString(s))

		// An identifier key only redacts values that look machine-generated, so
		// that a structural name sitting under the same key survives.
		case sc.id && looksLikeID(s):
			return jsontext.String(maskString(s))

		case sc.name && looksLikePersonName(s):
			return jsontext.String(maskedName)

		case sc.username && s != "":
			return jsontext.String(maskUsername(s))

		// An e-mail address is personal data wherever it turns up, including
		// under a key too generic to list in BodyKeys.
		case r.emails && reEmail.MatchString(s):
			return jsontext.String(maskedEmail)
		}

		return jsontext.String(replaceLiterals(s, r.values))
	case jsontext.KindNumber:
		if sc.all {
			// A masked number stays a number, so the inferred type is unchanged.
			return jsontext.Int(0)
		}
	default:
	}

	return tok.Clone()
}

// rePersonName matches two to four capitalised words of letters, allowing the
// apostrophes and hyphens that occur in real names.
var rePersonName = regexp.MustCompile(
	`^[A-Z][a-z]+(?:['\-][A-Z]?[a-z]+)*(?: [A-Z][a-z]+(?:['\-][A-Z]?[a-z]+)*){1,3}$`,
)

// looksLikePersonName reports whether s reads as a person's name rather than a
// capitalised phrase such as a company or a label. It is intentionally strict:
// a word like "APIs", with capitals past the first letter, disqualifies the
// whole string.
func looksLikePersonName(s string) bool { return rePersonName.MatchString(s) }

// reOpaqueID matches a run with no separators or spaces, the usual shape of a
// generated identifier.
var reOpaqueID = regexp.MustCompile(`^[0-9A-Za-z_-]{16,}$`)

// looksLikeID reports whether s has the shape of a machine-generated
// identifier rather than a human-chosen name.
func looksLikeID(s string) bool {
	if reUUID.MatchString(s) || reHexUUID.MatchString(s) {
		return true
	}

	// Require a digit: an all-letter run of this length reads as prose, not as
	// an identifier worth redacting.
	return reOpaqueID.MatchString(s) && strings.ContainsAny(s, "0123456789")
}

// maskUsername replaces a handle with "johndoe", matching the case of the
// original and keeping a leading "@".
func maskUsername(s string) string {
	at := ""
	if rest, ok := strings.CutPrefix(s, "@"); ok {
		at, s = "@", rest
	}

	if s == "" {
		return at
	}

	hasLower := strings.ContainsFunc(s, unicode.IsLower)

	switch {
	case !hasLower && strings.ContainsFunc(s, unicode.IsUpper):
		return at + maskedUsernameUpper
	case unicode.IsUpper(rune(s[0])):
		return at + maskedUsernameTitle
	default:
		return at + maskedUsername
	}
}

// replaceLiterals masks any of values occurring in s, whatever key holds it.
func replaceLiterals(s string, values []string) string {
	for _, v := range values {
		if v != "" && strings.Contains(s, v) {
			s = strings.ReplaceAll(s, v, maskString(v))
		}
	}

	return s
}

var (
	reEmail = regexp.MustCompile(`^[\w.+-]+@[\w-]+\.[\w.]{2,}$`)
	reUUID  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// A UUID is also valid without its separators, and is recognised as one when
	// a format is inferred, so it has to keep that shape when masked.
	reHexUUID = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
	reMAC     = regexp.MustCompile(`^(?:[0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
)

// The placeholders, all describing the same fictional person.
const (
	maskedEmail         = "john.doe@example.com"
	maskedName          = "John Doe"
	maskedUsername      = "johndoe"
	maskedUsernameTitle = "JohnDoe"
	maskedUsernameUpper = "JOHNDOE"
	maskedUUID          = "00000000-0000-0000-0000-000000000000"
	maskedHexUUID       = "00000000000000000000000000000000"
	maskedMAC           = "00:00:00:00:00:00"
)

// RedactsBodyKey reports whether [Interactions.Mask] redacts values held by key.
//
// A masked string announces itself — the zero UUID and the asterisk run are
// recognisable by [IsMasked] — but a masked number does not: a redacted count
// becomes 0, which is indistinguishable from a genuine 0. The key is the only
// remaining signal, so callers documenting a masked recording can ask about it
// directly rather than guessing from the value.
func RedactsBodyKey(key string) bool {
	return lowerSet(defaultMasker().BodyKeys)[strings.ToLower(key)]
}

// IsMasked reports whether s is a placeholder produced by masking.
//
// Such a value carries no information a reader does not already have: the
// inferred format says "uuid" far more usefully than an example of
// "00000000-0000-0000-0000-000000000000" does. Callers building documentation
// from a masked recording can use this to leave the example out entirely.
func IsMasked(s string) bool {
	switch strings.TrimPrefix(s, "@") {
	case maskedEmail, maskedName, maskedUUID, maskedHexUUID, maskedMAC,
		maskedUsername, maskedUsernameTitle, maskedUsernameUpper:
		return true
	}

	// A run of asterisks, optionally still carrying its scheme.
	s = strings.TrimPrefix(s, "Bearer ")
	s = strings.TrimPrefix(s, "Basic ")

	return s != "" && strings.Trim(s, "*") == ""
}

// maskString replaces s with a value of the same shape, so that a format
// inferred from the masked recording matches the one inferred from the original.
func maskString(s string) string {
	const (
		bearer = "Bearer "
		basic  = "Basic "
	)

	switch {
	case s == "":
		return s
	case reEmail.MatchString(s):
		return maskedEmail
	case reUUID.MatchString(s):
		return maskedUUID
	case reHexUUID.MatchString(s):
		return maskedHexUUID
	case reMAC.MatchString(s):
		return maskedMAC
	case strings.HasPrefix(s, bearer):
		return bearer + strings.Repeat("*", len(s)-len(bearer))
	case strings.HasPrefix(s, basic):
		return basic + strings.Repeat("*", len(s)-len(basic))
	default:
		return strings.Repeat("*", len(s))
	}
}
