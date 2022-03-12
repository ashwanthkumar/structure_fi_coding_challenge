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

	// Ref - https://go.dev/blog/pprof
	"github.com/gin-contrib/pprof"
	"github.com/rcrowley/go-metrics"

	"github.com/ashwanthkumar/structure_fi_coding_challenge/binance"
	docs "github.com/ashwanthkumar/structure_fi_coding_challenge/docs"
	"github.com/ashwanthkumar/structure_fi_coding_challenge/store"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var AppVersion = "0.0.1-dev"
var BuildTimestamp = "running locally"
var StartTime = time.Now()

func main() {
	metrics.RegisterDebugGCStats(metrics.DefaultRegistry)
	metrics.RegisterRuntimeMemStats(metrics.DefaultRegistry)
	go metrics.CaptureDebugGCStats(metrics.DefaultRegistry, 1*time.Minute)
	go metrics.CaptureRuntimeMemStats(metrics.DefaultRegistry, 1*time.Minute)
	go metrics.Log(metrics.DefaultRegistry, 30*time.Minute, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	streamsManager := binance.NewStreamsManager()

	allSymbols, err := binance.GetAllSymbols()
	if err != nil {
		log.Fatalf("%v", err)
	}

	datastore := store.NewStore(allSymbols)
	go startConsumingPriceStream(allSymbols, streamsManager, datastore)

	if !strings.Contains(AppVersion, "-dev") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	pprof.Register(router)

	docs.SwaggerInfo.Version = AppVersion
	docs.SwaggerInfo.Host = "localhost:8080" // TODO(ashwanthkumar): need to make this so it works on non-local deployments too
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := router.Group("/api/v1")
	{
		v1.GET("/z/info", AppInfo())
		v1.GET("/z/metrics", MetricsInfo())
		v1.GET("/symbols", ReturnAllSymbols(allSymbols, datastore))
		v1.GET("/:symbol", ReturnSymbolInfo(datastore))
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Starting Server at: %s", srv.Addr)
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
	log.Printf("Setting up trade stream connections")
	symbolTradeStreams := make([]string, 0)
	for _, symbol := range allSymbols {
		symbolTradeStreams = append(symbolTradeStreams, strings.ToLower(symbol)+"@trade")
	}

	err := streamsManager.Open(symbolTradeStreams)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
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
