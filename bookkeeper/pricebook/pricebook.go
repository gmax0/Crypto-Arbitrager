package pricebook

import (
    "strings"
    "strconv"
    _"../../client/coinbasepro"

    "github.com/buger/jsonparser"
    "github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {

}

const (
    CoinbasePro = 1
    Poloniex    = 2
)

// Generalize maker + taker fee structures later on
// 
type Fees struct {
    Maker float32
    Taker float32
}

//TODO: Representation of various price points at map[string]map[string]interface{}
type PriceBook struct {
    Exchange int
    Graph    map[string]map[string]float64
}

/*******************************************************************************/

//
func NewPricebook(exchange int, pricePairs []string) *PriceBook {
    graph := ParsePricePairs(pricePairs)
    p := &PriceBook{ Exchange: exchange, Graph: graph }
    return p
}

//
func ParsePricePairs(pricePairs []string) map[string]map[string]float64 {
    var m map[string]map[string]float64
    m = make(map[string]map[string]float64)

    for _, pair := range pricePairs {
        split := strings.Split(pair, "-")
        if m[split[0]] == nil {
            m[split[0]] = make(map[string]float64)
        }
        m[split[0]][split[1]] = 0

        if m[split[1]] == nil {
            m[split[1]] = make(map[string]float64)
        }
        m[split[1]][split[0]] = 0
    }

    return m  
}

//
func (p *PriceBook) ProcessPriceDump(msg []byte) error {
    log.Info("CALLED")
    switch(p.Exchange) {
    case 1:
        err := processCoinbaseProDump(p, msg)
        if err != nil {
            return err
        }
    default:
    }
    return nil
}

//
func processCoinbaseProDump(p *PriceBook, msg []byte) error {
    productId, err := jsonparser.GetString(msg, "product_id")
    log.Info(productId)
    if err != nil {
        log.Info(err)
    }

    splitProductId := strings.Split(productId, "-")
    src := splitProductId[0]
    dest := splitProductId[1]

    ask, err := jsonparser.GetString(msg, "asks", "[0]", "[0]")
    if err != nil {
        log.Info(err)
    }
    bid, err := jsonparser.GetString(msg, "bids", "[0]", "[0]")
    if err != nil {
        log.Info(err)
    }

    p.Graph[src][dest], _ = strconv.ParseFloat(bid, 64)
    p.Graph[dest][src], _ = strconv.ParseFloat(ask, 64)

    return nil
}