
Our local fork of https://code.google.com/p/go/source/browse/publicsuffix/?repo=net

Last update Jun: 2015-08-10

To check for a new version see: https://github.com/publicsuffix/list

To update run:
```bash
$ go run gen.go -version "3c153e3947dd6ee130c3f8706955281fa9923ed4" -url "https://raw.githubusercontent.com/publicsuffix/list/3c153e3947dd6ee130c3f8706955281fa9923ed4/public_suffix_list.dat"       > table.go
$ go run gen.go -version "3c153e3947dd6ee130c3f8706955281fa9923ed4" -url "https://raw.githubusercontent.com/publicsuffix/list/3c153e3947dd6ee130c3f8706955281fa9923ed4/public_suffix_list.dat" -test > table_test.go
```

See: [publicsuffix GoDoc](https://godoc.org/github.com/atomx/publicsuffix) and [uri GoDoc](https://godoc.org/github.com/atomx/publicsuffix/uri)

