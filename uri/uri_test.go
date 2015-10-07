package uri

import (
	"testing"
)

func TestURI(t *testing.T) {
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

		[]string{"foo.invalid", ""},
		[]string{"foo.bar.invalid", ""},

		// Cases found by https://github.com/dvyukov/go-fuzz
		[]string{"-.i-", ""},
		[]string{"0..in", ""},
		[]string{"[0.0", ""},
		[]string{"a b.com", ""},
		[]string{"-a.com", ""},

		// From: https://en.wikipedia.org/wiki/Hostname
		// "a subsequent specification (RFC 1123) permitted hostname labels to start with digits"
		[]string{"00l.com", "00l.com"},
		[]string{"000.com", "000.com"},

		[]string{"test.githubusercontent.com", "test.githubusercontent.com"},

		[]string{"http://测试.com/", "xn--0zwm56d.com"},
		[]string{"http://测试.com:80/", "xn--0zwm56d.com"},
	}

	for _, c := range cases {
		domain, err := EffectiveTLDPlusOne(c[0])
		if domain != c[1] {
			t.Errorf("%s: %s != \"%s\" (%v)\n", c[0], domain, c[1], err)
		}
	}
}

func BenchmarkURI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EffectiveTLDPlusOne("https://www.example.com/foobar?test")
	}
}
