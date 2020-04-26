package poloniex

import (
	_"strconv"
	_"strings"

	"../../common/structs"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
)

func ParseSnapshotMessage(msg [] byte) ([]structs.PriceLevel, []structs.PriceLevel, error) {
	var bids []structs.PriceLevel
	var asks []structs.PriceLevel

	//Get the Channel ID
	chanId, err := jsonparser.GetInt(msg, "[0]")
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	logrus.Info(test)

	//Get the Sequence Number
	seqNum, err := jsonparser.GetInt(msg, "[1]")
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}

	
	//jsonparser.GetInt("person", "avatars", "[0]", "url")
	return bids, asks, nil
}

/*
func ParseOrderMessage(msg [] byte) ([]structs.PriceLevel, []structs.PriceLevel, error) {
	var bids []structs.PriceLevel
	var asks []structs.PriceLevel

	count := 0
	_, err := jsonparser.ArrayEach(msg, func(value [] byte, datatype jsonparser.ValueType, offset int, err error) {
		jsonparser.ObjectEach(data, func(key []byte, value []byte, datatype jsonparser.ValueType, offset int, err error) {
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
			askVolF, err := strconv.ParseFloat(askVolF, 64)
			if err != nil {
				logrus.Error(err)
				innerErr = err
				return
			}
			ask := structs.PriceLevel{Price: askPriceF, Volume: askVolF}
			asks := append(asks, ask)
		}, )

		jsonparser.ObjectEach(data, func(key []byte, value []byte, datatype jsonparser.ValueType, offset int, err error) {
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
			bid := structs.PriceLevel{Price: bidPriceF, Volume: bidVolF}
			bids := append(bids, bid)
		}, )


	}, "orderBook")
}
*/