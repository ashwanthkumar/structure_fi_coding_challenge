package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/ashwanthkumar/structure_fi_coding_challenge/binance"
	"github.com/ashwanthkumar/structure_fi_coding_challenge/store"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	table := store.NewStore()

	allSymbols, err := binance.GetAllSymbols()
	if err != nil {
		log.Fatalf("%v", err)
	}

	symbolTradeStreams := make([]string, 0)
	for _, symbol := range allSymbols {
		symbolTradeStreams = append(symbolTradeStreams, strings.ToLower(symbol)+"@trade")
	}

	streamsManager := binance.NewStreamsManager()
	streamsManager.Open(symbolTradeStreams)
	for {
		select {
		case msg, ok := <-streamsManager.MessageBroadcast:
			if ok {
				// log.Printf("Message: %s\n", string(msg))
				table.Add(msg.Symbol, msg.Price)
			}
			// messages that we get
		case err, ok := <-streamsManager.ErrorBroadcast:
			if ok {
				// errors that we get while reading the data
				log.Fatalf("[ERROR]: %v\n", err)
			}
		case <-interrupt:
			log.Print("Shutting down all the connections")
			streamsManager.Close()
			log.Print("Good Bye!")
			return
		}
	}
}
