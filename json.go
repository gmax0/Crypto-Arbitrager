package main

import (
	"fmt"
	"log"
	_"os"
	_"io"
	"io/ioutil"
)

func main() {
	// jsonFile, err := os.Open("./config/test-subscribe.json")
	jsonFile, err := ioutil.ReadFile("./config/test-subscribe.json")
	if err != nil {
		fmt.Println(err)
	}

	log.Print(string(jsonFile))
	//defer jsonFile.Close()
}