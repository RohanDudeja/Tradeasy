package controller

import (
	"Tradeasy/internal/services/stock_exchange"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ExecuteBuyOrder(c *gin.Context) {
	var buyReq stock_exchange.OrderRequest

	if err := c.BindJSON(&buyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}
	msg, err := stock_exchange.BuyOrder(buyReq)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, msg)
	}
}
func ExecuteSellOrder(c *gin.Context) {
	var sellReq stock_exchange.OrderRequest

	if err := c.BindJSON(&sellReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}
	msg, err := stock_exchange.SellOrder(sellReq)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, msg)
	}
}
func DeleteBuyOrder(c *gin.Context) {
	id := c.Params.ByName("order_id")
	msg, err := stock_exchange.DeleteBuyOrder(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, msg)
	}

}
func DeleteSellOrder(c *gin.Context) {
	id := c.Params.ByName("order_id")
	msg, err := stock_exchange.DeleteSellOrder(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, msg)
	}
}

func ViewMarketDepth(c *gin.Context) {
	stockName := c.Params.ByName("stock_name")
	msg, err := stock_exchange.ViewMarketDepth(stockName)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	} else {
		c.JSON(http.StatusOK, msg)
	}
}
