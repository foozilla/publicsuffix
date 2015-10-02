
Our local fork of https://code.google.com/p/go/source/browse/publicsuffix/?repo=net

Last update: 2015-10-01

To check for a new version see: https://github.com/publicsuffix/list

To update run:
```bash
go run gen.go -v -version "master"       > table.go
go run gen.go -v -version "master" -test > table_test.go
```

See: [publicsuffix GoDoc](https://godoc.org/github.com/atomx/publicsuffix) and [uri GoDoc](https://godoc.org/github.com/atomx/publicsuffix/uri)

