// +build ignore

// Run using: go run web.go
// For a demo see: http://dubbelboer.com:8090
package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/atomx/publicsuffix/uri"
)

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
	d.TLDPlusOne, err = uri.EffectiveTLDPlusOne(d.Hostname)

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
