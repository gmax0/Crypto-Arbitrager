package orderbook

import (
	"fmt"
	"strconv"
	"container/heap"

	"../../common/constants"
	"github.com/buger/jsonparser"
)

/*******************************************************************************/

type BidHeap []*PriceLevel

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
	pricelevel := x.(*PriceLevel)
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

/*******************************************************************************/

type AskHeap []*PriceLevel

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
	pricelevel := x.(*PriceLevel)
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

/*******************************************************************************/

type PriceVol struct {
	Price  float64
	Volume float64
}
type OrderBook struct {
	Exchange  int
	PricePair string
	Bids BidHeap
	Ask  AskHeap
}

type PriceLevel struct {
	Price float64
	Volume float64
	Index int
}

/*******************************************************************************/

func NewOrderBook(exchange int, pricepair string, msg []byte) (*OrderBook, error) {

	if exchange == constants.CoinbasePro {
		var buyPriceLevels BidHeap
		var askPriceLevels AskHeap

		//Initialize bid heap
		i := 0
		jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			bidPrice, err  := jsonparser.GetString(value, "[0]")
			bidVol, err    := jsonparser.GetString(value, "[1]")
			bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
			bidVolF, err   := strconv.ParseFloat(bidVol, 64)

			pricelevel := &PriceLevel{Price: bidPriceF, Volume: bidVolF, Index: i}
			i++
			buyPriceLevels = append(buyPriceLevels, pricelevel)
		}, "bids")

		heap.Init(&buyPriceLevels)
		for _, level := range buyPriceLevels {
			fmt.Println(level)
		}
		fmt.Println(len(buyPriceLevels))

		for buyPriceLevels.Len() > 0 {
			pricelevel := heap.Pop(&buyPriceLevels).(*PriceLevel)
			fmt.Println(pricelevel)
		}

		//Initialize ask heap
		i = 0
		jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			askPrice, err  := jsonparser.GetString(value, "[0]")
			askVol, err    := jsonparser.GetString(value, "[1]")
			askPriceF, err := strconv.ParseFloat(askPrice, 64)
			askVolF, err   := strconv.ParseFloat(askVol, 64)

			pricelevel := &PriceLevel{Price: askPriceF, Volume: askVolF, Index: i}
			i++
			askPriceLevels = append(askPriceLevels, pricelevel)
		}, "asks")

		heap.Init(&askPriceLevels)
		for _, level := range askPriceLevels {
			fmt.Println(level)
		}
		fmt.Println(len(askPriceLevels))

		for askPriceLevels.Len() > 0 {
			pricelevel := heap.Pop(&askPriceLevels).(*PriceLevel)
			fmt.Println(pricelevel)
		}

	}


	return nil, nil
}

/*
func NewOrderBook(exchange int, pricepair string, msg []byte) (*OrderBook, error) {
	ob := &OrderBook{Exchange: exchange, PricePair: pricepair}
	ob.Bids = make([]PriceVol, PriceDepth) 
	ob.Asks = make([]PriceVol, PriceDepth)

	if exchange == constants.CoinbasePro {
		for i := 0; i < PriceDepth; i++ {
			bidPrice, err := jsonparser.GetString(msg, "bids", fmt.Sprintf("[%d]", i), "[0]")
			bidVol, err := jsonparser.GetString(msg, "bids", fmt.Sprintf("[%d]", i), "[1]")
			askPrice, err := jsonparser.GetString(msg, "asks", fmt.Sprintf("[%d]", i), "[0]")
			askVol, err := jsonparser.GetString(msg, "asks", fmt.Sprintf("[%d]", i), "[1]")

			bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
			bidVolF, err := strconv.ParseFloat(bidVol, 64)
			askPriceF, err := strconv.ParseFloat(askPrice, 64)
			askVolF, err := strconv.ParseFloat(askVol, 64)
			
			if err != nil {
				return nil, err
			}

			ob.Bids[i].Price = bidPriceF
			ob.Bids[i].Volume = bidVolF

			ob.Asks[i].Price = askPriceF
			ob.Asks[i].Volume = askVolF
		}
	}

	return ob, nil
}
*/