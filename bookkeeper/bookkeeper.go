package bookkeeper

import (
	"fmt"

	"../common/constants"
	"../common/structs"
	coinbaseProParser "../parser/coinbasepro"
	"./orderbook"
	"github.com/sirupsen/logrus"
)

/*******************************************************************************/

//Create
type UpdateMessage struct {
}

type Bookkeeper struct {
	books   map[string]map[int]*orderbook.OrderBookTreap //Maps price-pair -> exchange # -> orderbook
	outChan chan<- structs.PriceUpdate                   //Outgoing Max Bid or Min Ask update messages
}

/*******************************************************************************/

func NewBookkeeper(channel chan<- structs.PriceUpdate) *Bookkeeper {
	b := make(map[string]map[int]*orderbook.OrderBookTreap)
	bk := &Bookkeeper{books: b, outChan: channel}
	return bk
}

// InitBook will initialize an orderbook entry within the Bookkeeper's Bookkeeper.books map
// if not already initialized. An error is thrown if an existing entry is found.
// 	- pricepair must be in "A-B" format, analogous to A/B
// 	- See commons/constants for exchange values
func (bk *Bookkeeper) InitBook(exchange int, pricepair string, msg []byte) error {
	if bk.books[pricepair][exchange] != nil {
		return fmt.Errorf("Existing orderbook exists for pricepair: %s, exchange %d", pricepair, exchange)
	}

	ob, err := orderbook.NewOrderBookTreap(exchange, pricepair, msg)
	if err != nil {
		return err
	}

	if bk.books[pricepair] == nil {
		bk.books[pricepair] = make(map[int]*orderbook.OrderBookTreap)
	}
	bk.books[pricepair][exchange] = ob

	return nil
}

// TODO: Refactor orderbook calls to use Bid/Ask as argument

// ProcessPriceUpdate
func (bk *Bookkeeper) ProcessPriceUpdate(exchange int, pricepair string, msg []byte) {
	switch exchange {
	case constants.CoinbasePro:
		bidUpdates, askUpdates, err := coinbaseProParser.ParseUpdateMessage(msg)
		if err != nil {
			logrus.Error(err)
			return
		}
		for _, bu := range bidUpdates {
			bk.books[pricepair][exchange].UpsertBidPriceLevel(structs.PriceLevel{Price: bu.Price, Volume: bu.Volume})
		}
		for _, au := range askUpdates {
			bk.books[pricepair][exchange].UpsertAskPriceLevel(structs.PriceLevel{Price: au.Price, Volume: au.Volume})
		}
	}
}

func (bk *Bookkeeper) GetBooks() map[string]map[int]*orderbook.OrderBookTreap {
	return bk.books
}
