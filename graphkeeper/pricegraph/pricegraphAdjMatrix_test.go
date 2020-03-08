package pricegraph

import (
	"testing"

	"../../common/constants"
)

func TestNewGraphAdjMatrix(t *testing.T) {
	tickers := []string{"ETH", "USD", "BTC", "XTZ"}
	g := NewGraphAdjMatrix(constants.CoinbasePro, tickers)

	if &(g.TickerKeys[0]) == &tickers[0] {
		t.Errorf("TickerKeys and tickers map to the same memory address")
	}
	for i, key := range g.TickerKeys {
		if tickers[i] != key {
			t.Errorf("Expected ticker value: %s, got: %s", tickers[i], key)
		}
	}

	for i := 0; i < len(tickers); i++ {
		for j := 0; j < len(tickers); j++ {
			if g.Matrix[i][j] != 0 {
				t.Errorf("Initial graph position %d,%d does not have a 0 value.", i, j)
			}
		}
	}
}

func TestCloneGraph(t *testing.T) {
	tickers := []string{"ETH", "USD", "BTC", "XTZ"}
	g1 := NewGraphAdjMatrix(constants.CoinbasePro, tickers)
	g2 := g1.CloneGraph()

	if &(g1.TickerKeys[0]) == &(g2.TickerKeys[0]) {
		t.Errorf("g1.TickerKeys and g2.TickerKeys map to the same memory address")
	}
	for i, key := range g1.TickerKeys {
		if g2.TickerKeys[i] != key {
			t.Errorf("Expected ticker value: %s, got: %s", key, g2.TickerKeys[i])
		}
	}

	if &(g1.Matrix[0][0]) == &(g2.Matrix[0][0]) {
		t.Errorf("g1.Matrix and g2.Matrix map to the same memory address")
	}
	for i := 0; i < len(tickers); i++ {
		for j := 0; j < len(tickers); j++ {
			if g1.Matrix[i][j] != g2.Matrix[i][j] {
				t.Errorf("Mismatched graph values at %d, %d", i, j)
			}
		}
	}
}
