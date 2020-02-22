package orderbook

import (
	"math/rand"
	"strconv"
	"time"

	"../../common/constants"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"github.com/steveyen/gtreap" //Note that this is an immutable treap, TODO: add a size() function to this
)

/* TODO:
*	- Consider modifying get price level functions to returning structs instead of pointers...esp if any of this is to be used in parallel
*
 */

func init() {
	rand.Seed(time.Now().UnixNano()) //Look into alternatives to derive probalistic key
}

type OrderBookTreap struct {
	Exchange  int
	PricePair string
	bids      *gtreap.Treap //Treap of *priceLevel
	asks      *gtreap.Treap //Treap of *priceLevel
}

/*******************************************************************************/

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

// Returns a pointer to the Ask Price Level if it exists in the Asks Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) GetAskPriceLevel(price float64) *priceLevel {
	pl := ob.asks.Get(&priceLevel{Price: price})
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) GetMinAskPriceLevel() *priceLevel {
	pl := ob.asks.Min()
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) DeleteAskPriceLevel(price float64) {
	ob.asks = ob.asks.Delete(&priceLevel{Price: price})
}

func (ob *OrderBookTreap) InsertAskPriceLevel(price float64, volume float64) {
	foundPl := ob.GetAskPriceLevel(price)
	if foundPl != nil {
		//Already exists
		return
	}
	insertPl := &priceLevel{Price: price, Volume: volume}
	ob.asks = ob.asks.Upsert(insertPl, rand.Int())
}

func (ob *OrderBookTreap) UpdateAskPriceLevel(price float64, volume float64) {
	ob.DeleteAskPriceLevel(price)
	ob.asks = ob.asks.Upsert(&priceLevel{Price: price, Volume: volume}, rand.Int())
}

// Returns a pointer to the Bid Price Level if it exists in the Bids Treap
// O(logN) assuming probalistic tree height is achieved
func (ob *OrderBookTreap) GetBidPriceLevel(price float64) *priceLevel {
	pl := ob.bids.Get(&priceLevel{Price: price})
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) GetMaxBidPriceLevel() *priceLevel {
	pl := ob.bids.Max()
	if pl == nil {
		return nil
	}
	return pl.(*priceLevel)
}

func (ob *OrderBookTreap) DeleteBidPriceLevel(price float64) {
	ob.bids = ob.bids.Delete(&priceLevel{Price: price})
}

func (ob *OrderBookTreap) UpdateBidPriceLevel(price float64, volume float64) {
	ob.DeleteBidPriceLevel(price)
	ob.bids = ob.bids.Upsert(&priceLevel{Price: price, Volume: volume}, rand.Int())
}

func (ob *OrderBookTreap) InsertBidPriceLevel(price float64, volume float64) {
	foundPl := ob.GetBidPriceLevel(price)
	if foundPl != nil {
		//Already exists
		return
	}
	insertPl := &priceLevel{Price: price, Volume: volume}
	ob.bids = ob.bids.Upsert(insertPl, rand.Int())
}

/*******************************************************************************/

// Known issue with callback error handling: https://github.com/buger/jsonparser/issues/129
func NewOrderBookTreap(exchange int, pricepair string, msg []byte) (*OrderBookTreap, error) {
	if exchange == constants.CoinbasePro {
		bidTreap := gtreap.NewTreap(priceLevelAscendingCompare)
		askTreap := gtreap.NewTreap(priceLevelAscendingCompare)

		//Initialize bid treap
		i := 0
		var innerErr error
		_, err := jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			if innerErr != nil {
				//Skip callback iteration if an error was detected previously...
				return
			}
			bidPrice, err := jsonparser.GetString(value, "[0]")
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			bidVol, err := jsonparser.GetString(value, "[1]")
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			bidVolF, err := strconv.ParseFloat(bidVol, 64)
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}

			pricelevel := &priceLevel{Price: bidPriceF, Volume: bidVolF, Index: i}
			i++
			bidTreap = bidTreap.Upsert(pricelevel, rand.Int())
		}, "bids")

		//Handle ArrayEach error
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		//Handle ArrayEach callback error
		if innerErr != nil {
			return nil, innerErr
		}

		//Initialize ask Treap
		i = 0
		_, err = jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			if innerErr != nil {
				//Skip callback iteration if an error was detected previously...
				return
			}
			askPrice, err := jsonparser.GetString(value, "[0]")
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			askVol, err := jsonparser.GetString(value, "[1]")
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			askPriceF, err := strconv.ParseFloat(askPrice, 64)
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			askVolF, err := strconv.ParseFloat(askVol, 64)
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}

			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}

			pricelevel := &priceLevel{Price: askPriceF, Volume: askVolF, Index: i}
			i++
			askTreap = askTreap.Upsert(pricelevel, rand.Int())
		}, "asks")

		//Handle ArrayEach error
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		//Handle ArrayEach callback error
		if innerErr != nil {
			return nil, innerErr
		}

		ob := &OrderBookTreap{Exchange: exchange, PricePair: pricepair, bids: bidTreap, asks: askTreap}
		return ob, nil
	}
	return nil, nil
}
