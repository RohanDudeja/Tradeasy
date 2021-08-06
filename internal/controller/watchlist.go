package controller

import (
	"Tradeasy/internal/services/watchlist"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateWatchlist(c *gin.Context) {
	var req watchlist.CreateRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res, err := watchlist.CreateWatchlist(req)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func AddStockEntry(c *gin.Context) {
	var req watchlist.AddStockRequest
	watchlistId := c.Params.ByName("watchlist_id")
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	i, _ := strconv.ParseInt(watchlistId, 10, 64)
	res, er := watchlist.AddStockEntry(req, int(i))
	if er != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func DeleteStockEntry(c *gin.Context) {
	var req watchlist.DeleteStockRequest
	watchlistId := c.Params.ByName("watchlist_id")
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	i, _ := strconv.ParseInt(watchlistId, 10, 64)
	res, er := watchlist.DeleteStockEntry(req, int(i))
	if er != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func SortWatchlist(c *gin.Context) {
	var req watchlist.SortRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	res, er := watchlist.SortWatchlist(req)
	if er != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}
