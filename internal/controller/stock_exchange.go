package controller

import (
	"Tradeasy/internal/services/stock_exchange"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteBuyOrder(c *gin.Context) {
	var buyReq stock_exchange.OrderRequest
	c.BindJSON(&buyReq)
	msg, err := stock_exchange.BuyOrder(buyReq)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, msg)
	}
}
func ExecuteSellOrder(c *gin.Context) {
	var sellReq stock_exchange.OrderRequest
	c.BindJSON(&sellReq)
	msg, err := stock_exchange.SellOrder(sellReq)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, msg)
	}
}
func DeleteBuyOrder(c *gin.Context) {
	id := c.Params.ByName("order_id")
	msg, err := stock_exchange.DeleteBuyOrder(id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, msg)
	}

}
func DeleteSellOrder(c *gin.Context) {
	id := c.Params.ByName("order_id")
	msg, err := stock_exchange.DeleteSellOrder(id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, msg)
	}

}
