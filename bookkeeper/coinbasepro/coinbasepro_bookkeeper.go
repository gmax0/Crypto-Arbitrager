package coinbase

import (
	"errors"
	"time"

	"../../bookkeeper"
	"../../common/constants"
	coinbaseproParser "../../parser/coinbasepro"
	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
)

type CoinbaseproBookkeeper struct {
	B                  *bookkeeper.Bookkeeper
	UpdateLastReceived time.Time
}

func NewCoinbaseproBookkeeper() *CoinbaseproBookkeeper {
	bk := bookkeeper.NewBookkeeper(constants.CoinbasePro)
	return &CoinbaseproBookkeeper{B: bk, UpdateLastReceived: time.Now()}
}

func (cbbk *CoinbaseproBookkeeper) HandleMessage(msg []byte) error {
	//Get message type
	msgType, err := jsonparser.GetString(msg, "type")
	if err != nil {
		log.Errorf("CoinbasePro Bookkeeper: Could not locate key: type")
		return err
	}
	//Get product_id and canoniclized price pair value
	pricePair, err := jsonparser.GetString(msg, "product_id")
	if err != nil {
		log.Errorf("CoinbasePro Bookkeeper: Could not locate key: product_id")
		return err
	}
	pricePairInt := constants.CBPricePairToInt[pricePair]
	if pricePairInt == 0 {
		log.Errorf("CoinbasePro Bookkeeper: Pricepair %d")
	}

	if msgType == "snapshot" {
		//Check if orderbook already exists for this pricePair
		if cbbk.B.BookExists(pricePairInt) {
			log.Infof("CoinbasePro Bookkeeper: PricePair %d already exists...deleting", pricePairInt)
			cbbk.B.ClearBook(pricePairInt)
		}
		//Parse the message for asks and bids
		bids, asks, err := coinbaseproParser.ParseSnapshotMessage(msg)
		if err != nil {
			log.Errorf("CoinbasePro Bookkeeper: PricePair %d, unable to parse snapshot message")
			log.Error(err)
			return err
		}

		//Create a new orderbook for this pricePair
		log.Infof("CoinbasePro Bookkeeper: Initializing orderbook for PricePair %d", pricePairInt)
		cbbk.B.InitBook(pricePairInt, bids, asks)
	} else if msgType == "l2update" {
		//Check if orderbook already exists for this pricePair
		if !cbbk.B.BookExists(pricePairInt) {
			log.Errorf("CoinbasePro Bookkeeper: l2update message received for uninitialized PricePair %d", pricePairInt)
			return errors.New("Error placeholder 1")
		}

		//Check if update message is out of order
		timestamp, err := jsonparser.GetString(msg, "time")
		if err != nil {
			log.Error(err)
			return errors.New("Error placeholder 2")
		}
		time, err := time.Parse(time.RFC3339, timestamp)
		if time.Before(cbbk.UpdateLastReceived) {
			//Out of order update message
			log.Errorf("CoinbasePro Bookkeeper: Received out of order l2update message for PricePair %d", pricePairInt)
		}

		cbbk.UpdateLastReceived = time

		//Parse the message for asks and bids
		bids, asks, err := coinbaseproParser.ParseUpdateMessage(msg)

		//Update the orderbook for this pricePair
		cbbk.B.ProcessPriceUpdate(pricePairInt, bids, asks)
	}
	return nil
}
