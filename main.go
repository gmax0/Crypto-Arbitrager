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
    "github.com/buger/jsonparser"
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
    jsonFile2, err := ioutil.ReadFile("./testdata/coinbasepro/test-l2-unsubscribe.json")
    if err != nil {
        log.Fatal(err)
        return
    }

    //
    cbp_pricePairs := []string{"ETH-BTC", "ETH-USD"}

    //Initialize CoinbasePro Pricebook
    cbp_pb := pricebook.NewPricebook(pricebook.CoinbasePro, cbp_pricePairs)
    log.Info(cbp_pb)


    /*
	err = cbp_client.Subscribe(jsonFile)
	if err != nil {
		log.Fatal("coinbasepro Client write error:", err)
		return
	}
    */
    cbp_client.SetSubscribeMessage(jsonFile)
    cbp_client.SetUnsubscribeMessage(jsonFile2)

	// go client.
	go cbp_client.StartStreaming(c1, interrupt)

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
                log.Error("Could not get message type")
            }
            log.Info(msgType)

            if msgType == "snapshot" {
                err = cbp_pb.ProcessPriceDump(message)
                if err != nil {
                    log.Error("Error processing snapshot message for CoinbasePro")
                }
            }
		case <-interrupt:
			log.Println("interrupt")
            err = cbp_client.StopStreaming()
            if err != nil {
                log.Error("StopStreaming() error: ", err)
            }

			err = cbp_client.CloseConnection()
			if err != nil {
				log.Error("coinbasepro Client write close error:", err)
			}

            log.Info("Max Chan Size Reached: " , maxSizeReached)
            log.Info("Messages Received: ", msgReceived)
            log.Info(cbp_pb)
			return
		}
	}
}
