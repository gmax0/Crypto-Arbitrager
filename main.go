package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	coinbasepro "./client/coinbasepro"
    pricebook "./bookkeeper/pricebook"
	_ "./config"

    "github.com/sirupsen/logrus"
	_ "github.com/spf13/viper"
    json "github.com/buger/jsonparser"
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

	//Setup coinbasepro Client Thread
	cbp_client, err := coinbasepro.NewClient("ws-feed.pro.coinbase.com")
	if err != nil {
		log.Fatal("Unable to initialize coinbasepro Client:", err)
		return
	}
	defer cbp_client.CloseUnderlyingConnection()

	fmt.Println(cbp_client)

	//Setup JSON Message
    jsonFile, err := ioutil.ReadFile("./testdata/coinbasepro/test-l2-subscribe.json")
	if err != nil {
		log.Fatal(err)
		return
	}

    //
    cbp_pricePairs := []string{"ETH-BTC", "ETH-USD"}

    //Initialize CoinbasePro Pricebook
    cbp_pb := pricebook.NewPricebook(pricebook.CoinbasePro, cbp_pricePairs)
    log.Info(cbp_pb)


	err = cbp_client.Subscribe(jsonFile)
	if err != nil {
		log.Fatal("coinbasepro Client write error:", err)
		return
	}

	// go client.
	go cbp_client.StreamMessages(c1)

    maxSizeReached := 0
	for {
		select {
		case <-c1:
			message := <-c1
			chanSize := len(c1)
            if chanSize > maxSizeReached {
                maxSizeReached = chanSize
            }
			log.Info("Channel size: ", chanSize)
            log.Trace(message)
            msgType, err := json.GetString(message, "type")
            if err != nil {
                log.Error("Could not get message type")
            }
            log.Info(msgType)

            
		case <-interrupt:
			log.Println("interrupt")
			err = cbp_client.CloseConnection()
			if err != nil {
				log.Error("coinbasepro Client write close error:", err)
			}
            log.Info("Max Chan Size Reached: " , maxSizeReached)
			return
		}
	}
}
