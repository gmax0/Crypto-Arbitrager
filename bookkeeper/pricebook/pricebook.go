package pricebook

//Temporary Struct definition for price representations between coins
type Edge struct {
    High float32
    Low  float32
}

//TODO: Representation of various price points at map[string]map[string]interface{}
type PriceBook struct {
    Exchange string
    Graph    map[string]map[string]Edge
}

//
func NewPricebook(exchange string, pricePairs [][]string) {

}

func initGraph(pricePairs [][]string) {


}