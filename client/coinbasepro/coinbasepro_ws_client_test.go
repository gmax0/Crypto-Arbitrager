package coinbasepro

import (
    "encoding/json"
    "testing"
    "io/ioutil"
)

//Remove dependency on test file later on
func TestParsePricePairs(t *testing.T) {
    jsonFile, err := ioutil.ReadFile("../../testdata/coinbasepro/test-subscribe.json")
    if err != nil {
        t.Error("Could not read test file.")
    }
    // t.Log(jsonFile)

    var subscriptionMsg SubscriptionMessage

    err = json.Unmarshal(jsonFile, &subscriptionMsg)
    if err != nil {
        t.Error("Unmarshal error.")
    }

    t.Log(subscriptionMsg)

    pairMap := ParsePricePairs(subscriptionMsg)
    t.Log(len(*pairMap))
    /*
    for k,v := range *pairMap {
        t.Log(k)
        for k2, v2 := range v {
            t.Log(k2)
            t.Log(v2)
        }
    }
    */
    for k,v := range *pairMap {
        t.Log(k)
        t.Log(v)
    }
}