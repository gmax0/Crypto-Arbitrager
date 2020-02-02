package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	coinbasepro "./client/coinbasepro"
	_ "./config"
	_ "github.com/spf13/viper"
)

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
	jsonFile, err := ioutil.ReadFile("./config/test-subscribe.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cbp_client.Subscribe(jsonFile)
	if err != nil {
		log.Fatal("coinbasepro Client write error:", err)
		return
	}

	// go client.
	go cbp_client.StreamMessages(c1)

	for {
		select {
		case <-c1:
			message := <-c1
			message = []byte("TEST")
			chanSize := len(c1)
			log.Println(chanSize)
			log.Println(message)
		case <-interrupt:
			log.Println("interrupt")
			err = cbp_client.CloseConnection()
			if err != nil {
				log.Println("coinbasepro Client write close error:", err)
			}
			return
		}
	}
}
