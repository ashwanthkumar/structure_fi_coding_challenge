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
// @Success      200  {string}  main.SymbolsResponse
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
