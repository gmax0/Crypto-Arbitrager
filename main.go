package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "os"
    "os/signal"

    _"./client"
    "./client/CoinbasePro"
    _"github.com/spf13/viper"
    _"./config"

)

func worker() {

}

func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    c1 := make(chan []byte)

    //Setup CoinbasePro Client Thread
    cbp_client, err := CoinbasePro.NewClient("ws-feed.pro.coinbase.com")
    if err != nil {
        log.Fatal("Unable to initialize CoinbasePro Client:", err)
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
        log.Fatal("CoinbasePro Client write error:", err)
        return
    }


    // go client.
    go cbp_client.StreamMessages(c1)

    for {
        select {
            case <- c1:
                message := <-c1
                message = []byte("TEST")
                log.Println(message) 
            case <-interrupt:
                log.Println("interrupt")
                err = cbp_client.CloseConnection()
                if (err != nil) {
                    log.Println("CoinbasePro Client write close error:", err)
                }
                return
        }
    }
}