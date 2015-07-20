package publicsuffix

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

var (
	hostEndRegexp = regexp.MustCompile("[^a-z0-9\\.\\-]")
)

// EffectiveTLDPlusOne returns the effective top level domain plus one more label.
// For example, "http://www.example.com/foobar" will be "example.com".
func Example_EffectiveTLDPlusOne(u string) (string, error) {
	u = strings.ToLower(u)

	if strings.HasPrefix(u, "http%3a%2f%2f") ||
		strings.HasPrefix(u, "https%3a%2f%2f") ||
		strings.HasPrefix(u, "%2f%2f") {
		var err error
		u, err = url.QueryUnescape(u)
		if err != nil {
			return "", err
		}

		u = strings.ToLower(u)
	}

	// Trim http:// https:// or // from the start.
	if strings.HasPrefix(u, "http://") {
		u = u[len("http://"):]
	} else if strings.HasPrefix(u, "https://") {
		u = u[len("https://"):]
	} else if strings.HasPrefix(u, "//") {
		u = u[len("//"):]
	}

	// A TLD+1 needs to be at least 4 characters (g.cn for example).
	if len(u) < 4 {
		return "", fmt.Errorf("invalid domain")
	}

	if u[0] == '.' {
		u = u[1:]
	}

	// IPv6?
	if u[0] == '[' {
		i := strings.Index(u, "]")
		if net.ParseIP(u[1:i]) != nil {
			return u[1:i], nil
		}
	}

	// Trim everything after the first non hostname character.
	ii := hostEndRegexp.FindStringIndex(u)

	if len(ii) > 0 {
		u = u[0:ii[0]]
	}

	i := len(u) - 1

	// A TLD+1 needs to be at least 4 characters (g.cn for example).
	if i < 3 { // Note the - 1 above.
		return "", fmt.Errorf("invalid domain")
	}

	// Some web clients really fuck up and somehow end a domain with a '.', remove it.
	if u[i] == '.' {
		u = u[0:i]

		// We removed 1 character so check the length again.
		if i < 4 {
			return "", fmt.Errorf("invalid domain")
		}
	}

	if strings.IndexByte(u, '.') < 1 {
		return "", fmt.Errorf("invalid domain")
	}

	// Check if it's an IP.
	if net.ParseIP(u) != nil {
		return u, nil
	}

	return EffectiveTLDPlusOne(u)
}

func TestExample(t *testing.T) {
	cases := [][]string{
		[]string{"example.com", "example.com"},
		[]string{"example.com/", "example.com"},
		[]string{"www.example.com", "example.com"},
		[]string{"www.example.com/", "example.com"},
		[]string{"www.example.com/foobar", "example.com"},
		[]string{"//www.example.com/foobar", "example.com"},
		[]string{"http://www.example.com/foobar", "example.com"},
		[]string{"https://www.example.com/foobar", "example.com"},

		[]string{"example.com:80", "example.com"},
		[]string{"example.com:80/", "example.com"},
		[]string{"www.example.com:80", "example.com"},
		[]string{"www.example.com:80/", "example.com"},
		[]string{"www.example.com:80/foobar", "example.com"},
		[]string{"//www.example.com:80/foobar", "example.com"},
		[]string{"http://www.example.com:80/foobar", "example.com"},
		[]string{"https://www.example.com:80/foobar", "example.com"},

		[]string{"https://www.example.com/foobar?test", "example.com"},
		[]string{"https://www.example.com/foobar/test", "example.com"},
		[]string{"https://www.example.com/foobar/test?foo", "example.com"},
		[]string{"https://www.example.com/foobar/test?foo&bar", "example.com"},
		[]string{"https://www.example.com/foobar/test?foo=bar", "example.com"},
		[]string{"https://www.example.com/foobar/test?foo=bar&bar=foo", "example.com"},

		[]string{"https://www.example.com/foobar?test=foo:bar", "example.com"},

		[]string{"likes%26fb_source%3Dother_multiline%26action", ""},

		[]string{"example.com%2Ftest", "example.com"},

		[]string{".org", ""},
		[]string{"org", ""},

		[]string{"127.0.0.1", "127.0.0.1"},
		[]string{"http://127.0.0.1", "127.0.0.1"},
		[]string{"[2001:4860:0:2001::68]", "2001:4860:0:2001::68"},
		[]string{"http://[2001:4860:0:2001::68]", "2001:4860:0:2001::68"},

		[]string{"http://www.example.com./", "example.com"},

		[]string{"http%3A%2F%2Fexample.com", "example.com"},

		[]string{"http://example.com%2F", "example.com"},

		[]string{"[example.com/]", ""},
	}

	for _, c := range cases {
		domain, err := Example_EffectiveTLDPlusOne(c[0])
		if domain != c[1] {
			t.Errorf("%s: %s != %s (%v)\n", c[0], domain, c[1], err)
		}
	}
}
