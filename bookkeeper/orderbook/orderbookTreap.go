package orderbook

import (
	"math/rand"
	"strconv"
	"time"

	"../../common/constants"
	"github.com/buger/jsonparser"
	"github.com/steveyen/gtreap" //Note that this is an immutable treap
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func priceLevelAscendingCompare(a, b interface{}) int {
	if (*a.(*priceLevel)).Price < (*b.(*priceLevel)).Price {
		return -1
	} else if (*a.(*priceLevel)).Price > (*b.(*priceLevel)).Price {
		return 1
	} else {
		return 0
	}
}

func priceLevelDescendingCompare(a, b interface{}) int {
	if (*a.(*priceLevel)).Price < (*b.(*priceLevel)).Price {
		return 1
	} else if (*a.(*priceLevel)).Price > (*b.(*priceLevel)).Price {
		return -1
	} else {
		return 0
	}
}

type OrderBookTreap struct {
	Exchange  int
	PricePair string
	bids      *gtreap.Treap //Treap of *priceLevel
	asks      *gtreap.Treap //Treap of *priceLevel
}

// Returns a pointer to the Ask Price Level if it exists in the Asks Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) getAskPriceLevel(price float64) *priceLevel {
	pl := ob.asks.Get(&priceLevel{Price: price})
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) getMinAskPriceLevel() *priceLevel {
	pl := ob.asks.Min()
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) deleteAskPriceLevel(price float64) {
	ob.asks = ob.asks.Delete(&priceLevel{Price: price})
}

func (ob *OrderBookTreap) updateAskPriceLevel(price float64, volume float64) {
	ob.deleteAskPriceLevel(price)
	ob.asks = ob.asks.Upsert(&priceLevel{Price: price, Volume: volume}, rand.Int())
}

// Returns a pointer to the Bid Price Level if it exists in the Bids Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) getBidPriceLevel(price float64) *priceLevel {
	pl := ob.bids.Get(&priceLevel{Price: price})
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) getMaxBidPriceLevel() *priceLevel {
	pl := ob.bids.Max()
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) deleteBidPriceLevel(price float64) {
	ob.bids = ob.bids.Delete(&priceLevel{Price: price})
}

func (ob *OrderBookTreap) updateBidPriceLevel(price float64, volume float64) {
	ob.deleteBidPriceLevel(price)
	ob.bids = ob.bids.Upsert(&priceLevel{Price: price, Volume: volume}, rand.Int())
}

/*******************************************************************************/

func NewOrderBookTreap(exchange int, pricepair string, msg []byte) (*OrderBookTreap, error) {
	if exchange == constants.CoinbasePro {
		bidTreap := gtreap.NewTreap(priceLevelAscendingCompare)
		askTreap := gtreap.NewTreap(priceLevelAscendingCompare)

		//Initialize bid treap
		i := 0
		jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			bidPrice, err := jsonparser.GetString(value, "[0]")
			bidVol, err := jsonparser.GetString(value, "[1]")
			bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
			bidVolF, err := strconv.ParseFloat(bidVol, 64)

			pricelevel := &priceLevel{Price: bidPriceF, Volume: bidVolF, Index: i}
			i++
			bidTreap = bidTreap.Upsert(pricelevel, rand.Int())
		}, "bids")

		//Initialize ask Treap
		i = 0
		jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			askPrice, err := jsonparser.GetString(value, "[0]")
			askVol, err := jsonparser.GetString(value, "[1]")
			askPriceF, err := strconv.ParseFloat(askPrice, 64)
			askVolF, err := strconv.ParseFloat(askVol, 64)

			pricelevel := &priceLevel{Price: askPriceF, Volume: askVolF, Index: i}
			i++
			askTreap = askTreap.Upsert(pricelevel, rand.Int())
		}, "asks")

		ob := &OrderBookTreap{Exchange: exchange, PricePair: pricepair, bids: bidTreap, asks: askTreap}
		return ob, nil
	}
	return nil, nil
}
