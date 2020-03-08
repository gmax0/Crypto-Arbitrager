package constants

const (
	CoinbasePro = 1
	Poloniex    = 2

	BidUpdate = true
	AskUpdate = false
)

//More canoniclization
var TickerNums map[string]int  //Ticker Symbol -> int
var TickerNames map[int]string //int -> Ticker Symbol

func init() {
	TickerNums = make(map[string]int)
	// TickerNums[]
	TickerNames = make(map[int]string)
	// TickerNames[]
}
