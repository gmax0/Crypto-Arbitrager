package constants

const (
	CoinbasePro = 1
	Poloniex    = 2

	BidUpdate = true
	AskUpdate = false

	BTC_ETH = 1
	BTC_USD = 2
	ETH_USD = 3
)

//More canoniclization
var CBPricePairToInt map[string]int //Ticker Symbol -> int

func init() {
	CBPricePairToInt = make(map[string]int)
	CBPricePairToInt["ETH-BTC"] = BTC_ETH
	CBPricePairToInt["BTC-USD"] = BTC_USD
	CBPricePairToInt["ETH-USD"] = ETH_USD
}
