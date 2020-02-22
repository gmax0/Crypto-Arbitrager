package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"./bookkeeper"
	coinbasepro "./client/coinbasepro"
	"./common/constants"
	_ "./config"

	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	//Read configs in here
}

func worker() {

}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c1 := make(chan []byte, 1000)
	// c1 := make(chan []byte)

	//Setup Bookkeeper
	_ = make(chan []byte)
	bk := bookkeeper.NewBookkeeper()

	//Setup coinbasepro Client Thread
	cbpClient, err := coinbasepro.NewClient("ws-feed.pro.coinbase.com")
	if err != nil {
		log.Fatal("Unable to initialize coinbasepro Client:", err)
		return
	}
	defer cbpClient.CloseUnderlyingConnection()

	fmt.Println(cbpClient)

	//Setup JSON Message
	jsonFile, err := ioutil.ReadFile("./testhelpers/testdata/coinbasepro/test-l2-subscribe.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	jsonFile2, err := ioutil.ReadFile("./testhelpers/testdata/coinbasepro/test-l2-unsubscribe.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	//
	// cbpPricePairs := []string{"ETH-BTC", "ETH-USD"}

	//Initialize CoinbasePro Pricebook
	//cbpPb := pricebook.NewPricebook(pricebook.CoinbasePro, cbpPricePairs)
	//log.Info(cbpPb)

	/*
		err = cbpClient.Subscribe(jsonFile)
		if err != nil {
			log.Fatal("coinbasepro Client write error:", err)
			return
		}
	*/
	cbpClient.SetSubscribeMessage(jsonFile)
	cbpClient.SetUnsubscribeMessage(jsonFile2)

	// go client.
	go cbpClient.StartStreaming(c1, interrupt)

	maxSizeReached := 0
	msgReceived := 0
	for {
		select {
		case message := <-c1:
			chanSize := len(c1)
			msgReceived++
			if chanSize > maxSizeReached {
				maxSizeReached = chanSize
			}
			log.Info("Channel size: ", chanSize)
			log.Debug(message)

			msgType, err := jsonparser.GetString(message, "type")
			if err != nil {
				log.Error(err)
				return
			}

			if msgType == "snapshot" {
				pricePair, err := jsonparser.GetString(message, "product_id")
				if err != nil {
					log.Error(err)
					return
				}

				//Initialize the orderbook for price pair on exchange
				err = bk.InitBook(constants.CoinbasePro, pricePair, message)
				if err != nil {
					log.Error(err)
				}
			} else if msgType == "l2update" {
				pricePair, err := jsonparser.GetString(message, "product_id")
				if err != nil {
					log.Error(err)
					return
				}
				bk.ProcessPriceUpdate(constants.CoinbasePro, pricePair, message)

			}
		case <-interrupt:
			log.Println("interrupt")
			err = cbpClient.StopStreaming()
			if err != nil {
				log.Error("StopStreaming() error: ", err)
			}

			err = cbpClient.CloseConnection()
			if err != nil {
				log.Error("coinbasepro Client write close error:", err)
			}

			log.Info("Max Chan Size Reached: ", maxSizeReached)
			log.Info("Messages Received: ", msgReceived)
			log.Info(bk.GetBooks())

			log.Info(bk.GetBooks())
			return
		}
	}
}
