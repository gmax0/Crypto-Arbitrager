package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
	"io/ioutil"
	"fmt"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "ws-feed.pro.coinbase.com", "Coinbase Pro WS Address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	u := url.URL{Scheme: "wss", Host: *addr, Path: ""}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	jsonFile, err := ioutil.ReadFile("./config/test-subscribe.json")
	if err != nil {
		fmt.Println(err)
	}

	err = c.WriteMessage(websocket.TextMessage, jsonFile)
	if err != nil {
		log.Println("write:", err)
		return
	}


	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			log.Println("DONE")
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}

	}


}
