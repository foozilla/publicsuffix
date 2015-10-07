package uri

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/atomx/publicsuffix"
	"github.com/miekg/dns/idn"
)

// nonHostname returns the index of the first character not belonging to the hostname.
func nonHostname(s string) int {
	for i := 0; i < len(s); {
		r, l := utf8.DecodeRuneInString(s[i:])

		// https://en.wikipedia.org/wiki/Hostname#Restrictions_on_valid_host_names
		// "The Internet standards (Requests for Comments) for protocols mandate that component hostname labels may
		// contain only the ASCII letters 'a' through 'z' (in a case-insensitive manner), the digits '0' through '9', and the hyphen ('-')."
		// We skip 'A' - 'Z' because we only deal with lower case strings.
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.' && r != '-' {
			return i
		}

		i += l
	}
	return -1
}

// EffectiveTLDPlusOne returns the effective top level domain plus one more label.
// For example, "http://www.example.com/foobar" will be "example.com".
// It will only return valid ICANN domain names or IP addresses (both v4 and v6).
func EffectiveTLDPlusOne(u string) (string, error) {
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
		return "", fmt.Errorf("uri: invalid domain %q", u)
	}

	if u[0] == '.' {
		u = u[1:]
	}

	// IPv6?
	if u[0] == '[' {
		i := strings.Index(u, "]")
		if i == -1 {
			return "", fmt.Errorf("uri: invalid domain %q", u)
		}
		if net.ParseIP(u[1:i]) != nil {
			return u[1:i], nil
		}
	}

	// Trim everything after the first non hostname character.
	ii := nonHostname(u)

	if ii != -1 {
		u = u[0:ii]
	}

	// IE11 doesn't use Punycode in referrers, so encode it here first.
	// No need to check if this is needed, idn.ToPunycode has this check already.
	u = idn.ToPunycode(u)

	i := len(u) - 1

	// A TLD+1 needs to be at least 4 characters (g.cn for example).
	if i < 3 { // Note the - 1 above.
		return "", fmt.Errorf("uri: invalid domain %q", u)
	}

	if strings.Contains(u, "..") {
		return "", fmt.Errorf("uri: invalid domain %q", u)
	}

	// Some web clients really fuck up and somehow end a domain with a '.', remove it.
	if u[i] == '.' {
		u = u[0:i]

		// We removed 1 character so check the length again.
		if i < 4 {
			return "", fmt.Errorf("uri: invalid domain %q", u)
		}
	}

	if strings.IndexByte(u, '.') < 1 {
		return "", fmt.Errorf("uri: invalid domain %q", u)
	}

	// Check if it's an IP.
	if net.ParseIP(u) != nil {
		return u, nil
	}

	suffix, _, matched := publicsuffix.PublicSuffix(u)
	if !matched {
		 return "", fmt.Errorf("uri: no tld match found for domain %q", u)
	}

	if len(u) <= len(suffix) {
		return "", fmt.Errorf("uri: cannot derive eTLD+1 for domain %q", u)
	}
	i = len(u) - len(suffix) - 1
	if u[i] != '.' {
		return "", fmt.Errorf("uri: invalid public suffix %q for domain %q", suffix, u)
	}

	u = u[1+strings.LastIndex(u[:i], "."):]

	if c := u[0]; (c < 'a' || c > 'z') && (c < '0' || c > '9') {
		return "", fmt.Errorf("uri: invalid domain %q", u)
	}

	return u, nil
}
