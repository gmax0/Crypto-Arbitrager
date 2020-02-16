package orderbook

import (
	"strconv"
	"testing"

	"../../common/constants"
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

		pl := cbOb.getMaxBidPriceLevel()
		if (*pl).Price != bidPriceF {
			t.Errorf("Price level %d, expected bid price: %f, got %f", i, bidPriceF, pl.Price)
		}
		if (*pl).Volume != bidVolF {
			t.Errorf("Price level %d, expected bid vol: %f, got %f", i, bidVolF, pl.Volume)
		}

		cbOb.deleteBidPriceLevel((*pl).Price)

		i++
	}, "bids")

	i = 0
	jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		askPrice, err := jsonparser.GetString(value, "[0]")
		askVol, err := jsonparser.GetString(value, "[1]")
		askPriceF, err := strconv.ParseFloat(askPrice, 64)
		askVolF, err := strconv.ParseFloat(askVol, 64)

		pl := cbOb.getMinAskPriceLevel()
		if (*pl).Price != askPriceF {
			t.Errorf("Price level %d, expected ask price: %f, got %f", i, askPriceF, pl.Price)
		}
		if (*pl).Volume != askVolF {
			t.Errorf("Price level %d, expected ask vol: %f, got %f", i, askVolF, pl.Volume)
		}

		cbOb.deleteAskPriceLevel((*pl).Price)
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
	pl := cbOb.getMinAskPriceLevel()
	if pl != nil {
		t.Error("Expected pl nil")
		return
	}
	pl = cbOb.getMaxBidPriceLevel()
	if pl != nil {
		t.Error("Expected pl nil")
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

	cbOb.insertAskPriceLevel(100.1, 1)
	pl = cbOb.getMinAskPriceLevel()
	if price := (*pl).Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.updateAskPriceLevel(100.1, 2)
	pl = cbOb.getMinAskPriceLevel()
	if price := (*pl).Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.insertAskPriceLevel(10.9, 1)
	pl = cbOb.getMinAskPriceLevel()
	if price := (*pl).Price; price != 10.9 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.deleteAskPriceLevel(10.9)
	pl = cbOb.getMinAskPriceLevel()
	if price := (*pl).Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 2 {
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
	cbOb.insertBidPriceLevel(100.1, 1)
	pl = cbOb.getMaxBidPriceLevel()
	if price := (*pl).Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.updateBidPriceLevel(100.1, 2)
	pl = cbOb.getMaxBidPriceLevel()
	if price := (*pl).Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.insertBidPriceLevel(200, 1)
	pl = cbOb.getMaxBidPriceLevel()
	if price := (*pl).Price; price != 200 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 1 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

	cbOb.deleteBidPriceLevel(200)
	pl = cbOb.getMaxBidPriceLevel()
	if price := (*pl).Price; price != 100.1 {
		t.Errorf("Expected volume %f, got %f", 100.1, price)
	}
	if volume := (*pl).Volume; volume != 2 {
		t.Errorf("Expected volume %f, got %f", 100.1, volume)
	}

}
