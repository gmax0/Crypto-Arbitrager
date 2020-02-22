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

func TestNewOrderBookTreap(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
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

		cbOb.DeleteBidPriceLevel(pl.Price)

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

		cbOb.DeleteAskPriceLevel(pl.Price)
		i++
	}, "asks")
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
			expect Min Ask: [100.1, 2]
		INS [10.9 , 1]
			expect Min Ask: [10.9, 1]
		DEL [10.9, 1]
			expect Min Ask: [100.1, 2]
	*/

	cbOb.InsertAskPriceLevel(100.1, 1)
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.UpdateAskPriceLevel(100.1, 2)
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.InsertAskPriceLevel(10.9, 1)
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 10.9 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.DeleteAskPriceLevel(10.9)
	pl = cbOb.GetMinAskPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
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
	cbOb.InsertBidPriceLevel(100.1, 1)
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.UpdateBidPriceLevel(100.1, 2)
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.InsertBidPriceLevel(200, 1)
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 200 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.DeleteBidPriceLevel(200)
	pl = cbOb.GetMaxBidPriceLevel()
	if price := pl.Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := pl.Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
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

	cbOb.InsertAskPriceLevel(100, 1)
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

	cbOb.UpdateAskPriceLevel(100, 0)
	pl = cbOb.GetMinAskPriceLevel()
	if pl != (structs.PriceLevel{}) {
		t.Error("Expect empty pl")
		return
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

	cbOb.DeleteAskPriceLevel(100)
}
