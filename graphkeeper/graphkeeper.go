package graphkeeper

import (
	"./pricegraph"
)

/*  TODO:
 *		- Determine how to enable cross-exchange arbitrage with this model...may need to
 *		  refactor the design of Graphkeeper
 */

type Graphkeeper struct {
	Graphs map[int]pricegraph.GraphAdjMatrix //Key: constants.Exchange
}

func NewGraphkeeper() *Graphkeeper {
	return &Graphkeeper{Graphs: make(map[int]pricegraph.GraphAdjMatrix)}
}
