package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ashwanthkumar/structure_fi_coding_challenge/binance"
	"github.com/ashwanthkumar/structure_fi_coding_challenge/store"
	"github.com/gin-gonic/gin"
)

var AppVersion = "0.0.1-dev"
var BuildTimestamp = time.Now().Format(time.RFC3339)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	streamsManager := binance.NewStreamsManager()

	allSymbols, err := binance.GetAllSymbols()
	if err != nil {
		log.Fatalf("%v", err)
	}

	datastore := store.NewStore(allSymbols)
	go startConsumingPriceStream(allSymbols, streamsManager, datastore)

	router := gin.Default()
	router.GET("/:symbol", func(c *gin.Context) {
		symbol := c.Param("symbol")
		stat, present := datastore.Get(symbol)
		if present {
			c.JSON(http.StatusOK, stat)
		} else {
			c.JSON(http.StatusNotFound, []int{})
		}
	})
	router.GET("/symbols", func(c *gin.Context) {
		response := make(map[string]interface{})
		response["all"] = allSymbols
		response["active"] = datastore.Symbols()
		c.JSON(http.StatusOK, response)
	})
	router.GET("/version", func(c *gin.Context) {
		response := make(map[string]interface{})
		response["version"] = AppVersion
		response["timestamp"] = BuildTimestamp
		c.JSON(http.StatusOK, response)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 2)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	streamsManager.Close()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func startConsumingPriceStream(allSymbols []string, streamsManager binance.StreamsManager, datastore store.Store) {
	symbolTradeStreams := make([]string, 0)
	for _, symbol := range allSymbols {
		symbolTradeStreams = append(symbolTradeStreams, strings.ToLower(symbol)+"@trade")
	}

	streamsManager.Open(symbolTradeStreams)
	for {
		select {
		case msg, ok := <-streamsManager.MessageBroadcast:
			if ok {
				// log.Printf("Message: %s\n", string(msg))
				datastore.Add(msg.Symbol, msg.Price)
			}
			// messages that we get
		case err, ok := <-streamsManager.ErrorBroadcast:
			if ok {
				// errors that we get while reading the data
				log.Fatalf("[ERROR]: %v\n", err)
			}
		}
	}
}
