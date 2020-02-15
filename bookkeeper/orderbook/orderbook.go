package orderbook

import (
	"container/heap"
	"strconv"

	"../../common/constants"
	"github.com/buger/jsonparser"
)

/*******************************************************************************/

//Price levels in BidHeap shall be in descending order
type BidHeap []*priceLevel

func (bidH BidHeap) Len() int { return len(bidH) }

func (bidH BidHeap) Less(i, j int) bool {
	return bidH[i].Price > bidH[j].Price
}

func (bidH BidHeap) Swap(i, j int) {
	bidH[i], bidH[j] = bidH[j], bidH[i]
	bidH[i].Index = i
	bidH[j].Index = j
}

func (bidH *BidHeap) Push(x interface{}) {
	n := len(*bidH)
	pricelevel := x.(*priceLevel)
	pricelevel.Index = n
	*bidH = append(*bidH, pricelevel)
}

func (bidH *BidHeap) Pop() interface{} {
	old := *bidH
	n := len(old)
	pricelevel := old[n-1]
	old[n-1] = nil
	pricelevel.Index = -1
	*bidH = old[0 : n-1]
	return pricelevel
}

func (bidH *BidHeap) Top() *priceLevel {
	return (*bidH)[0]
}

/*******************************************************************************/

//Price levels in AskHeap shall be in ascending order
type AskHeap []*priceLevel

func (askH AskHeap) Len() int { return len(askH) }

func (askH AskHeap) Less(i, j int) bool {
	return askH[i].Price < askH[j].Price
}

func (askH AskHeap) Swap(i, j int) {
	askH[i], askH[j] = askH[j], askH[i]
	askH[i].Index = i
	askH[j].Index = j
}

func (askH *AskHeap) Push(x interface{}) {
	n := len(*askH)
	pricelevel := x.(*priceLevel)
	pricelevel.Index = n
	*askH = append(*askH, pricelevel)
}

func (askH *AskHeap) Pop() interface{} {
	old := *askH
	n := len(old)
	pricelevel := old[n-1]
	old[n-1] = nil
	pricelevel.Index = -1
	*askH = old[0 : n-1]
	return pricelevel
}

func (askH *AskHeap) Top() *priceLevel {
	return (*askH)[0]
}

/*******************************************************************************/

type OrderBook struct {
	Exchange  int
	PricePair string
	Bids      *BidHeap
	Asks      *AskHeap
}

type priceLevel struct {
	Price  float64
	Volume float64
	Index  int
}

/*******************************************************************************/

//Testing Helpers
func popAskHeap(ah *AskHeap) priceLevel {
	return *(heap.Pop(ah).(*priceLevel))
}

func popBidHeap(bh *BidHeap) priceLevel {
	return *(heap.Pop(bh).(*priceLevel))
}

func (ob *OrderBook) getLowestAsk() priceLevel {
	return *(*(ob.Asks))[0]
}

func (ob *OrderBook) getHighestBid() priceLevel {
	return *(*(ob.Bids))[0]
}

/*******************************************************************************/

func NewOrderBook(exchange int, pricepair string, msg []byte) (*OrderBook, error) {

	if exchange == constants.CoinbasePro {
		var bh BidHeap
		var ah AskHeap

		//Initialize bid heap
		i := 0
		jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			bidPrice, err := jsonparser.GetString(value, "[0]")
			bidVol, err := jsonparser.GetString(value, "[1]")
			bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
			bidVolF, err := strconv.ParseFloat(bidVol, 64)

			pricelevel := &priceLevel{Price: bidPriceF, Volume: bidVolF, Index: i}
			i++
			bh = append(bh, pricelevel)
		}, "bids")

		heap.Init(&bh)

		//Initialize ask heap
		i = 0
		jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			askPrice, err := jsonparser.GetString(value, "[0]")
			askVol, err := jsonparser.GetString(value, "[1]")
			askPriceF, err := strconv.ParseFloat(askPrice, 64)
			askVolF, err := strconv.ParseFloat(askVol, 64)

			pricelevel := &priceLevel{Price: askPriceF, Volume: askVolF, Index: i}
			i++
			ah = append(ah, pricelevel)
		}, "asks")

		heap.Init(&ah)

		ob := &OrderBook{Exchange: constants.CoinbasePro, PricePair: pricepair, Bids: &bh, Asks: &ah}
		return ob, nil
	}

	return nil, nil
}

func (ob *OrderBook) locatePriceLevel() {

}

// BatchUpdate will perform N number of updates/deletes on a Heap. Once all actions are completed,
// the heap invariant will be restored.
func (ob *OrderBook) BatchUpdate() {

}
