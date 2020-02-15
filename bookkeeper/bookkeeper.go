package bookkeeper

import (
	"fmt"

	"./orderbook"
)

//
type UpdateMessage struct {
}

type Bookkeeper struct {
	books map[string]map[int]*orderbook.OrderBook //Maps price-pair -> exchange # -> orderbook

}

func NewBookkeeper() *Bookkeeper {
	b := make(map[string]map[int]*orderbook.OrderBook)
	bk := &Bookkeeper{books: b}
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

	ob, err := orderbook.NewOrderBook(exchange, pricepair, msg)
	if err != nil {
		return err
	}

	if bk.books[pricepair] == nil {
		bk.books[pricepair] = make(map[int]*orderbook.OrderBook)
	}
	bk.books[pricepair][exchange] = ob

	return nil
}

func (bk *Bookkeeper) GetBooks() map[string]map[int]*orderbook.OrderBook {
	return bk.books
}
