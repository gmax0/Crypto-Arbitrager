package pricebook

import (
    "strings"
    _"../../client/coinbasepro"
)

func init() {

}

const (
    CoinbasePro = "cbp"
    Poloniex    = "plnx"
)

//Temporary Struct definition for price representations between coins
type Edge struct {
    High float32
    Low  float32
}

// Generalize maker + taker fee structures later on
// 
type Fees struct {
    Maker float32
    Taker float32
}

//TODO: Representation of various price points at map[string]map[string]interface{}
type PriceBook struct {
    Exchange string
    Graph    map[string]map[string]Edge
}

/*******************************************************************************/

//
func NewPricebook(exchange string, pricePairs []string) *PriceBook {

    graph := ParsePricePairs(pricePairs)
    p := &PriceBook{ Exchange: exchange, Graph: graph }
    return p

}

//
func ParsePricePairs(pricePairs []string) map[string]map[string]Edge {

    var m map[string]map[string]Edge
    m = make(map[string]map[string]Edge)

    for _, pair := range pricePairs {
        split := strings.Split(pair, "-")
        if m[split[0]] == nil {
            m[split[0]] = make(map[string]Edge)
        }
        m[split[0]][split[1]] = Edge{ High: -1, Low: -1}

        if m[split[1]] == nil {
            m[split[1]] = make(map[string]Edge)
        }
        m[split[1]][split[0]] = Edge{ High: -2, Low: -2}
    }

    return m  
    
}

/*
func (p *PriceBook) ProcessPriceDump(exchangeFlg int, msg []byte) error {
    switch(exchangeFlg) {
    case 1:
        err := processCoinbaseProDump(msg)
        if err != nil {
            return err
        }
    default:
    }
    return nil
}
*/

/*
//
// TODO: Optimize JSON unmarshalling if it becomes a bottle neck here.
func processCoinbaseProDump(p *Pricebook, msg []byte) error {
    var snapshot coinbasepro.L2SnapshotMessage

    err = json.Unmarshal(msg, &snapshot)
    if err != nil {
        return err
    }


}
*/