package bookkeeper

/*
func TestInitBook(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}
	c := make(chan structs.PriceUpdate, 1000)
	bk := NewBookkeeper(c)

	err = bk.InitBook(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err != nil {
		t.Error(err)
		return
	}
	if bk.GetBooks()["ETH-BTC"] == nil {
		t.Errorf("Expected entry for pricepair: %s, got nil", "ETH-BTC")
		return
	}
	if bk.GetBooks()["ETH-BTC"][constants.CoinbasePro] == nil {
		t.Errorf("Expected entry for pricepair: %s exchange: %d, got nil", "ETH-BTC", constants.CoinbasePro)
	}

	//Expect error from attempting to insert existing pricepair/exchange combination
	err = bk.InitBook(constants.CoinbasePro, "ETH-BTC", snapshotResp)
	if err == nil {
		t.Errorf("Expected error from attempting to init existing pricepair + exchange")
		return
	}
	t.Log(err)
}
*/
