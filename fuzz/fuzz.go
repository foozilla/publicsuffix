package fuzz

import (
	"fmt"

	"github.com/atomx/publicsuffix/uri"
)

func Fuzz(data []byte) int {
	if len(data) > 64 {
		return -1
	}

	d, err := uri.EffectiveTLDPlusOne(string(data))
	if err == nil {
		println(fmt.Sprintf("%q\t%q", string(data), d))
	}

	return 1
}
