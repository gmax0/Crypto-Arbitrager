package orderbook

import (
	"math/rand"
	"time"

	"../../common/constants"
	"../../common/structs"
	coinbaseProParser "../../parser/coinbasepro"
	"github.com/steveyen/gtreap" //Note that this is an immutable treap, TODO: add a size() function to this
)

/* NOTES:
 *	- Though gtreap.Treap is thread-safe, OrderBookTreap is NOT due to MaxBid, MinAsk, BidSize, AskSize
 * TODO:
 *	- Consider modifying get price level functions to returning structs instead of pointers...esp if any of this is to be used in parallel
 *	- Use canoniclized structures to update price levels
 */

func init() {
	rand.Seed(time.Now().UnixNano()) //Look into alternatives to derive probalistic key
}

type OrderBookTreap struct {
	Exchange  int
	PricePair string
	maxBid    structs.PriceLevel //When creating an instance of OrderBookTreap, make sure the price is Min float64
	minAsk    structs.PriceLevel //When creating an instance of OrderBookTreap, make sure the price is Max float64
	bidSize   int
	askSize   int
	bids      *gtreap.Treap //Treap of structs.PriceLevel
	asks      *gtreap.Treap //Treap of structs.PriceLevel
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
func (ob *OrderBookTreap) getAskPriceLevel(pl structs.PriceLevel) structs.PriceLevel {
	foundPl := ob.asks.Get(pl)
	if foundPl == nil {
		return structs.PriceLevel{}
	}
	return foundPl.(structs.PriceLevel)
}

func (ob *OrderBookTreap) getMinAskPriceLevel() structs.PriceLevel {
	pl := ob.asks.Min()
	if pl == nil {
		return structs.PriceLevel{}
	}
	return pl.(structs.PriceLevel)
}

//Amortize
func (ob *OrderBookTreap) GetMinAskPriceLevel() structs.PriceLevel {
	return ob.minAsk
}

func (ob *OrderBookTreap) DeleteAskPriceLevel(pl structs.PriceLevel) {
	foundPl := ob.getAskPriceLevel(pl)

	if foundPl != (structs.PriceLevel{}) {
		ob.asks = ob.asks.Delete(structs.PriceLevel{Price: pl.Price})
		ob.askSize-- //Decrement size as we've deleted an existing price level

		if ob.minAsk == (structs.PriceLevel{}) {
			return //This condition should never execute
		}
		if pl.Price == ob.minAsk.Price {
			newMinAskPriceLevel := ob.getMinAskPriceLevel()
			ob.minAsk = newMinAskPriceLevel
		}
	}
}

func (ob *OrderBookTreap) UpsertAskPriceLevel(pl structs.PriceLevel) {
	ob.DeleteAskPriceLevel(pl)

	if pl.Volume == float64(0) {
		return
	}

	ob.asks = ob.asks.Upsert(structs.PriceLevel{Price: pl.Price, Volume: pl.Volume}, rand.Int())
	ob.askSize++ //Increment size as we've inserted a new price level

	if ob.minAsk == (structs.PriceLevel{}) {
		ob.minAsk = pl
		return
	}
	if pl.Price < ob.minAsk.Price {
		ob.minAsk = pl //Update Max Bid accordingly
	}
}

/*******************************************************************************/

// Returns a copy of the struct equivalent to the Bid Price Level if it exists in the Bids Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) getBidPriceLevel(pl structs.PriceLevel) structs.PriceLevel {
	foundPl := ob.bids.Get(pl)
	if foundPl == nil {
		return structs.PriceLevel{}
	}
	return foundPl.(structs.PriceLevel)
}

func (ob *OrderBookTreap) getMaxBidPriceLevel() structs.PriceLevel {
	pl := ob.bids.Max()
	if pl == nil {
		return structs.PriceLevel{}
	}
	return pl.(structs.PriceLevel)
}

// Amortized
func (ob *OrderBookTreap) GetMaxBidPriceLevel() structs.PriceLevel {
	return ob.maxBid
}

// DeleteBidPriceLevel will delete the price level with a price equal to the price of the argument
// If no matching price level is found, nothing will be modified
// If a matching price level is found, it will be deleted, the bidSize counter decremented, and the
// maxBid updated accordingly if the price deleted was the previous maxBid (Note that the new maxBid may be an empty price level,
// indicating an empty treap)
func (ob *OrderBookTreap) DeleteBidPriceLevel(pl structs.PriceLevel) {
	foundPl := ob.getBidPriceLevel(pl)

	if foundPl != (structs.PriceLevel{}) {
		ob.bids = ob.bids.Delete(structs.PriceLevel{Price: pl.Price})
		ob.bidSize-- //Decrement size as we've deleted an existing price level

		if ob.maxBid == (structs.PriceLevel{}) {
			return //This condition should never execute
		}
		if pl.Price == ob.maxBid.Price {
			newMaxBidPriceLevel := ob.getMaxBidPriceLevel()
			ob.maxBid = newMaxBidPriceLevel
		}
	}
}

func (ob *OrderBookTreap) UpsertBidPriceLevel(pl structs.PriceLevel) {
	ob.DeleteBidPriceLevel(pl)

	if pl.Volume == float64(0) {
		return
	}

	ob.bids = ob.bids.Upsert(structs.PriceLevel{Price: pl.Price, Volume: pl.Volume}, rand.Int())
	ob.bidSize++ //Increment size as we've inserted a new price level

	if ob.maxBid == (structs.PriceLevel{}) {
		ob.maxBid = pl
		return
	}
	if pl.Price > ob.maxBid.Price {
		ob.maxBid = pl //Update Max Bid accordingly
	}
}

/*******************************************************************************/

// Known issue with callback error handling: https://github.com/buger/jsonparser/issues/129
func NewOrderBookTreap(exchange int, pricepair string, msg []byte) (*OrderBookTreap, error) {
	var bids []structs.PriceLevel
	var asks []structs.PriceLevel
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

	mBid := structs.PriceLevel{}
	mAsk := structs.PriceLevel{}

	ob := &OrderBookTreap{Exchange: exchange, PricePair: pricepair, maxBid: mBid, minAsk: mAsk, bidSize: 0, askSize: 0,
		bids: bidTreap, asks: askTreap}

	for _, b := range bids {
		pricelevel := structs.PriceLevel{Price: b.Price, Volume: b.Volume}
		ob.UpsertBidPriceLevel(pricelevel)
	}
	for _, a := range asks {
		pricelevel := structs.PriceLevel{Price: a.Price, Volume: a.Volume}
		ob.UpsertAskPriceLevel(pricelevel)
	}

	return ob, nil
}
