package binance

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type StreamsManager struct {
	// List of all the underlying Websocket connections
	Connections []*websocket.Conn
	// A common consumer Channel where all the websocket connections messages are relayed,
	// this way we can consume the messages across all the connections in a single place
	MessageBroadcast chan ([]byte)
	// A common consumer Channel where all the errors from the server are relayed
	ErrorBroadcast chan (error)

	// internal state to send pong frames across all connections that we maintain
	keepAliveTimer *time.Ticker
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Minute
)

func NewStreamsManager() StreamsManager {
	return StreamsManager{
		MessageBroadcast: make(chan []byte),
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
		streams := streamsInLowerCase[i*limit : (i*limit)+limit]
		connection, err := OpenStream(streams)
		if err != nil {
			return err
		}
		connection.SetReadDeadline(time.Now().Add(pongWait))
		connection.SetPongHandler(func(string) error { connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		SM.Connections = append(SM.Connections, connection)

		go func() {
			for {
				_, message, err := connection.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
						SM.ErrorBroadcast <- err
						return
					}
					break
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				SM.MessageBroadcast <- message
			}
		}()
	}

	// background task to send pong frames every 5 minutes to all the connections that we maintain
	SM.keepAliveTimer = time.NewTicker(5 * time.Minute)
	go func() {
		for {
			<-SM.keepAliveTimer.C

			for _, connection := range SM.Connections {
				err := connection.WriteMessage(websocket.TextMessage, []byte(""))
				if err != nil {
					log.Print("ERROR: write:", err)
					SM.ErrorBroadcast <- err
					return
				}
			}
		}
	}()

	return nil
}

func (SM StreamsManager) Close() {
	if SM.keepAliveTimer != nil {
		SM.keepAliveTimer.Stop()
	}
}
