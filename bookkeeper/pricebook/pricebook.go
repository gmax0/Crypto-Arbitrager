package pricebook

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

//
func NewPricebook(exchange string, emptyPairGraph map[string]map[string]Edge) *PriceBook {
    p := &PriceBook{ Exchange: exchange, Graph: emptyPairGraph }
    return p
}

// ProcessPriceDump will process an exchange's initial level-2 order book price dump as a canoniclized message
func ProcessPriceDump() {
    
}
