package orderbook

import (
	"testing"
    "../../common/constants"
	"../../testhelpers"
)

func TestNewOrderBook(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
	}
    orderbook, err := NewOrderBook(constants.CoinbasePro, "ETH-BTC", snapshotResp)
    t.Log(orderbook)
    /*
	orderbook, err := NewOrderBook(constants.CoinbasePro, "ETH-BTC", snapshotResp)

	t.Log(orderbook)
    t.Log(err)
    */


}
