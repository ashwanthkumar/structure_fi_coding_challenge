package store

import (
	"os"
	"testing"

	"github.com/ashwanthkumar/structure_fi_coding_challenge/binance"
	"github.com/stretchr/testify/assert"
)

func TestNoCollisionAgainstAllSymbolsDataset(t *testing.T) {
	data, err := os.ReadFile("../binance/symbols_response.json")
	assert.NoError(t, err)

	dataAsJsonString := string(data)
	allSymbols, err := binance.ParseSymbolsResponse(dataAsJsonString)
	assert.NoError(t, err)

	customMap := NewMapWithPHF(allSymbols)
	for _, symbol := range allSymbols {
		customMap.Set(symbol, &Stat{Symbol: symbol})
	}

	// now that we've set all some value to all the keys, we should never have an empty slot in out Values
	// because if we do, we have had a collision that overwrote the values elsewhere and hence a particular
	// slot is free.
	for idx, stat := range customMap.Values {
		assert.NotNil(t, stat, "map at key %d is nil, we have a collision", idx)
	}
}

func TestMapGetSetIsWorkingForValidInputs(t *testing.T) {
	keys := []string{"A", "B", "C"}
	customMap := NewMapWithPHF(keys)
	for _, k := range keys {
		customMap.Set(k, &Stat{Symbol: k})
	}

	for _, k := range keys {
		stat, present := customMap.Get(k)
		assert.True(t, present)
		assert.Equal(t, &Stat{Symbol: k}, stat)
	}
}
