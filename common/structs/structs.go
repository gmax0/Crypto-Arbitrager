package structs

type Bid struct {
	Price  float64
	Volume float64
}

type Ask struct {
	Price  float64
	Volume float64
}

type PriceLevel struct {
	Price  float64
	Volume float64
	Index  int //TODO Remove this when OrderBookHeap is scrapped
}

// PriceUpdate is passed from bookkeeper through its outgoing channel
// Used for whenever a Min Ask or Max Bid changes for a given price pair
type PriceUpdate struct {
	UpdateType string
	Exchange   int
	PricePair  string
	UpdateAsk  Ask
	UpdateBid  Bid
}
