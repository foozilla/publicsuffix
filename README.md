
Our local fork of https://code.google.com/p/go/source/browse/publicsuffix/?repo=net

Last update Jun: 2014-09-02 07:27 +0200

To check for a new version see: http://hg.mozilla.org/mozilla-central/filelog/29fbfc1b31aa/netwerk/dns/effective_tld_names.dat

To update run:
```bash
$ go run gen.go -version "f2c25ddbd1cf 2014-09-02 07:27 +0200"       >table.go
$ go run gen.go -version "f2c25ddbd1cf 2014-09-02 07:27 +0200" -test >table_test.go
```

See: [GoDoc](https://godoc.org/code.google.com/p/go.net/publicsuffix)

Example:
```go
var (
  hostEndRegexp = regexp.MustCompile("[^a-z0-9\\.\\-]")
)


// EffectiveTLDPlusOne returns the effective top level domain plus one more label.
// For example, "http://www.example.com/foobar" will be "example.com".
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

	// A TLD+1 needs to be at least 4 characters.
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

	// A TLD+1 needs to be at least 4 characters.
	if i < 3 {
		return "", fmt.Errorf("invalid domain")
	}

	// Some web clients really fuck up and somehow end a domain with a '.', remove it.
	if u[i] == '.' {
		u = u[0:i]
	}

	if strings.IndexByte(u, '.') < 1 {
		return "", fmt.Errorf("invalid domain")
	}

	// Check if it's an IP.
	if net.ParseIP(u) != nil {
		return u, nil
	}

  return publicsuffix.EffectiveTLDPlusOne(u)
}
```

