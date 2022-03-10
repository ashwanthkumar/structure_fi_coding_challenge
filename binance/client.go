package binance

import (
	"github.com/valyala/fastjson"
)

func GetAllSymbols() ([]string, error) {
	url := "https://api.binance.com/api/v3/exchangeInfo"

	responseJsonAsString, err := Get(url)
	if err != nil {
		return nil, err
	}

	return parseSymbolsResponse(responseJsonAsString)
}

func parseSymbolsResponse(responseJsonAsString string) ([]string, error) {
	var p fastjson.Parser
	v, err := p.Parse(responseJsonAsString)
	if err != nil {
		return nil, err
	}

	symbols := v.GetArray("symbols")
	allSymbols := make([]string, len(symbols))
	for index, symbol := range symbols {
		s := string(symbol.GetStringBytes("symbol"))
		allSymbols[index] = s
	}

	return allSymbols, nil
}
