package main

import (
	"net/http"

	"github.com/ashwanthkumar/structure_fi_coding_challenge/store"
	"github.com/gin-gonic/gin"
)

type SymbolsResponse struct {
	All    []string `json:"all" example:["NEOBTC"]"`
	Active []string `json:"active" example:["ETHBTC"]`
}

// Return all symbols
// @Summary  Return all symbols and the active symbols for which we have data in our datastore
// @Schemes
// @Description  Return all symbols and the active symbols for which we have data in our datastore
// @Tags         symbols
// @Produce      json
// @Success      200  {object}  main.SymbolsResponse
// @Router       /symbols [get]
func ReturnAllSymbols(allSymbols []string, datastore store.Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		response := SymbolsResponse{
			All:    allSymbols,
			Active: datastore.Symbols(),
		}
		c.JSON(http.StatusOK, response)
	}
}

// Return symbol specific info
// @Summary      Return symbol specific info
// @Param        symbol  path  string  true  "Ticker Symbol"
// @Description  Return symbol specific info
// @Tags         symbols
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

type VersionInfoResponse struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
}

// Return service version info
// @Summary      Return service version info
// @Description  Return service version info
// @Tags         Ops
// @Produce      json
// @Success      200  {object}  main.VersionInfoResponse
// @Router       /version [get]
func VersionInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		response := VersionInfoResponse{
			Version:   AppVersion,
			BuildTime: BuildTimestamp,
		}
		c.JSON(http.StatusOK, response)
	}
}
