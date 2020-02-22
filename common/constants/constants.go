package constants

import (
	"regexp"
)

const (
	CoinbasePro = 1
	Poloniex    = 2
)

var CoinbaseProReg *(regexp.Regexp)

func init() {
	CoinbaseProReg = regexp.MustCompile("[^a-zA-Z0-9.]+")
}
