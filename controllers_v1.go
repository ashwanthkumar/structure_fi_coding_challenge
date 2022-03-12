package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/ashwanthkumar/structure_fi_coding_challenge/store"
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
)

type SymbolsResponse struct {
	TotalOccurrences uint64   `json:"totalMessages" example:[15000]`
	Active           []string `json:"active" example:["ETHBTC"]`
	All              []string `json:"all" example:["NEOBTC"]"`
}

// Return all symbols
// @Summary  Return all symbols and the active symbols for which we have data in our datastore
// @Schemes
// @Description  Return all symbols and the active symbols for which we have data in our datastore
// @Tags         Solution
// @Produce      json
// @Success      200  {object}  main.SymbolsResponse
// @Router       /symbols [get]
func ReturnAllSymbols(allSymbols []string, datastore store.Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		activeSymbols, totalOccurrences := datastore.Symbols()
		response := SymbolsResponse{
			TotalOccurrences: totalOccurrences,
			Active:           activeSymbols,
			All:              allSymbols,
		}
		c.JSON(http.StatusOK, response)
	}
}

// Return symbol specific info
// @Summary      Return symbol specific info
// @Param        symbol  path  string  true  "Ticker Symbol"
// @Description  Return symbol specific info
// @Tags         Solution
// @Produce      json
// @Success      200  {object}  store.Stat
// @Router       /{symbol} [get]
func ReturnSymbolInfo(datastore store.Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		stat, present := datastore.Get(symbol)
		if present {
			c.JSON(http.StatusOK, stat)
		} else {
			c.JSON(http.StatusNotFound, []int{})
		}
	}
}

type AppInfoResponse struct {
	GitSha              string `json:"gitSha"`
	BuildTime           string `json:"buildTime"`
	StartTime           string `json:"startTime"`
	RunningTime         string `json:"runningTime"`
	HeapMemoryAllocated string `json:"heapMemoryUsage"`
	SysMemoryAllocated  string `json:"sysMemoryUsage"`
	MoreInfo            string `json:"moreInfo"`
}

// Return runtime service version info
// @Summary      Return runtime service version info
// @Description  Return runtime service version info
// @Tags         Ops
// @Produce      json
// @Success      200  {object}  main.AppInfoResponse
// @Router       /z/info [get]
func AppInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		rfc2822 := "Mon Jan 02 15:04:05 -0700 2006"
		runningTime := time.Since(StartTime)
		response := AppInfoResponse{
			GitSha:              AppVersion,
			BuildTime:           BuildTimestamp,
			StartTime:           StartTime.Format(rfc2822),
			RunningTime:         runningTime.String(),
			HeapMemoryAllocated: fmt.Sprintf("%d MiB", m.Alloc/1024/1024),
			SysMemoryAllocated:  fmt.Sprintf("%d MiB", m.Sys/1024/1024),
			MoreInfo:            "/api/v1/z/metrics",
		}
		c.JSON(http.StatusOK, response)
	}
}

// Return metrics that are captured
// @Summary      Return metrics that are captured
// @Description  Return metrics that are captured
// @Tags         Ops
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Router       /z/metrics [get]
func MetricsInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, metrics.DefaultRegistry.GetAll())
	}
}
