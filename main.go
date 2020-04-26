package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	cbpBookkeeper "./bookkeeper/coinbasepro"

	_ "./client/coinbasepro"
	websocket "./client/websocket"
	_ "./config"

	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper"
)

var log = logrus.New()

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
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
	c2 := make(chan []byte, 1000)
	// c1 := make(chan []byte)

	//Setup Bookkeepers
	//cBk := make(chan structs.PriceUpdate, 1000)
	cbpBk := cbpBookkeeper.NewCoinbaseproBookkeeper()

	//Setup coinbasepro Client Thread
	cbpClient, err := websocket.NewClient("ws-feed.pro.coinbase.com")
	poloClient, err := websocket.NewClient("api2.poloniex.com")
	if err != nil {
		log.Fatal("Unable to initialize coinbasepro Client:", err)
		return
	}
	defer cbpClient.CloseUnderlyingConnection()
	defer poloClient.CloseUnderlyingConnection()

	fmt.Println("CBP")
	fmt.Println(cbpClient)
	fmt.Println("Polo")
	fmt.Println(poloClient)

	//Setup JSON CoinbasePro
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

	//Setup JSON Poloniex
	jsonPolo, err := ioutil.ReadFile("./testhelpers/testdata/poloniex/test-ticker-sub.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	jsonPolo2, err := ioutil.ReadFile("./testhelpers/testdata/poloniex/test-ticker-unsub.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	cbpClient.SetSubscribeMessage(jsonFile)
	cbpClient.SetUnsubscribeMessage(jsonFile2)
	poloClient.SetSubscribeMessage(jsonPolo)
	poloClient.SetUnsubscribeMessage(jsonPolo2)

	go cbpClient.StartStreaming(c1, interrupt)
	//go poloClient.StartStreaming(c2, interrupt)

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

			cbpBk.HandleMessage(message)
		case message := <-c2:
			chanSize := len(c2)
			msgReceived++
			if chanSize > maxSizeReached {
				maxSizeReached = chanSize
			}
			log.Info("POLO Channel size: ", chanSize)
			log.Info(string(message))
			log.Debug(message)

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

			//Test log
			log.Info((cbpBk.B.GetBooks(3).GetMaxBidPriceLevel()))
			log.Info((cbpBk.B.GetBooks(3).GetMinAskPriceLevel()))

			return
		}
	}
}
