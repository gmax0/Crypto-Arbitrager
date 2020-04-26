package orderbook

import (
	"strconv"
	"testing"

	"../../common/constants"
	"../../testhelpers"
	"github.com/buger/jsonparser"
)

func TestNewOrderBookHeap(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
	}
	cb_ob, err := NewOrderBook(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBook")
	}

	if len(*cb_ob.Asks) != 8959 {
		t.Errorf("Expected Asks size: %d, got %d", 8959, len(*cb_ob.Asks))
	}
	if len(*cb_ob.Bids) != 2685 {
		t.Errorf("Expected Bids size: %d, got %d", 2685, len(*cb_ob.Bids))
	}

	//Ensure descending order
	i := 0
	jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		bidPrice, err := jsonparser.GetString(value, "[0]")
		bidVol, err := jsonparser.GetString(value, "[1]")
		bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
		bidVolF, err := strconv.ParseFloat(bidVol, 64)

		pl := cb_ob.getHighestBid()
		if pl.Price != bidPriceF {
			t.Errorf("Price level %d, expected bid price: %f, got %f", i, bidPriceF, pl.Price)
		}
		if pl.Volume != bidVolF {
			t.Errorf("Price level %d, expected bid vol: %f, got %f", i, bidVolF, pl.Volume)
		}

		popBidHeap(cb_ob.Bids)

		i++
	}, "bids")

	//Ensure ascending order
	i = 0
	jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		askPrice, err := jsonparser.GetString(value, "[0]")
		askVol, err := jsonparser.GetString(value, "[1]")
		askPriceF, err := strconv.ParseFloat(askPrice, 64)
		askVolF, err := strconv.ParseFloat(askVol, 64)

		pl := cb_ob.getLowestAsk()
		if pl.Price != askPriceF {
			t.Errorf("Price level %d, expected ask price: %f, got %f", i, askPriceF, pl.Price)
		}
		if pl.Volume != askVolF {
			t.Errorf("Price level %d, expected ask vol: %f, got %f", i, askVolF, pl.Volume)
		}

		popAskHeap(cb_ob.Asks)
		i++
	}, "asks")
}
