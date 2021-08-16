package controller

import (
	"Tradeasy/internal/services/watchlist"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func CreateWatchlist(c *gin.Context) {
	var req watchlist.CreateRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error in the Given Order Request Body"})
	}
	res, err := watchlist.CreateWatchlist(req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func AddStockEntry(c *gin.Context) {
	var req watchlist.AddStockRequest
	watchlistId := c.Params.ByName("watchlist_id")
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error in the Given Order Request Body"})
	}
	i, _ := strconv.ParseInt(watchlistId, 10, 64)
	res, err := watchlist.AddStockEntry(req, int(i))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func DeleteStockEntry(c *gin.Context) {
	var req watchlist.DeleteStockRequest
	watchlistId := c.Params.ByName("watchlist_id")
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Error in the Given Order Request Body"})
	}
	i, _ := strconv.ParseInt(watchlistId, 10, 64)
	res, err := watchlist.DeleteStockEntry(req, int(i))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
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
	res, err := watchlist.SortWatchlist(req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
