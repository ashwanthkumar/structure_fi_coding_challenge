package binance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllSymbols(t *testing.T) {
	data, err := os.ReadFile("./symbols_response.json")
	assert.NoError(t, err)

	dataAsJsonString := string(data)
	allSymbols, err := ParseSymbolsResponse(dataAsJsonString)
	assert.NoError(t, err)

	assert.Equal(t, 1961, len(allSymbols))
	assert.Contains(t, allSymbols, "WTCETH")
	assert.Contains(t, allSymbols, "BURGERBUSD")
}

func TestGetAllSymbolsWithInvalidInputShouldFail(t *testing.T) {
	dataAsJsonString := "invalid json input"
	_, err := ParseSymbolsResponse(dataAsJsonString)
	assert.Error(t, err)
}
