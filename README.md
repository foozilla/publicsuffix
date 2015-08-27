
Our local fork of https://code.google.com/p/go/source/browse/publicsuffix/?repo=net

Last update Jun: 2015-08-27

To check for a new version see: https://github.com/publicsuffix/list

To update run:
```bash
go run gen.go -version "master" -url "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat"       > table.go
go run gen.go -version "master" -url "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat" -test > table_test.go
```

See: [publicsuffix GoDoc](https://godoc.org/github.com/atomx/publicsuffix) and [uri GoDoc](https://godoc.org/github.com/atomx/publicsuffix/uri)

