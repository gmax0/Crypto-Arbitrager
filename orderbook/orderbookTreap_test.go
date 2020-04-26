package orderbook

import (
	"strconv"
	"testing"

	"../../common/constants"
	"../../common/structs"
	"../../testhelpers"
	"github.com/buger/jsonparser"
)

//TODO: Benchmark tests for treap operations

/*
func TestNewOrderBookTreapEmpty(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	//Initialize, populate treaps based off of L2 Snapshot Message contents
	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
		return
	}
}
*/

func TestNewOrderBookTreap(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	//Initialize, populate treaps based off of L2 Snapshot Message contents
	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
		return
	}

	//Check that sizes have been updated accordingly
	if cbOb.bidSize != 2685 {
		t.Errorf("Expected initial BidSize of 2685, got %d", cbOb.bidSize)
		return
	}
	if cbOb.askSize != 8959 {
		t.Errorf("Expected initial AskSize of 8959, got %d", cbOb.askSize)
		return
	}
	//Check initial Min/Max price values
	if cbOb.maxBid.Price != 188.51 {
		t.Errorf("Expected initial MaxBid.Price of 188.51, got %f", cbOb.maxBid.Price)
		return
	}
	if cbOb.minAsk.Price != 188.52 {
		t.Errorf("Expected initial MinAsk.Price of 188.52, got %f", cbOb.minAsk.Price)
		return
	}

	i := 0
	jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		bidPrice, err := jsonparser.GetString(value, "[0]")
		bidVol, err := jsonparser.GetString(value, "[1]")
		bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
		bidVolF, err := strconv.ParseFloat(bidVol, 64)

		pl := cbOb.GetMaxBidPriceLevel()
		if pl.Price != bidPriceF {
			t.Errorf("Price level %d, expected bid price: %f, got %f", i, bidPriceF, pl.Price)
		}
		if pl.Volume != bidVolF {
			t.Errorf("Price level %d, expected bid vol: %f, got %f", i, bidVolF, pl.Volume)
		}

		cbOb.DeleteBidPriceLevel(pl)

		i++
	}, "bids")

	i = 0
	jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		askPrice, err := jsonparser.GetString(value, "[0]")
		askVol, err := jsonparser.GetString(value, "[1]")
		askPriceF, err := strconv.ParseFloat(askPrice, 64)
		askVolF, err := strconv.ParseFloat(askVol, 64)

		pl := cbOb.GetMinAskPriceLevel()
		if pl.Price != askPriceF {
			t.Errorf("Price level %d, expected ask price: %f, got %f", i, askPriceF, pl.Price)
		}
		if pl.Volume != askVolF {
			t.Errorf("Price level %d, expected ask vol: %f, got %f", i, askVolF, pl.Volume)
		}

		cbOb.DeleteAskPriceLevel(pl)
		i++
	}, "asks")

	if cbOb.bidSize != 0 {
		t.Error("Expected final bidSize 0")
		return
	}
	if cbOb.askSize != 0 {
		t.Error("Expected final askSize 0")
		return
	}
}

func TestTreapOperations(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response-empty.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
		return
	}

	//Both Bids + Ask Treaps should initially be empty
	pl := cbOb.GetMinAskPriceLevel()
	if pl != (structs.PriceLevel{}) {
		t.Error("Expected pl empty")
		return
	}
	pl = cbOb.GetMaxBidPriceLevel()
	if pl != (structs.PriceLevel{}) {
		t.Error("Expected pl empty")
		return
	}

	//Ask Treap Operations:
	/*
		INS [100.1, 1]
			expect Min Ask: [100.1, 1]
		UPD [100.1, 1] -> [100.1, 2]
			expect Min Ask: [100, 2]
		INS [10.9 , 1]
			expect Min Ask: [10.9, 1]
		DEL [10.9, 1]
			expect Min Ask: [100.1, 2]
	*/

	cbOb.UpsertAskPriceLevel(structs.PriceLevel{Price: 100.1, Volume: 1})
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected price %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 1.0, volume)
	}
	if size := cbOb.askSize; size != 1 {
		t.Errorf("Expected ask treap size: 1, got %d", size)
	}

	cbOb.UpsertAskPriceLevel(structs.PriceLevel{Price: 100.1, Volume: 2})
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected price %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 2.0, volume)
	}
	if size := cbOb.askSize; size != 1 {
		t.Errorf("Expected ask treap size: 2, got %d", size)
	}

	cbOb.UpsertAskPriceLevel(structs.PriceLevel{Price: 10.9, Volume: 1})
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 10.9 {
		t.Errorf("Expected price %f, got %f", 10.9, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 1.0, volume)
	}
	if size := cbOb.askSize; size != 2 {
		t.Errorf("Expected ask treap size: 2, got %d", size)
	}

	cbOb.DeleteAskPriceLevel(structs.PriceLevel{Price: 10.9})
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected price %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}
	if size := cbOb.askSize; size != 1 {
		t.Errorf("Expected ask treap size: 2, got %d", size)
	}

	//Bid Treap Operations:
	/*
		INS [100.1, 1]
			expect Max Bid: [100.1, 1]
		UPD [100.1, 1] -> [100.1, 2]
			expect Max Bid: [100.1, 2]
		INS [200, 1]
			expect Max Bid: [200, 1]
		DEL [200, 1]
			expect Max Bid: [100.1, 2]
	*/
	cbOb.UpsertBidPriceLevel(structs.PriceLevel{Price: 100.1, Volume: 1})
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}
	if size := cbOb.bidSize; size != 1 {
		t.Errorf("Expected bid treap size: 1, got %d", size)
	}

	cbOb.UpsertBidPriceLevel(structs.PriceLevel{Price: 100.1, Volume: 2})
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}
	if size := cbOb.bidSize; size != 1 {
		t.Errorf("Expected bid treap size: 1, got %d", size)
	}

	cbOb.UpsertBidPriceLevel(structs.PriceLevel{Price: 200, Volume: 1})
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 200 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}
	if size := cbOb.bidSize; size != 2 {
		t.Errorf("Expected bid treap size: 2, got %d", size)
	}

	cbOb.DeleteBidPriceLevel(structs.PriceLevel{Price: 200})
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}
	if size := cbOb.bidSize; size != 1 {
		t.Errorf("Expected bid treap size: 1, got %d", size)
	}

}

func TestTreapUpdateZeroVolume(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response-empty.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
		return
	}

	cbOb.UpsertAskPriceLevel(structs.PriceLevel{Price: 100, Volume: 1})
	pl := cbOb.GetMinAskPriceLevel()
	if pl == (structs.PriceLevel{}) {
		t.Error("Expected non empty pl")
		return
	}
	if cbOb.GetMinAskPriceLevel().Price != 100.0 {
		t.Errorf("Expected price level price: %f, got %f", 100.0, cbOb.GetMinAskPriceLevel().Price)
		return
	}
	if cbOb.GetMinAskPriceLevel().Volume != 1.0 {
		t.Errorf("Expected price level volume: %f, got %f", 100.0, cbOb.GetMinAskPriceLevel().Volume)
		return
	}
	if size := cbOb.askSize; size != 1 {
		t.Errorf("Expected ask treap size: 1, got %d", size)
	}

	cbOb.UpsertAskPriceLevel(structs.PriceLevel{Price: 100, Volume: 0})
	pl = cbOb.GetMinAskPriceLevel()
	if pl != (structs.PriceLevel{}) {
		t.Error("Expect empty pl")
		return
	}
	if size := cbOb.askSize; size != 0 {
		t.Errorf("Expected ask treap size: 0, got %d", size)
	}
}

func TestTreapDeleteNonExistent(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response-empty.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
		return
	}

	cbOb.DeleteAskPriceLevel(structs.PriceLevel{Price: 0.0})
}
