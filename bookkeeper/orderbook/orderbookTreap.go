package orderbook

import (
	"math/rand"
	"time"

	"../../common/constants"
	"../../common/structs"
	coinbaseProParser "../../parser/coinbasepro"
	"github.com/steveyen/gtreap" //Note that this is an immutable treap, TODO: add a size() function to this
)

/* TODO:
*	- Consider modifying get price level functions to returning structs instead of pointers...esp if any of this is to be used in parallel
*	- Use canoniclized structures to update price levels
 */

func init() {
	rand.Seed(time.Now().UnixNano()) //Look into alternatives to derive probalistic key
}

type OrderBookTreap struct {
	Exchange  int
	PricePair string
	bids      *gtreap.Treap //Treap of *priceLevel
	asks      *gtreap.Treap //Treap of *priceLevel
	MaxBid    structs.Bid
	MinAsk    structs.Ask
}

/*******************************************************************************/

func priceLevelAscendingCompare(a, b interface{}) int {
	if a.(structs.PriceLevel).Price < b.(structs.PriceLevel).Price {
		return -1
	} else if a.(structs.PriceLevel).Price > b.(structs.PriceLevel).Price {
		return 1
	} else {
		return 0
	}
}

// Returns a copy of the struct equivalent to the Ask Price Level if it exists in the Asks Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) getAskPriceLevel(price float64) structs.PriceLevel {
	pl := ob.asks.Get(structs.PriceLevel{Price: price})
	if pl == nil {
		return structs.PriceLevel{}
	}
	return pl.(structs.PriceLevel)
}

func (ob *OrderBookTreap) GetMinAskPriceLevel() structs.PriceLevel {
	pl := ob.asks.Min()
	if pl == nil {
		return structs.PriceLevel{}
	}
	return pl.(structs.PriceLevel)
}

func (ob *OrderBookTreap) DeleteAskPriceLevel(price float64) {
	ob.asks = ob.asks.Delete(structs.PriceLevel{Price: price})
}

func (ob *OrderBookTreap) InsertAskPriceLevel(price float64, volume float64) {
	foundPl := ob.getAskPriceLevel(price)
	if foundPl != (structs.PriceLevel{}) {
		//Already exists
		return
	}
	insertPl := structs.PriceLevel{Price: price, Volume: volume}
	ob.asks = ob.asks.Upsert(insertPl, rand.Int())
}

func (ob *OrderBookTreap) UpdateAskPriceLevel(price float64, volume float64) {
	ob.DeleteAskPriceLevel(price)
	if volume == float64(0) {
		return
	}
	ob.asks = ob.asks.Upsert(structs.PriceLevel{Price: price, Volume: volume}, rand.Int())
}

// Returns a copy of the struct equivalent to the Bid Price Level if it exists in the Bids Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) getBidPriceLevel(price float64) structs.PriceLevel {
	pl := ob.bids.Get(structs.PriceLevel{Price: price})
	if pl == nil {
		return structs.PriceLevel{}
	}
	return pl.(structs.PriceLevel)
}

func (ob *OrderBookTreap) GetMaxBidPriceLevel() structs.PriceLevel {
	pl := ob.bids.Max()
	if pl == nil {
		return structs.PriceLevel{}
	}
	return pl.(structs.PriceLevel)
}

func (ob *OrderBookTreap) DeleteBidPriceLevel(price float64) {
	ob.bids = ob.bids.Delete(structs.PriceLevel{Price: price})
}

func (ob *OrderBookTreap) UpdateBidPriceLevel(price float64, volume float64) {
	ob.DeleteBidPriceLevel(price)
	if volume == float64(0) {
		return
	}
	ob.bids = ob.bids.Upsert(structs.PriceLevel{Price: price, Volume: volume}, rand.Int())
}

func (ob *OrderBookTreap) InsertBidPriceLevel(price float64, volume float64) {
	foundPl := ob.getBidPriceLevel(price)
	if foundPl != (structs.PriceLevel{}) {
		//Already exists
		return
	}
	insertPl := structs.PriceLevel{Price: price, Volume: volume}
	ob.bids = ob.bids.Upsert(insertPl, rand.Int())
}

/*******************************************************************************/

// Known issue with callback error handling: https://github.com/buger/jsonparser/issues/129
func NewOrderBookTreap(exchange int, pricepair string, msg []byte) (*OrderBookTreap, error) {
	var bids []structs.Bid
	var asks []structs.Ask
	var err error

	switch exchange {
	case constants.CoinbasePro:
		bids, asks, err = coinbaseProParser.ParseSnapshotMessage(msg)
		break
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	bidTreap := gtreap.NewTreap(priceLevelAscendingCompare)
	askTreap := gtreap.NewTreap(priceLevelAscendingCompare)

	for _, b := range bids {
		pricelevel := structs.PriceLevel{Price: b.Price, Volume: b.Volume}
		bidTreap = bidTreap.Upsert(pricelevel, rand.Int())
	}
	for _, a := range asks {
		pricelevel := structs.PriceLevel{Price: a.Price, Volume: a.Volume}
		askTreap = askTreap.Upsert(pricelevel, rand.Int())
	}

	ob := &OrderBookTreap{Exchange: exchange, PricePair: pricepair, bids: bidTreap, asks: askTreap}
	return ob, nil
}
