// Package origin implements tools and methods to compare and perform simple pattern-matching
// on the [origin] header of a request on the server-side, specifically in the context of [cross-origin
// resource sharing] (CORS).
//
// [origin]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin
// [cross-origin resource sharing]: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
package origin

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// wildcard symbols.
const (
	wildcard = "*"
	anyValue = "*://*:*"
)

// Standard ports for common web protocols.
var knownPorts = map[string]string{
	"https":  "443",
	"wss":    "443",
	"http":   "80",
	"ws":     "80",
	"ftp":    "23",
	"gopher": "70",
}

// Split is similar to net.SplitHostPort, but accounts for the
// scheme (protocol) as well.
func Split(origin string) (scheme, host, port string, err error) {
	var u *url.URL

	u, err = url.Parse(origin)
	if err != nil {
		err = fmt.Errorf("invalid origin: %v", err)
		return
	}

	scheme, host, port = u.Scheme, u.Hostname(), u.Port()
	if scheme == "" {
		err = errors.New("invalid origin: missing scheme")
		return
	}

	if port == "" {
		var ok bool
		port, ok = knownPorts[scheme]
		if !ok {
			err = errors.New("invalid origin: missing port")
			return
		}
	}

	return
}

// splitPattern is similar to Split, but supports wildcard characters
// in scheme, hostname and port.
func splitPattern(pattern string) (scheme, host, port string, err error) {
	if pattern == wildcard {
		scheme, host, port = wildcard, wildcard, wildcard
		return
	}

	const sep = "://"

	parts := strings.SplitN(pattern, sep, 2)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid pattern: missing scheme")
	}

	scheme, host = parts[0], parts[1]

	if strings.Contains(host, ":") {
		host, port, err = net.SplitHostPort(host)
		if err != nil {
			err = fmt.Errorf("invalid pattern: %v", err)
		}
		return
	}

	var ok bool
	port, ok = knownPorts[scheme]
	if !ok {
		err = errors.New("invalid origin: missing port")
		return
	}

	return
}

// normalize readies a string for comparison.
func normalize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return s
}

// matchHostname matches a hostname against pattern.
func matchHostname(origin, pattern string) (bool, error) {
	a := strings.Split(normalize(pattern), ".")
	b := strings.Split(normalize(origin), ".")

	for i := len(a) - 1; i >= 0; i-- {
		if len(b) < i+1 {
			return false, nil
		}
		if a[i] == wildcard || b[i] == wildcard {
			continue
		}
		if a[i] != b[i] {
			return false, nil
		}
	}
	return true, nil
}

func matchString(origin, pattern string) (bool, error) {
	if origin == "" {
		return false, nil
	}

	if pattern == wildcard {
		return true, nil
	}

	return strings.EqualFold(origin, pattern), nil
}

// Match returns true if the scheme, hostname and port
// of origin match the ones in the given pattern.
//
// Both the origin and the pattern must be formatted as:
// 	scheme://hostname:port
//
// Pattern may contain a wildcard "*" in any of the three
// components. For example, "https://*.example.com:*" will consider
// any subdomain of "example.com" on any port number as a match,
// provided that the scheme is HTTPS.
//
// The port number may be omitted in either the origin or pattern
// when the scheme has a known standard port number. For example,
// "https://example.com" and "https://example.com:443" are a match.
//
// The special pattern value "*" is equivalent to "*://*:*", and
// matches with any non-empty and valid origin.
func Match(origin, pattern string) (bool, error) {
	os, oh, op, err := Split(origin)
	if err != nil {
		return false, err
	}

	if pattern == "" {
		return false, errors.New("pattern cannot be an empty string")
	}
	if pattern == wildcard || pattern == anyValue {
		return true, nil
	}

	ps, ph, pp, err := splitPattern(pattern)
	if err != nil {
		return false, err
	}

	if ok, err := matchString(os, ps); !ok || err != nil {
		return false, err
	}

	if ok, err := matchHostname(oh, ph); !ok || err != nil {
		return false, err
	}

	if ok, err := matchString(op, pp); !ok || err != nil {
		return false, err
	}

	return true, nil
}

// Patterns represents a list of trusted origins or patterns
// such as "https://example.com", or "*://*.example.com:*".
type Patterns []string

// Match returns true if any of the patterns in p matches
// with origin.
func (p Patterns) Match(origin string) (bool, error) {
	if origin == "" {
		return false, nil
	}

	for _, item := range p {
		ok, err := Match(origin, item)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// Get returns the value of the origin header in r.
//
// Am empty string is returned if the value in the header in "null",
// indicating an [opaque origin].
//
// [opaque origin]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin#directives
func Get(r *http.Request) string {
	str := r.Header.Get("Origin")
	if strings.EqualFold(str, "null") {
		str = ""
	}
	return str
}
