package orderbook

import (
	"strconv"
	"testing"

	"../../common/constants"
	"../../testhelpers"
	"github.com/buger/jsonparser"
)

func TestNewOrderBookTreap(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
	}

	cbOb, err := NewOrderBookTreap(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error("Error initializing new OrderBookTreap")
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
