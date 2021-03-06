package binance

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rcrowley/go-metrics"
	"github.com/recws-org/recws"
	"github.com/valyala/fastjson"
)

type StreamsManager struct {
	// List of all the underlying Websocket connections
	Connections []recws.RecConn
	// A common consumer Channel where all the websocket connections messages are relayed,
	// this way we can consume the messages across all the connections in a single place
	MessageBroadcast chan (StreamMessage)
	// A common consumer Channel where all the errors from the server are relayed
	ErrorBroadcast chan (error)

	// internal state to send pong frames across all connections that we maintain
	keepAliveTimer *time.Ticker
}

type StreamMessage struct {
	Symbol string
	Price  float64
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewStreamsManager() StreamsManager {
	return StreamsManager{
		MessageBroadcast: make(chan StreamMessage, 500),
		ErrorBroadcast:   make(chan error),
	}
}

func (SM StreamsManager) Open(streamsInLowerCase []string) error {
	// Based on the documentation described at https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams.md#websocket-limits
	// it looks like we can only have1024 streams in a single connection, irrespective of how we open streams (using /streams/<streamname>) or
	// manually (read slowly) subscribe to multiple streams using SUBSCRIBE messages, the API seems to block us at 1024 streams.
	limit := 1024
	nrOfConnections := len(streamsInLowerCase) / limit
	for i := 0; i <= nrOfConnections; i++ {
		maxLimit := (i * limit) + limit
		if i == nrOfConnections {
			maxLimit = len(streamsInLowerCase) - 1
		}
		streams := streamsInLowerCase[i*limit : maxLimit]
		ws := OpenStream(streams)
		SM.Connections = append(SM.Connections, ws)

		go func() {
			// re-using the json parser per go-routine that consumes from the websocket connection
			var jsonParser fastjson.Parser

			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
						SM.ErrorBroadcast <- err
						return
					}
					break
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				streamMessage, err := parseMessageFromTradeStream(jsonParser, string(message))
				if err != nil {
					SM.ErrorBroadcast <- err
				} else {
					SM.MessageBroadcast <- *streamMessage
					metrics.GetOrRegisterMeter("consumption-rate", nil).Mark(1)
				}
			}
		}()
	}

	// background task to send pong frames every 5 minutes to all the connections that we maintain
	SM.keepAliveTimer = time.NewTicker(30 * time.Second)
	go func() {
		for {
			<-SM.keepAliveTimer.C

			for _, connection := range SM.Connections {
				err := connection.WriteMessage(websocket.PongMessage, []byte(""))
				if err != nil {
					log.Print("ERROR: write:", err)
					SM.ErrorBroadcast <- err
					return
				}
			}
		}
	}()

	log.Printf("Websocket streams are setup and we're consuming: %d trade streams across %d websocket connections", len(streamsInLowerCase), len(SM.Connections))
	return nil
}

func (SM StreamsManager) Close() {
	if SM.keepAliveTimer != nil {
		SM.keepAliveTimer.Stop()
	}
}

func parseMessageFromTradeStream(jsonParser fastjson.Parser, message string) (*StreamMessage, error) {
	v, err := jsonParser.Parse(message)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(string(v.GetStringBytes("data", "p")), 64)
	if err != nil {
		return nil, err
	}

	symbol := string(v.GetStringBytes("data", "s"))

	streamMessage := &StreamMessage{
		Symbol: symbol,
		Price:  price,
	}
	return streamMessage, nil
}
