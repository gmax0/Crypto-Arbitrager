package bookkeeper

import (
	"testing"

	"../common/constants"
	"../common/structs"
	"../testhelpers"
)

func TestInitBook(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}
	c := make(chan structs.PriceUpdate, 1000)
	bk := NewBookkeeper(c)

	bk.InitBook(constants.CoinbasePro, "ETH-BTC", snapshotResp)

	//Expect error from attempting to insert existing pricepair/exchange combination

}
