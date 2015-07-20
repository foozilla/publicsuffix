// +build ignore

// Run using: go run web.go
// For a demo see: http://dubbelboer.com:8090
package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	. "."
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

type data struct {
	Hostname   string
	TLDPlusOne string
	Error      string
}

var content = template.Must(template.New("content").Parse(`<!doctype html>
<html>
<head>
<meta charset=utf-8>
<title>publicsuffix</title>
<style>
input {
width: 100%;
}
</style>
</head>
<body>
<form action="" method=get>
<p>
<input type=text name=hostname value="{{.Hostname}}">
</p>
<p>
<input type=submit value=Lookup>
</p>
</form>
<p>
{{if .Error}}
{{.Error}}
{{else}}
{{.TLDPlusOne}}
{{end}}
</p>
`))

func index(w http.ResponseWriter, r *http.Request) {
	d := data{
		Hostname: r.FormValue("hostname"),
	}

	log.Printf("%22s | %s\n", r.RemoteAddr, d.Hostname)

	var err error
	d.TLDPlusOne, err = Example_EffectiveTLDPlusOne(d.Hostname)

	if err != nil {
		d.Error = err.Error()
	}

	w.Header().Set("Content-Type", "text/html")

	if err := content.Execute(w, d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", index)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Printf("[ERR] %v", err)
	}
}
