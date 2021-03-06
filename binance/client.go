package binance

import (
	"net/url"
	"strings"

	"github.com/recws-org/recws"
	"github.com/valyala/fastjson"
)

func GetAllSymbols() ([]string, error) {
	url := "https://api.binance.com/api/v3/exchangeInfo"

	responseJsonAsString, err := http_Get(url)
	if err != nil {
		return nil, err
	}

	return ParseSymbolsResponse(responseJsonAsString)
}

func ParseSymbolsResponse(responseJsonAsString string) ([]string, error) {
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

func OpenStream(streams []string) recws.RecConn {
	addr := "stream.binance.com:9443"
	path := "/stream"
	u := url.URL{Scheme: "wss", Host: addr, Path: path}
	if len(streams) > 0 {
		query := "streams=" + strings.Join(streams, "/")
		u.RawQuery = query
	}

	// log.Printf("Opening websocket connection to %s", u.String())
	ws := recws.RecConn{}
	ws.Dial(u.String(), nil)
	return ws
}
