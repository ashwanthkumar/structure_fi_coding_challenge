package binance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestParseStreamMessage(t *testing.T) {
	var p fastjson.Parser
	inputJson := `{"stream":"runeusdt@trade","data":{"e":"trade","E":1646952851073,"s":"RUNEUSDT","t":45289958,"p":"5.30400000","q":"2.00000000","b":610097169,"a":610096915,"T":1646952851072,"m":false,"M":true}}`
	streamMessage, err := parseMessageFromTradeStream(p, inputJson)

	assert.NoError(t, err)
	assert.Equal(t, "RUNEUSDT", streamMessage.Symbol)
	assert.Equal(t, 5.30400000, streamMessage.Price)
}
