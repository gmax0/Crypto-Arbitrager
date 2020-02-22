package coinbasepro

import (
	"strconv"
	"strings"

	"../../common/structs"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
)

func coinbaseSplitter(c rune) bool {
	return c == '[' || c == '"' || c == ']' || c == ','
}

// TODO: cleanup log statements

// ParseSnapshotMessage will parse the JSON Snapshot Message received through the
// CoinbasePro Websocket Feed and return the data in canoniclized structures
// * Known issue with callback error handling: https://github.com/buger/jsonparser/issues/129
func ParseSnapshotMessage(msg []byte) ([]structs.Bid, []structs.Ask, error) {
	var bids []structs.Bid
	var asks []structs.Ask

	var innerErr error

	//Parse Bid data
	_, err := jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		if innerErr != nil {
			//Skip callback iteration if an error was detected previously...
			return
		}

		bidPrice, err := jsonparser.GetString(value, "[0]")
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}
		bidVol, err := jsonparser.GetString(value, "[1]")
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}
		bidPriceF, err := strconv.ParseFloat(bidPrice, 64)
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}
		bidVolF, err := strconv.ParseFloat(bidVol, 64)
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}

		bid := &structs.Bid{Price: bidPriceF, Volume: bidVolF}
		bids = append(bids, *bid)
	}, "bids")

	//Handle ArrayEach error
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	//Handle ArrayEach callback error
	if innerErr != nil {
		return nil, nil, innerErr
	}

	//Parse Ask data
	_, err = jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		if innerErr != nil {
			//Skip callback iteration if an error was detected previously...
			return
		}
		askPrice, err := jsonparser.GetString(value, "[0]")
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}
		askVol, err := jsonparser.GetString(value, "[1]")
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}
		askPriceF, err := strconv.ParseFloat(askPrice, 64)
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}
		askVolF, err := strconv.ParseFloat(askVol, 64)
		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}

		if err != nil {
			logrus.Error(err)
			innerErr = err
			return
		}

		ask := &structs.Ask{Price: askPriceF, Volume: askVolF}
		asks = append(asks, *ask)
	}, "asks")

	//Handle ArrayEach error
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	//Handle ArrayEach callback error
	if innerErr != nil {
		return nil, nil, innerErr
	}

	return bids, asks, nil
}

// TODO: cleanup log statements

// ParseUpdateMessage will parse the JSON Update Message received through the
// CoinbasePro Websocket Feed and return the data in canoniclized structures
// Returns variable sized []structs.Bid, []structs.Ask in case multiple updates are present
// in a single message
func ParseUpdateMessage(msg []byte) ([]structs.Bid, []structs.Ask, error) {
	var bidUpdates []structs.Bid
	var askUpdates []structs.Ask

	var innerErr error

	_, err := jsonparser.ArrayEach(msg, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
		if innerErr != nil {
			return
		}

		//Example value passed to callback: ["buy", "100", "1"]
		rawUpdate := string(value)
		logrus.Debug(rawUpdate)

		parsedUpdate := strings.FieldsFunc(rawUpdate, coinbaseSplitter)
		logrus.Info(parsedUpdate)

		price, innerErr := strconv.ParseFloat(parsedUpdate[1], 64)
		if innerErr != nil {
			logrus.Error(innerErr)
			return
		}
		volume, innerErr := strconv.ParseFloat(parsedUpdate[2], 64)
		if innerErr != nil {
			logrus.Error(innerErr)
			return
		}

		logrus.Info(parsedUpdate[1])
		logrus.Info(parsedUpdate[2])

		if parsedUpdate[0] == "buy" {
			logrus.Info("BUY!")
			bu := &structs.Bid{Price: price, Volume: volume}
			bidUpdates = append(bidUpdates, *bu)
			return
		}

		if parsedUpdate[0] == "sell" {
			logrus.Info("SELL!")
			au := &structs.Ask{Price: price, Volume: volume}
			askUpdates = append(askUpdates, *au)
			return
		}
	}, "changes")

	//Handle ArrayEach outer error
	if err != nil {
		return nil, nil, err
	}
	if innerErr != nil {
		return nil, nil, innerErr
	}

	return bidUpdates, askUpdates, nil
}
