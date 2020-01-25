package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "os"
    "os/signal"

    "./client/CoinbasePro"
    _"github.com/spf13/viper"
    _"./config"

)


func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)


    client, err := CoinbasePro.NewClient("ws-feed.pro.coinbase.com")
    if err != nil {
        log.Fatal("Unable to initialize CoinbasePro Client:", err)
        return
    }
    defer client.CloseUnderlyingConnection()

    fmt.Println(client)

    //Setup JSON Message
    jsonFile, err := ioutil.ReadFile("./config/test-subscribe.json")
    if err != nil {
        fmt.Println(err)
        return
    }

    err = client.Subscribe(jsonFile)
    if err != nil {
        log.Fatal("CoinbasePro Client write error:", err)
        return
    }

    go client.StreamMessages()

    for {
        select {
            case <-interrupt:
                log.Println("interrupt")
                err = client.CloseConnection()
                if (err != nil) {
                    log.Println("CoinbasePro Client write close error:", err)
                }
                return
        }
    }
}