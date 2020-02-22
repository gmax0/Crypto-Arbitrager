package coinbasepro

import (
	"strconv"
	"testing"

	"../../testhelpers"
	"github.com/buger/jsonparser"
)

func TestParseSnapshotMessage(t *testing.T) {
	snapshotResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-snapshot-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	bids, asks, err := ParseSnapshotMessage(snapshotResp)
	if err != nil {
		t.Error("Not expecting error.")
	}

	i := 0
	var innerErr error
	_, err = jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		if innerErr != nil {
			//Skip callback iteration if an error was detected previously...
			return
		}

		bidPrice, err := jsonparser.GetString(value, "[0]")
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}
		bidVol, err := jsonparser.GetString(value, "[1]")
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}
		bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}
		bidVolF, err := strconv.ParseFloat(bidVol, 64)
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}

		if bidPriceF != bids[i].Price {
			t.Errorf("Expected price: %f, got %f", bidPriceF, bids[i].Price)
		}
		if bidVolF != bids[i].Volume {
			t.Errorf("Expected volume: %f, got %f", bidVolF, bids[i].Volume)
		}
		i++
	}, "bids")

	i = 0
	_, err = jsonparser.ArrayEach(snapshotResp, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		if innerErr != nil {
			//Skip callback iteration if an error was detected previously...
			return
		}
		askPrice, err := jsonparser.GetString(value, "[0]")
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}
		askVol, err := jsonparser.GetString(value, "[1]")
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}
		askPriceF, err := strconv.ParseFloat(askPrice, 64)
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}
		askVolF, err := strconv.ParseFloat(askVol, 64)
		if err != nil {
			t.Error(err)
			innerErr = err
			return
		}

		if askPriceF != asks[i].Price {
			t.Errorf("Expected price: %f, got %f", askPriceF, asks[i].Price)
		}
		if askVolF != asks[i].Volume {
			t.Errorf("Expected volume: %f, got %f", askVolF, asks[i].Volume)
		}
		i++
	}, "asks")
}

func TestParseUpdateMessage(t *testing.T) {
	//Single Update Test
	updateResp, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-update-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	bidUpdates, askUpdates, err := ParseUpdateMessage(updateResp)
	if err != nil {
		t.Error(err)
		return
	}

	if len(bidUpdates) != 1 {
		t.Errorf("Expected Bid Update length: %d, got: %d", 1, len(bidUpdates))
		return
	}
	if bidUpdates[0].Price != 263.11 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 263.11, bidUpdates[0].Price)
		return
	}
	if bidUpdates[0].Volume != 11.75216755 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 11.75216755, bidUpdates[0].Volume)
		return
	}
	if len(askUpdates) != 0 {
		t.Errorf("Expected Ask Update length: %d, got: %d", 0, len(askUpdates))
		return
	}

	//Multi Update Test
	updateResp, err = testhelpers.GetMsgFromFile("../../testhelpers/testdata/coinbasepro/test-l2-update-multi-response.json")
	if err != nil {
		t.Error("Could not read test file.")
		return
	}

	bidUpdates, askUpdates, err = ParseUpdateMessage(updateResp)
	if err != nil {
		t.Error(err)
		return
	}

	if len(bidUpdates) != 2 {
		t.Errorf("Expected Bid Update length: %d, got: %d", 2, len(bidUpdates))
		return
	}
	if bidUpdates[0].Price != 263.11 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 263.11, bidUpdates[0].Price)
		return
	}
	if bidUpdates[0].Volume != 1 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 1.0, bidUpdates[0].Volume)
		return
	}
	if bidUpdates[1].Price != 263.11 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 263.11, bidUpdates[0].Price)
		return
	}
	if bidUpdates[1].Volume != 11.75216755 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 11.75216755, bidUpdates[0].Volume)
		return
	}

	if len(askUpdates) != 1 {
		t.Errorf("Expected Ask Update length: %d, got: %d", 1, len(askUpdates))
		return
	}
	if askUpdates[0].Price != 264 {
		t.Errorf("Expected Ask Update price: %f, got: %f", 264.0, bidUpdates[0].Price)
		return
	}
	if askUpdates[0].Volume != 10.0 {
		t.Errorf("Expected Bid Update price: %f, got: %f", 10.0, bidUpdates[0].Volume)
		return
	}

}
