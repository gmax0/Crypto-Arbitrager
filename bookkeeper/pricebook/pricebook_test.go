package pricebook

import (
    "testing"
)

func containsVertices(destMap map[string]Edge, keys []string, t *testing.T) {
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

func TestParsePricePairs(t *testing.T) {
    pricePairs := []string{"ETH-BTC", "ETH-USD"}

    pairMap := ParsePricePairs(pricePairs)
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

