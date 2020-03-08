package pricegraph

import (
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type GraphAdjMatrix struct {
	Exchange          int
	TickerKeys        []string //Indexes into the 2-D slice.
	Matrix            [][]float64
	TransformedMatrix [][]float64 //Master Graph will not have this set.
	m                 sync.Mutex
}

func NewGraphAdjMatrix(exchange int, tickers []string) *GraphAdjMatrix {
	g := &GraphAdjMatrix{Exchange: exchange}

	//Deep Copy Ticker
	t := make([]string, len(tickers))
	for i, ticker := range tickers {
		t[i] = ticker
	}

	//Initialize Zero-Valued Adjacency Matrix
	m := make([][]float64, len(tickers))

	for i := 0; i < len(tickers); i++ {
		m[i] = make([]float64, len(tickers))
	}

	g.TickerKeys = t
	g.Matrix = m

	return g
}

//CloneGraph performs a deep copy of the GraphAdjMatrix that the function was called on
func (g *GraphAdjMatrix) CloneGraph() *GraphAdjMatrix {
	cg := &GraphAdjMatrix{Exchange: g.Exchange}

	l := len(g.TickerKeys)

	//Deep Copy TickerKeys
	t := make([]string, l)
	for i, ticker := range g.TickerKeys {
		t[i] = ticker
	}

	//Deep Copy Adjacency Matrix
	m := make([][]float64, l)
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			m[i] = make([]float64, l)
			m[i][j] = g.Matrix[i][j]
		}
	}

	cg.TickerKeys = t
	cg.Matrix = m

	return cg
}

//TODO: Right/Left Justify grid contents
func (g *GraphAdjMatrix) PrintGraph() {
	var b strings.Builder
	for i := range g.Matrix {
		for j := range g.Matrix[i] {
			b.WriteString(strconv.FormatFloat(g.Matrix[i][j], 'f', -1, 64) + " ")
		}
		b.WriteString("\n")
	}
	logrus.Info("\n" + b.String())
}
