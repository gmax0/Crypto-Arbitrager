package bookkeeper

import (
	"fmt"

	"../common/structs"
	"../orderbook"
)

/*******************************************************************************/

type Bookkeeper struct {
	Exchange int
	books    map[int]*orderbook.OrderBookTreap //Maps price-pair -> orderbook (See common/constants.go for possible values)
	//outChan  chan<- structs.PriceUpdate        //Non-blocking, Outgoing Max Bid or Min Ask update messages
}

/*******************************************************************************/

func NewBookkeeper(exchange int) *Bookkeeper {
	return &Bookkeeper{Exchange: exchange, books: make(map[int]*orderbook.OrderBookTreap)}
}

func (bk *Bookkeeper) InitBook(pricepair int, bids []structs.PriceLevel, asks []structs.PriceLevel) error {
	if bk.books[pricepair] != nil {
		return fmt.Errorf("Existing orderbook exists for pricepair: %s, exchange %d", pricepair, bk.Exchange)
	}

	ob, err := orderbook.NewOrderBookTreap(bids, asks)
	if err != nil {
		return err
	}

	bk.books[pricepair] = ob

	return nil
}

func (bk *Bookkeeper) ProcessPriceUpdate(pricepair int, bidUpdates []structs.PriceLevel, askUpdates []structs.PriceLevel) {
	for _, bu := range bidUpdates {
		bk.books[pricepair].UpsertBidPriceLevel(structs.PriceLevel{Price: bu.Price, Volume: bu.Volume})
	}
	for _, au := range askUpdates {
		bk.books[pricepair].UpsertAskPriceLevel(structs.PriceLevel{Price: au.Price, Volume: au.Volume})
	}
}

func (bk *Bookkeeper) ClearBook(pricepair int) {
	bk.books[pricepair] = nil
}

func (bk *Bookkeeper) BookExists(pricepair int) bool {
	return bk.books[pricepair] != nil
}

func (bk *Bookkeeper) GetBooks(pricepair int) *orderbook.OrderBookTreap {
	return bk.books[pricepair]
}
