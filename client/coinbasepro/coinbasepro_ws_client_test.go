package coinbasepro

import (
    "encoding/json"
    "testing"
    "io/ioutil"

    "../../bookkeeper/pricebook"
)

func containsVertices(destMap map[string]pricebook.Edge, keys []string, t *testing.T) {
    expectedKeySet := make(map[string]bool)
    for _,s := range keys {
        expectedKeySet[s] = true
    }
    for k,_ := range destMap {
        delete(expectedKeySet, k)
    }
    if len(expectedKeySet) != 0 {
        var remainingKeys string
        for k,_ := range expectedKeySet {
            remainingKeys += k 
            remainingKeys += " "
        }
        t.Errorf("Missing destination vertices: %s", remainingKeys)
    }
}

func getMsgFromFile(filepath string) ([]byte, error) {
    jsonFile, err := ioutil.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    return jsonFile, nil
}

//Remove dependency on test file later on
func TestParsePricePairs(t *testing.T) {
    jsonFile, err := getMsgFromFile("../../testdata/coinbasepro/test-l2-subscribe.json")
    if err != nil {
        t.Error("Could not read test file.")
    }

    var subscriptionMsg SubscriptionMessage

    err = json.Unmarshal(jsonFile, &subscriptionMsg)
    if err != nil {
        t.Error("Unmarshal error.")
    }

    pairMap := ParsePricePairs(subscriptionMsg)
    gotLen := len(pairMap)
    if gotLen != 3 {
        t.Errorf("len(*pairMap) = %d, want 3", gotLen)
    }

    for k,v := range pairMap {
        switch fromCoin := k; fromCoin {
        case "ETH":
            expectedKeys := []string{"USD", "BTC"}
            containsVertices(v, expectedKeys, t)
        case "BTC":
            expectedKeys := []string{"ETH"}
            containsVertices(v, expectedKeys, t)
        case "USD":
            expectedKeys := []string{"ETH"}
            containsVertices(v, expectedKeys, t)
        }
    }
}

func TestParseL2Message(t *testing.T) {
    jsonFile, err := getMsgFromFile("../../testdata/coinbasepro/test-l2-snapshot-response.json")
    if err != nil {
        t.Error("Could not read test file.")
    }

    var l2SnapshotMessage L2SnapshotMessage

    err = json.Unmarshal(jsonFile, &l2SnapshotMessage)
    if err != nil {
        t.Error("Unmarshal error.")
    }

    if l2SnapshotMessage.MessageType != "snapshot" {
        t.Errorf("Expected MessageType: %s, got %s", "snapshot", l2SnapshotMessage.MessageType)
    }
    if l2SnapshotMessage.ProductId != "ETH-USD" {
        t.Errorf("Expected ProductId: %s, got %s", "ETH-USD", l2SnapshotMessage.ProductId)
    }

    if l2SnapshotMessage.Asks[0][0] != "188.52" && l2SnapshotMessage.Asks[0][1] != "14.06613634" {
        t.Errorf("First Ask Price mismatched. Got %s, %s", l2SnapshotMessage.Asks[0][0], l2SnapshotMessage.Asks[0][1])
    }
    if l2SnapshotMessage.Bids[0][0] != "188.51" && l2SnapshotMessage.Bids[0][1] != "0.26060000" {
        t.Errorf("First Bid Price mismatched. Got %s, %s", l2SnapshotMessage.Bids[0][0], l2SnapshotMessage.Bids[0][1])
    }


}

func BenchmarkParseL2Message(b *testing.B) {
    jsonFile, err := getMsgFromFile("../../testdata/coinbasepro/test-l2-snapshot-response.json")
    if err != nil {
        b.Error("Could not read test file.")
    }
    b.ResetTimer()

    var l2SnapshotMessage L2SnapshotMessage

    for i := 0; i < b.N; i++ {
        err = json.Unmarshal(jsonFile, &l2SnapshotMessage)
        if err != nil {
            b.Error("Unmarshal error.")
        }
    }
}