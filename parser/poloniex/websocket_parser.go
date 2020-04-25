package poloniex

import (
	"strconv"
	"strings"

	"../../common/structs"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
)

func poloSplitter(c rune) bool {
	return c == '[' || c == '"' || c == ']' || c == ','
}

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